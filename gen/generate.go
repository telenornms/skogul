// Package gen consists of auto-generated protobuf code for the Junos
// streaming telemetry interface and USP (User Services Platform).
//
// # Regenerating Protocol Buffer Code
//
// This project uses protoc with protoc-gen-go and protoc-gen-go-vtproto
// for protocol buffer code generation. vtprotobuf provides optimized
// serialization performance.
//
// Requirements:
//   - protoc (protocol buffer compiler)
//   - Go toolchain (plugins are installed automatically)
//
// To regenerate the protocol buffer code:
//
//	make generate
//
// Or directly:
//
//	./gen/generate.sh
//
// # What Happens During Generation
//
// 1. Extracts proto files from tarballs in gen/tar-balls/
//   - junos-telemetry-interface-25.2R1.8-EVO.tar.gz
//   - usp-interface-1-1.tar.gz
//
// 2. Injects go_package options into extracted proto files
//
// 3. Removes proto files matching /(gnmi|sr_|Gnmi)/ pattern
//
// 4. Installs protoc plugins via go install
//
// 5. Runs protoc to generate:
//   - Standard protobuf Go code (.pb.go files)
//   - vtprotobuf optimized code (_vtproto.pb.go files)
//
// Generated code is placed in:
//   - gen/junos/telemetry/
//   - gen/usp/
package gen

//go:generate ./generate.sh
