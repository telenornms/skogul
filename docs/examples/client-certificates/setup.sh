#!/bin/bash

set -e

echo "==> Cleaning up ..."

rm -rf {certsdb,private} index.txt* serial* *.pem

mkdir -p {certsdb,private}
touch index.txt

echo "==> Generating keys for a CA which we want to trust"
openssl req -batch -new -newkey rsa:4096 -keyout private/ca_key.pem -out ca_csr.pem -nodes -config ./CA.conf -subj "/O=Corp/CN=localhost/C=NO/ST=Svalbard/L=Svalbard/OU=Office"
echo ""

echo "==> Signing CA"
openssl ca -batch -create_serial -out ca_cert.pem -days 365 -keyfile private/ca_key.pem -selfsign -extensions v3_ca_has_san -config ./CA.conf -infiles ca_csr.pem
echo ""

echo "==> Generating CSR for Server"
openssl req -newkey rsa:4096 -keyout private/server_key.pem -out server_csr.pem -nodes -days 90 -config ./server.conf -subj "/O=Corp/CN=Server/C=NO/ST=Svalbard/L=Svalbard/OU=Office"
echo ""

echo "==> Generating CSR for Alice"
openssl req -newkey rsa:4096 -keyout private/alice_key.pem -out alice_csr.pem -nodes -days 90 -config ./alice.conf -subj "/O=Corp/CN=Alice/C=NO/ST=Svalbard/L=Svalbard/OU=Office"
echo ""

echo "==> Generating CSR for Bob"
openssl req -newkey rsa:4096 -keyout private/bob_key.pem -out bob_csr.pem -nodes -days 90 -subj "/CN=Bob"
echo ""

echo "==> Sign server certificate using CA"
openssl ca -batch -config ./CA.conf -extensions v3_req -out server_cert.pem -infiles server_csr.pem
echo ""

echo "==> Sign Alice's certificate using CA"
openssl ca -batch -config ./CA.conf -extensions v3_req -out alice_cert.pem -infiles alice_csr.pem
echo ""

echo "==> Bob signs his own certificate"
openssl x509 -req -in bob_csr.pem -signkey private/bob_key.pem -out bob_cert.pem -days 90
echo ""

echo "==> DEBUG Verify that Alice's cert is signed by CA"
openssl verify -CAfile ca_cert.pem alice_cert.pem
echo ""

echo "==> DEBUG Verify that Bob's cert is NOT signed by CA"
! openssl verify -CAfile ca_cert.pem bob_cert.pem
echo ""

echo "==> Summary"
echo ""

echo "Alice is good"
echo "Bob is bad"
echo ""

echo "CA signed Alice's cert"
echo "CA did NOT sign Bob's cert"
echo ""

echo "Alice should be granted access"
echo "Bob should not be granted access"
echo ""

echo "==> Done"
