// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package settingappliance

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

// applianceSettings is the settings object the GET returns; settingsUpdate is
// the one the PUT takes.
type (
	applianceSettings = sdk.ListApplianceSettings200ResponseApplianceSettings
	settingsUpdate    = sdk.UpdateApplianceSettingsRequestApplianceSettings
)

// getSettings reads the appliance settings. The error is formatted for a
// diagnostic already; the response is returned for its status code.
func getSettings(ctx context.Context, client *sdk.APIClient) (*applianceSettings, *http.Response, error) {
	result, httpResp, err := client.ApplianceSettingsAPI.ListApplianceSettings(ctx).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		return nil, httpResp, err
	}

	if result == nil || result.ApplianceSettings == nil {
		return nil, httpResp, errors.New("the response carries no applianceSettings object")
	}

	return result.ApplianceSettings, httpResp, nil
}

// updateSettings sends the settings, and the answer must confirm them.
// Morpheus answers an update it refuses -- a default role that does not
// exist, a setting it does not know -- with HTTP 200 and {"success": false,
// "errors": {...}}, so the body decides, as it did for the SDKv2 resource. It
// renders the success flag on every answer, so an answer without it -- no
// body, or a body that is not Morpheus's, such as a proxy's JSON page --
// confirms nothing and is an error too: Create and Update must not record the
// digests of credentials the appliance may not have taken. (The SDKv2
// resource failed on a missing body but took a body without the flag as
// success.) The error is formatted for a diagnostic already.
func updateSettings(ctx context.Context, client *sdk.APIClient, body settingsUpdate) error {
	result, httpResp, err := client.ApplianceSettingsAPI.UpdateApplianceSettings(ctx).
		UpdateApplianceSettingsRequest(sdk.UpdateApplianceSettingsRequest{ApplianceSettings: &body}).
		Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		return err
	}

	// The generated client leaves the result nil when the body is empty.
	if result == nil {
		return errors.New("the appliance answered the update with an empty body, " +
			"so the settings cannot be confirmed as applied")
	}

	if result.Success == nil {
		return errors.New("the appliance's answer to the update has no success flag, " +
			"so the settings cannot be confirmed as applied")
	}

	if !*result.Success {
		return fmt.Errorf("the appliance refused the settings: %s",
			errfmt.RefusalMessage(result.AdditionalProperties))
	}

	return nil
}

// known reports whether a value is neither null nor unknown.
func known(v attr.Value) bool {
	return !v.IsNull() && !v.IsUnknown()
}

// settingsBody builds the PUT body from the plan: on create (state nil) every
// setting the plan knows -- the configured ones, since an unconfigured
// optional+computed setting is unknown on create -- and on update the ones
// whose planned value differs from the state. It never carries the write-only
// credentials (see addCredentials) or exchangeUrl (see the package
// documentation), and reports whether it carries anything: Morpheus refuses an
// empty update ("No data provided").
//
// The registration flag and the two default role ids travel together.
// Morpheus validates the role ids it is sent whenever it is sent
// registrationEnabled=true, and nothing it is not sent, so a change to any of
// the three sends all three that are known.
func settingsBody(plan, state *SettingApplianceModel) (settingsUpdate, bool) {
	var body settingsUpdate

	create := state == nil
	if create {
		state = &SettingApplianceModel{}
	}

	send := false

	// pick reports whether a planned value is to be sent: known and, on an
	// update, different from the state.
	pick := func(p, s attr.Value) bool {
		if !known(p) || (!create && p.Equal(s)) {
			return false
		}

		send = true

		return true
	}

	str := func(p, s types.String) *string {
		if !pick(p, s) {
			return nil
		}

		return p.ValueStringPointer()
	}

	nullable := func(dst *sdk.NullableString, p, s types.String) {
		if v := str(p, s); v != nil {
			dst.Set(v)
		}
	}

	boolean := func(p, s types.Bool) *bool {
		if !pick(p, s) {
			return nil
		}

		return p.ValueBoolPointer()
	}

	integer := func(p, s types.Int64) *int64 {
		if !pick(p, s) {
			return nil
		}

		return p.ValueInt64Pointer()
	}

	// numeral is integer for the settings the PUT types as strings
	// (proxyPort, the password-policy and session values); the appliance
	// stores every setting as one.
	numeral := func(p, s types.Int64) *string {
		v := integer(p, s)
		if v == nil {
			return nil
		}

		n := strconv.FormatInt(*v, 10)

		return &n
	}

	body.ApplianceUrl = str(plan.ApplianceUrl, state.ApplianceUrl)
	nullable(&body.InternalApplianceUrl, plan.InternalApplianceUrl, state.InternalApplianceUrl)
	nullable(&body.CorsAllowed, plan.ApiAllowedOrigins, state.ApiAllowedOrigins)
	body.DockerPrivilegedMode = boolean(plan.DockerPrivilegedMode, state.DockerPrivilegedMode)
	body.PasswordMinLength = numeral(plan.PasswordMinLength, state.PasswordMinLength)
	body.PasswordMinUpperCase = numeral(plan.PasswordMinUpperCase, state.PasswordMinUpperCase)
	body.PasswordMinNumbers = numeral(plan.PasswordMinNumbers, state.PasswordMinNumbers)
	body.PasswordMinSymbols = numeral(plan.PasswordMinSymbols, state.PasswordMinSymbols)
	body.UserBrowserSessionTimeout = numeral(plan.UserBrowserSessionTimeout, state.UserBrowserSessionTimeout)
	body.UserBrowserSessionWarning = numeral(plan.UserBrowserSessionWarning, state.UserBrowserSessionWarning)
	body.ExpirePwdDays = integer(plan.ExpirePasswordAfterDays, state.ExpirePasswordAfterDays)
	body.DisableAfterAttempts = integer(plan.DisableUserAfterAttempts, state.DisableUserAfterAttempts)
	body.DisableAfterDaysInactive = integer(plan.DisableUserAfterDaysInactive, state.DisableUserAfterDaysInactive)
	body.WarnUserDaysBefore = integer(plan.WarnUserDaysBeforeDisable, state.WarnUserDaysBeforeDisable)
	body.SmtpMailFrom = str(plan.SmtpFromAddress, state.SmtpFromAddress)
	body.SmtpServer = str(plan.SmtpServer, state.SmtpServer)
	body.SmtpPort = integer(plan.SmtpPort, state.SmtpPort)
	body.SmtpSSL = boolean(plan.SmtpUseSsl, state.SmtpUseSsl)
	body.SmtpTLS = boolean(plan.SmtpUseTls, state.SmtpUseTls)
	body.SmtpUser = str(plan.SmtpUsername, state.SmtpUsername)
	nullable(&body.ProxyHost, plan.ProxyHost, state.ProxyHost)

	if v := numeral(plan.ProxyPort, state.ProxyPort); v != nil {
		body.ProxyPort.Set(v)
	}

	body.ProxyUser = str(plan.ProxyUser, state.ProxyUser)
	nullable(&body.ProxyDomain, plan.ProxyDomain, state.ProxyDomain)
	nullable(&body.ProxyWorkstation, plan.ProxyWorkstation, state.ProxyWorkstation)
	nullable(&body.NoProxy, plan.NoProxy, state.NoProxy)
	body.CurrencyProvider = str(plan.CurrencyProvider, state.CurrencyProvider)
	body.TwilioAccountSid = str(plan.TwilioAccountSid, state.TwilioAccountSid)
	body.TwilioSmsFrom = str(plan.TwilioSmsFrom, state.TwilioSmsFrom)
	body.CloudSyncIntervalSeconds = integer(plan.CloudSyncIntervalSeconds, state.CloudSyncIntervalSeconds)
	body.ClusterSyncIntervalSeconds = integer(plan.ClusterSyncIntervalSeconds, state.ClusterSyncIntervalSeconds)
	body.UsageRetainmentPeriod = integer(plan.UsageRetainmentPeriod, state.UsageRetainmentPeriod)
	body.InvoiceRetainmentPeriod = integer(plan.InvoiceRetainmentPeriod, state.InvoiceRetainmentPeriod)
	body.IncidentRetainmentPeriod = integer(plan.IncidentRetainmentPeriod, state.IncidentRetainmentPeriod)
	body.StatsRetainmentPeriod = integer(plan.StatsRetainmentPeriod, state.StatsRetainmentPeriod)
	body.ReportsRetainmentPeriod = integer(plan.ReportsRetainmentPeriod, state.ReportsRetainmentPeriod)
	body.HttpBlacklistHosts = str(plan.HttpBlacklistHosts, state.HttpBlacklistHosts)
	body.HttpApprovelistHosts = str(plan.HttpApprovelistHosts, state.HttpApprovelistHosts)
	body.NoAgent = boolean(plan.NoAgent, state.NoAgent)
	body.AgentSSLVerify = boolean(plan.AgentSslVerify, state.AgentSslVerify)
	body.DisableSSHPasswordAuth = boolean(plan.DisableSshPasswordAuth, state.DisableSshPasswordAuth)
	body.DefaultLocale = str(plan.DefaultLocale, state.DefaultLocale)
	body.DefaultVdiGateway = integer(plan.DefaultVdiGatewayId, state.DefaultVdiGatewayId)
	body.MaxOptionListSize = integer(plan.MaxOptionListSize, state.MaxOptionListSize)
	nullable(&body.DashboardsToDisplay, plan.DashboardsToDisplay, state.DashboardsToDisplay)
	body.EnableMurph = boolean(plan.EnableMurph, state.EnableMurph)

	if create ||
		!plan.RegistrationEnabled.Equal(state.RegistrationEnabled) ||
		!plan.DefaultRoleId.Equal(state.DefaultRoleId) ||
		!plan.DefaultUserRoleId.Equal(state.DefaultUserRoleId) {
		if known(plan.RegistrationEnabled) {
			body.RegistrationEnabled = plan.RegistrationEnabled.ValueBoolPointer()
			send = true
		}

		if known(plan.DefaultRoleId) {
			body.DefaultRoleId = plan.DefaultRoleId.ValueInt64Pointer()
			send = true
		}

		if known(plan.DefaultUserRoleId) {
			body.DefaultUserRoleId = plan.DefaultUserRoleId.ValueInt64Pointer()
			send = true
		}
	}

	return body, send
}

// addCredentials adds the configured write-only credentials to the body and
// reports whether it added any. A credential the configuration does not set is
// not sent, and the appliance keeps the one it holds.
func addCredentials(body *settingsUpdate, config *SettingApplianceModel) bool {
	added := false

	if known(config.SmtpPassword) {
		body.SmtpPassword = config.SmtpPassword.ValueStringPointer()
		added = true
	}

	if known(config.ProxyPassword) {
		body.ProxyPassword = config.ProxyPassword.ValueStringPointer()
		added = true
	}

	if known(config.CurrencyProviderApiKey) {
		body.CurrencyKey.Set(config.CurrencyProviderApiKey.ValueStringPointer())
		added = true
	}

	if known(config.TwilioAuthToken) {
		body.TwilioAuthToken = config.TwilioAuthToken.ValueStringPointer()
		added = true
	}

	return added
}

// mapSettings copies the appliance settings into the model. The write-only
// credentials are never read back: the API returns a mask and a hash for the
// passwords and a mask for the Twilio token, and the currency key, which it
// returns in clear, is write-only all the same. The numeric settings are
// strings on the wire and integers here; see int64Setting and
// statsRetainmentPeriod.
//
// Morpheus stores a setting emptied with "" as none (its EncryptedString type
// writes an empty value as NULL) and reports it null. So a text setting the
// API reports null or empty keeps the "" the model holds -- the configured
// value on create and update, the prior state on a refresh -- and "" clears a
// setting without the apply reporting an inconsistent result.
func mapSettings(m *SettingApplianceModel, s *applianceSettings) diag.Diagnostics {
	var diags diag.Diagnostics

	// setting reads one of the numeric settings into an integer.
	setting := func(name string, v sdk.NullableString) types.Int64 {
		return int64Setting(name, v.Get(), &diags)
	}

	// text reads a text setting over the value the model holds; see above.
	text := func(cur types.String, v *string) types.String {
		if (v == nil || *v == "") && known(cur) && cur.ValueString() == "" {
			return cur
		}

		return convert.StrToType(v)
	}

	m.ApplianceUrl = text(m.ApplianceUrl, s.ApplianceUrl)
	m.InternalApplianceUrl = text(m.InternalApplianceUrl, s.InternalApplianceUrl)
	m.ApiAllowedOrigins = text(m.ApiAllowedOrigins, s.CorsAllowed)
	m.RegistrationEnabled = convert.BoolToType(s.RegistrationEnabled)
	m.DefaultRoleId = int64Setting("default_role_id", s.DefaultRoleId, &diags)
	m.DefaultUserRoleId = setting("default_user_role_id", s.DefaultUserRoleId)
	m.DockerPrivilegedMode = convert.BoolToType(s.DockerPrivilegedMode)
	m.PasswordMinLength = setting("password_min_length", s.PasswordMinLength)
	m.PasswordMinUpperCase = setting("password_min_upper_case", s.PasswordMinUpperCase)
	m.PasswordMinNumbers = setting("password_min_numbers", s.PasswordMinNumbers)
	m.PasswordMinSymbols = setting("password_min_symbols", s.PasswordMinSymbols)
	m.UserBrowserSessionTimeout = setting("user_browser_session_timeout", s.UserBrowserSessionTimeout)
	m.UserBrowserSessionWarning = setting("user_browser_session_warning", s.UserBrowserSessionWarning)
	m.ExpirePasswordAfterDays = setting("expire_password_after_days", s.ExpirePwdDays)
	m.DisableUserAfterAttempts = setting("disable_user_after_attempts", s.DisableAfterAttempts)
	m.DisableUserAfterDaysInactive = setting("disable_user_after_days_inactive", s.DisableAfterDaysInactive)
	m.WarnUserDaysBeforeDisable = setting("warn_user_days_before_disable", s.WarnUserDaysBefore)
	m.SmtpFromAddress = text(m.SmtpFromAddress, s.SmtpMailFrom.Get())
	m.SmtpServer = text(m.SmtpServer, s.SmtpServer.Get())
	m.SmtpPort = setting("smtp_port", s.SmtpPort)
	m.SmtpUseSsl = convert.BoolToType(s.SmtpSSL)
	m.SmtpUseTls = convert.BoolToType(s.SmtpTLS)
	m.SmtpUsername = text(m.SmtpUsername, s.SmtpUser.Get())
	m.ProxyHost = text(m.ProxyHost, s.ProxyHost.Get())
	m.ProxyPort = setting("proxy_port", s.ProxyPort)
	m.ProxyUser = text(m.ProxyUser, s.ProxyUser.Get())
	m.ProxyDomain = text(m.ProxyDomain, s.ProxyDomain.Get())
	m.ProxyWorkstation = text(m.ProxyWorkstation, s.ProxyWorkstation.Get())
	m.NoProxy = text(m.NoProxy, s.NoProxy.Get())
	m.CurrencyProvider = text(m.CurrencyProvider, s.CurrencyProvider.Get())
	m.TwilioAccountSid = text(m.TwilioAccountSid, s.TwilioAccountSid.Get())
	m.TwilioSmsFrom = text(m.TwilioSmsFrom, s.TwilioSmsFrom.Get())
	m.CloudSyncIntervalSeconds = setting("cloud_sync_interval_seconds", s.CloudSyncIntervalSeconds)
	m.ClusterSyncIntervalSeconds = setting("cluster_sync_interval_seconds", s.ClusterSyncIntervalSeconds)
	m.UsageRetainmentPeriod = setting("usage_retainment_period", s.UsageRetainmentPeriod)
	m.InvoiceRetainmentPeriod = setting("invoice_retainment_period", s.InvoiceRetainmentPeriod)
	m.IncidentRetainmentPeriod = setting("incident_retainment_period", s.IncidentRetainmentPeriod)
	m.StatsRetainmentPeriod = statsRetainmentPeriod(s.StatsRetainmentPeriod.Get(), &diags)
	m.ReportsRetainmentPeriod = setting("reports_retainment_period", s.ReportsRetainmentPeriod)
	m.HttpBlacklistHosts = text(m.HttpBlacklistHosts, s.HttpBlacklistHosts.Get())
	m.HttpApprovelistHosts = text(m.HttpApprovelistHosts, s.HttpApprovelistHosts.Get())
	m.NoAgent = convert.BoolToType(s.NoAgent)
	m.AgentSslVerify = convert.BoolToType(s.AgentSSLVerify)
	m.DisableSshPasswordAuth = convert.BoolToType(s.DisableSSHPasswordAuth)
	m.DefaultLocale = text(m.DefaultLocale, s.DefaultLocale.Get())
	m.DefaultVdiGatewayId = setting("default_vdi_gateway_id", s.DefaultVdiGateway)
	m.MaxOptionListSize = setting("max_option_list_size", s.MaxOptionListSize)
	m.DashboardsToDisplay = text(m.DashboardsToDisplay, s.DashboardsToDisplay.Get())
	m.EnableMurph = convert.BoolToType(s.EnableMurph)

	// Write-only: never in state.
	m.SmtpPassword = types.StringNull()
	m.ProxyPassword = types.StringNull()
	m.CurrencyProviderApiKey = types.StringNull()
	m.TwilioAuthToken = types.StringNull()

	return diags
}

// int64Setting reads a numeric setting -- a port, a role id, a policy value,
// an interval or a retention period -- which the appliance stores as a
// string. An absent or empty value is unset. One that is not an integer -- set
// to something else outside Terraform -- is reported as a warning and read as
// unset, so it does not make the settings unreadable; a configured value then
// plans an update that puts it right.
func int64Setting(name string, v *string, diags *diag.Diagnostics) types.Int64 {
	if v == nil || strings.TrimSpace(*v) == "" {
		return types.Int64Null()
	}

	n, err := strconv.ParseInt(strings.TrimSpace(*v), 10, 64)
	if err != nil {
		diags.AddAttributeWarning(
			path.Root(name),
			"Appliance setting is not an integer",
			fmt.Sprintf("The appliance reports %s as %q, which is not an integer; it is read as unset.", name, *v),
		)

		return types.Int64Null()
	}

	return types.Int64Value(n)
}

// statsRetainmentPeriod reads the stats retention. The appliance renders it as
// a string ("30" when none is set), which the SDK decodes as an integer or a
// string, depending on what it finds.
func statsRetainmentPeriod(
	v *sdk.ListApplianceSettings200ResponseApplianceSettingsStatsRetainmentPeriod,
	diags *diag.Diagnostics,
) types.Int64 {
	switch {
	case v == nil:
		return types.Int64Null()
	case v.Int32 != nil:
		return types.Int64Value(int64(*v.Int32))
	default:
		return int64Setting("stats_retainment_period", v.String, diags)
	}
}
