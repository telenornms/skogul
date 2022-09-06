/*
 * skogul, SQL receiver
 *
 * Copyright (c) 2022 Telenor Norge AS
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

package receiver

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql" // Imported for side effect/mysql support
	_ "github.com/lib/pq"
	"github.com/telenornms/skogul"
	"reflect"
	"time"
)

var sqlLog = skogul.Logger("receiver", "sql")

type SQL struct {
	ConnStr  string            `doc:"Connection string to use for database. Slight variations between database engines. For MySQL typically user:password@tcp(host:port)/database. For  MySQL, you need to add parseTime=true at the end to successfully parse a time column, e.g foo:bar@tcp(db2)/blatti?parseTime=true" example:"mysql: 'root:lol@/mydb' postgres: 'user=pqgotest dbname=pqgotest sslmode=verify-full'"`
	Query    string            `doc:"Query run for each metric. Any column named 'time' will be used as the metric time stamp."`
	Metadata []string          `doc:"Array of which columns to treat as metadata, the rest will be data fields."`
	Driver   string            `doc:"Database driver/system. Currently suported: mysql and postgres."`
	Interval skogul.Duration   `doc:"How often to run the query. Set to negative value to run it just once."`
	Handler  skogul.HandlerRef `doc:"Handler to use for data transmission."`
}

// Start the SQL receiver and never return
// This is still a monstrosity
func (s *SQL) Start() error {
	db, err := sql.Open(s.Driver, s.ConnStr)
	if err != nil {
		return fmt.Errorf("couldn't initialize SQL connection: %w", err)
	}
	stmt, err := db.Prepare(s.Query)
	if err != nil {
		return fmt.Errorf("couldn't create prepared statemet from query:  %w", err)
	}
	err = db.Ping()
	if err != nil {
		return fmt.Errorf("couldn't ping the database: %w", err)
	}
	tRawBytes := reflect.TypeOf(sql.RawBytes{})
	tString := reflect.TypeOf("")

	// Need a reverse map here to quickly check if columns are metadata
	// or not
	isMetadata := make(map[string]bool)
	for _, name := range s.Metadata {
		isMetadata[name] = true
	}

	// We want to sleep even if (specially if) we do a continue
	// anywhere, but not on initial startup, so we abuse for ;; a bit.
	for ; ; time.Sleep(s.Interval.Duration) {
		c := skogul.Container{}
		c.Metrics = make([]*skogul.Metric, 0)
		rows, err := stmt.Query()
		if err != nil {
			sqlLog.WithError(err).Error("couldn't run query")
			continue
		}
		columnt, err := rows.ColumnTypes()
		if err != nil {
			fmt.Printf("error: %v\n", err)
			continue
		}

		for rows.Next() {
			// We scan into values, but need to prepare it
			values := make([]interface{}, len(columnt))

			// Allocates type-specific values to scan into,
			// including a work-around for the mysql driver (at
			// least?) returns sql.RawBytes even for regular
			// varchar() data, which is rather annoying.
			for idx := range columnt {
				t := columnt[idx].ScanType()
				if t == tRawBytes {
					t = tString
				}
				values[idx] = reflect.New(t).Interface()
			}
			err = rows.Scan(values...)
			if err != nil {
				sqlLog.WithError(err).Error("Scan error")
				continue
			}

			metric := skogul.Metric{}
			metric.Metadata = make(map[string]interface{})
			metric.Data = make(map[string]interface{})

			// Store data where we actually want it
			for idx := range columnt {
				name := columnt[idx].Name()
				oldValue := reflect.ValueOf(values[idx])
				newValue := reflect.Indirect(oldValue).Interface()

				if isMetadata[name] {
					metric.Metadata[name] = newValue
				} else if name == "time" {
					ts, ok := newValue.(sql.NullTime)
					if !ok {
						sqlLog.Warnf("Unable to parse time column as timestamp. Value: %#v", newValue)
						// I considered storing
						// this as either metadata
						// or data, but in the case
						// where this would happen,
						// I couldn't really see a
						// good outcome of
						// accidentally creating
						// Influx tags for each
						// time stamp for example.
						// metric.Data[name] = newValue
						continue
					}
					metric.Time = &ts.Time
				} else {
					metric.Data[name] = newValue
				}
			}
			c.Metrics = append(c.Metrics, &metric)
		}
		if err := rows.Close(); err != nil {
			sqlLog.Errorf("couldn't close rows objects, this is really strange: %v", err)
		}

		if err := s.Handler.H.TransformAndSend(&c); err != nil {
			sqlLog.Errorf("Failed to transform and send metrics: %v", err)
		}

		// With 0 or negative interval we just run this once and return
		if s.Interval.Duration < time.Nanosecond {
			return nil
		}
	}
}
