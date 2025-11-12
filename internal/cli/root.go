package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

const Version = "0.1.0"

var (
	// Global flags
	Verbose    bool
	ConfigFile string
	JSON       bool
	NoColor    bool
)

func init() {
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Enable verbose output")
	rootCmd.PersistentFlags().StringVarP(&ConfigFile, "config", "c", "systemconfig.json", "Config file path")
	rootCmd.PersistentFlags().BoolVarP(&JSON, "json", "j", false, "Output as JSON")
	rootCmd.PersistentFlags().BoolVar(&NoColor, "no-color", false, "Disable colored output")

	// Add subcommands
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(systemCmd)
	rootCmd.AddCommand(serviceCmd)
	rootCmd.AddCommand(certCmd)
	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(coreCmd)
}

var rootCmd = &cobra.Command{
	Use:   "mbaigo",
	Short: "mbaigo - CLI for Arrowhead Framework systems",
	Long: `mbaigo is a command-line tool for managing Arrowhead Framework systems.

It provides commands for:
  - Initializing and managing systems
  - Registering and discovering services
  - Managing certificates
  - Generating configuration and code templates`,
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Run: func(cmd *cobra.Command, args []string) {
		if JSON {
			fmt.Printf(`{"version":"%s"}%s`, Version, "\n")
		} else {
			fmt.Printf("mbaigo version %s\n", Version)
		}
	},
}
