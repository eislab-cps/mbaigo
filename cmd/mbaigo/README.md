# mbaigo CLI

Command-line tool for managing Arrowhead Framework systems.

## Installation

### From Source

```bash
go install github.com/eislab-cps/mbaigo/cmd/mbaigo@latest
```

### Build Locally

```bash
# Build to ./bin directory
make build

# Or manually
go build -o bin/mbaigo ./cmd/mbaigo
```

### Install to System Path

```bash
# macOS/Linux
sudo cp bin/mbaigo /usr/local/bin/

# Or add to your PATH
export PATH=$PATH:$(pwd)/bin
```

## Quick Start

```bash
# Initialize a new system
mbaigo init --name MySystem --cloud LocalCloud

# Validate configuration
mbaigo config validate

# View system info
mbaigo system info

# Generate an asset template
mbaigo generate asset --name TemperatureSensor

# Generate example application
mbaigo generate example > main.go
```

## Available Commands

- `mbaigo init` - Initialize a new system configuration
- `mbaigo config` - Configuration management (show, validate)
- `mbaigo system` - System management (info, start template)
- `mbaigo service` - Service management (list, add)
- `mbaigo cert` - Certificate management (generate-key, create-csr)
- `mbaigo generate` - Code generation (asset, example)
- `mbaigo version` - Show version information

## Documentation

- [CLI Reference](../../docs/CLI-REFERENCE.MD) - Complete command documentation
- [Getting Started](../../docs/GETTING-STARTED.MD) - Step-by-step tutorial
- [Examples](../../examples/) - Working examples

## Global Flags

```
  -c, --config string   Config file path (default "systemconfig.json")
  -j, --json            Output as JSON
      --no-color        Disable colored output
  -v, --verbose         Enable verbose output
```

## Examples

### Initialize System with Custom Ports

```bash
mbaigo init --name SensorSystem --cloud Factory --http-port 9090 --https-port 9443
```

### Interactive Mode

```bash
mbaigo init --interactive
```

### JSON Output for Scripting

```bash
mbaigo service list --json | jq '.services[] | select(.asset=="TempSensor")'
```

### Generate Certificates

```bash
mbaigo cert generate-key
mbaigo cert create-csr --common-name MySystem
```

### Add Service to Asset

```bash
mbaigo service add --asset MySensor --name Temperature
```

## Help

Get help for any command:

```bash
mbaigo --help
mbaigo init --help
mbaigo cert --help
```

## Development

The CLI is built with [Cobra](https://github.com/spf13/cobra) and located in `/cmd/mbaigo` following Go standard project layout.

Source code structure:

```
cmd/mbaigo/
├── main.go              # Entry point
└── README.md            # This file

internal/cli/
├── root.go              # Root command and version
├── init.go              # System initialization
├── config.go            # Configuration commands
├── system.go            # System management
├── service.go           # Service management
├── cert.go              # Certificate commands
├── generate.go          # Code generation
└── utils.go             # Helper functions
```

## License

See [LICENSE](../../LICENSE) in the repository root.
