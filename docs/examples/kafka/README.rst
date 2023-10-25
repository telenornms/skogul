Kafka examples
==============

Warning: The Kafka sender has been used extensively, but not on a diverse
set of Kafka buses. The receiver is much less tested.

kafka_to_stdout.json
--------------------

Connects to a Kafka bus as a consumer and prints data to stdout, note the
parser needs to be "json1", not "json".

tester_to_kafka.json
--------------------

Connects to a Kafka bus as a producer and posts data, each metric is
encoded and sent as an independent Kafka message.

http_to_kafka.json
------------------

Listens for HTTP requests and forwards it to a Kafka bus. Very useful for
connecting applications that can't easily implement/maintain a Kafka
connection with Kafka.
