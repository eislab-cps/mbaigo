# mbaigo Examples

This directory contains example applications demonstrating how to use the mbaigo library.

## Available Examples

### [Three Sensors](./three-sensors/)
A practical multi-asset example showing:
- Temperature sensor with configurable min/max range
- Pressure sensor with simulated readings
- Controller with boolean status
- Asset factory pattern for dynamic instantiation
- Configuration-driven asset creation
- Traits for runtime configuration

**Complexity**: Beginner
**Prerequisites**: None (works standalone)
**Run time**: < 5 minutes to understand
**Best for**: Learning asset implementation and configuration

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
- `main.go` - The application code
- `README.md` - Specific instructions and explanations
- `systemconfig.json` - Generated configuration (after first run)

To run any example:

```bash
cd examples/<example-name>
go run main.go
```

## Learning Path

1. Start with **simple** to understand basic concepts
2. Explore **three-sensors** to learn asset configuration and factory pattern
3. Move to **consumer-provider** to see service interaction
4. Check the [documentation](../docs/) and [GETTING-STARTED.md](../GETTING-STARTED.md) for detailed guides

## Creating Your Own System

Use these examples as templates for your own Arrowhead systems:

1. Copy an example directory
2. Modify the UnitAsset implementation for your domain
3. Update service definitions
4. Run and test

See the main [README](../README.md) for complete API documentation.
