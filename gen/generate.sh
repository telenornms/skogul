#!/bin/bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
cd "$PROJECT_ROOT"

if ! command -v protoc &>/dev/null; then
    echo "Error: protoc is not installed or not in PATH" >&2
    echo "Install it from: https://github.com/protocolbuffers/protobuf/releases" >&2
    exit 1
fi
if ! command -v go &>/dev/null; then
    echo "Error: go is not installed or not in PATH" >&2
    exit 1
fi

PATH="$(go env GOPATH)/bin:$PATH"
export PATH

# These Juniper protos from 23.2R1 were removed in 25.2R1.8-EVO, but do not
# collide with any protos included there and are included for backward
# compatibility with older Junos devices.
LEGACY_PROTOS=(
    "lsp_stats.proto"
    "qmon.proto"
    "ancpd_oc.proto"
    "authd_oc.proto"
    "bbe-smgd_ancp_stats_oc.proto"
    "bbe-smgd_pppoe_stats_oc.proto"
    "bbe-smgd_rsmon_debug_oc.proto"
    "bbe-smgd_rsmon_stats_oc.proto"
    "bbe-smgd_smd_queue_stats_oc.proto"
    "bbe-smgd_sub_mgmt_network_stats_oc.proto"
    "dcd_oc.proto"
    "eventd.proto"
    "jdhcpd_oc.proto"
    "jl2tpd_oc.proto"
    "jpppd_oc.proto"
    "kmd_render.proto"
    "mib2d_nd6_oc.proto"
    "mib2d_oc.proto"
    "pfed_oc.proto"
    "pfe_ifl_oc.proto"
    "pfe_npu_resource.proto"
    "pfe_port_oc.proto"
    "xmlproxyd_show_local_interface_oc.proto"
    "rpd_loc_rib_oc.proto"
    "smid_oc.proto"
    "kernel-ifstate-render.proto"
    "ipsec_telemetry.proto"
    "svcset_telemetry.proto"
    "session_telemetry.proto"
    "bbe-statsd-telemetry_oc.proto"
    "jkhmd_oc.proto"
    "jdiameterd_render.proto"
    "sr_te_per_lsp_transit_stats.proto"
    "sr_te_per_lsp_ingress_stats.proto"
    "jkdsd_oc.proto"
    "jkdsd_cpu_oc.proto"
    "jkhmd_resiliency_render.proto"
    "spu_cpu_util.proto"
    "pfe_ifl_family_v4_stats_oc.proto"
    "pfe_ifl_family_v6_stats_oc.proto"
    "saegw-upad_oc.proto"
    "nasd_oc.proto"
    "ngapd_oc.proto"
    "pfe_page_drop_oc.proto"
    "pfe-junos-slice-egr-qstats-render.proto"
    "chassisd-junos-state-poe-render.proto"
    "chassisd-junos-state-chassis-render.proto"
    "mib2d-junos-state-interfaces-render.proto"
    "xmlproxyd-junos-openconfig-system-render.proto"
    "sysd-junos-openconfig-system-render.proto"
    "pbj.proto"
)

# V23.2R1 protos removed in 25.2R1.8-EVO that have conflicts (need to pick
# whether to keep v23 or v25)
# chassisd_oc.proto,
# mib2d_arp_oc.proto,
# rmopd_render.proto,
# alarmd_oc.proto,
# cosd_oc.proto,
# spu_flow_stats.proto

WORK_DIR=$(mktemp -d)
trap 'rm -rf "$WORK_DIR"' EXIT

echo "==> Extracting proto files..."
mkdir -p "$WORK_DIR/junos" "$WORK_DIR/junos-legacy" "$WORK_DIR/usp"
tar xzf gen/tar-balls/junos-telemetry-interface-25.2R1.8-EVO.tar.gz -C "$WORK_DIR/junos"
tar xzf gen/tar-balls/junos-telemetry-interface-23.2R1.tar.gz -C "$WORK_DIR/junos-legacy"
tar xzf gen/tar-balls/usp-interface-1-1.tar.gz -C "$WORK_DIR/usp"

# Restore legacy protos from 23.2R1 (only if not present in current version)
for proto in "${LEGACY_PROTOS[@]}"; do
    src="$WORK_DIR/junos-legacy/junos-telemetry-interface/$proto"
    dst="$WORK_DIR/junos/$proto"
    [ -f "$src" ] && [ ! -f "$dst" ] && cp "$src" "$dst"
done

echo "==> Injecting go_package options..."
for proto_file in "$WORK_DIR/junos"/*.proto; do
    [ -f "$proto_file" ] && ! grep -q "^option go_package" "$proto_file" &&
        sed -i.bak '/^syntax = /a\
option go_package = "github.com/telenornms/skogul/gen/junos/telemetry";
' "$proto_file" && rm -f "${proto_file}.bak"
done
for proto_file in "$WORK_DIR/usp"/*.proto; do
    [ -f "$proto_file" ] && ! grep -q "^option go_package" "$proto_file" &&
        sed -i.bak '/^syntax = /a\
option go_package = "github.com/telenornms/skogul/gen/usp";
' "$proto_file" && rm -f "${proto_file}.bak"
done

# Remove administrative protos (gnmi, sr_, Gnmi patterns)
find "$WORK_DIR/junos" -name "*.proto" | grep -E '(gnmi|sr_|Gnmi)' | xargs rm -f || true
# State message collision within junos-telemetry-interface-25.2R1.8-EVO
# (collides with transceiver.proto which is prioritised in this case)
rm -f "$WORK_DIR/junos/ddosd-junos-state-ddos-protection-render.proto" || true

echo "==> Installing protoc plugins..."
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install github.com/planetscale/vtprotobuf/cmd/protoc-gen-go-vtproto@latest

echo "==> Generating Go code with protoc..."
rm -rf gen/junos/telemetry gen/usp
mkdir -p gen/junos/telemetry gen/usp

if ! protoc \
    --go_out=gen/junos/telemetry --go_opt=paths=source_relative \
    --go-vtproto_out=gen/junos/telemetry --go-vtproto_opt=paths=source_relative,features=marshal+unmarshal+size+pool \
    -I "$WORK_DIR/junos" \
    "$WORK_DIR/junos"/*.proto; then
    echo "Error: protoc failed for junos telemetry. Check proto file syntax." >&2
    exit 1
fi

if ! protoc \
    --go_out=gen/usp --go_opt=paths=source_relative \
    --go-vtproto_out=gen/usp --go-vtproto_opt=paths=source_relative,features=marshal+unmarshal+size+pool \
    -I "$WORK_DIR/usp" \
    "$WORK_DIR/usp"/*.proto; then
    echo "Error: protoc failed for USP. Check proto file syntax." >&2
    exit 1
fi

echo "==> Done: gen/junos/telemetry/, gen/usp/"
