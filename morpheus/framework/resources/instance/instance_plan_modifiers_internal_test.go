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
