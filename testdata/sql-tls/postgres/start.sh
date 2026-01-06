#!/bin/bash
set -e

mkdir -p /var/lib/postgresql/certs
sync 2>/dev/null || true
sleep 2

copy_with_verify() {
    local src="$1"
    local dst="$2"
    local max_retries=10
    local retry=0

    while [ $retry -lt $max_retries ]; do
        rm -f "$dst"
        dd if="$src" of="$dst" bs=1 2>/dev/null
        local src_size
        local dst_size
        src_size=$(wc -c <"$src")
        dst_size=$(wc -c <"$dst")

        if [ "$src_size" = "$dst_size" ] && [ "$src_size" -gt 0 ]; then
            if echo "$dst" | grep -q 'key\.pem$'; then
                grep -q "END.*PRIVATE KEY" "$dst" && return 0
            else
                grep -q "END CERTIFICATE" "$dst" && return 0
            fi
        fi

        retry=$((retry + 1))
        sleep 1
        sync 2>/dev/null || true
    done

    echo "ERROR: Failed to copy $src after $max_retries attempts"
    return 1
}

copy_with_verify /certs/ca.pem /var/lib/postgresql/certs/ca.pem
copy_with_verify /certs/postgres-server.pem /var/lib/postgresql/certs/postgres-server.pem
copy_with_verify /certs/postgres-server-key.pem /var/lib/postgresql/certs/postgres-server-key.pem

chown -R postgres:postgres /var/lib/postgresql/certs
chmod 600 /var/lib/postgresql/certs/postgres-server-key.pem
chmod 644 /var/lib/postgresql/certs/ca.pem
chmod 644 /var/lib/postgresql/certs/postgres-server.pem

exec docker-entrypoint.sh "$@"
