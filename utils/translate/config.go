// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package translate

import (
	"context"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// ResourceConfig holds the transformation rules for a single resource.
// It uses the same YAML structure as the code-spec pipeline's config.yaml.
type ResourceConfig struct {
	// Templates are named attribute definitions that can be injected via moves.
	Templates map[string]TemplateConfig `yaml:"templates"`
	// Moves are ordered rename/restructure operations (array of {from: to} maps).
	Moves []map[string]string `yaml:"moves"`
	// Removes are paths to strip from the schema/response.
	Removes []string `yaml:"removes"`
	// Overrides modify attribute metadata (not used at runtime for translation).
	Overrides map[string]Override `yaml:"overrides"`
	// Versions holds version-specific adjustments.
	Versions []VersionOverride `yaml:"versions"`
	// Paths defines the API endpoints for each CRUD operation.
	Paths map[string]string `yaml:"paths"`
	// Envelope defines how request/response bodies are wrapped.
	Envelope *EnvelopeConfig `yaml:"envelope"`
	// Discriminator defines how anyOf variants are identified on read.
	Discriminator *DiscriminatorConfig `yaml:"discriminator"`

	// --- Hybrid extensions (runtime-only directives) ---

	// PreserveOnRead lists TF field paths that should always use the plan value
	// on read, because the API doesn't return them (write-only fields).
	PreserveOnRead []string `yaml:"preserve_on_read"`
	// ComputedFields defines fields whose values are derived from a different
	// path in the API response (structural differences).
	ComputedFields map[string]ComputedFieldConfig `yaml:"computed_fields"`
}

// TemplateConfig represents a template attribute definition.
// We only care about the template name and its implied type for runtime
// type conversions (e.g., template-bool implies on/off string conversion).
type TemplateConfig struct {
	Bool    *TemplateField `yaml:"bool"`
	String  *TemplateField `yaml:"string"`
	Int64   *TemplateField `yaml:"int64"`
	Dynamic *TemplateField `yaml:"dynamic"`
}

// TemplateField holds the field definition within a template.
type TemplateField struct {
	Description              string `yaml:"description"`
	ComputedOptionalRequired string `yaml:"computed_optional_required"`
}

// Override holds metadata overrides for a field (mostly for code-gen, not runtime).
type Override struct {
	OneOf                    []any  `yaml:"one_of"`
	Description              string `yaml:"description"`
	ComputedOptionalRequired string `yaml:"computed_optional_required"`
	Sensitive                *bool  `yaml:"sensitive"`
	WriteOnly                *bool  `yaml:"write_only"`
	Default                  any    `yaml:"default"`
	StateForUnknown          bool   `yaml:"state_for_unknown"`
	RequiresReplace          bool   `yaml:"requires_replace"`
}

// VersionOverride holds version-specific transform adjustments.
type VersionOverride struct {
	Constraint  string              `yaml:"constraint"`
	Description string              `yaml:"description"`
	Moves       []map[string]string `yaml:"moves"`
	Removes     []string            `yaml:"removes"`
}

// EnvelopeConfig defines how request/response bodies are wrapped.
type EnvelopeConfig struct {
	Request  string `yaml:"request"`  // Key that wraps the POST/PUT body (e.g., "zone")
	Response string `yaml:"response"` // Key that wraps the GET response (e.g., "zone")
	List     string `yaml:"list"`     // Key for list responses (e.g., "zones")
}

// DiscriminatorConfig defines how anyOf variants are identified on read.
type DiscriminatorConfig struct {
	// Field is the TF schema field name used to determine the variant (e.g., "cloud_type_code").
	Field string `yaml:"field"`
	// Variants maps discriminator values to variant field names.
	// e.g., {"amazon": "config_aws", "standard": "config_hvm"}
	Variants map[string]string `yaml:"variants"`
}

// ComputedFieldConfig defines how a TF field is derived from the API response
// when the response structure differs from the request structure.
type ComputedFieldConfig struct {
	// From is a dot-notation path with optional array index notation.
	// e.g., "groups.0.id" means response.groups[0].id
	From string `yaml:"from"`
	// Type is the expected value type: "string", "int64", "bool", "float64".
	Type string `yaml:"type"`
}

// PostReadHook is called after generic unmarshal to apply resource-specific fixups.
// It receives the raw API response map (already envelope-unwrapped) and can mutate
// the state model directly.
type PostReadHook func(ctx context.Context, raw map[string]any, state any, plan any) error

// PostWriteHook is called after generic marshal to mutate the request body
// before it's sent to the API.
type PostWriteHook func(ctx context.Context, body map[string]any, model any) error

// ParseConfig parses a single resource config.yaml content.
// The YAML is expected to have a single top-level key (the resource name).
func ParseConfig(data []byte) (*ResourceConfig, error) {
	var raw map[string]ResourceConfig
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	if len(raw) != 1 {
		return nil, fmt.Errorf("config must have exactly 1 top-level key, got %d", len(raw))
	}

	for _, cfg := range raw {
		return &cfg, nil
	}

	return nil, fmt.Errorf("empty config")
}

// CompiledConfig is the pre-processed form of ResourceConfig ready for runtime use.
// It pre-computes forward and inverse move mappings.
type CompiledConfig struct {
	raw *ResourceConfig

	// forwardMoves are processed moves for unmarshal (API→TF direction).
	// Each entry is {from, to} representing a rename.
	forwardMoves []movePair
	// inverseMoves are the reverse for marshal (TF→API direction).
	inverseMoves []movePair
	// templateTypes maps TF field paths to their template type (for type conversion).
	templateTypes map[string]templateType
	// removes is the set of paths to filter out.
	removes map[string]bool
	// preserveOnRead is the set of TF field paths that should use plan value on read.
	preserveOnRead map[string]bool

	// Hooks (runtime-only, not from YAML)
	postRead  PostReadHook
	postWrite PostWriteHook
}

type movePair struct {
	from string
	to   string
}

type templateType int

const (
	templateNone    templateType = iota
	templateBool                 // bool <-> "on"/"off" string
	templateDynamic              // dynamic <-> map[string]any
)

// Compile pre-processes a ResourceConfig into a CompiledConfig.
func Compile(cfg *ResourceConfig) *CompiledConfig {
	cc := &CompiledConfig{
		raw:            cfg,
		templateTypes:  make(map[string]templateType),
		removes:        make(map[string]bool),
		preserveOnRead: make(map[string]bool),
	}

	// Process moves
	for _, m := range cfg.Moves {
		for from, to := range m {
			if strings.HasPrefix(from, "template-") {
				// Template insertion: record the type conversion needed
				tplName := strings.TrimPrefix(from, "template-")
				if tpl, ok := cfg.Templates[tplName]; ok {
					if tpl.Bool != nil {
						cc.templateTypes[to] = templateBool
					} else if tpl.Dynamic != nil {
						cc.templateTypes[to] = templateDynamic
					}
				}
				// Template moves don't participate in rename mapping
				continue
			}

			cc.forwardMoves = append(cc.forwardMoves, movePair{from: from, to: to})
			// Inverse: swap from/to for marshal direction
			cc.inverseMoves = append(cc.inverseMoves, movePair{from: to, to: from})
		}
	}

	// Reverse the inverse moves order (must undo in reverse)
	for i, j := 0, len(cc.inverseMoves)-1; i < j; i, j = i+1, j-1 {
		cc.inverseMoves[i], cc.inverseMoves[j] = cc.inverseMoves[j], cc.inverseMoves[i]
	}

	// Process removes
	for _, r := range cfg.Removes {
		cc.removes[r] = true
	}

	// Process preserve_on_read
	for _, p := range cfg.PreserveOnRead {
		cc.preserveOnRead[p] = true
	}

	return cc
}
