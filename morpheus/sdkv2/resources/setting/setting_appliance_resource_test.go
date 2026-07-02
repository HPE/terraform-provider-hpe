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
// test mutates them and registers a t.Cleanup that restores them afterwards, so
// that an interrupted or failed run never leaves the shared appliance pointed at
// the test's placeholder proxy/URL values.
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

	resp, err := legacyClient.GetApplianceSettings(&sdklegacy.Request{})
	if err != nil {
		t.Fatalf("failed to read appliance settings: %v", err)
	}

	result, ok := resp.Result.(*sdklegacy.GetApplianceSettingsResult)
	if !ok || result.ApplianceSettings == nil {
		t.Fatal("appliance settings not found in response")
	}

	restore := applianceRestoreBody(result.ApplianceSettings)

	t.Cleanup(func() {
		if _, err := legacyClient.UpdateApplianceSettings(&sdklegacy.Request{
			Body: map[string]any{"applianceSettings": restore},
		}); err != nil {
			t.Errorf("failed to restore appliance settings: %v", err)
		}
	})
}

// applianceRestoreBody converts appliance settings read from the API into the
// map shape accepted by the appliance-settings PUT. Empty strings are preserved
// on purpose so that fields set by the test (for example a proxy host) are
// cleared again on restore. Role IDs are only included when registration is
// enabled, mirroring the RequiredWith constraints. Password fields are omitted
// because the API only returns password hashes, which cannot be replayed as
// plaintext.
func applianceRestoreBody(s *sdklegacy.ApplianceSettings) map[string]any {
	m := map[string]any{
		"applianceUrl":         s.ApplianceURL,
		"internalApplianceUrl": s.InternalApplianceURL,
		"corsAllowed":          s.CorsAllowed,
		"registrationEnabled":  s.RegistrationEnabled,
		"dockerPrivilegedMode": s.DockerPrivilegedMode,
		"smtpMailFrom":         s.SMTPMailFrom,
		"smtpServer":           s.SMTPServer,
		"smtpPort":             s.SMTPPort,
		"smtpSSL":              s.SMTPSSL,
		"smtpTLS":              s.SMTPTLS,
		"smtpUser":             s.SMTPUser,
		"proxyHost":            s.ProxyHost,
		"proxyPort":            s.ProxyPort,
		"proxyUser":            s.ProxyUser,
		"proxyDomain":          s.ProxyDomain,
		"proxyWorkstation":     s.ProxyWorkstation,
		"currencyProvider":     s.CurrencyProvider,
		"currencyKey":          s.CurrencyKey,
	}

	if s.RegistrationEnabled {
		m["defaultRoleId"] = s.DefaultRoleID
		m["defaultUserRoleId"] = s.DefaultUserRoleID
	}

	return m
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
