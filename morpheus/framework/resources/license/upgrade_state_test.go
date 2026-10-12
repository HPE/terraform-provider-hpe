// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package license

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// v0State builds a legacy (SDKv2-shaped, schema version 0) license state with
// the given id value, a non-empty key, and every other attribute null.
func v0State(t *testing.T, id tftypes.Value) tfsdk.State {
	t.Helper()

	s := licenseSchemaV0()
	objType, ok := s.Type().TerraformType(context.Background()).(tftypes.Object)
	if !ok {
		t.Fatalf("v0 schema type is %T, want tftypes.Object", s.Type().TerraformType(context.Background()))
	}

	attrs := make(map[string]tftypes.Value, len(objType.AttributeTypes))
	for name, typ := range objType.AttributeTypes {
		attrs[name] = tftypes.NewValue(typ, nil)
	}
	attrs["id"] = id
	attrs["key"] = tftypes.NewValue(tftypes.String, "SECRETKEY1234567")
	attrs["product_tier"] = tftypes.NewValue(tftypes.String, "enterprise")

	return tfsdk.State{Schema: s, Raw: tftypes.NewValue(objType, attrs)}
}

func ptrState(s tfsdk.State) *tfsdk.State { return &s }

func TestUpgradeLicenseStateV0toV1(t *testing.T) {
	ctx := context.Background()

	cases := []struct {
		name    string
		id      tftypes.Value
		wantID  string
		wantErr string
	}{
		{
			name:   "carries id and details, drops key",
			id:     tftypes.NewValue(tftypes.String, "ABCDEFGH"),
			wantID: "ABCDEFGH",
		},
		{
			name:    "null id is refused",
			id:      tftypes.NewValue(tftypes.String, nil),
			wantErr: "has no id",
		},
		{
			name:    "empty id is refused",
			id:      tftypes.NewValue(tftypes.String, ""),
			wantErr: "has no id",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := resource.UpgradeStateRequest{State: ptrState(v0State(t, tc.id))}
			current := LicenseResourceSchema(ctx)
			resp := resource.UpgradeStateResponse{
				State: tfsdk.State{Schema: current},
			}

			upgradeLicenseStateV0toV1(ctx, req, &resp)

			if tc.wantErr != "" {
				if !resp.Diagnostics.HasError() {
					t.Fatalf("expected an error containing %q, got none", tc.wantErr)
				}
				if detail := resp.Diagnostics.Errors()[0].Detail(); !strings.Contains(detail, tc.wantErr) {
					t.Fatalf("error detail %q does not contain %q", detail, tc.wantErr)
				}

				return
			}

			if resp.Diagnostics.HasError() {
				t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
			}

			var upgraded LicenseModel
			if diags := resp.State.Get(ctx, &upgraded); diags.HasError() {
				t.Fatalf("decoding upgraded state: %v", diags)
			}
			if upgraded.Id.ValueString() != tc.wantID {
				t.Fatalf("upgraded id = %q, want %q", upgraded.Id.ValueString(), tc.wantID)
			}
			if !upgraded.Key.IsNull() {
				t.Fatalf("key must upgrade to null, got %v", upgraded.Key)
			}
			if upgraded.ProductTier.ValueString() != "enterprise" {
				t.Fatalf("product_tier should carry over, got %q", upgraded.ProductTier.ValueString())
			}
		})
	}
}
