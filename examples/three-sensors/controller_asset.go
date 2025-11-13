package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/eislab-cps/mbaigo/pkg/components"
	"github.com/eislab-cps/mbaigo/pkg/forms"
	"github.com/eislab-cps/mbaigo/pkg/usecases"
)

// ControllerTraits defines the configurable parameters for this asset
type ControllerTraits struct {
	// Add your custom configuration fields here
}

// Controller implements the UnitAsset interface
type Controller struct {
	Name        string
	Details     map[string][]string
	ServicesMap components.Services
	CervicesMap components.Cervices
	Traits      ControllerTraits

	// Domain-specific fields
	lastCommand string
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
	case "control":
		// Return controller status
		signal := forms.SignalB_v1a{
			Value:     true, // Controller is active
			Timestamp: time.Now(),
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
		lastCommand: "none",
	}
}
