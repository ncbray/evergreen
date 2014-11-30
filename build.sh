#!/bin/bash
set -e
echo "### Building binary"
go install evergreen/cmd/regenerate
echo "### Generating sources"
bin/regenerate -replace
echo "### Formatting sources"
go fmt ./...
