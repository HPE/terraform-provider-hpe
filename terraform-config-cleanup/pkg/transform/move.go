package transform

import (
	"fmt"

	"github.com/HPE/terraform-config-cleanup/pkg/parser"
)

// ApplyMoves applies move operations to a resource
// Moves are specified as source: destination pairs
// Example: config.noAgent.Bool: config.noAgent (extracts Bool from union type)
func ApplyMoves(resource *parser.Resource, moves []interface{}) error {
	for _, moveOp := range moves {
		// Parse the move operation (it's a map with one entry: source -> destination)
		moveMap, ok := moveOp.(map[string]interface{})
		if !ok {
			return fmt.Errorf("invalid move operation format: expected map")
		}

		// Each move is a single key-value pair
		for source, destValue := range moveMap {
			destination, ok := destValue.(string)
			if !ok {
				return fmt.Errorf("invalid move destination: expected string")
			}

			if err := applyMove(resource, source, destination); err != nil {
				return fmt.Errorf("failed to apply move %s -> %s: %w", source, destination, err)
			}
		}
	}
	return nil
}

// applyMove applies a single move operation
func applyMove(resource *parser.Resource, source, destination string) error {
	// Get the value from source path
	value, exists := resource.GetAttributeValue(source)
	if !exists {
		// Source doesn't exist - this is okay, might not be present in all cases
		return nil
	}

	// Set the value at destination path
	if err := resource.SetAttributeValue(destination, value); err != nil {
		return fmt.Errorf("failed to set destination: %w", err)
	}

	// Remove the source path (we've moved it)
	if err := resource.RemoveAttributeValue(source); err != nil {
		return fmt.Errorf("failed to remove source: %w", err)
	}

	return nil
}
