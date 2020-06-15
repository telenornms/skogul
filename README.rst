
.. image:: https://goreportcard.com/badge/github.com/telenornms/skogul
   :target: https://goreportcard.com/report/github.com/telenornms/skogul

.. image:: https://godoc.org/github.com/telenornms/skogul?status.svg
   :target: https://godoc.org/github.com/telenornms/skogul

.. image:: https://cloud.drone.io/api/badges/telenornms/skogul/status.svg
   :target: https://cloud.drone.io/telenornms/skogul

======================================
Skogul - generic metric/data collector
======================================

Skogul is a generic tool for moving metric data around. It can serve as a
collector of data, but is primarily designed to be a framework for building
bridges between data collectors and storage engines.

This repository contains the Skogul library/package, and ``cmd/skogul``,
which parses a JSON-config to set up Skogul.

A copy of the auto-generated manual for skogul is also provided, which is
aimed at end-users. See ``skogul.rst`` (or ``man ./skogul.1``).

.. contents:: Table of contents
   :depth: 2
   :local:

Quickstart - RPM
----------------

If you're on CentOS/RHEL 6 or newer, you should use our naively built RPM,
available at https://github.com/telenornms/skogul/releases/latest.

It will install skogul, set up a systemd service and install a simple
configuration in ``/etc/skogul/default.json``.

There's also a 64-bit linux build there, which should work for most
non-RPM-based installations. But we make no attempt to really maintain it.

Quickstart - source
-------------------

Building from source is not difficult. First you need Golang. Get it at 
https://golang.org/dl/ (I think you want go 1.13 or newer).

Building ``skogul``, including cloning::

   $ git clone https://github.com/telenornms/skogul
   (...)
   $ make
   $ make install

You don't have to use make - there's very little magic involved for regular
building, but it will ensure ``-version`` works, along with building the
auto-generated documentation.

Running ``make install`` installs the binary and default configuration, but
does NOT install any systemd unit or similar.

Also see ``make help`` for other make targets.

About
-----

Skogul is written to solve a myriad of issues that typically arise when
dealing with metric data and complex systems. It can be used for very
simple setups, and expanded to large, multi-datacenter infrastructures with
a mixture of new and old systems attached to it.

To accomplish this, you set up chains that define how data is received, how
it is treated, where it goes and what happens if something goes wrong.

A Skogul chain is built from one or more independent receivers which
receive data and pass it on to a sender. A sender can either transmit data
to an external source (including an other Skogul instance), or add some
internal routing logic before passing it on to one or more other senders.

.. image:: docs/imgs/basic.png

Unlike most APIs or collectors of metrics, Skogul does NOT have a
preference when it comes to storage engine. It is explicitly designed to
disconnect the task of how data is collected from how it is stored.

The rationale is that the problem of writing an efficient snmp collector
should not be tightly coupled to where you store the data. And where you
store the data should not be tightly coupled with how you receive it, or
what you do with it.

This enables an organization to gradually shift from older to newer stacks,
as Skogul can both receive data on old and new transport mechanisms,
and store it both in new and old systems. That way, older collectors can
continue working how they always how worked, but send data to Skogul.
During testing/maturing, Skogul can store the data in both legacy system
and replacement system. When the legacy system is removed, no change is
needed on the side of the collector.

Extra care has been put into making it trivial to write senders and
receivers. For example, an author of a new sender only has to add tags
to their data structure to have that exposed as documentation.

See the package documentation over at godoc for development-related
documentation: 
https://godoc.org/github.com/telenornms/skogul

End-user documentation is found in the manual page, which Skogul can
generate on demand, or you can review a copy on github: 
https://github.com/telenornms/skogul/blob/primary/skogul.rst

More discussion on architecture can be found in `docs/`.

Performance
-----------

Skogul is meant to scale well. Early tests on a laptop proved to work very
well:

.. image:: docs/imgs/skogul-rates.png

The above graph is from a very simple test on a laptop (with a quad core
i7), using the provided tester to write data to influxdb. It demonstrates
that despite well-known weaknesses at the time (specially in the
influx-writer), we're able to push roughly 600-800k values/s through
Skogul. This has since been exceeded.

The laptop in question was using about 150-190% CPU for skogul and 400% for
InfluxDB, the rest went to the testers. No real attempt at tuning was done,
but a few different scenarios were tested.

Note that the general values/s is decent both with a ton of values for each
metric, and just a handful of values per metric, but tons of metrics per
containers.

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

A few sources for more documentation:

- docs/CODE_OF_CONDUCT.md
- docs/CONTRIBUTING
- docs/CODING
- doc.go

Testing
.......

To run test cases, ``go test`` can be run. This can be used either in
individual directories, or at the top directory, with ``go test -short ./...``
(note the triple dots. This is a go-ism for recursive behavior). Tests are
run automatically if you create a pull request.

The ``-short`` argument disables integration tests that would otherwise
fail unless you've set up a compatible postgres and mysql database locally.

To produce coverage analysis, use::

   $ cd skogul
   $ go test -short ./... -covermode=count -coverprofile=coverage.out
   $ go tool cover -html coverage.out
   // Opens a browser with coverage anlysis

Tests are extracted from ``*_test.go`` files, and start with the name
``Test`` followed by a function or data structure, optionally followed by
an underscore and an arbitrary name to support multiple tests of the same
function/type. E.g. ``TestValidate()``, ``TestHTTP_foobar()`` etc.

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

Documentation comes in two forms. One is aimed at end-users. This is
provided mainly by adding proper labels to your data structures (see any
sender or receiver implementation), and through hard-coded text found in
``cmd/skogul/main.go``. In addition to this, stand-alone examples of setups
are provided in the ``examples/`` directory.

For development, documentation is written and maintained using code
comments and runnable examples, following the ``godoc`` approach. Some
architecture comments are kept in ``docs/``, but by and large,
documentation should be consumed from godoc.

See https://godoc.org/github.com/telenornms/skogul for the online
version, or use ``go doc github.com/telenornms/skogul`` or similar,
as you would any other go package.

Examples are part of the test suite and thus extracted from ``*_test.go``
where applicable.

Roadmap
-------

We are doing frequent releases on github, with an ambition of reaching a
1.0 version within some reasonable time frame, I'm guessing 2020. It
doesn't really mean much.

Short term work is defined in milestones on github.

Overall, the core modules and the scaffolding is getting pretty good. The
new config engine is still receiving period updates, but the actual
configuration hasn't changed much.

Future work to get us to 1.0 will be rounding out the new logrus-based
logging by both rewriting the log receiver and overhauling each module to
make our approach to logging consistent across all modules.

Similarly, test cases need to be refreshed. Tests are written very
isolated, and a good bit of spaghetti-logic has arisen. We have decent
coverage, but it's getting trickier to scale test case writing.

Other than that, there are modules to be written and features to be added
which are mostly a matter of what needs arise.
