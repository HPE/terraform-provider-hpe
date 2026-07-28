// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

// Package specs embeds the per-resource and per-datasource config.yaml files
// that are the single source of truth for both code generation (go-split) and
// the runtime translation layer.
package specs

import "embed"

//go:embed resources/*/config.yaml datasources/*/config.yaml
var FS embed.FS

// ResourceConfig returns the raw config.yaml bytes for a named resource.
func ResourceConfig(name string) ([]byte, error) {
	return FS.ReadFile("resources/" + name + "/config.yaml")
}

// DataSourceConfig returns the raw config.yaml bytes for a named datasource.
func DataSourceConfig(name string) ([]byte, error) {
	return FS.ReadFile("datasources/" + name + "/config.yaml")
}
