// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package settingappliance_test

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/settingappliance"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/skip"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)

	testhelpers.WriteMergedResults()
	os.Exit(code)
}

// The tests here change settings that apply to the whole appliance. None runs
// in parallel: another test restoring its own snapshot of the same settings
// would overwrite the values under test.

// defaultExchangeURL is the plugin exchange URL every settings update resets
// the appliance to.
const defaultExchangeURL = "https://share.morpheusdata.com"

// snapshotAndRestoreApplianceSettings captures the appliance settings before a
// test changes them and registers a cleanup that puts them back, so that a
// failed or interrupted run never leaves the shared appliance pointed at the
// test's placeholder proxy or URL values. It skips the test on an appliance
// whose plugin exchange URL is not the default: every settings update resets
// it, and the API cannot put it back. It returns the snapshot.
func snapshotAndRestoreApplianceSettings(t *testing.T) *sdk.ListApplianceSettings200ResponseApplianceSettings {
	t.Helper()

	ctx := context.Background()

	client, err := testhelpers.NewClientForServer(ctx, "")
	if err != nil {
		t.Fatalf("creating the API client: %v", err)
	}

	result, _, err := client.ApplianceSettingsAPI.ListApplianceSettings(ctx).Execute()
	if err != nil {
		t.Fatalf("reading the appliance settings: %v", err)
	}

	if result == nil || result.ApplianceSettings == nil {
		t.Fatal("no applianceSettings in the response")
	}

	settings := result.ApplianceSettings

	if url := settings.ExchangeUrl.Get(); url != nil && *url != "" && *url != defaultExchangeURL {
		t.Skipf("the appliance's plugin exchange URL is %q: every settings update resets it to %s, "+
			"and the API cannot put it back", *url, defaultExchangeURL)
	}

	restore := applianceRestoreBody(settings)

	t.Cleanup(func() {
		putApplianceSettings(t, client, restore)
	})

	return settings
}

// putApplianceSettings sends a raw settings body, reporting a failure or a
// refusal (HTTP 200 with success false) as a test error.
func putApplianceSettings(t *testing.T, client *sdk.APIClient, settings map[string]any) {
	t.Helper()

	body := sdk.UpdateApplianceSettingsRequest{
		ApplianceSettings: &sdk.UpdateApplianceSettingsRequestApplianceSettings{
			AdditionalProperties: settings,
		},
	}

	result, _, err := client.ApplianceSettingsAPI.UpdateApplianceSettings(context.Background()).
		UpdateApplianceSettingsRequest(body).
		Execute()
	if err != nil {
		t.Errorf("updating the appliance settings: %v", err)

		return
	}

	if result != nil && result.Success != nil && !*result.Success {
		t.Errorf("the appliance refused the settings: %v", result.AdditionalProperties)
	}
}

// clearTestPasswordsAfter skips the test when the appliance already holds one
// of the passwords the test sets -- the API returns only their hashes, so they
// could not be put back -- and otherwise clears them when the test ends.
// passwords are the API's names (smtpPassword, proxyPassword). Call it after
// snapshotAndRestoreApplianceSettings, with the snapshot: cleanups run last in,
// first out, so the passwords are cleared before the snapshot, which leaves
// them alone, is restored.
func clearTestPasswordsAfter(
	t *testing.T,
	settings *sdk.ListApplianceSettings200ResponseApplianceSettings,
	passwords ...string,
) {
	t.Helper()

	hashes := map[string]sdk.NullableString{
		"smtpPassword":  settings.SmtpPasswordHash,
		"proxyPassword": settings.ProxyPasswordHash,
	}

	cleared := make(map[string]any, len(passwords))

	for _, name := range passwords {
		hash, ok := hashes[name]
		if !ok {
			t.Fatalf("no hash is known for %s", name)
		}

		if v := hash.Get(); v != nil && *v != "" {
			t.Skipf("the appliance holds a %s, which this test could not put back: the API returns only its hash", name)
		}

		cleared[name] = ""
	}

	client, err := testhelpers.NewClientForServer(context.Background(), "")
	if err != nil {
		t.Fatalf("creating the API client: %v", err)
	}

	t.Cleanup(func() {
		putApplianceSettings(t, client, cleared)
	})
}

// applianceRestoreBody converts settings read from the API into the PUT body
// that puts them back. It is a raw map, sent as the request model's additional
// properties, because the typed model cannot carry what a restore needs: empty
// strings, which clear what a test set (a proxy host, say), and the ports and
// role ids as the strings the appliance stores. The two default role ids are
// replayed whenever the appliance had them, whether or not registration is on:
// the appliance keeps them either way, and a test may have moved them with
// registration off. The passwords are left out: the API returns only their
// hashes, which cannot be replayed.
func applianceRestoreBody(s *sdk.ListApplianceSettings200ResponseApplianceSettings) map[string]any {
	str := func(v *string) string {
		if v == nil {
			return ""
		}

		return *v
	}

	flag := func(v *bool) bool { return v != nil && *v }

	m := map[string]any{
		"applianceUrl":         str(s.ApplianceUrl),
		"internalApplianceUrl": str(s.InternalApplianceUrl),
		"corsAllowed":          str(s.CorsAllowed),
		"registrationEnabled":  flag(s.RegistrationEnabled),
		"dockerPrivilegedMode": flag(s.DockerPrivilegedMode),
		"smtpMailFrom":         str(s.SmtpMailFrom.Get()),
		"smtpServer":           str(s.SmtpServer.Get()),
		"smtpPort":             str(s.SmtpPort.Get()),
		"smtpSSL":              flag(s.SmtpSSL),
		"smtpTLS":              flag(s.SmtpTLS),
		"smtpUser":             str(s.SmtpUser.Get()),
		"proxyHost":            str(s.ProxyHost.Get()),
		"proxyPort":            str(s.ProxyPort.Get()),
		"proxyUser":            str(s.ProxyUser.Get()),
		"proxyDomain":          str(s.ProxyDomain.Get()),
		"proxyWorkstation":     str(s.ProxyWorkstation.Get()),
		"noProxy":              str(s.NoProxy.Get()),
		"currencyProvider":     str(s.CurrencyProvider.Get()),
		"currencyKey":          str(s.CurrencyKey.Get()),
	}

	if id := str(s.DefaultRoleId); id != "" {
		m["defaultRoleId"] = id
	}

	if id := str(s.DefaultUserRoleId.Get()); id != "" {
		m["defaultUserRoleId"] = id
	}

	// The stats retention too, but only a value the appliance accepts: it
	// refuses the whole update otherwise.
	if n, ok := statsRetentionOf(s); ok && (n == 30 || n == 60 || n == 90) {
		m["statsRetainmentPeriod"] = strconv.FormatInt(n, 10)
	}

	// The settings TestAccMorpheusSettingAppliance_moreSettings changes, each
	// only with a value the appliance accepts, for the same reason. Unset, the
	// API reports the default the appliance applies (20 minutes, 90 days,
	// 1000, no host), so replaying it changes nothing.
	if n, ok := numericSetting(s.UserBrowserSessionTimeout); ok && n >= 1 {
		m["userBrowserSessionTimeout"] = strconv.FormatInt(n, 10)
	}

	if n, ok := numericSetting(s.ReportsRetainmentPeriod); ok && n >= 0 {
		m["reportsRetainmentPeriod"] = strconv.FormatInt(n, 10)
	}

	if n, ok := numericSetting(s.MaxOptionListSize); ok && n >= 1 && n <= maxOptionListSize {
		m["maxOptionListSize"] = strconv.FormatInt(n, 10)
	}

	m["httpBlacklistHosts"] = str(s.HttpBlacklistHosts.Get())

	return m
}

// maxOptionListSize is the largest option-list size the appliance accepts, in
// thousands of entries.
const maxOptionListSize = 10000

// numericSetting reads a numeric setting from settings read through the API,
// which renders it as a string.
func numericSetting(v sdk.NullableString) (int64, bool) {
	s := v.Get()
	if s == nil {
		return 0, false
	}

	n, err := strconv.ParseInt(strings.TrimSpace(*s), 10, 64)

	return n, err == nil
}

// importSingletonStep imports the settings into the resource at addr. The
// import ID identifies nothing: the id is "1" whatever it is.
func importSingletonStep(addr, importID string, ignore ...string) resource.TestStep {
	return resource.TestStep{
		ResourceName:            addr,
		ImportState:             true,
		ImportStateId:           importID,
		ImportStateVerify:       true,
		ImportStateVerifyIgnore: append(ignore, "write_only_revision"),
	}
}

func TestAccMorpheusSettingApplianceExampleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.Settings)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	// The example enables tenant registration with default roles whose ids
	// must exist on the appliance, and points the appliance at placeholder
	// URLs and a placeholder proxy for the duration. Skip unless explicitly
	// opted in.
	if skip.SkipByDefault(t) {
		t.Skip("set RUN_SKIPPED_BY_DEFAULT to run; needs existing role ids and changes appliance-wide settings")
	}

	// Safety net: snapshot the shared appliance settings now and restore them
	// when the test finishes, so a failure or interruption cannot leave the
	// appliance reconfigured with the placeholder values. The example's two
	// passwords cannot be restored, only cleared, so the test runs only on an
	// appliance that has neither.
	settings := snapshotAndRestoreApplianceSettings(t)
	clearTestPasswordsAfter(t, settings, "smtpPassword", "proxyPassword")

	providerConfig := testhelpers.ProviderBlock()

	resourceConfig, err := settingappliance.RenderSettingApplianceConfig(t, nil)
	if err != nil {
		t.Fatal(err)
	}

	const addr = "hpe_morpheus_setting_appliance.example"

	values := settingappliance.ExampleValues()

	var checks []resource.TestCheckFunc

	for attribute, key := range map[string]string{
		"api_allowed_origins":    "ApiAllowedOrigins",
		"appliance_url":          "ApplianceUrl",
		"currency_provider":      "CurrencyProvider",
		"default_role_id":        "DefaultRoleId",
		"default_user_role_id":   "DefaultUserRoleId",
		"docker_privileged_mode": "DockerPrivilegedMode",
		"internal_appliance_url": "InternalApplianceUrl",
		"no_proxy":               "NoProxy",
		"proxy_domain":           "ProxyDomain",
		"proxy_host":             "ProxyHost",
		"proxy_port":             "ProxyPort",
		"proxy_user":             "ProxyUser",
		"proxy_workstation":      "ProxyWorkstation",
		"registration_enabled":   "RegistrationEnabled",
		"smtp_from_address":      "SmtpFromAddress",
		"smtp_port":              "SmtpPort",
		"smtp_server":            "SmtpServer",
		"smtp_use_ssl":           "SmtpUseSsl",
		"smtp_use_tls":           "SmtpUseTls",
		"smtp_username":          "SmtpUsername",
	} {
		checks = append(checks, resource.TestCheckResourceAttr(addr, attribute, values[key]))
	}

	// The credentials are write-only: sent, never stored in state.
	for _, attribute := range []string{"smtp_password", "proxy_password", "currency_provider_api_key"} {
		checks = append(checks, resource.TestCheckNoResourceAttr(addr, attribute))
	}

	checks = append(checks, resource.TestCheckResourceAttrSet(addr, "write_only_revision"))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + resourceConfig,
				Check:  resource.ComposeAggregateTestCheckFunc(checks...),
			},
			// Plan after apply: the credentials are not re-sent.
			{
				Config:   providerConfig + resourceConfig,
				PlanOnly: true,
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
// success envelope (no applianceSettings object), which once failed every
// in-place change after it had been applied on the appliance. This
// create-then-update test asserts the second apply completes cleanly and reads
// the new value back, then imports the settings.
func TestAccMorpheusSettingAppliance_updateEnvelope(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.Settings)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	// api_allowed_origins is a global appliance setting: put it back afterwards.
	snapshotAndRestoreApplianceSettings(t)

	providerConfig := testhelpers.ProviderBlock()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + settingApplianceOriginsConfig("https://tf-acc-create.example.com"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						applianceUpdateEnvelopeAddr, "api_allowed_origins", "https://tf-acc-create.example.com",
					),
					resource.TestCheckResourceAttr(applianceUpdateEnvelopeAddr, "id", "1"),
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
			// Any import ID adopts the settings; "1" is the SDKv2 resource's.
			importSingletonStep(applianceUpdateEnvelopeAddr, "1"),
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
// create and on an in-place change. The resource requires both ids alongside
// registration_enabled whatever its value. Registration stays off: enabling it
// opens the appliance to sign-ups for the duration, and the role ids take the
// same path either way.
func TestAccMorpheusSettingAppliance_roleIds(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.Settings)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	// The default roles are global appliance settings: put them back
	// afterwards (the roles created here are destroyed with the config).
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
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + settingApplianceRoleIDsConfig(name, "a"),
				Check:  rolesMatch("a"),
			},
			// A fresh pair of roles: the ids change in place and must be sent,
			// with the registration flag, not dropped.
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

const applianceWriteOnlyAddr = "hpe_morpheus_setting_appliance.write_only"

// TestAccMorpheusSettingAppliance_writeOnly runs the write-only change
// detection matrix on smtp_password -- sent on create; a change re-sends it; a
// change to something else does not; an empty plan after each -- then the
// import forms. The value sent is observed in-process.
func TestAccMorpheusSettingAppliance_writeOnly(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.Settings)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	settings := snapshotAndRestoreApplianceSettings(t)

	// The test's SMTP password cannot be restored, only cleared: run only on
	// an appliance that has none.
	clearTestPasswordsAfter(t, settings, "smtpPassword")

	rec := testhelpers.RecordRequests(t)
	providerConfig := testhelpers.ProviderBlock()

	c := testhelpers.WriteOnlyCase{
		Address:   applianceWriteOnlyAddr,
		Attribute: "smtp_password",
		APIPath:   "/api/appliance-settings",
		APIKey:    "applianceSettings.smtpPassword",
		Config: func(value, extra string, unrelated bool) string {
			origins := "https://tf-acc-wo-a.example.com"
			if unrelated {
				origins = "https://tf-acc-wo-b.example.com"
			}

			return providerConfig + fmt.Sprintf(`
resource "hpe_morpheus_setting_appliance" "write_only" {
  api_allowed_origins = %q
  smtp_password       = %q
%s}
`, origins, value, extra)
		},
	}

	const v1, v2 = "tf-acc-Smtp-1!", "tf-acc-Smtp-2!"

	steps := testhelpers.WriteOnlyValueSteps(c, rec, v1, v2)
	steps = append(steps,
		importSingletonStep(applianceWriteOnlyAddr, "1", "smtp_password"),
		importSingletonStep(applianceWriteOnlyAddr, "1#smtp_password="+testhelpers.Sha256Hex(v2), "smtp_password"),
		resource.TestStep{
			ResourceName:  applianceWriteOnlyAddr,
			ImportState:   true,
			ImportStateId: "1#not_an_attribute=abc",
			ExpectError:   regexp.MustCompile(`(?s)not_an_attribute.*is\s+not\s+a\s+write-only\s+attribute`),
		},
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps:                    steps,
	})
}

// releasedSDKv2Version is the provider release the upgrade test starts from:
// the latest published release, every one of which has the SDKv2 resource.
// Move it to the last 2.2 release once that is published: the release most
// users will upgrade from.
const releasedSDKv2Version = "2.1.0"

const (
	applianceUpgradeAddr = "hpe_morpheus_setting_appliance.upgrade"
	upgradeOrigins       = "https://tf-acc-upgrade.example.com"
)

// settingApplianceUpgradeConfig sets one setting and no credentials. A
// credential would replace the shared appliance's SMTP password, which cannot
// be put back, and the released resource's plan never settles once a password
// is set: it compares a plain SHA-256 of the configured value with the
// appliance's HMAC.
const settingApplianceUpgradeConfig = `
resource "hpe_morpheus_setting_appliance" "upgrade" {
  api_allowed_origins = "` + upgradeOrigins + `"
}
`

// expectOnlyRevisionPlanned is a plan check: the resource is planned to keep
// every value of its prior state except write_only_revision. Every attribute
// of either side is compared, so one planned unknown, or present on one side
// only, fails it, as does a plan without a prior and a planned object.
type expectOnlyRevisionPlanned struct{ addr string }

func (e expectOnlyRevisionPlanned) CheckPlan(
	_ context.Context,
	req plancheck.CheckPlanRequest,
	resp *plancheck.CheckPlanResponse,
) {
	for _, rc := range req.Plan.ResourceChanges {
		if rc.Address != e.addr {
			continue
		}

		before, beforeOk := rc.Change.Before.(map[string]any)
		after, afterOk := rc.Change.After.(map[string]any)

		if !beforeOk || !afterOk || len(before) == 0 {
			resp.Error = fmt.Errorf("%s: the plan has no prior and planned object to compare", e.addr)

			return
		}

		keys := make(map[string]struct{}, len(before))
		for key := range before {
			keys[key] = struct{}{}
		}

		for key := range after {
			keys[key] = struct{}{}
		}

		for key := range keys {
			if key == "write_only_revision" {
				continue
			}

			if was, planned := before[key], after[key]; !reflect.DeepEqual(planned, was) {
				resp.Error = fmt.Errorf("%s: %s is planned as %v, was %v", e.addr, key, planned, was)

				return
			}
		}

		return
	}

	resp.Error = fmt.Errorf("%s is not in the plan", e.addr)
}

// TestAccMorpheusSettingAppliance_upgradeFromSDKv2 is the regression test for
// the move from the SDKv2 resource: the state a released provider writes
// upgrades in place to this build. The released resource writes version-0
// state ("id" "1", the ports and role ids as strings, the API's password hash
// in smtp_password, a stats_retainment_period it never read); this build
// upgrades it, and its first plan is the one in-place update the move
// documents -- write_only_revision, which records the credentials' digests --
// with every other value kept. The plan after it is empty.
//
// The released resource's Create also sends tenant registration and Docker
// privileged mode as off, whatever the configuration says (fixed by this
// build); the snapshot puts them back.
func TestAccMorpheusSettingAppliance_upgradeFromSDKv2(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.Settings)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	snapshotAndRestoreApplianceSettings(t)

	config := testhelpers.ProviderBlock() + settingApplianceUpgradeConfig
	thisBuild := testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil)

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				// The released provider, whose resource is the SDKv2 one.
				ExternalProviders: map[string]resource.ExternalProvider{
					"hpe": {Source: "HPE/hpe", VersionConstraint: releasedSDKv2Version},
				},
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(applianceUpgradeAddr, "id", "1"),
					resource.TestCheckResourceAttr(applianceUpgradeAddr, "api_allowed_origins", upgradeOrigins),
				),
			},
			{
				// This build, same configuration and state: the state is
				// upgraded, and the plan is one in-place update that sets
				// write_only_revision and keeps everything else.
				ProtoV6ProviderFactories: thisBuild,
				Config:                   config,
				ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(applianceUpgradeAddr, plancheck.ResourceActionUpdate),
					plancheck.ExpectUnknownValue(applianceUpgradeAddr, tfjsonpath.New("write_only_revision")),
					expectOnlyRevisionPlanned{addr: applianceUpgradeAddr},
				}},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(applianceUpgradeAddr, "id", "1"),
					resource.TestCheckResourceAttrSet(applianceUpgradeAddr, "stats_retainment_period"),
					resource.TestCheckResourceAttrSet(applianceUpgradeAddr, "write_only_revision"),
					resource.TestCheckResourceAttr(applianceUpgradeAddr, "api_allowed_origins", upgradeOrigins),
				),
			},
			{
				// And the plan after it is empty.
				ProtoV6ProviderFactories: thisBuild,
				Config:                   config,
				PlanOnly:                 true,
			},
		},
	})
}

// statsRetentionOf reads the stats retention from settings read through the
// API, which renders it as a string ("30" when none is set).
func statsRetentionOf(s *sdk.ListApplianceSettings200ResponseApplianceSettings) (int64, bool) {
	v := s.StatsRetainmentPeriod.Get()

	switch {
	case v == nil:
		return 0, false
	case v.Int32 != nil:
		return int64(*v.Int32), true
	case v.String != nil:
		n, err := strconv.ParseInt(strings.TrimSpace(*v.String), 10, 64)

		return n, err == nil
	default:
		return 0, false
	}
}

// higherStatsRetention is the next value the appliance accepts above current,
// and false when current is already the highest.
func higherStatsRetention(current int64) (int64, bool) {
	for _, days := range []int64{30, 60, 90} {
		if days > current {
			return days, true
		}
	}

	return 0, false
}

const applianceStatsRetentionAddr = "hpe_morpheus_setting_appliance.stats_retention"

// TestAccMorpheusSettingAppliance_statsRetainmentPeriod sets the stats
// retention, which the SDKv2 resource accepted but never sent: the value
// reaches the appliance and reads back, and the plan after it is empty. It
// only ever raises the retention: lowering it on a shared appliance could have
// the appliance purge stats while the test runs. The snapshot puts the
// original value back.
func TestAccMorpheusSettingAppliance_statsRetainmentPeriod(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.Settings)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	settings := snapshotAndRestoreApplianceSettings(t)

	current, ok := statsRetentionOf(settings)
	if !ok {
		t.Skip("the appliance does not report its stats retention")
	}

	target, ok := higherStatsRetention(current)
	if !ok {
		t.Skipf("the stats retention is already %d days, the most the appliance accepts; the test only raises it",
			current)
	}

	config := testhelpers.ProviderBlock() + fmt.Sprintf(`
resource "hpe_morpheus_setting_appliance" "stats_retention" {
  stats_retainment_period = %d
}
`, target)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.TestCheckResourceAttr(
					applianceStatsRetentionAddr, "stats_retainment_period", strconv.FormatInt(target, 10)),
			},
			{
				Config:   config,
				PlanOnly: true,
			},
		},
	})
}

const applianceMoreSettingsAddr = "hpe_morpheus_setting_appliance.more_settings"

// TestAccMorpheusSettingAppliance_moreSettings round-trips settings of each
// wire type the SDKv2 resource did not have -- the browser session timeout
// (an integer the PUT types as a string), the reports retention and the
// option-list size (integers) and the HTTP deny list (text) -- on create and
// on an in-place change, then imports them. They are the ones a shared
// appliance can lend a test: the timeout and the retention only rise, the
// option-list size moves by a thousand or two entries, the deny list gains a
// host in the reserved .invalid domain, and the snapshot puts each back. The
// password policy, user deactivation, sync intervals, Agent and SSH flags,
// locale, Twilio and the rest are left to the unit tests.
func TestAccMorpheusSettingAppliance_moreSettings(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.Settings)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	settings := snapshotAndRestoreApplianceSettings(t)

	timeout, timeoutOk := numericSetting(settings.UserBrowserSessionTimeout)
	reports, reportsOk := numericSetting(settings.ReportsRetainmentPeriod)
	listSize, listSizeOk := numericSetting(settings.MaxOptionListSize)

	if !timeoutOk || !reportsOk || !listSizeOk {
		t.Skip("the appliance does not report its session timeout, reports retention or option-list size")
	}

	if listSize+2 > maxOptionListSize {
		t.Skipf("the option-list size is already %d, too close to the most the appliance accepts", listSize)
	}

	denyList := func(host string) string {
		if v := settings.HttpBlacklistHosts.Get(); v != nil && strings.TrimSpace(*v) != "" {
			return *v + "," + host
		}

		return host
	}

	config := func(step int64, hosts string) string {
		return testhelpers.ProviderBlock() + fmt.Sprintf(`
resource "hpe_morpheus_setting_appliance" "more_settings" {
  user_browser_session_timeout = %d
  reports_retainment_period    = %d
  max_option_list_size         = %d
  http_blacklist_hosts         = %q
}
`, timeout+step, reports+step, listSize+step, hosts)
	}

	check := func(step int64, hosts string) resource.TestCheckFunc {
		return resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr(applianceMoreSettingsAddr,
				"user_browser_session_timeout", strconv.FormatInt(timeout+step, 10)),
			resource.TestCheckResourceAttr(applianceMoreSettingsAddr,
				"reports_retainment_period", strconv.FormatInt(reports+step, 10)),
			resource.TestCheckResourceAttr(applianceMoreSettingsAddr,
				"max_option_list_size", strconv.FormatInt(listSize+step, 10)),
			resource.TestCheckResourceAttr(applianceMoreSettingsAddr, "http_blacklist_hosts", hosts),
		)
	}

	steps := []resource.TestStep{
		{
			Config: config(1, denyList("tf-acc-a.invalid")),
			Check:  check(1, denyList("tf-acc-a.invalid")),
		},
		// All four change in place.
		{
			Config: config(2, denyList("tf-acc-b.invalid")),
			ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{
				plancheck.ExpectResourceAction(applianceMoreSettingsAddr, plancheck.ResourceActionUpdate),
			}},
			Check: check(2, denyList("tf-acc-b.invalid")),
		},
		{
			Config:   config(2, denyList("tf-acc-b.invalid")),
			PlanOnly: true,
		},
		// The import reads the same values back.
		{
			ResourceName:            applianceMoreSettingsAddr,
			ImportState:             true,
			ImportStateId:           "1",
			ImportStateVerify:       true,
			ImportStateVerifyIgnore: []string{"write_only_revision"},
		},
	}

	// "" clears a text setting. Morpheus stores an emptied setting as none and
	// reports it null, which must still read back as the "" configured, then
	// and on the refresh after. Only on an appliance whose deny list was
	// empty, which clearing it leaves as it was.
	if v := settings.HttpBlacklistHosts.Get(); v == nil || strings.TrimSpace(*v) == "" {
		steps = append(steps,
			resource.TestStep{
				Config: config(2, ""),
				Check:  check(2, ""),
			},
			resource.TestStep{
				Config:   config(2, ""),
				PlanOnly: true,
			},
		)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps:                    steps,
	})
}
