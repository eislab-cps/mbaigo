# Three Sensors Example

This example demonstrates a complete service-oriented architecture with automatic service registration, discovery, and consumption.

**What you'll build:**
- **Temperature Sensor** - Provider that registers a temperature service
- **Pressure Sensor** - Provider that registers a pressure service
- **Controller** - Consumer that discovers and calls both services

**Architecture:** 5 separate processes demonstrating the Arrowhead Framework:
1. ESR (Service Registry)
2. Orchestrator (Service Discovery)
3. Temperature Provider
4. Pressure Provider
5. Controller (Consumer)

## Quick Start

**CRITICAL SETUP:**

```bash
# Terminal 1: ESR (MUST run from mbaigo root!)
cd /path/to/mbaigo
mbaigo core start esr

# Terminal 2: Orchestrator (MUST run from mbaigo root!)
cd /path/to/mbaigo
mbaigo core start orchestrator

# Terminal 3: Application (run from examples directory)
cd /path/to/mbaigo/examples/three-sensors
go run *.go
```

**Why?** Core systems must run from root to avoid picking up application configs!

**Notes:**
- **First run**: ESR and Orchestrator create config files and exit. Run each command again to start them.
- Certificate authentication is disabled for this tutorial (no CA system needed)
- Service registration works automatically once core systems are running

## Docker Quick Start (Recommended)

Prefer Docker? Start everything with one command:

```bash
cd /path/to/mbaigo/examples/three-sensors
docker-compose up --build

# Or use the Makefile:
make docker-up
```

This automatically:
- Builds all services from source (no pre-built images)
- Starts ESR on port 20102
- Starts Orchestrator on port 20103
- Starts three-sensors application on port 8080
- Configures networking and health checks
- Handles startup order with dependencies

**Test the deployment:**

```bash
# Check ESR status
curl http://localhost:20102/serviceregistrar/registry/status

# View all registered services
curl http://localhost:20102/serviceregistrar/registry/query

# Test temperature sensor
curl http://localhost:8080/MySystem/TempSensor1/temperature

# Test pressure sensor
curl http://localhost:8080/MySystem/PressureSensor1/pressure

# Test controller
curl http://localhost:8080/MySystem/Controller1/control
```

**View logs:**
```bash
docker-compose logs -f
```

**Stop everything:**
```bash
docker-compose down
```

**Clean slate (remove volumes):**
```bash
docker-compose down -v
```

### Docker Architecture

All services run in an isolated network:
- ESR container: `three-sensors-esr`
- Orchestrator container: `three-sensors-orchestrator`
- Application container: `three-sensors-app`
- Network: `three-sensors-network`
- Volumes: `esr-data`, `orchestrator-data`, `app-data`

Services communicate internally using container names (e.g., `http://esr:20102`) and are accessible from the host via port mapping.

**See [DOCKER-README.md](./DOCKER-README.md) for complete Docker documentation including:**
- Detailed container information
- Troubleshooting guide
- Production deployment tips
- Development workflow

### Architecture Overview

```mermaid
graph TB
    subgraph cloud["Local Cloud Infrastructure"]
        ESR["ESR - Service Registry
        Port 20102"]
        ORCH["Orchestrator
        Port 20103"]
    end

    subgraph apps["Application Layer"]
        TEMP["Temperature Provider
        Port 8080"]
        PRESS["Pressure Provider
        Port 8081"]
        CTRL["Controller Consumer
        Port 8082"]
    end

    TEMP -->|"1. Register Service"| ESR
    PRESS -->|"1. Register Service"| ESR

    ORCH -->|Query Services| ESR

    CTRL -->|"2. Discovery Request"| ORCH
    ORCH -.->|"Service Location"| CTRL

    CTRL -->|"3. Direct Consumption"| TEMP
    CTRL -->|"3. Direct Consumption"| PRESS

    style ESR fill:#e1f5ff,stroke:#333,stroke-width:2px
    style ORCH fill:#e1f5ff,stroke:#333,stroke-width:2px
    style TEMP fill:#fff4e6,stroke:#333,stroke-width:2px
    style PRESS fill:#fff4e6,stroke:#333,stroke-width:2px
    style CTRL fill:#f0f0f0,stroke:#333,stroke-width:2px
```

### Prerequisites

- Go 1.24.4 or later
- Clone this repository
- Build the mbaigo CLI: `cd mbaigo && make build`

### Step 1: Start Core Systems

**CRITICAL:** Run core systems from **mbaigo root directory**, NOT from examples/three-sensors!

The mbaigo CLI includes embedded core systems. Start these FIRST before running application systems.

**Terminal 1 - Start ESR (Service Registry):**

```bash
# MUST be run from mbaigo root directory!
cd /path/to/mbaigo
mbaigo core start esr
# Or use: ./bin/mbaigo core start esr (if not installed globally)
```

**Why from root?** Running from examples directory will cause config conflicts!

- **Starts immediately** (no compilation needed!)
- First run: Creates `systemconfig.json` in `systems/esr/` and exits
- Run again: Starts on port 20102
- Tracks all registered services

**Verify it's running:**
```bash
curl http://localhost:20102/serviceregistrar/registry/status
# Should return: "lead Service Registrar since..."
```

**Terminal 2 - Start Orchestrator:**

```bash
# MUST be run from mbaigo root directory!
cd /path/to/mbaigo
mbaigo core start orchestrator
# Or use: ./bin/mbaigo core start orchestrator (if not installed globally)
```

- **Starts immediately** (embedded!)
- First run: Creates config and exits
- Run again: Starts on port 20103
- Handles service discovery

### Step 2: Start Application Systems

Now start the application systems that will register with and discover each other through the core systems.

**Terminal 3 - Start All Application Systems:**

For this example, run all three assets together in one process:

```bash
# Run from the three-sensors directory
cd /path/to/mbaigo/examples/three-sensors
go run *.go
```

This starts:
- Temperature Provider (registers temperature service)
- Pressure Provider (registers pressure service)
- Controller (discovers and consumes both services)

All services register with ESR and can be discovered through the Orchestrator.

**Summary of terminal setup:**
- **Terminal 1:** ESR (from mbaigo root)
- **Terminal 2:** Orchestrator (from mbaigo root)
- **Terminal 3:** Application (from examples/three-sensors)

### Step 3: Verify the Setup

**Check Service Registry:**

Open in browser or use curl: `http://localhost:20102/serviceregistrar/registry/query`

You should see all registered services with:
- Service ID, definition, system name
- IP addresses and ports
- Expiration times (TTL)

**Test Services:**

```bash
# Temperature
curl http://localhost:8080/MySystem/TempSensor1/temperature

# Pressure
curl http://localhost:8080/MySystem/PressureSensor1/pressure

# Controller
curl http://localhost:8080/MySystem/Controller1/control
```

All running from port 8080 since they're in the same process.

### Service Flow

```mermaid
sequenceDiagram
    participant Temp as Temperature Provider
    participant Press as Pressure Provider
    participant ESR as ESR Registry
    participant Orch as Orchestrator
    participant Ctrl as Controller Consumer

    Note over Temp,Press: 1. Service Registration
    Temp->>ESR: POST /register (temperature)
    ESR-->>Temp: ID: 1, expires in 30s
    Press->>ESR: POST /register (pressure)
    ESR-->>Press: ID: 2, expires in 30s

    Note over ESR: Registry contains temperature & pressure

    Note over Ctrl,Orch: 2. Service Discovery
    Ctrl->>Orch: POST /squest (find temperature)
    Orch->>ESR: POST /query (search: temperature)
    ESR-->>Orch: ServiceRecord ID:1
    Orch-->>Ctrl: http://localhost:8080/...

    Ctrl->>Orch: POST /squest (find pressure)
    Orch->>ESR: POST /query (search: pressure)
    ESR-->>Orch: ServiceRecord ID:2
    Orch-->>Ctrl: http://localhost:8081/...

    Note over Ctrl: Cache service URLs

    Note over Temp,Ctrl: 3. Service Consumption
    Ctrl->>Temp: GET /temperature
    Temp-->>Ctrl: value:25 celsius
    Ctrl->>Press: GET /pressure
    Press-->>Ctrl: value:1.025 bar

    Note over Ctrl: Process sensor data

    Note over Temp,Press: Re-registration every ~25s
    Temp->>ESR: PUT /register ID:1
    Press->>ESR: PUT /register ID:2
```

**Flow Details:**

1. **Registration:**
   - Providers register services with ESR
   - Re-registration happens every ~25 seconds

2. **Discovery:**
   - Controller asks Orchestrator for services
   - Orchestrator queries ESR and returns provider URLs
   - Controller caches service locations

3. **Consumption:**
   - Controller makes HTTP GET requests to providers
   - Providers return sensor data
   - Controller processes the data

### Stopping All Processes

Press `Ctrl+C` in each terminal (reverse order recommended):
1. Application systems (Terminal 3)
2. Orchestrator (Terminal 2)
3. ESR (Terminal 1)

Systems gracefully shutdown:
- Unregister services from ESR
- Close HTTP servers
- Clean up resources

**Note:** You can also run the three application assets as separate processes by creating individual main files for each. This example runs them together for simplicity.

### Common Issues

**Port Already in Use:**
```bash
# Find and kill process using port 8080
lsof -ti:8080 | xargs kill -9
```

**ESR Not Found:**
- Ensure ESR is running on port 20102
- Check `systemconfig.json` has correct ESR URL:
  ```json
  {
    "coreSystems": [{
      "coreSystem": "serviceregistrar",
      "url": "http://localhost:20102/serviceregistrar/registry"
    }]
  }
  ```

**Service Not Found:**
- Verify provider is running and registered
- Check ESR: `curl http://localhost:20102/serviceregistrar/registry/query`
- Ensure service names match exactly

---

## Files

- **main.go** - Application entry point
  - Creates the system
  - Registers asset factories
  - Loads configuration
  - Starts HTTP server

- **temperaturesensor_asset.go** - Temperature sensor
  - Configurable min/max range via Traits
  - Simulated readings with variation
  - Returns SignalA_v1a form

- **pressuresensor_asset.go** - Pressure sensor
  - Configurable min/max range via Traits
  - Simulated readings with variation
  - Returns SignalA_v1a form

- **controller_asset.go** - Controller
  - Returns boolean status (active/inactive)
  - Returns SignalB_v1a form

- **systemconfig.json** - System configuration (auto-generated)
  - Defines three assets with their services
  - Configures traits (min/max values)
  - Sets HTTP port and core system URLs

## Key Concepts

### 1. Asset Factory Pattern

```go
// Register factories for dynamic asset creation
usecases.RegisterAssetFactory("TemperatureSensor", NewTemperatureSensor)
usecases.RegisterAssetFactory("PressureSensor", NewPressureSensor)
usecases.RegisterAssetFactory("Controller", NewController)
```

Assets are instantiated at runtime based on configuration.

### 2. Configuration-Driven Assets

The `systemconfig.json` file defines assets:

```json
{
  "name": "TempSensor1",
  "details": {
    "type": ["TemperatureSensor"]
  },
  "traits": [
    {
      "minValue": 0,
      "maxValue": 50
    }
  ]
}
```

### 3. Service Implementation

Each asset implements the `Serving()` method:

```go
func (a *TemperatureSensor) Serving(w http.ResponseWriter, r *http.Request, servicePath string) {
    switch servicePath {
    case "temperature":
        signal := forms.SignalA_v1a{
            Value:     temp,
            Unit:      "celsius",
            Timestamp: time.Now(),
        }
        packed, _ := usecases.Pack(signal.NewForm(), "application/json")
        w.Header().Set("Content-Type", "application/json")
        w.Write(packed)
    }
}
```

### 4. Traits for Asset Configuration

Traits allow runtime configuration:

```go
type TemperatureSensorTraits struct {
    MinValue float64 `json:"minValue"`
    MaxValue float64 `json:"maxValue"`
}
```

## Customization Ideas

### Replace Simulation with Real Sensors

**DS18B20 Temperature Sensor:**

```go
import "github.com/yryz/ds18b20"

func (a *TemperatureSensor) Serving(w http.ResponseWriter, r *http.Request, servicePath string) {
    sensors, _ := ds18b20.Sensors()
    temp, _ := ds18b20.Temperature(sensors[0])

    signal := forms.SignalA_v1a{
        Value:     temp,
        Unit:      "celsius",
        Timestamp: time.Now(),
    }
    // ...
}
```

### Add State Management

```go
type Controller struct {
    state      string
    lastUpdate time.Time
}

func (a *Controller) Serving(w http.ResponseWriter, r *http.Request, servicePath string) {
    if r.Method == "POST" {
        var cmd ControlCommand
        json.NewDecoder(r.Body).Decode(&cmd)
        a.executeCommand(cmd)
    }
    // Return current state
}
```

### Add Error Handling

```go
func (a *PressureSensor) Serving(w http.ResponseWriter, r *http.Request, servicePath string) {
    pressure, err := a.readSensor()
    if err != nil {
        http.Error(w, "sensor read failed", http.StatusInternalServerError)
        return
    }
    // ... continue with valid reading
}
```

## Architecture Notes

This example follows the mbaigo architecture:

```mermaid
graph TD
    System["System
    MySystem"]
    Husk["Husk
    Metadata & Config"]

    subgraph assets["UnitAssets"]
        TempAsset["TemperatureSensor
        TempSensor1"]
        PressAsset["PressureSensor
        PressureSensor1"]
        CtrlAsset["Controller
        Controller1"]
    end

    subgraph services["Services - HTTP Endpoints"]
        TempService["temperature
        GET"]
        PressService["pressure
        GET"]
        CtrlService["control
        GET"]
    end

    subgraph forms["Forms - Data Exchange"]
        SignalA["SignalA_v1a
        Value + Unit"]
        SignalB["SignalB_v1a
        Boolean Status"]
    end

    System --> Husk
    System --> TempAsset
    System --> PressAsset
    System --> CtrlAsset

    TempAsset --> TempService
    PressAsset --> PressService
    CtrlAsset --> CtrlService

    TempService -.->|returns| SignalA
    PressService -.->|returns| SignalA
    CtrlService -.->|returns| SignalB

    style System fill:#e8f4f8,stroke:#333,stroke-width:3px
    style Husk fill:#fff9e6,stroke:#333,stroke-width:2px
    style TempAsset fill:#e6f3ff,stroke:#333,stroke-width:2px
    style PressAsset fill:#e6f3ff,stroke:#333,stroke-width:2px
    style CtrlAsset fill:#e6f3ff,stroke:#333,stroke-width:2px
    style TempService fill:#f0f0f0,stroke:#333,stroke-width:1px
    style PressService fill:#f0f0f0,stroke:#333,stroke-width:1px
    style CtrlService fill:#f0f0f0,stroke:#333,stroke-width:1px
    style SignalA fill:#e8ffe8,stroke:#333,stroke-width:1px
    style SignalB fill:#e8ffe8,stroke:#333,stroke-width:1px
```

**Component Hierarchy:**

1. **System** - Container for all components
2. **Husk** - System metadata (ports, description, etc.)
3. **UnitAssets** - Individual devices/sensors/actuators
4. **Services** - HTTP endpoints exposed by assets
5. **Forms** - Standardized data exchange formats (SignalA_v1a, SignalB_v1a)

## Makefile Commands

Convenience commands for easy operation:

```bash
# Docker operations (recommended)
make docker-up       # Start all services with Docker
make docker-down     # Stop all services
make docker-logs     # View logs from all containers
make docker-restart  # Restart services
make docker-clean    # Stop and remove volumes

# Native Go (requires core systems running separately)
make run            # Run application with go run
make test           # Test all service endpoints

# Help
make help           # Show all available commands
```

## CLI Reference

```bash
# List core systems (embedded ones marked [EMBEDDED])
mbaigo core list

# Start embedded core systems
mbaigo core start esr
mbaigo core start orchestrator
```

**Install mbaigo globally:**

```bash
cd /path/to/mbaigo
make build
sudo make install
# Now you can use 'mbaigo' from anywhere!
```

**Benefits of embedded core systems:**
- Single binary deployment (10MB)
- No runtime compilation
- Instant startup
- No separate go.mod management
- Consistent versions

## Next Steps

1. **Modify simulations** - Change random value generation
2. **Add new services** - Extend assets with additional endpoints
3. **Create new assets** - Additional sensor types (humidity, light, etc.)
4. **Integrate hardware** - Replace simulations with real I/O (GPIO, I2C, etc.)
5. **Add persistence** - Store readings in database
6. **Use HTTPS** - Enable mutual TLS authentication
7. **Separate processes** - Create dedicated applications for each provider:
   ```
   examples/temperature-provider/    # Own main.go, config, port
   examples/pressure-provider/       # Own main.go, config, port
   examples/controller-consumer/     # Own main.go, config, port
   ```

## Learn More

- [Main README](../../README.md) - Full mbaigo documentation
- [Getting Started Guide](../../GETTING-STARTED.md) - Tutorial guide
- [CLI Reference](../../docs/CLI-REFERENCE.MD) - Complete CLI docs
- [Architecture](../../docs/ARCHITECTURE.MD) - System architecture details
- [Examples](../README.md) - More example applications
