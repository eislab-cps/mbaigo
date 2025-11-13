package cli

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/eislab-cps/mbaigo/pkg/core/esr"
	"github.com/eislab-cps/mbaigo/pkg/core/orchestrator"
	"github.com/spf13/cobra"
)

var coreCmd = &cobra.Command{
	Use:   "core",
	Short: "Core system management commands",
	Long:  "Commands for managing Arrowhead core systems (ESR, Orchestrator, etc.)",
}

func init() {
	coreCmd.AddCommand(coreListCmd)
	coreCmd.AddCommand(coreStartCmd)
}

var coreListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available core systems",
	Long:  "Display all available core systems that can be started",
	RunE: func(cmd *cobra.Command, args []string) error {
		systems, err := getAvailableCoreSystems()
		if err != nil {
			return fmt.Errorf("failed to list core systems: %w", err)
		}

		if JSON {
			return printJSON(map[string]interface{}{
				"coreSystems": systems,
			})
		}

		fmt.Println("Available Core Systems")
		fmt.Println(repeat("=", 60))

		if len(systems) == 0 {
			fmt.Println("No core systems available")
			return nil
		}

		for _, sys := range systems {
			fmt.Printf("  • %s\n", sys.Name)
			if sys.Description != "" {
				fmt.Printf("    %s\n", sys.Description)
			}
			if sys.Port != 0 {
				fmt.Printf("    Default port: %d\n", sys.Port)
			}
			fmt.Println()
		}

		fmt.Println(repeat("=", 60))
		fmt.Println("\nUsage:")
		fmt.Println("  mbaigo core start <system-name>")

		return nil
	},
}

var coreStartCmd = &cobra.Command{
	Use:   "start [system-name]",
	Short: "Start a core system",
	Long: `Start a specific core system.

Available systems:
  esr          - Ephemeral Service Registry
  orchestrator - Service Orchestrator

Example:
  mbaigo core start esr
  mbaigo core start orchestrator`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		systemName := args[0]

		if Verbose {
			fmt.Printf("Starting core system: %s\n", systemName)
		}

		return startCoreSystem(systemName)
	},
}

// CoreSystemInfo holds information about a core system
type CoreSystemInfo struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Port        int    `json:"port,omitempty"`
}

// getAvailableCoreSystems returns the list of embedded core systems
func getAvailableCoreSystems() ([]CoreSystemInfo, error) {
	// Return embedded core systems directly
	// No need to scan systems/ directory - it's only for runtime data now
	coreSystems := []CoreSystemInfo{
		{
			Name:        "esr",
			Description: "Ephemeral Service Registry - tracks available services [EMBEDDED]",
			Port:        20102,
		},
		{
			Name:        "orchestrator",
			Description: "Orchestrator - handles service discovery [EMBEDDED]",
			Port:        20103,
		},
	}

	return coreSystems, nil
}

// startCoreSystem starts a specific core system (embedded)
func startCoreSystem(systemName string) error {
	fmt.Printf("Starting %s (embedded)...\n", systemName)
	fmt.Println(repeat("=", 60))

	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup signal handler for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Start the core system in a goroutine
	errChan := make(chan error, 1)
	go func() {
		var err error
		switch systemName {
		case "esr":
			err = esr.Start(ctx)
		case "orchestrator":
			err = orchestrator.Start(ctx)
		default:
			// Fall back to external systems directory if not embedded
			err = startExternalSystem(systemName)
		}
		errChan <- err
	}()

	// Wait for either completion or interrupt signal
	select {
	case err := <-errChan:
		if err != nil {
			return fmt.Errorf("core system %s failed: %w", systemName, err)
		}
		return nil
	case <-sigChan:
		fmt.Println("\nReceived interrupt signal, shutting down...")
		cancel()
		// Wait for graceful shutdown
		if err := <-errChan; err != nil {
			return fmt.Errorf("error during shutdown: %w", err)
		}
		return nil
	}
}

// startExternalSystem handles attempts to start non-embedded systems
func startExternalSystem(systemName string) error {
	return fmt.Errorf("core system '%s' not found. Available embedded systems: esr, orchestrator", systemName)
}
