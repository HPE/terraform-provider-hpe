// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instance

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// volumeElementType mirrors the shape of a single volumes element closely enough
// to exercise the comparison logic: an identity attribute the API supplies, and
// an attribute the practitioner sets.
var volumeElementType = tftypes.Object{
	AttributeTypes: map[string]tftypes.Type{
		"id":   tftypes.Number,
		"name": tftypes.String,
	},
}

var volumesListType = tftypes.List{ElementType: volumeElementType}

// resourceValue wraps a volumes list in an object, as it appears in the plan and
// state of a resource.
func resourceValue(volumes tftypes.Value) tftypes.Value {
	return tftypes.NewValue(
		tftypes.Object{
			AttributeTypes: map[string]tftypes.Type{"volumes": volumesListType},
		},
		map[string]tftypes.Value{"volumes": volumes},
	)
}

func volume(id, name tftypes.Value) tftypes.Value {
	return tftypes.NewValue(
		volumeElementType,
		map[string]tftypes.Value{"id": id, "name": name},
	)
}

func volumes(elements ...tftypes.Value) tftypes.Value {
	return tftypes.NewValue(volumesListType, elements)
}

func num(i int64) tftypes.Value    { return tftypes.NewValue(tftypes.Number, i) }
func str(s string) tftypes.Value   { return tftypes.NewValue(tftypes.String, s) }
func unknownNum() tftypes.Value    { return tftypes.NewValue(tftypes.Number, tftypes.UnknownValue) }
func unknownString() tftypes.Value { return tftypes.NewValue(tftypes.String, tftypes.UnknownValue) }

// TestUnitAttributeChanged verifies that only changes the practitioner actually
// made count as a change.
//
// The framework marks computed attributes as unknown whenever anything on the
// resource changes, so an unknown value carries no intent and must not be
// treated as a modification. A value the practitioner set, a change in the
// number of elements, or a newly added element must all be treated as changes.
func TestUnitAttributeChanged(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		plan  tftypes.Value
		state tftypes.Value
		want  bool
	}{
		{
			name:  "identical",
			plan:  resourceValue(volumes(volume(num(671), str("root")))),
			state: resourceValue(volumes(volume(num(671), str("root")))),
			want:  false,
		},
		{
			name: "only computed attributes unknown",
			// What the framework produces on an unrelated edit: identity is
			// unknown, the practitioner's value is untouched.
			plan:  resourceValue(volumes(volume(unknownNum(), str("root")))),
			state: resourceValue(volumes(volume(num(671), str("root")))),
			want:  false,
		},
		{
			name:  "practitioner changed a value",
			plan:  resourceValue(volumes(volume(unknownNum(), str("data")))),
			state: resourceValue(volumes(volume(num(671), str("root")))),
			want:  true,
		},
		{
			name: "element added",
			plan: resourceValue(volumes(
				volume(unknownNum(), str("root")),
				volume(unknownNum(), unknownString()),
			)),
			state: resourceValue(volumes(volume(num(671), str("root")))),
			want:  true,
		},
		{
			name:  "element removed",
			plan:  resourceValue(volumes(volume(unknownNum(), str("root")))),
			state: resourceValue(volumes(volume(num(671), str("root")), volume(num(672), str("data")))),
			want:  true,
		},
		{
			name: "elements reordered",
			plan: resourceValue(volumes(
				volume(unknownNum(), str("data")),
				volume(unknownNum(), str("root")),
			)),
			state: resourceValue(volumes(
				volume(num(671), str("root")),
				volume(num(672), str("data")),
			)),
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := attributeChanged(tt.plan, tt.state, "volumes"); got != tt.want {
				t.Errorf("attributeChanged() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestUnitAttributeChangedMissingAttribute verifies that an attribute which
// cannot be resolved is reported as changed, so that the caller leaves the plan
// alone rather than restoring a value it cannot verify.
func TestUnitAttributeChangedMissingAttribute(t *testing.T) {
	t.Parallel()

	plan := resourceValue(volumes(volume(num(1), str("root"))))
	state := resourceValue(volumes(volume(num(1), str("root"))))

	if !attributeChanged(plan, state, "does_not_exist") {
		t.Error("attributeChanged() = false for an unresolvable attribute, want true")
	}
}

// TestUnitFillUnknownsFromState verifies that unknown values are taken from
// prior state, and that a value with no counterpart in state is left unknown so
// the subsequent comparison still reports a difference.
func TestUnitFillUnknownsFromState(t *testing.T) {
	t.Parallel()

	plan := volumes(volume(unknownNum(), str("root")))
	state := volumes(volume(num(671), str("root")))

	filled, err := fillUnknownsFromState(plan, state)
	if err != nil {
		t.Fatalf("fillUnknownsFromState() error = %v", err)
	}

	if !filled.Equal(state) {
		t.Errorf("filled plan = %v, want it to equal state %v", filled, state)
	}
}

// TestUnitFillUnknownsFromStateLeavesUnmatchedUnknown verifies that an unknown
// with no counterpart in prior state stays unknown.
func TestUnitFillUnknownsFromStateLeavesUnmatchedUnknown(t *testing.T) {
	t.Parallel()

	plan := volumes(
		volume(unknownNum(), str("root")),
		volume(unknownNum(), str("data")),
	)
	state := volumes(volume(num(671), str("root")))

	filled, err := fillUnknownsFromState(plan, state)
	if err != nil {
		t.Fatalf("fillUnknownsFromState() error = %v", err)
	}

	if filled.Equal(state) {
		t.Error("filled plan equals state, want a difference for the added element")
	}
}

// TestUnitMatchesRestoreRule verifies which attribute paths a rule covers.
//
// A rule without fields covers only the root attribute. A rule with fields
// covers those fields within a collection element, and nothing else — in
// particular it must not match sibling attributes that the post-apply read
// sources from the API.
func TestUnitMatchesRestoreRule(t *testing.T) {
	t.Parallel()

	volumes := restoreRule{attribute: "volumes", fields: []string{"id"}}
	labels := restoreRule{attribute: "labels"}

	root := func(n string) *tftypes.AttributePath {
		return tftypes.NewAttributePath().WithAttributeName(n)
	}
	elem := func(n string, i int, f string) *tftypes.AttributePath {
		return tftypes.NewAttributePath().WithAttributeName(n).WithElementKeyInt(i).WithAttributeName(f)
	}

	tests := []struct {
		name string
		path *tftypes.AttributePath
		rule restoreRule
		want bool
	}{
		{"volume id covered", elem("volumes", 0, "id"), volumes, true},
		{"volume id at higher index covered", elem("volumes", 3, "id"), volumes, true},
		{"volume storage_profile NOT covered", elem("volumes", 0, "storage_profile"), volumes, false},
		{"volume controller_mount_point NOT covered", elem("volumes", 0, "controller_mount_point"), volumes, false},
		{"volumes root NOT covered when fields set", root("volumes"), volumes, false},
		{"different collection not covered", elem("network_interfaces", 0, "id"), volumes, false},
		{"labels root covered", root("labels"), labels, true},
		{"labels element not covered", elem("labels", 0, "id"), labels, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := matchesRestoreRule(tt.path, tt.rule); got != tt.want {
				t.Errorf("matchesRestoreRule() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestUnitRestoreFromStateNeverPinsNull is the regression test for the failure
// this design exists to avoid.
//
// Volume storage_profile is null in prior state on a post-apply read, but the
// API supplies a value. Pinning the null produced "Provider produced
// inconsistent result after apply: .volumes[0].storage_profile: was null, but
// now cty.StringVal(...)". The restore must leave such values unknown: both
// because storage_profile is not a covered field, and because a null is never
// pinned even where it is covered.
func TestUnitRestoreFromStateNeverPinsNull(t *testing.T) {
	t.Parallel()

	planVolumes := volumes(volume(unknownNum(), unknownString()))
	stateVolumes := volumes(volume(num(671), tftypes.NewValue(tftypes.String, nil)))

	plan := resourceValue(planVolumes)
	state := resourceValue(stateVolumes)

	// "name" stands in for a value the read path sources from the API: it is
	// null in prior state, so it must not be pinned.
	rules := []restoreRule{{attribute: "volumes", fields: []string{"id", "name"}}}

	got, err := restoreFromState(plan, state, rules)
	if err != nil {
		t.Fatalf("restoreFromState() error = %v", err)
	}

	idPath := tftypes.NewAttributePath().WithAttributeName("volumes").WithElementKeyInt(0).WithAttributeName("id")
	namePath := tftypes.NewAttributePath().WithAttributeName("volumes").WithElementKeyInt(0).WithAttributeName("name")

	gotID, err := valueAtPath(got, idPath)
	if err != nil {
		t.Fatalf("id path: %v", err)
	}
	if !gotID.Equal(num(671)) {
		t.Errorf("id = %v, want it restored to 671", gotID)
	}

	gotName, err := valueAtPath(got, namePath)
	if err != nil {
		t.Fatalf("name path: %v", err)
	}
	if gotName.IsKnown() {
		t.Errorf("name = %v, want it left unknown rather than pinned to null", gotName)
	}
}

// TestUnitRestoreFromStateLeavesUncoveredAttributesAlone verifies that an
// attribute outside the rule's fields is untouched even when prior state holds a
// perfectly good value for it.
func TestUnitRestoreFromStateLeavesUncoveredAttributesAlone(t *testing.T) {
	t.Parallel()

	plan := resourceValue(volumes(volume(unknownNum(), unknownString())))
	state := resourceValue(volumes(volume(num(671), str("kvm-cache-none"))))

	rules := []restoreRule{{attribute: "volumes", fields: []string{"id"}}}

	got, err := restoreFromState(plan, state, rules)
	if err != nil {
		t.Fatalf("restoreFromState() error = %v", err)
	}

	namePath := tftypes.NewAttributePath().WithAttributeName("volumes").WithElementKeyInt(0).WithAttributeName("name")
	gotName, err := valueAtPath(got, namePath)
	if err != nil {
		t.Fatalf("name path: %v", err)
	}

	if gotName.IsKnown() {
		t.Errorf("uncovered attribute = %v, want it left unknown", gotName)
	}
}
