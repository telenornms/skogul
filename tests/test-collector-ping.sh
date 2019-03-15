#!/bin/bash

me="$$"
offset="$(( ( $RANDOM % 4 ) + 1))"
generate() {
    timeit=$(date --iso-8601=seconds -d @$1)
    cat <<-_EOF_
{
    "src": "test",
    "timestamp": "$timeit",
    "template": {
        "src": "test-collector-pid-$me"
    },
    "metrics": [
_EOF_
    COMMA=" "
    for row in {1..50}; do
        for num in {1..4}; do 
            latency4="$offset.$(( $RANDOM % 1000 ))"
            latency6="$offset.$(( $RANDOM % 1000 ))"
    cat <<_EOF_
        ${COMMA}{
            "metadata":  {
                "switch": "row$row-n$num"
            },
            "data": {
                "latency4": $latency4,
                "latency6": $latency6
            }
        }
_EOF_
        COMMA=","
        done
    done
cat <<_EOF_
    ]
}
_EOF_
}

startd=$(( $(date +%s) - 3600 ))
for a in {0..3600}; do
    generate $(( $startd + $a )) | POST -H Content-type:\ application/json http://[::1]:8080/api/write/collector > /dev/null
    echo -n "."
done
