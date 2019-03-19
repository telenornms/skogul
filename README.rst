======================================
Skogul - generic metric/data collector
======================================

Skogul is a collector of various data, but unlike most similar tools, it
does NOT provide a backend storage, and the aim is to provide a generic,
reasonably future-proof API that can be built upon using various backend
writers, transformers and more.

It is very much a work in progress, and is aimed to handle both simple
installations where there's roughly a 1:1 between input sources and storage
backends, and large enterprise installations where there can be hundreds of
different input sources all being routed and transformed based on
site-specific needs.

The first use-case is expected to be Gondul
(https://github.com/gathering/gondul), where it will provide a shim-layer
between SNMP collectors, ping collectors and DHCP event data collectors;
and postgresql and influxdb as storage backends.

The rationale is that the problem of writing an efficient snmp collector
should not be tightly coupled to where you store the data. And where you
store the data should not be tightly coupled with how you receive it, or
what you do with it.

At present time, it's not suited for much more than looking at the general
development of the architecture. As such, build-instructions and more are
explicitly left out.

Name
----

Skogul is a Valkyrie. After extensive research (5 minutes on wikipedia with
a cross-check on duckduckgo), this name was selected because it is
reasonably unique and is also a valkyrie, like Gondul, a sister-project.

Being a whole week old, a name change was due, so you might find references
to Gollector here and there, if I suck at grep(1).

