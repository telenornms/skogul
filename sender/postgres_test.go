/*
 * skogul, postgres tests
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
	"github.com/KristianLyng/skogul/config"
	"github.com/KristianLyng/skogul/sender"
	"testing"
)

var postgresBase = `
{
	"senders": {
		"postgres": {
			"type": "sql",
			"driver": "postgres"
			%s
		}
	}
}`

func postgresTestAuto(t *testing.T, url string) *config.Config {
	conf, err := config.Bytes([]byte(fmt.Sprintf(postgresBase, url)))
	if conf == nil {
		t.Errorf("Bytes(\"%s\" failed", url)
	}
	if err != nil {
		t.Errorf("Bytes(\"%s\" failed: %v", url, err)
	}
	return conf
}

func postgresTestAutoNeg(t *testing.T, url string) {
	conf, err := config.Bytes([]byte(fmt.Sprintf(postgresBase, url)))
	if conf != nil {
		t.Errorf("Bytes(\"%s\" succeeded, but expected failure. Val: %v", url, conf)
	}
	if err == nil {
		t.Errorf("Bytes(\"%s\" succeeded, but expected failure. Val: %v", url, conf)
	}
}

func TestSQL_postgres_auto(t *testing.T) {
	postgresTestAutoNeg(t, ``)
	postgresTestAuto(t, `,"connstr":"something","query": "blatti"`)
}

func TestSQL_postgres_json(t *testing.T) {
	cnf := postgresTestAuto(t, `,"connstr":"database=skogul sslmode=disable user=postgres password=finnlandshette","query":"INSERT INTO test (ts, meta,data) VALUES(${timestamp.timestamp},${json.metadata},${json.data})"`)
	if cnf == nil {
		t.Errorf("Failed to build configuration")
		return
	}
	s, ok := cnf.Senders["postgres"].Sender.(*sender.SQL)
	if !ok {
		t.Errorf("Failed to cast postgres sender to SQL sender?")
		return
	}
	container := getValidContainer()
	err := s.Send(container)
	if err != nil {
		t.Errorf("Failed to send to postgres: %v", err)
	}
}
