# Docker Deployment for Three Sensors Example

This document provides detailed instructions for running the three-sensors example with Docker.

## Overview

The Docker setup includes:
- **ESR** (Service Registry) - Port 20102
- **Orchestrator** (Service Discovery) - Port 20103
- **Three Sensors Application** - Port 8080
  - Temperature Sensor
  - Pressure Sensor
  - Controller

All services are built from source (no pre-built images) and run in an isolated Docker network.

## Prerequisites

- Docker Engine 20.10+
- Docker Compose 1.29+
- No other services running on ports 8080, 20102, 20103

## Quick Start

### Option 1: Using Docker Compose

```bash
cd /path/to/mbaigo/examples/three-sensors
docker-compose up --build
```

### Option 2: Using Makefile (Recommended)

```bash
cd /path/to/mbaigo/examples/three-sensors
make docker-up
```

Both commands will:
1. Build mbaigo CLI from source
2. Build three-sensors application
3. Start ESR container
4. Start Orchestrator container
5. Start application container
6. Configure networking and dependencies

## Verifying the Deployment

### 1. Check Container Status

```bash
docker-compose ps
```

Expected output:
```
NAME                       STATUS    PORTS
three-sensors-esr          Up        0.0.0.0:20102->20102/tcp
three-sensors-orchestrator Up        0.0.0.0:20103->20103/tcp
three-sensors-app          Up        0.0.0.0:8080->8080/tcp
```

### 2. Check ESR Status

```bash
curl http://localhost:20102/serviceregistrar/registry/status
```

Expected: `lead Service Registrar since...`

### 3. View Registered Services

```bash
curl http://localhost:20102/serviceregistrar/registry/query
```

Should show 3 registered services: temperature, pressure, control

### 4. Test Service Endpoints

**Temperature:**
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

**Pressure:**
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

**Controller:**
```bash
curl http://localhost:8080/MySystem/Controller1/control
```

Expected response:
```json
{
  "status": true,
  "timestamp": "2025-11-13T10:00:00Z"
}
```

## Container Details

### ESR (Service Registry)

```yaml
Container: three-sensors-esr
Image: Built from ../../Dockerfile
Port: 20102
Volume: esr-data:/app/systems/esr
Health Check: wget http://localhost:20102/serviceregistrar/registry/status
```

**Purpose:** Tracks all available services in the local cloud

### Orchestrator

```yaml
Container: three-sensors-orchestrator
Image: Built from ../../Dockerfile
Port: 20103
Volume: orchestrator-data:/app/systems/orchestrator
Health Check: wget http://localhost:20103/orchestrator/orchestration/status
Depends On: esr (healthy)
```

**Purpose:** Provides service discovery and orchestration

### Application

```yaml
Container: three-sensors-app
Image: Built from ./Dockerfile
Port: 8080
Volume: app-data:/app
Health Check: wget http://localhost:8080/MySystem/TempSensor1/temperature
Depends On: esr (healthy), orchestrator (healthy)
```

**Purpose:** Runs the three sensor assets (temperature, pressure, controller)

## Common Operations

### View Logs

```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f esr
docker-compose logs -f orchestrator
docker-compose logs -f application

# Using Makefile
make docker-logs
```

### Stop Services

```bash
# Stop containers (preserve volumes)
docker-compose down

# Using Makefile
make docker-down
```

### Restart Services

```bash
# Restart all
docker-compose restart

# Restart specific service
docker-compose restart application

# Using Makefile
make docker-restart
```

### Rebuild and Restart

```bash
# After code changes
docker-compose up -d --build

# Or rebuild specific service
docker-compose up -d --build application
```

### Clean Everything

```bash
# Stop and remove volumes
docker-compose down -v

# Using Makefile
make docker-clean
```

## Networking

All containers run in the `three-sensors-network` bridge network.

**Internal communication:**
- ESR: `http://esr:20102`
- Orchestrator: `http://orchestrator:20103`
- Application: `http://application:8080`

**External access (from host):**
- ESR: `http://localhost:20102`
- Orchestrator: `http://localhost:20103`
- Application: `http://localhost:8080`

## Volumes

Persistent data is stored in Docker volumes:

```bash
# List volumes
docker volume ls | grep three-sensors

# Inspect a volume
docker volume inspect three-sensors_esr-data

# Remove all volumes (after docker-compose down)
docker volume rm three-sensors_esr-data \
                 three-sensors_orchestrator-data \
                 three-sensors_app-data
```

## Troubleshooting

### Services Won't Start

**Check logs:**
```bash
docker-compose logs esr
docker-compose logs orchestrator
docker-compose logs application
```

**Common issues:**
- Port conflicts: Another service is using 8080, 20102, or 20103
- Docker daemon not running
- Insufficient disk space

### Port Conflicts

**Find what's using the port:**
```bash
lsof -ti:20102
lsof -ti:20103
lsof -ti:8080
```

**Change ports in docker-compose.yml:**
```yaml
ports:
  - "8102:20102"  # Use 8102 externally instead of 20102
```

### Health Checks Failing

**Disable health checks temporarily:**
Edit `docker-compose.yml` and comment out healthcheck sections.

**Increase timeout:**
```yaml
healthcheck:
  timeout: 10s  # Increase from 3s
  retries: 20   # Increase from 10
```

### Service Not Registering

**Check ESR is reachable from application:**
```bash
docker-compose exec application wget -O- http://esr:20102/serviceregistrar/registry/status
```

**Check logs for registration errors:**
```bash
docker-compose logs application | grep "registration"
```

### Rebuild After Code Changes

```bash
# Stop everything
docker-compose down

# Rebuild with no cache
docker-compose build --no-cache

# Start again
docker-compose up -d
```

## Performance

**Container resource usage:**
```bash
docker stats
```

Typical usage:
- ESR: ~50MB RAM, <1% CPU
- Orchestrator: ~50MB RAM, <1% CPU
- Application: ~60MB RAM, <1% CPU

**Total overhead:** ~160MB RAM, negligible CPU when idle

## Production Deployment

For production, consider:

1. **Use specific image tags**
2. **Set resource limits:**
   ```yaml
   deploy:
     resources:
       limits:
         cpus: '0.5'
         memory: 512M
   ```
3. **Enable TLS/mTLS**
4. **Use secrets for sensitive data**
5. **Configure proper logging**
6. **Set up monitoring**
7. **Use container orchestration (Kubernetes, Docker Swarm)**
8. **Implement backup strategy for volumes**

## Development Workflow

### 1. Make Code Changes

Edit Go files in the three-sensors directory.

### 2. Rebuild Application

```bash
docker-compose up -d --build application
```

### 3. View Logs

```bash
docker-compose logs -f application
```

### 4. Test Changes

```bash
curl http://localhost:8080/MySystem/TempSensor1/temperature
```

### 5. Iterate

Repeat steps 1-4 until satisfied.

## Alternative: Running from Root

You can also run the complete setup from the mbaigo root:

```bash
cd /path/to/mbaigo
docker-compose up --build
```

This uses the main `docker-compose.yml` which is identical but with slightly different container names.

## See Also

- [Main README](README.md) - Complete three-sensors documentation
- [Root DOCKER.md](../../DOCKER.md) - Docker documentation for entire project
- [Architecture](../../docs/ARCHITECTURE.MD) - System architecture details
- [Main README](../../README.md) - Project overview
