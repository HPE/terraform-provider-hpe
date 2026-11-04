// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package settingappliance

import (
	"context"
	"math/big"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// testConfig is a configuration of the resource that sets the given attributes
// and leaves every other one null.
func testConfig(ctx context.Context, set map[string]tftypes.Value) tfsdk.Config {
	s := SettingApplianceResourceSchema(ctx)
	objectType := s.Type().TerraformType(ctx).(tftypes.Object)

	values := make(map[string]tftypes.Value, len(objectType.AttributeTypes))
	for name, typ := range objectType.AttributeTypes {
		values[name] = tftypes.NewValue(typ, nil)
	}

	for name, v := range set {
		values[name] = v
	}

	return tfsdk.Config{Schema: s, Raw: tftypes.NewValue(objectType, values)}
}

// validateConfig runs the resource's config validators on a configuration.
func validateConfig(ctx context.Context, config tfsdk.Config) diag.Diagnostics {
	req := resource.ValidateConfigRequest{Config: config}
	resp := &resource.ValidateConfigResponse{}

	for _, v := range (&Resource{}).ConfigValidators(ctx) {
		v.ValidateResource(ctx, req, resp)
	}

	return resp.Diagnostics
}

// The registration flag and the two default role ids are configured together,
// as the SDKv2 resource required (each was RequiredWith the other two).
func TestUnitRegistrationAndDefaultRolesRequiredTogether(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	flag := func(v bool) tftypes.Value { return tftypes.NewValue(tftypes.Bool, v) }
	id := func(v int64) tftypes.Value { return tftypes.NewValue(tftypes.Number, big.NewFloat(float64(v))) }
	unknownID := tftypes.NewValue(tftypes.Number, tftypes.UnknownValue)

	cases := []struct {
		name    string
		set     map[string]tftypes.Value
		wantErr bool
	}{
		{"none", nil, false},
		{"all three", map[string]tftypes.Value{
			"registration_enabled": flag(true), "default_role_id": id(5), "default_user_role_id": id(10),
		}, false},
		{"registration off with both roles", map[string]tftypes.Value{
			"registration_enabled": flag(false), "default_role_id": id(5), "default_user_role_id": id(10),
		}, false},
		{"registration alone", map[string]tftypes.Value{
			"registration_enabled": flag(true),
		}, true},
		{"registration off alone", map[string]tftypes.Value{
			"registration_enabled": flag(false),
		}, true},
		{"one role alone", map[string]tftypes.Value{
			"default_role_id": id(5),
		}, true},
		{"both roles without registration", map[string]tftypes.Value{
			"default_role_id": id(5), "default_user_role_id": id(10),
		}, true},
		// A role created in the same apply: decided once it is known.
		{"a role not known yet", map[string]tftypes.Value{
			"registration_enabled": flag(true), "default_role_id": unknownID, "default_user_role_id": id(10),
		}, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			diags := validateConfig(ctx, testConfig(ctx, tc.set))
			if got := diags.HasError(); got != tc.wantErr {
				t.Errorf("error: got %v, want %v (%v)", got, tc.wantErr, diags)
			}
		})
	}
}

// An empty smtp_from_address is refused when planning: Morpheus reports a
// default address whenever none is set, so the apply could not read the empty
// one back.
func TestUnitEmptySMTPFromAddressRejected(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	cases := []struct {
		name    string
		value   tftypes.Value
		wantErr bool
	}{
		{"empty", tftypes.NewValue(tftypes.String, ""), true},
		{"an address", tftypes.NewValue(tftypes.String, "ops@example.com"), false},
		{"not set", tftypes.NewValue(tftypes.String, nil), false},
		{"not known yet", tftypes.NewValue(tftypes.String, tftypes.UnknownValue), false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			diags := validateConfig(ctx, testConfig(ctx, map[string]tftypes.Value{"smtp_from_address": tc.value}))
			if got := diags.HasError(); got != tc.wantErr {
				t.Errorf("error: got %v, want %v (%v)", got, tc.wantErr, diags)
			}
		})
	}
}

// stats_retainment_period takes the values the appliance accepts, 30, 60 or 90
// days; anything else fails when planning instead of at the appliance.
func TestUnitStatsRetainmentPeriodValues(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	attr, ok := SettingApplianceResourceSchema(ctx).Attributes["stats_retainment_period"].(schema.Int64Attribute)
	if !ok {
		t.Fatal("stats_retainment_period is not an integer attribute")
	}

	for _, tc := range []struct {
		days    int64
		wantErr bool
	}{{30, false}, {60, false}, {90, false}, {45, true}, {0, true}} {
		req := validator.Int64Request{Path: path.Root("stats_retainment_period"), ConfigValue: types.Int64Value(tc.days)}
		resp := &validator.Int64Response{}

		for _, v := range attr.Validators {
			v.ValidateInt64(ctx, req, resp)
		}

		if got := resp.Diagnostics.HasError(); got != tc.wantErr {
			t.Errorf("%d days: error %v, want %v", tc.days, got, tc.wantErr)
		}
	}
}
