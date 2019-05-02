/*
 * skogul, postgres sender
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
	"log"
	"sync"

	skogul "github.com/KristianLyng/skogul/pkg"
	"github.com/lib/pq"
)

/*
Postgres sender writes to a postgres-database, at the moment using the simplest
schema imaginable with ts (timestamp), metadata (jsonb), data (jsonb). Future
versions will most likely be less stupid schema-wise.
*/
type Postgres struct {
	ConnStr string
	db      *sql.DB
	mux     sync.Mutex
}

/*
Init will connect to the database, ping it and set things up.

Running it is optional. It will be run when the first metric passes
through the Postgres sender, but it might still be a good idea to
run it "yourself" so you can decide what to do before you start receiving
metrics (e.g.: if used as part of a cluster, it might be better to exit/die
if the storage isn't present so upstream skogul instances can use other
members of the cluster - your milage may vary)
*/
func (pqs *Postgres) Init() error {
	var err error
	pqs.db, err = sql.Open("postgres", pqs.ConnStr)
	if err != nil {
		log.Print(err)
		return err
	}
	err = pqs.db.Ping()
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}

func (pqs *Postgres) checkInit() error {
	if pqs.db == nil {
		pqs.mux.Lock()
		defer pqs.mux.Unlock()
		if pqs.db == nil {
			err := pqs.Init()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

/*
Send will send to the Postgres database, after first ensuring
the connection is OK.

Currently it will expect a table named "test" which accepts
"ts", "metadata" and "data" in json-encoded format for metadata and data.

Future versions should probably be smarter, and auto-create the
required tables if they are missing.
*/
func (pqs *Postgres) Send(c *skogul.Container) error {
	var er error
	er = pqs.checkInit()
	if er != nil {
		log.Print(er)
		return er
	}
	txn, err := pqs.db.Begin()
	if err != nil {
		log.Print(err)
		return err
	}

	stmt, err := txn.Prepare(pq.CopyIn("test", "ts", "metadata", "data"))
	if err != nil {
		log.Print(err)
		return err
	}

	for _, m := range c.Metrics {
		meta, er1 := json.Marshal(m.Metadata)
		if er1 != nil {
			log.Print(err)
			txn.Rollback()
			return er1
		}
		data, er2 := json.Marshal(m.Data)
		if er2 != nil {
			log.Print(err)
			txn.Rollback()
			return er2
		}
		_, err = stmt.Exec(m.Time, meta, data)
		if err != nil {
			log.Print(err)
			txn.Rollback()
			return err
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		log.Print(err)
		txn.Rollback()
		return err
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
