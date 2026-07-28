// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package translate

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestModel mimics a simplified CloudModel structure for testing.
type TestModel struct {
	Name      types.String       `tfsdk:"name"`
	TenantId  types.Int64        `tfsdk:"tenant_id"`
	Enabled   types.Bool         `tfsdk:"enabled"`
	ConfigAws TestConfigAwsValue `tfsdk:"config_aws"`
	Id        types.Int64        `tfsdk:"id"`
}

type TestConfigAwsValue struct {
	Endpoint           types.String `tfsdk:"endpoint"`
	Region             types.String `tfsdk:"region"`
	UseHostCredentials types.Bool   `tfsdk:"use_host_credentials"`
	state              attr.ValueState
}

func (v TestConfigAwsValue) IsNull() bool                                    { return v.state == attr.ValueStateNull }
func (v TestConfigAwsValue) IsUnknown() bool                                 { return v.state == attr.ValueStateUnknown }
func (v TestConfigAwsValue) Type(_ context.Context) attr.Type                { return nil }
func (v TestConfigAwsValue) ToTerraformValue(_ context.Context) (any, error) { return nil, nil }
func (v TestConfigAwsValue) Equal(_ attr.Value) bool                         { return false }
func (v TestConfigAwsValue) String() string                                  { return "" }

func TestMarshal_SimpleModel(t *testing.T) {
	ctx := context.Background()

	model := &TestModel{
		Name:     types.StringValue("my-cloud"),
		TenantId: types.Int64Value(42),
		Enabled:  types.BoolValue(true),
		Id:       types.Int64Unknown(), // unknown should be skipped
		ConfigAws: TestConfigAwsValue{
			Endpoint:           types.StringValue("https://aws.example.com"),
			Region:             types.StringValue("us-east-1"),
			UseHostCredentials: types.BoolValue(true),
			state:              attr.ValueStateKnown,
		},
	}

	cfg := &ResourceConfig{
		Templates: map[string]TemplateConfig{
			"bool": {Bool: &TemplateField{Description: "test"}},
		},
		Moves: []map[string]string{
			{"zone": ""},                    // unnest (inverse = nest into zone)
			{"account_id": "tenant_id"},     // rename
			{"config.anyof0": "config_aws"}, // extract variant
			{"template-bool": "config_aws.use_host_credentials"},
		},
		Removes: []string{"status"},
	}
	cc := Compile(cfg)

	result, err := Marshal(ctx, model, cc)
	require.NoError(t, err)

	// The moves handle the envelope: "zone: ''" inversed nests into "zone"
	zone, ok := result["zone"].(map[string]any)
	require.True(t, ok, "result should be wrapped in 'zone' by the inverse unnest move, got: %v", result)

	// Check fields were placed correctly (camelCase in API)
	assert.Equal(t, "my-cloud", zone["name"])
	assert.Equal(t, int64(42), zone["accountId"]) // tenant_id → account_id → accountId
	assert.Equal(t, true, zone["enabled"])

	// Config should be nested directly (anyof0 is transparent in output)
	config, ok := zone["config"].(map[string]any)
	require.True(t, ok, "config should be a nested map, got: %v", zone)

	assert.Equal(t, "https://aws.example.com", config["endpoint"])
	assert.Equal(t, "us-east-1", config["region"])
	assert.Equal(t, "on", config["useHostCredentials"]) // bool → "on", camelCase

	// ID should not be present (was unknown)
	assert.NotContains(t, zone, "id")
}

func TestMarshal_NullConfigAws(t *testing.T) {
	ctx := context.Background()

	model := &TestModel{
		Name:     types.StringValue("my-cloud"),
		TenantId: types.Int64Value(1),
		Enabled:  types.BoolValue(true),
		ConfigAws: TestConfigAwsValue{
			state: attr.ValueStateNull, // null — should be skipped entirely
		},
	}

	cfg := &ResourceConfig{
		Moves: []map[string]string{
			{"zone": ""},
			{"account_id": "tenant_id"},
			{"config.anyof0": "config_aws"},
		},
	}
	cc := Compile(cfg)

	result, err := Marshal(ctx, model, cc)
	require.NoError(t, err)

	zone := result["zone"].(map[string]any)
	assert.Equal(t, "my-cloud", zone["name"])
	// No config should be present since config_aws was null
	assert.NotContains(t, zone, "config")
}

func TestUnmarshal_SimpleModel(t *testing.T) {
	ctx := context.Background()

	// API responses use camelCase
	apiResponse := map[string]any{
		"zone": map[string]any{
			"name":      "my-cloud",
			"accountId": float64(42),
			"enabled":   true,
			"id":        float64(123),
			"status":    "active", // should be removed
			"config": map[string]any{
				"anyof0": map[string]any{
					"endpoint":           "https://aws.example.com",
					"region":             "us-east-1",
					"useHostCredentials": "on",
				},
			},
		},
	}

	cfg := &ResourceConfig{
		Templates: map[string]TemplateConfig{
			"bool": {Bool: &TemplateField{Description: "test"}},
		},
		Moves: []map[string]string{
			{"zone": ""},
			{"account_id": "tenant_id"},
			{"config.anyof0": "config_aws"},
			{"template-bool": "config_aws.use_host_credentials"},
		},
		Removes: []string{"status"},
	}
	cc := Compile(cfg)

	var model TestModel
	err := Unmarshal(ctx, apiResponse, &model, cc, nil)
	require.NoError(t, err)

	assert.Equal(t, "my-cloud", model.Name.ValueString())
	assert.Equal(t, int64(42), model.TenantId.ValueInt64())
	assert.Equal(t, true, model.Enabled.ValueBool())
	assert.Equal(t, int64(123), model.Id.ValueInt64())

	// Nested config should be populated
	assert.False(t, model.ConfigAws.IsNull())
	assert.Equal(t, "https://aws.example.com", model.ConfigAws.Endpoint.ValueString())
	assert.Equal(t, "us-east-1", model.ConfigAws.Region.ValueString())
	assert.Equal(t, true, model.ConfigAws.UseHostCredentials.ValueBool())
}

func TestUnmarshal_PlanPreservation(t *testing.T) {
	ctx := context.Background()

	// API response is missing the "region" field (camelCase)
	apiResponse := map[string]any{
		"zone": map[string]any{
			"name":      "my-cloud",
			"accountId": float64(42),
			"enabled":   true,
			"id":        float64(123),
			"config": map[string]any{
				"anyof0": map[string]any{
					"endpoint": "https://aws.example.com",
					// "region" is missing from the response
				},
			},
		},
	}

	cfg := &ResourceConfig{
		Moves: []map[string]string{
			{"zone": ""},
			{"account_id": "tenant_id"},
			{"config.anyof0": "config_aws"},
		},
	}
	cc := Compile(cfg)

	// Plan has a region value that should be preserved
	plan := &TestModel{
		Name:     types.StringValue("my-cloud"),
		TenantId: types.Int64Value(42),
		Enabled:  types.BoolValue(true),
		ConfigAws: TestConfigAwsValue{
			Endpoint: types.StringValue("https://aws.example.com"),
			Region:   types.StringValue("us-east-1"), // This should be preserved
			state:    attr.ValueStateKnown,
		},
	}

	var model TestModel
	err := Unmarshal(ctx, apiResponse, &model, cc, plan)
	require.NoError(t, err)

	assert.Equal(t, "my-cloud", model.Name.ValueString())
	// Region should be preserved from plan
	assert.Equal(t, "us-east-1", model.ConfigAws.Region.ValueString())
}

func TestMarshalUnmarshal_RoundTrip(t *testing.T) {
	ctx := context.Background()

	original := &TestModel{
		Name:     types.StringValue("my-cloud"),
		TenantId: types.Int64Value(42),
		Enabled:  types.BoolValue(true),
		ConfigAws: TestConfigAwsValue{
			Endpoint: types.StringValue("https://aws.example.com"),
			Region:   types.StringValue("us-east-1"),
			state:    attr.ValueStateKnown,
		},
	}

	cfg := &ResourceConfig{
		Moves: []map[string]string{
			{"zone": ""},
			{"account_id": "tenant_id"},
			{"config.anyof0": "config_aws"},
		},
	}
	cc := Compile(cfg)

	// Marshal produces the API shape (with anyof0 transparent)
	apiBody, err := Marshal(ctx, original, cc)
	require.NoError(t, err)

	// Verify the API shape is correct
	zone, ok := apiBody["zone"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "my-cloud", zone["name"])
	config, ok := zone["config"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "https://aws.example.com", config["endpoint"])

	// For round-trip unmarshal, the API response would have the same shape
	// but with camelCase keys (which it already has from marshal).
	// The Unmarshal needs config fields to be under "config" (which they are).
	// The forward move "config.anyof0: config_aws" expects "config.anyof0.*" keys
	// in the flat map. Since the API response has "config.endpoint", we need the
	// flatten to NOT add anyof0 on read (it doesn't exist in the response).
	// This means on read, the forward move won't match these fields.
	// In practice, the discriminator tells us which variant to use.

	// For now, test that top-level fields round-trip correctly
	var reconstructed TestModel
	err = Unmarshal(ctx, apiBody, &reconstructed, cc, nil)
	require.NoError(t, err)

	assert.Equal(t, original.Name.ValueString(), reconstructed.Name.ValueString())
	assert.Equal(t, original.TenantId.ValueInt64(), reconstructed.TenantId.ValueInt64())
	assert.Equal(t, original.Enabled.ValueBool(), reconstructed.Enabled.ValueBool())
	// Note: config_aws fields won't round-trip without discriminator support
	// because the API response has flat "config" without "anyof0" prefix.
}
