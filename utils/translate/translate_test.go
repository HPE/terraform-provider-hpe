// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package translate

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseConfig(t *testing.T) {
	yaml := []byte(`
cloud:
  templates:
    bool:
      bool:
        description: A boolean field
        computed_optional_required: optional
  moves:
    - zone: ""
    - account_id: tenant_id
    - config.anyof0: config_aws
    - template-bool: config_aws.use_host_credentials
  removes:
    - status
    - date_created
  envelope:
    request: zone
    response: zone
  discriminator:
    field: cloud_type_code
    variants:
      amazon: config_aws
      standard: config_hvm
  paths:
    create: POST /api/zones
    read: GET /api/zones/{id}
`)

	cfg, err := ParseConfig(yaml)
	require.NoError(t, err)

	assert.Len(t, cfg.Moves, 4)
	assert.Len(t, cfg.Removes, 2)
	assert.Contains(t, cfg.Templates, "bool")
	assert.Equal(t, "zone", cfg.Envelope.Request)
	assert.Equal(t, "zone", cfg.Envelope.Response)
	assert.Equal(t, "cloud_type_code", cfg.Discriminator.Field)
	assert.Equal(t, "config_aws", cfg.Discriminator.Variants["amazon"])
}

func TestFlatten(t *testing.T) {
	nested := map[string]any{
		"zone": map[string]any{
			"name": "my-cloud",
			"config": map[string]any{
				"apiUrl":     "https://vcenter/sdk",
				"datacenter": "DC1",
			},
			"enabled": true,
		},
	}

	flat := Flatten(nested)

	assert.Equal(t, "my-cloud", flat["zone.name"])
	assert.Equal(t, "https://vcenter/sdk", flat["zone.config.apiUrl"])
	assert.Equal(t, "DC1", flat["zone.config.datacenter"])
	assert.Equal(t, true, flat["zone.enabled"])
	assert.Len(t, flat, 4)
}

func TestUnflatten(t *testing.T) {
	flat := map[string]any{
		"zone.name":              "my-cloud",
		"zone.config.apiUrl":     "https://vcenter/sdk",
		"zone.config.datacenter": "DC1",
		"zone.enabled":           true,
	}

	nested := Unflatten(flat)

	zone := nested["zone"].(map[string]any)
	assert.Equal(t, "my-cloud", zone["name"])
	assert.Equal(t, true, zone["enabled"])

	config := zone["config"].(map[string]any)
	assert.Equal(t, "https://vcenter/sdk", config["apiUrl"])
	assert.Equal(t, "DC1", config["datacenter"])
}

func TestFlattenUnflattenRoundTrip(t *testing.T) {
	original := map[string]any{
		"zone": map[string]any{
			"name": "test",
			"config": map[string]any{
				"field1": "value1",
				"field2": float64(42),
			},
		},
	}

	flat := Flatten(original)
	reconstructed := Unflatten(flat)

	zone := reconstructed["zone"].(map[string]any)
	assert.Equal(t, "test", zone["name"])
	config := zone["config"].(map[string]any)
	assert.Equal(t, "value1", config["field1"])
	assert.Equal(t, float64(42), config["field2"])
}

func TestTransformForRead_UnnestMove(t *testing.T) {
	cfg := &ResourceConfig{
		Moves: []map[string]string{
			{"zone": ""},
		},
	}
	cc := Compile(cfg)

	input := map[string]any{
		"zone.name":    "my-cloud",
		"zone.enabled": true,
		"zone.config":  "value",
	}

	result := cc.TransformForRead(input)

	assert.Equal(t, "my-cloud", result["name"])
	assert.Equal(t, true, result["enabled"])
	assert.Equal(t, "value", result["config"])
	assert.NotContains(t, result, "zone.name")
}

func TestTransformForRead_Rename(t *testing.T) {
	cfg := &ResourceConfig{
		Moves: []map[string]string{
			{"account_id": "tenant_id"},
		},
	}
	cc := Compile(cfg)

	input := map[string]any{
		"account_id": float64(42),
		"name":       "test",
	}

	result := cc.TransformForRead(input)

	assert.Equal(t, float64(42), result["tenant_id"])
	assert.Equal(t, "test", result["name"])
	assert.NotContains(t, result, "account_id")
}

func TestTransformForRead_AnyOfExtract(t *testing.T) {
	cfg := &ResourceConfig{
		Moves: []map[string]string{
			{"config.anyof0": "config_aws"},
		},
	}
	cc := Compile(cfg)

	input := map[string]any{
		"config.anyof0.endpoint": "https://aws.example.com",
		"config.anyof0.region":   "us-east-1",
		"name":                   "test",
	}

	result := cc.TransformForRead(input)

	assert.Equal(t, "https://aws.example.com", result["config_aws.endpoint"])
	assert.Equal(t, "us-east-1", result["config_aws.region"])
	assert.Equal(t, "test", result["name"])
}

func TestTransformForRead_Removes(t *testing.T) {
	cfg := &ResourceConfig{
		Removes: []string{"status", "date_created"},
	}
	cc := Compile(cfg)

	input := map[string]any{
		"name":         "test",
		"status":       "active",
		"date_created": "2024-01-01",
		"enabled":      true,
	}

	result := cc.TransformForRead(input)

	assert.Equal(t, "test", result["name"])
	assert.Equal(t, true, result["enabled"])
	assert.NotContains(t, result, "status")
	assert.NotContains(t, result, "date_created")
}

func TestTransformForWrite_InverseUnnest(t *testing.T) {
	cfg := &ResourceConfig{
		Moves: []map[string]string{
			{"zone": ""},
		},
	}
	cc := Compile(cfg)

	input := map[string]any{
		"name":    "my-cloud",
		"enabled": true,
	}

	result := cc.TransformForWrite(input)

	assert.Equal(t, "my-cloud", result["zone.name"])
	assert.Equal(t, true, result["zone.enabled"])
	assert.NotContains(t, result, "name")
}

func TestTransformForWrite_InverseRename(t *testing.T) {
	cfg := &ResourceConfig{
		Moves: []map[string]string{
			{"account_id": "tenant_id"},
		},
	}
	cc := Compile(cfg)

	input := map[string]any{
		"tenant_id": float64(42),
		"name":      "test",
	}

	result := cc.TransformForWrite(input)

	assert.Equal(t, float64(42), result["account_id"])
	assert.Equal(t, "test", result["name"])
	assert.NotContains(t, result, "tenant_id")
}

func TestTransformForWrite_InverseAnyOfMerge(t *testing.T) {
	cfg := &ResourceConfig{
		Moves: []map[string]string{
			{"config.anyof0": "config_aws"},
		},
	}
	cc := Compile(cfg)

	input := map[string]any{
		"config_aws.endpoint": "https://aws.example.com",
		"config_aws.region":   "us-east-1",
		"name":                "test",
	}

	result := cc.TransformForWrite(input)

	assert.Equal(t, "https://aws.example.com", result["config.anyof0.endpoint"])
	assert.Equal(t, "us-east-1", result["config.anyof0.region"])
	assert.NotContains(t, result, "config_aws.endpoint")
}

func TestTypeConversionsForWrite_BoolToOnOff(t *testing.T) {
	cfg := &ResourceConfig{
		Templates: map[string]TemplateConfig{
			"bool": {Bool: &TemplateField{Description: "test"}},
		},
		Moves: []map[string]string{
			{"template-bool": "config_aws.use_host_credentials"},
			{"template-bool": "config_aws.bypass_proxy"},
		},
	}
	cc := Compile(cfg)

	flat := map[string]any{
		"config_aws.use_host_credentials": true,
		"config_aws.bypass_proxy":         false,
	}

	cc.ApplyTypeConversionsForWrite(flat)

	assert.Equal(t, "on", flat["config_aws.use_host_credentials"])
	assert.Equal(t, "off", flat["config_aws.bypass_proxy"])
}

func TestTypeConversionsForRead_OnOffToBool(t *testing.T) {
	cfg := &ResourceConfig{
		Templates: map[string]TemplateConfig{
			"bool": {Bool: &TemplateField{Description: "test"}},
		},
		Moves: []map[string]string{
			{"template-bool": "config_aws.use_host_credentials"},
		},
	}
	cc := Compile(cfg)

	flat := map[string]any{
		"config_aws.use_host_credentials": "on",
	}

	cc.ApplyTypeConversionsForRead(flat)

	assert.Equal(t, true, flat["config_aws.use_host_credentials"])
}

func TestCompileConfig_FullCloudPattern(t *testing.T) {
	// Simulate the cloud config pattern
	cfg := &ResourceConfig{
		Templates: map[string]TemplateConfig{
			"bool": {Bool: &TemplateField{Description: "test"}},
		},
		Moves: []map[string]string{
			{"zone": ""},                // unnest zone wrapper
			{"account_id": "tenant_id"}, // rename
			{"agent_mode": "agent_install_mode"},
			{"config.anyof0": "config_aws"},
			{"template-bool": "config_aws.use_host_credentials"},
		},
		Removes: []string{
			"status",
			"date_created",
			"uuid",
		},
		Envelope: &EnvelopeConfig{
			Request:  "zone",
			Response: "zone",
		},
	}

	cc := Compile(cfg)

	// Forward (read): API response → TF model
	apiFlat := map[string]any{
		"zone.name":                   "my-cloud",
		"zone.account_id":             float64(1),
		"zone.agent_mode":             "cloudInit",
		"zone.config.anyof0.endpoint": "https://aws.example.com",
		"zone.status":                 "active",
		"zone.date_created":           "2024-01-01",
		"zone.uuid":                   "abc-123",
	}

	tfFlat := cc.TransformForRead(apiFlat)

	assert.Equal(t, "my-cloud", tfFlat["name"])
	assert.Equal(t, float64(1), tfFlat["tenant_id"])
	assert.Equal(t, "cloudInit", tfFlat["agent_install_mode"])
	assert.Equal(t, "https://aws.example.com", tfFlat["config_aws.endpoint"])
	assert.NotContains(t, tfFlat, "status")
	assert.NotContains(t, tfFlat, "date_created")
	assert.NotContains(t, tfFlat, "uuid")

	// Inverse (write): TF model → API request
	tfInput := map[string]any{
		"name":                            "my-cloud",
		"tenant_id":                       float64(1),
		"agent_install_mode":              "cloudInit",
		"config_aws.endpoint":             "https://aws.example.com",
		"config_aws.use_host_credentials": true,
	}

	// Apply type conversions first
	cc.ApplyTypeConversionsForWrite(tfInput)
	assert.Equal(t, "on", tfInput["config_aws.use_host_credentials"])

	apiResult := cc.TransformForWrite(tfInput)

	assert.Equal(t, "my-cloud", apiResult["zone.name"])
	assert.Equal(t, float64(1), apiResult["zone.account_id"])
	assert.Equal(t, "cloudInit", apiResult["zone.agent_mode"])
	assert.Equal(t, "https://aws.example.com", apiResult["zone.config.anyof0.endpoint"])
	assert.Equal(t, "on", apiResult["zone.config.anyof0.use_host_credentials"])
}

func TestVersion_NoOverrides(t *testing.T) {
	cfg := &ResourceConfig{
		Moves: []map[string]string{
			{"zone": ""},
		},
		Removes: []string{"status"},
	}
	cc := Compile(cfg)

	resolved := cc.ResolveForVersion("8.0.0")
	assert.Equal(t, cc, resolved) // Same object when no version overrides
}

func TestVersion_WithMatchingOverride(t *testing.T) {
	cfg := &ResourceConfig{
		Moves: []map[string]string{
			{"zone": ""},
		},
		Removes: []string{"status"},
		Versions: []VersionOverride{
			{
				Constraint: ">= 8.1.0",
				Removes:    []string{"legacy_field"},
			},
		},
	}
	cc := Compile(cfg)

	resolved := cc.ResolveForVersion("8.2.0")
	assert.True(t, resolved.removes["status"])
	assert.True(t, resolved.removes["legacy_field"])

	// Version below threshold should not get the override
	base := cc.ResolveForVersion("8.0.0")
	assert.True(t, base.removes["status"])
	assert.False(t, base.removes["legacy_field"])
}
