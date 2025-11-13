#!/usr/bin/env bash
# Run source code linters and security checks

set -e

echo "Checking code formatting..."
if [ -n "$(gofmt -l .)" ]; then
    echo "Error: Code is not gofmt'ed!"
    echo "Run: gofmt -w ."
    exit 1
fi

echo "Running go vet..."
go vet $(go list ./... | grep -v /tmp)

echo "Running gosec (security scanner)..."
gosec -quiet -fmt=golint -exclude-dir="tmp" ./...

echo "Running staticcheck..."
staticcheck ./...

echo "Running govulncheck..."
govulncheck -test ./...

echo ""
echo "All linters passed!"
