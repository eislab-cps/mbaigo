# Migration Guide: Reorganization to Go Standard Project Layout

This document explains the changes made to reorganize mbaigo according to the Go Standard Project Layout.

## What Changed

### Directory Structure

**Before:**
```
mbaigo/
├── components/
├── forms/
├── usecases/
├── tests/
├── go.mod
├── LICENSE
├── Makefile
└── README.md
```

**After:**
```
mbaigo/
├── pkg/
│   ├── components/
│   ├── forms/
│   └── usecases/
├── internal/            (for future private code)
├── examples/            (new)
├── docs/                (new)
├── scripts/             (new)
├── tests/
├── .github/workflows/
├── go.mod
├── LICENSE
├── Makefile
└── README.md
```

## Import Path Changes

All import paths now include the `/pkg/` prefix:

### Before
```go
import (
    "github.com/sdoque/mbaigo/components"
    "github.com/sdoque/mbaigo/forms"
    "github.com/sdoque/mbaigo/usecases"
)
```

### After
```go
import (
    "github.com/sdoque/mbaigo/pkg/components"
    "github.com/sdoque/mbaigo/pkg/forms"
    "github.com/sdoque/mbaigo/pkg/usecases"
)
```

## Migration Steps for Existing Code

### 1. Update Import Statements

Replace all imports in your code:

```bash
# Using sed (macOS)
find . -name "*.go" -exec sed -i '' 's|github.com/sdoque/mbaigo/components|github.com/sdoque/mbaigo/pkg/components|g' {} \;
find . -name "*.go" -exec sed -i '' 's|github.com/sdoque/mbaigo/forms|github.com/sdoque/mbaigo/pkg/forms|g' {} \;
find . -name "*.go" -exec sed -i '' 's|github.com/sdoque/mbaigo/usecases|github.com/sdoque/mbaigo/pkg/usecases|g' {} \;

# Using sed (Linux)
find . -name "*.go" -exec sed -i 's|github.com/sdoque/mbaigo/components|github.com/sdoque/mbaigo/pkg/components|g' {} \;
find . -name "*.go" -exec sed -i 's|github.com/sdoque/mbaigo/forms|github.com/sdoque/mbaigo/pkg/forms|g' {} \;
find . -name "*.go" -exec sed -i 's|github.com/sdoque/mbaigo/usecases|github.com/sdoque/mbaigo/pkg/usecases|g' {} \;
```

### 2. Update Dependencies

```bash
go mod tidy
```

### 3. Rebuild and Test

```bash
go build ./...
go test ./...
```

## New Features

### Examples Directory

Working example applications are now in `/examples`:
- `examples/simple/` - Basic system example
- `examples/consumer-provider/` - Service interaction example

Run examples:
```bash
cd examples/simple
go run main.go
```

### Documentation

Comprehensive documentation is now in `/docs`:
- `docs/ARCHITECTURE.MD` - System architecture
- `docs/USECASES.MD` - Use cases documentation
- `docs/MIGRATION-GUIDE.MD` - This guide

### Build Scripts

Build and test scripts are now in `/scripts`:
- `scripts/test.sh` - Run tests
- `scripts/lint.sh` - Run linters
- `scripts/coverage-report.sh` - Generate coverage
- `scripts/install-tools.sh` - Install dev tools

Use via Makefile:
```bash
make test
make lint
make runchecks
```

## Benefits of New Structure

### 1. Clear Public API

The `/pkg` directory clearly signals which packages are intended for external use.

### 2. Future Private Code

The `/internal` directory is ready for private implementation code that shouldn't be imported by external projects.

### 3. Better Examples

Runnable example applications help users understand how to use the library.

### 4. Organized Documentation

All documentation is centralized in `/docs` for easy reference.

### 5. Maintainable Scripts

Build and test logic is extracted from Makefile into versioned scripts.

### 6. Standard Layout

Following community standards makes the project more approachable for Go developers.

## Backwards Compatibility

**Important:** This reorganization introduces **breaking changes** to import paths.

If you maintain code using mbaigo:
1. Update all imports as shown above
2. Run `go mod tidy`
3. Test your code thoroughly

## Questions or Issues?

If you encounter problems during migration:
1. Check that all imports have been updated
2. Run `go mod tidy` to update dependencies
3. Verify tests pass with `go test ./...`
4. Review the examples in `/examples` for reference

## Timeline

This reorganization was completed on: 2025-11-11

All internal tests pass after reorganization. The public API remains unchanged - only import paths have changed.
