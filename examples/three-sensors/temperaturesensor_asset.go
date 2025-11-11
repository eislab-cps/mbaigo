package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/sdoque/mbaigo/pkg/components"
	"github.com/sdoque/mbaigo/pkg/forms"
	"github.com/sdoque/mbaigo/pkg/usecases"
)

// TemperatureSensorTraits defines the configurable parameters for this asset
type TemperatureSensorTraits struct {
	MinValue float64 `json:"minValue"`
	MaxValue float64 `json:"maxValue"`
}

// TemperatureSensor implements the UnitAsset interface
type TemperatureSensor struct {
	Name        string
	Details     map[string][]string
	ServicesMap components.Services
	CervicesMap components.Cervices
	Traits      TemperatureSensorTraits

	// Domain-specific fields
	currentTemp float64
}

// Implement the UnitAsset interface
func (a *TemperatureSensor) GetName() string                  { return a.Name }
func (a *TemperatureSensor) GetServices() components.Services { return a.ServicesMap }
func (a *TemperatureSensor) GetCervices() components.Cervices { return a.CervicesMap }
func (a *TemperatureSensor) GetDetails() map[string][]string  { return a.Details }
func (a *TemperatureSensor) GetTraits() any                   { return a.Traits }

// Serving handles HTTP requests for this asset's services
func (a *TemperatureSensor) Serving(w http.ResponseWriter, r *http.Request, servicePath string) {
	switch servicePath {
	case "temperature":
		// Simulate temperature reading (customize this for real sensor)
		// Generate a value between MinValue and MaxValue
		temp := a.Traits.MinValue + (a.Traits.MaxValue-a.Traits.MinValue)*0.5
		// Add some variation
		temp += (0.5 - float64(time.Now().UnixNano()%100)/200) * 2

		signal := forms.SignalA_v1a{
			Value:     temp,
			Unit:      "celsius",
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

// NewTemperatureSensor creates a new instance from configuration
func NewTemperatureSensor(config usecases.ConfigurableAsset) components.UnitAsset {
	var traits TemperatureSensorTraits
	if len(config.Traits) > 0 {
		json.Unmarshal(config.Traits[0], &traits)
	}

	return &TemperatureSensor{
		Name:        config.Name,
		Details:     config.Details,
		ServicesMap: usecases.MakeServiceMap(config.Services),
		Traits:      traits,
		currentTemp: 22.0, // Initial value
	}
}
