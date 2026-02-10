#!/bin/bash
# FD Leak Regression Test — proves that the resp.Body.Close() fix
# prevents TCP file descriptor leaks in the HTTP sender.
#
# Outputs a CSV report to stdout (redirect to file for graphing).
# Human-readable progress goes to stderr.
set -e

cd "$(dirname "${BASH_SOURCE[0]}")"

COMPOSE=$(command -v docker &>/dev/null && docker compose version &>/dev/null && echo "docker compose" || echo "podman-compose")

log()  { echo -e "\033[0;33m[INFO]\033[0m $1" >&2; }

cleanup() {
    $COMPOSE down -v --remove-orphans 2>/dev/null || true
}
trap cleanup EXIT

PAYLOAD='{"metrics":[{"timestamp":"2024-01-01T00:00:00Z","metadata":{"host":"fdtest"},"data":{"value":1}}]}'

ROUNDS=5
REQUESTS_PER_ROUND=200
DOWN_ROUNDS=3
DOWN_REQUESTS=100

# ── helpers ──────────────────────────────────────────────────────────

send_requests() {
    local port=$1 count=$2
    for _ in $(seq 1 "$count"); do
        curl -s -o /dev/null -w '' -X POST \
            -H "Content-Type: application/json" \
            -d "$PAYLOAD" \
            "http://localhost:${port}/" 2>/dev/null || true
    done
}

get_fds() {
    local container=$1
    docker exec "$container" sh -c 'ls /proc/1/fd 2>/dev/null | wc -l' 2>/dev/null | tr -d '[:space:]'
}

get_socket_fds() {
    local container=$1
    docker exec "$container" sh -c 'ls -la /proc/1/fd 2>/dev/null | grep socket | wc -l' 2>/dev/null | tr -d '[:space:]'
}

get_goroutines() {
    local port=$1
    curl -s "http://localhost:${port}/debug/pprof/goroutine?debug=0" 2>/dev/null \
        | head -1 | grep -oE 'goroutine profile: total [0-9]+' | grep -oE '[0-9]+' || echo ""
}

sample() {
    # Collect all six values for both binaries at once
    local old_fds old_sock old_gorout new_fds new_sock new_gorout
    old_fds=$(get_fds fd-leak-skogul-old)
    old_sock=$(get_socket_fds fd-leak-skogul-old)
    old_gorout=$(get_goroutines 6061)
    new_fds=$(get_fds fd-leak-skogul-new)
    new_sock=$(get_socket_fds fd-leak-skogul-new)
    new_gorout=$(get_goroutines 6062)
    echo "${old_fds},${old_sock},${old_gorout},${new_fds},${new_sock},${new_gorout}"
}

# ── build & start ────────────────────────────────────────────────────

log "Checking for binaries..."
if [ ! -f skogul-old ] || [ ! -f skogul-new ]; then
    cat >&2 <<'USAGE'

ERROR: skogul-old and skogul-new binaries must exist in this directory.

Build them with (from project root):
  GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o testdata/fd-leak-test/skogul-new ./cmd/skogul
  git stash && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o testdata/fd-leak-test/skogul-old ./cmd/skogul && git stash pop

USAGE
    exit 1
fi

log "Starting containers..."
$COMPOSE build --quiet 2>&2
$COMPOSE up -d 2>&2

log "Waiting for services to be ready..."
sleep 3

for port in 8081 8082; do
    for _ in $(seq 1 20); do
        if curl -s -o /dev/null -w "%{http_code}" -X POST \
            -H "Content-Type: application/json" \
            -d "$PAYLOAD" \
            "http://localhost:${port}/" 2>/dev/null | grep -q "^2"; then
            break
        fi
        sleep 1
    done
done

# ── CSV header ───────────────────────────────────────────────────────

echo "phase,round,cumulative_requests,old_fds,old_sockets,old_goroutines,new_fds,new_sockets,new_goroutines"

# ── Phase 1: Target UP ──────────────────────────────────────────────

log "Phase 1: Target UP — sending $((ROUNDS * REQUESTS_PER_ROUND)) requests per binary"

cumulative=0

# Baseline measurement before any load
vals=$(sample)
echo "up,0,${cumulative},${vals}"
log "  baseline: $vals"

for round in $(seq 1 $ROUNDS); do
    send_requests 8081 $REQUESTS_PER_ROUND
    send_requests 8082 $REQUESTS_PER_ROUND
    cumulative=$((cumulative + REQUESTS_PER_ROUND))
    sleep 1

    vals=$(sample)
    echo "up,${round},${cumulative},${vals}"
    log "  round ${round} (${cumulative} reqs): ${vals}"
done

# ── Phase 2: Target DOWN ────────────────────────────────────────────

log "Phase 2: Target DOWN — sending $((DOWN_ROUNDS * DOWN_REQUESTS)) failing requests per binary"

$COMPOSE stop target 2>/dev/null
sleep 2

# Baseline for down phase (continues cumulative count)
vals=$(sample)
echo "down,0,${cumulative},${vals}"
log "  baseline (target stopped): $vals"

for round in $(seq 1 $DOWN_ROUNDS); do
    send_requests 8081 $DOWN_REQUESTS
    send_requests 8082 $DOWN_REQUESTS
    cumulative=$((cumulative + DOWN_REQUESTS))
    sleep 1

    vals=$(sample)
    echo "down,${round},${cumulative},${vals}"
    log "  round ${round} (${cumulative} reqs): ${vals}"
done

log "Done. Pipe stdout to a .csv file for graphing."
