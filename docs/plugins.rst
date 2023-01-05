Plugins
=======

Plugins are a new concept that is currently going through testing.

See https://github.com/telenornms/skogul-plugin-example for an example
plugin demonstrating how to use this new support.

The basic idea is simple and solid: Allow loading out-of-tree modules so
complex modules with potentially problematic dependency chains doesn't have
to be shipped together with "mainline" Skogul.

The basic module design of Skogul supports this with practically zero
changes, so the core design shouldn't change much, however there are
several minor, but important details to iron out, mainly related to
release-maintenance.

Expect changes in how the repository is organized. Mainly:

1. Expect the main skogul repository to SHRINK in size as modules such as
   juniper telemetry support is moved out.
2. Expect the main RELEASES to include code from more than one package so
   end-user-experience remains the same.
3. Expect the actual plugin loading to be actually intelligent, not just a
   static list of modules to load.
4. More focus on API stability.

This document will be updated as the API for plugins are ironed out.

