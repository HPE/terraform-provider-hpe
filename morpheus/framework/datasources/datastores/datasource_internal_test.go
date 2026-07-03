// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package datastores

import (
	"regexp"
	"testing"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

// TestUnitDatastoreFieldValue exercises the field value extraction used for
// filter matching.
func TestUnitDatastoreFieldValue(t *testing.T) {
	t.Parallel()

	vis := "public"
	refType := "ComputeZone"

	code := sdk.NewNullableString(sdk.PtrString("ds-code-1"))

	ds := sdk.ListDatastores200ResponseAllOfDatastoresInner{
		Id:         42,
		Name:       "test-datastore",
		Code:       *code,
		Type:       "nfs",
		Status:     "provisioned",
		Visibility: &vis,
		RefType:    &refType,
	}

	tests := []struct {
		name    string
		field   string
		wantVal string
		wantOK  bool
	}{
		{"name field", "name", "test-datastore", true},
		{"code field", "code", "ds-code-1", true},
		{"type field", "type", "nfs", true},
		{"status field", "status", "provisioned", true},
		{"visibility field", "visibility", "public", true},
		{"associated_resource_type field (ComputeZone->Cloud)", "associated_resource_type", "Cloud", true},
		{"unknown field", "unknown", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gotVal, gotOK := datastoreFieldValue(&ds, tt.field)
			if gotOK != tt.wantOK {
				t.Errorf("datastoreFieldValue(%q) ok = %v, want %v", tt.field, gotOK, tt.wantOK)
			}
			if gotVal != tt.wantVal {
				t.Errorf("datastoreFieldValue(%q) = %q, want %q", tt.field, gotVal, tt.wantVal)
			}
		})
	}
}

// TestUnitDatastoreFieldValueNilOptionals verifies behaviour when optional
// pointer fields are nil.
func TestUnitDatastoreFieldValueNilOptionals(t *testing.T) {
	t.Parallel()

	ds := sdk.ListDatastores200ResponseAllOfDatastoresInner{
		Id:     1,
		Name:   "min",
		Type:   "generic",
		Status: "provisioned",
		// Visibility, RefType are nil
	}

	tests := []struct {
		name   string
		field  string
		wantOK bool
	}{
		{"visibility nil", "visibility", false},
		{"associated_resource_type nil", "associated_resource_type", false},
		{"code unset", "code", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, gotOK := datastoreFieldValue(&ds, tt.field)
			if gotOK != tt.wantOK {
				t.Errorf("datastoreFieldValue(%q) ok = %v, want %v", tt.field, gotOK, tt.wantOK)
			}
		})
	}
}

// TestUnitDatastoreMatchesFilters tests the filter logic across multiple
// filter blocks (AND semantics) with regex (OR within a block).
func TestUnitDatastoreMatchesFilters(t *testing.T) {
	t.Parallel()

	vis := "private"
	ds := sdk.ListDatastores200ResponseAllOfDatastoresInner{
		Id:         10,
		Name:       "prod-vmfs-01",
		Type:       "vmfs",
		Status:     "provisioned",
		Visibility: &vis,
	}

	tests := []struct {
		name    string
		filters []compiledFilter
		want    bool
	}{
		{
			"no filters matches everything",
			nil,
			true,
		},
		{
			"single matching filter",
			[]compiledFilter{
				{field: "type", res: []*regexp.Regexp{regexp.MustCompile("^vmfs$")}},
			},
			true,
		},
		{
			"single non-matching filter",
			[]compiledFilter{
				{field: "type", res: []*regexp.Regexp{regexp.MustCompile("^nfs$")}},
			},
			false,
		},
		{
			"OR within block matches",
			[]compiledFilter{
				{field: "type", res: []*regexp.Regexp{
					regexp.MustCompile("^nfs$"),
					regexp.MustCompile("^vmfs$"),
				}},
			},
			true,
		},
		{
			"AND across blocks - both match",
			[]compiledFilter{
				{field: "type", res: []*regexp.Regexp{regexp.MustCompile("vmfs")}},
				{field: "status", res: []*regexp.Regexp{regexp.MustCompile("provisioned")}},
			},
			true,
		},
		{
			"AND across blocks - second fails",
			[]compiledFilter{
				{field: "type", res: []*regexp.Regexp{regexp.MustCompile("vmfs")}},
				{field: "status", res: []*regexp.Regexp{regexp.MustCompile("failed")}},
			},
			false,
		},
		{
			"regex partial match",
			[]compiledFilter{
				{field: "name", res: []*regexp.Regexp{regexp.MustCompile("prod-.*-01")}},
			},
			true,
		},
		{
			"filter on absent optional field",
			[]compiledFilter{
				{field: "associated_resource_type", res: []*regexp.Regexp{regexp.MustCompile(".*")}},
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := datastoreMatchesFilters(&ds, tt.filters)
			if got != tt.want {
				t.Errorf("datastoreMatchesFilters() = %v, want %v", got, tt.want)
			}
		})
	}
}
