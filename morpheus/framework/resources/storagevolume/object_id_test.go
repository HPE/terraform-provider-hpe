// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package storagevolume

import "testing"

func TestObjectID(t *testing.T) {
	cases := []struct {
		name   string
		obj    map[string]interface{}
		wantID int64
		wantOk bool
	}{
		{"present float64 id", map[string]interface{}{"id": float64(42)}, 42, true},
		{"zero id is still present", map[string]interface{}{"id": float64(0)}, 0, true},
		{"nil map", nil, 0, false},
		{"missing id key", map[string]interface{}{"name": "vol"}, 0, false},
		{"non-numeric id", map[string]interface{}{"id": "42"}, 0, false},
		{"nil id value", map[string]interface{}{"id": nil}, 0, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			id, ok := objectID(tc.obj)
			if id != tc.wantID || ok != tc.wantOk {
				t.Fatalf("objectID(%v) = (%d, %v), want (%d, %v)",
					tc.obj, id, ok, tc.wantID, tc.wantOk)
			}
		})
	}
}
