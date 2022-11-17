/*
 * skogul, cast transformer
 *
 * Copyright (c) 2019-2020 Telenor Norge AS
 * Author(s):
 *  - Kristian Lyngstøl <kly@kly.no>
 *  - Håkon Solbjørg <hakon.solbjorg@telenor.com>
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
	"fmt"
	"math/big"
	"net"
	"strconv"

	"github.com/telenornms/skogul"
)

type Cast struct {
	MetadataStrings    []string `doc:"List of metadatafields that should be strings"`
	MetadataInts       []string `doc:"List of metadatafields that should be integers"`
	MetadataFloats     []string `doc:"List of metadatafields that should be 64-bit floats"`
	MetadataFlatFloats []string `doc:"List of metadatafields that are floats which should be expressed as plain, non-exponential numbers in text. E.g.: Large serial numbers will be written as plain numbers, not 1.1231215e+10. If the field is a non-float, it will be left as is."`
	MetadataIpToDec    []string `doc:"List of metadatafields containing IP addresses that should be decimals"`
	DataStrings        []string `doc:"List of datafields that should be strings"`
	DataInts           []string `doc:"List of datafields that should be integers"`
	DataFloats         []string `doc:"List of datafields that should be 64-bit floats"`
	DataFlatFloats     []string `doc:"List of metadatafields that are floats which should be expressed as plain, non-exponential numbers in text. E.g.: Large serial numbers will be written as plain numbers, not 1.1231215e+10. If the field is a non-float, it will be left as is."`
	DataIpToDec        []string `doc:"List of datafields containing IP addresses that should be decimals"`
}

// Transform enforces the Cast rules
func (cast *Cast) Transform(c *skogul.Container) error {
	for mi := range c.Metrics {
		if c.Metrics[mi].Data != nil {
			for _, value := range cast.DataStrings {
				if c.Metrics[mi].Data[value] != nil {
					_, ok := c.Metrics[mi].Data[value].(string)
					if ok {
						continue
					}
					c.Metrics[mi].Data[value] = fmt.Sprintf("%v", c.Metrics[mi].Data[value])
				}
			}
			for _, value := range cast.DataFloats {
				if c.Metrics[mi].Data[value] != nil {
					_, ok := c.Metrics[mi].Data[value].(float64)
					if ok {
						continue
					}
					var tmp float64
					_, err := fmt.Sscanf(fmt.Sprintf("%v", c.Metrics[mi].Data[value]), "%f", &tmp)
					if err != nil {
						return err
					}
					c.Metrics[mi].Data[value] = tmp
				}
			}
			for _, value := range cast.DataInts {
				if c.Metrics[mi].Data[value] != nil {
					_, ok := c.Metrics[mi].Data[value].(int)
					if ok {
						continue
					}
					var tmp int
					_, err := fmt.Sscanf(fmt.Sprintf("%v", c.Metrics[mi].Data[value]), "%d", &tmp)
					if err != nil {
						return err
					}
					c.Metrics[mi].Data[value] = tmp
				}
			}
			for _, value := range cast.DataFlatFloats {
				if c.Metrics[mi].Data[value] != nil {
					f, ok := c.Metrics[mi].Data[value].(float64)
					if !ok {
						continue
					}
					c.Metrics[mi].Data[value] = strconv.FormatFloat(f, 'f', -1, 64)
				}
			}
			for _, value := range cast.DataIpToDec {
				if c.Metrics[mi].Data[value] != nil {
					s, ok := c.Metrics[mi].Data[value].(string)
					if !ok {
						continue
					}
					c.Metrics[mi].Data[value] = cast.Inet6Aton(s)
				}
			}
		}
		if c.Metrics[mi].Metadata == nil {
			continue
		}
		for _, value := range cast.MetadataStrings {
			if c.Metrics[mi].Metadata[value] != nil {
				_, ok := c.Metrics[mi].Metadata[value].(string)
				if ok {
					continue
				}
				c.Metrics[mi].Metadata[value] = fmt.Sprintf("%v", c.Metrics[mi].Metadata[value])
			}
		}
		for _, value := range cast.MetadataFloats {
			if c.Metrics[mi].Metadata[value] != nil {
				_, ok := c.Metrics[mi].Metadata[value].(float64)
				if ok {
					continue
				}
				var tmp float64
				_, err := fmt.Sscanf(fmt.Sprintf("%v", c.Metrics[mi].Metadata[value]), "%f", &tmp)
				if err != nil {
					return err
				}
				c.Metrics[mi].Metadata[value] = tmp
			}
		}
		for _, value := range cast.MetadataInts {
			if c.Metrics[mi].Metadata[value] != nil {
				_, ok := c.Metrics[mi].Metadata[value].(int)
				if ok {
					continue
				}
				var tmp int
				_, err := fmt.Sscanf(fmt.Sprintf("%v", c.Metrics[mi].Metadata[value]), "%d", &tmp)
				if err != nil {
					return err
				}
				c.Metrics[mi].Metadata[value] = tmp
			}
		}
		for _, value := range cast.MetadataFlatFloats {
			if c.Metrics[mi].Metadata[value] != nil {
				f, ok := c.Metrics[mi].Metadata[value].(float64)
				if !ok {
					continue
				}
				c.Metrics[mi].Metadata[value] = strconv.FormatFloat(f, 'f', -1, 64)
			}
		}
		for _, value := range cast.MetadataIpToDec {
			if c.Metrics[mi].Metadata[value] != nil {
				s, ok := c.Metrics[mi].Metadata[value].(string)
				if !ok {
					continue
				}
				c.Metrics[mi].Metadata[value] = cast.Inet6Aton(s)
			}
		}
	}
	return nil
}

// Convers IPv4 or IPv6 to decimals
func (cast *Cast) Inet6Aton(ipAddress string) *big.Int {
	ip := net.ParseIP(ipAddress)

	if ip == nil {
		x := big.NewInt(0)
		return x
	}

	ipInt := big.NewInt(0)
	if ip.To4() != nil {
		ipInt.SetBytes(ip.To4())
		return ipInt
	}

	ipInt.SetBytes(ip.To16())
	return ipInt
}
