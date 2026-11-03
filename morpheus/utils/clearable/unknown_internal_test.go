// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package clearable

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// unknownSchema has one attribute for each thing UnknownComputed has to get
// right: a computed id kept by UseStateForUnknown, computed-only volatile
// attributes, an Optional+Computed attribute the user never set, one with a
// default, one with RequiresReplace and no UseStateForUnknown, a configured
// nested object with a computed child, a computed nested object kept by
// UseStateForUnknown, and a set of nested objects with computed children.
func unknownSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name":         schema.StringAttribute{Required: true},
			"description":  schema.StringAttribute{Optional: true, Computed: true},
			"display_name": schema.StringAttribute{Optional: true, Computed: true},
			"status":       schema.StringAttribute{Computed: true},
			"last_updated": schema.StringAttribute{Computed: true},
			"kind": schema.StringAttribute{
				Optional:      true,
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"tier": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("basic"),
			},
			"config": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"mode":        schema.StringAttribute{Optional: true},
					"external_id": schema.StringAttribute{Computed: true},
				},
			},
			"owner": schema.SingleNestedAttribute{
				Optional:      true,
				Computed:      true,
				PlanModifiers: []planmodifier.Object{objectplanmodifier.UseStateForUnknown()},
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{Computed: true},
				},
			},
			"tags": schema.SetNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{Required: true},
						"id": schema.StringAttribute{
							Computed:      true,
							PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
						},
						"value": schema.StringAttribute{Computed: true},
					},
				},
			},
		},
	}
}

var (
	configType = tftypes.Object{AttributeTypes: map[string]tftypes.Type{"mode": tftypes.String, "external_id": tftypes.String}}
	ownerType  = tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String}}
	tagType    = tftypes.Object{AttributeTypes: map[string]tftypes.Type{
		"name": tftypes.String, "id": tftypes.String, "value": tftypes.String,
	}}
)

func tag(name, id, value string) tftypes.Value {
	return tftypes.NewValue(tagType, map[string]tftypes.Value{"name": str(name), "id": str(id), "value": str(value)})
}

// unknownState is a fully populated prior state for unknownSchema.
func unknownState(ctx context.Context, t *testing.T) tfsdk.State {
	t.Helper()

	return tfsdk.State{Schema: unknownSchema(), Raw: object(ctx, t, unknownSchema(), map[string]tftypes.Value{
		"id":           str("42"),
		"name":         str("n"),
		"description":  str("d"),
		"display_name": str("N"),
		"status":       str("offline"),
		"last_updated": str("yesterday"),
		"kind":         str("small"),
		"tier":         str("basic"),
		"config":       tftypes.NewValue(configType, map[string]tftypes.Value{"mode": str("auto"), "external_id": str("ext")}),
		"owner":        tftypes.NewValue(ownerType, map[string]tftypes.Value{"id": str("7")}),
		"tags":         tftypes.NewValue(tftypes.Set{ElementType: tagType}, []tftypes.Value{tag("env", "t1", "prod")}),
	})}
}

// unknownConfig is the configuration that goes with unknownState once the
// user has removed description: name, the tier default, config.mode and the
// tag name are set; everything else is null.
func unknownConfig(ctx context.Context, t *testing.T, extra map[string]tftypes.Value) tfsdk.Config {
	t.Helper()

	set := map[string]tftypes.Value{
		"name": str("n"),
		"tier": str("basic"),
		"config": tftypes.NewValue(configType, map[string]tftypes.Value{
			"mode": str("auto"), "external_id": tftypes.NewValue(tftypes.String, nil),
		}),
		"tags": tftypes.NewValue(tftypes.Set{ElementType: tagType}, []tftypes.Value{
			tftypes.NewValue(tagType, map[string]tftypes.Value{
				"name":  str("env"),
				"id":    tftypes.NewValue(tftypes.String, nil),
				"value": tftypes.NewValue(tftypes.String, nil),
			}),
		}),
	}
	for k, v := range extra {
		set[k] = v
	}

	return tfsdk.Config{Schema: unknownSchema(), Raw: object(ctx, t, unknownSchema(), set)}
}

// clearOnlyRequest is the ModifyPlan request for "description removed,
// nothing else changed": the proposed plan is the prior state.
func clearOnlyRequest(ctx context.Context, t *testing.T) (resource.ModifyPlanRequest, *resource.ModifyPlanResponse) {
	t.Helper()

	st := unknownState(ctx, t)
	req := resource.ModifyPlanRequest{
		Config: unknownConfig(ctx, t, nil),
		State:  st,
		Plan:   tfsdk.Plan(st),
	}

	return req, &resource.ModifyPlanResponse{Plan: req.Plan}
}

func planValue(ctx context.Context, t *testing.T, plan tfsdk.Plan, p path.Path) attr.Value {
	t.Helper()

	var v attr.Value

	mustNoDiags(t, plan.GetAttribute(ctx, p, &v))

	return v
}

func TestUnknownComputedMarksWhatTheFrameworkWould(t *testing.T) {
	ctx := context.Background()
	req, resp := clearOnlyRequest(ctx, t)

	mustNoDiags(t, UnknownComputed(ctx, req, resp))

	// Computed attributes the configuration leaves null become unknown:
	// the volatile computed-only ones, the Optional+Computed one the user
	// never set, and the computed child of the configured nested object.
	for _, p := range []path.Path{
		path.Root("status"), path.Root("last_updated"), path.Root("display_name"),
		path.Root("config").AtName("external_id"),
	} {
		if v := planValue(ctx, t, resp.Plan, p); !v.IsUnknown() {
			t.Errorf("%s = %v, want unknown", p, v)
		}
	}

	// The removed description is Optional+Computed too; UnknownComputed
	// treats it like any other, and Mark then decides.
	if v := planValue(ctx, t, resp.Plan, path.Root("description")); !v.IsUnknown() {
		t.Errorf("description = %v, want unknown before Mark", v)
	}

	// UseStateForUnknown keeps its say, on a string and on a whole object.
	if v, ok := planValue(ctx, t, resp.Plan, path.Root("id")).(types.String); !ok || v.ValueString() != "42" {
		t.Errorf("id = %v, want the prior 42", planValue(ctx, t, resp.Plan, path.Root("id")))
	}

	if v := planValue(ctx, t, resp.Plan, path.Root("owner").AtName("id")); v.IsUnknown() || v.IsNull() {
		t.Errorf("owner.id = %v, want the prior value", v)
	}

	// A default is never marked; a configured attribute is never touched.
	for p, want := range map[string]string{"tier": "basic", "name": "n"} {
		if v, ok := planValue(ctx, t, resp.Plan, path.Root(p)).(types.String); !ok || v.ValueString() != want {
			t.Errorf("%s = %v, want %q", p, planValue(ctx, t, resp.Plan, path.Root(p)), want)
		}
	}

	if v, ok := planValue(ctx, t, resp.Plan, path.Root("config").AtName("mode")).(types.String); !ok || v.ValueString() != "auto" {
		t.Errorf("config.mode = %v, want auto", planValue(ctx, t, resp.Plan, path.Root("config").AtName("mode")))
	}

	// RequiresReplace with no UseStateForUnknown fires, exactly as it does
	// on any other update of a resource with such an attribute.
	if v := planValue(ctx, t, resp.Plan, path.Root("kind")); !v.IsUnknown() {
		t.Errorf("kind = %v, want unknown", v)
	}

	if !resp.RequiresReplace.Contains(path.Root("kind")) {
		t.Errorf("RequiresReplace = %v, want kind", resp.RequiresReplace)
	}

	if len(resp.RequiresReplace) != 1 {
		t.Errorf("RequiresReplace = %v, want kind only", resp.RequiresReplace)
	}

	// Inside a set element nothing is marked, as in the framework's own
	// pass: the element is keyed by its value, and the configuration's
	// element — computed children null — never matches the plan's.
	tags, _, err := tftypes.WalkAttributePath(resp.Plan.Raw, tftypes.NewAttributePath().WithAttributeName("tags"))
	if err != nil {
		t.Fatal(err)
	}

	var elems []tftypes.Value
	if err := tags.(tftypes.Value).As(&elems); err != nil || len(elems) != 1 {
		t.Fatalf("tags = %v (%v), want one element", tags, err)
	}

	var fields map[string]tftypes.Value
	if err := elems[0].As(&fields); err != nil {
		t.Fatal(err)
	}

	for _, name := range []string{"id", "value"} {
		if !fields[name].IsKnown() || fields[name].IsNull() {
			t.Errorf("tags[0].%s = %v, want the prior value", name, fields[name])
		}
	}
}

func TestUnknownComputedIsIdempotent(t *testing.T) {
	ctx := context.Background()
	req, resp := clearOnlyRequest(ctx, t)

	mustNoDiags(t, UnknownComputed(ctx, req, resp))
	once := resp.Plan.Raw

	// The framework ran its pass, or we did: running again changes nothing
	// and asks for no second replacement.
	req.Plan = resp.Plan

	mustNoDiags(t, UnknownComputed(ctx, req, resp))

	if !resp.Plan.Raw.Equal(once) {
		t.Errorf("second run changed the plan:\n%s\n%s", once, resp.Plan.Raw)
	}

	if len(resp.RequiresReplace) != 1 {
		t.Errorf("RequiresReplace = %v, want kind once", resp.RequiresReplace)
	}
}

func TestUnknownComputedOnCreateAndDestroyDoesNothing(t *testing.T) {
	ctx := context.Background()
	st := unknownState(ctx, t)
	nullRaw := tftypes.NewValue(unknownSchema().Type().TerraformType(ctx), nil)

	// Create: no state.
	req := resource.ModifyPlanRequest{
		Config: unknownConfig(ctx, t, nil),
		State:  tfsdk.State{Schema: st.Schema, Raw: nullRaw},
		Plan:   tfsdk.Plan(st),
	}
	resp := &resource.ModifyPlanResponse{Plan: req.Plan}

	mustNoDiags(t, UnknownComputed(ctx, req, resp))

	if !resp.Plan.Raw.Equal(st.Raw) {
		t.Error("a create plan was modified")
	}

	// Destroy: no plan.
	req.State = st
	req.Plan = tfsdk.Plan{Schema: st.Schema, Raw: nullRaw}
	resp = &resource.ModifyPlanResponse{Plan: req.Plan}

	mustNoDiags(t, UnknownComputed(ctx, req, resp))

	if !resp.Plan.Raw.IsNull() {
		t.Error("a destroy plan was modified")
	}
}

// The end-to-end case that found the gap: a load balancer pool's status and
// last_updated changed on the appliance during a clear-only update, and the
// apply failed with "Provider produced inconsistent result".
func TestPlanMarksComputedUnknownWhenTheClearIsTheOnlyChange(t *testing.T) {
	ctx := context.Background()
	req, resp := clearOnlyRequest(ctx, t)

	attrs := Attributes{{Name: "description", Mode: Clear}}

	// The previous apply recorded description as configured.
	private := &fakePrivate{}
	mustNoDiags(t, Record(ctx, unknownConfig(ctx, t, map[string]tftypes.Value{"description": str("d")}), private, attrs))

	pending, diags := modifyPlan(ctx, req, resp, private, attrs)
	mustNoDiags(t, diags)

	if !pending.Clears("description") {
		t.Fatalf("pending = %v, want description", pending.Names())
	}

	if v := planValue(ctx, t, resp.Plan, path.Root("description")); !v.IsNull() {
		t.Errorf("description = %v, want null", v)
	}

	for _, name := range []string{"status", "last_updated"} {
		if v := planValue(ctx, t, resp.Plan, path.Root(name)); !v.IsUnknown() {
			t.Errorf("%s = %v, want unknown", name, v)
		}
	}

	if v, ok := planValue(ctx, t, resp.Plan, path.Root("id")).(types.String); !ok || v.ValueString() != "42" {
		t.Errorf("id = %v, want the prior 42", planValue(ctx, t, resp.Plan, path.Root("id")))
	}
}

// A plan with nothing pending must stay exactly as Terraform proposed it:
// marking computed attributes unknown there would turn every no-op plan
// into an update.
func TestPlanWithNothingPendingLeavesComputedAlone(t *testing.T) {
	ctx := context.Background()
	req, resp := clearOnlyRequest(ctx, t)
	before := resp.Plan.Raw

	// No record: description was never configured, so nothing is pending.
	pending, diags := modifyPlan(ctx, req, resp, &fakePrivate{}, Attributes{{Name: "description", Mode: Clear}})
	mustNoDiags(t, diags)

	if len(pending) != 0 {
		t.Fatalf("pending = %v, want none", pending.Names())
	}

	if !resp.Plan.Raw.Equal(before) {
		t.Error("a plan with nothing pending was modified")
	}
}
