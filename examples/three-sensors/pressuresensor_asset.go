package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/sdoque/mbaigo/pkg/components"
	"github.com/sdoque/mbaigo/pkg/forms"
	"github.com/sdoque/mbaigo/pkg/usecases"
)

// PressureSensorTraits defines the configurable parameters for this asset
type PressureSensorTraits struct {
	MinValue float64 `json:"minValue"`
	MaxValue float64 `json:"maxValue"`
}

// PressureSensor implements the UnitAsset interface
type PressureSensor struct {
	Name        string
	Details     map[string][]string
	ServicesMap components.Services
	CervicesMap components.Cervices
	Traits      PressureSensorTraits

	// Domain-specific fields
	currentPressure float64
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
	case "pressure":
		// Simulate pressure reading (customize this for real sensor)
		// Generate a value between MinValue and MaxValue
		pressure := a.Traits.MinValue + (a.Traits.MaxValue-a.Traits.MinValue)*0.5
		// Add some variation
		pressure += (0.5 - float64(time.Now().UnixNano()%100)/200) * 0.05

		signal := forms.SignalA_v1a{
			Value:     pressure,
			Unit:      "bar",
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

// NewPressureSensor creates a new instance from configuration
func NewPressureSensor(config usecases.ConfigurableAsset) components.UnitAsset {
	var traits PressureSensorTraits
	if len(config.Traits) > 0 {
		json.Unmarshal(config.Traits[0], &traits)
	}

	return &PressureSensor{
		Name:            config.Name,
		Details:         config.Details,
		ServicesMap:     usecases.MakeServiceMap(config.Services),
		Traits:          traits,
		currentPressure: 1.0, // Initial value (1 bar = atmospheric pressure)
	}
}
