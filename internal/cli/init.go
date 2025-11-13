package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	initSystemName  string
	initLocalCloud  string
	initHTTPPort    int
	initHTTPSPort   int
	initInteractive bool
)

func init() {
	initCmd.Flags().StringVarP(&initSystemName, "name", "n", "", "System name (required)")
	initCmd.Flags().StringVarP(&initLocalCloud, "cloud", "l", "", "Local cloud name (required)")
	initCmd.Flags().IntVar(&initHTTPPort, "http-port", 8080, "HTTP port")
	initCmd.Flags().IntVar(&initHTTPSPort, "https-port", 8443, "HTTPS port")
	initCmd.Flags().BoolVarP(&initInteractive, "interactive", "i", false, "Interactive mode")

	initCmd.MarkFlagRequired("name")
	initCmd.MarkFlagRequired("cloud")
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new mbaigo system configuration",
	Long: `Initialize a new Arrowhead system by creating a systemconfig.json file.

This command creates a basic configuration file with:
  - System name and local cloud configuration
  - Protocol and port bindings
  - Core system placeholders
  - Empty unit assets array

Example:
  mbaigo init --name MySystem --cloud FactoryCloud
  mbaigo init -n SensorSystem -l Building1 --http-port 9090`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if initInteractive {
			return runInteractiveInit()
		}
		return runInit()
	},
}

func runInit() error {
	config := map[string]interface{}{
		"systemname": initSystemName,
		"localcloud": initLocalCloud,
		"protocolsNports": map[string]int{
			"http":  initHTTPPort,
			"https": initHTTPSPort,
		},
		"coreSystems": []map[string]string{
			{
				"coreSystem": "serviceregistrar",
				"url":        "https://localhost:8000",
			},
			{
				"coreSystem": "orchestrator",
				"url":        "https://localhost:8001",
			},
		},
		"unit_assets": []interface{}{},
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Check if file exists
	if _, err := os.Stat(ConfigFile); err == nil {
		return fmt.Errorf("config file already exists: %s (use --config to specify a different file)", ConfigFile)
	}

	if err := os.WriteFile(ConfigFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	if !JSON {
		printSuccess(fmt.Sprintf("Created configuration file: %s", ConfigFile))
		fmt.Println("\nNext steps:")
		fmt.Println("  1. Review and edit the configuration file")
		fmt.Println("  2. Add unit assets with: mbaigo generate asset")
		fmt.Println("  3. Validate config with: mbaigo config validate")
		fmt.Println("  4. Start your system")
	} else {
		fmt.Printf(`{"status":"success","file":"%s"}%s`, ConfigFile, "\n")
	}

	return nil
}

func runInteractiveInit() error {
	// Interactive prompts for configuration
	fmt.Println("Interactive System Initialization")
	fmt.Println()

	// Prompt for system name
	if initSystemName == "" {
		fmt.Print("System name: ")
		fmt.Scanln(&initSystemName)
	}

	// Prompt for local cloud
	if initLocalCloud == "" {
		fmt.Print("Local cloud name: ")
		fmt.Scanln(&initLocalCloud)
	}

	// Prompt for HTTP port
	fmt.Printf("HTTP port (default: %d): ", initHTTPPort)
	var portInput string
	fmt.Scanln(&portInput)
	if portInput != "" {
		fmt.Sscanf(portInput, "%d", &initHTTPPort)
	}

	// Prompt for HTTPS port
	fmt.Printf("HTTPS port (default: %d): ", initHTTPSPort)
	fmt.Scanln(&portInput)
	if portInput != "" {
		fmt.Sscanf(portInput, "%d", &initHTTPSPort)
	}

	fmt.Println()
	return runInit()
}
