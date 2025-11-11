package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
)

func colorize(color, text string) string {
	if NoColor {
		return text
	}
	return color + text + colorReset
}

func printSuccess(msg string) {
	fmt.Println(colorize(colorGreen, "✓ " + msg))
}

func printError(msg string) {
	fmt.Println(colorize(colorRed, "✗ " + msg))
}

func printWarning(msg string) {
	fmt.Println(colorize(colorYellow, "⚠ " + msg))
}

func printInfo(msg string) {
	fmt.Println(colorize(colorCyan, "ℹ " + msg))
}

func repeat(s string, count int) string {
	return strings.Repeat(s, count)
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func getString(m map[string]interface{}, key string, defaultVal string) string {
	if val, ok := m[key].(string); ok {
		return val
	}
	return defaultVal
}

func getInt(m map[string]interface{}, key string, defaultVal int) int {
	if val, ok := m[key].(float64); ok {
		return int(val)
	}
	return defaultVal
}

func printJSON(data interface{}) error {
	output, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	fmt.Println(string(output))
	return nil
}

func loadConfig() (map[string]interface{}, error) {
	data, err := os.ReadFile(ConfigFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file '%s': %w (use 'mbaigo init' to create)", ConfigFile, err)
	}

	var config map[string]interface{}
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return config, nil
}

func saveConfig(config map[string]interface{}) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(ConfigFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}
