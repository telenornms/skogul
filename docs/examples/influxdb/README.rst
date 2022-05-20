Skogul InfluxDB examples
========================

Consider skimming through the basics/ examples first.

tester_to_influxdb.json
-----------------------

This example writes test-data to to the "skogul" database of a local
InfluxDB instance, on the measurement also called "skogul".

http_to_influx.json
-------------------

Accepts unencrypted and unauthenticated data on localhost, batches up to
1000 metrics before writing them to a local InfluxDB, database "testdb",
measurement "demo", with a 10s timeout set.

Batching will wait a maximum of 5 seconds to collect 1000 metrics before
sending. Even batching for 1s is useful and dramatically improves
performance.

The influx sender will send the entire batch as a single HTTP POST.

https_influx.json
-----------------

Same as above, but with TLS certificates and basic authentication
