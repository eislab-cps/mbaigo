package main

import (
	"encoding/json"
	"net/http"

	"github.com/sdoque/mbaigo/pkg/components"
	"github.com/sdoque/mbaigo/pkg/forms"
	"github.com/sdoque/mbaigo/pkg/usecases"
)

// PressureSensorTraits defines the configurable parameters for this asset
type PressureSensorTraits struct {
	// Add your custom configuration fields here
	// Example:
	// MinValue float64 `json:"minValue"`
	// MaxValue float64 `json:"maxValue"`
}

// PressureSensor implements the UnitAsset interface
type PressureSensor struct {
	Name        string
	Details     map[string][]string
	ServicesMap components.Services
	CervicesMap components.Cervices
	Traits      PressureSensorTraits

	// Add your domain-specific fields here
}

// Implement the UnitAsset interface
func (a *PressureSensor) GetName() string                  { return a.Name }
func (a *PressureSensor) GetServices() components.Services { return a.ServicesMap }
func (a *PressureSensor) GetCervices() components.Cervices { return a.CervicesMap }
func (a *PressureSensor) GetDetails() map[string][]string  { return a.Details }
func (a *PressureSensor) GetTraits() any                   { return a.Traits }

// Serving handles HTTP requests for this asset's services
func (a *PressureSensor) Serving(w http.ResponseWriter, r *http.Request, servicePath string) {
	switch servicePath {
	case "example":
		// Handle the "example" service
		signal := forms.SignalA_v1a{
			Value: 42.0,
			Unit:  "units",
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

// NewPressureSensor creates a new instance from configuration
func NewPressureSensor(config usecases.ConfigurableAsset) components.UnitAsset {
	var traits PressureSensorTraits
	if len(config.Traits) > 0 {
		json.Unmarshal(config.Traits[0], &traits)
	}

	return &PressureSensor{
		Name:        config.Name,
		Details:     config.Details,
		ServicesMap: usecases.MakeServiceMap(config.Services),
		Traits:      traits,
		// Initialize your custom fields here
	}
}

func init() {
	// Register the factory function
	usecases.RegisterAssetFactory("PressureSensor", NewPressureSensor)
}
