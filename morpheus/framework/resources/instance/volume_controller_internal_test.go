// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instance

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// TestUnitReconcileControllerMountPoint covers MORPH-18424: a configured
// "-1:bus:type:unit" is a request for a new controller, which Morpheus reports
// back with the id it assigned. State keeps the configured form when the
// request was satisfied, and the API value otherwise.
func TestUnitReconcileControllerMountPoint(t *testing.T) {
	t.Parallel()

	s := func(v string) *string { return &v }

	tests := []struct {
		name       string
		configured types.String
		api        *string
		want       *string
	}{
		{"satisfied: keep configured form", types.StringValue("-1:1:6:0"), s("1559905:1:6:0"), s("-1:1:6:0")},
		{"different bus: report the API", types.StringValue("-1:1:6:0"), s("1559902:0:4:1"), s("1559902:0:4:1")},
		{"different unit: report the API", types.StringValue("-1:1:6:0"), s("1559905:1:6:1"), s("1559905:1:6:1")},
		{"existing controller: report the API", types.StringValue("1559902:0:4:1"), s("1559902:0:4:1"), s("1559902:0:4:1")},
		{"nothing configured: report the API", types.StringNull(), s("1559902:0:4:0"), s("1559902:0:4:0")},
		{"unknown configured: report the API", types.StringUnknown(), s("1559902:0:4:0"), s("1559902:0:4:0")},
		{"API reports nothing: nothing", types.StringValue("-1:1:6:0"), nil, nil},
		{"malformed configured value: report the API", types.StringValue("-1"), s("1559905:1:6:0"), s("1559905:1:6:0")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := reconcileControllerMountPoint(tt.configured, tt.api)
			switch {
			case got == nil && tt.want == nil:
			case got == nil || tt.want == nil || *got != *tt.want:
				t.Errorf("reconcileControllerMountPoint(%v, %v) = %v, want %v",
					tt.configured, deref(tt.api), deref(got), deref(tt.want))
			}
		})
	}
}
