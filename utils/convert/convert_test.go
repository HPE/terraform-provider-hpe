// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package convert

import (
	"math"
	"testing"
)

func TestInt32ToType(t *testing.T) {
	t.Parallel()

	if got := Int32ToType(nil); !got.IsNull() {
		t.Errorf("nil should map to a null Int64, got %v", got)
	}

	for _, want := range []int32{0, 7, -7, math.MaxInt32, math.MinInt32} {
		v := want

		got := Int32ToType(&v)
		if got.IsNull() {
			t.Errorf("%d: unexpected null", want)

			continue
		}

		if got.ValueInt64() != int64(want) {
			t.Errorf("%d: widened to %d", want, got.ValueInt64())
		}
	}
}
