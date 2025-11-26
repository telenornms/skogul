#!/bin/bash
# Script to regenerate protobuf code for Junos telemetry and USP interfaces
# This script handles both GNU sed (Linux) and BSD sed (macOS)

set -e

cd "$(dirname "$0")"

# Cross-platform in-place sed
# BSD sed (macOS) requires -i '' while GNU sed requires just -i
sed_inplace() {
    if sed --version 2>/dev/null | grep -q GNU; then
        sed -i "$@"
    else
        sed -i '' "$@"
    fi
}

echo "Generating Junos telemetry protobuf code..."

# Clean and prepare output directory
rm -f junos/telemetry/*pb.go
mkdir -p junos/telemetry

# Extract proto files
tar xzf tar-balls/junos-telemetry-interface-23.2R1.tar.gz

# Generate Go code from proto files, skipping gnmi/sr_/Gnmi patterns
for a in junos-telemetry-interface/*.proto; do
    if echo "${a}" | grep -Eqv '/(gnmi|sr_|Gnmi)'; then
        protoc --gogo_out=junos/telemetry --gogo_opt=M="${PWD}/junos-telemetry-interface/" -Ijunos-telemetry-interface/ "${a}"
    else
        echo "skipping ${a}"
    fi
done

# Fix package names in generated files
sed_inplace 's/^package.*/package telemetry/g' junos/telemetry/*.go

echo "Generating USP protobuf code..."

# Clean and prepare output directory
rm -f usp/*pb.go
mkdir -p usp

# Extract proto files
tar xf tar-balls/usp-interface-1-1.tar.gz

# Generate Go code
protoc --gogo_out=usp --gogo_opt=M="${PWD}/usp" usp-record-1-1.proto
protoc --gogo_out=usp --gogo_opt=M="${PWD}/usp" usp-msg-1-1.proto

echo "Done!"
