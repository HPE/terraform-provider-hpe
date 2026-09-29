// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package tenant

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// v0State builds a legacy (SDKv2-shaped, schema version 0) tenant state with
// the given id value and every other attribute null.
func v0State(t *testing.T, id tftypes.Value) tfsdk.State {
	t.Helper()

	s := tenantSchemaV0()
	objType, ok := s.Type().TerraformType(context.Background()).(tftypes.Object)
	if !ok {
		t.Fatalf("v0 schema type is %T, want tftypes.Object", s.Type().TerraformType(context.Background()))
	}

	attrs := make(map[string]tftypes.Value, len(objType.AttributeTypes))
	for name, typ := range objType.AttributeTypes {
		attrs[name] = tftypes.NewValue(typ, nil)
	}
	attrs["id"] = id

	return tfsdk.State{Schema: s, Raw: tftypes.NewValue(objType, attrs)}
}

func TestUpgradeTenantStateV0toV1(t *testing.T) {
	ctx := context.Background()

	cases := []struct {
		name    string
		id      tftypes.Value
		wantID  int64
		wantErr string
	}{
		{
			name:   "string id converts to int64",
			id:     tftypes.NewValue(tftypes.String, "42"),
			wantID: 42,
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
		{
			name:    "non-numeric id is refused",
			id:      tftypes.NewValue(tftypes.String, "abc"),
			wantErr: "could not convert legacy string id abc",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := resource.UpgradeStateRequest{State: ptrState(v0State(t, tc.id))}
			current := TenantResourceSchema(ctx)
			resp := resource.UpgradeStateResponse{
				State: tfsdk.State{Schema: current},
			}

			upgradeTenantStateV0toV1(ctx, req, &resp)

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

			var upgraded TenantModel
			if diags := resp.State.Get(ctx, &upgraded); diags.HasError() {
				t.Fatalf("decoding upgraded state: %v", diags)
			}
			if upgraded.Id.ValueInt64() != tc.wantID {
				t.Fatalf("upgraded id = %d, want %d", upgraded.Id.ValueInt64(), tc.wantID)
			}
			if !upgraded.ParentId.IsNull() || !upgraded.Master.IsNull() {
				t.Fatalf("framework-only attributes should upgrade to null: parent_id=%v master=%v",
					upgraded.ParentId, upgraded.Master)
			}
		})
	}
}

func ptrState(s tfsdk.State) *tfsdk.State { return &s }
