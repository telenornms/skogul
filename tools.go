//go:build tools
// +build tools

// Package tools tracks dependencies for code generation and build tools.
// These are installed via `go install` in gen/generate.sh.
// This file ensures that `go mod tidy` won't remove tool dependencies.
package tools

import (
	_ "github.com/planetscale/vtprotobuf/cmd/protoc-gen-go-vtproto"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
)
