Junos Telemetry and Skogul
==========================

To get Juniper's streaming telemetry working with Skogul you need a Junos
device that supports Telemetry, and some relevant sensors. The matrix of
what devices support what sensors is beyond the scope of this document, but
the following is a snippet of junos config that works on MX960 with Junos
18.4 at the very least (and Junos17 / MX480).

This document is meant to provide a starting point, and is NOT meant as
reference documentation. Please see the official Juniper documentation for
that.

Junos config
------------

::

   > show configuration services analytics
   streaming-server telemetry_target1 {
       remote-address 192.168.1.10;
       remote-port 1234;
   }
   streaming-server telemetry_target2 {
       remote-address 172.16.0.10;
       remote-port 1234;
   }
   export-profile export_often {
       local-address 192.168.5.5;
       local-port 20006;
       reporting-rate 6;
       payload-size 3000;
       format gpb;
       transport udp;
   }

   export-profile export_rarely {
       local-address 192.168.5.5;
       local-port 20030;
       reporting-rate 300;
       payload-size 3000;
       format gpb;
       transport udp;
   }
   sensor junos_system_linecard_intf-exp {
       server-name [ telemetry_target1 telemetry_target2 ];
       export-name export_often;
       resource /junos/system/linecard/intf-exp/;
   }
   sensor junos_system_linecard_interface {
       server-name [ telemetry_target1 telemetry_target2 ];
       export-name export_often;
       resource /junos/system/linecard/interface/;
   }
   sensor junos_system_linecard_optics {
       server-name [ telemetry_target1 telemetry_target2 ];
       export-name export_rarely;
       resource /junos/system/linecard/optics/;
   }

`streaming-server` defines the servers to receive the traffic, two
different servers in this example.

`export-profile` defines how often to send data (reporting-rate), and
payload-size. You probably don't need local-address/local-port or
payload-size. Payload-size is relevant if your device uses jumboframes but
your server does not. (Note that payload-size is currently not directly
mapped to actual UDP payload, so payload-size of 3000 on junos 18.4 results
in roughly udp packages of 1400 bytes).

`sensor` is device specific and sets up which sensors are active and send
telemetry.

Skogul config
-------------

Basic reception of Telemetry is very straight forward::

   {
     "receivers": {
         "udp": {
           "type": "udp",
           "address": ":1234",
           "handler": "protobuf",
           "packetsize": 9000
         }
     },
     "handlers": {
       "protobuf": {
         "parser": "protobuf",
         "transformers": [],
         "sender": "debug"
       }
     },
     "senders": {
       "debug": {
           "type": "debug"
       }
     }
   }

This will receive UDP-based telemetry on port 1234 and print the parsed
result.

Since printing to stdout is less useful, here's an example that also uses a
set of transformers to "flatten" the data and send it to influxdb::

  {
    "receivers": {
      "udp": {
        "type": "udp",
          "address": ":1234",
          "handler": "protobuf",
          "packetsize": 9000
      }
    },
    "handlers": {
      "protobuf": {
        "parser": "protobuf",
        "transformers": [
          "interfaceexp_stats",
          "interface_stats",
          "interface_info",
          "optics_diag",
          "flatten_egress_queue_info",
          "remove_egress_queue_info",
          "extract_if_name",
          "extract_measurement_name",
          "extract_measurement_name2",
          "flatten_systemId"
        ],
        "sender": "batch"
      }
    },
    "senders": {
      "batch": {
        "type": "batch",
        "interval": "2s",
        "threshold": 1000,
        "next": "influx"
      },
      "influx": {
        "type": "influx",
        "measurementfrommetadata": "measurement",
        "URL": "https://localhost:8086/write?db=skogul",
        "Timeout": "10s"
      }
    },
    "transformers": {
      "interfaceexp_stats": {
        "type": "split",
        "field": ["interfaceExp_stats"]
      },
      "interface_stats": {
        "type": "split",
        "field": ["interface_stats"]
      },
      "interface_info": {
        "type": "split",
        "field": ["interface_info"]
      },
      "optics_diag": {
        "type": "split",
        "field": ["Optics_diag"]
      },
      "flatten_egress_queue_info": {
        "type": "data",
        "flatten": [
          ["egress_queue_info"],
          ["egress_errors"],
          ["egress_stats"],
          ["ingress_errors"],
          ["ingress_stats"],
          ["egress_stats"],
          ["ingress_stats"],
          ["op_state"],
          ["optics_diag_stats"],
          ["optics_diag_stats__optics_lane_diag_stats"]
        ]
      },
      "remove_egress_queue_info": {
        "type": "data",
        "remove": [
          "egress_queue_info",
          "egress_errors",
          "egress_stats",
          "ingress_errors",
          "ingress_stats",
          "egress_stats",
          "ingress_stats",
          "op_state",
          "optics_diag_stats",
          "optics_diag_stats__optics_lane_diag_stats"
        ]
      },
      "extract_if_name": {
        "type": "metadata",
        "extractFromData": ["if_name"]
      },
      "extract_measurement_name": {
        "type": "replace",
        "regex": "^([^:]*):/([^:]*)/:.*$",
        "source": "sensorName",
        "destination": "measurement",
        "replacement": "$2"
      },  
      "extract_measurement_name2": {
        "type": "replace",
        "regex": "[^a-zA-Z]",
        "source": "measurement",
        "replacement": "_" 
      },
      "flatten_systemId": {
        "type": "replace",
        "regex": "-.*$",
        "source": "systemId"
      }
    }
  }

This seems like a big config, but it needs to be split up.

First, just one receiver: UDP. And one handler: protobuf. Two senders: The
batch sender to batch up data for 2 seconds before dumping it to influxdb.

The real "magic" is in the transformers. They are all executed in the order
specified. There are a few different types:

- The various "split" transformers iterate over an array of data and create
  individual metrics for each item. E.g.: One metric for each interface,
  instead of a single metric for all interfaces.
- "flatten_egress_queue_info" will "flatten" a nested data structure - e.g,
  instead of `` "foo": { "bar": x, "baz": y }``, you get `` "foo__bar": x,
  "foo__baz": y``. This tends to create a good bit of columns for things
  like queue stats (E.g.: one column for each stat for each queue), but
  works fine for things like egress_stats and ingress_stats. You could also
  use a second "split" for that.
- "remove_egress_queue_info" removes the fields that were just flattened to
  avoid duplicates.
- "extract_if_name" extracts the "if_name" field from "data" into
  "metadata". E.g.: Instead of being any value, it becomes a tag in
  influxdb which you can search for.
- "extract_measurement_name" uses the sensorName provided and extracts the
  actual sensor, then removes illegal characters. This is so we can get
  individual measurements for each sensor in influxdb.
- Lastly, "flatten_systemId" will clean up system names like
  "foo100b45r-re0:192.168.1.5" to just "foo100b45r" - this one might
  require local modifications.

Please note: This configuration is more complicated than it absolutely has
to be, but should provide a good example of a real-world use case.
