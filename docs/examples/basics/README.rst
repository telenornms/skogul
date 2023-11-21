Basic Skogul examples explained
===============================

This directory has a number of fundamental building blocks to help you
understand Skogul. They are useful to read through to quickly grasp the
basic concepts.

This file is tries to walk you through them in incremental steps. Gradually
less prose is provided.

default.json
------------

This is the default installed configuration file. It contains a single
receiver and handler.

It listens for HTTP on localhost (IPv6).

There are two named objects here, the receiver named "api", which is of
type "http". You can look up the HTTP receiver in the manual file to find
all configuration options.

The HTTP receiver has a map of URLs to handlers - you can use a single HTTP
receiver to receive data on different URLs and treat them differently. In
this case, only "/" is defined, and anything received on
http://localhost:8080/ is sent to the "myhandler" handler.

A handler defines how data is parsed - in this case, using the native
"skogul" parser which expects Skogul-formatted JSON data. No transformers
are listed, but a single sender is listed: "print".

There are a number of modules with few or zero options where you don't have
to define them. The "print"-sender is one such sender. If you look it up in
the manual you will see the "debug" sender, which is an other name for the
same. This is to save you from having to add a configuration like this::

        "senders": {
                "print": {
                        "type": "print"
                }
        }

You can still DO that though! If you were to define that print should do
something else, your configuration will take precedence.

tester_to_stdout.json
---------------------

This is a dummy-configuration which only serves to test that Skogul works
at all. The "tester" receiver generates random data for you and is very
useful for testing if your senders work, or during development. In this
example, it just prints the result to standard out. You can run it with::

        ./skogul -f tester_to_stdout.json

And see for yourself. Both the "print" sender and "test" receiver are
mainly useful for debugging and testing.

forward_and_dupe.json
---------------------

This is a demonstration of how you can use a dupe sender to debug a stream
of data. The dupe sender can be inserted anywhere.

In addition it demonstrates that two receivers can use the same handlers:
Both the test-receiver that automatically generates dummy-data and the
"api"-receiver that listens for data on HTTP use the handler called
"myhandler", and thus treat the data the same.

Two different methods of debugging is provided here as well: By using the
"print" (or "debug") sender, data is printed to standard output (e.g.:
journalctl). Because the print-sender doesn't require any configuration,
you don't have to define it anywhere either, though you can.

The "to-file" sender simply writes the data to a file, using the provided
encoder. "skogul" means this is encoded as json, with Skogul's format.

http_count.json
---------------

This introduces the count-sender, which does exactly what the name
suggests: it counts whatever passes through it. It is highly useful for
record-keeping and for debugging larger quantities of data.

To use the count sender you need at least two handlers: One regular
handler, and one to use for sending the statistics to. You can re-use the
same handler if you like, it doesn't create any loops since it'll just
count its own statistics.

We also use the "null"-sender here, which just discards data. The null
sender is useful whenever you need to specify a sender, but don't want to
keep the data.

https_basic_auth_count.json
---------------------------

Introduces authentication and certificates of the HTTP receiver, otherwise
much the same. Worth noting that you must provide a username and password
for each path you use. It's slightly cumbersome, but works well.

Skogul will warn you if you try to use basic authentication without TLS
encryption.

batch_simple.json
-----------------

This demonstrates the batch-sender. The test-receiver generates a single
metric every 0.4 seconds, the batch sender will try to collect up to 100
metrics before forwarding them, but wait a maximum of 1s. The effect is
that you will see containers printed to stdout with 2 and 3 metrics.

batch_burner.json
-----------------

This demonstrates how to deal with delays in the sender-path. This
is simulated with the "sleep" sender, which introduces a variable delay
before passing the data on - in this case, it passes it to the "null"
sender which just discards the data for the sake of demonstration. Think of
this as simulating a slow database or storage backend.

Additionally, two count-senders are used: One after the sleep-sender, to
measure the throughput of "good" data, and one to measure "burned" data.

Data "burning" is a feature of the batch sender. The batch sender can
smooth out variable response times, but it has a finite backlog. If this
backlog is full and no burner is defined, the batch sender will block,
which means data will pile up and possibly slow down your receivers as
well.

With a burner defined, you can chose an alternate path for this data. The
most useful thing to do with the data is often to just throw it away -
using the null-sender, but ours is a very common pattern: Send it to a
count-sender first to gather statistics on how often this happens.

The count sender identifies itself in metadata, so you can distinguish
between them by looking at the "skogul" metadata field, which will have the
value of the name you give the sender. Here, that will be count-burned and
count-ok.

The test-receiver is used without a delay, which means it will generate
data as fast as it can and quickly overwhelm the "sleep"-sender.

It is highly recommended to address what to do if data is coming too fast
explicitly, specially if you are using multiple senders. It allows you to
fail in a more controlled manner. The simplest way to do this is just add a
burner to a batch sender, with no further options. This is a valid, minimal
configuration doing just that::

        "senders": {
                "mybatch": {
                        "type": "batch",
                        "burner": "null",
                        "next": "next-sender"
                },
                "next-sender": {
                        ....
                }
        }


