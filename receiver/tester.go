/*
 * skogul, test receiver
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

package receiver

import (
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"time"

	"github.com/KristianLyng/skogul"
)

// Tester synthesise dummy data
type Tester struct {
	Metrics int64
	Values  int64
	Threads int
	Handler skogul.Handler
}

func (tst *Tester) generate(t time.Time) skogul.Container {
	c := skogul.Container{}
	c.Template = &skogul.Metric{}
	c.Template.Time = &t
	c.Metrics = make([]*skogul.Metric, tst.Metrics)
	for i := int64(0); i < tst.Metrics; i++ {
		m := skogul.Metric{}
		m.Metadata = map[string]interface{}{}
		m.Metadata["key1"] = i
		m.Data = map[string]interface{}{}
		for key := int64(0); key < tst.Values; key++ {
			m.Data[fmt.Sprintf("metric%d", key)] = rand.Int63()
		}
		c.Metrics[i] = &m
	}
	return c
}

// Start never returns.
func (tst *Tester) Start() error {
	for i := 1; i < tst.Threads; i++ {
		go tst.run()
	}
	tst.run()
	return nil
}

func (tst *Tester) run() {
	for {
		c := tst.generate(time.Now())
		for _, t := range tst.Handler.Transformers {
			t.Transform(&c)
		}
		err := tst.Handler.Sender.Send(&c)
		if err != nil {
			log.Print(err)
		}
	}
}

func init() {
	addAutoReceiver("test", NewTester, "Generate dummy-data, each container contains $m metrics and each metric $v values, multiplied by $t threads. All parameters are optional. Example: test:///?threads=4&metrics=2&values=12")
}

type myval struct {
	v url.Values
}

func (m myval) getDefault(key string, def int) (int, error) {
	values := m.v
	str := values.Get(key)
	result := def
	if str != "" {
		n, err := fmt.Sscanf(str, "%d", &result)
		if n != 1 || err != nil {
			return def, skogul.Error{Source: "tester sender", Reason: fmt.Sprintf("invalid parameter \"%s\". Value is %s", key, str), Next: err}
		}
	}
	return result, nil
}

/*
NewTester returns a new Tester receiver, building values/metrics from URL.
*/
func NewTester(ul url.URL, h skogul.Handler) skogul.Receiver {
	values := myval{v: ul.Query()}
	var metrics, vals, threads int
	var err error
	metrics, err = values.getDefault("metrics", 10)
	if err != nil {
		log.Print(err)
		return nil
	}
	vals, err = values.getDefault("values", 50)
	if err != nil {
		log.Print(err)
		return nil
	}
	threads, err = values.getDefault("threads", 4)
	if err != nil {
		log.Print(err)
		return nil
	}

	return &Tester{Metrics: int64(metrics), Values: int64(vals), Threads: threads, Handler: h}
}
