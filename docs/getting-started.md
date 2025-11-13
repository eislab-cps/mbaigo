# Getting Started with mbaigo

This guide will walk you through creating your first Arrowhead system using mbaigo in less than 30 minutes.

## Prerequisites

- Go 1.24.4 or later
- Basic understanding of Go programming
- (Optional) Access to Arrowhead Core Systems for full functionality

## Installation

### Option 1: Install from source

```bash
go install github.com/sdoque/mbaigo/cmd/mbaigo@latest
```

### Option 2: Build locally

```bash
git clone https://github.com/sdoque/mbaigo
cd mbaigo
make build  # Builds to ./bin/mbaigo

# Optional: install globally
sudo cp bin/mbaigo /usr/local/bin/

# Or add to PATH
export PATH=$PATH:$(pwd)/bin
```

### Verify Installation

```bash
mbaigo version
# Output: mbaigo version 0.1.0
```

## Quick Start: Your First System in 5 Minutes

### Step 1: Initialize a New System

```bash
# Create a new directory for your project
mkdir my-arrowhead-system
cd my-arrowhead-system

# Initialize the system configuration
mbaigo init --name TemperatureMonitor --cloud SmartBuilding

# Created configuration file: systemconfig.json
```

This creates a `systemconfig.json` file with basic structure.

### Step 2: Generate an Example Application

```bash
# Generate a working example
mbaigo generate example > main.go
```

### Step 3: Run Your System

```bash
# Initialize Go module
go mod init my-temperature-system
go mod tidy

# Run the system
go run main.go
```

You should see:

```
System started: TemperatureMonitor
Try: curl http://localhost:8080/TemperatureMonitor/example/random
```

### Step 4: Test Your Service

In another terminal:

```bash
curl http://localhost:8080/TemperatureMonitor/example/random
```

Expected output:

```json
{
  "version": "1a",
  "typename": "SignalA",
  "value": 42.73891234,
  "unit": "random"
}
```

**Congratulations!** You've just created and run your first Arrowhead system!

---

## Building a Real System

Now let's build something more practical: a temperature sensor system.

### Step 1: Create Your Project

```bash
mkdir temperature-system
cd temperature-system

# Initialize
mbaigo init --name TempMonitor --cloud Building1 --http-port 9090

# Initialize Go module
go mod init temperature-system
```

### Step 2: Generate Your Custom Asset

```bash
# Generate asset template
mbaigo generate asset --name TemperatureSensor

# This creates: temperaturesensor_asset.go
```

### Step 3: Customize the Asset

Edit `temperaturesensor_asset.go`:

```go
package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"time"

	"github.com/sdoque/mbaigo/pkg/components"
	"github.com/sdoque/mbaigo/pkg/forms"
	"github.com/sdoque/mbaigo/pkg/usecases"
)

// TemperatureSensorTraits defines configurable parameters
type TemperatureSensorTraits struct {
	MinValue float64 `json:"minValue"`
	MaxValue float64 `json:"maxValue"`
	Unit     string  `json:"unit"`
	Location string  `json:"location"`
}

// TemperatureSensor provides temperature readings
type TemperatureSensor struct {
	Name        string
	Details     map[string][]string
	ServicesMap components.Services
	CervicesMap components.Cervices
	Traits      TemperatureSensorTraits

	// Simulated current temperature
	currentTemp float64
}

func (a *TemperatureSensor) GetName() string                  { return a.Name }
func (a *TemperatureSensor) GetServices() components.Services { return a.ServicesMap }
func (a *TemperatureSensor) GetCervices() components.Cervices { return a.CervicesMap }
func (a *TemperatureSensor) GetDetails() map[string][]string  { return a.Details }
func (a *TemperatureSensor) GetTraits() any                   { return a.Traits }

func (a *TemperatureSensor) Serving(w http.ResponseWriter, r *http.Request, servicePath string) {
	switch servicePath {
	case "temperature":
		// Simulate temperature reading (you'd read from real sensor here)
		a.currentTemp = a.Traits.MinValue + rand.Float64()*(a.Traits.MaxValue-a.Traits.MinValue)

		signal := forms.SignalA_v1a{
			Value:     a.currentTemp,
			Unit:      a.Traits.Unit,
			Timestamp: time.Now().Unix(),
		}

		packed, err := usecases.Pack(signal.NewForm(), "application/json")
		if err != nil {
			http.Error(w, "failed to pack response", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(packed)

	default:
		http.Error(w, "unknown service path", http.StatusNotFound)
	}
}

func NewTemperatureSensor(config usecases.ConfigurableAsset) components.UnitAsset {
	var traits TemperatureSensorTraits
	if len(config.Traits) > 0 {
		json.Unmarshal(config.Traits[0], &traits)
	}

	return &TemperatureSensor{
		Name:        config.Name,
		Details:     config.Details,
		ServicesMap: usecases.MakeServiceMap(config.Services),
		Traits:      traits,
		currentTemp: 22.0, // Initial value
	}
}

func init() {
	usecases.RegisterAssetFactory("TemperatureSensor", NewTemperatureSensor)
}
```

### Step 4: Update Configuration

Edit `systemconfig.json` to add your asset:

```json
{
  "systemname": "TempMonitor",
  "localcloud": "Building1",
  "protocolsNports": {
    "http": 9090,
    "https": 9443
  },
  "coreSystems": [
    {
      "coreSystem": "serviceregistrar",
      "url": "https://localhost:8000"
    },
    {
      "coreSystem": "orchestrator",
      "url": "https://localhost:8001"
    }
  ],
  "unit_assets": [
    {
      "name": "TempSensor1",
      "details": {
        "location": ["room101"],
        "type": ["PT100"]
      },
      "services": [
        {
          "definition": "TemperatureReading",
          "subpath": "temperature",
          "registrationPeriod": 60,
          "details": {
            "datatype": ["float64"],
            "unit": ["celsius"]
          },
          "description": "provides current temperature reading",
          "costUnit": "credits"
        }
      ],
      "traits": [
        {
          "minValue": 18.0,
          "maxValue": 28.0,
          "unit": "celsius",
          "location": "room101"
        }
      ]
    }
  ]
}
```

Verify configuration:

```bash
mbaigo config validate
# Configuration is valid

mbaigo system info
# Shows your configured system
```

### Step 5: Create Main Application

Create `main.go`:

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/sdoque/mbaigo/pkg/components"
	"github.com/sdoque/mbaigo/pkg/usecases"
)

func main() {
	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create system
	sys := components.NewSystem("TempMonitor", ctx)

	// Load configuration
	rawAssets, err := usecases.Configure(&sys)
	if err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	// Load unit assets
	for _, rawAsset := range rawAssets {
		factory, ok := usecases.GetAssetFactory(rawAsset.Details["type"][0])
		if !ok {
			log.Printf("Warning: Unknown asset type: %s", rawAsset.Details["type"][0])
			continue
		}

		asset := factory(rawAsset)
		sys.UAssets[asset.GetName()] = &asset
	}

	// Setup servers and services
	go usecases.SetoutServers(&sys)
	go usecases.RegisterServices(&sys)

	fmt.Printf("Temperature Monitor System Started\n")
	fmt.Printf("System: %s\n", sys.Name)
	fmt.Printf("HTTP:   http://localhost:%d\n", sys.Husk.ProtoPort["http"])
	fmt.Printf("\nEndpoints:\n")
	for name, asset := range sys.UAssets {
		for subpath := range (*asset).GetServices() {
			fmt.Printf("  GET /%s/%s/%s\n", sys.Name, name, subpath)
		}
	}

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	fmt.Println("\n\nShutting down gracefully...")
	cancel()
}
```

### Step 6: Run Your System

```bash
# Download dependencies
go mod tidy

# Run the system
go run *.go
```

Output:

```
Temperature Monitor System Started
System: TempMonitor
HTTP:   http://localhost:9090

Endpoints:
  GET /TempMonitor/TempSensor1/temperature
```

### Step 7: Test Your Service

```bash
# Get temperature reading
curl http://localhost:9090/TempMonitor/TempSensor1/temperature

# Pretty print with jq
curl -s http://localhost:9090/TempMonitor/TempSensor1/temperature | jq
```

Output:

```json
{
  "version": "1a",
  "typename": "SignalA",
  "value": 23.451234,
  "unit": "celsius",
  "timestamp": 1699876543
}
```

---

## Adding Certificates (Production Setup)

For production systems, you need certificates for secure communication.

### Generate Private Key

```bash
mbaigo cert generate-key
# Generated private key: ./private_key.pem
# WARNING: Keep this file secure and do not share it!
```

### Create Certificate Signing Request

```bash
mbaigo cert create-csr
# Created CSR: ./csr.pem
# Common Name: TempMonitor
#
# Submit this CSR to your Certificate Authority to obtain a signed certificate
```

### Configure System with Certificates

Update your system to use the certificates (see examples/consumer-provider for full implementation).

---

## Next Steps

### Learn More

1. **Explore Examples:**
   ```bash
   cd examples/simple
   cat README.md
   go run main.go
   ```

2. **Read the Architecture:**
   - [Architecture Documentation](./ARCHITECTURE.MD)
   - [CLI Reference](./CLI-REFERENCE.MD)

3. **Add More Services:**
   ```bash
   mbaigo service add --asset TempSensor1 --name Humidity
   mbaigo service list
   ```

4. **Create More Assets:**
   ```bash
   mbaigo generate asset --name PressureSensor
   # Edit and integrate into your system
   ```

### Common Tasks

#### Viewing Configuration

```bash
# Show current config
mbaigo config show

# Validate
mbaigo config validate

# System info
mbaigo system info
```

#### Managing Services

```bash
# List all services
mbaigo service list

# Filter by asset
mbaigo service list --asset TempSensor1

# Add new service
mbaigo service add --asset TempSensor1 --name Status
```

#### Code Generation

```bash
# Generate new asset
mbaigo generate asset --name HumiditySensor

# Generate complete example
mbaigo generate example > example.go
```

---

## Troubleshooting

### "Config file not found"

```bash
# Make sure you're in the right directory
ls systemconfig.json

# Or specify the config file
mbaigo --config /path/to/systemconfig.json system info
```

### "Port already in use"

Change the port in `systemconfig.json`:

```json
{
  "protocolsNports": {
    "http": 9090,  // Changed from 8080
    "https": 9443
  }
}
```

### "Asset factory not found"

Make sure your asset file has the `init()` function:

```go
func init() {
	usecases.RegisterAssetFactory("YourAssetType", NewYourAsset)
}
```

---

## Getting Help

- **CLI Help:** `mbaigo --help`
- **Command Help:** `mbaigo <command> --help`
- **Documentation:** See [docs/](../docs/) directory
- **Examples:** See [examples/](../examples/) directory
- **Issues:** https://github.com/sdoque/mbaigo/issues

---

## Summary

You've learned how to:

- Install the mbaigo CLI
- Initialize a new system configuration
- Generate and customize unit assets
- Create a complete working system
- Test your services
- Manage certificates for production

**Ready to build more?** Check out the [examples](../examples/) directory for more advanced use cases!
