package convert

import (
	"context"
	"math/big"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

// TestValueToAnyUnknown is the regression guard for the panic reported against
// v1.6.0 (MORPH-16244).
//
// An unknown NumberValue holds a nil *big.Float, so ValueBigFloat().Float64()
// dereferences nil. The guard only covered null, which was enough while every
// value came from a literal and became a crash the moment one came from a
// variable, a count or for_each reference, or anything else deferred to apply.
//
// Every case here panics without the fix rather than merely returning the wrong
// thing, so a failure is unmissable.
func TestUnitValueToAnyUnknown(t *testing.T) {
	t.Parallel()

	elem := map[string]attr.Type{"id": types.NumberType}

	cases := map[string]attr.Value{
		"number":  types.NumberUnknown(),
		"string":  types.StringUnknown(),
		"bool":    types.BoolUnknown(),
		"int64":   types.Int64Unknown(),
		"float64": types.Float64Unknown(),
		"list":    types.ListUnknown(types.NumberType),
		"set":     types.SetUnknown(types.NumberType),
		"map":     types.MapUnknown(types.NumberType),
		"object":  types.ObjectUnknown(elem),
		"tuple":   types.TupleUnknown([]attr.Type{types.NumberType}),
	}

	for name, value := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := ValueToAny(context.Background(), value)
			require.NoError(t, err)
			require.Nil(t, got, "an unknown value carries nothing to convert")
		})
	}
}

// TestValueToAnyNestedUnknown covers the shape the customer actually hit.
//
// The validator guards the outer value, but `config = { templateId = var.x }`
// is an object whose type is known and only whose member is unknown. The outer
// guard passes, the walk reaches the member, and that is where it panicked.
func TestUnitValueToAnyNestedUnknown(t *testing.T) {
	t.Parallel()

	obj := types.ObjectValueMust(
		map[string]attr.Type{
			"templateId": types.NumberType,
			"noAgent":    types.BoolType,
		},
		map[string]attr.Value{
			"templateId": types.NumberUnknown(),
			"noAgent":    types.BoolValue(true),
		},
	)

	got, err := ValueToAny(context.Background(), obj)
	require.NoError(t, err)

	m, ok := got.(map[string]any)
	require.True(t, ok, "an object converts to a map")

	// The known sibling survives; the unknown member is absent rather than a
	// zero value, so a caller cannot mistake "not yet known" for "set to 0".
	require.Equal(t, true, m["noAgent"])
	require.Nil(t, m["templateId"])
}

// TestValueToAnyIntegerBecomesFloat pins current behaviour rather than
// endorsing it.
//
// A Terraform number is a big.Float, and this converts it to float64
// regardless of whether the practitioner wrote an integer. An image id of 187
// arrives back as float64(187), not int64(187).
//
// That is fine for JSON, which does not distinguish them, and matches what the
// customer saw working with a literal. It is recorded here because the
// conversion is lossy in one direction that does matter: see the precision
// case below.
func TestUnitValueToAnyIntegerBecomesFloat(t *testing.T) {
	t.Parallel()

	got, err := ValueToAny(
		context.Background(),
		types.NumberValue(big.NewFloat(187)),
	)
	require.NoError(t, err)

	require.IsType(t, float64(0), got,
		"a whole number still converts through float64")
	require.InDelta(t, 187, got, 0)
}

// TestValueToAnyLargeIntegerLosesPrecision records the edge the float64
// conversion cannot represent.
//
// Beyond 2^53 an integer has no exact float64, so the conversion reports a loss
// of precision and fails rather than silently rounding. That is the right
// answer, but it means a sufficiently large id cannot be passed through a
// dynamic attribute at all.
func TestUnitValueToAnyLargeIntegerLosesPrecision(t *testing.T) {
	t.Parallel()

	// 2^53 + 1: the first integer float64 cannot hold exactly.
	big53 := new(big.Float).SetInt64(1<<53 + 1)

	_, err := ValueToAny(context.Background(), types.NumberValue(big53))
	require.Error(t, err,
		"an integer beyond float64's exact range is refused, not rounded")
	require.Contains(t, err.Error(), "loss of precision")
}
