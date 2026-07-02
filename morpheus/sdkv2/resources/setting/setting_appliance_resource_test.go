// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package setting_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	sdklegacy "github.com/HPE/terraform-provider-hpe/internal/sdk/legacy"
	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/client"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/resources/setting"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/skip"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

// snapshotAndRestoreApplianceSettings captures the appliance settings before a
// test mutates them and registers a t.Cleanup that restores them afterwards.
// This is an independent safety net (in addition to the resource's own
// restore-on-destroy) so that an interrupted or failed run never leaves the
// shared appliance pointed at the test's placeholder proxy/URL values.
func snapshotAndRestoreApplianceSettings(t *testing.T) {
	t.Helper()

	url, ok := testhelpers.LookupProviderEnv("", "url")
	if !ok {
		t.Fatal("TF_VAR_testacc_morpheus_url must be set")
	}

	token, _ := testhelpers.LookupProviderEnv("", "access_token")
	username, _ := testhelpers.LookupProviderEnv("", "username")
	password, _ := testhelpers.LookupProviderEnv("", "password")
	tenantSubdomain, _ := testhelpers.LookupProviderEnv("", "tenant_subdomain")
	_, insecure := testhelpers.LookupProviderEnv("", "insecure")

	legacyClient := client.NewLegacyClient(
		context.Background(),
		url,
		username,
		password,
		tenantSubdomain,
		token,
		sdklegacy.WithInsecure(insecure),
	)

	restore, err := setting.SnapshotApplianceSettingsForTest(legacyClient)
	if err != nil {
		t.Fatalf("failed to snapshot appliance settings: %v", err)
	}

	t.Cleanup(func() {
		if err := restore(); err != nil {
			t.Errorf("failed to restore appliance settings: %v", err)
		}
	})
}

func TestAccMorpheusSettingApplianceExampleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.Settings)

	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	// This example asserts environment-specific resource IDs (default user and
	// account roles) that must be pre-seeded, and it mutates global appliance
	// settings (URL, proxy, SMTP). Skip unless explicitly opted in.
	if skip.SkipByDefault(t) {
		t.Skip("set RUN_SKIPPED_BY_DEFAULT to run; needs seeded role IDs and mutates global appliance settings")
	}

	// Safety net: snapshot the shared appliance settings now and restore them
	// when the test finishes, so a failure or interruption cannot leave the
	// appliance reconfigured with the placeholder proxy/URL values below.
	snapshotAndRestoreApplianceSettings(t)

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
