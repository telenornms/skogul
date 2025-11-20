#!/bin/bash
set -euo pipefail

# Script to extract protobuf tarballs and generate Go code using buf

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

cd "$PROJECT_ROOT"

echo "==> Cleaning old proto extraction directories..."
rm -rf gen/junos-telemetry-interface gen/usp-interface

echo "==> Extracting Junos telemetry proto files..."
mkdir -p gen/junos-telemetry-interface
tar xzf gen/tar-balls/junos-telemetry-interface-25.2R1.8-EVO.tar.gz -C gen/junos-telemetry-interface

echo "==> Extracting USP proto files..."
mkdir -p gen/usp-interface
tar xzf gen/tar-balls/usp-interface-1-1.tar.gz -C gen/usp-interface

echo "==> Injecting go_package options into proto files..."
for proto_file in gen/junos-telemetry-interface/*.proto; do
  if [ -f "$proto_file" ]; then
    # Check if uncommented go_package already exists
    if ! grep -q "^option go_package" "$proto_file"; then
      sed -i.bak '/^syntax = /a\
option go_package = "github.com/telenornms/skogul/gen/junos/telemetry";
' "$proto_file"
      rm -f "${proto_file}.bak"
    fi
  fi
done

# Remove administrative protos (gnmi, sr_, Gnmi patterns)
find gen/junos-telemetry-interface -name "*.proto" | grep -E '(gnmi|sr_|Gnmi)' | xargs rm -f || true
# ddosd-junos-state-ddos-protection-render.proto excluded via buf.yaml (State message collision)

echo "==> Cleaning old generated Go files..."
rm -rf gen/junos gen/usp

echo "==> Running buf generate..."
go run github.com/bufbuild/buf/cmd/buf generate

echo "==> Moving generated files to proper locations..."
mkdir -p gen/junos/telemetry
mkdir -p gen/usp

for file in gen/*.pb.go; do
  if [ -f "$file" ]; then
    if grep -q "^package usp$" "$file" 2>/dev/null; then
      mv "$file" gen/usp/
    else
      # otherwise it's a Junos telemetry file (package telemetry)
      mv "$file" gen/junos/telemetry/
    fi
  fi
done

echo "==> Protocol buffer generation complete!"
echo "    - gen/junos/telemetry/"
echo "    - gen/usp/"
