#!/bin/bash

target=${1:-enrich_large.json}
metrics=${2:-1000000}

echo Writing $metrics metrics to $target

echo '[' > $target
i=0;
comma=""
while [ $i -lt $metrics ]; do
	i=$(( i + 1 ))
	cat >> $target <<-_EOF_
	$comma
	{ "metadata": { "key1": $i }, "data": { "who": "I am $i", "serial": $RANDOM } }
	_EOF_
	comma=","
done
echo ']' >> $target
