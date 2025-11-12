# Dependency Management in mbaigo

## Overview

**mbaigo follows Go standard practices for dependency management.** All dependencies are managed through `go.mod` and `go.sum`, not through custom installation scripts.

## Philosophy

1. **Standard Go Modules**: Use `go.mod` for all dependencies
2. **No Global Requirements**: Build with just Go installed
3. **Optional Dev Tools**: Quality tools are optional, not required
4. **Reproducible Builds**: `go.mod` + `go.sum` ensure consistency
5. **Vendoring Optional**: Support vendoring without requiring it

## Managing Dependencies

### Download Dependencies

**Dependencies are automatically downloaded** when you run `go build`, `go test`, or `make build`.

There's no need to run a separate download command in modern Go projects.

```bash
# Just build - dependencies download automatically
go build ./...
make build

# Optional: Explicit download (rarely needed)
go mod download        # Download modules to cache
go mod verify          # Verify checksums
```

### Add New Dependencies

```bash
# Add a new dependency
go get github.com/some/package

# Clean up
go mod tidy
```

### Update Dependencies

```bash
# Update all dependencies
go get -u ./...
go mod tidy

# Update specific package
go get -u github.com/some/package
```

### View Dependencies

```bash
# List all dependencies
go list -m all

# Why is this dependency here?
go mod why github.com/some/package

# Dependency graph
go mod graph
```

## Vendoring (Optional)

Vendoring creates a local copy of all dependencies in the `vendor/` directory.

### When to Vendor

Consider vendoring if you:
- Need offline builds
- Want to freeze dependencies exactly
- CI environment requires it
- Company policy mandates it

### How to Vendor

```bash
# Create vendor directory
make vendor
# Or: go mod vendor

# Build using vendor
go build -mod=vendor ./...

# Remove vendor
make clean-vendor
# Or: rm -rf vendor
```

### .gitignore Note

The `vendor/` directory is typically **not** committed to Git in library projects, but **may be** committed in application projects. Currently mbaigo doesn't commit vendor/.

## Development Tools (Optional)

Some commands use optional tools for code quality checks:

```bash
# Install optional tools
make tools
```

This installs:
- **gocyclo** - Cyclomatic complexity analyzer
- **gosec** - Security vulnerability scanner
- **staticcheck** - Advanced static analysis
- **govulncheck** - Known vulnerability checker

**These are NOT required** to build or use mbaigo. They're only for development quality checks.

### Check Without Installing

```bash
# Use built-in Go tools only
go fmt ./...
go vet ./...
go test ./...
```

## Makefile Commands

### Standard Operations

```bash
make build         # Build all packages and CLI
make test          # Run tests
make vendor        # Vendor dependencies (optional)
```

### Development

```bash
make lint          # Run linters (requires optional tools)
make spellcheck    # Run spell checker (requires typos)
make runchecks     # Run all checks
```

### Optional

```bash
make vendor        # Create vendor directory
make tools         # Install optional dev tools
make clean-vendor  # Remove vendor directory
```

## CI/CD

For continuous integration:

```bash
# Download dependencies
go mod download

# Verify they match go.sum
go mod verify

# Build
go build ./...

# Test
go test ./...
```

No vendoring or custom tool installation required.

## Troubleshooting

### "Module not found"

```bash
go mod download
go mod tidy
```

### "go.sum mismatch"

```bash
# Regenerate checksums
rm go.sum
go mod tidy
```

### "Old version of dependency"

```bash
# Update specific module
go get -u github.com/some/package@latest
go mod tidy

# Or update all
go get -u ./...
go mod tidy
```

### "Want to use older Go version"

Edit `go.mod`:

```
module github.com/sdoque/mbaigo

go 1.21  // Change to your version
```

Then:

```bash
go mod tidy
```

## Comparison: Before vs After

### Before (Non-standard)

```bash
make installpkgs   # Custom script that mixes deps with tools
```

Problems:
- Confuses dependencies with development tools
- Not following Go conventions
- Harder for newcomers to understand

### After (Standard Go)

```bash
# Build (downloads dependencies automatically)
make build

# Optional dev tools (separate)
make tools         # Clearly labeled as optional
```

Benefits:
- Follows Go best practices
- Dependencies downloaded automatically
- No confusing "deps" command
- Clear separation of concerns
- Standard commands everyone knows
- Works with all Go tooling

## Best Practices

1. **Always commit go.mod and go.sum**
2. **Run `go mod tidy` before committing**
3. **Don't commit vendor/ for libraries** (optional for apps)
4. **Use `go get` for new dependencies**, not manual edits
5. **Verify deps periodically**: `go mod verify`
6. **Check for vulnerabilities**: `govulncheck ./...` (if installed)

## References

- [Go Modules Reference](https://go.dev/ref/mod)
- [Module maintenance](https://go.dev/doc/modules/managing-dependencies)
- [go.mod file reference](https://go.dev/doc/modules/gomod-ref)
