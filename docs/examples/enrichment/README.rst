Enrichment
==========

Enrichment is the process of adding more metadata to incoming data to make
it more useful. This can be adding the customer name to an interface,
adding site ID, mapping serial numbers to customer, and so forth.

Think of it is as doing an in-line SQL join.

To accomplish this, Skogul has a very simple Enrichment transformer, and a
corresponding "enrichmentupdater" sender.

For skogul to enrich data, the sort data must be expressed as a typical
Skogul container, so we can leverage the rest of Skogul's mechanisms.

In other words::

        {
                "metrics": [
                        {
                                "metadata": {
                                        "sysName": "foobar",
                                        "ifName": "eth0"
                                },
                                "data": {
                                        "customer": "something"
                                }
                        },
                        {
                                "metadata": {
                                        "sysName": "foobarx",
                                        "ifName": "eth2"
                                },
                                "data": {
                                        "customer": "blatti"
                                }
                        }
                ]
        }

This data tells Skogul to enrich metrics where the fields "sysName" and
"ifName" match, and ADD the metadata "customer". In other words, it ALL
works on metadata in the end but:

- Enrichment fields listed in "metadata" are what is matched for in metrics
- Enrichment fields listed in "data" are added to the metric

To use this data, you need to set up a transformer, in the example, it's
called "someEnricher". The transformer needs to know which fields to look
up: The example uses the test-sender, which has "key1", for the above
example data, the list would instead read ``[ "sysName", "ifName" ]``.

To actually load data, you need to set up a second (or third)
receiver-pipeline to load the data, then pass it on to the
"enrichmentupdater" sender. The enrichment updater sender needs to know
which transformer to update, and that's it.

This will let you load enrichment information from any type of source for
which there is a receiver. E.g.: From file, from a HTTP end-point, from a
bus.

The example has set up two methods:

1. The "bootstrap" receiver will read JSON data from a file on boot, once.
2. The "updater" receiver is a regular HTTP endpoint so you can POST
   updates to Skogul live. Updates are incremental.

*One small note*: The example also uses the "now" transformer. This just
adds a time stamp to the incoming data, so you don't have to provide a
timestamp, since Skogul will refuse to validate incoming data without
timestamps.

Feel free to test it - after startup you can modify the payload data in
docs/examples/payloads/enrichment.json and use HTTP POST to localhost to
see the updates take effect live.
