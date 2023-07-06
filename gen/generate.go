// Package gen consists of auto-generated protobuf code for the junos
// streaming telemetry interface. This file provides the commands to
// (re)generate the files.
package gen

//go:generate /bin/bash -c "rm -f junos/telemetry/*pb.go; mkdir -p junos/telemetry"
//go:generate /bin/bash -c "tar xzf tar-balls/junos-telemetry-interface-23.2R1.tar.gz"
//go:generate /bin/bash -c "for a in junos-telemetry-interface/*.proto; do if echo $DOLLAR{a} | egrep -qv '/(gnmi|sr_|Gnmi)'; then protoc --gogo_out=junos/telemetry --gogo_opt=M=$DOLLAR{PWD}junos-telemetry-interface/ -Ijunos-telemetry-interface/ $DOLLAR{a}; else echo skipping $DOLLAR{a}; fi; done"
//go:generate /bin/bash -c "sed -i 's/^package.*/package telemetry/g' junos/telemetry/*.go"

//go:generate /bin/bash -c "rm -f usp/*pb.go; mkdir -p usp"
//go:generate /bin/bash -c "tar xf tar-balls/usp-interface-1-1.tar.gz"
//go:generate /bin/bash -c "protoc --gogo_out=usp --gogo_opt=M=${PWD}/usp usp-record-1-1.proto"
//go:generate /bin/bash -c "protoc --gogo_out=usp --gogo_opt=M=${PWD}/usp usp-msg-1-1.proto"
