Templating
==========

The core skogul format, usually represented as JSON, supports the concept
of a template metric. This is intended to reduce the amount of duplicated
fields when a number of fields would otherwise be identical across multiple
metrics.

Typically, this can be:

- hostname
- timestamp
- software version
etc

To enable this, your receiver/handler needs to use the "template"
transformer.

Without templating, sending multiple metrics might look like this::

   {
      "metrics": [
         {
            "timestamp": "2019-03-14T10:00:00+0100",
            "metadata": {
               "poller": "ping-boxen1.example.com",
               "switch": "e13-1",
               "interface": "eth0"
            },
            "data": {
               "in_octets": 151,
               "out_octets": 1112
            }
         },
         {
            "timestamp": "2019-03-14T10:00:00+0100",
            "metadata": {
               "poller": "ping-boxen1.example.com",
               "switch": "e13-1",
               "interface": "eth1"
            },
            "data": {
               "in_octets": 123,
               "out_octets": 456
            }
         },
         {
            "timestamp": "2019-03-14T10:00:00+0100",
            "metadata": {
               "poller": "ping-boxen1.example.com",
               "switch": "e13-1",
               "interface": "eth2"
            },
            "data": {
               "in_octets": 420,
               "out_octets": 1337
            }
         },
         {
            "timestamp": "2019-03-14T10:00:00+0100",
            "metadata": {
               "poller": "ping-boxen1.example.com",
               "switch": "e13-1",
               "interface": "eth3"
            },
            "data": {
               "in_octets": 15,
               "out_octets": 112
            }
         }
      ]
   }

This is valid and fine, however, the more interfaces you have, the more
times timestamp, poller and switch is repeated.

With the template transformer enabled, you can send::

   {
      "template": {
            "timestamp": "2019-03-14T10:00:00+0100",
            "metadata": {
               "poller": "ping-boxen1.example.com",
               "switch": "e13-1",
            }
      },
      "metrics": [
         {
            "metadata": {
               "interface": "eth0"
            },
            "data": {
               "in_octets": 151,
               "out_octets": 1112
            }
         },
         {
            "metadata": {
               "interface": "eth1"
            },
            "data": {
               "in_octets": 123,
               "out_octets": 456
            }
         },
         {
            "metadata": {
               "interface": "eth2"
            },
            "data": {
               "in_octets": 420,
               "out_octets": 1337
            }
         },
         {
            "metadata": {
               "interface": "eth3"
            },
            "data": {
               "in_octets": 15,
               "out_octets": 112
            }
         }
      ]
   }

This works for ALL items, including timestamp, metadata and data.

It will also work with diverse metrics, there's no reason they have to look
similar.

One note: nested objects will NOT be merged, but overwritten. E.g.::

   {
      "template": {
            "timestamp": "2019-03-14T10:00:00+0100",
            "metadata": {
               "poller": "ping-boxen1.example.com",
               "geo": {
                       "country": "norway",
                       "customer": "someone"
               },
               "switch": "e13-1",
            }
      },
      "metrics": [
         {
            "metadata": {
               "interface": "eth0",
               "geo": {
                       "customer_pop": "hamar"
               }
            },
            "data": {
               "in_octets": 151,
               "out_octets": 1112
            }
         }
   }

Will not blend the two "geo" objects together. Instead, the template will
overwrite the metrics.


