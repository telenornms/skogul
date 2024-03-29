package sender

import (
	"fmt"
	"sync"
	"time"

	"github.com/gosnmp/gosnmp"
	"github.com/telenornms/skogul"
)

type SNMP struct {
	Port        uint16                 `doc:"Snmp port"`
	Community   string                 `doc:"Snmp communit field"`
	Version     string                 `doc:"Snmp version possible values: 2c, 3"`
	Target      string                 `doc:"Snmp target"`
	Oidmap      map[string]interface{} `doc:"Snmp oid to json field mapping"`
	Timeout     uint                   `doc:"Snmp timeout, default 5 seconds"`
	r           sync.Once
	err         error
	g           *gosnmp.GoSNMP
	SnmpTrapOID string `doc:"Value of the snmp trap oid pdu"`
}

/*
 * SNMP trap sender
 */
func (x *SNMP) init() {
	var version gosnmp.SnmpVersion
	if x.Version == "2c" {
		version = gosnmp.Version2c
	} else if x.Version == "3" {
		version = gosnmp.Version3
	} else {
		version = gosnmp.Version2c
	}

	if x.Timeout == 0 {
		x.Timeout = 5
	}

	x.g = &gosnmp.GoSNMP{
		Port:               x.Port,
		Transport:          "udp",
		Community:          x.Community,
		Version:            version,
		Timeout:            time.Duration(x.Timeout) * time.Second,
		Retries:            1,
		ExponentialTimeout: true,
		MaxOids:            gosnmp.MaxOids,
	}

	x.g.Target = x.Target
	x.err = x.g.Connect()
}

func (x *SNMP) Send(c *skogul.Container) error {
	x.r.Do(func() {
		x.init()
	})

	var errors []error
	for _, m := range c.Metrics {
		var pdutypes []gosnmp.SnmpPDU
		if x.SnmpTrapOID != "" {
			pdutypes = append(pdutypes, gosnmp.SnmpPDU{
				Value: x.SnmpTrapOID,
				Type:  gosnmp.ObjectIdentifier,
				Name:  ".1.3.6.1.6.3.1.1.4.1.0",
			})
		}
		for j, i := range m.Data {
			var pdutype gosnmp.SnmpPDU

			pduName := fmt.Sprintf("%s", x.Oidmap[j])

			switch i.(type) {
			case string:
				pdutype = gosnmp.SnmpPDU{
					Value: i,
					Name:  pduName,
					Type:  gosnmp.OctetString,
				}
			case bool:
				pdutype = gosnmp.SnmpPDU{
					Value: i,
					Name:  pduName,
					Type:  gosnmp.Boolean,
				}
			case float64:
				k := int(i.(float64))
				pdutype = gosnmp.SnmpPDU{
					Value: k,
					Name:  pduName,
					Type:  gosnmp.Integer,
				}
			default:
			}
			pdutypes = append(pdutypes, pdutype)
		}

		trap := gosnmp.SnmpTrap{}
		trap.Variables = pdutypes
		trap.IsInform = false
		trap.Enterprise = "no"
		trap.AgentAddress = "localhost"
		_, err := x.g.SendTrap(trap)
		if err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) == 0 {
		return nil
	}
	return fmt.Errorf("%d of %d traps failed, first error: %w", len(errors), len(c.Metrics), errors[0])
}
