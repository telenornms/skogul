
.. image:: https://goreportcard.com/badge/github.com/KristianLyng/skogul
   :target: https://goreportcard.com/report/github.com/KristianLyng/skogul

.. image:: https://godoc.org/github.com/KristianLyng/skogul?status.svg
   :target: https://godoc.org/github.com/KristianLyng/skogul

======================================
Skogul - generic metric/data collector
======================================

Skogul is a generic tool for moving metric data around. It can serve as a
collector of data, but is primarily designed to be a framework for building
bridges between data collectors and storage engines.

This repository contains the Skogul library/package, and ``cmd/skogul``,
which parses a JSON-config to set up Skogul.

.. contents:: Table of contents
   :depth: 2
   :local:

Quickstart
----------

You need to install a recent/decent version of Go. Either from your
favorite Linux distro, or through https://golang.org/dl/ .

Building ``skogul``, including cloning::

   $ git clone https://github.com/KristianLyng/skogul
   (...)
   $ cd skogul/cmd/skogul
   $ go build
   $ 
   # (No output from go build is good)

Alternatively, you can use ``go install`` instead of ``go build`` to
install to ``$GOPATH/bin``, which is typically ``~/go/bin``.

To use the locally imported/vendored packages instead of downloading them
directly, e.g. if a system does not have direct internet access or you wish
to take a local copy of the code in its entirety, including dependencies.
First make a vendored copy on an internet-attached computer - checksums in
the repo will be verified::

   $ cd skogul
   $ go mod vendor
   $
   ( skogul/vendor is now populated after a while )

Copy repo/directory to relevant computer, then run::

   $ cd skogul/cmd/skogul
   $ go build -mod vendor
   $

(or ``go install -mod vendor``)


About
-----

A Skogul chain is built from one or more independent receivers which
receive data and pass it on to a sender. A sender can either transmit data
to an external source (including an other Skogul instance), or make
internal changes to data before passing it on to one or more other senders.

.. image:: docs/basic.png

Unlike most APIs or collectors of metrics, Skogul does NOT have a
preference when it comes to storage engine. It is explicitly designed to
disconnect the task of how data is collected from how it is stored.

The rationale is that the problem of writing an efficient snmp collector
should not be tightly coupled to where you store the data. And where you
store the data should not be tightly coupled with how you receive it, or
what you do with it.

Extra care has been put into making it trivial to write senders and
receivers. For example, an author of a new sender only has to add tags
to their data structure to have that exposed as documentation.

See the package documentation over at godoc for more usage:
https://godoc.org/github.com/KristianLyng/skogul

More discussion on architecture can be found in `docs/`.

Performance
-----------

Skogul is meant to scale well. At present time, there are known flaws in
the implementation, but still, simple local testing on a laptop is able to
produce decent results.

.. image:: docs/skogul-rates.png

The above graph is from a very simple test on a laptop (with a quad core
i7), using the provided tester to write data to influxdb. It demonstrates
that despite well-known weaknesses (specially in the influx-writer), we're
able to push roughly 600-800k values/s through Skogul.

The laptop in question was using about 150-190% CPU for skogul and 400% for
InfluxDB, the rest went to the testers. No real attempt at tuning was done,
but a few different scenarios were tested.

Note that the general values/s is decent both with a ton of values for each
metric, and just a handful of values per metric, but tons of metrics per
containers.

As future work will introduce buffers and "batch aggregators" to make it
better equipped to handle irregular traffic, it's is expected and
acceptable that performance dips when the number of values per container
drops.

Update:

As of September 2019, TLS was enabled and Skogul was tested again, just for
TLS. Skogul was seen sending roughly 2 million key:values/s over HTTPS on
the same laptop. The batch sender has also proven to be very valuable.

Name
----

Skogul is a Valkyrie. After extensive research (5 minutes on Wikipedia with
a cross-check on duckduckgo), this name was selected because it is
reasonably unique and is also a Valkyrie, like Gondul, a sister-project.

Hacking
-------

There is little "exotic" about Skogul hacking, so the following sections
are aimed mostly at people who are unfamiliar with Go.


.. note::
   
   The majority of all documentation is kept in godoc source comments, and
   available either in the code directly, through ``go doc
   github.com/KristianLyng/skogul`` or  through the web, at
   https://godoc.org/github.com/KristianLyng/skogul . This includes, but is
   not limited to example code and API documentation.

Testing
.......

To run test cases, ``go test`` can be run. This can be used either in
individual directories, or at the top directory, with ``go test ./...``
(note the triple dots. This is a go-ism for recursive behavior).

To produce coverage analysis, use::

   $ cd skogul
   $ go test ./... -covermode=count -coverprofile=coverage.out
   $ go tool cover -html coverage.out
   // Opens a browser with coverage anlysis

Be aware that the MySQL sender does not do integration testing by default,
as that requires a working MySQL instance.

Tests are extracted from ``*_test.go`` files, and start with the name
``Test`` followed by a function or data structure, optionally followed by
an underscore and an arbitrary name to support multiple tests of the same
function/type. E.g. ``TestValidate()``, ``TestHTTP_foobar()`` etc.

Runnable examples follow the same style, but are named Example, not Test.

Formatting etc
..............

The "go report" at the top of this document is a decent test of
marginal OK-ish-ness.

Tools you should use:

- `gofmt`, to format code according to Go coding style. Use ``gofmt -d .``
  see local diff, or ``gofmt -w .`` to fix it.
- `golint` to lint your code. ``golint .``

Installing these tools is left as an exercise to the reader.

Documentation
.............

Documentation is written and maintained using code comments and runnable
examples, following the ``godoc`` approach. Some architecture comments are
kept in ``docs//``, but by and large, documentation should be consumed from
godoc.

See https://godoc.org/github.com/KristianLyng/skogul for the online
version, or use ``go doc github.com/KristianLyng/skogul`` or similar,
as you would any other go package.

Examples are part of the test suite and thus extracted from ``*_test.go``
where applicable.

Roadmap
-------

The configuration backend was just introduced. It took a few iterations,
and will most likely be updated slightly. This is the groundwork that is
required to ensure a healthy development environment.

This introduced a shift in focus. Previously, ``skogul-x2y`` was provided
as a binary to set up simple, but commonly used Skogul-chains. As Skogul
grew, this became a bottleneck, because exposing more complex configuration
was hard. As such, the idea was to write custom-binaries for more complex
chains.

With the new JSON-based configuration, this seems redundant. As such, focus
will be on simplifying the ``cmd/skogul`` binary user experience, and
streamlining development.

Immediately, that means work on documentation, re-writing a number of now
broken tests, and generally tweaking things to see how it works.

One thing that needs to be done, however, is provide better feedback on
invalid configuration. Including when options are provided that are not
used.

Time-wise, we hope to do a release in 2019 when we feel Skogul is mature
enough. It is already in use.
