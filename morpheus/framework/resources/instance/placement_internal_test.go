// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instance

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"

	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

// tfOrNull builds a tftypes value that is null when p is nil.
func numOrNull(p *int64) tftypes.Value {
	if p == nil {
		return tftypes.NewValue(tftypes.Number, nil)
	}

	return tftypes.NewValue(tftypes.Number, *p)
}

func strOrNull(p *string) tftypes.Value {
	if p == nil {
		return tftypes.NewValue(tftypes.String, nil)
	}

	return tftypes.NewValue(tftypes.String, *p)
}

// placementValidateReq builds a ValidateConfigRequest with a config_hvm block
// carrying the given placement/affinity/kvm values. A nil pointer means null;
// unknown is signalled by the matching *Unknown bool.
func placementValidateReq(
	ctx context.Context,
	t *testing.T,
	strategy *string,
	strategyUnknown bool,
	preferredHost *int64,
	preferredUnknown bool,
	affinityGroup *int64,
	kvmHost *int64,
) (fwresource.ValidateConfigRequest, *fwresource.ValidateConfigResponse) {
	t.Helper()

	s := getTestSchema(ctx)

	hvmType := tftypes.Object{AttributeTypes: map[string]tftypes.Type{
		"affinity_group_id":     tftypes.Number,
		"create_user":           tftypes.Bool,
		"image_id":              tftypes.Number,
		"kvm_host_id":           tftypes.Number,
		"nested_virtualization": tftypes.String,
		"no_agent":              tftypes.Bool,
		"placement_strategy":    tftypes.String,
		"preferred_host_id":     tftypes.Number,
		"resource_pool_id":      tftypes.String,
		"user_data":             tftypes.String,
	}}

	strategyVal := strOrNull(strategy)
	if strategyUnknown {
		strategyVal = tftypes.NewValue(tftypes.String, tftypes.UnknownValue)
	}

	preferredVal := numOrNull(preferredHost)
	if preferredUnknown {
		preferredVal = tftypes.NewValue(tftypes.Number, tftypes.UnknownValue)
	}

	hvmAttrs := map[string]tftypes.Value{
		"affinity_group_id":     numOrNull(affinityGroup),
		"create_user":           tftypes.NewValue(tftypes.Bool, nil),
		"image_id":              tftypes.NewValue(tftypes.Number, nil),
		"kvm_host_id":           numOrNull(kvmHost),
		"nested_virtualization": tftypes.NewValue(tftypes.String, nil),
		"no_agent":              tftypes.NewValue(tftypes.Bool, nil),
		"placement_strategy":    strategyVal,
		"preferred_host_id":     preferredVal,
		"resource_pool_id":      tftypes.NewValue(tftypes.String, nil),
		"user_data":             tftypes.NewValue(tftypes.String, nil),
	}

	hvmVal := tftypes.NewValue(hvmType, hvmAttrs)

	rootType := schemaRootTfType(ctx, s)
	rawAttrs := map[string]tftypes.Value{}
	for name, attr := range s.Attributes {
		if name == "config_hvm" {
			rawAttrs[name] = hvmVal
		} else {
			rawAttrs[name] = tftypes.NewValue(attr.GetType().TerraformType(ctx), nil)
		}
	}

	req := fwresource.ValidateConfigRequest{
		Config: tfsdk.Config{
			Schema: s,
			Raw:    tftypes.NewValue(rootType, rawAttrs),
		},
	}

	return req, &fwresource.ValidateConfigResponse{}
}

func sp(v string) *string { return &v }

func TestUnitHvmPlacementValidator(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	cases := []struct {
		name             string
		strategy         *string
		strategyUnknown  bool
		preferred        *int64
		preferredUnknown bool
		affinity         *int64
		kvmHost          *int64
		wantErr          string // error summary, "" if none
		wantWarn         string // warning summary, "" if none
	}{
		{
			name:      "preferred without strategy errors",
			strategy:  nil,
			preferred: ptr(int64(7)),
			wantErr:   "Missing placement strategy",
		},
		{
			name:      "failover without host errors",
			strategy:  sp("failover"),
			preferred: nil,
			wantErr:   "Missing preferred host for failover",
		},
		{
			name:      "failover with host ok",
			strategy:  sp("failover"),
			preferred: ptr(int64(7)),
		},
		{
			name:     "pinned without host ok",
			strategy: sp("pinned"),
		},
		{
			name:     "pinned with affinity warns",
			strategy: sp("pinned"),
			affinity: ptr(int64(3)),
			wantWarn: "Pin outranks affinity group",
		},
		{
			name:      "kvm_host differs from preferred warns",
			strategy:  sp("pinned"),
			preferred: ptr(int64(9)),
			kvmHost:   ptr(int64(5)),
			wantWarn:  "Initial placement differs from steady-state preference",
		},
		{
			name:      "kvm_host equals preferred no warn",
			strategy:  sp("pinned"),
			preferred: ptr(int64(5)),
			kvmHost:   ptr(int64(5)),
		},
		{
			name:            "unknown strategy skips checks",
			strategyUnknown: true,
			preferred:       ptr(int64(7)),
		},
		{
			name:             "unknown preferred skips checks",
			strategy:         sp("failover"),
			preferredUnknown: true,
		},
		{
			name:     "auto no host ok",
			strategy: sp("auto"),
		},
		{
			name:      "auto with host errors",
			strategy:  sp("auto"),
			preferred: ptr(int64(5)),
			wantErr:   "Preferred host is ignored under auto placement",
		},
		{
			name:      "auto with differing kvm_host does not warn about migration",
			strategy:  sp("auto"),
			preferred: ptr(int64(9)),
			kvmHost:   ptr(int64(5)),
			wantErr:   "Preferred host is ignored under auto placement",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			req, resp := placementValidateReq(ctx, t,
				tc.strategy, tc.strategyUnknown,
				tc.preferred, tc.preferredUnknown,
				tc.affinity, tc.kvmHost)

			hvmPlacementValidator{}.ValidateResource(ctx, req, resp)

			if tc.wantErr == "" && resp.Diagnostics.HasError() {
				t.Fatalf("unexpected error: %v", resp.Diagnostics.Errors())
			}
			if tc.wantErr != "" {
				found := false
				for _, d := range resp.Diagnostics.Errors() {
					if d.Summary() == tc.wantErr {
						found = true
					}
				}
				if !found {
					t.Fatalf("want error %q, got: %v", tc.wantErr, resp.Diagnostics.Errors())
				}
			}
			if tc.wantWarn == "" && resp.Diagnostics.WarningsCount() > 0 {
				t.Fatalf("unexpected warning: %v", resp.Diagnostics.Warnings())
			}
			if tc.wantWarn != "" {
				found := false
				for _, d := range resp.Diagnostics.Warnings() {
					if d.Summary() == tc.wantWarn {
						found = true
					}
				}
				if !found {
					t.Fatalf("want warning %q, got: %v", tc.wantWarn, resp.Diagnostics.Warnings())
				}
			}
		})
	}
}

// hvmModel builds an InstanceModel whose config_hvm carries the given
// placement values (nil = null).
func hvmModel(t *testing.T, strategy *string, preferred *int64) InstanceModel {
	t.Helper()

	ctx := context.Background()
	cfg := NewConfigHvmValueNull()
	attrs := map[string]attr.Value{}
	for k, typ := range cfg.AttributeTypes(ctx) {
		attrs[k] = nullOf(typ)
	}
	attrs["placement_strategy"] = types.StringNull()
	if strategy != nil {
		attrs["placement_strategy"] = types.StringValue(*strategy)
	}
	attrs["preferred_host_id"] = types.Int64Null()
	if preferred != nil {
		attrs["preferred_host_id"] = types.Int64Value(*preferred)
	}

	v, d := NewConfigHvmValue(cfg.AttributeTypes(ctx), attrs)
	if d.HasError() {
		t.Fatal(d)
	}

	return InstanceModel{ConfigHvm: v}
}

func nullOf(typ attr.Type) attr.Value {
	v, _ := typ.ValueFromTerraform(context.Background(), tftypes.NewValue(typ.TerraformType(context.Background()), nil))

	return v
}

func TestUnitPlannedPlacement(t *testing.T) {
	t.Parallel()

	if _, _, want := plannedPlacement(InstanceModel{ConfigHvm: NewConfigHvmValueNull()}); want {
		t.Fatal("null config_hvm must not request placement")
	}
	if _, _, want := plannedPlacement(hvmModel(t, nil, nil)); want {
		t.Fatal("no placement attributes must not request placement")
	}

	strategy, preferred, want := plannedPlacement(hvmModel(t, sp("pinned"), ptr(int64(7))))
	if !want || strategy != "pinned" || preferred == nil || *preferred != 7 {
		t.Fatalf("pinned+7: got %q %v %v", strategy, preferred, want)
	}

	strategy, preferred, want = plannedPlacement(hvmModel(t, sp("auto"), nil))
	if !want || strategy != "auto" || preferred != nil {
		t.Fatalf("auto: got %q %v %v", strategy, preferred, want)
	}
}

func TestUnitPlacementChanged(t *testing.T) {
	t.Parallel()

	pinned7 := hvmModel(t, sp("pinned"), ptr(int64(7)))

	if placementChanged(pinned7, pinned7) {
		t.Fatal("identical placement must not be a change")
	}
	if !placementChanged(hvmModel(t, sp("pinned"), ptr(int64(8))), pinned7) {
		t.Fatal("a different preferred host is a change")
	}
	if !placementChanged(hvmModel(t, sp("auto"), nil), pinned7) {
		t.Fatal("a different strategy is a change")
	}
	if placementChanged(InstanceModel{ConfigHvm: NewConfigHvmValueNull()}, pinned7) {
		t.Fatal("null config_hvm in the plan must not be a change")
	}
}

func TestUnitBuildManageHostPlacementRequest(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name      string
		strategy  string
		preferred *int64
		wantJSON  string
	}{
		{
			name:      "pinned with host",
			strategy:  "pinned",
			preferred: ptr(int64(42)),
			wantJSON:  `{"server":{"placementStrategy":"pinned","preferredParentServer":{"id":42}}}`,
		},
		{
			name:      "failover with host",
			strategy:  "failover",
			preferred: ptr(int64(7)),
			wantJSON:  `{"server":{"placementStrategy":"failover","preferredParentServer":{"id":7}}}`,
		},
		{
			name:      "auto omits preferred",
			strategy:  "auto",
			preferred: ptr(int64(99)),
			wantJSON:  `{"server":{"placementStrategy":"auto"}}`,
		},
		{
			name:      "pinned without host",
			strategy:  "pinned",
			preferred: nil,
			wantJSON:  `{"server":{"placementStrategy":"pinned"}}`,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			req := buildManageHostPlacementRequest(tc.strategy, tc.preferred)
			b, err := json.Marshal(req)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			if string(b) != tc.wantJSON {
				t.Fatalf("JSON mismatch:\n got: %s\nwant: %s", b, tc.wantJSON)
			}
		})
	}
}

func TestUnitReadHostPlacementMapping(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name          string
		server        *sdk.GetHost200ResponseServer
		wantStrategy  types.String
		wantPreferred types.Int64
	}{
		{
			name: "present",
			server: func() *sdk.GetHost200ResponseServer {
				s := &sdk.GetHost200ResponseServer{}
				s.PlacementStrategy = *sdk.NewNullableString(sdk.PtrString("pinned"))
				s.PreferredParentServer = &sdk.AddBaremetalHost200ResponseServerPreferredParentServer{
					Id: ptr(int64(42)),
				}

				return s
			}(),
			wantStrategy:  types.StringValue("pinned"),
			wantPreferred: types.Int64Value(42),
		},
		{
			name: "null placement",
			server: func() *sdk.GetHost200ResponseServer {
				return &sdk.GetHost200ResponseServer{}
			}(),
			wantStrategy:  types.StringNull(),
			wantPreferred: types.Int64Null(),
		},
		{
			name: "strategy set no preferred",
			server: func() *sdk.GetHost200ResponseServer {
				s := &sdk.GetHost200ResponseServer{}
				s.PlacementStrategy = *sdk.NewNullableString(sdk.PtrString("failover"))

				return s
			}(),
			wantStrategy:  types.StringValue("failover"),
			wantPreferred: types.Int64Null(),
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			gotStrategy, gotPreferred := mapServerPlacement(tc.server)

			if !gotStrategy.Equal(tc.wantStrategy) {
				t.Fatalf("strategy: got %v want %v", gotStrategy, tc.wantStrategy)
			}
			if !gotPreferred.Equal(tc.wantPreferred) {
				t.Fatalf("preferred: got %v want %v", gotPreferred, tc.wantPreferred)
			}
		})
	}
}
