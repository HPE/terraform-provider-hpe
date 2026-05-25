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
)

func TestAccMorpheusSettingBackupExampleOk(t *testing.T) {
	if capabilities.Missing(t, capabilities.Settings) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Skip("Skipping due to API error")
	// diagnostic_summary="Not found in response: BackupSettings"

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
