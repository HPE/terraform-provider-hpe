package config

// CleanupConfig defines transformations to apply to generated Terraform config
type CleanupConfig struct {
	// Resource type this config applies to (e.g., "hpe_morpheus_instance")
	ResourceType string `yaml:"-"`

	// Moves define operations to move/extract values (like go-split moves)
	// Simple syntax: - source: destination
	Moves []interface{} `yaml:"moves"`

	// Removes define attributes/blocks to remove (like go-split removes)
	Removes []interface{} `yaml:"removes"`

	// Overrides define value transformations (like go-split overrides)
	Overrides map[string]OverrideOperation `yaml:"overrides"`

	// RemoveNullValues removes all attributes with null values
	RemoveNullValues bool `yaml:"remove_null_values"`

	// RemoveEmptyBlocks removes blocks that are empty after cleanup
	RemoveEmptyBlocks bool `yaml:"remove_empty_blocks"`
}

// OverrideOperation sets a specific value at a path
type OverrideOperation struct {
	// Value to set
	Value interface{} `yaml:"value"`
}
