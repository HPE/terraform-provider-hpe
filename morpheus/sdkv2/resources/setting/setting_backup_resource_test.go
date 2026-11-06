// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package setting_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/resources/setting"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/skip"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

func TestAccMorpheusSettingBackupExampleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.Settings)

	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	// The default backup schedule and storage bucket must exist on the appliance
	// (TF_ACC_BACKUP_SCHEDULE_ID / TF_ACC_BACKUP_BUCKET_ID; the rendered example's
	// ids otherwise), and the test mutates global backup settings. Skip unless
	// explicitly opted in.
	if skip.SkipByDefault(t) {
		t.Skip("set RUN_SKIPPED_BY_DEFAULT to run; needs an existing backup schedule and storage bucket (TF_ACC_BACKUP_SCHEDULE_ID, TF_ACC_BACKUP_BUCKET_ID) and mutates global settings")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	scheduleID := envOr("TF_ACC_BACKUP_SCHEDULE_ID", "3")
	bucketID := envOr("TF_ACC_BACKUP_BUCKET_ID", "17")

	resourceConfig, err := setting.RenderSettingBackupConfig(t, map[string]string{
		"Name":                         name,
		"DefaultBackupScheduleId":      scheduleID,
		"DefaultBackupStorageBucketId": bucketID,
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_setting_backup.example",
			"backup_appliance",
			"false",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_setting_backup.example",
			"create_backups",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_setting_backup.example",
			"default_backup_schedule_id",
			scheduleID,
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_setting_backup.example",
			"default_backup_storage_bucket_id",
			bucketID,
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_setting_backup.example",
			"retention_days",
			"21",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_setting_backup.example",
			"scheduled_backups",
			"true",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			// Apply
			{
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				Check:              checkFn,
			},
			// Plan after apply
			{
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
		},
	})
}

// envOr returns the environment variable's value, or the fallback when it is unset or empty.
func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return fallback
}
