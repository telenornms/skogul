Receiving Juniper Telemetry with Skogul
=======================================

Skogul comes bundled with support for parsing Juniper streaming telemetry.

This directory contains a single Skogul configuration, split into three
separate files. It can be loaded by pointing skogul to the directory,
instead of just a single file.

The RPM installation will handle this as long as files are stored in
``/etc/skogul/conf.d/``.

This is done because it will let you copy the transformers.json verbatim,
update it when the example update, and still provide your own setup for how
to receive the data and what to do with it.

See also https://kly.no/posts/2020_01_13_Skogul_Junos_Telemetry_and_InfluxDB.html

The biggest task is "flattening" the data, as the telemetry data is a
nested tree structure. This example provides a set of transformers that
does this for you.

Configuring Junos
-----------------

Before you start working with Skogul, you need to set up your router.

See ``junos-conf.conf`` for example configuration which sets up 4 sensors,
3 different intervals and sends to two servers. You can obviously simplify
this.

receiver.json
-------------

This defines the receiver, which is of type "udp" and is not very special.
It does a bit of tuning which you probably don't need unless you receive a
significant amount of data.

It then hands things of to the "juniperxform" transformer.

transformers.json
-----------------

Try running with a sender "print" and with and without the transformers and
you will quickly see the point of this.

This is really the meat of the example - the "juniperxform" handler and the
related transformer-configurations. It's somewhat verbose, but flattens the
tree structures Juniper telemetry into something InfluxDB can deal with.
There are a few custom details here that may not be relevant for you: We
remove the trailing "-.*$" from systemId for instance, because we don't
want to have "-re0" and "-re1", we also extract sensor name and make it
more friendly (e.g.: no special characters), and saves it as "measurement".

It's a bit verbose, but each transformer is simple.

Once done, we pass everything off to the "juniper-data" sender.

senders.json
------------

This defines the "juniper-data" sender, mentioned above. Here, it is
InfluxDB, but the whole point of separating this into different files is
that it is up to you where to store this.

Converting ints to floats is done because, well, if you get "0" and treat
it as an integer, then get "0.5" later, Influx will not be happy. Casting
correctly is an alternative, but way more work. We also use
"measurementfrommetadata" to pick InfluxDB measurement based on the
metadata field "measurement" - which is based on the name of the sensor,
e.g.: "junos_system_linecard_interface".

Using it
--------

To use this, configure your router(s), copy all three files to
``/etc/skogul/conf.d/``, adapt ``receiver.json`` and ``senders.json`` to
your liking, then start it up. If we add more transformation (e.g.: for
better performance, or new sensors), you should only need to update
``receiver.json``, not the rest.
