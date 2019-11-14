Coding
======

This document describes architectural aspects of developing (for) Skogul.
You should read CODE_OF_CONDUCT.md and CONTRIBUTING first.

This document is meant to guide, not dictate.

Division of labour
------------------

Senders, receivers, transformers and parsers are written to be isolated and
modularized. Currently, parsers do now have a module map since none of them
have configuration options, but the others do. All can be considered
modules.

A module should not be concerned with the inner workings of any other
module. A module should never modify data outside its own data. And a
sender should never rely on data that is not directly sent to it.

Ideally, it should be impossible for a module to access any data that it
shouldn't have access to. Alas, that is not always achievable. As such,
even if there are global data structures available, do not use them in a
module.

More specifically, the division of labor is:

Common
......

If it is possible to configure a module in a way that is invalid, you must
check for this. Where reasonable, this should be done in a Verify()
function _and_ at run-time. The run-time check can be skipped if it is
expensive.

The Verify() function of a module should never change state, and must be
safe to repeat. It should verify the validity of the configuration, but it
should accept configuration that currently doesn't work, but might work in
the future. E.g.: If you rely on a database, Verify() can check that the
required connection paramaters are present, but should not fail if
connecting to the database is unsuccessful. At the very least, that should
be configurable. This is to ensure that Skogul will start even if the
external systems are not currently available.

Receivers
.........

Receivers is where data/metrics originates. Typically from outside sources,
but theoretically from internal sources (e.g.: test-data or log-data).

Receivers should not parse that data if at all possible. That is the job of
a parser, as configured through a handler. This is to ensure maximum
flexibility in the future: If we can accept data on UDP we can accept it
both as JSON-encoded data and telemetry data.

Receivers should, by default, use multiple go routines where it makes
sense. Receivers should, where possible, expose sender-errors to the
caller, unless configured to do otherwise.

Receivers should never use global scope. Multiple instances of the same
receiver implementation, with different configuration, might be configured.

Transformers
............

This one is tricky.

Transformers can mutate containers and metrics. The mutated result should
be valid.

Transformers must be concurrency-safe.

Transformers can copy data, modify data and more. But it is important that
the resulting metrics do not have multiple references to the same data. If
they do, this creates a problem if an other transformer subsequently tries
to modify the content of one of the references - it will unintentionally
also modify the other reference.

Transformers should never depend on other special configuration items.
E.g.: It is illegal to write a transformer that assumes that the templating
transformer will be run later.

Senders
.......

Senders. Never. Modify. Containers. Ever.

Senders must be concurrency-safe. They can and will be accessed in multiple
go-routines at the same time.

Other than that, they are fairly unrestrained.

Logging
-------

We currently use logrus.

Each module should acquire a logger that indicates the type of code it is,
and the name of the implementation. Use additional fields where it makes
sense.

Ideally, what to log should be self-explanatory from the regular log
levels. They range from:

error
      Something is very wrong. E.g.: Failure to connect to required
      services. Expect loss of service.

warn
      Something is somewhat wrong. Expect reduced service. E.g.: Connection
      to external service was OK, but the service didn't accept our data.
      Disconnected from service, but reconnect worked.  Also useful for
      configuration issues that are valid, but likely problematic. E.g.:
      transmitting passwords in plain-text. 

info
      Moderate issues and general information. E.g.: Connection successful
      for persistent connections (e.g.: database connections).
      Initialization chatter.

debug
      We get chatty. Details of config parsing/validation. We want debug to
      be something you can run in production if you are seeing periodic
      issues, but not something you need or generally want.

trace/verbose
      Used for development and bug-hunting. Details of code paths that are
      valid, etc. Should not be needed for any regular user.

Some items are left out: request-oriented logs.

Receivers should not log OK requests by default, but logging failed
requests is acceptable. Logging OK requests can be provided as a
configuration option. This is to ensure that things do not explode on
high-throughput installations.

Similar practices should apply to senders and transformers.

If a module returns an error, it should either NOT log, or log at trace. If
a module uses an other module and the child-module fails, the parent-module
should:

1. Not log, or log at trace, if it propagates the error upward (e.g.:
   returns error).
2. Log Info at most if it knows how to handle an error and that succeeds,
   e.g. if the Fallback sender first tries sender A, which fails, then
   tries sender B, that succeeds, this is typically Info.
3. Log Warn at most if it can't recover from the error and isn't
   propagating it. E.g.: The batch sender.

These rules are not absolute and certainly not implemented all over yet,
but should be a starting point.
