// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package setting

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// MORPH-20485: PUT /api/backup-settings reads the default storage bucket and the
// default schedule from inside "backupSettings". The resource used to place both
// beside it, where the appliance silently ignored them and every plan re-proposed
// the change.
func TestBackupSettingsRequestBody(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		raw  map[string]any
		want map[string]any
	}{
		{
			name: "bucket and schedule are nested under backupSettings",
			raw: map[string]any{
				"scheduled_backups":                true,
				"create_backups":                   true,
				"backup_appliance":                 false,
				"retention_days":                   21,
				"default_backup_storage_bucket_id": 17,
				"default_backup_schedule_id":       3,
			},
			want: map[string]any{
				"backupSettings": map[string]any{
					"backupsEnabled":       true,
					"createBackups":        true,
					"backupAppliance":      false,
					"retentionCount":       21,
					"defaultStorageBucket": map[string]any{"id": 17},
					"defaultSchedule":      map[string]any{"id": 3},
				},
			},
		},
		{
			name: "unset bucket and schedule are omitted, not sent as id 0",
			raw: map[string]any{
				"scheduled_backups": false,
				"create_backups":    false,
				"backup_appliance":  true,
				"retention_days":    7,
			},
			want: map[string]any{
				"backupSettings": map[string]any{
					"backupsEnabled":  false,
					"createBackups":   false,
					"backupAppliance": true,
					"retentionCount":  7,
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			d := schema.TestResourceDataRaw(t, ResourceSettingBackup().Schema, tc.raw)

			req, diags := backupSettingsRequest(d)
			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}
			if !reflect.DeepEqual(req.Body, tc.want) {
				t.Fatalf("request body\n got: %#v\nwant: %#v", req.Body, tc.want)
			}
			if len(req.Body) != 1 {
				t.Fatalf("the body must carry only backupSettings; got keys %v", keys(req.Body))
			}
		})
	}
}

func keys(m map[string]any) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}

	return out
}
