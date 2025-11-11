# Three Sensors Example

This example demonstrates a complete mbaigo system with three different asset types:

- **Temperature Sensor** - Simulated analog sensor reading temperature in Celsius
- **Pressure Sensor** - Simulated analog sensor reading pressure in bar
- **Controller** - Digital controller providing boolean status

## Running the Example

```bash
cd examples/three-sensors
go run *.go
```

## Expected Output

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

💡 Test with:
  curl http://localhost:8080/MySystem/TempSensor1/temperature
  curl http://localhost:8080/MySystem/PressureSensor1/pressure
  curl http://localhost:8080/MySystem/Controller1/control

Press Ctrl+C to stop...
```

## Testing the Services

In another terminal:

```bash
# Test temperature sensor
curl http://localhost:8080/MySystem/TempSensor1/temperature
# Response: {"value":25,"unit":"celsius","timestamp":"2025-11-11T23:35:58+01:00","version":"SignalA_v1.0"}

# Test pressure sensor
curl http://localhost:8080/MySystem/PressureSensor1/pressure
# Response: {"value":1.025,"unit":"bar","timestamp":"2025-11-11T23:36:02+01:00","version":"SignalA_v1.0"}

# Test controller
curl http://localhost:8080/MySystem/Controller1/control
# Response: {"value":true,"timestamp":"2025-11-11T23:36:02+01:00","version":"SignalB_v1.0"}
```

## Files

- **main.go** - Application entry point that:
  - Creates the system
  - Registers asset factories
  - Loads configuration
  - Starts HTTP server

- **temperaturesensor_asset.go** - Temperature sensor implementation
  - Configurable min/max range via Traits
  - Simulated readings with variation
  - Returns SignalA_v1a form

- **pressuresensor_asset.go** - Pressure sensor implementation
  - Configurable min/max range via Traits
  - Simulated readings with variation
  - Returns SignalA_v1a form

- **controller_asset.go** - Controller implementation
  - Returns boolean status (active/inactive)
  - Returns SignalB_v1a form

- **systemconfig.json** - System configuration
  - Defines three assets with their services
  - Configures traits (min/max values for sensors)
  - Sets HTTP port to 8080

## Key Concepts Demonstrated

### 1. Asset Factory Pattern

```go
// Register factories for dynamic asset creation
usecases.RegisterAssetFactory("TemperatureSensor", NewTemperatureSensor)
usecases.RegisterAssetFactory("PressureSensor", NewPressureSensor)
usecases.RegisterAssetFactory("Controller", NewController)
```

### 2. Configuration-Driven Assets

The `systemconfig.json` file defines assets that are instantiated at runtime:

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

Each asset implements the `Serving()` method to handle HTTP requests:

```go
func (a *TemperatureSensor) Serving(w http.ResponseWriter, r *http.Request, servicePath string) {
    switch servicePath {
    case "temperature":
        // Create signal with sensor data
        signal := forms.SignalA_v1a{
            Value:     temp,
            Unit:      "celsius",
            Timestamp: time.Now(),
        }
        // Pack and send response
        packed, _ := usecases.Pack(signal.NewForm(), "application/json")
        w.Header().Set("Content-Type", "application/json")
        w.Write(packed)
    }
}
```

### 4. Traits for Asset Configuration

Traits allow runtime configuration of asset behavior:

```go
type TemperatureSensorTraits struct {
    MinValue float64 `json:"minValue"`
    MaxValue float64 `json:"maxValue"`
}
```

## Customization Ideas

### Replace Simulation with Real Sensors

**Temperature Sensor (DS18B20 example):**
```go
import "github.com/yryz/ds18b20"

func (a *TemperatureSensor) Serving(w http.ResponseWriter, r *http.Request, servicePath string) {
    // Read from real sensor
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

**Controller with command handling:**
```go
type Controller struct {
    // ... existing fields
    state      string
    lastUpdate time.Time
}

func (a *Controller) Serving(w http.ResponseWriter, r *http.Request, servicePath string) {
    switch servicePath {
    case "control":
        if r.Method == "POST" {
            // Handle control commands
            var cmd ControlCommand
            json.NewDecoder(r.Body).Decode(&cmd)
            a.executeCommand(cmd)
        }
        // Return current state
    }
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

1. **System** - Container for all components
2. **Husk** - System metadata (ports, description, etc.)
3. **UnitAssets** - Individual devices/sensors/actuators
4. **Services** - HTTP endpoints exposed by assets
5. **Forms** - Standardized data exchange formats (SignalA_v1a, SignalB_v1a)

## Next Steps

1. **Modify the simulation** - Change the random value generation logic
2. **Add new services** - Extend assets with additional endpoints
3. **Add new assets** - Create additional sensor types
4. **Integrate real hardware** - Replace simulations with actual I/O
5. **Add persistence** - Store readings in a database
6. **Add MQTT** - Publish readings to an MQTT broker
7. **Connect to Arrowhead** - Integrate with Arrowhead Framework core systems

## Learn More

- [Getting Started Guide](../../GETTING-STARTED.md)
- [CLI Reference](../../docs/cli-reference.md)
- [mbaigo Documentation](../../README.md)
