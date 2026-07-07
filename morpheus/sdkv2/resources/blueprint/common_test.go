// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package blueprint

import (
	"reflect"
	"testing"
)

func TestMatchTemplatesWithSchema(t *testing.T) {
	cases := []struct {
		name      string
		templates []int64 // IDs returned by the API
		declared  []any   // IDs from config/state (SDKv2 stores TypeInt as int)
		want      []int64
	}{
		{
			name:      "MORPH-10324: declared ints are not doubled with phantom zeros",
			templates: []int64{1, 2},
			declared:  []any{int(1), int(2)},
			want:      []int64{1, 2},
		},
		{
			name:      "MORPH-10323: no phantom 0 entries / no perpetual diff",
			templates: []int64{105, 106},
			declared:  []any{int(105), int(106)},
			want:      []int64{105, 106},
		},
		{
			name:      "int64 declared elements are also accepted",
			templates: []int64{1, 2},
			declared:  []any{int64(1), int64(2)},
			want:      []int64{1, 2},
		},
		{
			name:      "server-added template is appended in API order",
			templates: []int64{1, 2},
			declared:  []any{int(1)},
			want:      []int64{1, 2},
		},
		{
			name:      "server-removed template is dropped, not zero-filled",
			templates: []int64{1},
			declared:  []any{int(1), int(2)},
			want:      []int64{1},
		},
		{
			name:      "declared order preserved, extras appended deterministically",
			templates: []int64{1, 2, 3},
			declared:  []any{int(3), int(1)},
			want:      []int64{3, 1, 2},
		},
		{
			name:      "empty declared returns all API templates in order",
			templates: []int64{1, 2},
			declared:  []any{},
			want:      []int64{1, 2},
		},
		{
			name:      "empty API returns empty, no zero padding",
			templates: []int64{},
			declared:  []any{int(1)},
			want:      []int64{},
		},
		{
			name:      "unexpected element types are skipped",
			templates: []int64{1, 2},
			declared:  []any{"1", int(2)},
			want:      []int64{2, 1},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := matchTemplatesWithSchema(tc.templates, tc.declared)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("matchTemplatesWithSchema(%v, %v) = %v, want %v",
					tc.templates, tc.declared, got, tc.want)
			}
		})
	}
}
