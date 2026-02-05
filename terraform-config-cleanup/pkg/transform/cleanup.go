package transform

import (
	"github.com/HPE/terraform-config-cleanup/pkg/parser"
)

// RemoveNullValues removes all attributes with null values from the resource
func RemoveNullValues(resource *parser.Resource) error {
	var toRemove []string

	// Find all null attributes at the top level
	for path, attr := range resource.Attributes {
		if attr.Value == nil {
			toRemove = append(toRemove, path)
			continue
		}

		// Check nested attributes in maps
		if m, ok := attr.Value.(map[string]interface{}); ok {
			removeNullFromMap(m)
		}
	}

	// Remove null attributes
	for _, path := range toRemove {
		resource.RemoveAttribute(path)
	}

	return nil
}

// removeNullFromMap recursively removes null values from a map
func removeNullFromMap(m map[string]interface{}) {
	var toRemove []string

	for key, value := range m {
		if value == nil {
			toRemove = append(toRemove, key)
			continue
		}

		// Recursively handle nested maps
		if nested, ok := value.(map[string]interface{}); ok {
			removeNullFromMap(nested)
		}
	}

	// Remove null keys
	for _, key := range toRemove {
		delete(m, key)
	}
}

// RemoveEmptyBlocks removes blocks that are empty after cleanup
func RemoveEmptyBlocks(resource *parser.Resource) error {
	var toRemove []string

	for path, attr := range resource.Attributes {
		if m, ok := attr.Value.(map[string]interface{}); ok {
			if len(m) == 0 {
				toRemove = append(toRemove, path)
			}
		}
	}

	// Remove empty blocks
	for _, path := range toRemove {
		resource.RemoveAttribute(path)
	}

	return nil
}
