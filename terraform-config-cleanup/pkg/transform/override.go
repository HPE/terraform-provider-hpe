package transform

import (
	"fmt"

	"github.com/HPE/terraform-config-cleanup/pkg/config"
	"github.com/HPE/terraform-config-cleanup/pkg/parser"
)

// ApplyOverrides applies override operations to a resource
// Overrides set specific values at given paths
func ApplyOverrides(resource *parser.Resource, overrides map[string]config.OverrideOperation) error {
	for path, override := range overrides {
		if err := applyOverride(resource, path, override); err != nil {
			return fmt.Errorf("failed to apply override at %s: %w", path, err)
		}
	}
	return nil
}

// applyOverride applies a single override operation
func applyOverride(resource *parser.Resource, path string, override config.OverrideOperation) error {
	// Set the value at the specified path
	if err := resource.SetAttributeValue(path, override.Value); err != nil {
		return fmt.Errorf("failed to set value: %w", err)
	}
	return nil
}
