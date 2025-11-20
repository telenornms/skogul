// Package gen consists of auto-generated protobuf code for the Junos
// streaming telemetry interface and USP (User Services Platform).
//
// # Regenerating Protocol Buffer Code
//
// This project uses buf (https://buf.build) for protocol buffer code generation
// with vtprotobuf (https://github.com/planetscale/vtprotobuf) for optimized
// serialization performance.
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
// 4. Runs buf generate to create:
//   - Standard protobuf Go code (.pb.go files)
//   - vtprotobuf optimized code (_vtproto.pb.go files)
//
// Generated code is placed in:
//   - gen/junos/telemetry/
//   - gen/usp/
//
// # Configuration Files
//
// - buf.yaml: Defines the buf workspace and modules
// - buf.gen.yaml: Configures code generation plugins and options
// - gen/generate.sh: Orchestrates extraction and generation
//
// # Migration from gogo/protobuf
//
// This project was migrated from github.com/gogo/protobuf (now deprecated)
// to the official google.golang.org/protobuf with vtprotobuf for performance.
// The migration maintains API compatibility while using modern protobuf tooling.
package gen

//go:generate ./generate.sh
