// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package network

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/networkconfig"
)

// TestUnitClearableClassificationIsComplete holds clearableAttributes and
// notClearable against the schema: every optional attribute is in exactly one
// of them, every name exists, every exclusion has a reason, and nothing that
// cannot be removed is declared clearable. A new optional attribute fails
// this test until it is classified.
func TestUnitClearableClassificationIsComplete(t *testing.T) {
	t.Parallel()

	problems := clearable.Check(NetworkResourceSchema(context.Background()).Attributes, clearableAttributes, notClearable)
	if len(problems) != 0 {
		t.Fatalf("clearable classification:\n  %s", strings.Join(problems, "\n  "))
	}
}

// networkObject builds a value of the resource's object type with the named
// attributes set and every other one null.
func networkObject(ctx context.Context, t *testing.T, set map[string]tftypes.Value) tftypes.Value {
	t.Helper()

	objType, ok := NetworkResourceSchema(ctx).Type().TerraformType(ctx).(tftypes.Object)
	if !ok {
		t.Fatal("schema type is not an object")
	}

	vals := make(map[string]tftypes.Value, len(objType.AttributeTypes))

	for name, typ := range objType.AttributeTypes {
		if v, ok := set[name]; ok {
			vals[name] = v

			continue
		}

		vals[name] = tftypes.NewValue(typ, nil)
	}

	return tftypes.NewValue(objType, vals)
}

// TestUnitModifyPlanClearOfImmutableCidrIsAReplacement covers the rule for an
// attribute that is Optional but not Computed, which needs no private-state
// record: cidr removed from the configuration is planned null by Terraform,
// and on a type whose cidrEditable flag is off Morpheus would accept and
// ignore the update, so the plan is a replacement — the same rule as for a
// changed value. On a type that edits the CIDR in place it is an update.
func TestUnitModifyPlanClearOfImmutableCidrIsAReplacement(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	schema := NetworkResourceSchema(ctx)

	cases := map[string]struct {
		typeCode string
		replace  bool
	}{
		"amazonSubnet keeps its CIDR: replacement": {typeCode: "amazonSubnet", replace: true},
		"host edits its CIDR: update":              {typeCode: "host", replace: false},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			common := map[string]tftypes.Value{
				"name":      tftypes.NewValue(tftypes.String, "n"),
				"cloud_id":  tftypes.NewValue(tftypes.Number, 1),
				"group_id":  tftypes.NewValue(tftypes.Number, 1),
				"type_code": tftypes.NewValue(tftypes.String, tc.typeCode),
			}

			withCidr := map[string]tftypes.Value{
				"id":   tftypes.NewValue(tftypes.Number, 7),
				"cidr": tftypes.NewValue(tftypes.String, "10.0.0.0/24"),
			}
			for k, v := range common {
				withCidr[k] = v
			}

			// The configuration and the plan have no cidr; the state does.
			req := resource.ModifyPlanRequest{
				Config: tfsdk.Config{Schema: schema, Raw: networkObject(ctx, t, common)},
				State:  tfsdk.State{Schema: schema, Raw: networkObject(ctx, t, withCidr)},
				Plan:   tfsdk.Plan{Schema: schema, Raw: networkObject(ctx, t, common)},
			}
			resp := &resource.ModifyPlanResponse{Plan: req.Plan}

			// An unconfigured resource: no appliance, so the type comes from
			// the generated table and there is no private state.
			(&Resource{}).ModifyPlan(ctx, req, resp)

			if resp.Diagnostics.HasError() {
				t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
			}

			replaced := false
			for _, p := range resp.RequiresReplace {
				if p.Equal(path.Root("cidr")) {
					replaced = true
				}
			}

			if replaced != tc.replace {
				t.Fatalf("cidr requires replace = %v, want %v (RequiresReplace: %v)", replaced, tc.replace, resp.RequiresReplace)
			}

			var cidr types.String

			if diags := resp.Plan.GetAttribute(ctx, path.Root("cidr"), &cidr); diags.HasError() {
				t.Fatal(diags)
			}

			if !cidr.IsNull() {
				t.Fatalf("planned cidr = %v, want null", cidr)
			}
		})
	}
}

// TestUnitPoolClearRefusedWhereThePoolIsTheCluster: removing zone_pool_id is
// refused on the HVM types, whose plugin resolves the cluster from the pool
// for every operation, and allowed elsewhere.
func TestUnitPoolClearRefusedWhereThePoolIsTheCluster(t *testing.T) {
	t.Parallel()

	clearing := clearable.Pending{"zone_pool_id": {Name: "zone_pool_id"}}
	other := clearable.Pending{"description": {Name: "description"}}

	hvm := networkconfig.TypeSpec{Code: "ovsPortGroup", Name: "HVM Standard Network"}
	host := networkconfig.TypeSpec{Code: "host", Name: "Host Network"}

	if d := poolClearRefusal(&hvm, clearing); d == nil {
		t.Error("ovsPortGroup: want a refusal")
	} else if !strings.Contains(d.Detail(), "HVM Standard Network") || !strings.Contains(d.Detail(), "(ovsPortGroup)") {
		t.Errorf("ovsPortGroup: detail = %q, want the type named", d.Detail())
	}

	for name, tc := range map[string]struct {
		typ     *networkconfig.TypeSpec
		pending clearable.Pending
	}{
		"host clears its pool":           {&host, clearing},
		"hvm clearing something else":    {&hvm, other},
		"type unknown at plan time":      {nil, clearing},
		"every listed type is creatable": {nil, nil},
	} {
		if d := poolClearRefusal(tc.typ, tc.pending); d != nil {
			t.Errorf("%s: unexpected refusal: %s", name, d.Detail())
		}
	}

	// The list names real types, or it silently protects nothing.
	for code := range poolIsCluster {
		if _, ok := networkconfig.Lookup(code); !ok {
			t.Errorf("poolIsCluster names %q, which the catalogue does not have", code)
		}
	}
}
