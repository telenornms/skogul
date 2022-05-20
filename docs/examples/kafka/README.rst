Kafka examples
==============

WARNING: The Kafka implementation is very much "MVP" and hast not _yet_
been tested extensively.

kafka_to_stdout.json
--------------------

Connects to a Kafka bus as a consumer and prints data to stdout, note the
parser needs to be "json1", not "json".

tester_to_kafka.json
--------------------

Connects to a Kafka bus as a producer and posts data, each metric is
encoded and sent as an independent Kafka message.
