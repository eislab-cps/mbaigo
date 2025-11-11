// Package main demonstrates a simple Arrowhead system with a single unit asset
// that provides a random number service.
package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/sdoque/mbaigo/pkg/components"
	"github.com/sdoque/mbaigo/pkg/forms"
	"github.com/sdoque/mbaigo/pkg/usecases"
)

// Traits defines the configurable parameters for our randomizer asset
type Traits struct {
	MinValue float64 `json:"minValue"`
	MaxValue float64 `json:"maxValue"`
}

// Randomizer is a simple unit asset that generates random numbers
type Randomizer struct {
	Name        string
	Details     map[string][]string
	ServicesMap components.Services
	CervicesMap components.Cervices
	Traits      Traits
}

// Implement the UnitAsset interface
func (r *Randomizer) GetName() string                  { return r.Name }
func (r *Randomizer) GetServices() components.Services { return r.ServicesMap }
func (r *Randomizer) GetCervices() components.Cervices { return r.CervicesMap }
func (r *Randomizer) GetDetails() map[string][]string  { return r.Details }
func (r *Randomizer) GetTraits() any                   { return r.Traits }

// Serving handles HTTP requests for this asset's services
func (r *Randomizer) Serving(w http.ResponseWriter, req *http.Request, servicePath string) {
	if servicePath != "random" {
		http.Error(w, "unknown service path", http.StatusNotFound)
		return
	}

	// Generate random number and create response
	value := rand.Float64()*(r.Traits.MaxValue-r.Traits.MinValue) + r.Traits.MinValue

	signal := forms.SignalA_v1a{
		Value: value,
		Unit:  "float64",
	}

	// Pack and send response
	packed, err := usecases.Pack(signal.NewForm(), "application/json")
	if err != nil {
		http.Error(w, "failed to pack response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(packed)
}

func main() {
	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a new Arrowhead system
	sys := components.NewSystem("RandomizerSystem", ctx)

	// Configure the Husk (middleware layer)
	sys.Husk = &components.Husk{
		Description: "Simple randomizer system example",
		ProtoPort:   map[string]int{"http": 8080},
	}

	// Create and configure the randomizer unit asset
	service := &components.Service{
		Definition:  "RandomNumber",
		SubPath:     "random",
		RegPeriod:   60,
		Description: "returns a random float64 number",
		Details:     map[string][]string{"type": {"float64"}},
	}

	randomizer := &Randomizer{
		Name:        "randomizer",
		Details:     map[string][]string{"location": {"local"}},
		ServicesMap: components.Services{"random": service},
		Traits:      Traits{MinValue: 0.0, MaxValue: 100.0},
	}

	// Add asset to system
	asset := components.UnitAsset(randomizer)
	sys.UAssets[randomizer.GetName()] = &asset

	// Load or create configuration
	if _, err := usecases.Configure(&sys); err != nil {
		log.Printf("Configuration warning: %v", err)
	}

	// Start services
	go usecases.SetoutServers(&sys)
	go usecases.RegisterServices(&sys)

	fmt.Println("Randomizer system started on port 8080")
	fmt.Printf("Try: curl http://localhost:8080/%s/%s/%s\n",
		sys.Name, randomizer.GetName(), "random")

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	fmt.Println("\nShutting down gracefully...")
	cancel()
}
