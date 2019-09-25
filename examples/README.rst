=============================
Skogul configuration examples
=============================


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
