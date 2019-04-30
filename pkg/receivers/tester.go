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

package receivers

import (
	"fmt"
	"github.com/KristianLyng/skogul/pkg"
	"log"
	"math/rand"
	"net/url"
	"time"
)

// Tester synthesise dummy data
type Tester struct {
	Metrics int64
	Values  int64
	Handler skogul.Handler
}

func (tst *Tester) generate(t time.Time) skogul.Container {
	c := skogul.Container{}
	c.Template.Time = &t
	c.Metrics = make([]skogul.Metric, tst.Metrics)
	for i := int64(0); i < tst.Metrics; i++ {
		m := skogul.Metric{}
		m.Metadata = map[string]interface{}{}
		m.Metadata["key1"] = i
		m.Data = map[string]interface{}{}
		for key := int64(0); key < tst.Values; key++ {
			m.Data[fmt.Sprintf("metric%d", key)] = rand.Int63()
		}
		c.Metrics[i] = m
	}
	return c
}

// Start never returns.
func (tst *Tester) Start() error {
	for {
		c := tst.generate(time.Now())
		err := tst.Handler.Sender.Send(&c)
		if err != nil {
			log.Print(err)
		}
	}
	return skogul.Error{Reason: "Shouldn't reach this"}
}

func init() {
	addAutoReceiver("test", NewTester, "Generate dummy-data, each container contains $m metrics and each metric $v values, format: tester://$m/$v")
}

/*
NewTester returns a new Tester receiver, building values/metrics from URL.
*/
func NewTester(ul url.URL, h skogul.Handler) skogul.Receiver {

	var host int64
	var path int64
	n, err := fmt.Sscanf(ul.Host, "%d", &host)
	if n != 1 || err != nil {
		log.Fatalf("a Invalid URL for Tester %s (n: %d err: %v)", ul.Host, n, err)
	}
	n, err = fmt.Sscanf(ul.Path, "/%d", &path)
	if n != 1 || err != nil {
		log.Fatalf("a Invalid URL for Tester %s (n: %d err: %v)", ul.Path, n, err)
	}

	return &Tester{Metrics: host, Values: path, Handler: h}
}
