# Consumer-Provider Example

This example demonstrates a complete service-oriented interaction:

- **Provider System**: Offers a temperature reading service
- **Consumer System**: Discovers and consumes the temperature service
- Shows service registration, discovery, and consumption flow

## Prerequisites

You need running Arrowhead Core Systems:
- Service Registrar (default: https://localhost:8000)
- Orchestrator (default: https://localhost:8001)

## Running

### Start the Provider

```bash
cd examples/consumer-provider
go run provider/main.go
```

### Start the Consumer

In another terminal:

```bash
cd examples/consumer-provider
go run consumer/main.go
```

## What This Demonstrates

1. **Service Registration**: Provider registers with Service Registrar
2. **Service Discovery**: Consumer queries Orchestrator for providers
3. **Service Consumption**: Consumer makes HTTP request to provider
4. **Form Serialization**: Using SignalA_v1a for data exchange
5. **Mutual TLS**: Secure communication between systems (if configured)

## Configuration

Both systems will create `systemconfig.json` files in their directories. You can customize:
- System names
- Port numbers
- Core system URLs
- Service details
