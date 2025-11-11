#!/usr/bin/env bash
# Run spell checker on code and documentation

set -e

if ! command -v typos &> /dev/null; then
    echo "Error: typos not found"
    echo "Install with: cargo install typos-cli"
    echo "Or visit: https://github.com/crate-ci/typos"
    exit 1
fi

echo "Running spell checker..."
typos .

echo ""
echo "✓ Spell check passed!"
