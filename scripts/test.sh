#!/usr/bin/env bash
# Run tests with coverage report

set -e

echo "Running tests with race detection and coverage..."
go test -v -race -coverprofile=".cover.out" $(go list ./... | grep -v /tmp | grep -v /systems/)

echo ""
echo "Test coverage report:"
go tool cover -func=.cover.out
