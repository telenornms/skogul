/*
 * Copyright (c) 2024 Telenor Norge AS
 * Author(s):
 *  - Hans Rafaelsen <hans.rafaelsen@telenor.no>
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

package transformer

import (
	"errors"
	"fmt"
	"hash/fnv"
	"slices"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/telenornms/skogul"
)

type Bucket_entry struct {
	RoutersLock *sync.RWMutex
	Routers     map[string]*Router_entry
}

type Router_entry struct {
	IFsLock    *sync.RWMutex
	Interfaces map[string]*Interface_entry
}

type Interface_entry struct {
	CoutersLocks *sync.RWMutex
	Counters     map[uint64]DataPoint
}

type DataPoint struct {
	Time  time.Time
	Value uint64
}

var lookup_tbl map[string]uint64
var lookup_tbl_rev map[uint64]string
var postfix_tbl map[string]string
var bitfields_tbl map[string]bool
var deltafields_tbl map[string]bool

const num_buckets = 32

var rate_buckets [num_buckets]*Bucket_entry

func hash_bucket(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32() % num_buckets
}

func get_value(v any, ts time.Time) (DataPoint, bool) {
	pf, ok := v.(float64)
	if !ok {
		return DataPoint{}, false
	}
	return DataPoint{Value: uint64(pf), Time: ts}, true
}

func (rate *Rate) compute_rate(routerName string, interName string, now time.Time, queue int, points map[string]any) {
	bucketNum := hash_bucket(routerName)
	cm := rate_buckets[bucketNum]
	cm.RoutersLock.RLock()
	router, ok := cm.Routers[routerName]
	if ok {
		cm.RoutersLock.RUnlock()
		router.IFsLock.RLock()
		inter, ok := router.Interfaces[interName]
		router.IFsLock.RUnlock()
		if ok {
			// Could have used a readlock until we are ready
			// to write, but
			inter.CoutersLocks.Lock()
			for key, point := range points {
				okey := key
				if queue != -1 {
					key = fmt.Sprintf("%d%s", queue, key)
				}
				_, ok = lookup_tbl[key]
				if !ok {
					// Not a value to compute rate for
					continue
				}
				old, ok := inter.Counters[lookup_tbl[key]]
				if ok {
					// Chekc if update or duplicate
					duration := now.Sub(old.Time).Seconds()
					if duration <= 1 {
						// FIXME: Need to check why we are getting +/- 0.004 in difference and not equal
						continue
					}
					if old.Time.Equal(now) {
						// FIXME: Do we need this as the preivious check catches this case
						continue
					}
					pv, ok := get_value(point, now)
					if !ok {
						rateLog.Log(logrus.WarnLevel, fmt.Sprintf("%s:%s failed to get point for %s", routerName, interName, key))
						continue
					}
					inter.Counters[lookup_tbl[key]] = pv
					// compute rate
					if pv.Value >= old.Value {
						delta := pv.Value - old.Value
						if !now.Equal(old.Time) {
							// Update new data point
							inter.Counters[lookup_tbl[key]] = pv
						}
						var r float64
						_, ok = deltafields_tbl[okey]
						if ok {
							_, ok = bitfields_tbl[okey]
							if ok {
								r = float64(delta * 8)
							} else {
								r = float64(delta)
							}

						} else {
							_, ok = bitfields_tbl[okey]
							if ok {
								r = float64(delta*8) / float64(duration)
							} else {
								r = float64(delta) / float64(duration)
							}
						}
						rate_name, ok := postfix_tbl[key]
						if !ok {
							rateLog.Log(logrus.ErrorLevel, fmt.Sprintf("failed to lookup rate_name for %s", key))
							continue
						}
						points[rate_name] = r
					} else {
						// Counter has wrapped around
						if !now.Equal(old.Time) {
							// Update new data point
							inter.Counters[lookup_tbl[key]] = pv
						}
					}

				} else {
					pv, ok := get_value(point, now)
					if !ok {
						rateLog.Log(logrus.WarnLevel, fmt.Sprintf("%s:%s failed to get point for %s", routerName, interName, key))
						continue
					}
					inter.Counters[lookup_tbl[key]] = pv
				}
			}
			inter.CoutersLocks.Unlock()
		} else {
			// New interName
			l := sync.RWMutex{}
			inter := Interface_entry{CoutersLocks: &l,
				Counters: make(map[uint64]DataPoint, len(points))}
			for key, point := range points {
				if queue != -1 {
					key = fmt.Sprintf("%d%s", queue, key)
				}
				_, ok = lookup_tbl[key]
				if !ok {
					// Not a value to compute rate for
					continue
				}
				pv, ok := get_value(point, now)
				if !ok {
					rateLog.Log(logrus.WarnLevel, fmt.Sprintf("%s:%s failed to get point for %s", routerName, interName, key))
					continue
				}
				inter.Counters[lookup_tbl[key]] = pv
			}
			router.IFsLock.Lock()
			router.Interfaces[interName] = &inter
			router.IFsLock.Unlock()
		}
	} else {
		cm.RoutersLock.RUnlock()
		// New router
		l := sync.RWMutex{}
		inter := Interface_entry{CoutersLocks: &l,
			Counters: make(map[uint64]DataPoint, len(points))}
		for key, point := range points {
			if queue != -1 {
				key = fmt.Sprintf("%d%s", queue, key)
			}
			_, ok = lookup_tbl[key]
			if !ok {
				// Not a value to compute rate for
				continue
			}
			pv, ok := get_value(point, now)
			if !ok {
				rateLog.Log(logrus.WarnLevel, fmt.Sprintf("%s:%s failed to get point for %s", routerName, interName, key))
				continue
			}
			inter.Counters[lookup_tbl[key]] = pv
		}
		lr := sync.RWMutex{}
		router := Router_entry{IFsLock: &lr, Interfaces: make(map[string]*Interface_entry)}
		router.Interfaces[interName] = &inter
		cm.RoutersLock.Lock()
		cm.Routers[routerName] = &router
		cm.RoutersLock.Unlock()
	}
	return
}

func create_lookup_tbl(rate_names, delta_names [][]string, bitfilds_names []string, qos_names []string) {
	lookup_tbl = make(map[string]uint64)
	lookup_tbl_rev = make(map[uint64]string)
	postfix_tbl = make(map[string]string)
	bitfields_tbl = make(map[string]bool)
	deltafields_tbl = make(map[string]bool)

	c := 0
	for i := range num_buckets {
		l := sync.RWMutex{}
		rate_buckets[i] = &Bucket_entry{RoutersLock: &l, Routers: make(map[string]*Router_entry)}
	}
	for _, k := range rate_names {
		if slices.Contains(qos_names, k[0]) {
			for j := range 7 {
				// Different names for different qos classes
				name := fmt.Sprintf("%d%s", j, k[0])
				lookup_tbl[name] = uint64(c)
				lookup_tbl_rev[uint64(c)] = name
				c++
				postfix_tbl[name] = k[1]
			}
		} else {
			postfix_tbl[k[0]] = k[1]
			lookup_tbl[k[0]] = uint64(c)
			lookup_tbl_rev[uint64(c)] = k[0]
			c++
		}
	}
	for _, k := range delta_names {
		if slices.Contains(qos_names, k[0]) {
			for j := range 7 {
				// Different names for different qos classes
				name := fmt.Sprintf("%d%s", j, k[0])
				lookup_tbl[name] = uint64(c)
				deltafields_tbl[name] = true
				lookup_tbl_rev[uint64(c)] = name
				c++
				postfix_tbl[name] = k[1]
			}
		} else {
			postfix_tbl[k[0]] = k[1]
			lookup_tbl[k[0]] = uint64(c)
			deltafields_tbl[k[0]] = true
			lookup_tbl_rev[uint64(c)] = k[0]
			c++
		}
	}
	for _, k := range bitfilds_names {
		bitfields_tbl[k] = true
	}
}

var rateLog = skogul.Logger("transformer", "rate")

type Rate struct {
	RateFields  [][]string `doc:"Fields to compute reate for"`
	DeltaFields [][]string `doc:"Fields to compute delta for"`
	BitFields   []string   `doc:"Fields convert from bytes to bits"`
	QosFields   []string   `doc:"Fields fields with quality classes"`
	SystemId    string     `doc:"String used for naming the system field"`
	InterfaceId string     `doc:"String used for naming the interface filed"`
	QoSClass    string     `doc:"String used for naming the QoS class field"`
	once        sync.Once
	err         error
}

// Transform
func (rate *Rate) Transform(c *skogul.Container) error {
	rate.err = nil
	for _, m := range c.Metrics {
		systemId, ok := m.Metadata[rate.SystemId].(string)
		if !ok {
			rateLog.Log(logrus.ErrorLevel, "Metadata missing or wrong type systemId field")
			rate.err = errors.New(fmt.Sprintf("Metadata missing or wrong type systemId field"))
			continue
		}
		if_name, ok := m.Metadata[rate.InterfaceId].(string)
		if !ok {
			// Seems like it is common that if_name is missing for metrics that have info
			// for the full router.
			continue
		}
		queue_number := float64(-1) // No queue number
		queue, ok := m.Metadata[rate.QoSClass]
		if ok {
			queue_number, ok = queue.(float64)
			if !ok {
				rateLog.Log(logrus.ErrorLevel, fmt.Sprintf("queue_number is not a float %s.%s", systemId, if_name))
				rate.err = errors.New(fmt.Sprintf("queue_number is not a float %s.%s", systemId, if_name))
			}
		}
		rate.compute_rate(systemId, if_name, *m.Time, int(queue_number), m.Data)
	}
	return rate.err
}

// Verify checks that the required variables are set
// Initialize various tables used
func (rate *Rate) Verify() error {
	if len(rate.RateFields) == 0 {
		return skogul.MissingArgument("No ratefields given")
	}
	create_lookup_tbl(rate.RateFields, rate.DeltaFields, rate.BitFields, rate.QosFields)
	return nil
}
