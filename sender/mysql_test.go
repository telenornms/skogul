/*
 * skogul, mysql tests
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

package sender_test

// These tests require mariadb or MySQL to run, with a database "skogul",
// user root and password "lol" (....). To get that up and running, the
// simplest solution is a docker container with mariadb:
//
// Working test-setup for MySQL/Mariadb tests:
// docker run -ti -e MYSQL_ROOT_PASSWORD=lol -n kek --network=host mariadb
// docker exec -ti kek mysql
// MariaDB [(none)]> create database skogul
// MariaDB [(none)]> use skogul
// MariaDB [skogul]> create table test (timestamp varchar(100), src varchar(100), name varchar(100), data varchar(100));

import (
	"flag"
	"fmt"
	"github.com/KristianLyng/skogul"
	"github.com/KristianLyng/skogul/config"
	"github.com/KristianLyng/skogul/sender"
	"os"
	"testing"
	"time"
)

var mysqlBase = `
{
	"senders": {
		"mysql": {
			"type": "sql",
			"driver": "mysql"
			%s
		}
	}
}`

func mysqlTestAuto(t *testing.T, url string) {
	conf, err := config.Bytes([]byte(fmt.Sprintf(mysqlBase, url)))
	if conf == nil {
		t.Errorf("Bytes(\"%s\" failed", url)
	}
	if err != nil {
		t.Errorf("Bytes(\"%s\" failed: %v", url, err)
	}
}

var flag_mysql = flag.Bool("mysql", false, "Test mysql")

func TestMain(m *testing.M) {
	flag.Parse()
	// call flag.Parse() here if TestMain uses flags
	os.Exit(m.Run())
}
func mysqlTestAutoNeg(t *testing.T, url string) {
	conf, err := config.Bytes([]byte(fmt.Sprintf(mysqlBase, url)))
	if conf != nil {
		t.Errorf("Bytes(\"%s\" succeeded, but expected failure. Val: %v", url, conf)
	}
	if err == nil {
		t.Errorf("Bytes(\"%s\" succeeded, but expected failure. Val: %v", url, conf)
	}
}

func TestSQL_auto(t *testing.T) {
	mysqlTestAutoNeg(t, ``)
	mysqlTestAutoNeg(t, `,"connstr": "something"`)
	mysqlTestAutoNeg(t, `,"query": "something"`)
	mysqlTestAuto(t, `,"connstr":"something","query": "blatti"`)
	mysqlTestAuto(t, `,"connstr":"foo:bar@/blatt", "query":"foo%20bar"`)
}

func TestSQL(t *testing.T) {
	m := sender.SQL{}
	s, err := m.GetQuery()
	if err == nil {
		t.Errorf("m.GetQuery() succeeded despite query not being created")
	}
	if s != "" {
		t.Errorf("m.GetQuery() returned data despite query not being created. Got %s.", s)
	}
	query := "INSERT INTO test VALUES(${timestamp.timestamp},${metadata.src},${name},${data});"
	connStr := "root:lol@/skogul"
	m = sender.SQL{Query: query, ConnStr: connStr, Driver: "mysql"}
	err = m.Init()
	if err != nil {
		t.Errorf("SQL.Init failed: %v", err)
	}
	want := "INSERT INTO test VALUES(?,?,?,?);"
	var got string
	got, err = m.GetQuery()
	if err != nil {
		t.Errorf("SQL.getQuery() failed: %v", err)
	}
	if want != got {
		t.Errorf("SQL.Init wanted %s got %s", want, got)
	}

	createTable := "create table test (timestamp varchar(100), src varchar(100), name varchar(100), data varchar(100));"
	if *flag_mysql == false {
		t.Log("WARNING: Skipping MySQL integration tests. Use `go test -mysql' to run them.")
		return
	}
	t.Logf("Using MySQL connection string %s", connStr)
	t.Logf("Assuming database name skogul and table ala: %s", createTable)

	c := skogul.Container{}
	me := skogul.Metric{}
	n := time.Now()
	me.Time = &n
	me.Metadata = make(map[string]interface{})
	me.Data = make(map[string]interface{})
	me.Metadata["src"] = "Test"
	me.Data["name"] = "Foo Bar"
	me.Data["data"] = "something"
	c.Metrics = []*skogul.Metric{&me}

	err = m.Send(&c)
	if err != nil {
		t.Errorf("SQL.Send failed: %v", err)
	}
	me.Data = make(map[string]interface{})
	me.Data["name"] = "Foo Bar"
	err = m.Send(&c)
	if err != nil {
		t.Errorf("SQL.Send failed: %v", err)
	}
	me.Time = nil
	err = m.Send(&c)
	if err != nil {
		t.Errorf("SQL.Send failed: %v", err)
	}
}
