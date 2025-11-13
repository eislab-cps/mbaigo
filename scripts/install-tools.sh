#!/usr/bin/env bash
# Install optional development tools
# These are NOT required for building or using mbaigo
# Dependencies are managed via go.mod

set -e

echo "Installing optional development tools..."
echo ""
echo "Note: These tools are optional for code quality checks."
echo "They are NOT required to build or use mbaigo."
echo "All dependencies are managed via go.mod"
echo ""

# Check if tools are already installed
check_tool() {
    if command -v "$1" &> /dev/null; then
        echo "  $1 already installed"
        return 0
    else
        return 1
    fi
}

# Install gocyclo
if ! check_tool gocyclo; then
    echo "  Installing gocyclo (cyclomatic complexity)..."
    go install github.com/fzipp/gocyclo/cmd/gocyclo@latest
fi

# Install gosec
if ! check_tool gosec; then
    echo "  Installing gosec (security scanner)..."
    go install github.com/securego/gosec/v2/cmd/gosec@latest
fi

# Install staticcheck
if ! check_tool staticcheck; then
    echo "  Installing staticcheck (static analysis)..."
    go install honnef.co/go/tools/cmd/staticcheck@latest
fi

# Install govulncheck
if ! check_tool govulncheck; then
    echo "  Installing govulncheck (vulnerability scanner)..."
    go install golang.org/x/vuln/cmd/govulncheck@latest
fi

echo ""
echo "Optional tools installed!"
echo ""
echo "Additional optional tools:"
echo "  - typos (spell checker): cargo install typos-cli"
echo "    https://github.com/crate-ci/typos"
echo ""
echo "To manage project dependencies, use standard Go commands:"
echo "  - go mod download    # Download dependencies"
echo "  - go mod tidy        # Add missing and remove unused modules"
echo "  - go mod vendor      # Make vendored copy of dependencies"
echo "  - go mod verify      # Verify dependencies have expected content"
