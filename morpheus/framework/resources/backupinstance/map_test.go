// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package backupinstance

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

// TestMapBackupToStateStorageProviderId covers MORPH-16232: the API may return a
// null storage provider even when one was requested, so the configured value
// must be preserved rather than nulled out after apply.
func TestMapBackupToStateStorageProviderId(t *testing.T) {
	id := int64(2)

	tests := []struct {
		name   string
		backup *sdk.GetBackups200ResponseBackup
		plan   BackupInstanceModel
		want   types.Int64
	}{
		{
			name:   "nil provider, known plan preserves plan",
			backup: &sdk.GetBackups200ResponseBackup{},
			plan:   BackupInstanceModel{StorageProviderId: types.Int64Value(1)},
			want:   types.Int64Value(1),
		},
		{
			name:   "nil provider, unknown plan nulls",
			backup: &sdk.GetBackups200ResponseBackup{},
			plan:   BackupInstanceModel{StorageProviderId: types.Int64Unknown()},
			want:   types.Int64Null(),
		},
		{
			name:   "nil provider, null plan nulls",
			backup: &sdk.GetBackups200ResponseBackup{},
			plan:   BackupInstanceModel{StorageProviderId: types.Int64Null()},
			want:   types.Int64Null(),
		},
		{
			name: "non-nil provider overrides plan",
			backup: &sdk.GetBackups200ResponseBackup{
				StorageProvider: &sdk.GetBackups200ResponseBackupStorageProvider{Id: &id},
			},
			plan: BackupInstanceModel{StorageProviderId: types.Int64Value(1)},
			want: types.Int64Value(2),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mapBackupToState(tt.backup, tt.plan)
			if !got.StorageProviderId.Equal(tt.want) {
				t.Errorf("StorageProviderId = %v, want %v", got.StorageProviderId, tt.want)
			}
		})
	}
}
