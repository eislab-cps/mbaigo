# mbaigo

**Build Arrowhead Framework systems in Go**

mbaigo is a Go library for building service-oriented IoT and cyber-physical systems using the [Arrowhead Framework](https://arrowhead.eu/). Create distributed systems with automatic service registration and discovery in minutes.

## Quick Start (3 minutes)

### 1. Install CLI

```bash
git clone https://github.com/eislab-cps/mbaigo.git
cd mbaigo
make build
```

### 2. Start Core Systems

```bash
# Terminal 1: Start Service Registry
./bin/mbaigo core start esr

# Terminal 2: Start Orchestrator
./bin/mbaigo core start orchestrator
```

### 3. Run Example Application

```bash
# Terminal 3: Start sensors
cd examples/three-sensors
go run *.go
```

Test it:
```bash
curl http://localhost:8080/MySystem/TempSensor1/temperature
# Response: {"value":25,"unit":"celsius",...}
```

**That's it!** You now have a complete Arrowhead local cloud with service discovery.

## CLI Features

The `mbaigo` CLI provides everything you need:

```bash
# Core Systems Management
./bin/mbaigo core start esr           # Start Service Registry
./bin/mbaigo core start orchestrator  # Start Orchestrator
./bin/mbaigo core list                # List available core systems

# Project Management
./bin/mbaigo init --name MySystem     # Create new project
./bin/mbaigo config show              # View configuration
./bin/mbaigo config validate          # Validate config
./bin/mbaigo system info              # Show system details

# Code Generation
./bin/mbaigo generate asset --name Sensor  # Generate asset template
./bin/mbaigo generate example              # Generate complete example

# Service Management
./bin/mbaigo service list             # List all services
./bin/mbaigo service add              # Add service to asset
```

## Create Your Own System

```bash
# 1. Create project directory
mkdir my-sensor-system && cd my-sensor-system

# 2. Initialize with CLI
~/mbaigo/bin/mbaigo init --name MySensorSystem --cloud Factory

# 3. Generate asset
~/mbaigo/bin/mbaigo generate asset --name TemperatureSensor

# 4. Setup and run
go mod init my-sensor-system
go mod tidy
go run *.go
```

Your system will start with auto-generated endpoints!

## Core Concepts

**System Hierarchy:**
```
System
├── Host        (physical device)
├── Husk        (security/middleware)
└── UnitAssets  (sensors, actuators, controllers)
    ├── Services    (provides to others)
    └── Cervices    (consumes from others)
```

**Service Flow:**
1. Providers register services with ESR
2. Consumers query Orchestrator for services
3. Direct HTTP/HTTPS communication between systems

## Project Structure

```
mbaigo/
├── cmd/mbaigo/          # CLI tool
├── pkg/                 # Public library
│   ├── components/      # System structure (System, Host, Husk, UnitAsset)
│   ├── forms/           # Data exchange schemas (SignalA, ServiceRecord, etc.)
│   └── usecases/        # Business logic (registration, discovery, consumption)
├── examples/            # Working examples
│   └── three-sensors/   # Temperature, pressure, controller example
└── docs/                # Detailed documentation
```

## Examples

### Three Sensors Example

Complete multi-process service-oriented architecture:

```bash
cd examples/three-sensors
go run *.go  # Single process

# Or run as 5 separate processes:
# Terminal 1: ./bin/mbaigo core start esr
# Terminal 2: ./bin/mbaigo core start orchestrator
# Terminal 3: go run main.go temperaturesensor_asset.go
# Terminal 4: go run main.go pressuresensor_asset.go
# Terminal 5: go run main.go controller_asset.go
```

See [examples/three-sensors/README.md](./examples/three-sensors/README.md) for details.

## What You Get with the CLI

- **⚡ Instant Core Systems** - Embedded ESR and Orchestrator (no separate installation)
- **🚀 Code Generation** - Generate asset templates and examples
- **🔧 Project Management** - Initialize, configure, validate projects
- **📊 System Monitoring** - View system info, list services
- **🎯 Zero Config Start** - Default configs for immediate development

## Development

```bash
make build         # Build CLI
make test          # Run tests
make lint          # Run linters
make runchecks     # Run all checks
```

## Documentation

- **[Getting Started Guide](./docs/GETTING-STARTED.MD)** - Detailed tutorial
- **[Architecture](./docs/ARCHITECTURE.MD)** - System design with diagrams
- **[CLI Reference](./docs/CLI-REFERENCE.MD)** - Complete CLI documentation
- **[Examples](./examples/)** - Working example applications

## Key Features

- **Embedded Core Systems** - ESR and Orchestrator built into 10MB binary
- **Service-Oriented** - Automatic registration and discovery
- **Configuration-Driven** - JSON-based system configuration
- **Secure** - Mutual TLS with X.509 certificates
- **Extensible** - Interface-based design for custom assets
- **Standards-Based** - Implements Arrowhead Framework patterns

## Requirements

- Go 1.21 or later
- No external dependencies for basic usage
- (Optional) Arrowhead Core Systems for multi-cloud deployments

## Status

⚠️ **Early Development** - The code is in early stages and not production ready.

## Use Cases

- Industrial IoT systems
- Smart building automation
- Distributed sensor networks
- Edge computing applications
- System-of-systems integration

## Learn More

- [Arrowhead Framework](https://arrowhead.eu/)
- [Arrowhead GitHub](https://github.com/eclipse-arrowhead)
- [Eclipse Arrowhead Documentation](https://eclipse-arrowhead.github.io/core-java-spring/)

## License

[To be determined]

---

**Made with mbaigo** - Rapid development of distributed cyber-physical systems
