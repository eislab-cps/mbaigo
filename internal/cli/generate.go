package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	genAssetName string
	genOutput    string
)

func init() {
	generateCmd.AddCommand(genAssetCmd)
	generateCmd.AddCommand(genExampleCmd)

	genAssetCmd.Flags().StringVarP(&genAssetName, "name", "n", "", "Asset name (required)")
	genAssetCmd.Flags().StringVarP(&genOutput, "output", "o", "", "Output file (default: <name>_asset.go)")
	genAssetCmd.MarkFlagRequired("name")
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Code generation commands",
	Long:  "Generate code templates for assets and systems",
}

var genAssetCmd = &cobra.Command{
	Use:   "asset",
	Short: "Generate a unit asset template",
	Long:  "Generate a Go source file template for implementing a custom unit asset",
	RunE: func(cmd *cobra.Command, args []string) error {
		if genOutput == "" {
			genOutput = fmt.Sprintf("%s_asset.go", strings.ToLower(genAssetName))
		}

		template := fmt.Sprintf(`package main

import (
	"encoding/json"
	"net/http"

	"github.com/sdoque/mbaigo/pkg/components"
	"github.com/sdoque/mbaigo/pkg/forms"
	"github.com/sdoque/mbaigo/pkg/usecases"
)

// %sTraits defines the configurable parameters for this asset
type %sTraits struct {
	// Add your custom configuration fields here
	// Example:
	// MinValue float64 ` + "`json:\"minValue\"`" + `
	// MaxValue float64 ` + "`json:\"maxValue\"`" + `
}

// %s implements the UnitAsset interface
type %s struct {
	Name        string
	Details     map[string][]string
	ServicesMap components.Services
	CervicesMap components.Cervices
	Traits      %sTraits

	// Add your domain-specific fields here
}

// Implement the UnitAsset interface
func (a *%s) GetName() string                  { return a.Name }
func (a *%s) GetServices() components.Services { return a.ServicesMap }
func (a *%s) GetCervices() components.Cervices { return a.CervicesMap }
func (a *%s) GetDetails() map[string][]string  { return a.Details }
func (a *%s) GetTraits() any                   { return a.Traits }

// Serving handles HTTP requests for this asset's services
func (a *%s) Serving(w http.ResponseWriter, r *http.Request, servicePath string) {
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

// New%s creates a new instance from configuration
func New%s(config usecases.ConfigurableAsset) components.UnitAsset {
	var traits %sTraits
	if len(config.Traits) > 0 {
		json.Unmarshal(config.Traits[0], &traits)
	}

	return &%s{
		Name:        config.Name,
		Details:     config.Details,
		ServicesMap: usecases.MakeServiceMap(config.Services),
		Traits:      traits,
		// Initialize your custom fields here
	}
}

func init() {
	// Register the factory function
	usecases.RegisterAssetFactory("%s", New%s)
}
`, genAssetName, genAssetName, genAssetName, genAssetName, genAssetName,
			genAssetName, genAssetName, genAssetName, genAssetName, genAssetName, genAssetName,
			genAssetName, genAssetName, genAssetName, genAssetName, genAssetName, genAssetName)

		if err := os.WriteFile(genOutput, []byte(template), 0644); err != nil {
			return fmt.Errorf("failed to write file: %w", err)
		}

		if JSON {
			return printJSON(map[string]string{"file": genOutput, "asset": genAssetName})
		}

		printSuccess(fmt.Sprintf("Generated asset template: %s", genOutput))
		fmt.Println("\nNext steps:")
		fmt.Println("  1. Edit the generated file and implement your logic")
		fmt.Println("  2. Add the asset to systemconfig.json")
		fmt.Println("  3. Import and use in your main.go")

		return nil
	},
}

var genExampleCmd = &cobra.Command{
	Use:   "example",
	Short: "Generate a complete example application",
	Long:  "Generate a complete example application with a working system",
	RunE: func(cmd *cobra.Command, args []string) error {
		config, err := loadConfig()
		if err != nil {
			return err
		}

		systemName := "MySystem"
		if name, ok := config["systemname"].(string); ok {
			systemName = name
		}

		example := fmt.Sprintf(`package main

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

// ExampleAsset provides a simple random number service
type ExampleAsset struct {
	Name        string
	ServicesMap components.Services
	Details     map[string][]string
}

func (a *ExampleAsset) GetName() string                  { return a.Name }
func (a *ExampleAsset) GetServices() components.Services { return a.ServicesMap }
func (a *ExampleAsset) GetCervices() components.Cervices { return components.Cervices{} }
func (a *ExampleAsset) GetDetails() map[string][]string  { return a.Details }
func (a *ExampleAsset) GetTraits() any                   { return nil }

func (a *ExampleAsset) Serving(w http.ResponseWriter, r *http.Request, servicePath string) {
	signal := forms.SignalA_v1a{
		Value: rand.Float64() * 100,
		Unit:  "random",
	}

	packed, _ := usecases.Pack(signal.NewForm(), "application/json")
	w.Header().Set("Content-Type", "application/json")
	w.Write(packed)
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create system
	sys := components.NewSystem("%s", ctx)

	// Configure
	if _, err := usecases.Configure(&sys); err != nil {
		log.Fatalf("Configuration error: %%v", err)
	}

	// Add example asset
	service := &components.Service{
		Definition:  "RandomNumber",
		SubPath:     "random",
		RegPeriod:   60,
		Description: "returns a random number",
	}

	asset := components.UnitAsset(&ExampleAsset{
		Name:        "example",
		ServicesMap: components.Services{"random": service},
		Details:     map[string][]string{"type": {"example"}},
	})

	sys.UAssets["example"] = &asset

	// Start services
	go usecases.SetoutServers(&sys)
	go usecases.RegisterServices(&sys)

	fmt.Printf("System started: %%s\n", sys.Name)
	if port, ok := sys.Husk.ProtoPort["http"]; ok {
		fmt.Printf("Try: curl http://localhost:%%d/%%s/example/random\n", port, sys.Name)
	}

	// Wait for signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	fmt.Println("\nShutting down...")
	cancel()
}
`, systemName)

		filename := "example_main.go"
		if err := os.WriteFile(filename, []byte(example), 0644); err != nil {
			return fmt.Errorf("failed to write file: %w", err)
		}

		if JSON {
			return printJSON(map[string]string{"file": filename})
		}

		printSuccess(fmt.Sprintf("Generated example application: %s", filename))
		fmt.Println("\nRun with: go run " + filename)

		return nil
	},
}
