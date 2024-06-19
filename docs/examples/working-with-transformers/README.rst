Developing transformer configuration
====================================

When dealing with input data, the relatively easy part is getting the data,
but you often have to transform it significantly before it is suited for
storage. An example of this is juniper telemetry data, but the same
fundamental method applies to any semi-complex data, e.g.: if you have an
application that outputs statistics as a nested data structure.

The general approach to this is the same:

1. Set up a receiver for the appropriate data. This depends on the source.
   Use a handler that simply prints to standard out to verify that your
   receiver work. E.g.: the "print" sender. See the basics examples for
   several examples of this.
2. Once you see data on standard output, switch form the "print" sender to
   a file sender, use the "skogul" encoder. This will let you write json
   data to disk, which is easier to work with.
3. Once you have a sufficiently large data set, you can stop skogul for
   now, and start reading from disk instead.

This directory contains an example payload (``payload.json``) from Juniper
telemetry which is fairly nested.

The ``step1.json`` just reads that file and prints it. Which is the first
step. Then we start applying various transformers.

In ``step2.json`` we demonstrate basic splitting per interface and how to
extract ``if_name`` from the data fields to metadata. In ``step3.json``, we
further split per lane.
