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
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"time"

	"github.com/KristianLyng/skogul"
)

// Tester synthesise dummy data. Use New("tester:///?metrics=foo) to create it.
type Tester struct {
	flags   *flag.FlagSet
	Metrics *int64
	Values  *int64
	Threads *int
	Delay   *time.Duration
	Handler skogul.Handler
}

func (tst *Tester) generate(t time.Time) skogul.Container {
	c := skogul.Container{}
	c.Template = &skogul.Metric{}
	c.Template.Time = &t
	c.Metrics = make([]*skogul.Metric, *tst.Metrics)
	for i := int64(0); i < *tst.Metrics; i++ {
		m := skogul.Metric{}
		m.Metadata = map[string]interface{}{}
		m.Metadata["key1"] = i
		m.Data = map[string]interface{}{}
		for key := int64(0); key < *tst.Values; key++ {
			m.Data[fmt.Sprintf("metric%d", key)] = rand.Int63()
		}
		c.Metrics[i] = &m
	}
	return c
}

// Start never returns.
func (tst *Tester) Start() error {
	if tst.flags == nil {
		tst.flags = testerFlags()
	}
	for i := 1; i < *tst.Threads; i++ {
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
		if tst.Delay != nil && *tst.Delay != 0 {
			time.Sleep(*tst.Delay)
		}
	}
}

func init() {
	n := AutoReceiver{
		Init:  newTester,
		Help:  "Generate dummy-data, each container contains $m metrics and each metric $v values, multiplied by $t threads. A delay of $d is inserted between \"runs\". All parameters are optional. Example: test:///?threads=4&metrics=2&values=12&delay=1s",
		Flags: testerFlags,
	}
	newAutoReceiver("test", &n)
}

func testerFlags() *flag.FlagSet {
	x := allocTester()
	return x.flags
}

func allocTester() *Tester {
	t := Tester{}
	fs := flag.NewFlagSet("tester", flag.ExitOnError)
	t.Metrics = fs.Int64("metrics", 10, "Number of metrics per container")
	t.Values = fs.Int64("values", 50, "Number of values per metric")
	t.Threads = fs.Int("threads", 4, "Number of go routines to run in parallel")
	t.Delay = fs.Duration("delay", 0, "Delay between containers are produced")
	t.flags = fs
	return &t
}

/*
newTester returns a new Tester receiver, building values/metrics from URL.
*/
func newTester(ul url.URL, h skogul.Handler) skogul.Receiver {
	t := allocTester()
	err := URLParse(ul, t.flags)
	if err != nil {
		log.Printf("%v", err)
		return nil
	}
	t.Handler = h
	return t
}
