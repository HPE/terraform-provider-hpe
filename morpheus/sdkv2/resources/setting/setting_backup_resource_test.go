// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package setting_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/resources/setting"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/skip"
)

func TestAccMorpheusSettingBackupExampleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHave(t, capabilities.Settings)

	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	// This example asserts environment-specific resource IDs (default backup
	// schedule and storage bucket) that must be pre-seeded, and it mutates
	// global backup settings. Skip unless explicitly opted in.
	if skip.SkipByDefault(t) {
		t.Skip("set RUN_SKIPPED_BY_DEFAULT to run; needs seeded backup schedule/storage bucket IDs and mutates global settings")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := setting.RenderSettingBackupConfig(t, map[string]string{
		"Name": name,
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
			"3",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_setting_backup.example",
			"default_backup_storage_bucket_id",
			"17",
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
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), sdkv2morpheus.Provider()),
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
