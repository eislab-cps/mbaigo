# Makefile for mbaigo
# Uses scripts in /scripts directory following Go standard project layout

.PHONY: help build test lint spellcheck runchecks analyse deps vendor tools clean

# Default target
help:
	@echo "mbaigo - Makefile targets:"
	@echo ""
	@echo "Development:"
	@echo "  make deps          - Download Go module dependencies"
	@echo "  make vendor        - Vendor dependencies (optional)"
	@echo "  make build         - Build all packages and CLI"
	@echo "  make test          - Run tests with coverage"
	@echo "  make lint          - Run linters and static analysis"
	@echo "  make spellcheck    - Run spell checker"
	@echo "  make runchecks     - Run all checks (test + lint + spellcheck)"
	@echo ""
	@echo "Analysis:"
	@echo "  make analyse       - Generate detailed coverage report"
	@echo ""
	@echo "Tools (optional):"
	@echo "  make tools         - Install optional development tools"
	@echo ""
	@echo "Cleanup:"
	@echo "  make clean         - Clean up temporary files"
	@echo "  make clean-vendor  - Remove vendor directory"
	@echo ""
	@echo "See scripts/README.md for more details"

# Download Go module dependencies
deps:
	@echo "Downloading Go module dependencies..."
	@go mod download
	@go mod verify
	@echo "✓ Dependencies downloaded and verified"

# Vendor dependencies (optional, creates vendor/ directory)
vendor:
	@echo "Vendoring dependencies..."
	@go mod vendor
	@echo "✓ Dependencies vendored to ./vendor"

# Build all packages and the CLI
build:
	@echo "Building packages..."
	@go build ./pkg/... ./internal/... ./cmd/...
	@echo "Building CLI..."
	@mkdir -p bin
	@go build -o bin/mbaigo ./cmd/mbaigo
	@echo "✓ Build complete: bin/mbaigo"

# Run tests and log the test coverage
test:
	@./scripts/test.sh

# Run source code linters and catch common errors
lint:
	@./scripts/lint.sh

# Run spell checker on code and comments
spellcheck:
	@./scripts/spellcheck.sh

# Run all checks
runchecks: test lint spellcheck
	@echo ""
	@echo "✓ All checks passed!"

# Generate detailed coverage report
analyse:
	@./scripts/coverage-report.sh

# Install optional development tools (gosec, staticcheck, etc.)
tools:
	@./scripts/install-tools.sh

# Clean up temporary files
clean:
	@go clean
	@rm -f .cover.out cover.html
	@rm -rf bin/
	@rm -f pkg/**/systemconfig.json
	@rm -f examples/**/systemconfig.json
	@rm -f tests/systemconfig.json
	@echo "✓ Cleaned up temporary files"

# Remove vendor directory
clean-vendor:
	@rm -rf vendor
	@echo "✓ Removed vendor directory"
