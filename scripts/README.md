# Build and Development Scripts

This directory contains shell scripts for common development tasks.

## Managing Dependencies

**mbaigo uses standard Go modules.** All dependencies are declared in `go.mod`.

```bash
# Download dependencies
go mod download

# Add missing and remove unused modules
go mod tidy

# Verify dependencies
go mod verify

# Create vendored copy (optional)
go mod vendor
```

Or use the Makefile:

```bash
make build     # Build (downloads dependencies automatically)
make vendor    # Vendor dependencies (optional)
```

## Available Scripts

### `test.sh`
Run all tests with race detection and generate coverage report.

```bash
./scripts/test.sh
# Or: make test
```

### `lint.sh`
Run all linters and static analysis tools:
- `gofmt` - Code formatting check
- `go vet` - Go static analysis
- `gosec` - Security vulnerability scanner
- `staticcheck` - Advanced static analysis
- `govulncheck` - Known vulnerability database check

```bash
./scripts/lint.sh
# Or: make lint
```

**Note:** Linting tools are optional. Install with `make tools` or `./scripts/install-tools.sh`

### `coverage-report.sh`
Generate detailed HTML coverage report and show cyclomatic complexity.

```bash
./scripts/coverage-report.sh
# Or: make analyse
```

Output: `cover.html` (open in browser)

### `spellcheck.sh`
Run spell checker on code and documentation.

```bash
./scripts/spellcheck.sh
# Or: make spellcheck
```

Requires: [typos](https://github.com/crate-ci/typos) (install with `cargo install typos-cli`)

### `install-tools.sh`
Install **optional** development tools for code quality checks.

```bash
./scripts/install-tools.sh
# Or: make tools
```

Installs (globally):
- gocyclo - Cyclomatic complexity
- gosec - Security scanning
- staticcheck - Static analysis
- govulncheck - Vulnerability scanning

**These are NOT required** to build or use mbaigo. They're optional tools for development.

## Usage with Make

All scripts can be invoked via the Makefile:

```bash
# Build
make build          # Build all packages and CLI
make vendor         # Vendor dependencies (optional)

# Development
make test           # Run tests
make lint           # Run linters
make spellcheck     # Run spell checker
make runchecks      # Run all checks

# Analysis
make analyse        # Generate coverage report

# Tools
make tools          # Install optional development tools

# Cleanup
make clean          # Clean temporary files
make clean-vendor   # Remove vendor directory
```

**Note:** Dependencies are downloaded automatically by Go. No separate command needed.

## Prerequisites

### Required
- Go 1.24.4 or later
- Bash shell

### Optional
- Development tools (install with `make tools`)
- typos spell checker (install with `cargo install typos-cli`)

## Dependency Management Philosophy

**mbaigo follows Go best practices:**

1. **Dependencies in go.mod**: All runtime dependencies are in `go.mod`
2. **No global tools required**: You can build with just Go installed
3. **Optional dev tools**: Quality tools (linters, etc.) are optional
4. **Vendor optional**: Vendoring is supported but not required
5. **Reproducible builds**: `go.mod` and `go.sum` ensure reproducibility

### When to Vendor

Consider vendoring if:
- You need offline builds
- You want to freeze dependencies exactly
- Your CI environment requires it

```bash
make vendor
```

This creates a `vendor/` directory with all dependencies.

## CI/CD Integration

These scripts are used by GitHub Actions workflows in `.github/workflows/`.

The CI environment:
1. Uses `go mod download` to get dependencies
2. Installs linting tools as needed
3. Runs tests and checks
4. Does NOT require vendoring
