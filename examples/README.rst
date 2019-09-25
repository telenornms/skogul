=============================
Skogul configuration examples
=============================

This directory contains various Skogul configuration examples. Some are
provided to demonstrate core Skogul functionality, some to serve as a
template which can be used almost verbatim for production.

http_to_influx.json
===================

Sets up a single receiver on HTTP that uses a handler that parses input as
json, executes the templater transformer, then executes two different
senders in series.

The batch sender, which is the first in the chain, simply waits until it
has a reasonable amount of metrics before it passes them on. If it does not
have the desired quantity, it will still forward the metrics after an
interval has expired.

The batch sender will allow you to protect your databases. Most databases
(and storage systems in general) are more efficient if you can forward more
data at the same time, instead of many small chunks.

The next sender is the influxdb-sender, which writes to influxdb over HTTP
and the line protocol of InfluxDB.

https_basic_auth_count.json
---------------------------

Demonstrates basic authentication with TLS, sends data to a batch sender,
then counter, then discards it. The counter periodically emits data to the
"debugger" handler - which uses the "print" sender to simply output the
data.

In other words: Accepts data on basic-auth protected HTTPS and emits
performance statistics without ever storing the data.

https_influx.json
-----------------

Accepts data on HTTPS, protected by basic authentication and TLS, collects
up to 1000 metrics for up to 5 seconds, then writes the data to InfluxDB.

tester_to_http.json
-------------------

Generates test-data - 500 metrics, each with 10 key/values. Then batches up
to 1000 metrics for up to 5 seconds. The backoff sender will send data
immediately, but if it fails, it will wait for 100ms, then 200ms, then
400ms, etc, for up to 10 attempts. Data is sent to a basic auth-protected
HTTPS endpoint, but certificates are not verified.
