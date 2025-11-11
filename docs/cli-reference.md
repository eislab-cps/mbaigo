# mbaigo CLI Reference

The `mbaigo` command-line tool helps you manage Arrowhead Framework systems quickly and efficiently.

## Installation

```bash
# Install from source
go install github.com/sdoque/mbaigo/cmd/mbaigo@latest

# Or build locally
make build           # Outputs to ./bin/mbaigo
./bin/mbaigo --help

# Or manually
go build -o bin/mbaigo ./cmd/mbaigo
```

## Global Flags

All commands support these flags:

- `--verbose, -v` - Enable verbose output for debugging
- `--config, -c` - Specify config file path (default: systemconfig.json)
- `--json, -j` - Output as JSON instead of formatted text
- `--no-color` - Disable colored output

## Commands

### `mbaigo init`

Initialize a new system configuration file.

```bash
# Basic initialization
mbaigo init --name MySystem --cloud FactoryCloud

# With custom ports
mbaigo init -n SensorSystem -l Building1 --http-port 9090 --https-port 9443

# Interactive mode
mbaigo init --interactive
```

**Flags:**
- `--name, -n` - System name (required)
- `--cloud, -l` - Local cloud name (required)
- `--http-port` - HTTP port (default: 8080)
- `--https-port` - HTTPS port (default: 8443)
- `--interactive, -i` - Interactive configuration mode

**Creates:**
- `systemconfig.json` with basic structure
- Empty unit_assets array ready for your assets

---

### `mbaigo config`

Manage system configuration.

#### `mbaigo config show`

Display the current configuration.

```bash
# Show as formatted table
mbaigo config show

# Show as JSON
mbaigo config show --json
```

#### `mbaigo config validate`

Validate the configuration file.

```bash
mbaigo config validate
```

**Checks:**
- Required fields present
- Valid JSON structure
- Proper field types

---

### `mbaigo system`

System management commands.

#### `mbaigo system info`

Show system information.

```bash
# Display system info
mbaigo system info

# JSON output
mbaigo system info --json
```

**Shows:**
- System name and local cloud
- Configured endpoints
- Number of unit assets
- Services per asset

#### `mbaigo system start`

Generate a template main.go file to start your system.

```bash
mbaigo system start > main.go
```

**Generates:**
- Complete working main.go
- Graceful shutdown handling
- Service registration setup
- HTTP server configuration

---

### `mbaigo service`

Service management commands.

#### `mbaigo service list`

List all configured services.

```bash
# List all services
mbaigo service list

# Filter by asset
mbaigo service list --asset TempSensor

# JSON output
mbaigo service list --json
```

**Flags:**
- `--asset, -a` - Filter by asset name

**Displays:**
- Asset name
- Service definition
- Subpath
- Registration period

#### `mbaigo service add`

Add a service to an asset.

```bash
mbaigo service add --asset MySensor --name Temperature
```

**Flags:**
- `--asset, -a` - Asset name (required)
- `--name, -n` - Service name (required)

**Creates:**
- New service definition in systemconfig.json
- Default registration period of 60 seconds
- Template details and description

---

### `mbaigo cert`

Certificate management commands.

#### `mbaigo cert generate-key`

Generate an ECDSA P256 private key.

```bash
# Generate in current directory
mbaigo cert generate-key

# Specify output directory
mbaigo cert generate-key --output ./certs
```

**Flags:**
- `--output, -o` - Output directory (default: .)

**Creates:**
- `private_key.pem` - ECDSA P256 private key (keep secure!)

#### `mbaigo cert create-csr`

Create a Certificate Signing Request.

```bash
# Create CSR using system name
mbaigo cert create-csr

# Specify common name
mbaigo cert create-csr --common-name MySystem

# Specify directory
mbaigo cert create-csr --output ./certs
```

**Flags:**
- `--output, -o` - Directory with private_key.pem (default: .)
- `--common-name, -n` - Common name (defaults to system name from config)

**Requires:**
- Existing `private_key.pem` in output directory

**Creates:**
- `csr.pem` - Certificate Signing Request to submit to CA

---

### `mbaigo generate`

Code generation commands.

#### `mbaigo generate asset`

Generate a unit asset template.

```bash
# Generate asset code
mbaigo generate asset --name TemperatureSensor

# Specify output file
mbaigo generate asset -n PressureSensor -o pressure.go
```

**Flags:**
- `--name, -n` - Asset name (required)
- `--output, -o` - Output file (default: <name>_asset.go)

**Generates:**
- Complete UnitAsset implementation
- Traits structure for configuration
- Serving() method template
- Factory function for loading from config
- Registration helper

#### `mbaigo generate example`

Generate a complete example application.

```bash
mbaigo generate example > main.go
```

**Generates:**
- Working system with example asset
- Random number service
- Complete setup and teardown
- Ready to run immediately

---

### `mbaigo version`

Show version information.

```bash
mbaigo version

# JSON output
mbaigo version --json
```

---

## Common Workflows

### Starting a New Project

```bash
# 1. Initialize configuration
mbaigo init --name MySystem --cloud LocalCloud

# 2. Generate your first asset
mbaigo generate asset --name Sensor1

# 3. Add the asset to config (edit systemconfig.json)

# 4. Generate certificates
mbaigo cert generate-key
mbaigo cert create-csr

# 5. Generate a main.go starter
mbaigo system start > main.go

# 6. Run your system
go run main.go
```

### Adding a New Service

```bash
# 1. Add service to existing asset
mbaigo service add --asset Sensor1 --name Temperature

# 2. List to verify
mbaigo service list

# 3. Edit service details in systemconfig.json

# 4. Restart system
```

### Validating Configuration

```bash
# Check configuration
mbaigo config validate

# View full config
mbaigo config show

# Check system info
mbaigo system info
```

## Output Formats

### Table Output (default)

Human-readable formatted tables with colors:

```bash
$ mbaigo service list
📡 Configured Services
================================================================================
ASSET                SERVICE              SUBPATH         REG PERIOD
--------------------------------------------------------------------------------
TempSensor           TemperatureReading   temperature     60s
PressureSensor       PressureReading      pressure        60s
--------------------------------------------------------------------------------
Total: 2 services
```

### JSON Output

Machine-readable JSON for scripting:

```bash
$ mbaigo service list --json
{
  "services": [
    {
      "asset": "TempSensor",
      "definition": "TemperatureReading",
      "subpath": "temperature",
      "registrationPeriod": 60
    }
  ]
}
```

## Environment Variables

The CLI respects standard Go environment variables:

- `GOOS`, `GOARCH` - For cross-compilation
- `GOPATH` - Go workspace location

## Exit Codes

- `0` - Success
- `1` - Error occurred

## Tips and Best Practices

1. **Use JSON output for scripting:**
   ```bash
   mbaigo service list --json | jq '.services[] | select(.asset=="Sensor1")'
   ```

2. **Validate before deploying:**
   ```bash
   mbaigo config validate && go build
   ```

3. **Keep certificates secure:**
   ```bash
   chmod 600 private_key.pem
   ```

4. **Use verbose mode for debugging:**
   ```bash
   mbaigo --verbose system info
   ```

5. **Separate configs for environments:**
   ```bash
   mbaigo --config dev-config.json system info
   mbaigo --config prod-config.json system info
   ```

## Getting Help

```bash
# General help
mbaigo --help

# Command-specific help
mbaigo init --help
mbaigo service --help
mbaigo cert create-csr --help
```

## See Also

- [Getting Started Guide](./getting-started.md)
- [Architecture Documentation](./architecture.md)
- [Examples](../examples/)
