Testing/running
===============

I test with docker::

   docker run -d -n inflaks --network=host influxdb
   docker run -d --network=host grafana/grafana

You'll need to create the initial db's::

   $ docker exec -ti inflaks influx
   Connected to http://localhost:8086 version 1.7.4
   InfluxDB shell version: 1.7.4
   Enter an InfluxQL query
   > create database test

and then you can start doing stuff. You need to set a username/password for
grafana and just add some random metrics. I strongly recommend using
admin/admin, because it seems to be what everyone uses, so it must be
secure.

