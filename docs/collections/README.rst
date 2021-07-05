Skogul config collections
=========================

This directory contains various configuration examples that you can use
more or less verbatim to enable support for various devices/collectors.

Each collection should have its own README.rst explaining the contents and
how to use it.

Ideally you just copy the json-config files to /etc/skogul/conf.d and
things should work immediately - but you may have to change ports and such
to suite your needs.

Sender
------

The sender used for these collections is typically "default" - if you don't
define "default" yourself, it will be the same as logging to stdout which
is nice for testing but unhelpful for production.

It is therefor expected that you define the "default" sender yourself.


