/*
 * skogul, sender automation
 *
 * Copyright (c) 2019-2020 Telenor Norge AS
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
	"github.com/telenornms/skogul"
)

// Auto maps sender-names to sender implementation, used for auto
// configuration.
var Auto skogul.ModuleMap

func init() {
	Auto.Add(skogul.Module{
		Name:    "backoff",
		Aliases: []string{"retry"},
		Alloc:   func() interface{} { return &Backoff{} },
		Help:    "Forwards data to the next sender, retrying after a delay upon failure. For each retry, the delay is doubled. Gives up after the set number of retries.",
	})
	Auto.Add(skogul.Module{
		Name:    "batch",
		Aliases: []string{"batcher"},
		Alloc:   func() interface{} { return &Batch{} },
		Help:    "Accepts metrics and puts them in a shared container. When the container either has a set number of metrics (Threshold), or a timeout occurs, the entire container is forwarded. This allows down-stream senders to work with larger batches of metrics at a time, which is frequently more efficient. A side effect of this is that down-stream errors are not propogated upstream. That means any errors need to be dealt with down stream, or they will be ignored.",
	})
	Auto.Add(skogul.Module{
		Name:    "counter",
		Aliases: []string{"count"},
		Alloc:   func() interface{} { return &Counter{} },
		Help:    "Accepts metrics, counts them and passes them on. Then emits statistics to the Stats-handler on an interval. Useful for housekeeping, highly recommended both during testing and production for internal Skogul-metrics.",
	})
	Auto.Add(skogul.Module{
		Name:     "debug",
		Aliases:  []string{"print"},
		Alloc:    func() interface{} { return &Debug{} },
		Help:     "Prints received metrics to stdout.",
		AutoMake: true,
	})
	Auto.Add(skogul.Module{
		Name:    "detacher",
		Aliases: []string{"detach"},
		Alloc:   func() interface{} { return &Detacher{} },
		Help:    "Returns OK without waiting for the next sender to finish. The detached part is single-threaded.",
	})
	Auto.Add(skogul.Module{
		Name:    "dupe",
		Aliases: []string{"dup", "duplicate"},
		Alloc:   func() interface{} { return &Dupe{} },
		Help:    "Sends the same metrics to all senders listed in Next.",
	})
	Auto.Add(skogul.Module{
		Name:    "errdiverter",
		Aliases: []string{"errordiverter", "errdivert", "errordivert"},
		Alloc:   func() interface{} { return &ErrDiverter{} },
		Help:    "Forwards data to next sender. If an error is returned, the error is converted into a Skogul container and sent to the err-handler. This provides the means of logging errors through regular skogul-chains. See the logrus receiver for a more solid approach to diverting all log messages, instead of individually failed containers.",
	})
	Auto.Add(skogul.Module{
		Name:  "fanout",
		Alloc: func() interface{} { return &Fanout{} },
		Help:  "Fan out (load balance) to a fixed number of threads before passing data on. This is rarely needed, as receivers should do this.",
	})
	Auto.Add(skogul.Module{
		Name:  "fallback",
		Alloc: func() interface{} { return &Fallback{} },
		Help:  "Tries the senders provided in Next, in order. E.g.: if the first responds OK, the second will never get data. Useful for diverting traffic to alternate paths upon failure.",
	})
	Auto.Add(skogul.Module{
		Name:  "file",
		Alloc: func() interface{} { return &File{} },
		Help:  "Writes metrics to a file.",
	})
	Auto.Add(skogul.Module{
		Name:  "forwardfail",
		Alloc: func() interface{} { return &ForwardAndFail{} },
		Help:  "Forwards metrics, but always returns failure. Useful in complex failure handling involving e.g. fallback sender, where it might be used to write log or stats on failure while still propogating a failure upward.",
	})
	Auto.Add(skogul.Module{
		Name:    "http",
		Aliases: []string{"https"},
		Alloc:   func() interface{} { return &HTTP{} },
		Help:    "Sends Skogul-formatted JSON-data to a HTTP endpoint (e.g.: an other Skogul instance?). Highly useful in scenarios with multiple data collection methods spread over several servers.",
	})
	Auto.Add(skogul.Module{
		Name:    "influx",
		Aliases: []string{"influxdb"},
		Alloc:   func() interface{} { return &InfluxDB{} },
		Help:    "Send to a InfluxDB HTTP endpoint. The sender can either always send the data to a single measurement, send it to a measurement extracted from the metadata of a metric, or a combination where the \"measurement\" serves as a default measurement to use if the metric doesn't have the key presented in \"measurementfrommetadata\".",
	})
	Auto.Add(skogul.Module{
		Name:  "log",
		Alloc: func() interface{} { return &Log{} },
		Help:  "Logs a message, mainly useful for enriching debug information in conjunction with, for example, dupe and debug.",
	})
	Auto.Add(skogul.Module{
		Name:    "mnr",
		Aliases: []string{"m&r"},
		Alloc:   func() interface{} { return &MnR{} },
		Help:    "Sends M&R line format to a TCP endpoint.",
	})
	Auto.Add(skogul.Module{
		Name:  "mqtt",
		Alloc: func() interface{} { return &MQTT{} },
		Help:  "Publishes received metrics to an MQTT broker/topic.",
	})
	Auto.Add(skogul.Module{
		Name:  "sql",
		Alloc: func() interface{} { return &SQL{} },
		Help:  "Execute a SQL query for each received metric, using a template. Any query can be run, and if multiple metrics are present in the same container, they are all executed in a single transaction, which means the batch-sender will greatly increase performance. Supported engines are MySQL/MariaDB and Postgres.",
	})
	Auto.Add(skogul.Module{
		Name:     "null",
		Alloc:    func() interface{} { return &Null{} },
		Help:     "Discards all data. Mainly useful for testing the receive-pipeline decoupled from storage.",
		AutoMake: true,
	})
	Auto.Add(skogul.Module{
		Name:  "sleep",
		Alloc: func() interface{} { return &Sleeper{} },
		Help:  "Injects a random delay before passing data on. Mainly for testing.",
	})
	Auto.Add(skogul.Module{
		Name:    "splunk",
		Aliases: []string{"hec"},
		Alloc:   func() interface{} { return &Splunk{HTTP: &HTTP{}} },
		Help:    "A sender to Splunk HEC",
	})
	Auto.Add(skogul.Module{
		Name:  "net",
		Alloc: func() interface{} { return &Net{} },
		Help:  "Sends json data to a network endpoint.",
	})
	Auto.Add(skogul.Module{
		Name:   "switch",
		Alloc:  func() interface{} { return &Switch{} },
		Help:   "Sends data selectively based on metedata.",
		Extras: []interface{}{Match{}},
	})
	Auto.Add(skogul.Module{
		Name:    "enrichmentupdater",
		Aliases: []string{"eupdater"},
		Alloc:   func() interface{} { return &EnrichmentUpdater{} },
		Help:    "Updates the enrichment database of an enrichment transformer.",
	})
	Auto.Add(skogul.Module{
		Name:  "test",
		Alloc: func() interface{} { return &Test{} },
		Help:  "Used for internal testing. Basically just discards data but provides an internal counter of received data",
	})
	Auto.Add(skogul.Module{
		Name:  "kafka",
		Alloc: func() interface{} { return &Kafka{} },
		Help:  "EXPERIMENTAL Kafka sender",
	})

}
