// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package datastore

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// TestUnitValidateBmaasConfig pins the plan-time contract for BMaaS datastores:
// when config_alletramp_bmaas is set, storage_server and resource_pool are
// required and the datastore must be Cloud-scoped. Non-BMaaS configs are never
// affected.
func TestUnitValidateBmaasConfig(t *testing.T) {
	t.Parallel()

	bmaasSet := ConfigAlletrampBmaasValue{
		ProtocolType: types.StringValue("iSCSI"),
		state:        attr.ValueStateKnown,
	}
	storageServer := StorageServerValue{Id: types.Int64Value(1), state: attr.ValueStateKnown}
	resourcePool := ResourcePoolValue{Id: types.Int64Value(1), state: attr.ValueStateKnown}

	// Empty blocks (e.g. storage_server = {}) are non-null but carry a null id,
	// which the API would reject; validation must catch this at plan time.
	storageServerNoID := StorageServerValue{Id: types.Int64Null(), state: attr.ValueStateKnown}
	resourcePoolNoID := ResourcePoolValue{Id: types.Int64Null(), state: attr.ValueStateKnown}

	tests := []struct {
		name         string
		model        DatastoreModel
		wantErrCount int
	}{
		{
			name:         "non-bmaas config is skipped",
			model:        DatastoreModel{ConfigAlletrampBmaas: NewConfigAlletrampBmaasValueNull()},
			wantErrCount: 0,
		},
		{
			name: "complete cloud-scoped bmaas is valid",
			model: DatastoreModel{
				ConfigAlletrampBmaas:   bmaasSet,
				StorageServer:          storageServer,
				ResourcePool:           resourcePool,
				AssociatedResourceType: types.StringValue(associatedResourceTypeCloud),
			},
			wantErrCount: 0,
		},
		{
			name: "missing storage_server",
			model: DatastoreModel{
				ConfigAlletrampBmaas:   bmaasSet,
				StorageServer:          NewStorageServerValueNull(),
				ResourcePool:           resourcePool,
				AssociatedResourceType: types.StringValue(associatedResourceTypeCloud),
			},
			wantErrCount: 1,
		},
		{
			name: "missing resource_pool",
			model: DatastoreModel{
				ConfigAlletrampBmaas:   bmaasSet,
				StorageServer:          storageServer,
				ResourcePool:           NewResourcePoolValueNull(),
				AssociatedResourceType: types.StringValue(associatedResourceTypeCloud),
			},
			wantErrCount: 1,
		},
		{
			name: "empty storage_server block (null id) is rejected",
			model: DatastoreModel{
				ConfigAlletrampBmaas:   bmaasSet,
				StorageServer:          storageServerNoID,
				ResourcePool:           resourcePool,
				AssociatedResourceType: types.StringValue(associatedResourceTypeCloud),
			},
			wantErrCount: 1,
		},
		{
			name: "empty resource_pool block (null id) is rejected",
			model: DatastoreModel{
				ConfigAlletrampBmaas:   bmaasSet,
				StorageServer:          storageServer,
				ResourcePool:           resourcePoolNoID,
				AssociatedResourceType: types.StringValue(associatedResourceTypeCloud),
			},
			wantErrCount: 1,
		},
		{
			name: "cluster scope is rejected",
			model: DatastoreModel{
				ConfigAlletrampBmaas:   bmaasSet,
				StorageServer:          storageServer,
				ResourcePool:           resourcePool,
				AssociatedResourceType: types.StringValue(associatedResourceTypeCluster),
			},
			wantErrCount: 1,
		},
		{
			name: "all three requirements missing",
			model: DatastoreModel{
				ConfigAlletrampBmaas:   bmaasSet,
				StorageServer:          NewStorageServerValueNull(),
				ResourcePool:           NewResourcePoolValueNull(),
				AssociatedResourceType: types.StringValue(associatedResourceTypeCluster),
			},
			wantErrCount: 3,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			diags := validateBmaasConfig(tc.model)
			if diags.ErrorsCount() != tc.wantErrCount {
				t.Fatalf("got %d error diagnostics, want %d: %v",
					diags.ErrorsCount(), tc.wantErrCount, diags)
			}
		})
	}
}
