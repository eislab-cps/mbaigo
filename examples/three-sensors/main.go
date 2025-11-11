package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/sdoque/mbaigo/pkg/components"
	"github.com/sdoque/mbaigo/pkg/usecases"
)

func main() {
	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create system
	sys := components.NewSystem("MySystem", ctx)

	// Initialize Husk (system container/description)
	sys.Husk = &components.Husk{
		Description: "My mbaigo system",
		Details:     map[string][]string{"LocalCloud": {"LocalCloud"}},
		ProtoPort:   map[string]int{"http": 8080},
	}

	// Register asset factories (these functions are defined in *_asset.go files)
	usecases.RegisterAssetFactory("TemperatureSensor", NewTemperatureSensor)
	usecases.RegisterAssetFactory("PressureSensor", NewPressureSensor)
	usecases.RegisterAssetFactory("Controller", NewController)

	// Create a dummy asset map to satisfy Configure() requirements
	// (Configure() calls setupDefaultConfig which needs at least one asset,
	// even though the config already exists and the dummy won't be used)
	sys.UAssets = make(map[string]*components.UnitAsset)
	dummyConfig := usecases.ConfigurableAsset{
		Name:    "dummy",
		Details: map[string][]string{"type": {"TemperatureSensor"}},
		Services: []components.Service{
			{SubPath: "dummy", Definition: "dummy"},
		},
	}
	dummyAsset := NewTemperatureSensor(dummyConfig)
	sys.UAssets["dummy"] = &dummyAsset

	// Load configuration
	rawAssets, err := usecases.Configure(&sys)
	if err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	// Load unit assets from configuration
	sys.UAssets = make(map[string]*components.UnitAsset)

	for _, rawAsset := range rawAssets {
		// Unmarshal the raw JSON into ConfigurableAsset
		var ca usecases.ConfigurableAsset
		if err := json.Unmarshal(rawAsset, &ca); err != nil {
			log.Printf("Error unmarshaling asset: %v", err)
			continue
		}

		// Get the asset type from details
		assetTypes, ok := ca.Details["type"]
		if !ok || len(assetTypes) == 0 {
			log.Printf("Warning: Asset %s has no type specified", ca.Name)
			continue
		}
		assetType := assetTypes[0]

		// Get the factory function
		factory, ok := usecases.GetAssetFactory(assetType)
		if !ok {
			log.Printf("Warning: Unknown asset type: %s", assetType)
			continue
		}

		// Create the asset using the factory
		asset := factory(ca)
		sys.UAssets[asset.GetName()] = &asset
	}

	// Setup servers and start services
	go usecases.SetoutServers(&sys)
	go usecases.RegisterServices(&sys)

	fmt.Println("\n🚀 System Started!")
	fmt.Println("=" + repeat("=", 60))
	fmt.Printf("System Name:  %s\n", sys.Name)

	// Get LocalCloud from Husk details if available
	if localCloud, ok := sys.Husk.Details["LocalCloud"]; ok && len(localCloud) > 0 {
		fmt.Printf("Local Cloud:  %s\n", localCloud[0])
	}

	fmt.Printf("HTTP:         http://localhost:%d\n", sys.Husk.ProtoPort["http"])
	fmt.Println("\n📍 Available Endpoints:")

	for name, asset := range sys.UAssets {
		services := (*asset).GetServices()
		for subpath := range services {
			fmt.Printf("  GET /%s/%s/%s\n", sys.Name, name, subpath)
		}
	}

	fmt.Println(repeat("=", 60))
	fmt.Println("\n💡 Test with:")
	for name, asset := range sys.UAssets {
		services := (*asset).GetServices()
		for subpath := range services {
			fmt.Printf("  curl http://localhost:%d/%s/%s/%s\n",
				sys.Husk.ProtoPort["http"], sys.Name, name, subpath)
			break // Just show one example per asset
		}
	}
	fmt.Println("\nPress Ctrl+C to stop...")

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	fmt.Println("\n\n🛑 Shutting down gracefully...")
	cancel()
}

func repeat(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}
