#!/bin/bash
set -e

GIT_DESCRIBE="$(git describe --always --tag --dirty)"
OS=$(uname -s | tr A-Z a-z)
ARCH=$(uname -m)
V=${1:-unknown}
if [ "x${V}" = "xunknown" ] && [ "x${GIT_DESCRIBE}" != "x" ]; then
	V=${GIT_DESCRIBE}
fi

rm -fr dist/
mkdir -p dist/share/doc/skogul
mkdir -p dist/bin
mkdir -p dist/share/man/man1
mkdir -p dist/share/doc/skogul

go build ./cmd/skogul
./skogul -make-man > dist/share/doc/skogul/skogul.rst
rst2man < dist/share/doc/skogul/skogul.rst > dist/share/man/man1/skogul.1
cp skogul dist/bin
cp LICENSE dist/share/doc/skogul
cd dist
tar cvf ../skogul-${V}.${OS}-${ARCH}.tar.bz2 .
