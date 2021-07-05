Juniper Telemetry and skogul
============================

These two JSON config files can be used together to enable juniper
telemetry parsing on skogul. The files are split in two for your
convenience - the ``junos-telemetry-receiver.json`` specifies what port to
listen on and how - which is assumed to be something you might want to
tweak. The ``junos-telemetry-transformers.json`` specifies the transformers
and handler which you probably don't need to twiddle with.

The sender used is "default" - if you have an empty skogul config that just
means it will be print to stdout, but you are free to change this as need
be.

So to use this you need:

1. Enable telemetry on your router - see later chapter.
2. Copy both files to ``/etc/skogul/conf.d``
3. (Re)start skogul
4. Verify you get telemetry in your skogul logs
5. Redefine the "default" sender to something sensible like influxdb. An
   example that is useful is provided.

Sender example
--------------

An example sender can be::
        
      { "senders": {
        "influx": {
           "type": "influx",
           "measurementfrommetadata": "measurement",
           "URL": "https://your-influx-host:andport/write?db=db",
           "Timeout": "4s"
        }
      }

The important bit here is really the "measurementfrommetadata" since it
will store different sensors' data on different measurements in influxdb.

Example Juniper-config
----------------------

This is JUST an example, you probably want to change it.

::

   streaming-server telemetry_1 {
       remote-address 192.168.0.10;
       remote-port 3300;
   }
   streaming-server telemetry_2 {
       remote-address 192.168.0.20;
       remote-port 3300;
   }
   export-profile export_fast {
       reporting-rate 6;
       payload-size 3000;
       format gpb;
       transport udp;
   }
   export-profile export_medium {
       reporting-rate 30;
       payload-size 3000;
       format gpb;
       transport udp;
   }
   export-profile export_slow {
       reporting-rate 300;
       payload-size 3000;
       format gpb;
       transport udp;
   }

   sensor xx_linecard_intf-exp {
       server-name [ telemetry_1 telemetry_2 ];
       export-name export_medium;
       resource /junos/system/linecard/intf-exp/;
   }
   sensor xx_linecard_optics {
       server-name [ telemetry_1 telemetry_2 ];
       export-name export_slow;
       resource /junos/system/linecard/optics/;
   }
   sensor xx_interfaces_interface {
       server-name [ telemetry_1 telemetry_2 ];
       export-name export_fast;
       resource /junos/system/linecard/interface/;
   }
   sensor xx_interfaces_interface_subinterfaces {
       server-name [ telemetry_1 telemetry_2 ];
       export-name export_medium;
       resource /junos/system/linecard/interface/logical/usage/;
   }

Further reading
---------------

This document and the exmaples are based on a blog post that goes into
further detail:
https://kly.no/posts/2020_01_13_Skogul_Junos_Telemetry_and_InfluxDB.html
