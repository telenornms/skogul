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
	"flag"
	"fmt"
	"github.com/KristianLyng/skogul"
	"github.com/KristianLyng/skogul/sender"
	"testing"
	"time"
)

var ftestMysql = flag.Bool("test.mysql", false, "Enable integration tests for mysql")

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

func TestMysql_auto(t *testing.T) {
	mysqlTestAutoNeg(t, "mysql:///")
	mysqlTestAutoNeg(t, "mysql:///?connstr=something")
	mysqlTestAutoNeg(t, "mysql:///?query=something")
	mysqlTestAutoNeg(t, "mysql://")
	mysqlTestAuto(t, "mysql://?connstr=something&query=blatti")
	mysqlTestAuto(t, "mysql:///?connstr=something&query=blatti")
	mysqlTestAuto(t, "mysql:///?connstr=foo:bar@/blatt&query=foo%20bar")
}
func TestMysql(t *testing.T) {
	if *ftestMysql != true {
		fmt.Printf("WARNING: Skipping MySQL integration tests. Use -test.mysql to execute them.\n")
		return
	}
	m := sender.Mysql{}
	s, err := m.GetQuery()
	if err == nil {
		t.Errorf("m.GetQuery() succeeded despite query not being created")
	}
	if s != "" {
		t.Errorf("m.GetQuery() returned data despite query not being created. Got %s.", s)
	}
	m = sender.Mysql{Query: "INSERT INTO test VALUES(${timestamp.timestamp},${metadata.src},${name},${data});", ConnStr: "root:lol@/skogul"}
	err = m.Init()
	if err != nil {
		t.Errorf("Mysql.Init failed: %v", err)
	}
	want := "INSERT INTO test VALUES(?,?,?,?);"
	var got string
	got, err = m.GetQuery()
	if err != nil {
		t.Errorf("Mysql.getQuery() failed: %v", err)
	}
	if want != got {
		t.Errorf("Mysql.Init wanted %s got %s", want, got)
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
		t.Errorf("Mysql.Send failed: %v", err)
	}
	me.Data = make(map[string]interface{})
	me.Data["name"] = "Foo Bar"
	err = m.Send(&c)
	if err != nil {
		t.Errorf("Mysql.Send failed: %v", err)
	}
	me.Time = nil
	err = m.Send(&c)
	if err != nil {
		t.Errorf("Mysql.Send failed: %v", err)
	}
}

// Basic MySQL example, using user root (bad idea) and password "lol"
// (voted most secure password of 2019), connecting to the database
// "skogul". Also demonstrates printing of the query.
//
// Will, obviously, require a database to be running.
func ExampleMysql() {
	m := sender.Mysql{Query: "INSERT INTO test VALUES(${timestamp.timestamp},${metadata.src},${name},${data});", ConnStr: "root:lol@/skogul"}
	m.Init()
	str, err := m.GetQuery()
	if err != nil {
		panic(err)
	}
	fmt.Println(str)
	// Output:
	// INSERT INTO test VALUES(?,?,?,?);
}
