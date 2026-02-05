package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/goccy/go-yaml"
)

// LoadConfigs loads all YAML config files from a directory
func LoadConfigs(configDir string) (map[string]CleanupConfig, error) {
	configs := make(map[string]CleanupConfig)

	entries, err := os.ReadDir(configDir)
	if err != nil {
		return nil, fmt.Errorf("read config dir: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if !strings.HasSuffix(entry.Name(), ".yaml") && !strings.HasSuffix(entry.Name(), ".yml") {
			continue
		}

		configPath := filepath.Join(configDir, entry.Name())
		fileConfigs, err := loadConfigFile(configPath)
		if err != nil {
			return nil, fmt.Errorf("load %s: %w", entry.Name(), err)
		}

		// Merge configs from this file
		for resourceType, config := range fileConfigs {
			configs[resourceType] = config
		}
	}

	return configs, nil
}

// loadConfigFile loads a single YAML config file
// File format:
// resource_type_name:
//
//	extracts: [...]
//	removes: [...]
//	...
func loadConfigFile(path string) (map[string]CleanupConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	// Parse as map[string]CleanupConfig
	var configs map[string]CleanupConfig
	if err := yaml.Unmarshal(data, &configs); err != nil {
		return nil, fmt.Errorf("parse yaml: %w", err)
	}

	// Set ResourceType field for each config
	for resourceType, config := range configs {
		config.ResourceType = resourceType
		configs[resourceType] = config
	}

	return configs, nil
}

// GetConfigForResource gets the cleanup config for a specific resource type
func GetConfigForResource(configs map[string]CleanupConfig, resourceType string) (CleanupConfig, bool) {
	config, exists := configs[resourceType]
	return config, exists
}
