# Makefile for three-sensors example

.PHONY: help run docker-up docker-down docker-logs docker-clean docker-restart test

# Default target
help:
	@echo "Three Sensors Example - Makefile targets:"
	@echo ""
	@echo "Native Go:"
	@echo "  make run           - Run with go run (requires core systems running)"
	@echo "  make test          - Test all endpoints"
	@echo ""
	@echo "Docker:"
	@echo "  make docker-up     - Start all services with Docker (recommended)"
	@echo "  make docker-down   - Stop all services"
	@echo "  make docker-logs   - View logs from all services"
	@echo "  make docker-restart - Restart all services"
	@echo "  make docker-clean  - Stop and remove all volumes"
	@echo ""

# Run natively (assumes core systems are already running)
run:
	@echo "Starting three-sensors application..."
	@echo "Make sure ESR and Orchestrator are running first!"
	@go run *.go

# Start everything with Docker
docker-up:
	@echo "Starting three-sensors example with Docker Compose..."
	@docker-compose up -d --build
	@echo ""
	@echo "Services started!"
	@echo ""
	@echo "Services running:"
	@echo "  ESR:          http://localhost:20102/serviceregistrar/registry/status"
	@echo "  Orchestrator: http://localhost:20103/orchestrator/orchestration/status"
	@echo "  Temperature:  http://localhost:8080/MySystem/TempSensor1/temperature"
	@echo "  Pressure:     http://localhost:8080/MySystem/PressureSensor1/pressure"
	@echo "  Controller:   http://localhost:8080/MySystem/Controller1/control"
	@echo ""
	@echo "View logs:  make docker-logs"
	@echo "Stop:       make docker-down"

# Stop Docker services
docker-down:
	@echo "Stopping Docker services..."
	@docker-compose down
	@echo "Services stopped"

# View Docker logs
docker-logs:
	@docker-compose logs -f

# Restart Docker services
docker-restart:
	@echo "Restarting services..."
	@docker-compose restart
	@echo "Services restarted"

# Clean everything (remove volumes)
docker-clean:
	@echo "Stopping services and removing volumes..."
	@docker-compose down -v
	@echo "Services stopped and volumes removed"

# Test all endpoints
test:
	@echo "Testing three-sensors endpoints..."
	@echo ""
	@echo "1. Checking ESR status..."
	@curl -s http://localhost:20102/serviceregistrar/registry/status && echo "" || echo "ERROR: ESR not responding"
	@echo ""
	@echo "2. Testing temperature sensor..."
	@curl -s http://localhost:8080/MySystem/TempSensor1/temperature | head -c 100 && echo "..." || echo "ERROR: Temperature service not responding"
	@echo ""
	@echo "3. Testing pressure sensor..."
	@curl -s http://localhost:8080/MySystem/PressureSensor1/pressure | head -c 100 && echo "..." || echo "ERROR: Pressure service not responding"
	@echo ""
	@echo "4. Testing controller..."
	@curl -s http://localhost:8080/MySystem/Controller1/control | head -c 100 && echo "..." || echo "ERROR: Controller service not responding"
	@echo ""
	@echo "All tests complete"
