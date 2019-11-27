#!/bin/bash

set -e

mkdir -p dist/share/doc/skogul
mkdir -p dist/share/man/man1

./dist/skogul -make-man > dist/share/doc/skogul/skogul.rst
rst2man < dist/share/doc/skogul/skogul.rst > dist/share/man/man1/skogul.1
