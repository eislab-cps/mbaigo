#!/usr/bin/env bash
# Generate detailed coverage report

set -e

if [ ! -f ".cover.out" ]; then
    echo "No coverage data found. Running tests first..."
    ./scripts/test.sh
fi

echo "Generating HTML coverage report..."
go tool cover -html=".cover.out" -o="cover.html"

echo ""
echo "COVERAGE SUMMARY"
echo "===================="
go tool cover -func=.cover.out

echo ""
echo "CYCLOMATIC COMPLEXITY"
echo "===================="
gocyclo -avg -top 10 .

echo ""
echo "HTML coverage report: cover.html"
