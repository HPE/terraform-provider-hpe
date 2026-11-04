// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package settingappliance

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// upgradeFixture runs the upgrader on testdata/state_v0.json -- the stored
// instance of the SDKv2 resource, of which Terraform hands the provider the
// attributes object as the raw state -- after edit has changed its attributes.
// It decodes the raw state exactly as the framework does for an upgrade:
// against the version-0 schema, ignoring attributes the schema lacks.
func upgradeFixture(t *testing.T, edit func(attrs map[string]any)) (resourceModel, diag.Diagnostics) {
	t.Helper()

	ctx := context.Background()

	raw, err := os.ReadFile("testdata/state_v0.json")
	if err != nil {
		t.Fatal(err)
	}

	var instance struct {
		Attributes map[string]any `json:"attributes"`
	}

	if err := json.Unmarshal(raw, &instance); err != nil {
		t.Fatal(err)
	}

	if edit != nil {
		edit(instance.Attributes)
	}

	attrs, err := json.Marshal(instance.Attributes)
	if err != nil {
		t.Fatal(err)
	}

	prior := settingApplianceSchemaV0()
	priorType := prior.Type().TerraformType(ctx).(tftypes.Object)

	rawState := tfprotov6.RawState{JSON: attrs}

	value, err := rawState.UnmarshalWithOpts(priorType, tfprotov6.UnmarshalOpts{
		ValueFromJSONOpts: tftypes.ValueFromJSONOpts{IgnoreUndefinedAttributes: true},
	})
	if err != nil {
		t.Fatalf("decode of the SDKv2 state against the version-0 schema failed: %v", err)
	}

	req := resource.UpgradeStateRequest{State: &tfsdk.State{Schema: prior, Raw: value}}
	resp := resource.UpgradeStateResponse{State: tfsdk.State{Schema: resourceSchema(ctx)}}

	upgradeStateV0toV1(ctx, req, &resp)

	var upgraded resourceModel
	if !resp.Diagnostics.HasError() {
		resp.Diagnostics.Append(resp.State.Get(ctx, &upgraded)...)
	}

	return upgraded, resp.Diagnostics
}

func TestUnitUpgradeStateV0toV1(t *testing.T) {
	t.Parallel()

	got, diags := upgradeFixture(t, nil)
	if diags.HasError() || diags.WarningsCount() > 0 {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}

	strs := map[string]struct{ got, want string }{
		"appliance_url":          {got.ApplianceUrl.ValueString(), "https://morpheus.test.local"},
		"internal_appliance_url": {got.InternalApplianceUrl.ValueString(), "https://pxemorpheus.test.local"},
		"api_allowed_origins":    {got.ApiAllowedOrigins.ValueString(), "demo"},
		"smtp_from_address":      {got.SmtpFromAddress.ValueString(), "testemail@test.local"},
		"smtp_server":            {got.SmtpServer.ValueString(), "smtp01.test.local"},
		"smtp_username":          {got.SmtpUsername.ValueString(), "jsmith"},
		"proxy_host":             {got.ProxyHost.ValueString(), "10.0.0.100"},
		"proxy_user":             {got.ProxyUser.ValueString(), "jsmith"},
		"proxy_domain":           {got.ProxyDomain.ValueString(), "test.local"},
		"proxy_workstation":      {got.ProxyWorkstation.ValueString(), "work123"},
		"no_proxy":               {got.NoProxy.ValueString(), "10.0.0.0/8,localhost,127.0.0.1"},
		"currency_provider":      {got.CurrencyProvider.ValueString(), "fixer"},
	}
	for name, c := range strs {
		if c.got != c.want {
			t.Errorf("%s: got %q, want %q", name, c.got, c.want)
		}
	}

	ints := map[string]struct{ got, want int64 }{
		"smtp_port":            {got.SmtpPort.ValueInt64(), 465},
		"proxy_port":           {got.ProxyPort.ValueInt64(), 3128},
		"default_role_id":      {got.DefaultRoleId.ValueInt64(), 5},
		"default_user_role_id": {got.DefaultUserRoleId.ValueInt64(), 10},
	}
	for name, c := range ints {
		if c.got != c.want {
			t.Errorf("%s: got %d, want %d", name, c.got, c.want)
		}
	}

	if !got.RegistrationEnabled.ValueBool() || got.DockerPrivilegedMode.ValueBool() ||
		!got.SmtpUseSsl.ValueBool() || !got.SmtpUseTls.ValueBool() {
		t.Errorf("booleans must carry across: registration=%v docker=%v ssl=%v tls=%v",
			got.RegistrationEnabled, got.DockerPrivilegedMode, got.SmtpUseSsl, got.SmtpUseTls)
	}

	// The SDKv2 state held the API's password hashes and the currency key in
	// clear; write-only values are never in state.
	for name, v := range map[string]interface{ IsNull() bool }{
		"smtp_password":             got.SmtpPassword,
		"proxy_password":            got.ProxyPassword,
		"currency_provider_api_key": got.CurrencyProviderApiKey,
	} {
		if !v.IsNull() {
			t.Errorf("%s must be null after the upgrade", name)
		}
	}

	// Null, so the first plan after the upgrade records the digests.
	if !got.WriteOnlyRevision.IsNull() {
		t.Error("write_only_revision must be null after the upgrade")
	}

	// The id the SDKv2 resource set is kept.
	if got.Id.ValueString() != "1" {
		t.Errorf("id: got %v, want \"1\"", got.Id)
	}

	// The settings the SDKv2 resource did not have start null; the refresh
	// after the upgrade reads them.
	for name, v := range map[string]interface{ IsNull() bool }{
		"password_min_length":              got.PasswordMinLength,
		"password_min_upper_case":          got.PasswordMinUpperCase,
		"password_min_numbers":             got.PasswordMinNumbers,
		"password_min_symbols":             got.PasswordMinSymbols,
		"user_browser_session_timeout":     got.UserBrowserSessionTimeout,
		"user_browser_session_warning":     got.UserBrowserSessionWarning,
		"expire_password_after_days":       got.ExpirePasswordAfterDays,
		"disable_user_after_attempts":      got.DisableUserAfterAttempts,
		"disable_user_after_days_inactive": got.DisableUserAfterDaysInactive,
		"warn_user_days_before_disable":    got.WarnUserDaysBeforeDisable,
		"twilio_account_sid":               got.TwilioAccountSid,
		"twilio_sms_from":                  got.TwilioSmsFrom,
		"twilio_auth_token":                got.TwilioAuthToken,
		"cloud_sync_interval_seconds":      got.CloudSyncIntervalSeconds,
		"cluster_sync_interval_seconds":    got.ClusterSyncIntervalSeconds,
		"usage_retainment_period":          got.UsageRetainmentPeriod,
		"invoice_retainment_period":        got.InvoiceRetainmentPeriod,
		"incident_retainment_period":       got.IncidentRetainmentPeriod,
		"reports_retainment_period":        got.ReportsRetainmentPeriod,
		"http_blacklist_hosts":             got.HttpBlacklistHosts,
		"http_approvelist_hosts":           got.HttpApprovelistHosts,
		"no_agent":                         got.NoAgent,
		"agent_ssl_verify":                 got.AgentSslVerify,
		"disable_ssh_password_auth":        got.DisableSshPasswordAuth,
		"default_locale":                   got.DefaultLocale,
		"default_vdi_gateway_id":           got.DefaultVdiGatewayId,
		"max_option_list_size":             got.MaxOptionListSize,
		"dashboards_to_display":            got.DashboardsToDisplay,
		"enable_murph":                     got.EnableMurph,
	} {
		if !v.IsNull() {
			t.Errorf("%s must be null after the upgrade", name)
		}
	}
}

// Released SDKv2 versions predate no_proxy: their state lacks the attribute
// altogether, and must upgrade all the same.
func TestUnitUpgradeStateV0toV1WithoutNoProxy(t *testing.T) {
	t.Parallel()

	got, diags := upgradeFixture(t, func(attrs map[string]any) {
		delete(attrs, "no_proxy")
	})
	if diags.HasError() {
		t.Fatal(diags)
	}

	if !got.NoProxy.IsNull() {
		t.Errorf("no_proxy: got %v, want null", got.NoProxy)
	}

	if got.ApplianceUrl.ValueString() != "https://morpheus.test.local" {
		t.Errorf("appliance_url: got %v", got.ApplianceUrl)
	}
}

// The SDKv2 resource read unset settings back as empty strings: empty ports
// and role ids become null; empty strings are carried.
func TestUnitUpgradeStateV0toV1UnsetValues(t *testing.T) {
	t.Parallel()

	got, diags := upgradeFixture(t, func(attrs map[string]any) {
		for _, name := range []string{
			"smtp_port", "proxy_port", "default_role_id", "default_user_role_id",
			"proxy_host", "smtp_password", "proxy_password", "currency_provider_api_key",
		} {
			attrs[name] = ""
		}
	})
	if diags.HasError() || diags.WarningsCount() > 0 {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}

	for name, v := range map[string]interface{ IsNull() bool }{
		"smtp_port":            got.SmtpPort,
		"proxy_port":           got.ProxyPort,
		"default_role_id":      got.DefaultRoleId,
		"default_user_role_id": got.DefaultUserRoleId,
	} {
		if !v.IsNull() {
			t.Errorf("%s: an empty value must upgrade to null", name)
		}
	}

	if got.ProxyHost.IsNull() || got.ProxyHost.ValueString() != "" {
		t.Errorf("proxy_host: an empty string must be carried, got %v", got.ProxyHost)
	}
}

// A port or role id that is not an integer is dropped with a warning, not an
// error: the refresh after the upgrade reads the appliance's value again.
func TestUnitUpgradeStateV0toV1NonNumeric(t *testing.T) {
	t.Parallel()

	got, diags := upgradeFixture(t, func(attrs map[string]any) {
		attrs["smtp_port"] = "smtp"
	})
	if diags.HasError() {
		t.Fatal(diags)
	}

	if diags.WarningsCount() != 1 {
		t.Errorf("want one warning, got %v", diags)
	}

	if !got.SmtpPort.IsNull() {
		t.Errorf("smtp_port: got %v, want null", got.SmtpPort)
	}

	if got.ProxyPort.ValueInt64() != 3128 {
		t.Errorf("proxy_port must be unaffected, got %v", got.ProxyPort)
	}
}

// The SDKv2 resource never set stats_retainment_period, but a state that has
// a value carries it.
func TestUnitUpgradeStateV0toV1CarriesStatsRetainmentPeriod(t *testing.T) {
	t.Parallel()

	got, diags := upgradeFixture(t, func(attrs map[string]any) {
		attrs["stats_retainment_period"] = 60
	})
	if diags.HasError() {
		t.Fatal(diags)
	}

	if got.StatsRetainmentPeriod.ValueInt64() != 60 {
		t.Errorf("stats_retainment_period: got %v, want 60", got.StatsRetainmentPeriod)
	}
}
