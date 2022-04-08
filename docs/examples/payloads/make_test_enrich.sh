#!/bin/bash

target=${1:-enrich_large.json}
metrics=${2:-10000}

echo Writing $metrics metrics to $target

echo '' > $target
i=0;
comma=""
while [ $i -lt $metrics ]; do
	i=$(( i + 1 ))
	cat >> $target <<-_EOF_
	{ 
		"metadata": { 
			"key1": $i, 
			"id": "test"
		}, 
		"data": { 
			"who": "I am $i", 
			"serial": $RANDOM, 
			"lol": "kek", 
			"something": "big", 
			"well": "bigger", 
			"wellwell": "bigger", 
			"wellwell2": "bigger", 
			"wel": "biaaaaagger", 
			"wl": "biggasfer", 
			"wal": "biggasfer", 
			"w": "biggasfer", 
			"v2": "biasfgger", 
			"vell": "bigger", 
			"vellwell": "bigger", 
			"vellwell2": "bigger", 
			"vel": "biaaaaagger", 
			"vl": "biggasfer", 
			"val": "biggasfer", 
			"v": "biggasfer", 
			"w2": "biasfgger", 
			"xell": "bigger", 
			"xellwell": "bigger", 
			"xellwell2": "bigger", 
			"xel": "biaaaaagger", 
			"xl": "biggasfer", 
			"xal": "biggasfer", 
			"x": "biggasfer", 
			"z2": "biasfgger", 
			"zell": "bigger", 
			"zellwell": "bigger", 
			"zellwell2": "bigger", 
			"zel": "biaaaaagger", 
			"zl": "biggasfer", 
			"yl": "biggasfer", 
			"y": "biggasfer", 
			"y2": "biasfgger", 
			"a number follows": 123 
		} 
	}
	_EOF_
done
