// Package gen consists of auto-generated protobuf code for the junos
// streaming telemetry interface. This file provides the commands to
// (re)generate the files, assuming you have downloaded the junos telemetry
// interface files, as provided by Juniper under the Apache 2.0 License.
package gen

//go:generate /bin/bash -c "rm -f *pb.go"
//go:generate /bin/bash -c "for a in junos-telemetry-interface/*.proto; do protoc --go_out=. -Ijunos-telemetry-interface/ $DOLLAR{a}; done"
//go:generate /bin/bash -c "sed -i 's/^package.*/package gen/g' *.go"
