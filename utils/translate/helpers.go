// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package translate

import (
	"fmt"
	"io/fs"
	"net/http"
	"path"
	"strings"
)

// FromSDKClient creates a translate.Client from an existing SDK client's
// HTTP client and base URL. This allows reusing the auth round-trippers.
func FromSDKClient(baseURL string, httpClient *http.Client, opts ...ClientOption) *Client {
	allOpts := append([]ClientOption{WithHTTPClient(httpClient)}, opts...)

	return NewClient(baseURL, allOpts...)
}

// LoadConfigsFromFS loads all config.yaml files from an embedded filesystem.
// It expects the layout: resources/<name>/config.yaml
func LoadConfigsFromFS(client *Client, fsys fs.FS, dir string) error {
	entries, err := fs.ReadDir(fsys, dir)
	if err != nil {
		return fmt.Errorf("reading config directory %q: %w", dir, err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		configPath := path.Join(dir, entry.Name(), "config.yaml")
		data, err := fs.ReadFile(fsys, configPath)
		if err != nil {
			continue // Skip resources without config.yaml
		}

		if err := client.RegisterResource(entry.Name(), data); err != nil {
			return fmt.Errorf("registering resource %q: %w", entry.Name(), err)
		}
	}

	return nil
}

// MustParseConfig parses a config.yaml and panics on error (for init-time use).
func MustParseConfig(data []byte) *ResourceConfig {
	cfg, err := ParseConfig(data)
	if err != nil {
		panic(fmt.Sprintf("failed to parse config: %v", err))
	}

	return cfg
}

// ResourcePath returns the API path pattern for a resource based on the envelope config.
// This is used when paths aren't explicitly configured.
func ResourcePath(envelope string) string {
	if envelope == "" {
		return ""
	}
	// Pluralize: zone → zones, environment → environments
	if strings.HasSuffix(envelope, "s") {
		return "/api/" + envelope
	}

	return "/api/" + envelope + "s"
}
