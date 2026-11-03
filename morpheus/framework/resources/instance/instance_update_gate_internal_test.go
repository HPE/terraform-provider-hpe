// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instance

import (
	"context"
	"sort"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

// The update gates are tested against values shaped by the REAL resource
// schema, exactly as Terraform hands them to Update. Every root attribute the
// schema declares is present in each value (null unless a case sets it), so
// the tests cannot drift from the schema: adding an attribute to the resource
// without classifying it in the gate lists shows up in
// TestUnitUpdateGateClassifiesEveryRootAttribute.

// instanceObjectType is the resource's root object type, derived from the
// schema rather than hand-written.
func instanceObjectType(t *testing.T) tftypes.Object {
	t.Helper()

	typ := InstanceResourceSchema(context.Background()).Type().TerraformType(context.Background())
	obj, ok := typ.(tftypes.Object)
	if !ok {
		t.Fatalf("resource type is %T, want tftypes.Object", typ)
	}

	return obj
}

// instanceValue builds a full resource value: every root attribute null except
// those in set. Attributes named in unknown are set to unknown instead.
func instanceValue(t *testing.T, obj tftypes.Object, set map[string]tftypes.Value, unknown ...string) tftypes.Value {
	t.Helper()

	vals := make(map[string]tftypes.Value, len(obj.AttributeTypes))
	for name, typ := range obj.AttributeTypes {
		vals[name] = tftypes.NewValue(typ, nil)
	}
	for name, v := range set {
		if _, ok := obj.AttributeTypes[name]; !ok {
			t.Fatalf("attribute %q is not in the schema", name)
		}
		vals[name] = v
	}
	for _, name := range unknown {
		vals[name] = tftypes.NewValue(obj.AttributeTypes[name], tftypes.UnknownValue)
	}

	return tftypes.NewValue(obj, vals)
}

// gateVolume builds one volumes element as it appears in a plan (planned=true)
// or in state. The practitioner sets root_volume/name/size/storage_type_id/
// datastore_id and omits the optional-and-computed controller_mount_point and
// storage_profile -- so in the plan those two are UNKNOWN while in state they
// hold the values the appliance reported. This is the shape that defeated the
// old model-level gate.
func gateVolume(t *testing.T, obj tftypes.Object, planned bool, size int64) tftypes.Value {
	t.Helper()

	elemType := obj.AttributeTypes["volumes"].(tftypes.List).ElementType.(tftypes.Object)
	vals := make(map[string]tftypes.Value, len(elemType.AttributeTypes))
	for name, typ := range elemType.AttributeTypes {
		vals[name] = tftypes.NewValue(typ, nil)
	}
	vals["root_volume"] = tftypes.NewValue(tftypes.Bool, true)
	vals["name"] = tftypes.NewValue(tftypes.String, "root")
	vals["size"] = tftypes.NewValue(tftypes.Number, size)
	vals["storage_type_id"] = tftypes.NewValue(tftypes.Number, int64(1))
	vals["datastore_id"] = tftypes.NewValue(tftypes.Number, int64(38658))
	vals["id"] = tftypes.NewValue(tftypes.Number, int64(4242))
	vals["actual_size"] = tftypes.NewValue(tftypes.Number, size)
	if planned {
		vals["controller_mount_point"] = tftypes.NewValue(tftypes.String, tftypes.UnknownValue)
		vals["storage_profile"] = tftypes.NewValue(tftypes.String, tftypes.UnknownValue)
	} else {
		vals["controller_mount_point"] = tftypes.NewValue(tftypes.String, "/dev/sda")
		vals["storage_profile"] = tftypes.NewValue(tftypes.String, "default")
	}

	return tftypes.NewValue(elemType, vals)
}

func gateVolumes(t *testing.T, obj tftypes.Object, elems ...tftypes.Value) tftypes.Value {
	t.Helper()

	return tftypes.NewValue(obj.AttributeTypes["volumes"], elems)
}

func gateStr(s string) tftypes.Value { return tftypes.NewValue(tftypes.String, s) }
func gateNum(i int64) tftypes.Value  { return tftypes.NewValue(tftypes.Number, i) }
func gateBool(b bool) tftypes.Value  { return tftypes.NewValue(tftypes.Bool, b) }
func gateStrs(t *testing.T, obj tftypes.Object, attr string, ss ...string) tftypes.Value {
	t.Helper()
	vals := make([]tftypes.Value, 0, len(ss))
	for _, s := range ss {
		vals = append(vals, gateStr(s))
	}

	return tftypes.NewValue(obj.AttributeTypes[attr], vals)
}

// gateTimeouts builds a timeouts block value with the given create timeout.
func gateTimeouts(t *testing.T, obj tftypes.Object, create string) tftypes.Value {
	t.Helper()

	typ := obj.AttributeTypes["timeouts"].(tftypes.Object)
	vals := make(map[string]tftypes.Value, len(typ.AttributeTypes))
	for name, at := range typ.AttributeTypes {
		vals[name] = tftypes.NewValue(at, nil)
	}
	vals["create"] = gateStr(create)

	return tftypes.NewValue(typ, vals)
}

// baseline is a realistic existing instance: the common attributes set, a
// volume with omitted optional-and-computed children, and a timeouts block.
func baseline(t *testing.T, obj tftypes.Object, planned bool) map[string]tftypes.Value {
	t.Helper()

	return map[string]tftypes.Value{
		"id":               gateNum(2213),
		"name":             gateStr("web-1"),
		"description":      gateStr("web server"),
		"instance_context": gateStr("dev"),
		"group_id":         gateNum(1),
		"plan_id":          gateNum(20),
		"labels":           gateStrs(t, obj, "labels", "a", "b"),
		"volumes":          gateVolumes(t, obj, gateVolume(t, obj, planned, 10)),
		"timeouts":         gateTimeouts(t, obj, "45m"),
	}
}

func with(m map[string]tftypes.Value, k string, v tftypes.Value) map[string]tftypes.Value {
	out := make(map[string]tftypes.Value, len(m)+1)
	for kk, vv := range m {
		out[kk] = vv
	}
	out[k] = v

	return out
}

// TestUnitUpdateGates is the exhaustive table for both gates. Each case is a
// single-attribute change on top of the baseline (unless noted) and states
// which appliance calls it must and must not trigger.
func TestUnitUpdateGates(t *testing.T) {
	t.Parallel()

	obj := instanceObjectType(t)
	state := instanceValue(t, obj, baseline(t, obj, false))

	cases := []struct {
		name       string
		plan       tftypes.Value
		wantPUT    bool
		wantResize bool
	}{
		// ---- the reason this ticket exists ----
		{
			name:    "timeouts-only change makes no call",
			plan:    instanceValue(t, obj, with(baseline(t, obj, true), "timeouts", gateTimeouts(t, obj, "90m"))),
			wantPUT: false, wantResize: false,
		},
		{
			name:    "no change at all makes no call",
			plan:    instanceValue(t, obj, baseline(t, obj, true)),
			wantPUT: false, wantResize: false,
		},
		{
			// The old gate failed exactly here: unknown volume children in the plan.
			name:    "unknown optional-and-computed volume children are not a change",
			plan:    instanceValue(t, obj, baseline(t, obj, true)),
			wantPUT: false, wantResize: false,
		},

		// ---- attributes the instance update call carries ----
		{name: "name", plan: instanceValue(t, obj, with(baseline(t, obj, true), "name", gateStr("web-2"))), wantPUT: true},
		{
			name:    "description",
			plan:    instanceValue(t, obj, with(baseline(t, obj, true), "description", gateStr("changed"))),
			wantPUT: true,
		},
		{
			name:    "instance_context",
			plan:    instanceValue(t, obj, with(baseline(t, obj, true), "instance_context", gateStr("prod"))),
			wantPUT: true,
		},
		{name: "group_id", plan: instanceValue(t, obj, with(baseline(t, obj, true), "group_id", gateNum(2))), wantPUT: true},
		{
			name:    "labels added",
			plan:    instanceValue(t, obj, with(baseline(t, obj, true), "labels", gateStrs(t, obj, "labels", "a", "b", "c"))),
			wantPUT: true,
		},
		{
			name:    "labels removed",
			plan:    instanceValue(t, obj, with(baseline(t, obj, true), "labels", gateStrs(t, obj, "labels", "a"))),
			wantPUT: true,
		},
		{
			name:    "labels reordered is not a change (set)",
			plan:    instanceValue(t, obj, with(baseline(t, obj, true), "labels", gateStrs(t, obj, "labels", "b", "a"))),
			wantPUT: false,
		},

		// ---- attributes the resize call carries ----
		{name: "plan_id", plan: instanceValue(t, obj, with(baseline(t, obj, true), "plan_id", gateNum(21))), wantResize: true},
		{
			name:       "volume resized",
			plan:       instanceValue(t, obj, with(baseline(t, obj, true), "volumes", gateVolumes(t, obj, gateVolume(t, obj, true, 20)))),
			wantResize: true,
		},
		{
			name: "volume added",
			plan: instanceValue(t, obj, with(baseline(t, obj, true), "volumes",
				gateVolumes(t, obj, gateVolume(t, obj, true, 10), gateVolume(t, obj, true, 10)))),
			wantResize: true,
		},

		// ---- provider-only and separately-handled attributes never trigger either ----
		{name: "wait_for_ip_address", plan: instanceValue(t, obj, with(baseline(t, obj, true), "wait_for_ip_address", gateBool(true)))},
		{name: "evars (own operations)", plan: instanceValue(t, obj, baseline(t, obj, true), "evars")},
		{name: "ports (replacement)", plan: instanceValue(t, obj, baseline(t, obj, true), "ports")},
		{
			name: "computed-only attributes unknown",
			plan: instanceValue(t, obj, baseline(t, obj, true), "status", "compute_servers", "connection_info", "container_id"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if got := isUpdateCallNeeded(tc.plan, state); got != tc.wantPUT {
				t.Errorf("isUpdateCallNeeded = %v, want %v", got, tc.wantPUT)
			}
			if got := isResizeNeeded(tc.plan, state); got != tc.wantResize {
				t.Errorf("isResizeNeeded = %v, want %v", got, tc.wantResize)
			}
		})
	}
}

// TestUnitUpdateGateClassifiesEveryRootAttribute is the drift guard: every
// root attribute in the schema must be deliberately classified as carried by
// the update call, carried by the resize call, or explicitly ignored. Adding an
// attribute to the resource without adding it to one of these lists fails
// here, so the gates cannot silently miss a new API-affecting attribute.
func TestUnitUpdateGateClassifiesEveryRootAttribute(t *testing.T) {
	t.Parallel()

	// Attributes that must NOT trigger either call, with the reason.
	ignored := map[string]string{
		// provider-only
		"timeouts":            "provider-only",
		"wait_for_ip_address": "provider-only",
		// applied through dedicated operations
		"evars": "own add/update and delete operations",
		// create-only: RequiresReplace, never reach Update
		"ports":             "RequiresReplace",
		"config_aws":        "RequiresReplace",
		"config_bmaas":      "RequiresReplace",
		"host_name":         "RequiresReplace",
		"network_domain_id": "RequiresReplace",
		"server_uuid":       "RequiresReplace",
		"server_uuids":      "RequiresReplace",
		"user_group":        "RequiresReplace",
		// create-only by contract but not marked RequiresReplace (pre-existing;
		// not applied on update today, tracked separately)
		"cloud_id":         "create-only",
		"config":           "create-only",
		"config_hvm":       "create-only",
		"config_vmware":    "create-only",
		"instance_type_id": "create-only",
		"layout_id":        "create-only",
		"layout_size":      "create-only",
		"task_set_id":      "create-only",
		// computed-only: never set by the practitioner
		"id":              "computed",
		"status":          "computed",
		"compute_servers": "computed",
		"connection_info": "computed",
		"container_id":    "computed",
	}

	obj := instanceObjectType(t)
	classified := map[string]bool{}
	for _, a := range putAttributes {
		classified[a] = true
	}
	for _, a := range resizeAttributes {
		classified[a] = true
	}
	for a := range ignored {
		if classified[a] {
			t.Errorf("%q is both ignored and gated", a)
		}
		classified[a] = true
	}

	var missing []string
	for name := range obj.AttributeTypes {
		if !classified[name] {
			missing = append(missing, name)
		}
	}
	sort.Strings(missing)
	if len(missing) > 0 {
		t.Errorf("root attributes not classified for the update gates (add to putAttributes, "+
			"resizeAttributes, or the ignored list with a reason): %v", missing)
	}

	for a := range classified {
		if _, ok := obj.AttributeTypes[a]; !ok {
			t.Errorf("%q is classified but no longer in the schema", a)
		}
	}
}

// TestUnitResizeRequestOmitsUnchangedServicePlanOptions pins that the resize
// request only carries service_plan_options when they changed. Previously an
// (all-null) options object was attached unconditionally, so any resize -- say
// for a volume change -- also told the appliance to apply service plan options
// it was never asked to touch, and the resize path's own "is a resize needed"
// check was always true.
func TestUnitResizeRequestOmitsUnchangedServicePlanOptions(t *testing.T) {
	t.Parallel()

	same := NewServicePlanOptionsValueNull()
	plan := InstanceModel{ServicePlanOptions: same}
	state := InstanceModel{ServicePlanOptions: same}

	req := &sdk.ResizeInstanceRequest{}
	addServicePlanOptionsToResizeRequest(plan, state, req)
	if req.ServicePlanOptions != nil {
		t.Fatal("unchanged service_plan_options must not be attached to the resize request")
	}

	// A real change is attached, with the memory converted to bytes.
	changed, d := NewServicePlanOptionsValue(
		ServicePlanOptionsValue{}.AttributeTypes(context.Background()),
		map[string]attr.Value{
			"cores_per_socket": types.Int64Null(),
			"max_cores":        types.Int64Null(),
			"max_memory":       types.Int64Value(2048),
		},
	)
	if d.HasError() {
		t.Fatal(d)
	}
	plan.ServicePlanOptions = changed

	req = &sdk.ResizeInstanceRequest{}
	addServicePlanOptionsToResizeRequest(plan, state, req)
	if req.ServicePlanOptions == nil || req.ServicePlanOptions.MaxMemory == nil {
		t.Fatal("changed service_plan_options must be attached")
	}
	if got := *req.ServicePlanOptions.MaxMemory; got != 2048<<20 {
		t.Fatalf("max_memory: got %d bytes, want %d", got, 2048<<20)
	}
}
