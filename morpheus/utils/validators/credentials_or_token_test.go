// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// identityObject builds an identity block value with the given credential
// attributes, defaulting the rest to null.
func identityObject(t *testing.T, set map[string]attr.Value) basetypes.ObjectValue {
	t.Helper()

	attributeTypes := map[string]attr.Type{
		"client_id":        types.StringType,
		"client_secret":    types.StringType,
		"issuer_url":       types.StringType,
		"token_issuer_url": types.StringType,
		"iam_token":        types.StringType,
		"location":         types.StringType,
	}

	values := map[string]attr.Value{}
	for name := range attributeTypes {
		if v, ok := set[name]; ok {
			values[name] = v

			continue
		}

		values[name] = types.StringNull()
	}

	obj, diags := types.ObjectValue(attributeTypes, values)
	if diags.HasError() {
		t.Fatalf("could not build the object value: %v", diags)
	}

	return obj
}

func TestIdentityCredentialsOrTokenValidator(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		value     attr.Value
		wantError bool
	}{
		"no credentials at all": {
			value:     nil, // filled in below with only location set
			wantError: true,
		},
		"client id only": {
			value: identityObject(t, map[string]attr.Value{
				"client_id": types.StringValue("id"),
			}),
		},
		"client secret only": {
			value: identityObject(t, map[string]attr.Value{
				"client_secret": types.StringValue("secret"),
			}),
		},
		"issuer url only": {
			value: identityObject(t, map[string]attr.Value{
				"issuer_url": types.StringValue("https://issuer.example.invalid"),
			}),
		},
		// The Disconnected block's spelling of the same credential.
		"token issuer url only": {
			value: identityObject(t, map[string]attr.Value{
				"token_issuer_url": types.StringValue("https://issuer.example.invalid"),
			}),
		},
		"iam token only": {
			value: identityObject(t, map[string]attr.Value{
				"iam_token": types.StringValue("token"),
			}),
		},
		"full credentials": {
			value: identityObject(t, map[string]attr.Value{
				"client_id":     types.StringValue("id"),
				"client_secret": types.StringValue("secret"),
				"issuer_url":    types.StringValue("https://issuer.example.invalid"),
			}),
		},
		// An unknown credential may still turn out to be set, so the
		// configuration must be left alone until Terraform knows.
		"unknown client id": {
			value: identityObject(t, map[string]attr.Value{
				"client_id": types.StringUnknown(),
			}),
		},
		"unknown iam token": {
			value: identityObject(t, map[string]attr.Value{
				"iam_token": types.StringUnknown(),
			}),
		},
	}

	for name, testcase := range tests {
		tc := testcase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			value := tc.value
			if value == nil {
				value = identityObject(t, map[string]attr.Value{
					"location": types.StringValue("site-a"),
				})
			}

			obj, ok := value.(basetypes.ObjectValue)
			if !ok {
				t.Fatalf("test value is %T, want basetypes.ObjectValue", value)
			}

			resp := &validator.ObjectResponse{}
			IdentityCredentialsOrTokenValidator().ValidateObject(
				context.Background(),
				validator.ObjectRequest{
					Path:        path.Root("pce_identity"),
					ConfigValue: obj,
				},
				resp,
			)

			if got := resp.Diagnostics.HasError(); got != tc.wantError {
				t.Errorf("HasError = %v, want %v (diags: %v)",
					got, tc.wantError, resp.Diagnostics)
			}
		})
	}
}

// The two identity blocks name the issuer URL differently, and neither has the
// other's spelling. The credential list spans both, so it names an attribute
// that any given block does not have: that has to be skipped rather than
// treated as unset, and the block's own spelling has to still be found.
func TestIdentityCredentialsOrTokenValidatorPerBlockIssuerAttribute(t *testing.T) {
	t.Parallel()

	for block, issuerAttr := range map[string]string{
		"pce_identity":              "issuer_url",
		"pce_disconnected_identity": "token_issuer_url",
	} {
		name, attribute := block, issuerAttr

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Only the attributes this block actually has, so the other
			// block's issuer spelling is absent from the object entirely.
			obj, diags := types.ObjectValue(
				map[string]attr.Type{
					"client_id":     types.StringType,
					"client_secret": types.StringType,
					attribute:       types.StringType,
					"iam_token":     types.StringType,
				},
				map[string]attr.Value{
					"client_id":     types.StringNull(),
					"client_secret": types.StringNull(),
					attribute:       types.StringValue("https://issuer.example.invalid"),
					"iam_token":     types.StringNull(),
				},
			)
			if diags.HasError() {
				t.Fatalf("could not build the object value: %v", diags)
			}

			resp := &validator.ObjectResponse{}
			IdentityCredentialsOrTokenValidator().ValidateObject(
				context.Background(),
				validator.ObjectRequest{
					Path:        path.Root(name),
					ConfigValue: obj,
				},
				resp,
			)

			if resp.Diagnostics.HasError() {
				t.Errorf("HasError = true for %s set on %s, want false: %v",
					attribute, name, resp.Diagnostics)
			}
		})
	}
}

// A null or unknown block must not be reported: there is nothing to validate.
func TestIdentityCredentialsOrTokenValidatorIgnoresAbsentBlock(t *testing.T) {
	t.Parallel()

	objectType := map[string]attr.Type{
		"client_id": types.StringType,
		"iam_token": types.StringType,
	}

	for name, value := range map[string]basetypes.ObjectValue{
		"null":    types.ObjectNull(objectType),
		"unknown": types.ObjectUnknown(objectType),
	} {
		obj := value

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			resp := &validator.ObjectResponse{}
			IdentityCredentialsOrTokenValidator().ValidateObject(
				context.Background(),
				validator.ObjectRequest{
					Path:        path.Root("pce_identity"),
					ConfigValue: obj,
				},
				resp,
			)

			if resp.Diagnostics.HasError() {
				t.Errorf("HasError = true for a %s block, want false: %v",
					name, resp.Diagnostics)
			}
		})
	}
}
