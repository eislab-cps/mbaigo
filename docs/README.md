# mbaigo Documentation

Complete documentation for building Arrowhead Framework systems in Go.

## Quick Links

- **[Getting Started](./GETTING-STARTED.MD)** - Step-by-step tutorial from zero to running system
- **[Architecture](./ARCHITECTURE.MD)** - System design, diagrams, and patterns
- **[CLI Reference](./CLI-REFERENCE.MD)** - Complete command-line interface guide
- **[Main README](../README.md)** - Quick start and overview

## Documentation Structure

### Core Guides

| Document | Purpose | Audience |
|----------|---------|----------|
| [GETTING-STARTED.MD](./GETTING-STARTED.MD) | Hands-on tutorial | New users |
| [ARCHITECTURE.MD](./ARCHITECTURE.MD) | Design and diagrams | Developers |
| [CLI-REFERENCE.MD](./CLI-REFERENCE.MD) | CLI commands | All users |

### Package Documentation

- **Components** (`pkg/components`) - System structure (System, Host, Husk, UnitAsset)
- **Forms** (`pkg/forms`) - Data exchange schemas (ServiceRecord, SignalA, etc.)
- **Use Cases** (`pkg/usecases`) - Business logic (registration, discovery, consumption)

See [pkg.go.dev](https://pkg.go.dev/github.com/sdoque/mbaigo/pkg) for API reference.

## Getting Help

**New to mbaigo?**
1. Read the [main README](../README.md) for 3-minute quick start
2. Follow [GETTING-STARTED.MD](./GETTING-STARTED.MD) for detailed tutorial
3. Run the [examples](../examples/three-sensors/)

**Building a system?**
1. Use `./bin/mbaigo init` to create project
2. Check [CLI-REFERENCE.MD](./CLI-REFERENCE.MD) for commands
3. Read [ARCHITECTURE.MD](./ARCHITECTURE.MD) for design patterns

**Contributing?**
1. Understand the architecture
2. Run `make test` and `make lint`
3. Follow Go standard practices

## Examples

Working example applications in [/examples](../examples/):
- **three-sensors** - Complete multi-process SOA system

Each example has its own README with setup instructions.
