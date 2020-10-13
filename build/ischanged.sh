#!/bin/bash

# Check if a file is changed this year and if the copyright-header reflects
# it. Use in combination with find.

YEAR=$(date +%Y)
if [ $(git log --since=${YEAR}-01-01 --oneline "$1" | wc -l) -gt 0 ]; then
	if grep -q "Copyright.*${YEAR}" $1; then
		if [ ! -z "${VERBOSE}" ]; then
			echo "$1 Updated, but correct header";
		fi
	else
		echo "$1 Updated, missing copyright year"
	fi
else
	if [ ! -z "${VERBOSE}" ]; then
		echo "$1 Not updated"
	fi
fi

if ! grep Copyright $1 -q; then
	echo $1 missing Copyright-header all together
fi
