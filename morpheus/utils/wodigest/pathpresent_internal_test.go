// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package wodigest

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// An immutable write-only attribute under a variant block that is null in
// both the plan and the prior state must not be marked for replacement: the
// path does not exist and Terraform rejects the plan. The block in use is.
func TestUnitPathPresentSkipsNullVariantBlocks(t *testing.T) {
	t.Parallel()

	block := schema.SingleNestedAttribute{
		Optional:   true,
		Attributes: map[string]schema.Attribute{"pw": schema.StringAttribute{Optional: true, WriteOnly: true}},
	}
	s := schema.Schema{Attributes: map[string]schema.Attribute{"a": block, "b": block, "top": schema.StringAttribute{Optional: true}}}

	blockType := tftypes.Object{AttributeTypes: map[string]tftypes.Type{"pw": tftypes.String}}
	objType := tftypes.Object{AttributeTypes: map[string]tftypes.Type{"a": blockType, "b": blockType, "top": tftypes.String}}

	raw := tftypes.NewValue(objType, map[string]tftypes.Value{
		"a":   tftypes.NewValue(blockType, map[string]tftypes.Value{"pw": tftypes.NewValue(tftypes.String, nil)}),
		"b":   tftypes.NewValue(blockType, nil),
		"top": tftypes.NewValue(tftypes.String, nil),
	})

	plan := tfsdk.Plan{Schema: s, Raw: raw}
	state := tfsdk.State{Schema: s, Raw: raw}
	tr := Tracker{}

	if !tr.pathPresent(context.Background(), path.Root("a").AtName("pw"), plan, state) {
		t.Error("a.pw: the block in use must be present")
	}
	if tr.pathPresent(context.Background(), path.Root("b").AtName("pw"), plan, state) {
		t.Error("b.pw: a block null in plan and state must not be present")
	}
	if !tr.pathPresent(context.Background(), path.Root("top"), plan, state) {
		t.Error("top: a top-level attribute is always present")
	}
}

// A path is present when its block is set in the plan alone (a variant being
// switched to) or in the prior state alone (a variant being switched from),
// so a replacement can be hung on either side of a variant switch.
func TestUnitPathPresentPlanOrStateAlone(t *testing.T) {
	t.Parallel()

	block := schema.SingleNestedAttribute{
		Optional:   true,
		Attributes: map[string]schema.Attribute{"pw": schema.StringAttribute{Optional: true, WriteOnly: true}},
	}
	s := schema.Schema{Attributes: map[string]schema.Attribute{"a": block, "b": block}}

	blockType := tftypes.Object{AttributeTypes: map[string]tftypes.Type{"pw": tftypes.String}}
	objType := tftypes.Object{AttributeTypes: map[string]tftypes.Type{"a": blockType, "b": blockType}}

	set := tftypes.NewValue(blockType, map[string]tftypes.Value{"pw": tftypes.NewValue(tftypes.String, nil)})
	null := tftypes.NewValue(blockType, nil)

	aOnly := tftypes.NewValue(objType, map[string]tftypes.Value{"a": set, "b": null})
	bOnly := tftypes.NewValue(objType, map[string]tftypes.Value{"a": null, "b": set})

	plan := tfsdk.Plan{Schema: s, Raw: bOnly}   // switching to b
	state := tfsdk.State{Schema: s, Raw: aOnly} // from a
	tr := Tracker{}

	if !tr.pathPresent(context.Background(), path.Root("a").AtName("pw"), plan, state) {
		t.Error("a.pw: set in the prior state alone must be present")
	}
	if !tr.pathPresent(context.Background(), path.Root("b").AtName("pw"), plan, state) {
		t.Error("b.pw: set in the plan alone must be present")
	}

	none := tfsdk.Plan{Schema: s, Raw: tftypes.NewValue(objType, map[string]tftypes.Value{"a": null, "b": null})}
	if tr.pathPresent(context.Background(), path.Root("a").AtName("pw"), none, tfsdk.State{Schema: s, Raw: none.Raw}) {
		t.Error("a.pw: null on both sides must not be present")
	}
}

// importedHow names what happens to each configured kind of write-only value
// on an imported resource: recorded as applied on the next update (not
// re-sent), by replacement for the immutable ones, or both for a mixed
// tracker -- and only for the attributes actually configured.
func TestUnitImportedHow(t *testing.T) {
	t.Parallel()

	inPlace := Tracker{Attributes: map[string]string{"pw": ""}}
	if got := inPlace.importedHow([]string{"pw"}); !strings.Contains(got, "recorded as applied on the next update") ||
		strings.Contains(got, "replacement") || strings.Contains(got, "sent") {
		t.Errorf("in-place only: %q", got)
	}

	immutable := Tracker{Attributes: map[string]string{"ilo": ""}, Immutable: []string{"ilo"}}
	if got := immutable.importedHow([]string{"ilo"}); !strings.Contains(got, "only take effect through a replacement") {
		t.Errorf("immutable only: %q", got)
	}

	mixed := Tracker{Attributes: map[string]string{"ssh": "", "ilo": "", "console": ""}, Immutable: []string{"ilo", "console"}}

	got := mixed.importedHow([]string{"ssh", "console"})
	if !strings.HasPrefix(got, "ssh is recorded as applied on the next update") || !strings.Contains(got, "as is console, which") ||
		strings.Contains(got, "ilo") {
		t.Errorf("mixed, ilo not configured: %q", got)
	}

	// Only the configured in-place value: the immutable ones are not mentioned.
	if got := mixed.importedHow([]string{"ssh"}); strings.Contains(got, "replacement") {
		t.Errorf("mixed tracker, in-place value configured alone: %q", got)
	}
}

// planReplacement marks the immutable paths whose block is set -- the one in
// use, or both sides of a variant switch -- and errors when none is, so a
// changed immutable value is never silently ignored.
func TestUnitPlanReplacement(t *testing.T) {
	t.Parallel()

	block := schema.SingleNestedAttribute{
		Optional:   true,
		Attributes: map[string]schema.Attribute{"pw": schema.StringAttribute{Optional: true, WriteOnly: true}},
	}
	s := schema.Schema{Attributes: map[string]schema.Attribute{"a": block, "b": block}}

	blockType := tftypes.Object{AttributeTypes: map[string]tftypes.Type{"pw": tftypes.String}}
	objType := tftypes.Object{AttributeTypes: map[string]tftypes.Type{"a": blockType, "b": blockType}}

	set := tftypes.NewValue(blockType, map[string]tftypes.Value{"pw": tftypes.NewValue(tftypes.String, nil)})
	null := tftypes.NewValue(blockType, nil)

	aOnly := tftypes.NewValue(objType, map[string]tftypes.Value{"a": set, "b": null})
	bOnly := tftypes.NewValue(objType, map[string]tftypes.Value{"a": null, "b": set})
	none := tftypes.NewValue(objType, map[string]tftypes.Value{"a": null, "b": null})

	tr := Tracker{Resource: "host", Immutable: []string{"a.pw", "b.pw"}}
	ctx := context.Background()

	run := func(plan, state tftypes.Value) *resource.ModifyPlanResponse {
		req := resource.ModifyPlanRequest{Plan: tfsdk.Plan{Schema: s, Raw: plan}, State: tfsdk.State{Schema: s, Raw: state}}
		resp := &resource.ModifyPlanResponse{Plan: req.Plan}
		tr.planReplacement(ctx, req, resp)

		return resp
	}

	if resp := run(aOnly, aOnly); len(resp.RequiresReplace) != 1 || !resp.RequiresReplace[0].Equal(path.Root("a").AtName("pw")) ||
		resp.Diagnostics.HasError() {
		t.Errorf("block a in use: RequiresReplace=%v diags=%v, want a.pw only", resp.RequiresReplace, resp.Diagnostics)
	}

	if resp := run(bOnly, aOnly); len(resp.RequiresReplace) != 2 || resp.Diagnostics.HasError() {
		t.Errorf("variant switch a->b: RequiresReplace=%v diags=%v, want both paths", resp.RequiresReplace, resp.Diagnostics)
	}

	resp := run(none, none)
	if len(resp.RequiresReplace) != 0 || !resp.Diagnostics.HasError() ||
		!strings.Contains(resp.Diagnostics.Errors()[0].Summary(), "cannot be applied") {
		t.Errorf("no block set: RequiresReplace=%v diags=%v, want no paths and the error", resp.RequiresReplace, resp.Diagnostics)
	}
}
