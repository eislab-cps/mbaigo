package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var systemCmd = &cobra.Command{
	Use:   "system",
	Short: "System management commands",
	Long:  "Commands for managing Arrowhead systems",
}

func init() {
	systemCmd.AddCommand(systemInfoCmd)
	systemCmd.AddCommand(systemStartCmd)
}

var systemInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show system information",
	Long:  "Display information about the configured system",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Load config and display system information
		config, err := loadConfig()
		if err != nil {
			return err
		}

		if JSON {
			return printJSON(config)
		}

		fmt.Println("System Information")
		fmt.Println("=" + repeat("=", 60))

		if name, ok := config["systemname"].(string); ok {
			fmt.Printf("System Name:      %s\n", name)
		}
		if cloud, ok := config["localcloud"].(string); ok {
			fmt.Printf("Local Cloud:      %s\n", cloud)
		}

		if protocols, ok := config["protocolsNports"].(map[string]interface{}); ok {
			fmt.Println("\nEndpoints:")
			for proto, port := range protocols {
				fmt.Printf("  %s://localhost:%v\n", proto, port)
			}
		}

		if assets, ok := config["unit_assets"].([]interface{}); ok {
			fmt.Printf("\nUnit Assets:      %d\n", len(assets))
			for i, asset := range assets {
				if assetMap, ok := asset.(map[string]interface{}); ok {
					if name, ok := assetMap["name"].(string); ok {
						services := 0
						if svcs, ok := assetMap["services"].([]interface{}); ok {
							services = len(svcs)
						}
						fmt.Printf("  %d. %s (%d service%s)\n", i+1, name, services, pluralize(services))
					}
				}
			}
		}

		fmt.Println(repeat("=", 60))
		return nil
	},
}

var systemStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the system (example template)",
	Long: `Generate a template main.go file to start your system.

This creates a complete example showing how to:
  - Load the configuration
  - Create and configure the system
  - Setup HTTP/HTTPS servers
  - Register services
  - Handle graceful shutdown`,
	RunE: func(cmd *cobra.Command, args []string) error {
		template := `package main

import (
	"context"
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
	sys := components.NewSystem("%s", ctx)

	// Load configuration
	if _, err := usecases.Configure(&sys); err != nil {
		log.Fatalf("Configuration error: %%v", err)
	}

	// Setup servers and start services
	go usecases.SetoutServers(&sys)
	go usecases.RegisterServices(&sys)

	fmt.Printf("System started: %%s\n", sys.Name)
	fmt.Printf("HTTP endpoint: http://localhost:%%d\n", sys.Husk.ProtoPort["http"])

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	fmt.Println("\nShutting down gracefully...")
	cancel()
}
`
		config, err := loadConfig()
		if err != nil {
			return err
		}

		systemName := "MySystem"
		if name, ok := config["systemname"].(string); ok {
			systemName = name
		}

		output := fmt.Sprintf(template, systemName)

		if JSON {
			return printJSON(map[string]string{"template": output})
		}

		fmt.Println(output)

		fmt.Println("\nUsage:")
		fmt.Println("  Save this to main.go and run with: go run main.go")

		return nil
	},
}

func pluralize(count int) string {
	if count == 1 {
		return ""
	}
	return "s"
}
