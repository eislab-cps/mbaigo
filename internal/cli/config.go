package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configValidateCmd)
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configuration management commands",
	Long:  "Commands for managing system configuration files",
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Long:  "Display the current system configuration from systemconfig.json",
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := os.ReadFile(ConfigFile)
		if err != nil {
			return fmt.Errorf("failed to read config file: %w", err)
		}

		if JSON {
			// Just output the raw JSON
			fmt.Println(string(data))
		} else {
			// Pretty print the configuration
			var config map[string]interface{}
			if err := json.Unmarshal(data, &config); err != nil {
				return fmt.Errorf("failed to parse config: %w", err)
			}

			printConfigTable(config)
		}

		return nil
	},
}

var configValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate configuration file",
	Long:  "Validate the systemconfig.json file for correctness",
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := os.ReadFile(ConfigFile)
		if err != nil {
			return fmt.Errorf("failed to read config file: %w", err)
		}

		var config map[string]interface{}
		if err := json.Unmarshal(data, &config); err != nil {
			return fmt.Errorf("invalid JSON: %w", err)
		}

		// Validate required fields
		errors := []string{}

		if _, ok := config["systemname"]; !ok {
			errors = append(errors, "missing required field: systemname")
		}
		if _, ok := config["localcloud"]; !ok {
			errors = append(errors, "missing required field: localcloud")
		}
		if _, ok := config["protocolsNports"]; !ok {
			errors = append(errors, "missing required field: protocolsNports")
		}
		if _, ok := config["coreSystems"]; !ok {
			errors = append(errors, "missing required field: coreSystems")
		}
		if _, ok := config["unit_assets"]; !ok {
			errors = append(errors, "missing required field: unit_assets")
		}

		if len(errors) > 0 {
			if JSON {
				fmt.Printf(`{"valid":false,"errors":%s}%s`, toJSONArray(errors), "\n")
			} else {
				printError("Configuration validation failed:")
				for _, err := range errors {
					fmt.Printf("  ✗ %s\n", err)
				}
			}
			return fmt.Errorf("validation failed")
		}

		if JSON {
			fmt.Println(`{"valid":true}`)
		} else {
			printSuccess("✓ Configuration is valid")
		}

		return nil
	},
}

func printConfigTable(config map[string]interface{}) {
	fmt.Println("📋 System Configuration")
	fmt.Println("=" + repeat("=", 60))

	if systemName, ok := config["systemname"].(string); ok {
		fmt.Printf("System Name:    %s\n", systemName)
	}
	if localCloud, ok := config["localcloud"].(string); ok {
		fmt.Printf("Local Cloud:    %s\n", localCloud)
	}

	if protocols, ok := config["protocolsNports"].(map[string]interface{}); ok {
		fmt.Println("\nProtocols & Ports:")
		for proto, port := range protocols {
			fmt.Printf("  %s: %v\n", proto, port)
		}
	}

	if coreSystems, ok := config["coreSystems"].([]interface{}); ok {
		fmt.Println("\nCore Systems:")
		for _, cs := range coreSystems {
			if csMap, ok := cs.(map[string]interface{}); ok {
				name := csMap["coreSystem"]
				url := csMap["url"]
				fmt.Printf("  %s: %s\n", name, url)
			}
		}
	}

	if assets, ok := config["unit_assets"].([]interface{}); ok {
		fmt.Printf("\nUnit Assets:    %d configured\n", len(assets))
	}

	fmt.Println(repeat("=", 60))
}

func toJSONArray(arr []string) string {
	data, _ := json.Marshal(arr)
	return string(data)
}
