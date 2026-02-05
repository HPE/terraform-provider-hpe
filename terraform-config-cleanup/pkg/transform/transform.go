package transform

import (
	"fmt"

	"github.com/HPE/terraform-config-cleanup/pkg/config"
	"github.com/HPE/terraform-config-cleanup/pkg/parser"
)

// ApplyTransformations applies all transformations from a cleanup config to a resource
func ApplyTransformations(resource *parser.Resource, cfg config.CleanupConfig) error {
	// Apply transformations in order: moves, removes, overrides, then cleanup

	// 1. Apply moves first (these might expose or reorganize data)
	if len(cfg.Moves) > 0 {
		if err := ApplyMoves(resource, cfg.Moves); err != nil {
			return fmt.Errorf("failed to apply moves: %w", err)
		}
	}

	// 2. Apply removes (clean up unwanted attributes)
	if len(cfg.Removes) > 0 {
		if err := ApplyRemoves(resource, cfg.Removes); err != nil {
			return fmt.Errorf("failed to apply removes: %w", err)
		}
	}

	// 3. Apply overrides (set specific values)
	if len(cfg.Overrides) > 0 {
		if err := ApplyOverrides(resource, cfg.Overrides); err != nil {
			return fmt.Errorf("failed to apply overrides: %w", err)
		}
	}

	// 4. Remove null values if configured
	if cfg.RemoveNullValues {
		if err := RemoveNullValues(resource); err != nil {
			return fmt.Errorf("failed to remove null values: %w", err)
		}
	}

	// 5. Remove empty blocks if configured
	if cfg.RemoveEmptyBlocks {
		if err := RemoveEmptyBlocks(resource); err != nil {
			return fmt.Errorf("failed to remove empty blocks: %w", err)
		}
	}

	return nil
}

// TransformFile applies cleanup configurations to all resources in a Terraform file
func TransformFile(tfFile *parser.TerraformFile, configs map[string]config.CleanupConfig) error {
	for _, resource := range tfFile.Resources {
		// Check if we have a config for this resource type
		cfg, exists := configs[resource.Type]
		if !exists {
			// No cleanup config for this resource type, skip
			continue
		}

		if err := ApplyTransformations(resource, cfg); err != nil {
			return fmt.Errorf("failed to transform resource %s.%s: %w", resource.Type, resource.Name, err)
		}
	}

	return nil
}
