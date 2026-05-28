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

func TestAccMorpheusSettingMonitoringExampleOk(t *testing.T) {
	if capabilities.Missing(t, capabilities.All) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := setting.RenderSettingMonitoringConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_setting_monitoring.example",
			"morpheus_auto_create_checks",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_setting_monitoring.example",
			"morpheus_availability_precision",
			"4",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_setting_monitoring.example",
			"morpheus_availability_time_frame",
			"30",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_setting_monitoring.example",
			"morpheus_default_check_interval",
			"120",
		),

		// resource.TestCheckResourceAttr(
		// 	"hpe_morpheus_setting_monitoring.example",
		// 	"new_relic_license_key",
		// 	"ABC123",
		// ),
		//
		// resource.TestCheckResourceAttr(
		// 	"hpe_morpheus_setting_monitoring.example",
		// 	"new_relic_monitoring_enabled",
		// 	"true",
		// ),
		//
		// resource.TestCheckResourceAttr(
		// 	"hpe_morpheus_setting_monitoring.example",
		// 	"servicenow_close_incident_action",
		// 	"activity",
		// ),
		//
		// resource.TestCheckResourceAttr(
		// 	"hpe_morpheus_setting_monitoring.example",
		// 	"servicenow_integration_id",
		// 	"1",
		// ),
		//
		// resource.TestCheckResourceAttr(
		// 	"hpe_morpheus_setting_monitoring.example",
		// 	"servicenow_monitoring_enabled",
		// 	"true",
		// ),
		//
		// resource.TestCheckResourceAttr(
		// 	"hpe_morpheus_setting_monitoring.example",
		// 	"servicenow_new_incident_action",
		// 	"create",
		// ),
		//
		// resource.TestCheckResourceAttr(
		// 	"hpe_morpheus_setting_monitoring.example",
		// 	"servicenow_severity_critical_impact",
		// 	"low",
		// ),
		//
		// resource.TestCheckResourceAttr(
		// 	"hpe_morpheus_setting_monitoring.example",
		// 	"servicenow_severity_info_impact",
		// 	"high",
		// ),
		//
		// resource.TestCheckResourceAttr(
		// 	"hpe_morpheus_setting_monitoring.example",
		// 	"servicenow_severity_warning_impact",
		// 	"high",
		// ),
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
