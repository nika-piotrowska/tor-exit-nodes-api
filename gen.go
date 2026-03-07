// Package main contains go:generate directives for OpenAPI code generation.
package main

//go:generate go run github.com/ogen-go/ogen/cmd/ogen@latest --target internal/generated/oas --package oas --clean api/openapi.yaml
