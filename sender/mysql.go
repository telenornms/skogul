/*
 * skogul, mysql sender
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
	"log"
	"net/url"
	"os"
	"sync"

	"github.com/KristianLyng/skogul"
	_ "github.com/go-sql-driver/mysql" // Imported for side effect/mysql support
)

func init() {
	addAutoSender("mysql", newMysql, "Write to a MySQL database, use parameters connstr for connection string, and query for... query. E.g.: mysql:///?connstr=root:lol@/skogul&query=INSERT%20INTO%20test%20VALUES%28%24%7Btimestamp%2Etimestamp%7D%2C%27hei%27%2C%24%7Bmetadata%2Ekey1%7D%2C%24%7Bmetric1%7D%29")
}

// newMysql creates a new Mysql sender
func newMysql(ul url.URL) skogul.Sender {
	x := Mysql{}
	values := ul.Query()
	query := values.Get("query")
	conn := values.Get("connstr")
	if query == "" || conn == "" {
		log.Printf("Invalid url for mysql. Need connstr and query.")
		return nil
	}
	x.Query = query
	x.ConnStr = conn
	return &x
}

const (
	timestamp = iota
	metadata  = iota
	data      = iota
)

type dbElement struct {
	family int
	key    string
}

/*
Mysql sender accepts a ConnStr according to
https://github.com/go-sql-driver/mysql/, and a query which is expanded
using os.Expand(), allowing arbitrary queries to be executed.

Variable expansion is done through ${foo}, ${timestamp.timestamp} and
${metadata.foo}. This is done using a prepared statement and is thus
safe and relatively fast.
*/
type Mysql struct {
	ConnStr string
	Query   string
	q       string
	list    []dbElement
	db      *sql.DB
	once    sync.Once
}

/*
prep parses my.Query into q and populates my.list accordingly
*/
func (my *Mysql) prep() {
	//str := "INSERT INTO lol VALUES($date, $name, $foo)"

	mlen := len("metadata.")

	expander := func(element string) string {
		if element == "timestamp.timestamp" {
			my.list = append(my.list, dbElement{timestamp, element})
		} else if len(element) > mlen && element[0:mlen] == "metadata." {
			my.list = append(my.list, dbElement{metadata, element[mlen:]})
		} else {
			my.list = append(my.list, dbElement{data, element})
		}
		return "?"
	}
	my.q = os.Expand(my.Query, expander)
}

// GetQuery returns the parsed query, assuming there is one.
func (my *Mysql) GetQuery() (string, error) {
	if my.Query == "" {
		return "", skogul.Error{Source: "mysql sender", Reason: "No Query set, but GetQuery() called"}
	}
	err := my.Init()
	if err != nil {
		return "", skogul.Error{Source: "mysql sender", Reason: "Mysql.Init failed during GetQuery()", Next: err}
	}
	return my.q, nil
}

/*
Init will connect to the database, ping it and set things up. But only once.
*/
func (my *Mysql) Init() error {
	var er error
	my.once.Do(func() {
		er = my.init()
	})
	return er
}

func (my *Mysql) init() error {
	var err error
	my.db, err = sql.Open("mysql", my.ConnStr)
	if err != nil {
		log.Print(err)
		return err
	}
	err = my.db.Ping()
	if err != nil {
		log.Print(err)
		return err
	}
	my.prep()
	return nil
}

func (my *Mysql) exec(stmt *sql.Stmt, m *skogul.Metric) error {
	var vals []interface{}
	for _, e := range my.list {
		switch e.family {
		case timestamp:
			vals = append(vals, m.Time)
		case metadata:
			vals = append(vals, m.Metadata[e.key])
		case data:
			vals = append(vals, m.Data[e.key])
		}
	}
	_, err := stmt.Exec(vals...)
	return err
}

/*
Send will send to the MySQL database, after first ensuring
the connection is OK.
*/
func (my *Mysql) Send(c *skogul.Container) error {
	er := my.Init()
	if er != nil {
		log.Print(er)
		return er
	}
	txn, err := my.db.Begin()
	if err != nil {
		log.Print(err)
		return err
	}

	stmt, err := txn.Prepare(my.q)
	if err != nil {
		log.Print(err)
		return err
	}

	for _, m := range c.Metrics {
		err = my.exec(stmt, m)
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
