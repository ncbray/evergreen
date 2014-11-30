#!/bin/bash
set -e
echo "### Cleaning directory"
rm -r src/generated 2> /dev/null || true
echo "### Generating sources"
go run src/evergreen/cmd/regenerate/main.go
echo "### Running tests"
go test ./...
