API
===

The general API for the Skogul framework is aimed at being a
general-purpose "data"-moving API that can transform data, but also has to
be reasonably fast.

I intend to version the API.

Structure
---------

There are two parts to the API:

1. The container-format, which is generic, JSON-based (with a
   go-representation) and used for ALL endpoints.
2. Endpoints - if you write to /snmp you will expect different results
   than writing to /ping or /generic.

The latter part is not part of the spec except a single generic variant,
and while the skogul-code will probably provide multiple implementations
for convenience, the actual endpoints are considered site-specific.

This document deals with the general aspects of the API.

Container
---------

Root container::

   {
      "template": metric-object
      "metrics": array-of-metric-objects
   }

The metric object::

   {
      "timestamp": "RFC3339/ISO8601-based time",
      "metadata": {
         "key": "value",
         ....
      },
      "data": {
         "key": "value",
         ...
      }
   }

Metadata are identifying variables for the metrics, e.g.: Expected to be
indexed and not vary TOO much. Examples are server names, port names, e.g.

The template in the root-container is completely optional, and is provided
to allow you to post hundreds of metrics which all share a set of common
attributes.

Full example without template, single metric::

   {
      "metrics": [
         {
            "timestamp": "2019-03-14T10:00:00+0100",
            "metadata": {
               "switch": "e13-1",
               "source": "ping-boxen1.example.com"
            },
            "data": {
               "latency4": 1.51,
               "latency6": 0.1112
            }
         }
      ]
   }

A similar example, using a template and more metrics::

   {
      "template": {
         "timestamp": "2019-03-14T10:00:00+0100",
         "metadata": {
               "source": "ping-boxen1.example.com"
         }
      },
      "metrics": [
         {
            "metadata": { "switch": "e13-1", },
            "data": { "latency4": 1.51, "latency6": 0.1112 }
         },
         {
            "metadata": { "switch": "e13-2", },
            "data": { "latency4": 1.11, "latency6": 0.1312 }
         },
         {
            "metadata": { "switch": "e13-3", },
            "data": { "latency4": 0.51, "latency6": 0.1512 }
         },
         {
            "metadata": { "switch": "e13-4", },
            "data": { "latency4": 1.91, "latency6": 0.1912 }
         }
      ]
   }

While it is possible to write a more terse API, readability is prioritized
and frankly, if we need to write shorter variable names to get performance,
we're doing it wrong anyway.

Details
=======

Depth
-----

While the container format does NOT put a limit on the depth of the
metadata- and data- objects, using more than a single level is UNDEFINED.

E.g. the following is allowed, but the API does not guarantee that it will
work::

   [
      {
         "metadata": {
            "switch": "r1.noc"
         },
         "data": {
            "ae0": {
               "ifInOctets": 515
            },
            "ae2": {
               "ifInOctets": 525
            }
         }
      }
   ]

However, the following IS defined and MUST store all data fields::
   
   [
      {
         "metadata": {
            "switch": "r1.noc",
            "port": "ae0"
         },
         "data": {
            "ifInOctets": 515
         }
      },
      {
         "metadata": {
            "switch": "r1.noc",
            "port": "ae2"
         },
         "data": {
            "ifInOctets": 525
         }
      }
   ]

This is explicitly kept as such because it allows the API to account for
endpoints that accept the former output and produces the second.

The generic APIs will only work on a single level of depth.



