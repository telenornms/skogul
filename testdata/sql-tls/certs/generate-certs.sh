#!/bin/bash
# Generate test certificates for SQL TLS testing
set -e
cd "$(dirname "${BASH_SOURCE[0]}")"

rm -rf rsa && mkdir -p rsa
cd rsa

# CA
openssl genrsa -out ca-key.pem 2048 2>/dev/null
openssl req -new -x509 -days 365 -key ca-key.pem -out ca.pem -subj "/CN=Test CA/O=Skogul Test"

# Server certs (MySQL and PostgreSQL)
for srv in mysql postgres; do
    openssl genrsa -out ${srv}-server-key.pem 2048 2>/dev/null
    openssl req -new -key ${srv}-server-key.pem -subj "/CN=${srv}/O=Skogul Test" |
        openssl x509 -req -days 365 -CA ca.pem -CAkey ca-key.pem -CAcreateserial \
            -extfile <(echo "subjectAltName=DNS:${srv},DNS:localhost,IP:127.0.0.1") \
            -out ${srv}-server.pem 2>/dev/null
done

# Valid client cert
openssl genrsa -out client-key.pem 2048 2>/dev/null
openssl req -new -key client-key.pem -subj "/CN=testuser/O=Skogul Test" |
    openssl x509 -req -days 365 -CA ca.pem -CAkey ca-key.pem -CAcreateserial \
        -extfile <(echo "extendedKeyUsage=clientAuth") -out client.pem 2>/dev/null

# Wrong CA + client cert (for negative tests)
openssl genrsa -out wrong-ca-key.pem 2048 2>/dev/null
openssl req -new -x509 -days 365 -key wrong-ca-key.pem -out wrong-ca.pem -subj "/CN=Wrong CA"
openssl genrsa -out wrong-client-key.pem 2048 2>/dev/null
openssl req -new -key wrong-client-key.pem -subj "/CN=wronguser" |
    openssl x509 -req -days 365 -CA wrong-ca.pem -CAkey wrong-ca-key.pem -CAcreateserial \
        -extfile <(echo "extendedKeyUsage=clientAuth") -out wrong-client.pem 2>/dev/null

rm -f ./*.srl
chmod 600 ./*-key.pem
