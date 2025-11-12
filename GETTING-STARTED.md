# Getting Started with mbaigo

This guide will help you create your first mbaigo application.

## Prerequisites

- Go 1.21 or later
- Basic understanding of Go programming

## Installation

```bash
# Clone the repository
git clone https://github.com/eislab-cps/mbaigo.git
cd mbaigo

# Install dependencies
make deps

# Build the CLI
make build

# The mbaigo CLI will be available at ./bin/mbaigo
```

## Quick Start: Create Your First System

### 1. Create a New Project

Create a new directory for your project (outside the mbaigo library):

```bash
cd ~  # Or wherever you want your project
mkdir my-iot-system
cd my-iot-system
```

### 2. Initialize Your System

```bash
# Initialize with mbaigo CLI (use full path or add to PATH)
~/path/to/mbaigo/bin/mbaigo init --name MySystem --cloud LocalCloud

# This creates:
# - systemconfig.json (system configuration)
# - main.go (application entry point)
```

### 3. Initialize Go Module

```bash
go mod init my-iot-system
go mod tidy
```

### 4. Generate Assets

```bash
# Generate asset templates
~/path/to/mbaigo/bin/mbaigo generate asset --name TemperatureSensor
~/path/to/mbaigo/bin/mbaigo generate asset --name PressureSensor
~/path/to/mbaigo/bin/mbaigo generate asset --name Controller
```

This creates:
- `temperaturesensor_asset.go`
- `pressuresensor_asset.go`
- `controller_asset.go`

### 5. Configure Your Assets

Edit `systemconfig.json` to add your assets with their services and traits:

```json
{
  "systemName": "MySystem",
  "localCloud": "LocalCloud",
  "protocolsNports": {
    "http": 8080
  },
  "unit_assets": [
    {
      "name": "TempSensor1",
      "details": {
        "type": ["TemperatureSensor"]
      },
      "services": [
        {
          "subpath": "temperature",
          "definition": "temperature-reading",
          "details": {
            "Forms": ["SignalA_v1a"]
          }
        }
      ],
      "traits": [
        {
          "minValue": 0,
          "maxValue": 50
        }
      ]
    }
  ]
}
```

### 6. Run Your System

```bash
go run *.go
```

You should see:

```
🚀 System Started!
=============================================================
System Name:  MySystem
Local Cloud:  LocalCloud
HTTP:         http://localhost:8080

📍 Available Endpoints:
  GET /MySystem/TempSensor1/temperature
  GET /MySystem/PressureSensor1/pressure
  GET /MySystem/Controller1/control
=============================================================
```

### 7. Test Your Services

In another terminal:

```bash
# Test temperature sensor
curl http://localhost:8080/MySystem/TempSensor1/temperature

# Expected response:
# {"value":25,"unit":"celsius","timestamp":"2025-11-11T23:35:58+01:00","version":"SignalA_v1.0"}
```

## Example: Three Sensors

A complete working example is available in `examples/three-sensors/`:

```bash
cd examples/three-sensors
go run *.go
```

This example includes:
- Temperature sensor with configurable min/max range
- Pressure sensor with simulated readings
- Controller with boolean status
- Pre-configured systemconfig.json

## Customizing Your Assets

### Temperature Sensor Example

Edit `temperaturesensor_asset.go` to add real sensor logic:

```go
func (a *TemperatureSensor) Serving(w http.ResponseWriter, r *http.Request, servicePath string) {
    switch servicePath {
    case "temperature":
        // TODO: Replace simulation with real sensor reading
        // Example: temp := readFromDS18B20()
        temp := a.Traits.MinValue + (a.Traits.MaxValue-a.Traits.MinValue)*0.5

        signal := forms.SignalA_v1a{
            Value:     temp,
            Unit:      "celsius",
            Timestamp: time.Now(),
        }
        // ... pack and send response
    }
}
```

### Adding Custom Traits

Traits allow you to configure asset behavior:

```go
type TemperatureSensorTraits struct {
    MinValue      float64 `json:"minValue"`
    MaxValue      float64 `json:"maxValue"`
    SensorAddress string  `json:"sensorAddress"`  // Add custom fields
    SampleRate    int     `json:"sampleRate"`
}
```

## CLI Commands

### Configuration Management

```bash
# View current config
./bin/mbaigo config show

# Validate config
./bin/mbaigo config validate
```

### Asset Generation

```bash
# Generate a new asset
./bin/mbaigo generate asset --name HumiditySensor

# Generate example application
./bin/mbaigo generate example --name my-app
```

### Service Management

```bash
# Add a service to an asset
./bin/mbaigo service add --asset TempSensor1 --name status

# List all services
./bin/mbaigo service list
```

### Certificate Management (for production)

```bash
# Generate private key
./bin/mbaigo cert generate-key

# Create certificate signing request
./bin/mbaigo cert create-csr
```

### System Information

```bash
# View system details
./bin/mbaigo system info
```

## Project Structure

Your project should follow this structure:

```
my-iot-system/
├── main.go                      # Main application entry point
├── temperaturesensor_asset.go   # Asset implementations
├── pressuresensor_asset.go
├── controller_asset.go
├── systemconfig.json            # System configuration
├── go.mod                       # Go module definition
└── go.sum                       # Go dependencies
```

## Building for Production

```bash
# Build binary
go build -o my-system

# Run binary
./my-system
```

## Troubleshooting

### "Package not found"

```bash
go mod tidy
```

### "Port already in use"

Edit `systemconfig.json` and change the port:

```json
"protocolsNports": {
  "http": 9090
}
```

### "Asset factory not found"

Make sure `main.go` registers each asset type:

```go
usecases.RegisterAssetFactory("TemperatureSensor", NewTemperatureSensor)
usecases.RegisterAssetFactory("PressureSensor", NewPressureSensor)
```

### Service Registrar Warnings

These warnings are normal if you're not running Arrowhead Core Systems:

```
failed to find lead registrar: core system 'serviceregistrar' not found
```

Your HTTP services will still work locally. These warnings only affect Arrowhead Framework integration.

## Next Steps

1. **Customize Assets**: Add real sensor/actuator logic to your `*_asset.go` files
2. **Add More Assets**: Generate additional assets with `./bin/mbaigo generate asset`
3. **Configure Services**: Edit `systemconfig.json` to add services and customize behavior
4. **Deploy**: Build and deploy to your target device (Raspberry Pi, etc.)
5. **Integrate with Arrowhead**: Set up Arrowhead Core Systems for service discovery
6. **Secure**: Generate certificates for production deployments

## Documentation

- [CLI Reference](docs/CLI-REFERENCE.MD) - Complete CLI command documentation
- [Architecture](docs/ARCHITECTURE.MD) - System design and components
- [Examples](examples/) - Working example applications

## Get Help

```bash
# CLI help
./bin/mbaigo --help
./bin/mbaigo <command> --help

# Report issues
# https://github.com/eislab-cps/mbaigo/issues
```

---

**Ready to start?** Create your project directory and run `./bin/mbaigo init --name MySystem`! 🚀
