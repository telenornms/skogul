/*
 * skogul, common trivialities
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

package skogul

import (
	"crypto/x509"
	"fmt"
	"math"
	"os"
	"runtime"

	"github.com/sirupsen/logrus"
)

/*
bit patterns for 32-bit infinity. See ../common.go for details.

This is needed because the math package doesn't handle 32-bit variants.

The pattern is:

seeeeeee e0000000 00000000 00000000
01234567 01234567 01234567 01234567

Where s is sign, e is exponent and the rest is the value.

Exponents of all 1's means either Infinity or NaN.

Value of 0, exponent of 1's, mean infinity, sign will determine, well, the
sign (e.g. Inf vs -Inf).

Could also set this as:

	uvneginf = 0b11111111100000000000000000000000
	uvinf    = 0b01111111100000000000000000000000

Hex equivalents are just weird, so I left it as binary after first figuring
out the hex encoding...

	uvneginf = 0xFF800000
	uvinf = 0x7F800000

nanmask, while identical to uvinf, is there to explicitly distinguish
between checking if this is a nan or not, but isn't currently implemented.
*/
const (
	uvneginf = 0b11111111100000000000000000000000
	uvinf    = 0b01111111100000000000000000000000
	nanmask  = 0b01111111100000000000000000000000
	//           01234567012345670123456701234567
)

/*
Handler determines what a receiver will do with data received. It requires
a parser to interperet the raw data, 0 or more transformers to mutate
Containers, and a sender to call after data is parsed and mutated and ready
to be dealt with.

Whenever a new Container is created, it should pass that to a Handler, not
directly to a Sender. This goes for artificially created data too, e.g. if
a sender wants to emit statistics. This ensures that transformers can be
used in the future.

IngorePartialFailures should probably be renamed to removeinvalidmetrics or
something like that, as that's closer to what it does.

To make it configurable, a HandlerRef should be used.
*/
type Handler struct {
	parser                Parser
	Transformers          []Transformer
	Sender                Sender
	IgnorePartialFailures bool
}

// Parser is the interface for parsing arbitrary data into a Container
type Parser interface {
	Parse(data []byte) (*Container, error)
}

/*
Sender accepts data through Send() - and "sends it off". The canonical
sender is one that implements a storage backend or outgoing API. E.g.:
accept data, send to influx.

Senders are not allowed to modify the Container - there could be multiple
goroutines running with same Container. If modification is required, the
Sender needs to take a copy.

A sender should assume that the container has been validated, and is
non-null. Slightly counter to common sense, it is NOT recommended to
verify the input data again, since multiple senders are likely chained
and will thus likely redo the same verifications.

Senders that pass data off to other senders should use a SenderRef instead,
to facilitate configuration.
*/
type Sender interface {
	Send(c *Container) error
}

/*
Transformer mutates a collection before it is passed to a sender. Transformers
should be very fast, but are the only means to modifying the data.
*/
type Transformer interface {
	Transform(c *Container) error
}

/*
Receiver is how we get data. Receivers are responsible for getting raw data and the
outer boundaries of a Container, but should explicitly avoid parsing raw data.
This ensures that how data is transported is not bound by how it is parsed.
*/
type Receiver interface {
	Start() error
}

/*
Encoder is an *optional* way to encode data, it is used by senders where
data encoding can vary, but not all senders use it.
*/
type Encoder interface {
	Encode(c *Container) ([]byte, error)
	EncodeMetric(m *Metric) ([]byte, error)
}

/*
Verifier is an *optional* interface for modules. If
implemented, the configuration engine will issue Verify() after all
configuration is parsed. The module should never modify state upon
Verify(), but should simply check that internal state is usable.
*/
type Verifier interface {
	Verify() error
}

/*
Deprecated is an optional interface for modules which, if implemented,
allows modules to signal deprecated usage. If implemented, the Deprecated()
method will be called after configuration is loaded and verified. Returning
an error signals deprecation, and can be used both to deprecate specific
settings and entire modules. There is currently no plan on how to remove
deprecated modules over time.
*/
type Deprecated interface {
	Deprecated() error
}

/*
Stats is an optional interface for all skogul modules. It is used
by modules to export internal stats about them, such as received data,
specific errors, successes and sent metrics.
The GetStats() method is called by the stats package for each configured module
and the stats are sent onto the stats channel.
*/
type Stats interface {
	GetStats() *Metric
}

/*
SenderRef is a reference to a named sender. This is required to allow
references to be resolved after all senders are loaded. Wherever a
Sender is loaded from configuration, a SenderRef should be used in its
place. The maintenance of the sender is handled in the configuration
system.
*/
type SenderRef struct {
	S    Sender
	Name string
}

// HandlerRef references a named handler. Used whenever a handler is
// defined by configuration.
type HandlerRef struct {
	H    *Handler
	Name string
}

// TransformerRef is a string mapping to a Transformer.
// It is used during configuration/transformer setup.
type TransformerRef struct {
	T    Transformer
	Name string
}

// ParserRef is a string mapping to a Parser.
// It is used during configuration setup.
type ParserRef struct {
	P    Parser
	Name string
}

// EncoderRef is a string mapping to an Encoder.
type EncoderRef struct {
	E    Encoder
	Name string
}

// SetParser sets the parser to use for a Handler
func (h *Handler) SetParser(p Parser) error {
	if h.parser != nil {
		return fmt.Errorf("handler already has a parser set")
	}
	if p == nil {
		return fmt.Errorf("attempting to set parser to `nil'")
	}
	h.parser = p
	return nil
}

// Parse parses the bytes into a Container
func (h *Handler) Parse(b []byte) (*Container, error) {
	c, err := h.parser.Parse(b)
	if err != nil {
		return nil, fmt.Errorf("parsing failed: %w", err)
	}
	return c, nil
}

// Transform runs all available transformers
func (h *Handler) Transform(c *Container) error {
	for _, t := range h.Transformers {
		if err := t.Transform(c); err != nil {
			return err
		}
	}
	return nil
}

// Send validates the container and sends it to the configured sender
func (h *Handler) Send(c *Container) error {
	if err := c.Validate(h.IgnorePartialFailures); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}
	if err := h.Sender.Send(c); err != nil {
		return fmt.Errorf("sender failed: %w", err)
	}
	return nil
}

// Handle parses the byte array using the configured parser, issues
// transformers and sends the data off.
func (h *Handler) Handle(b []byte) error {
	c, err := h.Parse(b)
	if err != nil {
		return err
	}
	if err = h.TransformAndSend(c); err != nil {
		return err
	}
	return nil
}

// TransformAndSend transforms the already parsed container and sends the
// data off.
func (h *Handler) TransformAndSend(c *Container) error {
	if err := h.Transform(c); err != nil {
		return fmt.Errorf("transforming metrics failed: %w", err)
	}
	if err := h.Send(c); err != nil {
		return fmt.Errorf("sending metrics failed: %w", err)
	}
	return nil
}

// Verify the basic integrity of a handler. Quite shallow.
func (h Handler) Verify() error {
	if h.parser == nil {
		return fmt.Errorf("missing parser")
	}
	for i, t := range h.Transformers {
		if t == nil {
			return fmt.Errorf("nil-transformer %d for handler", i)
		}
	}
	if h.Sender == nil {
		return fmt.Errorf("missing sender")
	}
	return nil
}

// Logger returns a logrus.Entry pre-populated with standard Skogul fields.
// category is the typical family of the code/module:
// sender/receiver/parser/transformer/core, while implementation is the
// local implementation (http, json, protobuf, udp, etc).
func Logger(category, implementation string) *logrus.Entry {
	return logrus.WithField("category", category).WithField(category, implementation)
}

// AssertErrors counts the number of assert errors
var AssertErrors int

// Assert panics if x is false, useful mainly for doing error-checking for
// "impossible scenarios" we can't really handle anyway.
//
// Keep in mind that net/http steals panics, but you can check
// AssertErrors, which is incremented with each assert error encountered.
func Assert(x bool, v ...interface{}) {
	if !x {
		out := "assertion failed"
		pc, file, line, ok := runtime.Caller(1)
		if ok {
			out = fmt.Sprintf("%X:%s:%d assertion failed", pc, file, line)
		}
		AssertErrors++
		panic(fmt.Sprintf("%s %s", out, fmt.Sprint(v...)))
	}
}

// ExtractNestedObject extracts an object from a nested object structure. All intermediate objects has to map[string]interface{}
func ExtractNestedObject(object map[string]interface{}, keys []string) (map[string]interface{}, error) {
	if len(keys) == 1 {
		return object, nil
	}

	next, ok := object[keys[0]].(map[string]interface{})

	if !ok {
		return nil, fmt.Errorf("failed to cast nested object to map[string]interface{}")
	}

	return ExtractNestedObject(next, keys[1:])
}

// Secret is a common type that wraps a string where the contents of the string
// can be sensitive, such as a credential. The String() func will output `***` to prevent accidental exposure,
// but the raw contents can be `Expose()`d.
type Secret string

// String replaces the underlying data with the string "<redacted>"
// so that it is not accidentally revealed in logs or other debug related outputs.
func (s Secret) String() string {
	return "<redacted>"
}

// Expose must be called when the underlying secret is to be revealed,
// such as to the service that requires the data.
func (s Secret) Expose() string {
	return string(s)
}

// IsInf is a direct copy of Math.IsInf(), converted to 32-bit. This is
// some times needed when encodings use 32 bit floats (e.g.: Juniper
// telemetry). The following is the original documentation.
//
// IsInf reports whether f is an infinity, according to sign.
// If sign > 0, IsInf reports whether f is positive infinity.
// If sign < 0, IsInf reports whether f is negative infinity.
// If sign == 0, IsInf reports whether f is either infinity.
func IsInf(f float32, sign int) bool {
	// Test for infinity by comparing against maximum float.
	// To avoid the floating-point hardware, could use:
	x := math.Float32bits(f)
	return sign >= 0 && x == uvinf || sign <= 0 && x == uvneginf
	//        return sign >= 0 && f > math.MaxFloat32 || sign <= 0 && f < -math.SmallestNonzeroFloat32
}

// MissingArgument provides a standard error message for use in Verify to
// report missing arguments.
func MissingArgument(field string) error {
	return fmt.Errorf("missing required configuration option `%s'", field)
}

// GetCertPool accepts a path to a directory of certificate authorities to trust. Pass in the empty string to use system defaults
func GetCertPool(path string) (*x509.CertPool, error) {
	// this means "use system default"
	if path == "" {
		return nil, nil
	}
	cp := x509.NewCertPool()
	fd, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("unable to open custom root CA: %w", err)
	}
	defer func() {
		fd.Close()
	}()
	bytes := make([]byte, 1024000)
	n, err := fd.Read(bytes)
	if err != nil {
		return nil, fmt.Errorf("unable to read custom root CA: %w", err)
	}
	ok := cp.AppendCertsFromPEM(bytes[:n])
	if !ok {
		return nil, fmt.Errorf("unable to append certificate to root CA pool")
	}
	return cp, nil
}
