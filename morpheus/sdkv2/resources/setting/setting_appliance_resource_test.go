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

func TestAccMorpheusSettingApplianceExampleOk(t *testing.T) {
	if capabilities.Missing(t, capabilities.All) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Skip("Skipping due to API error")

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := setting.RenderSettingApplianceConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_setting_appliance.example",
			"api_allowed_origins",
			"demo",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_setting_appliance.example",
			"appliance_url",
			"https://morpheus.test.local",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_setting_appliance.example",
			"currency_provider",
			"fixer",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_setting_appliance.example",
			"currency_provider_api_key",
			"5a4b220e-6f9f-43da-a572-390c8f6afed8",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_setting_appliance.example",
			"default_role_id",
			"5",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_setting_appliance.example",
			"default_user_role_id",
			"10",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_setting_appliance.example",
			"docker_privileged_mode",
			"false",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_setting_appliance.example",
			"internal_appliance_url",
			"https://pxemorpheus.test.local",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_setting_appliance.example",
			"proxy_domain",
			"test.local",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_setting_appliance.example",
			"proxy_host",
			"10.0.0.100",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_setting_appliance.example",
			"proxy_password",
			"Password123456",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_setting_appliance.example",
			"proxy_port",
			"3128",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_setting_appliance.example",
			"proxy_user",
			"jsmith",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_setting_appliance.example",
			"proxy_workstation",
			"work123",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_setting_appliance.example",
			"registration_enabled",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_setting_appliance.example",
			"smtp_from_address",
			"testemail@test.local",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_setting_appliance.example",
			"smtp_password",
			"Password12",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_setting_appliance.example",
			"smtp_port",
			"465",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_setting_appliance.example",
			"smtp_server",
			"smtp01.test.local",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_setting_appliance.example",
			"smtp_use_ssl",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_setting_appliance.example",
			"smtp_use_tls",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_setting_appliance.example",
			"smtp_username",
			"jsmith",
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
