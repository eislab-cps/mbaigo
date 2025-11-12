# Building mbaigo

## Quick Start

```bash
# Build everything to ./bin/
make build

# Run the CLI
./bin/mbaigo version
```

## Build Targets

### Build CLI Only

```bash
# Using Make (recommended)
make build

# Or manually
mkdir -p bin
go build -o bin/mbaigo ./cmd/mbaigo
```

The binary will be created at: `./bin/mbaigo`

### Build All Packages

```bash
go build ./...
```

### Install Globally

```bash
# From source
go install github.com/sdoque/mbaigo/cmd/mbaigo@latest

# Or from local build
sudo cp bin/mbaigo /usr/local/bin/
```

## Binary Location

**All binaries are built into `./bin/`** which is excluded from Git via `.gitignore`.

| What | Location | Committed? |
|------|----------|------------|
| CLI binary | `./bin/mbaigo` | ❌ No (in .gitignore) |
| Source code | `./cmd/mbaigo/` | ✅ Yes |
| Dependencies | `go.mod`, `go.sum` | ✅ Yes |
| Vendor (optional) | `./vendor/` | ❌ No (in .gitignore) |

## Build Options

### With Debug Symbols

```bash
go build -o bin/mbaigo ./cmd/mbaigo
```

### Without Debug Symbols (smaller binary)

```bash
go build -ldflags="-s -w" -o bin/mbaigo ./cmd/mbaigo
```

### Cross-Compilation

```bash
# Linux AMD64
GOOS=linux GOARCH=amd64 go build -o bin/mbaigo-linux-amd64 ./cmd/mbaigo

# Linux ARM64
GOOS=linux GOARCH=arm64 go build -o bin/mbaigo-linux-arm64 ./cmd/mbaigo

# Windows AMD64
GOOS=windows GOARCH=amd64 go build -o bin/mbaigo-windows-amd64.exe ./cmd/mbaigo

# macOS Intel
GOOS=darwin GOARCH=amd64 go build -o bin/mbaigo-darwin-amd64 ./cmd/mbaigo

# macOS Apple Silicon
GOOS=darwin GOARCH=arm64 go build -o bin/mbaigo-darwin-arm64 ./cmd/mbaigo
```

## Development

### Build and Test

```bash
# Build (automatically downloads dependencies)
make build

# Test
make test

# All checks
make runchecks
```

### Clean Build

```bash
# Remove binaries and temporary files
make clean

# Rebuild
make build
```

## Makefile Targets

```bash
make help          # Show all available targets
make build         # Build CLI to ./bin/mbaigo
make test          # Run tests
make lint          # Run linters
make vendor        # Vendor dependencies (optional)
make clean         # Remove binaries and temp files
```

**Note:** Dependencies are downloaded automatically by Go. No separate `make deps` needed.

## Dependencies

Dependencies are managed via `go.mod`:

```bash
# View dependencies
go list -m all

# Update dependencies
go get -u ./...
go mod tidy

# Vendor (optional)
make vendor
```

See [DEPENDENCY-MANAGEMENT.md](./DEPENDENCY-MANAGEMENT.md) for details.

## Troubleshooting

### "Permission denied" when running binary

```bash
chmod +x bin/mbaigo
```

### Binary not found

Make sure you're in the project root:

```bash
pwd  # Should show .../mbaigo
ls bin/mbaigo  # Should exist
```

Or add to PATH:

```bash
export PATH=$PATH:$(pwd)/bin
mbaigo version
```

### Old binary cached

```bash
make clean
make build
```

## CI/CD

For continuous integration:

```bash
# Install dependencies
go mod download

# Build
mkdir -p bin
go build -o bin/mbaigo ./cmd/mbaigo

# Test
go test ./...
```

## Size Optimization

Compare binary sizes:

```bash
# Standard build
go build -o bin/mbaigo-debug ./cmd/mbaigo
ls -lh bin/mbaigo-debug

# Stripped build
go build -ldflags="-s -w" -o bin/mbaigo-stripped ./cmd/mbaigo
ls -lh bin/mbaigo-stripped

# With UPX compression (if installed)
upx --best --lzma bin/mbaigo-stripped
```

## Build for Distribution

Create release binaries for all platforms:

```bash
#!/bin/bash
# build-release.sh

VERSION="0.1.0"
PLATFORMS=("linux/amd64" "linux/arm64" "darwin/amd64" "darwin/arm64" "windows/amd64")

mkdir -p dist

for PLATFORM in "${PLATFORMS[@]}"; do
    GOOS="${PLATFORM%/*}"
    GOARCH="${PLATFORM#*/}"
    OUTPUT="dist/mbaigo-${VERSION}-${GOOS}-${GOARCH}"

    if [ "$GOOS" = "windows" ]; then
        OUTPUT="${OUTPUT}.exe"
    fi

    echo "Building for $GOOS/$GOARCH..."
    GOOS=$GOOS GOARCH=$GOARCH go build -ldflags="-s -w" -o "$OUTPUT" ./cmd/mbaigo

    if [ $? -eq 0 ]; then
        echo "✓ Built: $OUTPUT"
    else
        echo "✗ Failed to build for $GOOS/$GOARCH"
    fi
done
```

Run with:

```bash
chmod +x build-release.sh
./build-release.sh
ls -lh dist/
```
