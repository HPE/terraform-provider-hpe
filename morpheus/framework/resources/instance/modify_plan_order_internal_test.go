// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instance

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"
)

// nullObject builds a value of the object type with every attribute null,
// then applies set.
func nullObject(t *testing.T, typ tftypes.Type, set map[string]tftypes.Value) tftypes.Value {
	t.Helper()

	obj, ok := typ.(tftypes.Object)
	if !ok {
		t.Fatalf("type is %T, want tftypes.Object", typ)
	}

	vals := make(map[string]tftypes.Value, len(obj.AttributeTypes))
	for name, at := range obj.AttributeTypes {
		if v, ok := set[name]; ok {
			vals[name] = v

			continue
		}

		vals[name] = tftypes.NewValue(at, nil)
	}

	return tftypes.NewValue(obj, vals)
}

// clearOnlyInstancePlan is "labels removed, nothing else changed" on an
// instance with one volume. labels is Optional+Computed, so the proposed plan
// equals the prior state; the configuration has the volume without its
// computed id.
func clearOnlyInstancePlan(t *testing.T) (resource.ModifyPlanRequest, *resource.ModifyPlanResponse) {
	t.Helper()

	ctx := context.Background()
	s := InstanceResourceSchema(ctx)
	objType := s.Type().TerraformType(ctx)
	volType := objType.(tftypes.Object).AttributeTypes["volumes"].(tftypes.List).ElementType

	volume := func(id tftypes.Value) tftypes.Value {
		return tftypes.NewValue(tftypes.List{ElementType: volType}, []tftypes.Value{
			nullObject(t, volType, map[string]tftypes.Value{
				"id":   id,
				"name": tftypes.NewValue(tftypes.String, "root"),
				"size": tftypes.NewValue(tftypes.Number, 20),
			}),
		})
	}

	common := map[string]tftypes.Value{
		"name":      tftypes.NewValue(tftypes.String, "i"),
		"group_id":  tftypes.NewValue(tftypes.Number, 1),
		"cloud_id":  tftypes.NewValue(tftypes.Number, 2),
		"layout_id": tftypes.NewValue(tftypes.Number, 3),
		"plan_id":   tftypes.NewValue(tftypes.Number, 4),
	}

	stateSet := map[string]tftypes.Value{
		"id":      tftypes.NewValue(tftypes.Number, 100),
		"labels":  tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, []tftypes.Value{tftypes.NewValue(tftypes.String, "a")}),
		"volumes": volume(tftypes.NewValue(tftypes.Number, 671)),
		"connection_info": tftypes.NewValue(
			tftypes.List{ElementType: tftypes.String},
			[]tftypes.Value{tftypes.NewValue(tftypes.String, "10.0.0.5")},
		),
	}
	for k, v := range common {
		stateSet[k] = v
	}

	configSet := map[string]tftypes.Value{
		"volumes": volume(tftypes.NewValue(tftypes.Number, nil)),
	}
	for k, v := range common {
		configSet[k] = v
	}

	state := tfsdk.State{Schema: s, Raw: nullObject(t, objType, stateSet)}
	req := resource.ModifyPlanRequest{
		Config: tfsdk.Config{Schema: s, Raw: nullObject(t, objType, configSet)},
		State:  state,
		Plan:   tfsdk.Plan(state),
	}

	return req, &resource.ModifyPlanResponse{Plan: req.Plan}
}

func rawAt(t *testing.T, v tftypes.Value, steps ...tftypes.AttributePathStep) tftypes.Value {
	t.Helper()

	got, err := valueAtPath(v, tftypes.NewAttributePathWithSteps(steps))
	if err != nil {
		t.Fatalf("path %v: %v", steps, err)
	}

	return got
}

// TestUnitClearablePlanRunsBeforeTheRestore pins the order of the two plan
// passes: clearable.Plan marks the computed attributes unknown for a
// clear-only update, exactly as the framework does for any other update, and
// restoreUnchangedComputedAttributes then puts the stable ones back. The
// other way round the restore is undone and the plan shows the volume's id
// and the connection info as known after apply.
func TestUnitClearablePlanRunsBeforeTheRestore(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	volumeID := []tftypes.AttributePathStep{tftypes.AttributeName("volumes"), tftypes.ElementKeyInt(0), tftypes.AttributeName("id")}

	// The wrong order: restore, then the marking.
	req, resp := clearOnlyInstancePlan(t)
	restoreUnchangedComputedAttributes(ctx, req, resp)

	if diags := clearable.UnknownComputed(ctx, req, resp); diags.HasError() {
		t.Fatalf("UnknownComputed: %v", diags)
	}

	if v := rawAt(t, resp.Plan.Raw, volumeID...); v.IsKnown() {
		t.Fatalf("precondition: marking after the restore should leave volumes[0].id unknown, got %v", v)
	}

	// The order ModifyPlan uses: the marking, then the restore.
	req, resp = clearOnlyInstancePlan(t)

	if diags := clearable.UnknownComputed(ctx, req, resp); diags.HasError() {
		t.Fatalf("UnknownComputed: %v", diags)
	}

	restoreUnchangedComputedAttributes(ctx, req, resp)

	if v := rawAt(t, resp.Plan.Raw, volumeID...); !v.Equal(tftypes.NewValue(tftypes.Number, 671)) {
		t.Errorf("volumes[0].id = %v, want the prior 671 restored", v)
	}

	if v := rawAt(t, resp.Plan.Raw, tftypes.AttributeName("connection_info")); !v.IsKnown() {
		t.Errorf("connection_info = %v, want the prior value restored", v)
	}

	// Attributes nobody restores stay as the framework would leave them.
	if v := rawAt(t, resp.Plan.Raw, tftypes.AttributeName("compute_servers")); v.IsKnown() {
		t.Errorf("compute_servers = %v, want unknown", v)
	}
}
