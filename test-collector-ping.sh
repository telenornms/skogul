#!/bin/bash





generate() {
    time=$(date --iso-8601=seconds -d @$1)
    cat <<-_EOF_

{
    "src": "test",
    "metrics": [
_EOF_
    COMMA=" "
    for row in {1..5}; do
        for num in {1..4}; do 
            latency4="1.$(( $RANDOM % 1000 ))"
            latency6="1.$(( $RANDOM % 1000 ))"
    cat <<_EOF_
        ${COMMA}{
            "timestamp": "$time",
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
    data="$(generate $(( $startd + $a )))"
    o=$(POST http://localhost:8080/api/write/collector <<<"$data")
    if [ "x$o" != "xOK" ]; then
        echo $o
        echo -e "$data"
        echo $o
        read

        fi
sleep 2
done
