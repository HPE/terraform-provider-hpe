// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package setting_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus"
	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus"
	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/resources/setting"
)

func TestMain(m *testing.M) {
	code := m.Run()

	testhelpers.WriteMergedResults()

	os.Exit(code)
}

func TestAccMorpheusSettingProvisioningExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	resourceConfig, err := setting.RenderSettingProvisioningConfig(t, nil)
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_setting_provisioning.tf_example_provisioning_setting",
			"allow_zone_selection",
			"false",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_setting_provisioning.tf_example_provisioning_setting",
			"allow_host_selection",
			"false",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_setting_provisioning.tf_example_provisioning_setting",
			"require_environments",
			"false",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_setting_provisioning.tf_example_provisioning_setting",
			"show_pricing",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_setting_provisioning.tf_example_provisioning_setting",
			"hide_datastore_stats",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_setting_provisioning.tf_example_provisioning_setting",
			"cross_tenant_naming_policies",
			"false",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_setting_provisioning.tf_example_provisioning_setting",
			"cloudinit_username",
			"cloudinit",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			// Plan
			{
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: true,
				Check:              checkFn,
				PlanOnly:           true,
			},
			// Apply
			{
				Config: providerConfig + resourceConfig,
				Check:  checkFn,
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
