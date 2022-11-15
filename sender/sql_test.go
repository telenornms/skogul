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

// Mysql:
//
// These tests require mariadb or MySQL to run, with a database "skogul",
// user root and password "lol" (....). To get that up and running, the
// simplest solution is a docker container with mariadb:
//
// Working test-setup for MySQL/Mariadb tests:
// docker run -ti -e MYSQL_ROOT_PASSWORD=lol -n kek --network=host mariadb
// docker exec -ti kek mysql
// MariaDB [(none)]> create database skogul
// MariaDB [(none)]> use skogul
// MariaDB [skogul]> create table test (timestamp varchar(100) not null, src varchar(100) not null, name varchar(100) not null, data varchar(100) not null);
//
//
// Postgres:
//
// docker run -n kjeks --network=host -e POSTGRES_PASSWORD=finnlandshette -ti postgres
// docker exec -n kjeks psql -U postgres
// > create database skogul;
// > \c skogul
// > create table test (ts varchar(100) not null, meta varchar(100) not null, data varchar(100) not null);
//
// The "not null" bit is important: Some tests explicitly try sending data
// with missing fields to ensure that errors are caught.

import (
	"fmt"
	"github.com/telenornms/skogul"
	"github.com/telenornms/skogul/config"
	"github.com/telenornms/skogul/sender"
	"testing"
	"time"
)

var sqlBase = `
{
	"senders": {
		"sql": {
			"type": "sql",
			%s
		}
	}
}`

func sqlTestAuto(t *testing.T, url string) *config.Config {
	t.Helper()
	conf, err := config.Bytes([]byte(fmt.Sprintf(sqlBase, url)))
	if conf == nil {
		t.Errorf("Bytes(\"%s\" failed", url)
	}
	if err != nil {
		t.Errorf("Bytes(\"%s\" failed: %v", url, err)
	}
	return conf
}

func sqlTestAutoNeg(t *testing.T, url string) {
	t.Helper()
	conf, err := config.Bytes([]byte(fmt.Sprintf(sqlBase, url)))
	if conf != nil {
		t.Errorf("Bytes(\"%s\" succeeded, but expected failure. Val: %v", url, conf)
	}
	if err == nil {
		t.Errorf("Bytes(\"%s\" succeeded, but expected failure. Val: %v", url, conf)
	}
}

func sqlSender(t *testing.T, conf string) *sender.SQL {
	t.Helper()
	cnf := sqlTestAuto(t, conf)
	if cnf == nil {
		t.Errorf("Failed to build configuration")
		return nil
	}
	s, ok := cnf.Senders["sql"].Sender.(*sender.SQL)
	if !ok {
		t.Errorf("Failed to cast postgres sender to SQL sender?")
		return nil
	}
	return s
}

func getValidContainer() *skogul.Container {
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
	return &c
}

func TestSQL_auto(t *testing.T) {
	sqlTestAutoNeg(t, `"driver":"mysql"`)
	sqlTestAutoNeg(t, `"driver":"mysql","connstr": "something"`)
	sqlTestAutoNeg(t, `"driver":"mysql","query": "something"`)
	sqlTestAuto(t, `"driver":"mysql","connstr":"something","query": "blatti"`)
	sqlTestAuto(t, `"driver":"mysql","connstr":"foo:bar@/blatt", "query":"foo%20bar"`)
	sqlTestAutoNeg(t, `"driver":"postgres"`)
	sqlTestAuto(t, `"driver":"postgres","connstr":"something","query": "blatti"`)
}

func TestSQL_mysql_basic(t *testing.T) {
	s := sqlSender(t, `"driver":"mysql","connstr": "root:lol@/skogul", "query": "INSERT INTO test VALUES(${timestamp},${metadata.src},${name},${data});"`)
	if s == nil {
		t.Errorf("Failed to get sender")
	}
}

func TestSQL_postgres_basic(t *testing.T) {
	s := sqlSender(t, `"driver":"postgres","connstr":"database=skogul sslmode=disable user=postgres password=finnlandshette","query":"INSERT INTO test (ts, meta,data) VALUES(${timestamp},${json.metadata},${json.data})"`)
	if s == nil {
		t.Errorf("Failed to get sender")
	}
}

func TestSQL_mysql_connect(t *testing.T) {
	if testing.Short() {
		t.Skip("Short test: Not connecting to a MySQL database")
	}
	createTable := "create table test (timestamp varchar(100) not null, src varchar(100) not null, name varchar(100) not null, data varchar(100) not null);"
	t.Logf("Assuming database name skogul and table ala: %s", createTable)
	s := sqlSender(t, `"driver":"mysql","connstr": "root:lol@/skogul", "query": "INSERT INTO test VALUES(${timestamp},${metadata.src},${name},${data});"`)

	container := getValidContainer()

	if err := s.Send(container); err != nil {
		t.Errorf("SQL.Send failed: %v", err)
	}

	container.Metrics[0].Data = make(map[string]interface{})
	container.Metrics[0].Data["name"] = "Foo Bar"
	if err := s.Send(container); err == nil {
		t.Errorf("SQL.Send succeeded with missing data field")
	}

	container.Metrics[0].Time = nil
	if err := s.Send(container); err == nil {
		t.Errorf("SQL.Send succeeded with missing timestamp")
	}
}

func TestSQL_mysql_json(t *testing.T) {
	if testing.Short() {
		t.Skip("Short test: Not connecting to a MySQL database")
	}
	s := sqlSender(t, `"driver":"mysql","connstr":"root:lol@/skogul","query":"INSERT INTO test VALUES(${timestamp},'foo',${json.metadata},${json.data})"`)
	if s == nil {
		t.Errorf("Failed to get sender")
		t.Skip("Can't proceed without a sender")
	}
	container := getValidContainer()
	if err := s.Send(container); err != nil {
		t.Errorf("Failed to send to mysql: %v", err)
	}
}

func TestSQL_postgres_json(t *testing.T) {
	if testing.Short() {
		t.Skip("Short test: Not connecting to a MySQL database")
	}
	s := sqlSender(t, `"driver":"postgres","connstr":"database=skogul sslmode=disable user=postgres password=finnlandshette","query":"INSERT INTO test (ts, meta,data) VALUES(${timestamp},${json.metadata},${json.data})"`)
	if s == nil {
		t.Errorf("Failed to get sender")
		t.Skip("Can't proceed without a sender")
	}
	container := getValidContainer()
	if err := s.Send(container); err != nil {
		t.Errorf("Failed to send to postgres: %v", err)
	}
}
