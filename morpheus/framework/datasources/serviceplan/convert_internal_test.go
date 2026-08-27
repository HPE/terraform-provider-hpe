// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package serviceplan

import (
	"testing"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

func strp(s string) *string { return &s }
func i64p(i int64) *int64   { return &i }

// TestServicePlanFromListEntry pins the conversion the by-name path relies on.
//
// The listing entry and the single-item shape are generated from the same API
// object, so this asserts the fields actually survive re-encoding rather than
// silently arriving as zero values.
func TestServicePlanFromListEntry(t *testing.T) {
	t.Parallel()

	in := &sdk.ListServicePlans200ResponseAllOfServicePlansInner{
		Id:          i64p(407),
		Name:        strp("G1-Small"),
		Code:        strp("g1-small"),
		Description: strp("1 Core, 2GB Memory"),
		MaxMemory:   i64p(2147483648),
		MaxCores:    *sdk.NewNullableInt64(i64p(1)),
		ProvisionType: &sdk.ListServicePlans200ResponseAllOfServicePlansInnerProvisionType{
			Id:   i64p(6),
			Code: strp("vmware"),
		},
	}

	out, err := servicePlanFromListEntry(in)
	if err != nil {
		t.Fatalf("conversion failed: %v", err)
	}

	if out.Id == nil || *out.Id != *in.Id {
		t.Errorf("id = %v, want %d", out.Id, *in.Id)
	}

	if out.Name == nil || *out.Name != *in.Name {
		t.Errorf("name = %v, want %q", out.Name, *in.Name)
	}

	if out.Code == nil || *out.Code != *in.Code {
		t.Errorf("code = %v, want %q", out.Code, *in.Code)
	}

	// Numeric fields matter as much as the identifiers: a lost value looks like
	// "not set" rather than an error.
	if out.MaxMemory == nil || *out.MaxMemory != *in.MaxMemory {
		t.Errorf("maxMemory = %v, want %d", out.MaxMemory, *in.MaxMemory)
	}

	// provision_type_code is exposed by the schema and comes from this nested
	// structure, which is a separate generated type on each side.
	if out.ProvisionType == nil {
		t.Fatal("provisionType lost in conversion")
	}

	if out.ProvisionType.Code == nil || *out.ProvisionType.Code != *in.ProvisionType.Code {
		t.Errorf("provisionType.code = %v, want %q",
			out.ProvisionType.Code, *in.ProvisionType.Code)
	}
}

// TestServicePlanFromListEntryEmpty checks the conversion does not invent
// values for a sparse entry, since a plan with few fields set is ordinary.
func TestServicePlanFromListEntryEmpty(t *testing.T) {
	t.Parallel()

	out, err := servicePlanFromListEntry(
		&sdk.ListServicePlans200ResponseAllOfServicePlansInner{Id: i64p(1)},
	)
	if err != nil {
		t.Fatalf("conversion failed: %v", err)
	}

	if out.Id == nil || *out.Id != 1 {
		t.Errorf("id = %v, want 1", out.Id)
	}

	if out.Name != nil {
		t.Errorf("name = %v, want nil", out.Name)
	}

	if out.ProvisionType != nil {
		t.Errorf("provisionType = %v, want nil", out.ProvisionType)
	}
}
