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

import (
	"fmt"
	"github.com/KristianLyng/skogul/sender"
)

/*

FIXME: This needs to be re-done now that New() is gone in favor of json.
Should be fairly  trivial to do New(`{ "type": "mysql", "ConnStr": ...`),
but I'm not sure if we/I want to go down that road just yet.

func mysqlTestAuto(t *testing.T, url string) {
	m, err := sender.New(url)
	if m == nil {
		t.Errorf("New(\"%s\" failed", url)
	}
	if err != nil {
		t.Errorf("New(\"%s\" failed: %v", url, err)
	}
}
func mysqlTestAutoNeg(t *testing.T, url string) {
	m, err := sender.New(url)
	if m != nil {
		t.Errorf("New(\"%s\" succeeded, but expected failure. Val: %v", url, m)
	}
	if err == nil {
		t.Errorf("New(\"%s\" succeeded, but expected failure. Val: %v", url, m)
	}
}

func TestSql_auto(t *testing.T) {
	mysqlTestAutoNeg(t, "mysql:///")
	mysqlTestAutoNeg(t, "mysql:///?connstr=something")
	mysqlTestAutoNeg(t, "mysql:///?query=something")
	mysqlTestAutoNeg(t, "mysql://")
	mysqlTestAuto(t, "mysql://?connstr=something&query=blatti")
	mysqlTestAuto(t, "mysql:///?connstr=something&query=blatti")
	mysqlTestAuto(t, "mysql:///?connstr=foo:bar@/blatt&query=foo%20bar")
}
func TestSql(t *testing.T) {
	m := sender.Sql{}
	s, err := m.GetQuery()
	if err == nil {
		t.Errorf("m.GetQuery() succeeded despite query not being created")
	}
	if s != "" {
		t.Errorf("m.GetQuery() returned data despite query not being created. Got %s.", s)
	}
	query := "INSERT INTO test VALUES(${timestamp.timestamp},${metadata.src},${name},${data});"
	connStr := "root:lol@/skogul"
	m = sender.Sql{Query: query, ConnStr: connStr}
	err = m.Init()
	if err != nil {
		t.Errorf("Sql.Init failed: %v", err)
	}
	want := "INSERT INTO test VALUES(?,?,?,?);"
	var got string
	got, err = m.GetQuery()
	if err != nil {
		t.Errorf("Sql.getQuery() failed: %v", err)
	}
	if want != got {
		t.Errorf("Sql.Init wanted %s got %s", want, got)
	}

	createTable := "create table test (timestamp varchar(100), src varchar(100), name varchar(100), data varchar(100));"
	if testing.Short() == true {
		t.Log("WARNING: Skipping MySQL integration tests with -testing.short on.")
		return
	} else {
		t.Logf("Using MySQL connection string %s", connStr)
		t.Logf("Assuming database name skogul and table ala: %s", createTable)
		t.Logf("If you don't have a suitable MySQL db running, use -test.short")
	}
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
		t.Errorf("Sql.Send failed: %v", err)
	}
	me.Data = make(map[string]interface{})
	me.Data["name"] = "Foo Bar"
	err = m.Send(&c)
	if err != nil {
		t.Errorf("Sql.Send failed: %v", err)
	}
	me.Time = nil
	err = m.Send(&c)
	if err != nil {
		t.Errorf("Sql.Send failed: %v", err)
	}
}

*/
// Basic MySQL example, using user root (bad idea) and password "lol"
// (voted most secure password of 2019), connecting to the database
// "skogul". Also demonstrates printing of the query.
//
// Will, obviously, require a database to be running.
func ExampleSQL() {
	query := "INSERT INTO test VALUES(${timestamp.timestamp},${metadata.src},${name},${data});"
	connStr := "root:lol@/skogul"
	m := sender.SQL{Query: query, ConnStr: connStr, Driver: "mysql"}
	m.Init()
	str, err := m.GetQuery()
	if err != nil {
		panic(err)
	}
	fmt.Println(str)
	// Output:
	// INSERT INTO test VALUES(?,?,?,?);
}
