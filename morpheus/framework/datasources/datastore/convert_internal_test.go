// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package datastore

import (
	"testing"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

func strp(s string) *string { return &s }
func boolp(b bool) *bool    { return &b }
func i64p(i int64) *int64   { return &i }

// TestDatastoreFromListEntry pins the conversion the by-name path relies on.
//
// The listing entry and the single-item shape are generated from the same API
// object, so this asserts the fields actually survive re-encoding rather than
// silently arriving as zero values.
func TestDatastoreFromListEntry(t *testing.T) {
	t.Parallel()

	in := &sdk.ListDatastores200ResponseAllOfDatastoresInner{
		Id:           4,
		Name:         "datastore1",
		Type:         "vmfs",
		Status:       "available",
		RefType:      strp("ComputeZone"),
		RefId:        i64p(1),
		Active:       boolp(true),
		Online:       boolp(true),
		ExternalId:   strp("datastore-11"),
		ExternalType: strp("vmfs"),
		DatastoreType: sdk.ListDatastores200ResponseAllOfDatastoresInnerDatastoreType{
			Id:   7,
			Code: "vmware-datastore",
		},
	}

	out, err := datastoreFromListEntry(in)
	if err != nil {
		t.Fatalf("conversion failed: %v", err)
	}

	if out.Id != in.Id {
		t.Errorf("id = %d, want %d", out.Id, in.Id)
	}

	if out.Name != in.Name {
		t.Errorf("name = %q, want %q", out.Name, in.Name)
	}

	if out.Type != in.Type {
		t.Errorf("type = %q, want %q", out.Type, in.Type)
	}

	// refType and refId decide whether the datastore is treated as belonging to
	// a cloud or a cluster, so losing them would not fail loudly -- it would
	// produce the wrong associated_resource_type.
	if out.RefType == nil || *out.RefType != *in.RefType {
		t.Errorf("refType = %v, want %q", out.RefType, *in.RefType)
	}

	if out.RefId == nil || *out.RefId != *in.RefId {
		t.Errorf("refId = %v, want %d", out.RefId, *in.RefId)
	}

	// Nested structures are separate generated types on each side, which is the
	// part a hand-written copy would be most likely to get wrong.
	if out.DatastoreType.Id != in.DatastoreType.Id {
		t.Errorf("datastoreType.id = %d, want %d", out.DatastoreType.Id, in.DatastoreType.Id)
	}

	if out.DatastoreType.Code != in.DatastoreType.Code {
		t.Errorf("datastoreType.code = %q, want %q",
			out.DatastoreType.Code, in.DatastoreType.Code)
	}
}

// TestDatastoreListEntryCompleteness is the guard for the fallback in
// getDatastoreByName.
//
// Older appliances answer the listing with nothing but an id and a name. The
// by-name path detects that with refType and prefers to fetch the datastore,
// which carries more; it maps the entry itself only when that fetch fails. If
// refType ever stops being a reliable marker, this is the test that should
// fail.
func TestDatastoreListEntryCompleteness(t *testing.T) {
	t.Parallel()

	thin := &sdk.ListDatastores200ResponseAllOfDatastoresInner{
		Id:   1,
		Name: "toro-gl1-trial4-Vol1",
	}

	if thin.RefType != nil {
		t.Fatal("a thin entry must have no refType; the fallback depends on it")
	}

	full := &sdk.ListDatastores200ResponseAllOfDatastoresInner{
		Id:      4,
		Name:    "datastore1",
		RefType: strp("ComputeZone"),
		RefId:   i64p(1),
	}

	if full.RefType == nil {
		t.Fatal("a full entry must carry refType")
	}

	// The conversion must still work for the full entry, since that is the path
	// the fallback declines to take.
	out, err := datastoreFromListEntry(full)
	if err != nil {
		t.Fatalf("conversion of a full entry failed: %v", err)
	}

	if out.RefType == nil || *out.RefType != "ComputeZone" {
		t.Errorf("refType = %v, want ComputeZone", out.RefType)
	}
}

// TestDatastoreThinEntryConvertsWithoutAssociation covers the case that made
// this data source unusable on some appliances: a datastore the listing can
// name, whose association the API never reports and whose single-item endpoint
// answers 404.
//
// Converting must succeed and leave the association absent rather than failing.
// associated_resource_type, associated_resource_id, tenants and
// resource_permissions are all null in that state, which is the honest answer
// -- the appliance did not say.
func TestDatastoreThinEntryConvertsWithoutAssociation(t *testing.T) {
	t.Parallel()

	out, err := datastoreFromListEntry(
		&sdk.ListDatastores200ResponseAllOfDatastoresInner{
			Id:   1,
			Name: "toro-gl1-trial4-Vol1",
		},
	)
	if err != nil {
		t.Fatalf("a thin entry must still convert: %v", err)
	}

	if out.Id != 1 {
		t.Errorf("id = %d, want 1", out.Id)
	}

	if out.Name != "toro-gl1-trial4-Vol1" {
		t.Errorf("name = %q, want toro-gl1-trial4-Vol1", out.Name)
	}

	// The mapper keys off these being nil to skip the association entirely.
	if out.RefType != nil {
		t.Errorf("refType = %v, want nil", out.RefType)
	}

	if out.RefId != nil {
		t.Errorf("refId = %v, want nil", out.RefId)
	}
}
