# mbaigo

**Build Arrowhead Framework systems in Go**

mbaigo is a Go library for building service-oriented IoT and cyber-physical systems using the [Arrowhead Framework](https://arrowhead.eu/). Create distributed systems with automatic service registration and discovery in minutes.

## Quick Start

### Fastest Path: Simple Example (No Setup Required)

Get started immediately with a standalone example:

```bash
git clone https://github.com/eislab-cps/mbaigo.git
cd mbaigo/examples/simple
go run main.go
```

Test it:
```bash
curl http://localhost:8080/RandomizerSystem/randomizer/random
# Response: {"value":42.5,"unit":"float64",...}
```

**Perfect for:** Learning the basics, understanding the UnitAsset interface

### Full SOA Experience: Three Sensors with Docker

Experience complete service registration and discovery:

```bash
git clone https://github.com/eislab-cps/mbaigo.git
cd mbaigo/examples/three-sensors
docker-compose up --build
```

Test it:
```bash
curl http://localhost:8080/MySystem/TempSensor1/temperature
# Response: {"value":25,"unit":"celsius",...}
```

**Perfect for:** Understanding service-oriented architecture, production-like deployment

### Native Setup (For CLI Development)

If you want to develop with the CLI or run core systems natively:

```bash
# 1. Build CLI
git clone https://github.com/eislab-cps/mbaigo.git
cd mbaigo
make build

# Optional: Install globally
sudo make install

# 2. Start core systems (2 separate terminals from mbaigo root)
# Terminal 1:
mbaigo core start esr

# Terminal 2:
mbaigo core start orchestrator

# 3. Run example (Terminal 3)
cd examples/three-sensors
go run *.go
```

**Notes:**
- Core systems must run from mbaigo root directory
- First run creates config files and exits - run again to start
- Use `./bin/mbaigo` if not installed globally

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
# 1. Start core systems (in separate terminals, from mbaigo root)
cd /path/to/mbaigo
mbaigo core start esr         # Terminal 1

cd /path/to/mbaigo
mbaigo core start orchestrator # Terminal 2

# 2. Create project directory (in Terminal 3)
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

## Docker Deployment

Start core systems only:
```bash
cd mbaigo
docker-compose up --build
```

Or start a complete example with core systems + application:
```bash
cd mbaigo/examples/three-sensors
docker-compose up --build
```

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
1. Providers register services with ESR (Service Registry)
2. Consumers query Orchestrator for services
3. Direct HTTP/HTTPS communication between systems

## Project Structure

```
mbaigo/
├── cmd/mbaigo/          # CLI tool
├── pkg/                 # Public library
│   ├── components/      # System structure (System, Host, Husk, UnitAsset)
│   ├── forms/           # Data exchange schemas (SignalA, ServiceRecord, etc.)
│   ├── usecases/        # Business logic (registration, discovery, consumption)
│   └── core/            # Embedded core systems (ESR, Orchestrator)
├── examples/            # Working examples
│   ├── simple/          # Basic standalone example (no core systems)
│   └── three-sensors/   # Full SOA example with service discovery
├── docker/              # Pre-configured Docker configs for core systems
└── docs/                # Detailed documentation
```

## Examples

### 1. Simple Randomizer (Standalone - No Core Systems)

**Perfect for beginners** - Standalone HTTP service demonstrating the basics:

```bash
cd examples/simple
go run main.go

# Or with Docker:
docker-compose up --build
```

Test it:
```bash
curl http://localhost:8080/RandomizerSystem/randomizer/random
```

**What it demonstrates:**
- Basic system setup without core systems
- UnitAsset interface implementation
- HTTP service handling

See [examples/simple/README.md](./examples/simple/README.md) for details.

### 2. Three Sensors (Full Service-Oriented Architecture)

**Complete example** with service registration, discovery, and consumption:

```bash
# Docker (Recommended - all-in-one):
cd examples/three-sensors
docker-compose up --build

# Or Native Go (requires 3 terminals):
# Terminal 1 (from mbaigo root): mbaigo core start esr
# Terminal 2 (from mbaigo root): mbaigo core start orchestrator
# Terminal 3: cd examples/three-sensors && go run *.go
```

Test it:
```bash
curl http://localhost:8080/MySystem/TempSensor1/temperature
curl http://localhost:8080/MySystem/PressureSensor1/pressure
curl http://localhost:8080/MySystem/Controller1/control
```

**What it demonstrates:**
- Service registration with ESR
- Service discovery through Orchestrator
- Multiple assets in one system
- Service consumption (Cervices)
- Automatic re-registration

See [examples/three-sensors/README.md](./examples/three-sensors/README.md) for details.

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
- **Docker Ready** - Complete Docker Compose setup with one-command deployment
- **Service-Oriented** - Automatic registration and discovery
- **Configuration-Driven** - JSON-based system configuration
- **Secure** - Mutual TLS with X.509 certificates
- **Extensible** - Interface-based design for custom assets
- **Standards-Based** - Implements Arrowhead Framework patterns

## Requirements

- Go 1.24 or later
- Docker (optional, for containerized deployment)
- No external dependencies for basic usage

## Status

**Early Development** - The code is in early stages and not production ready.

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
