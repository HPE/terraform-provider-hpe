// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package datastoretypes

import (
	"regexp"
	"testing"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

func TestUnitDatastoreTypeFieldValue(t *testing.T) {
	t.Parallel()

	name := "HPE Alletra MP"
	code := "hpedatastore-alletra-mp-bmaas"
	extType := "alletra-mp"
	diskType := "ssd"
	creatableTrue := true
	isPluginFalse := false

	dt := &sdk.ListDatastoreTypes200ResponseDatastoreTypesInner{
		Name:         &name,
		Code:         &code,
		ExternalType: &extType,
		DiskType:     &diskType,
		Creatable:    &creatableTrue,
		IsPlugin:     &isPluginFalse,
	}

	tests := []struct {
		field  string
		want   string
		wantOK bool
	}{
		{"name", "HPE Alletra MP", true},
		{"code", "hpedatastore-alletra-mp-bmaas", true},
		{"external_type", "alletra-mp", true},
		{"disk_type", "ssd", true},
		{"creatable", "true", true},
		{"is_plugin", "false", true},
		{"unknown_field", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.field, func(t *testing.T) {
			t.Parallel()

			got, ok := datastoreTypeFieldValue(dt, tt.field)
			if ok != tt.wantOK {
				t.Fatalf("datastoreTypeFieldValue(%q) ok = %v, want %v", tt.field, ok, tt.wantOK)
			}
			if got != tt.want {
				t.Errorf("datastoreTypeFieldValue(%q) = %q, want %q", tt.field, got, tt.want)
			}
		})
	}
}

func TestUnitDatastoreTypeFieldValueNilFields(t *testing.T) {
	t.Parallel()

	dt := &sdk.ListDatastoreTypes200ResponseDatastoreTypesInner{}

	fields := []string{"name", "code", "external_type", "disk_type", "creatable", "is_plugin"}

	for _, field := range fields {
		t.Run(field, func(t *testing.T) {
			t.Parallel()

			_, ok := datastoreTypeFieldValue(dt, field)
			if ok {
				t.Errorf("datastoreTypeFieldValue(%q) with nil field should return false", field)
			}
		})
	}
}

func TestUnitDatastoreTypeMatchesFilters(t *testing.T) {
	t.Parallel()

	name := "HPE Alletra MP"
	code := "hpedatastore-alletra-mp-bmaas"
	creatableTrue := true
	isPluginTrue := true

	dt := &sdk.ListDatastoreTypes200ResponseDatastoreTypesInner{
		Name:      &name,
		Code:      &code,
		Creatable: &creatableTrue,
		IsPlugin:  &isPluginTrue,
	}

	tests := []struct {
		name    string
		filters []compiledFilter
		want    bool
	}{
		{
			name:    "no filters matches everything",
			filters: nil,
			want:    true,
		},
		{
			name: "single filter matches",
			filters: []compiledFilter{
				{field: "name", res: []*regexp.Regexp{regexp.MustCompile("Alletra")}},
			},
			want: true,
		},
		{
			name: "single filter no match",
			filters: []compiledFilter{
				{field: "name", res: []*regexp.Regexp{regexp.MustCompile("^NonExistent$")}},
			},
			want: false,
		},
		{
			name: "multiple filters AND logic - all match",
			filters: []compiledFilter{
				{field: "name", res: []*regexp.Regexp{regexp.MustCompile("Alletra")}},
				{field: "code", res: []*regexp.Regexp{regexp.MustCompile("bmaas")}},
			},
			want: true,
		},
		{
			name: "multiple filters AND logic - one fails",
			filters: []compiledFilter{
				{field: "name", res: []*regexp.Regexp{regexp.MustCompile("Alletra")}},
				{field: "code", res: []*regexp.Regexp{regexp.MustCompile("^nope$")}},
			},
			want: false,
		},
		{
			name: "OR within a block - any value matches",
			filters: []compiledFilter{
				{field: "name", res: []*regexp.Regexp{
					regexp.MustCompile("^NoMatch$"),
					regexp.MustCompile("HPE"),
				}},
			},
			want: true,
		},
		{
			name: "boolean filter true",
			filters: []compiledFilter{
				{field: "creatable", res: []*regexp.Regexp{regexp.MustCompile("^true$")}},
			},
			want: true,
		},
		{
			name: "boolean filter false does not match true",
			filters: []compiledFilter{
				{field: "creatable", res: []*regexp.Regexp{regexp.MustCompile("^false$")}},
			},
			want: false,
		},
		{
			name: "filter on missing field returns false",
			filters: []compiledFilter{
				{field: "external_type", res: []*regexp.Regexp{regexp.MustCompile(".*")}},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := datastoreTypeMatchesFilters(dt, tt.filters)
			if got != tt.want {
				t.Errorf("datastoreTypeMatchesFilters() = %v, want %v", got, tt.want)
			}
		})
	}
}
