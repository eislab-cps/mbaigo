package cli

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/sdoque/mbaigo/pkg/core/esr"
	"github.com/sdoque/mbaigo/pkg/core/orchestrator"
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
			fmt.Println("No core systems found in systems/ directory")
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
	Path        string `json:"path"`
	Description string `json:"description,omitempty"`
	Port        int    `json:"port,omitempty"`
}

// getAvailableCoreSystems scans the systems/ directory for core systems
func getAvailableCoreSystems() ([]CoreSystemInfo, error) {
	systemsDir := "systems"

	// Check if systems directory exists
	if _, err := os.Stat(systemsDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("systems directory not found (are you in the mbaigo root directory?)")
	}

	entries, err := os.ReadDir(systemsDir)
	if err != nil {
		return nil, err
	}

	var coreSystems []CoreSystemInfo

	// Map of known core systems with their descriptions and default ports
	knownSystems := map[string]CoreSystemInfo{
		"esr": {
			Name:        "esr",
			Description: "Ephemeral Service Registry - tracks available services [EMBEDDED]",
			Port:        20102,
		},
		"orchestrator": {
			Name:        "orchestrator",
			Description: "Orchestrator - handles service discovery [EMBEDDED]",
			Port:        20103,
		},
		"messenger": {
			Name:        "messenger",
			Description: "Messenger - event handling and logging",
			Port:        20104,
		},
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		name := entry.Name()

		// Skip non-core directories
		if name == ".git" || name == "." || name == ".." {
			continue
		}

		systemPath := filepath.Join(systemsDir, name)

		// Check if it has a main.go or .go files
		hasGoFiles := false
		files, err := os.ReadDir(systemPath)
		if err == nil {
			for _, f := range files {
				if filepath.Ext(f.Name()) == ".go" && !f.IsDir() {
					hasGoFiles = true
					break
				}
			}
		}

		if !hasGoFiles {
			continue
		}

		// Use known info if available, otherwise create basic entry
		if info, exists := knownSystems[name]; exists {
			info.Path = systemPath
			coreSystems = append(coreSystems, info)
		} else {
			coreSystems = append(coreSystems, CoreSystemInfo{
				Name: name,
				Path: systemPath,
			})
		}
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

// startExternalSystem starts a core system from the systems/ directory
func startExternalSystem(systemName string) error {
	systemPath := filepath.Join("systems", systemName)

	// Check if the system directory exists
	if _, err := os.Stat(systemPath); os.IsNotExist(err) {
		return fmt.Errorf("core system '%s' not found (not embedded and not in systems/ directory)", systemName)
	}

	return fmt.Errorf("external systems from systems/ directory are not supported in embedded mode. Use 'cd systems/%s && go run .' instead", systemName)
}
