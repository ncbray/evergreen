#!/bin/bash
set -e
echo "### Building binary"
go install tools/regenerate
echo "### Generating sources"
bin/regenerate -replace
echo "### Formatting sources"
go fmt ./...
