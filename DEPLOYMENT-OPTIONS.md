# mbaigo Deployment Options

This guide compares different ways to deploy and run mbaigo systems.

## Quick Comparison

| Feature | Docker Compose | Native Go | Production K8s |
|---------|---------------|-----------|----------------|
| **Setup Time** | 1 minute | 5-10 minutes | 30+ minutes |
| **Prerequisites** | Docker only | Go + mbaigo CLI | Kubernetes cluster |
| **Isolation** | Full | None | Full |
| **Port Conflicts** | Avoided | Manual | Avoided |
| **Resource Usage** | ~160MB RAM | ~100MB RAM | Varies |
| **Best For** | Quick start, development | Learning, debugging | Production, scale |
| **Multi-host** | Single host | Single host | Multi-host |

## Option 1: Docker Compose (Recommended for Quick Start)

### Advantages

- **Fastest setup** - One command deployment
- **No manual configuration** - Everything automated
- **Isolated environment** - No system conflicts
- **Reproducible** - Same setup every time
- **Easy cleanup** - Remove everything with one command
- **Version control** - Infrastructure as code

### Disadvantages

- **Requires Docker** - ~500MB Docker Engine
- **Slightly higher RAM** - Extra container overhead
- **Single host only** - Can't span multiple machines
- **Learning curve** - Need to understand Docker basics

### When to Use

- Quick demos and testing
- Development environment
- CI/CD pipelines
- Onboarding new developers
- Conference presentations

### When NOT to Use

- Production multi-host deployments
- Embedded systems without Docker
- Resource-constrained devices
- Debugging low-level networking

### Quick Start

**From project root:**
```bash
cd /path/to/mbaigo
docker-compose up --build
# or
make docker-up
```

**From three-sensors example:**
```bash
cd /path/to/mbaigo/examples/three-sensors
docker-compose up --build
# or
make docker-up
```

**Test:**
```bash
curl http://localhost:8080/MySystem/TempSensor1/temperature
```

**Stop:**
```bash
docker-compose down
# or
make docker-down
```

### Documentation

- [Root DOCKER.md](./DOCKER.md) - Complete Docker guide
- [Three-sensors DOCKER-README.md](./examples/three-sensors/DOCKER-README.md) - Example-specific guide

## Option 2: Native Go (Recommended for Learning)

### Advantages

- **Lower resource usage** - No container overhead
- **Faster iteration** - No rebuild for code changes
- **Better debugging** - Direct Go tools access
- **Educational** - See how everything works
- **No Docker dependency** - Works anywhere Go runs

### Disadvantages

- **Manual setup** - Multiple terminal windows
- **Port conflicts** - Need to manage manually
- **Directory discipline** - Must run from correct directories
- **More cleanup** - Multiple processes to kill
- **Config management** - Need to track config files

### When to Use

- Learning mbaigo architecture
- Developing new features
- Debugging core systems
- Performance profiling
- Running on embedded devices

### When NOT to Use

- Quick demos
- CI/CD pipelines
- Multiple concurrent projects
- Unfamiliar with terminal multiplexing

### Quick Start

**Build:**
```bash
cd /path/to/mbaigo
make build
sudo make install  # Optional: Install globally
```

**Terminal 1 - ESR:**
```bash
cd /path/to/mbaigo  # MUST be root!
mbaigo core start esr
```

**Terminal 2 - Orchestrator:**
```bash
cd /path/to/mbaigo  # MUST be root!
mbaigo core start orchestrator
```

**Terminal 3 - Application:**
```bash
cd /path/to/mbaigo/examples/three-sensors
go run *.go
```

**Test:**
```bash
curl http://localhost:8080/MySystem/TempSensor1/temperature
```

**Stop:**
- Press `Ctrl+C` in each terminal (reverse order)

### Documentation

- [README.md](./README.md) - Main documentation
- [TESTING.md](./TESTING.md) - Testing guide
- [examples/three-sensors/README.md](./examples/three-sensors/README.md) - Example guide

## Option 3: Production Kubernetes (Future)

### Advantages

- **Multi-host** - Span multiple servers
- **High availability** - Automatic failover
- **Auto-scaling** - Scale based on load
- **Load balancing** - Built-in
- **Service mesh** - Advanced networking
- **Monitoring** - Prometheus integration

### Disadvantages

- **Complex setup** - Steep learning curve
- **Infrastructure overhead** - Need K8s cluster
- **Higher costs** - More resources required
- **Maintenance** - Ongoing cluster management

### When to Use

- Production deployments
- High availability requirements
- Multiple environments (dev/staging/prod)
- Large-scale IoT deployments
- Enterprise deployments

### Status

**Coming Soon** - Kubernetes manifests and Helm charts are planned for future releases.

## Choosing the Right Option

### For Quick Start / Demos

```bash
# Choose Docker Compose
cd mbaigo
docker-compose up --build
```

**Why:** One command, no manual setup, guaranteed to work.

### For Learning / Development

```bash
# Choose Native Go
cd mbaigo
make build
# Then follow multi-terminal setup
```

**Why:** See exactly what's happening, easier debugging, faster iteration.

### For Production

**Current recommendation:**
- Use Docker Compose for single-host deployments
- Wait for Kubernetes support for multi-host deployments
- Or create your own deployment strategy

## Hybrid Approach

You can mix and match:

### Core Systems in Docker, Application Native

```bash
# Terminal 1: Start core systems with Docker
cd /path/to/mbaigo
docker-compose up esr orchestrator

# Terminal 2: Run application natively
cd /path/to/mbaigo/examples/three-sensors
go run *.go
```

**Why:** Core systems isolated, but application easy to debug.

### Development Iteration

```bash
# Initial setup with Docker
docker-compose up -d

# Develop and test application natively
cd examples/three-sensors
go run *.go

# When done, stop Docker
docker-compose down
```

**Why:** Best of both worlds - isolated core systems, native app development.

## Migration Path

### 1. Start with Docker

```bash
docker-compose up --build
```

Learn the system, understand the architecture.

### 2. Move to Native for Development

```bash
make build
# Multi-terminal setup
```

Develop features, understand internals.

### 3. Use Docker for CI/CD

```yaml
# .github/workflows/test.yml
- name: Test with Docker
  run: docker-compose up -d && make test
```

Automated testing with consistent environment.

### 4. Plan for Production

Evaluate:
- Single host → Docker Compose
- Multi-host → Plan for Kubernetes
- Embedded → Native Go

## Common Scenarios

### Scenario: Conference Demo

**Best choice:** Docker Compose
```bash
docker-compose up --build
```
**Why:** Reliable, fast, no surprises.

### Scenario: University Course

**Best choice:** Native Go
**Why:** Students learn the architecture, understand service-oriented design.

### Scenario: IoT Production Deployment

**Current:** Native Go on each device
**Future:** Kubernetes for orchestration
**Why:** Direct hardware access, efficient resource usage.

### Scenario: CI/CD Testing

**Best choice:** Docker Compose
```yaml
test:
  script:
    - docker-compose up -d
    - make test
    - docker-compose down
```
**Why:** Isolated, reproducible, automated.

### Scenario: Contributing to mbaigo

**Best choice:** Native Go
**Why:** Easy to modify code, test changes, debug issues.

## Resource Requirements

### Docker Compose

- **CPU:** 1-2 cores (minimal usage when idle)
- **RAM:** ~200MB total
  - ESR: ~50MB
  - Orchestrator: ~50MB
  - Application: ~60MB
  - Docker overhead: ~40MB
- **Disk:** ~2GB (including Docker images)
- **Network:** Local bridge, negligible bandwidth

### Native Go

- **CPU:** 1-2 cores (minimal usage when idle)
- **RAM:** ~120MB total
  - ESR: ~40MB
  - Orchestrator: ~40MB
  - Application: ~40MB
- **Disk:** ~50MB (binaries + configs)
- **Network:** Local loopback, negligible bandwidth

### Production Kubernetes

- **CPU:** 4+ cores (cluster + workloads)
- **RAM:** 4GB+ (cluster + workloads)
- **Disk:** 20GB+ (images, volumes, logs)
- **Network:** Depends on scale

## Troubleshooting

### Docker Issues

**Problem:** Port already in use
```bash
# Solution: Change port in docker-compose.yml
ports:
  - "8102:20102"  # Different external port
```

**Problem:** Services won't start
```bash
# Solution: Check logs
docker-compose logs esr
docker-compose logs orchestrator
```

### Native Issues

**Problem:** Config conflicts
```bash
# Solution: Always run core systems from root
cd /path/to/mbaigo  # NOT examples/!
mbaigo core start esr
```

**Problem:** Port conflicts
```bash
# Solution: Kill conflicting processes
lsof -ti:20102 | xargs kill -9
```

## See Also

- [README.md](./README.md) - Main documentation
- [DOCKER.md](./DOCKER.md) - Complete Docker guide
- [TESTING.md](./TESTING.md) - Testing guide
- [examples/](./examples/) - Example applications
- [docs/ARCHITECTURE.MD](./docs/ARCHITECTURE.MD) - Architecture details
