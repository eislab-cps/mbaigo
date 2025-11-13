# Testing Guide

This document provides step-by-step testing instructions to verify mbaigo works correctly.

## Prerequisites

- Go 1.21 or later installed
- Terminal with bash or zsh
- Port 8080, 20102, and 20103 available

## Complete Test Workflow

### 1. Build the Project

```bash
cd /path/to/mbaigo
make build
```

**Expected output:**
- `bin/mbaigo` binary created
- No build errors

### 2. First Run - Create Config Files

Both core systems need to run once to generate their configuration files.

**Terminal 1: ESR (Service Registry)**
```bash
cd /path/to/mbaigo  # MUST be in root directory
./bin/mbaigo core start esr
```

**Expected output:**
```
Starting esr (embedded)...
============================================================
```

The system will create `systems/esr/systemconfig.json` and exit with an error message about new config.

Run the command **again** to actually start ESR:
```bash
./bin/mbaigo core start esr
```

**Expected output:**
```
Starting esr (embedded)...
============================================================
2025/11/13 XX:XX:XX The system serviceregistrar is up with its web server available at http://XXX.XXX.XXX.XXX:20102/serviceregistrar
2025/11/13 XX:XX:XX Taking the service registry lead at ...
```

**Verify ESR is running:**
```bash
# In another terminal
curl http://localhost:20102/serviceregistrar/registry/status
```

**Expected output:**
```
lead Service Registrar since 2025-11-13 XX:XX:XX...
```

**Terminal 2: Orchestrator**
```bash
cd /path/to/mbaigo  # MUST be in root directory
./bin/mbaigo core start orchestrator
```

First run creates config and exits. Run **again** to start:
```bash
./bin/mbaigo core start orchestrator
```

**Expected output:**
```
Starting orchestrator (embedded)...
============================================================
2025/11/13 XX:XX:XX The system orchestrator is up with its web server available at http://XXX.XXX.XXX.XXX:20103/orchestrator
```

### 3. Run Example Application

**Terminal 3: Three Sensors Example**
```bash
cd /path/to/mbaigo/examples/three-sensors
go run *.go
```

**Expected output:**
```
2025/11/13 XX:XX:XX The system MySystem is up with its web server available at http://XXX.XXX.XXX.XXX:8080/MySystem
```

You may see some initial registration errors like:
```
registering service: registration request: Post "http://localhost:20102/serviceregistrar/registry/register": ...
```

These are normal during startup and will resolve once ESR is ready.

### 4. Verify Services

**Check ESR has registered services:**
```bash
curl http://localhost:20102/serviceregistrar/registry/query
```

**Expected output:** HTML page listing:
- Service ID: 1 with definition **temperature**
- Service ID: 2 with definition **pressure**
- Service ID: 3 with definition **control**

**Test Temperature Service:**
```bash
curl http://localhost:8080/MySystem/TempSensor1/temperature
```

**Expected output:**
```json
{
  "value": 25.123,
  "unit": "celsius",
  "timestamp": "2025-11-13T09:30:00Z"
}
```

**Test Pressure Service:**
```bash
curl http://localhost:8080/MySystem/PressureSensor1/pressure
```

**Expected output:**
```json
{
  "value": 1.025,
  "unit": "bar",
  "timestamp": "2025-11-13T09:30:00Z"
}
```

**Test Controller Service:**
```bash
curl http://localhost:8080/MySystem/Controller1/control
```

**Expected output:**
```json
{
  "status": true,
  "timestamp": "2025-11-13T09:30:00Z"
}
```

### 5. Cleanup

Stop all processes with `Ctrl+C` in each terminal:
1. Terminal 3 (application)
2. Terminal 2 (orchestrator)
3. Terminal 1 (ESR)

**Expected:** Systems unregister services and shutdown gracefully.

## Common Issues

### Port Already in Use

**Error:**
```
listen tcp :20102: bind: address already in use
```

**Solution:**
```bash
# Find and kill process using the port
lsof -ti:20102 | xargs kill -9
```

### Config File Conflicts

**Error:**
```
The system MySystem is up with its web server available at http://XXX:20102/serviceregistrar
```

ESR is identifying as "MySystem" instead of "serviceregistrar".

**Cause:** Running core systems from wrong directory (e.g., from examples/)

**Solution:** Always run core systems from mbaigo root:
```bash
cd /path/to/mbaigo  # NOT examples/three-sensors!
./bin/mbaigo core start esr
```

### Certificate Errors (Should Not Occur)

If you see:
```
certification failure: failed to send CSR: Post "http://localhost:20100/ca/certification/certify": connection refused
```

This means certificate authentication is still enabled. It should be commented out in:
- `pkg/core/esr/esr.go` (line ~92)
- `pkg/core/orchestrator/orchestrator.go` (line ~86)

## Automated Test Script

Run the automated config creation test:

```bash
cd /path/to/mbaigo
chmod +x test-quickstart.sh
./test-quickstart.sh
```

This verifies:
- Build succeeds
- ESR config is created
- Orchestrator config is created

You still need to manually start the systems in separate terminals to test the full workflow.

## Success Criteria

- ESR starts on port 20102
- Orchestrator starts on port 20103
- Application starts on port 8080
- All three services register with ESR
- Temperature, pressure, and control endpoints return valid JSON
- No certificate errors
- Systems shutdown gracefully

## Next Steps

After verifying the basic workflow:
1. Review [examples/three-sensors/README.md](./examples/three-sensors/README.md) for architecture details
2. Explore [docs/ARCHITECTURE.MD](./docs/ARCHITECTURE.MD) for system design
3. Try creating your own system with `mbaigo init`
