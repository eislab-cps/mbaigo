package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	serviceAssetName string
	serviceName      string
)

func init() {
	serviceCmd.AddCommand(serviceListCmd)
	serviceCmd.AddCommand(serviceAddCmd)

	serviceListCmd.Flags().StringVarP(&serviceAssetName, "asset", "a", "", "Filter by asset name")

	serviceAddCmd.Flags().StringVarP(&serviceAssetName, "asset", "a", "", "Asset name (required)")
	serviceAddCmd.Flags().StringVarP(&serviceName, "name", "n", "", "Service name (required)")
	serviceAddCmd.MarkFlagRequired("asset")
	serviceAddCmd.MarkFlagRequired("name")
}

var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "Service management commands",
	Long:  "Commands for managing services provided by unit assets",
}

var serviceListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all services",
	Long:  "List all services from configured unit assets",
	RunE: func(cmd *cobra.Command, args []string) error {
		config, err := loadConfig()
		if err != nil {
			return err
		}

		assets, ok := config["unit_assets"].([]interface{})
		if !ok {
			return fmt.Errorf("no unit assets found in configuration")
		}

		if JSON {
			services := []map[string]interface{}{}
			for _, asset := range assets {
				assetMap := asset.(map[string]interface{})
				assetName := assetMap["name"].(string)

				if serviceAssetName != "" && assetName != serviceAssetName {
					continue
				}

				if svcs, ok := assetMap["services"].([]interface{}); ok {
					for _, svc := range svcs {
						svcMap := svc.(map[string]interface{})
						svcMap["asset"] = assetName
						services = append(services, svcMap)
					}
				}
			}
			return printJSON(map[string]interface{}{"services": services})
		}

		fmt.Println("Configured Services")
		fmt.Println("=" + repeat("=", 80))
		fmt.Printf("%-20s %-20s %-15s %-10s\n", "ASSET", "SERVICE", "SUBPATH", "REG PERIOD")
		fmt.Println(repeat("-", 80))

		totalServices := 0
		for _, asset := range assets {
			assetMap := asset.(map[string]interface{})
			assetName := assetMap["name"].(string)

			if serviceAssetName != "" && assetName != serviceAssetName {
				continue
			}

			if svcs, ok := assetMap["services"].([]interface{}); ok {
				for _, svc := range svcs {
					svcMap := svc.(map[string]interface{})
					definition := getString(svcMap, "definition", "-")
					subpath := getString(svcMap, "subpath", "-")
					regPeriod := getInt(svcMap, "registrationPeriod", 0)

					fmt.Printf("%-20s %-20s %-15s %-10ds\n",
						truncate(assetName, 20),
						truncate(definition, 20),
						truncate(subpath, 15),
						regPeriod)
					totalServices++
				}
			}
		}

		fmt.Println(repeat("-", 80))
		fmt.Printf("Total: %d service%s\n", totalServices, pluralize(totalServices))

		return nil
	},
}

var serviceAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a service to an asset",
	Long:  "Add a new service definition to a unit asset",
	RunE: func(cmd *cobra.Command, args []string) error {
		config, err := loadConfig()
		if err != nil {
			return err
		}

		assets, ok := config["unit_assets"].([]interface{})
		if !ok {
			return fmt.Errorf("no unit assets found in configuration")
		}

		// Find the asset
		found := false
		for i, asset := range assets {
			assetMap := asset.(map[string]interface{})
			if assetMap["name"].(string) == serviceAssetName {
				found = true

				// Add new service
				newService := map[string]interface{}{
					"definition":         serviceName,
					"subpath":           serviceName,
					"registrationPeriod": 60,
					"details": map[string][]string{
						"type": {"service"},
					},
					"description": "Service description here",
				}

				if svcs, ok := assetMap["services"].([]interface{}); ok {
					assetMap["services"] = append(svcs, newService)
				} else {
					assetMap["services"] = []interface{}{newService}
				}

				assets[i] = assetMap
				break
			}
		}

		if !found {
			return fmt.Errorf("asset not found: %s", serviceAssetName)
		}

		config["unit_assets"] = assets

		if err := saveConfig(config); err != nil {
			return err
		}

		if JSON {
			return printJSON(map[string]string{"status": "success", "asset": serviceAssetName, "service": serviceName})
		}

		printSuccess(fmt.Sprintf("Added service '%s' to asset '%s'", serviceName, serviceAssetName))
		fmt.Println("\nNext: Edit the service details in systemconfig.json")

		return nil
	},
}
