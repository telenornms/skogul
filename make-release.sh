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
cp -r README.rst docs/* dist/share/doc/skogul/

awk < docs/NEWS -v ver="$V" '
$0 == ver {
	here=1
	herel=length($0)
	next
}
here == 1 && herel == length($0) && /^=+$/ {
	here = 2
	next
}
here == 1 {
	here=0
}

/^v[0-9]+\.[0-9]+\.[0-9]+[0-9a-z]*$/ {
	here = 0
}

here == 2 && /.+/ {
	lines++
}

here == 2 && lines
' > notes
cd dist
tar cvjf ../skogul-${V}.${OS}-${ARCH}.tar.bz2 .
