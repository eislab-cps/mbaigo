package main

import (
	"encoding/json"
	"net/http"

	"github.com/sdoque/mbaigo/pkg/components"
	"github.com/sdoque/mbaigo/pkg/forms"
	"github.com/sdoque/mbaigo/pkg/usecases"
)

// ControllerTraits defines the configurable parameters for this asset
type ControllerTraits struct {
	// Add your custom configuration fields here
	// Example:
	// MinValue float64 `json:"minValue"`
	// MaxValue float64 `json:"maxValue"`
}

// Controller implements the UnitAsset interface
type Controller struct {
	Name        string
	Details     map[string][]string
	ServicesMap components.Services
	CervicesMap components.Cervices
	Traits      ControllerTraits

	// Add your domain-specific fields here
}

// Implement the UnitAsset interface
func (a *Controller) GetName() string                  { return a.Name }
func (a *Controller) GetServices() components.Services { return a.ServicesMap }
func (a *Controller) GetCervices() components.Cervices { return a.CervicesMap }
func (a *Controller) GetDetails() map[string][]string  { return a.Details }
func (a *Controller) GetTraits() any                   { return a.Traits }

// Serving handles HTTP requests for this asset's services
func (a *Controller) Serving(w http.ResponseWriter, r *http.Request, servicePath string) {
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

// NewController creates a new instance from configuration
func NewController(config usecases.ConfigurableAsset) components.UnitAsset {
	var traits ControllerTraits
	if len(config.Traits) > 0 {
		json.Unmarshal(config.Traits[0], &traits)
	}

	return &Controller{
		Name:        config.Name,
		Details:     config.Details,
		ServicesMap: usecases.MakeServiceMap(config.Services),
		Traits:      traits,
		// Initialize your custom fields here
	}
}

func init() {
	// Register the factory function
	usecases.RegisterAssetFactory("Controller", NewController)
}
