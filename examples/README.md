# mbaigo Examples

This directory contains example applications demonstrating how to use the mbaigo library.

## Available Examples

### [Simple](./simple/)
The most basic example showing:
- Creating a system with one unit asset
- Implementing the UnitAsset interface
- Providing a simple HTTP service
- Handling graceful shutdown

**Complexity**: Beginner
**Prerequisites**: None
**Run time**: < 1 minute to understand

### [Consumer-Provider](./consumer-provider/)
A complete service-oriented example showing:
- Service registration with Service Registrar
- Service discovery via Orchestrator
- Service consumption between systems
- Form-based data exchange

**Complexity**: Intermediate
**Prerequisites**: Running Arrowhead Core Systems
**Run time**: 5-10 minutes to set up

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
2. Move to **consumer-provider** to see service interaction
3. Check the [documentation](../docs/) for detailed API reference

## Creating Your Own System

Use these examples as templates for your own Arrowhead systems:

1. Copy an example directory
2. Modify the UnitAsset implementation for your domain
3. Update service definitions
4. Run and test

See the main [README](../README.md) for complete API documentation.
