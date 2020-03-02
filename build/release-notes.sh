#!/bin/bash

set -u
set -e
test -e docs/NEWS || {
	echo >&2 "Can't find NEWS. Must run from the top-directory."
	exit 1
}

awk < docs/NEWS -v ver="$1" '
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
'
