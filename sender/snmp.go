package sender

import (
	"fmt"
	"sync"
	"time"

	"github.com/gosnmp/gosnmp"
	"github.com/telenornms/skogul"
)

var snmpLog = skogul.Logger("sender", "snmp")

type SNMP struct {
	Port        uint16         `doc:"Snmp port. Default: 162, the standard trap port."`
	Community   string         `doc:"Snmp communit field"`
	Version     string         `doc:"Snmp version possible values: 2c, 3"`
	Target      string         `doc:"Snmp target"`
	Oidmap      map[string]any `doc:"Snmp oid to json field mapping"`
	Timeout     uint           `doc:"Snmp timeout, default 5 seconds"`
	SnmpTrapOID string         `doc:"Value of the snmp trap oid pdu"`
	RetryConfig

	// mu protects g. gosnmp.GoSNMP is not safe for concurrent use:
	// SendTrap mutates the handle's message and request IDs as it goes.
	mu sync.Mutex
	g  *gosnmp.GoSNMP
}

const (
	defaultSNMPPort    = 162
	defaultSNMPTimeout = 5 * time.Second
)

// Verify checks the configuration of the SNMP sender.
func (x *SNMP) Verify() error {
	if x.Target == "" {
		return skogul.MissingArgument("Target")
	}
	return x.verifyRetry()
}

// snmpVersion maps the configured version onto gosnmp's, defaulting to 2c.
func (x *SNMP) snmpVersion() gosnmp.SnmpVersion {
	if x.Version == "3" {
		return gosnmp.Version3
	}
	return gosnmp.Version2c
}

/*
 * SNMP trap sender
 */

// connect opens the socket to the trap target. The handle is built from
// scratch each time rather than re-connected: gosnmp sets up per-
// connection state (message and request IDs, receive buffer) in
// Connect(), and re-using a handle whose connect failed leaves that state
// unset. Caller must hold x.mu.
func (x *SNMP) connect() error {
	port := x.Port
	if port == 0 {
		port = defaultSNMPPort
	}
	timeout := defaultSNMPTimeout
	if x.Timeout != 0 {
		timeout = time.Duration(x.Timeout) * time.Second
	}

	g := &gosnmp.GoSNMP{
		Target:             x.Target,
		Port:               port,
		Transport:          "udp",
		Community:          x.Community,
		Version:            x.snmpVersion(),
		Timeout:            timeout,
		Retries:            1,
		ExponentialTimeout: true,
		MaxOids:            gosnmp.MaxOids,
	}
	if err := g.Connect(); err != nil {
		return err
	}
	x.g = g
	return nil
}

// disconnect drops the socket so the next send reconnects from scratch.
// Caller must hold x.mu.
func (x *SNMP) disconnect() {
	if x.g == nil {
		return
	}
	if x.g.Conn != nil {
		x.g.Conn.Close()
	}
	x.g = nil
}

// buildTrap builds a trap for a single metric.
func (x *SNMP) buildTrap(m *skogul.Metric) gosnmp.SnmpTrap {
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

		oid, ok := x.Oidmap[j]
		if !ok {
			// Without a mapping there is no OID to send this
			// under, and a PDU named for a nil OID only fails
			// later, in gosnmp.
			snmpLog.WithField("field", j).Debug("No OID mapped for field, skipping")
			continue
		}
		pduName := fmt.Sprintf("%s", oid)

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
			// Appending the zero-value PDU here would send a
			// nameless, valueless variable.
			snmpLog.WithField("field", j).Debugf("Unsupported value type %T for SNMP, skipping", i)
			continue
		}
		pdutypes = append(pdutypes, pdutype)
	}

	trap := gosnmp.SnmpTrap{}
	trap.Variables = pdutypes
	trap.IsInform = false
	trap.Enterprise = "no"
	trap.AgentAddress = "localhost"
	return trap
}

func (x *SNMP) Send(c *skogul.Container) error {
	// Build all traps up front; the retries only re-send the ones that
	// failed, so successfully delivered traps are not duplicated.
	traps := make([]gosnmp.SnmpTrap, 0, len(c.Metrics))
	skipped := 0
	for _, m := range c.Metrics {
		trap := x.buildTrap(m)
		if len(trap.Variables) == 0 {
			// Nothing in the metric could be mapped to a PDU, and
			// SendTrap rejects a trap without at least one. Same
			// policy as the per-field skipping in buildTrap, so
			// count them and carry on rather than failing the
			// whole container.
			skipped++
			continue
		}
		traps = append(traps, trap)
	}
	if skipped > 0 {
		snmpLog.Warnf("Skipped %d of %d metrics with no mappable fields, check Oidmap and SnmpTrapOID", skipped, len(c.Metrics))
	}
	if len(traps) == 0 {
		return nil
	}
	pending := traps
	return retryNetwork(&x.RetryConfig, snmpLog, func() error {
		// Each attempt holds the lock: gosnmp.GoSNMP is not safe for
		// concurrent use, and the connect/send/disconnect sequence
		// mutates x.g. Backoff sleeps happen in retryNetwork,
		// outside the lock.
		x.mu.Lock()
		defer x.mu.Unlock()

		// Connecting here, rather than once at startup, is what lets
		// a retry recover: a target that was unreachable when the
		// first container arrived is reachable on a later attempt.
		if x.g == nil {
			if err := x.connect(); err != nil {
				return fmt.Errorf("unable to connect to SNMP target %s: %w", x.Target, err)
			}
		}

		var failed []gosnmp.SnmpTrap
		var firstErr error
		for _, trap := range pending {
			if _, err := x.g.SendTrap(trap); err != nil {
				failed = append(failed, trap)
				if firstErr == nil {
					firstErr = err
				}
			}
		}
		if firstErr == nil {
			return nil
		}
		total := len(pending)
		pending = failed
		if isTransientError(firstErr) {
			// Looks like the socket rather than the payload, so
			// redial on the next attempt instead of re-using a
			// handle that just failed.
			x.disconnect()
		}
		return fmt.Errorf("%d of %d traps failed, first error: %w", len(failed), total, firstErr)
	})
}
