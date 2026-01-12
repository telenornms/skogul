#!/bin/bash
# SQL TLS Integration Tests - tests mTLS with MySQL and PostgreSQL
set -e

cd "$(dirname "${BASH_SOURCE[0]}")"

COMPOSE=$(command -v docker &>/dev/null && docker compose version &>/dev/null && echo "docker compose" || echo "podman-compose")
SKOGUL="${SKOGUL_BIN:-../../skogul}"
PASSED=0 FAILED=0

log() { echo -e "\033[0;33m[INFO]\033[0m $1"; }
pass() {
    echo -e "\033[0;32m[PASS]\033[0m $1"
    PASSED=$((PASSED + 1))
}
fail() {
    echo -e "\033[0;31m[FAIL]\033[0m $1"
    FAILED=$((FAILED + 1))
}

cleanup() {
    pkill -f "skogul.*sql-tls" 2>/dev/null || true
    $COMPOSE down -v 2>/dev/null || true
}
trap cleanup EXIT

gen_sender_config() {
    local driver=$1 port=$2 connstr=$3 tls_opts=$4
    cat <<EOF
{"receivers":{"http":{"type":"http","address":"localhost:$port","handlers":{"/":"h"}}},"handlers":{"h":{"parser":"skogul","sender":"db"}},"senders":{"db":{"type":"sql","driver":"$driver","connstr":"$connstr","query":"INSERT INTO metrics (time,host,metric_name,value) VALUES(\${timestamp},\${metadata.host},\${metric_name},\${value})"$tls_opts}}}
EOF
}

gen_receiver_config() {
    local driver=$1 connstr=$2 tls_opts=$3
    cat <<EOF
{"receivers":{"db":{"type":"sql","driver":"$driver","connstr":"$connstr","query":"SELECT time,host,value FROM source_metrics LIMIT 1","handler":"h","interval":"-1s"$tls_opts}},"handlers":{"h":{"parser":"skogul","sender":"debug"}},"senders":{"debug":{"type":"debug"}}}
EOF
}

test_sender() {
    local name=$1 port=$2 config=$3 expect_success=$4
    echo "$config" >/tmp/skogul-test.json
    timeout 10 "$SKOGUL" -f /tmp/skogul-test.json &>/dev/null &
    local pid=$!
    sleep 2

    if ! kill -0 $pid 2>/dev/null; then
        if [ "$expect_success" = "false" ]; then
            pass "$name - rejected"
        else
            fail "$name - crashed"
        fi
        return
    fi
    local code
    code=$(curl -s -o /dev/null -w "%{http_code}" -X POST -H "Content-Type: application/json" \
        -d '{"metrics":[{"timestamp":"2024-01-01T00:00:00Z","metadata":{"host":"test"},"data":{"metric_name":"test","value":1}}]}' \
        "http://localhost:$port/" 2>/dev/null || echo "000")
    kill $pid 2>/dev/null
    wait $pid 2>/dev/null || true

    if [ "$expect_success" = "true" ] && [ "$code" -ge 200 ] && [ "$code" -lt 300 ]; then
        pass "$name"
    elif [ "$expect_success" = "false" ] && { [ "$code" -ge 400 ] || [ "$code" = "000" ]; }; then
        pass "$name - failed as expected"
    else
        fail "$name (HTTP $code)"
    fi
}

test_receiver() {
    local name=$1 config=$2
    echo "$config" >/tmp/skogul-test.json
    local out
    out=$(timeout 5 "$SKOGUL" -f /tmp/skogul-test.json 2>&1 || true)
    if echo "$out" | grep -q "host"; then
        pass "$name"
    else
        fail "$name"
    fi
}

log "Generating certificates..."
(cd certs && ./generate-certs.sh)
cp certs/rsa/*.pem certs/
chmod 600 certs/*-key.pem

log "Building skogul..."
[ -f "$SKOGUL" ] || (cd ../.. && make)

log "Starting containers..."
$COMPOSE up -d

for _ in $(seq 1 60); do
    $COMPOSE exec -T mysql healthcheck.sh --connect --innodb_initialized 2>/dev/null && break
    sleep 2
done

for _ in $(seq 1 60); do
    $COMPOSE exec -T postgres pg_isready -U postgres 2>/dev/null && break
    sleep 2
done
sleep 2

log "Running tests..."
MTLS=',"cafile":"certs/ca.pem","certfile":"certs/client.pem","keyfile":"certs/client-key.pem"'
SKIP=',"certfile":"certs/client.pem","keyfile":"certs/client-key.pem","insecure":true'
CAONLY=',"cafile":"certs/ca.pem"'
WRONG=',"cafile":"certs/ca.pem","certfile":"certs/wrong-client.pem","keyfile":"certs/wrong-client-key.pem"'

test_sender "MySQL mTLS" 18080 "$(gen_sender_config mysql 18080 'testuser@tcp(localhost:3306)/skogul_test' "$MTLS")" true
test_sender "MySQL skip-verify" 18081 "$(gen_sender_config mysql 18081 'testuser@tcp(localhost:3306)/skogul_test' "$SKIP")" true
test_sender "MySQL no-client-cert" 18082 "$(gen_sender_config mysql 18082 'testuser@tcp(localhost:3306)/skogul_test' "$CAONLY")" false
test_sender "MySQL wrong-cert" 18083 "$(gen_sender_config mysql 18083 'testuser@tcp(localhost:3306)/skogul_test' "$WRONG")" false

test_sender "PostgreSQL mTLS" 18090 "$(gen_sender_config postgres 18090 'host=localhost user=testuser dbname=skogul_test' "$MTLS")" true
test_sender "PostgreSQL skip-verify" 18091 "$(gen_sender_config postgres 18091 'host=localhost user=testuser dbname=skogul_test' "$SKIP")" true
test_sender "PostgreSQL no-client-cert" 18092 "$(gen_sender_config postgres 18092 'host=localhost user=testuser dbname=skogul_test' "$CAONLY")" false
test_sender "PostgreSQL wrong-cert" 18093 "$(gen_sender_config postgres 18093 'host=localhost user=testuser dbname=skogul_test' "$WRONG")" false

test_receiver "MySQL receiver" "$(gen_receiver_config mysql 'testuser@tcp(localhost:3306)/skogul_test?parseTime=true' "$MTLS")"
test_receiver "PostgreSQL receiver" "$(gen_receiver_config postgres 'host=localhost user=testuser dbname=skogul_test' "$MTLS")"

echo ""
echo "Results: $PASSED passed, $FAILED failed"
[ $FAILED -eq 0 ] && exit 0 || exit 1
