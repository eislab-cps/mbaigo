#!/bin/bash
# Test script for mbaigo quick start guide
# This script verifies the complete workflow works

set -e

echo "=========================================="
echo "Testing mbaigo Quick Start Guide"
echo "=========================================="
echo ""

# Check if we're in the right directory
if [ ! -f "go.mod" ] || [ ! -d "cmd/mbaigo" ]; then
    echo "ERROR: Must run from mbaigo root directory"
    exit 1
fi

echo "Running from mbaigo root directory"
echo ""

# Step 1: Build
echo "Step 1: Building mbaigo..."
make build
if [ ! -f "bin/mbaigo" ]; then
    echo "ERROR: bin/mbaigo not found after build"
    exit 1
fi
echo "Build successful"
echo ""

# Step 2: Start ESR (first run to create config)
echo "Step 2: Creating ESR config..."
timeout 3 ./bin/mbaigo core start esr || true
if [ ! -f "systems/esr/systemconfig.json" ]; then
    echo "ERROR: ESR config not created"
    exit 1
fi
echo "ESR config created at systems/esr/systemconfig.json"
echo ""

# Step 3: Start Orchestrator (first run to create config)
echo "Step 3: Creating Orchestrator config..."
timeout 3 ./bin/mbaigo core start orchestrator || true
if [ ! -f "systems/orchestrator/systemconfig.json" ]; then
    echo "ERROR: Orchestrator config not created"
    exit 1
fi
echo "Orchestrator config created at systems/orchestrator/systemconfig.json"
echo ""

echo "=========================================="
echo "All config files created successfully!"
echo "=========================================="
echo ""
echo "Next steps (manual):"
echo "1. Terminal 1: ./bin/mbaigo core start esr"
echo "2. Terminal 2: ./bin/mbaigo core start orchestrator"
echo "3. Terminal 3: cd examples/three-sensors && go run *.go"
echo "4. Test: curl http://localhost:8080/MySystem/TempSensor1/temperature"
echo ""
