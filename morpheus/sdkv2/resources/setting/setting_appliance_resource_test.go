// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package setting_test

import (
	"context"
	"fmt"
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
// cleared again on restore. The two default role ids are replayed whenever the
// appliance had them, whether or not registration is enabled: the appliance
// keeps them either way, and a test may have moved them with registration
// off. Password fields are omitted because the API only returns password
// hashes, which cannot be replayed as plaintext.
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
		"noProxy":              s.NoProxy,
		"currencyProvider":     s.CurrencyProvider,
		"currencyKey":          s.CurrencyKey,
	}

	if s.DefaultRoleID != "" {
		m["defaultRoleId"] = s.DefaultRoleID
	}

	if s.DefaultUserRoleID != "" {
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
			"example",
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
			"no_proxy",
			"10.0.0.0/8,localhost,127.0.0.1",
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

const applianceUpdateEnvelopeAddr = "hpe_morpheus_setting_appliance.update_envelope"

// settingApplianceOriginsConfig renders a minimal appliance-settings resource
// that sets only api_allowed_origins, so the create and update payloads carry
// that one field and nothing else on the shared appliance changes.
func settingApplianceOriginsConfig(origins string) string {
	return fmt.Sprintf(`
resource "hpe_morpheus_setting_appliance" "update_envelope" {
  api_allowed_origins = %q
}
`, origins)
}

// TestAccMorpheusSettingAppliance_updateEnvelope is the regression for the
// appliance-settings update path: PUT /api/appliance-settings returns only a
// success envelope (no applianceSettings object), which the create path already
// tolerated but the update path did not, so every in-place change was applied
// on the appliance and then failed with "Not found in response:
// ApplianceSettings", leaving state stale. This create-then-update test asserts
// the second apply completes cleanly and reads the new value back.
func TestAccMorpheusSettingAppliance_updateEnvelope(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.Settings)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	// api_allowed_origins is a global appliance setting: put it back afterwards.
	// Not parallel: another test restoring its own snapshot of these shared
	// settings at the same time would overwrite the value under test.
	snapshotAndRestoreApplianceSettings(t)

	providerConfig := testhelpers.ProviderBlock()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + settingApplianceOriginsConfig("https://tf-acc-create.example.com"),
				Check: resource.TestCheckResourceAttr(
					applianceUpdateEnvelopeAddr, "api_allowed_origins", "https://tf-acc-create.example.com",
				),
			},
			// The update path is the regression target: it must apply cleanly
			// and state must carry the new value.
			{
				Config: providerConfig + settingApplianceOriginsConfig("https://tf-acc-update.example.com"),
				Check: resource.TestCheckResourceAttr(
					applianceUpdateEnvelopeAddr, "api_allowed_origins", "https://tf-acc-update.example.com",
				),
			},
			// And the plan after the update is empty.
			{
				Config:   providerConfig + settingApplianceOriginsConfig("https://tf-acc-update.example.com"),
				PlanOnly: true,
			},
		},
	})
}

const applianceRoleIDsAddr = "hpe_morpheus_setting_appliance.role_ids"

// settingApplianceRoleIDsConfig renders two roles created for the run and an
// appliance-settings resource that names them as the default tenant and user
// roles with registration left off. The roles' addresses carry the pair label,
// so a step that renders another pair creates new roles (new ids) rather than
// renaming the first pair in place.
func settingApplianceRoleIDsConfig(name, pair string) string {
	return fmt.Sprintf(`
resource "hpe_morpheus_role" "tenant_%[2]s" {
  name      = "%[1]s-tenant-%[2]s"
  role_type = "tenant"
}

resource "hpe_morpheus_role" "user_%[2]s" {
  name      = "%[1]s-user-%[2]s"
  role_type = "user"
}

resource "hpe_morpheus_setting_appliance" "role_ids" {
  registration_enabled = false
  default_role_id      = hpe_morpheus_role.tenant_%[2]s.id
  default_user_role_id = hpe_morpheus_role.user_%[2]s.id
}
`, name, pair)
}

// TestAccMorpheusSettingAppliance_roleIds asserts that the two default role ids
// reach the appliance and round-trip into state while registration is off, on
// create and on an in-place change. The schema requires both ids alongside
// registration_enabled whatever its value, so this is an ordinary
// configuration; dropping them while registration was off left the appliance's
// own ids in state and the configuration in a perpetual diff. Registration
// stays off: enabling it opens the appliance to sign-ups for the duration, and
// the role ids take the same path either way.
func TestAccMorpheusSettingAppliance_roleIds(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.Settings)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	// The default roles are global appliance settings: put them back
	// afterwards (the roles created here are destroyed with the config). Not
	// parallel, for the same reason as the update-envelope test.
	snapshotAndRestoreApplianceSettings(t)

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())

	rolesMatch := func(pair string) resource.TestCheckFunc {
		return resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr(applianceRoleIDsAddr, "registration_enabled", "false"),
			resource.TestCheckResourceAttrPair(
				applianceRoleIDsAddr, "default_role_id", "hpe_morpheus_role.tenant_"+pair, "id"),
			resource.TestCheckResourceAttrPair(
				applianceRoleIDsAddr, "default_user_role_id", "hpe_morpheus_role.user_"+pair, "id"),
		)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + settingApplianceRoleIDsConfig(name, "a"),
				Check:  rolesMatch("a"),
			},
			// A fresh pair of roles: the ids change in place (registration_enabled
			// is unchanged, so no replacement) and must be sent, not dropped.
			{
				Config: providerConfig + settingApplianceRoleIDsConfig(name, "b"),
				Check:  rolesMatch("b"),
			},
			{
				Config:   providerConfig + settingApplianceRoleIDsConfig(name, "b"),
				PlanOnly: true,
			},
		},
	})
}
