# Simple Randomizer Example

This example demonstrates the most basic Arrowhead system using mbaigo:

- Creates a system with a single unit asset
- Provides a simple HTTP service that returns random numbers
- Shows the essential UnitAsset interface implementation
- **Standalone application - no core systems required or used**
- Perfect for learning the basics without ESR/Orchestrator complexity

## Quick Start

### Docker (Recommended)

```bash
cd examples/simple
docker-compose up --build
```

This automatically:
- Builds the application from source
- Starts the randomizer service on port 8080
- Configures health checks

**Test the service:**
```bash
curl http://localhost:8080/RandomizerSystem/randomizer/random
```

**Stop the service:**
```bash
docker-compose down
```

### Native Go

```bash
cd examples/simple
go run main.go
```

## Testing the Service

In another terminal:

```bash
curl http://localhost:8080/RandomizerSystem/randomizer/random
```

Expected response:
```json
{
  "version": "1a",
  "typename": "SignalA",
  "value": 42.123,
  "unit": "float64"
}
```

## Configuration

The `systemconfig.json` defines:
- System name: RandomizerSystem
- HTTP port: 8080
- Asset: randomizer with random number service
- Traits: min/max value range (0.0 to 100.0)

## What This Demonstrates

1. **System Creation**: Initialize an Arrowhead system with context
2. **Husk Configuration**: Set up middleware with protocol/port bindings
3. **UnitAsset Implementation**: Create a custom asset implementing the interface
4. **Service Definition**: Define services without registration (standalone mode)
5. **HTTP Serving**: Handle incoming service requests
6. **Graceful Shutdown**: Clean shutdown on SIGINT/SIGTERM

**Note**: This example runs in standalone mode without core systems (ESR/Orchestrator). You may see log messages about "failed to find lead registrar" - these are harmless and can be ignored. The service works perfectly without registration. See the `three-sensors` example for full service registration and discovery.

## Files

- **main.go** - Application entry point and randomizer implementation
- **systemconfig.json** - System configuration
- **Dockerfile** - Docker build configuration
- **docker-compose.yml** - Docker orchestration
