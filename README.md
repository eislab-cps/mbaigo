# mbaigo

**Build Arrowhead Framework systems in Go**

mbaigo is a Go library for building service-oriented IoT and cyber-physical systems using the [Arrowhead Framework](https://arrowhead.eu/). Create distributed systems with automatic service registration and discovery in minutes.

## Quick Start (3 minutes)

### 1. Install CLI

```bash
git clone https://github.com/eislab-cps/mbaigo.git
cd mbaigo
make build

# Optional: Install globally
sudo make install
# Now you can use 'mbaigo' from anywhere instead of './bin/mbaigo'
```

### 2. Start Core Systems

```bash
# Terminal 1: Start Service Registry
mbaigo core start esr

# Terminal 2: Start Orchestrator
mbaigo core start orchestrator
```

**Note:** Use `mbaigo` if installed globally, or `./bin/mbaigo` if running from the build directory.

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
mbaigo core start esr           # Start Service Registry
mbaigo core start orchestrator  # Start Orchestrator
mbaigo core list                # List available core systems

# Project Management
mbaigo init --name MySystem     # Create new project
mbaigo config show              # View configuration
mbaigo config validate          # Validate config
mbaigo system info              # Show system details

# Code Generation
mbaigo generate asset --name Sensor  # Generate asset template
mbaigo generate example              # Generate complete example

# Service Management
mbaigo service list             # List all services
mbaigo service add              # Add service to asset
```

## Create Your Own System

```bash
# 1. Start core systems (in separate terminals)
mbaigo core start esr
mbaigo core start orchestrator

# 2. Create project directory
mkdir my-sensor-system && cd my-sensor-system

# 3. Initialize with CLI
mbaigo init --name MySensorSystem --cloud Factory

# 4. Generate asset
mbaigo generate asset --name TemperatureSensor

# 5. Setup and run
go mod init my-sensor-system
go mod tidy
go run *.go
```

Your system will register with ESR and be discoverable through the Orchestrator!

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

Complete service-oriented architecture demonstrating service registration, discovery, and consumption:

**Prerequisites:** Start core systems first (ESR and Orchestrator)

```bash
# Terminal 1: Start Service Registry
mbaigo core start esr

# Terminal 2: Start Orchestrator
mbaigo core start orchestrator

# Terminal 3: Start application systems
cd examples/three-sensors
go run *.go
```

The application will register services with ESR and discover them through the Orchestrator.

See [examples/three-sensors/README.md](./examples/three-sensors/README.md) for complete details.

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
- **[IEEE Paper](./paper/)** - Academic paper with performance evaluation

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
