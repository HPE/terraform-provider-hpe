// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package clusternamespace

import (
	"testing"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

// TestActiveFromNamespaceList covers MORPH-16158: `active` is looked up from the
// namespace list on import, since the single-namespace GET does not return it.
func TestActiveFromNamespaceList(t *testing.T) {
	id1 := int64(1)
	id2 := int64(2)
	inactive := false
	active := true

	items := []sdk.GetClusterNamespaces200ResponseAllOfNamespacesInner{
		{Id: &id1, Active: &inactive},
		{Id: &id2, Active: &active},
	}

	tests := []struct {
		name    string
		id      int64
		wantVal bool
		wantOk  bool
	}{
		{"inactive namespace", 1, false, true},
		{"active namespace", 2, true, true},
		{"missing id", 99, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := activeFromNamespaceList(items, tt.id)
			if got != tt.wantVal || ok != tt.wantOk {
				t.Errorf("activeFromNamespaceList(%d) = (%v, %v), want (%v, %v)", tt.id, got, ok, tt.wantVal, tt.wantOk)
			}
		})
	}
}
