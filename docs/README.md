# mbaigo Documentation

Welcome to the mbaigo documentation. This directory contains comprehensive guides and references for using the mbaigo library.

## Available Documentation

### [Architecture](./architecture.md)
Complete architectural overview of mbaigo including:
- Core architecture principles
- Package organization
- Key architectural patterns
- Data flow and concurrency model
- Security considerations
- Performance characteristics

### [Use Cases](./usecases.md)
Detailed explanation of the use cases package:
- Authentication and certification
- Configuration management
- Service consumption and provision
- Service registration and discovery
- Cost handling
- Knowledge graph generation

### [Migration Guide](./migration-guide.md)
Guide for updating existing code after the Go Standard Project Layout reorganization:
- Import path changes
- Migration steps
- New features
- Breaking changes

## Quick Links

- [Main README](../README.md) - Project overview and getting started
- [Examples](../examples/) - Runnable example applications
- [Build Scripts](../scripts/) - Development and build scripts
- [API Reference](https://pkg.go.dev/github.com/sdoque/mbaigo) - Auto-generated API docs

## Package Documentation

### Components (`pkg/components`)
System structure models - the building blocks of Arrowhead systems:
- **System**: Root orchestrator
- **HostingDevice**: Physical/virtual machine representation
- **Husk**: Middleware and security layer
- **Service & Cervice**: Service definitions
- **UnitAsset**: Interface for domain implementations

### Forms (`pkg/forms`)
Data exchange schemas with versioning support:
- Service forms (registration, discovery)
- Signal forms (analog, digital)
- System forms (messages, costs)
- Certificate forms

### Use Cases (`pkg/usecases`)
Business logic and common system behaviors:
- Registration and service discovery
- Consumption and provision
- Configuration and authentication
- Documentation and knowledge graphs

## Getting Help

1. **Start with examples**: Check [/examples](../examples/) for working code
2. **Read the architecture**: Understand the design in [architecture.md](./architecture.md)
3. **Check use cases**: See how components work together in [usecases.md](./usecases.md)
4. **API reference**: Browse godoc for detailed API documentation

## Contributing to Documentation

When contributing documentation:
1. Use clear, concise language
2. Provide code examples where appropriate
3. Include diagrams for complex concepts
4. Keep documents focused on a single topic
5. Update this index when adding new documents

## Document Versioning

Documentation is versioned with the code. The current documentation corresponds to the latest version in the main branch.
