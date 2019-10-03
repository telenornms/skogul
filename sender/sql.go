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
	"log"
	"os"
	"sync"

	"github.com/KristianLyng/skogul"
	_ "github.com/go-sql-driver/mysql" // Imported for side effect/mysql support
	_ "github.com/lib/pq"
)

const (
	timestamp   = iota
	metadata    = iota
	data        = iota
	marshalData = iota
	marshalMeta = iota
)

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
VLAUES(${timestamp.timestamp},${metadata.foo},${someData})
to foo("INSERT INTO foo VALUES(?,?,?)", timestamp, foo, someData), so they
will be sensibly escaped.
*/
type SQL struct {
	ConnStr string `doc:"Connection string to use for database. Slight variations between database engines. For MySQL typically user:password@host/database." example:"mysql: 'root:lol@/mydb' postgres: 'user=pqgotest dbname=pqgotest sslmode=verify-full'"`
	Query   string `doc:"Query run for each metric. ${timestamp.timestamp} is expanded to the actual metric timestamp. ${metadata.KEY} will be expanded to the metadata with key name \"KEY\", other ${foo} will be expanded to data[foo]. \n\nIn addition, ${json.data} and ${json.metadata} will be expanded to the json-encoded representation of the data and metadata respectively.\n\nNote that this is sensibly escaped, so while it might seem like it is vulnerable to SQL injection, it should be safe." example:"INSERT INTO test VALUES(${timestamp.timestamp},${hei},${metadata.key1})"`
	Driver  string `doc:"Database driver/system. Currently suported: mysql and postgres."`
	q       string
	list    []dbElement
	db      *sql.DB
	once    sync.Once
}

/*
prep parses sq.Query into q and populates sq.list accordingly
*/
func (sq *SQL) prep() {
	mlen := len("metadata.")

	expander := func(element string) string {
		if element == "timestamp.timestamp" {
			sq.list = append(sq.list, dbElement{timestamp, element})
		} else if len(element) > mlen && element[0:mlen] == "metadata." {
			sq.list = append(sq.list, dbElement{metadata, element[mlen:]})
		} else if element == "json.metadata" {
			sq.list = append(sq.list, dbElement{marshalMeta, ""})
		} else if element == "json.data" {
			sq.list = append(sq.list, dbElement{marshalData, ""})
		} else {
			sq.list = append(sq.list, dbElement{data, element})
		}
		return "?"
	}
	sq.q = os.Expand(sq.Query, expander)
}

// GetQuery returns the parsed query, assuming there is one.
func (sq *SQL) GetQuery() (string, error) {
	if sq.Query == "" {
		return "", skogul.Error{Source: "sql sender", Reason: "No Query set, but GetQuery() called"}
	}
	err := sq.Init()
	if err != nil {
		return "", skogul.Error{Source: "sql sender", Reason: "SQL.Init failed during GetQuery()", Next: err}
	}
	return sq.q, nil
}

/*
Init will connect to the database, ping it and set things up. But only once.
*/
func (sq *SQL) Init() error {
	var er error
	sq.once.Do(func() {
		er = sq.init()
	})
	return er
}

func (sq *SQL) init() error {
	var err error
	sq.db, err = sql.Open(sq.Driver, sq.ConnStr)
	if err != nil {
		log.Print(err)
		return err
	}
	sq.prep()
	return nil
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
				return skogul.Error{Source: "db sender", Reason: "unable to marshal metadata into json", Next: err}
			}
			vals = append(vals, meta)
		case marshalData:
			data, err := json.Marshal(m.Data)
			if err != nil {
				return skogul.Error{Source: "db sender", Reason: "unable to marshal data into json", Next: err}
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
	if er := sq.Init(); er != nil {
		log.Print(er)
		return er
	}
	txn, err := sq.db.Begin()
	if err != nil {
		log.Print(err)
		return err
	}

	stmt, err := txn.Prepare(sq.q)
	if err != nil {
		log.Print(err)
		return err
	}

	for _, m := range c.Metrics {
		err = sq.exec(stmt, m)
		if err != nil {
			log.Print(err)
			txn.Rollback()
			return err
		}
	}

	err = stmt.Close()
	if err != nil {
		log.Print(err)
		txn.Rollback()
		return err
	}

	err = txn.Commit()
	if err != nil {
		log.Print(err)
		txn.Rollback()
		return err
	}
	return nil
}

// Verify ensures options are set, but currently doesn't check very well,
// since it is disallowed from connecting to a database and such.
func (sq *SQL) Verify() error {
	if sq.ConnStr == "" {
		return skogul.Error{Source: "sql sender", Reason: "ConnStr is empty"}
	}
	if sq.Query == "" {
		return skogul.Error{Source: "sql sender", Reason: "Query is empty"}
	}
	if sq.Driver == "" {
		return skogul.Error{Source: "sql sender", Reason: "Driver is empty"}
	}
	if sq.Driver != "mysql" && sq.Driver != "postgres" {
		return skogul.Error{Source: "sql sender", Reason: fmt.Sprintf("unsuported database driver %s - must be `mysql' or `postgres'", sq.Driver)}
	}
	return nil
}
