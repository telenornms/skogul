/*
 * skogul, sql sender
 *
 * Copyright (c) 2019 Telenor Norge AS
 * Author(s):
 *  - Kristian Lyngstøl <kly@kly.no>
 *
 * This library is free software; you can redistribute it and/or
 * modify it under the terms of the GNU Lesser General Public
 * License as published by the Free Software Foundation; either
 * version 2.1 of the License, or (at your option) any later version.
 *
 * This library is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
 * Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public
 * License along with this library; if not, write to the Free Software
 * Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA
 * 02110-1301  USA
 */

package sender

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	_ "github.com/go-sql-driver/mysql" // Imported for side effect/mysql support
	_ "github.com/lib/pq"
	"github.com/telenornms/skogul"
)

const (
	timestamp   = iota
	metadata    = iota
	data        = iota
	marshalData = iota
	marshalMeta = iota
)

var sqlLog = skogul.Logger("sender", "sql")

type dbElement struct {
	family int
	key    string
}

/*
SQL sender connects to a SQL Database, currently either MySQL(or Mariadb I
suppose) or Postgres. The Connection String for MySQL is specified at
https://github.com/go-sql-driver/mysql/ and postgres at
http://www.postgresql.org/docs/current/static/libpq-connect.html#LIBPQ-CONNSTRING
.

The query is expanded using os.Expand() and will fill in
timestamp/metadata/data. The sender will prep the query and
essentially covert INSERT INTO foo
VLAUES(${timestamp},${metadata.foo},${someData})
to foo("INSERT INTO foo VALUES(?,?,?)", timestamp, foo, someData), so they
will be sensibly escaped.
*/
type SQL struct {
	ConnStr string `doc:"Connection string to use for database. Slight variations between database engines. For MySQL typically user:password@host/database." example:"mysql: 'root:lol@/mydb' postgres: 'user=pqgotest dbname=pqgotest sslmode=verify-full'"`
	Query   string `doc:"Query run for each metric. The following expansions are made:\n\n${timestamp} is expanded to the actual metric timestamp.\n\n${metadata.KEY} will be expanded to the metadata with key name \"KEY\".\n\n${data.KEY} will be expanded to data[foo].\n\n${json.metadata} will be expanded to a json representation of all metadata.\n\n${json.data} will be expanded to a json representation of all data.\n\nFinally, ${KEY} is a shorthand for ${data.KEY}. Both methods are provided, to allow referencing data fields named \"metadata.\". E.g.: ${data.metadata.x} will match data[\"metadata.x\"], while ${metadata.x} will match metadata[\"x\"]." example:"INSERT INTO test VALUES(${timestamp},${hei},${metadata.key1})"`
	Driver  string `doc:"Database driver/system. Currently suported: mysql and postgres."`
	initErr error
	q       string
	list    []dbElement
	db      *sql.DB
	once    sync.Once
}

/*
prep parses sq.Query into q and populates sq.list accordingly.

It works by using os.Expand, with a custom function. For each element, we
always return ?/$x, thus building VALUES(?,?,?....), but it also records the
type of database element and the key in sq.list, so upon execution, we can
iterate over sq.list quickly to build the argument list.

Conveniently, the mysql driver doesn't understand $1, $2, $3 and the postgres
driver doesn't understand ?, ?, ?.

I love humans.
*/
func (sq *SQL) prep() {
	mlen := len("metadata.")
	dlen := len("data.")

	nElement := 0
	expander := func(element string) string {
		if element == "timestamp" {
			sq.list = append(sq.list, dbElement{timestamp, element})
		} else if len(element) > mlen && element[0:mlen] == "metadata." {
			sq.list = append(sq.list, dbElement{metadata, element[mlen:]})
		} else if element == "json.metadata" {
			sq.list = append(sq.list, dbElement{marshalMeta, ""})
		} else if element == "json.data" {
			sq.list = append(sq.list, dbElement{marshalData, ""})
		} else if len(element) > dlen && element[0:dlen] == "data." {
			sq.list = append(sq.list, dbElement{data, element[dlen:]})
		} else {
			sq.list = append(sq.list, dbElement{data, element})
		}
		if sq.Driver == "mysql" {
			return "?"
		}
		nElement++
		return fmt.Sprintf("$%d", nElement)
	}
	sq.q = os.Expand(sq.Query, expander)
}

func (sq *SQL) init() {
	sq.db, sq.initErr = sql.Open(sq.Driver, sq.ConnStr)
	if sq.initErr != nil {
		sqlLog.WithError(sq.initErr).WithField("driver", sq.Driver).Error("Failed to initialize SQL connection")
		return
	}
	sq.prep()
}

func (sq *SQL) exec(stmt *sql.Stmt, m *skogul.Metric) error {
	var vals []interface{}
	for _, e := range sq.list {
		switch e.family {
		case timestamp:
			vals = append(vals, m.Time)
		case metadata:
			vals = append(vals, m.Metadata[e.key])
		case data:
			vals = append(vals, m.Data[e.key])
		case marshalMeta:
			meta, err := json.Marshal(m.Metadata)
			if err != nil {
				return fmt.Errorf("unable to marshal metadata into json: %w", err)
			}
			vals = append(vals, meta)
		case marshalData:
			data, err := json.Marshal(m.Data)
			if err != nil {
				return fmt.Errorf("unable to marshal data into json: %w", err)
			}
			vals = append(vals, data)
		}
	}
	_, err := stmt.Exec(vals...)
	return err
}

/*
Send will send to the database, after first ensuring
the connection is OK.
*/
func (sq *SQL) Send(c *skogul.Container) error {
	sq.once.Do(func() {
		sq.init()
	})
	if sq.initErr != nil {
		return fmt.Errorf("database initialization failed: %w", sq.initErr)
	}
	txn, err := sq.db.Begin()
	if err != nil {
		return fmt.Errorf("beginning database transaction failed: %w", err)
	}
	defer func() {
		if err != nil {
			txn.Rollback()
		}
	}()

	stmt, err := txn.Prepare(sq.q)
	if err != nil {
		return err
	}

	for _, m := range c.Metrics {
		if err = sq.exec(stmt, m); err != nil {
			return err
		}
	}

	if err = stmt.Close(); err != nil {
		return err
	}

	if err = txn.Commit(); err != nil {
		return err
	}
	return nil
}

// Verify ensures options are set, but currently doesn't check very well,
// since it is disallowed from connecting to a database and such.
func (sq *SQL) Verify() error {
	if sq.ConnStr == "" {
		return skogul.MissingArgument("ConnStr")
	}
	if sq.Query == "" {
		return skogul.MissingArgument("Query")
	}
	if sq.Driver == "" {
		return skogul.MissingArgument("Driver")
	}
	if sq.Driver != "mysql" && sq.Driver != "postgres" {
		return fmt.Errorf("unsuported database driver %s - must be `mysql' or `postgres'", sq.Driver)
	}
	return nil
}
