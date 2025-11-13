# mbaigo Examples

This directory contains example applications demonstrating how to use the mbaigo library.

## Available Examples

### [Three Sensors](./three-sensors/)
A complete service-oriented architecture example showing:
- Temperature and pressure providers registering services
- Controller discovering and consuming services
- Service registration with ESR (Service Registry)
- Service discovery via Orchestrator
- Asset factory pattern for dynamic instantiation
- Configuration-driven asset creation
- Traits for runtime configuration

**Complexity**: Intermediate
**Prerequisites**: Running core systems (ESR and Orchestrator)
**Run time**: 5-10 minutes to set up
**Best for**: Learning complete SOA with service registration and discovery

### [Simple](./simple/)
The most basic example showing:
- Creating a system with one unit asset
- Implementing the UnitAsset interface
- Providing a simple HTTP service
- Handling graceful shutdown

**Complexity**: Beginner
**Prerequisites**: None
**Run time**: < 1 minute to understand
**Best for**: Understanding core concepts

### [Consumer-Provider](./consumer-provider/)
A complete service-oriented example showing:
- Service registration with Service Registrar
- Service discovery via Orchestrator
- Service consumption between systems
- Form-based data exchange

**Complexity**: Intermediate
**Prerequisites**: Running Arrowhead Core Systems
**Run time**: 5-10 minutes to set up
**Best for**: Understanding Arrowhead Framework integration

## Running Examples

Each example directory contains:
- `main.go` (or `*.go`) - The application code
- `README.md` - Specific instructions and explanations
- `systemconfig.json` - System configuration

**For full service-oriented architecture (recommended):**

⚠️ **IMPORTANT:** Core systems MUST run from **mbaigo root directory**!

```bash
# Terminal 1: Start ESR (from mbaigo root!)
cd /path/to/mbaigo
mbaigo core start esr

# Terminal 2: Start Orchestrator (from mbaigo root!)
cd /path/to/mbaigo
mbaigo core start orchestrator

# Terminal 3: Start example (from example directory)
cd /path/to/mbaigo/examples/<example-name>
go run *.go
```

**Notes:**
- Use `mbaigo` if installed globally (`sudo make install`), or `./bin/mbaigo` if running from build directory
- **Core systems from root, applications from their own directories** - this prevents config conflicts
- Some examples may work standalone without core systems, but you'll miss the service registration and discovery features

## Learning Path

**Recommended order:**

1. **Start with three-sensors** - The main example demonstrating complete service-oriented architecture with registration, discovery, and consumption
2. **Explore simple** (if available) - Understand basic system concepts without SOA complexity
3. **Try consumer-provider** (if available) - Advanced service interaction patterns
4. **Read the docs** - Check [documentation](../docs/) and [GETTING-STARTED.MD](../docs/GETTING-STARTED.MD) for detailed guides

**Quick start:** Jump straight to `three-sensors` - it has everything you need to understand mbaigo!

## Creating Your Own System

Use these examples as templates for your own Arrowhead systems:

1. Copy an example directory
2. Modify the UnitAsset implementation for your domain
3. Update service definitions
4. Run and test

See the main [README](../README.md) for complete API documentation.
