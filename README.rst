
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

This repository contains the Skogul library/package, and ``skogul-x2y``, a
binary that leverages the Skogul library to provide a bridge from X to Y,
e.g. from a HTTP receiver to an InfluxDB sender. ``skogul-x2y`` supports
most Skogul receivers and senders, but more complex chains can also be
created.

.. contents:: Table of contents
   :depth: 2
   :local:

Quickstart
----------

You need to install a recent/decent version of Go. Either from your
favorite Linux distro, or through https://golang.org/dl/ .

Building ``skogul-x2y``, including cloning::

   $ git clone https://github.com/KristianLyng/skogul
   (...)
   $ cd skogul/cmd/skogul-x2y
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

   $ cd skogul/cmd/skogul-x2y
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

The simplest use of Skogul is to use the ``cmd/skogul-x2y`` package, which
provides *limited* support for a number of "senders" and "receivers", which
can be arbitrarily matched. This will allow you to receive Skogul-formated
JSON data on HTTP, MQTT, local fifo, line-based TCP sockets and possibly
other sources in the (near) future, and pass them on to an other Skogul
instance over http, to influxdb, M&Rm or post it on a MQTT bus.

An example of the help-screen of ``skogul-x2y`` gives an idea of what you
can use it for::

   Usage of cmd/skogul-x2y/skogul-x2y:
     -help
           Print extensive help/usage
     -receiver string
           Where to receive data from. See -help for details. (default "http://[::1]:8080")
     -sender string
           Where to send data. See -help for details. (default "debug://")

   skogul-x2y sets up a skogul receiver, accepts data from it and passes it to the sender.

   Available senders:
     scheme:// | Description
   ------------+------------
        mnr:// | MNR sender sends M&R line format to an endpoint, optional DefaultGroup
               | is provided as the path element.
       mqtt:// | MQTT sender publishes received metrics to an MQTT broker/topic
      debug:// | Debug sender prints received metrics to stdout
       http:// | Post Skogul-formatted JSON to a HTTP endpoint
     influx:// | Send InfluxDB data to a HTTP endpoint, using the first element of the
               | path as db and second as measurement, e.g:
               | influx://host/db/measurement


   Available receivers:
     scheme:// | Description
   ------------+------------
       http:// | Listen for Skogul-formatted JSON on a HTTP endpoint
       fifo:// | Read from a FIFO on disk, reading one Skogul-formatted JSON per line.
               | fifo:///var/skogul/foo
       mqtt:// | Listen for Skogul-formatted JSON on a MQTT endpoint
        tcp:// | Listen for Skogul-formatted JSON on a line-separate tcp socket
       test:// | Generate dummy-data, each container contains $m metrics and each
               | metric $v values, format: test://$m/$v

skogul-x2y can also be used to test Skogul. Here's a very simple example
where data is moved from one Skogul instance to an other over HTTP, using
the "test receiver" to generate dummy data and the "counter receiver" to
instrument it on the other side. Similar can also be used to pipe data to
influx or M&R or any other sender.

.. image:: docs/self-test.png

While this 1-to-1 scenario is very useful and common, it is not really
where Skogul shines the most. The core idea behind Skogul is building
pipelines that starts with one or more receiver and builds a chain of
multiple senders. Each sender comes in one of two forms: largely "internal"
senders, and "terminal/external" senders. The latter is the most easily
understood sender: One that transmits the data to an external data source -
presumably for permanent storage. The internal sender will allow such
things as duplicating a metric to multiple other senders (e.g.: Send the
data to both influx and postgres), try sending first to one sender, then if
that fails, push to an other (e.g.: fallback / ha), and so on.

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

At present, the big missing piece is decent TLS-support. This is not
particularly hard code-wise, since golang provides all the support we need,
but has been delayed mainly for UX reasons: Experimentation with how to
provide options both for skogul-x2y and within the library in a manner
where each receiver/sender only require a minimal amount of boilerplate
code has taken some precedence.

Other than that, there are no huge plans beyond what the issues denote. New
receivers and senders will be added as needed (patches are welcome), and
the code will be maintained as it's more widely deployed.
