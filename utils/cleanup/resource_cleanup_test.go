// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package cleanup_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/HPE/terraform-provider-hpe/utils/cleanup"
)

func taintSchema() schema.Schema {
	return schema.Schema{Attributes: map[string]schema.Attribute{
		"id":   schema.Int64Attribute{Computed: true},
		"name": schema.StringAttribute{Required: true},
		"uuid": schema.StringAttribute{Computed: true},
		"config": schema.SingleNestedAttribute{Computed: true, Attributes: map[string]schema.Attribute{
			"mode": schema.StringAttribute{Computed: true},
		}},
		"tags": schema.ListAttribute{ElementType: types.StringType, Computed: true},
	}}
}

func TestUnitTaintResourceStateNullsUnknowns(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	s := taintSchema()
	objType := s.Type().TerraformType(ctx).(tftypes.Object)

	// The state a Create leaves behind when it wrote the plan before a long
	// poll: the required attribute known, every computed attribute unknown.
	planned := tftypes.NewValue(objType, map[string]tftypes.Value{
		"id":   tftypes.NewValue(tftypes.Number, tftypes.UnknownValue),
		"name": tftypes.NewValue(tftypes.String, "c1"),
		"uuid": tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
		"config": tftypes.NewValue(objType.AttributeTypes["config"], map[string]tftypes.Value{
			"mode": tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
		}),
		"tags": tftypes.NewValue(objType.AttributeTypes["tags"], tftypes.UnknownValue),
	})
	state := tfsdk.State{Schema: s, Raw: planned}

	var diags diag.Diagnostics
	cleanup.TaintResourceState(ctx, cleanup.TaintResourceStateConfig{
		ResourceType: "cluster", ResourceID: 42, StateWriter: &state, Diagnostics: &diags,
	})

	if diags.HasError() {
		t.Fatalf("unexpected errors: %v", diags.Errors())
	}
	if len(diags.Warnings()) != 1 {
		t.Fatalf("expected the single taint warning, got %d", len(diags.Warnings()))
	}
	if !state.Raw.IsFullyKnown() {
		t.Fatalf("state must hold no unknown values after taint: %s", state.Raw)
	}

	var got struct {
		ID     *int64  `tfsdk:"id"`
		Name   string  `tfsdk:"name"`
		UUID   *string `tfsdk:"uuid"`
		Config *struct {
			Mode *string `tfsdk:"mode"`
		} `tfsdk:"config"`
		Tags []string `tfsdk:"tags"`
	}
	if d := state.Get(ctx, &got); d.HasError() {
		t.Fatalf("get: %v", d)
	}
	if got.ID == nil || *got.ID != 42 || got.Name != "c1" {
		t.Errorf("id/name: %+v", got)
	}
	// Unknown scalars and the unknown list become null; the nested object was
	// itself known, so it stays an object whose unknown attribute is now null.
	if got.UUID != nil || len(got.Tags) != 0 {
		t.Errorf("computed attributes must be null, not unknown: %+v", got)
	}
	if got.Config == nil || got.Config.Mode != nil {
		t.Errorf("nested unknown must become null inside its object: %+v", got.Config)
	}
}

// A Create that never wrote state taints a null state; the helper must still
// produce an object holding only the id.
func TestUnitTaintResourceStateFromNullState(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	s := taintSchema()
	state := tfsdk.State{Schema: s, Raw: tftypes.NewValue(s.Type().TerraformType(ctx), nil)}

	var diags diag.Diagnostics
	cleanup.TaintResourceState(ctx, cleanup.TaintResourceStateConfig{
		ResourceType: "cluster", ResourceID: 7, StateWriter: &state, Diagnostics: &diags,
	})
	if diags.HasError() || !state.Raw.IsFullyKnown() {
		t.Fatalf("errors=%v fullyKnown=%v", diags.Errors(), state.Raw.IsFullyKnown())
	}
}
