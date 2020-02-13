package config

import (
	"github.com/telenornms/skogul"
	"github.com/telenornms/skogul/parser"
	"github.com/telenornms/skogul/receiver"
	"github.com/telenornms/skogul/sender"
)

// RunTestWithConfig runs skogul with only a debugger
// The debugger can be configured all the way through
// using a configuration file, and by overriding either
// the receiver, handler or sender. The ones not defined
// will be automatically defined with some defaults.
// If nothing is defined, the default behaviour is to
// * Read from /dev/stdin
// * Parse the input as "skogul-JSON" (container format)
// * Run no transformers
// * Print the result to stdout
func RunTestWithConfig(conf *Config) error {
	const debuggerName = "_debugger"

	debugSender := sender.Debug{}

	var debugHandlerRef skogul.HandlerRef
	// If no handler is specified, create a simple one.
	if conf.Handlers[debuggerName] == nil {
		debugHandler := skogul.Handler{
			Sender: &debugSender,
		}
		debugHandler.SetParser(parser.JSON{})

		debugHandlerRef = skogul.HandlerRef{
			H:    &debugHandler,
			Name: debuggerName,
		}
	} else {
		// If a _debugger handler is specified, re-use relevant pieces
		// of the config and re-apply them to relevant configurations
		_debugHandler := conf.Handlers[debuggerName]
		transformers := make([]skogul.Transformer, 0)

		for _, t := range _debugHandler.Transformers {
			transformers = append(transformers, t.T)
		}

		debugHandler := skogul.Handler{
			Sender:       _debugHandler.Sender.S,
			Transformers: transformers,
		}
		debugHandler.SetParser(parser.JSON{})
		debugHandlerRef = skogul.HandlerRef{
			H:    &debugHandler,
			Name: debuggerName,
		}
	}

	// If a receiver is specified, use that.
	// Otherwise, initialise a simple stdin reader.
	if conf.Receivers[debuggerName] == nil {
		recvHandler := Receiver{
			Receiver: &receiver.Stdin{
				Handler: debugHandlerRef,
			},
		}
		conf.Receivers[debuggerName] = &recvHandler
	}

	return conf.Receivers[debuggerName].Receiver.Start()
}
