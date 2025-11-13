# Docker Deployment Guide

This guide explains how to run mbaigo with Docker and Docker Compose.

## Quick Start

Start the complete service-oriented architecture with one command:

```bash
docker-compose up --build
```

This will:
1. Build the mbaigo CLI from source
2. Start ESR (Service Registry) on port 20102
3. Start Orchestrator on port 20103
4. Start three-sensors example on port 8080
5. All services will auto-register and discover each other

## Architecture

```
┌─────────────────────────────────────────┐
│          Docker Network                 │
│                                         │
│  ┌──────────┐  ┌──────────────┐       │
│  │   ESR    │  │ Orchestrator │       │
│  │ :20102   │  │   :20103     │       │
│  └────┬─────┘  └──────┬───────┘       │
│       │               │                │
│       │  ┌────────────▼─────────┐     │
│       └─►│   Three Sensors      │     │
│          │      :8080           │     │
│          └──────────────────────┘     │
└─────────────────────────────────────────┘
         │           │
         ▼           ▼
    Host:20102  Host:8080
```

## Services

### ESR (Service Registry)
- **Port**: 20102
- **Health**: http://localhost:20102/serviceregistrar/registry/status
- **Config**: Stored in `esr-data` volume
- **Purpose**: Tracks all available services in the local cloud

### Orchestrator
- **Port**: 20103
- **Health**: http://localhost:20103/orchestrator/orchestration/status
- **Config**: Stored in `orchestrator-data` volume
- **Purpose**: Provides service discovery and orchestration

### Three Sensors Application
- **Port**: 8080
- **Config**: Stored in `three-sensors-data` volume
- **Services**:
  - Temperature: http://localhost:8080/MySystem/TempSensor1/temperature
  - Pressure: http://localhost:8080/MySystem/PressureSensor1/pressure
  - Controller: http://localhost:8080/MySystem/Controller1/control

## Common Commands

### Start all services
```bash
docker-compose up --build
```

### Start in detached mode (background)
```bash
docker-compose up -d --build
```

### View logs
```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f esr
docker-compose logs -f orchestrator
docker-compose logs -f three-sensors
```

### Stop all services
```bash
docker-compose down
```

### Stop and remove volumes (clean slate)
```bash
docker-compose down -v
```

### Restart a specific service
```bash
docker-compose restart esr
docker-compose restart orchestrator
docker-compose restart three-sensors
```

### Rebuild a specific service
```bash
docker-compose up -d --build esr
```

## Testing the Deployment

### 1. Check Service Registry
```bash
curl http://localhost:20102/serviceregistrar/registry/query
```

Should show all registered services from the three-sensors application.

### 2. Test Temperature Service
```bash
curl http://localhost:8080/MySystem/TempSensor1/temperature
```

Expected response:
```json
{
  "value": 25.4,
  "unit": "celsius",
  "timestamp": "2025-11-13T10:00:00Z"
}
```

### 3. Test Pressure Service
```bash
curl http://localhost:8080/MySystem/PressureSensor1/pressure
```

Expected response:
```json
{
  "value": 1.025,
  "unit": "bar",
  "timestamp": "2025-11-13T10:00:00Z"
}
```

### 4. Check ESR Status
```bash
curl http://localhost:20102/serviceregistrar/registry/status
```

Expected response:
```
lead Service Registrar since 2025-11-13 10:00:00
```

## Networking

All services run in the `mbaigo-network` bridge network. Services can communicate using their service names:

- ESR is accessible at: `http://esr:20102`
- Orchestrator is accessible at: `http://orchestrator:20103`
- Three-sensors is accessible at: `http://three-sensors:8080`

From the host machine, services are accessible on localhost with mapped ports.

## Volumes

Persistent data is stored in Docker volumes:

- `esr-data`: ESR configuration and state
- `orchestrator-data`: Orchestrator configuration
- `three-sensors-data`: Application configuration

To inspect volumes:
```bash
docker volume ls
docker volume inspect mbaigo_esr-data
```

## Customization

### Environment Variables

You can customize service URLs by modifying the `environment` section in `docker-compose.yml`:

```yaml
three-sensors:
  environment:
    - ESR_URL=http://esr:20102/serviceregistrar/registry
    - ORCHESTRATOR_URL=http://orchestrator:20103/orchestrator/orchestration
```

### Port Mapping

Change external ports in the `ports` section:

```yaml
esr:
  ports:
    - "8102:20102"  # Access ESR on port 8102 instead of 20102
```

### Resource Limits

Add resource constraints to services:

```yaml
esr:
  deploy:
    resources:
      limits:
        cpus: '0.5'
        memory: 512M
      reservations:
        cpus: '0.25'
        memory: 256M
```

## Troubleshooting

### Service won't start

Check logs:
```bash
docker-compose logs esr
```

### Port already in use

Stop conflicting services or change port mapping in `docker-compose.yml`:
```yaml
ports:
  - "8102:20102"  # Use different external port
```

### Services can't communicate

Ensure all services are on the same network:
```bash
docker network inspect mbaigo_mbaigo-network
```

### Config issues

Reset volumes and restart:
```bash
docker-compose down -v
docker-compose up --build
```

### Health check failing

Wait for services to initialize (first run creates configs):
```bash
# Watch the startup sequence
docker-compose up
```

ESR and Orchestrator will exit on first run after creating configs, then restart automatically.

## Production Considerations

For production deployment:

1. **Use specific image tags** instead of building from source
2. **Enable TLS/mTLS** for secure communication
3. **Configure proper logging** with log aggregation
4. **Set resource limits** to prevent resource exhaustion
5. **Use secrets management** for sensitive configuration
6. **Implement monitoring** with health checks and metrics
7. **Configure restart policies** appropriately
8. **Use Docker Swarm or Kubernetes** for orchestration at scale

## Development Workflow

### 1. Make code changes

Edit source files in your local repository.

### 2. Rebuild affected services

```bash
# Rebuild everything
docker-compose up -d --build

# Or rebuild specific service
docker-compose up -d --build three-sensors
```

### 3. View logs
```bash
docker-compose logs -f three-sensors
```

### 4. Test changes
```bash
curl http://localhost:8080/MySystem/TempSensor1/temperature
```

## Multi-Architecture Builds

To build for different architectures (e.g., ARM for Raspberry Pi):

```bash
# Enable buildx
docker buildx create --use

# Build for multiple platforms
docker buildx build --platform linux/amd64,linux/arm64,linux/arm/v7 -t mbaigo:latest .
```

## See Also

- [Main README](README.md) - Project overview and native installation
- [Getting Started](docs/GETTING-STARTED.MD) - Detailed tutorial
- [Three Sensors Example](examples/three-sensors/README.md) - Example application details
- [Architecture](docs/ARCHITECTURE.MD) - System architecture and design
