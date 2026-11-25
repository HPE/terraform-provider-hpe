// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package settingappliance

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

// readFixture decodes testdata/settings_get.json -- a GET response shaped as
// the appliance renders it: every setting a string (ports and role ids
// included), unset ones null or empty, password masks with their hashes, the
// currency key in clear -- through the SDK model's own decoder.
func readFixture(t *testing.T) *applianceSettings {
	t.Helper()

	raw, err := os.ReadFile("testdata/settings_get.json")
	if err != nil {
		t.Fatal(err)
	}

	var resp sdk.ListApplianceSettings200Response
	if err := json.Unmarshal(raw, &resp); err != nil {
		t.Fatalf("the settings the appliance renders must decode: %v", err)
	}

	if resp.ApplianceSettings == nil {
		t.Fatal("no applianceSettings in the fixture")
	}

	return resp.ApplianceSettings
}

// unconfiguredPlan is a create plan in which nothing is configured: every
// optional+computed attribute is unknown, every write-only one null.
func unconfiguredPlan() SettingApplianceModel {
	return SettingApplianceModel{
		AgentSslVerify:               types.BoolUnknown(),
		ApiAllowedOrigins:            types.StringUnknown(),
		ApplianceUrl:                 types.StringUnknown(),
		CloudSyncIntervalSeconds:     types.Int64Unknown(),
		ClusterSyncIntervalSeconds:   types.Int64Unknown(),
		CurrencyProvider:             types.StringUnknown(),
		CurrencyProviderApiKey:       types.StringNull(),
		DashboardsToDisplay:          types.StringUnknown(),
		DefaultLocale:                types.StringUnknown(),
		DefaultRoleId:                types.Int64Unknown(),
		DefaultUserRoleId:            types.Int64Unknown(),
		DefaultVdiGatewayId:          types.Int64Unknown(),
		DisableUserAfterAttempts:     types.Int64Unknown(),
		DisableUserAfterDaysInactive: types.Int64Unknown(),
		DisableSshPasswordAuth:       types.BoolUnknown(),
		DockerPrivilegedMode:         types.BoolUnknown(),
		EnableMurph:                  types.BoolUnknown(),
		ExpirePasswordAfterDays:      types.Int64Unknown(),
		HttpApprovelistHosts:         types.StringUnknown(),
		HttpBlacklistHosts:           types.StringUnknown(),
		IncidentRetainmentPeriod:     types.Int64Unknown(),
		InternalApplianceUrl:         types.StringUnknown(),
		InvoiceRetainmentPeriod:      types.Int64Unknown(),
		MaxOptionListSize:            types.Int64Unknown(),
		NoAgent:                      types.BoolUnknown(),
		NoProxy:                      types.StringUnknown(),
		PasswordMinLength:            types.Int64Unknown(),
		PasswordMinNumbers:           types.Int64Unknown(),
		PasswordMinSymbols:           types.Int64Unknown(),
		PasswordMinUpperCase:         types.Int64Unknown(),
		ProxyDomain:                  types.StringUnknown(),
		ProxyHost:                    types.StringUnknown(),
		ProxyPassword:                types.StringNull(),
		ProxyPort:                    types.Int64Unknown(),
		ProxyUser:                    types.StringUnknown(),
		ProxyWorkstation:             types.StringUnknown(),
		RegistrationEnabled:          types.BoolUnknown(),
		ReportsRetainmentPeriod:      types.Int64Unknown(),
		SmtpFromAddress:              types.StringUnknown(),
		SmtpPassword:                 types.StringNull(),
		SmtpPort:                     types.Int64Unknown(),
		SmtpServer:                   types.StringUnknown(),
		SmtpUseSsl:                   types.BoolUnknown(),
		SmtpUseTls:                   types.BoolUnknown(),
		SmtpUsername:                 types.StringUnknown(),
		StatsRetainmentPeriod:        types.Int64Unknown(),
		TwilioAccountSid:             types.StringUnknown(),
		TwilioAuthToken:              types.StringNull(),
		TwilioSmsFrom:                types.StringUnknown(),
		UsageRetainmentPeriod:        types.Int64Unknown(),
		UserBrowserSessionTimeout:    types.Int64Unknown(),
		UserBrowserSessionWarning:    types.Int64Unknown(),
		WarnUserDaysBeforeDisable:    types.Int64Unknown(),
		WriteOnlyRevision:            types.Int64Unknown(),
	}
}

// currentState is the state of a resource whose settings are the fixture's.
func currentState(t *testing.T) SettingApplianceModel {
	t.Helper()

	var state SettingApplianceModel
	if diags := mapSettings(&state, readFixture(t)); diags.HasError() {
		t.Fatal(diags)
	}

	state.WriteOnlyRevision = types.Int64Value(1700000000)

	return state
}

// wire is the body as it goes on the wire.
func wire(t *testing.T, body settingsUpdate) map[string]any {
	t.Helper()

	raw, err := json.Marshal(sdk.UpdateApplianceSettingsRequest{ApplianceSettings: &body})
	if err != nil {
		t.Fatal(err)
	}

	var out struct {
		ApplianceSettings map[string]any `json:"applianceSettings"`
	}

	if err := json.Unmarshal(raw, &out); err != nil {
		t.Fatal(err)
	}

	if out.ApplianceSettings == nil {
		return map[string]any{}
	}

	return out.ApplianceSettings
}

func expectWire(t *testing.T, body settingsUpdate, want map[string]any) {
	t.Helper()

	if got := wire(t, body); !reflect.DeepEqual(got, want) {
		t.Errorf("body on the wire:\n got  %v\n want %v", got, want)
	}
}

// The SDKv2 resource sent registrationEnabled and dockerPrivilegedMode on
// every create, false when unset, so a first apply that set only, say, the
// SMTP server switched both off. Only the configured settings are sent now.
func TestUnitSettingsBodyCreateSendsOnlyTheConfiguredSettings(t *testing.T) {
	t.Parallel()

	plan := unconfiguredPlan()
	plan.SmtpServer = types.StringValue("smtp01.test.local")

	body, send := settingsBody(&plan, nil)
	if !send {
		t.Fatal("a configured setting must be sent")
	}

	expectWire(t, body, map[string]any{"smtpServer": "smtp01.test.local"})
}

// A configuration that sets nothing sends nothing: Morpheus refuses an empty
// update ("No data provided").
func TestUnitSettingsBodyNothingConfigured(t *testing.T) {
	t.Parallel()

	plan := unconfiguredPlan()

	body, send := settingsBody(&plan, nil)
	if send {
		t.Errorf("nothing to send, got %v", wire(t, body))
	}
}

func TestUnitSettingsBodyUpdateSendsOnlyWhatChanged(t *testing.T) {
	t.Parallel()

	state := currentState(t)

	plan := state
	plan.ApiAllowedOrigins = types.StringValue("https://tf-acc.example.com")
	plan.SmtpUseTls = types.BoolValue(true)

	body, send := settingsBody(&plan, &state)
	if !send {
		t.Fatal("a changed setting must be sent")
	}

	expectWire(t, body, map[string]any{
		"corsAllowed": "https://tf-acc.example.com",
		"smtpTLS":     true,
	})
}

// An update that changes no setting -- a credential re-send, or the one that
// records the digests after an upgrade -- has nothing else to send.
func TestUnitSettingsBodyUpdateWithNoSettingChanged(t *testing.T) {
	t.Parallel()

	state := currentState(t)

	plan := state
	plan.WriteOnlyRevision = types.Int64Unknown()

	if body, send := settingsBody(&plan, &state); send {
		t.Errorf("nothing changed, got %v", wire(t, body))
	}
}

// Morpheus validates the default role ids it is sent whenever it is sent
// registrationEnabled=true, and nothing it is not sent: a change to any of the
// three sends all three.
func TestUnitSettingsBodyRegistrationTravelsWithTheRoleIds(t *testing.T) {
	t.Parallel()

	state := currentState(t)
	state.DefaultUserRoleId = types.Int64Value(10)

	plan := state
	plan.DefaultRoleId = types.Int64Value(6)

	body, send := settingsBody(&plan, &state)
	if !send {
		t.Fatal("a changed role id must be sent")
	}

	expectWire(t, body, map[string]any{
		"registrationEnabled": false,
		"defaultRoleId":       float64(6),
		"defaultUserRoleId":   float64(10),
	})
}

// The PUT types smtpPort as an integer and proxyPort as a string.
func TestUnitSettingsBodyPortTypes(t *testing.T) {
	t.Parallel()

	plan := unconfiguredPlan()
	plan.SmtpPort = types.Int64Value(465)
	plan.ProxyPort = types.Int64Value(3128)

	body, send := settingsBody(&plan, nil)
	if !send {
		t.Fatal("configured ports must be sent")
	}

	expectWire(t, body, map[string]any{
		"smtpPort":  float64(465),
		"proxyPort": "3128",
	})
}

// A configured empty string is a value like any other, and is sent.
func TestUnitSettingsBodySendsAnEmptyString(t *testing.T) {
	t.Parallel()

	state := currentState(t)
	state.ProxyHost = types.StringValue("10.0.0.100")

	plan := state
	plan.ProxyHost = types.StringValue("")

	body, send := settingsBody(&plan, &state)
	if !send {
		t.Fatal("a configured empty string must be sent")
	}

	expectWire(t, body, map[string]any{"proxyHost": ""})
}

// The PUT types the password-policy and session values as strings, as it
// does proxyPort, and the other numeric settings as integers. A 0 is a value
// (it disables a policy) and is sent.
func TestUnitSettingsBodyNumericSettings(t *testing.T) {
	t.Parallel()

	plan := unconfiguredPlan()
	plan.PasswordMinLength = types.Int64Value(12)
	plan.PasswordMinUpperCase = types.Int64Value(2)
	plan.PasswordMinNumbers = types.Int64Value(3)
	plan.PasswordMinSymbols = types.Int64Value(0)
	plan.UserBrowserSessionTimeout = types.Int64Value(30)
	plan.UserBrowserSessionWarning = types.Int64Value(5)
	plan.ExpirePasswordAfterDays = types.Int64Value(90)
	plan.DisableUserAfterAttempts = types.Int64Value(5)
	plan.DisableUserAfterDaysInactive = types.Int64Value(0)
	plan.WarnUserDaysBeforeDisable = types.Int64Value(7)
	plan.CloudSyncIntervalSeconds = types.Int64Value(600)
	plan.ClusterSyncIntervalSeconds = types.Int64Value(120)
	plan.UsageRetainmentPeriod = types.Int64Value(180)
	plan.InvoiceRetainmentPeriod = types.Int64Value(365)
	plan.IncidentRetainmentPeriod = types.Int64Value(30)
	plan.ReportsRetainmentPeriod = types.Int64Value(0)
	plan.DefaultVdiGatewayId = types.Int64Value(3)
	plan.MaxOptionListSize = types.Int64Value(2000)

	body, send := settingsBody(&plan, nil)
	if !send {
		t.Fatal("configured settings must be sent")
	}

	expectWire(t, body, map[string]any{
		"passwordMinLength":          "12",
		"passwordMinUpperCase":       "2",
		"passwordMinNumbers":         "3",
		"passwordMinSymbols":         "0",
		"userBrowserSessionTimeout":  "30",
		"userBrowserSessionWarning":  "5",
		"expirePwdDays":              float64(90),
		"disableAfterAttempts":       float64(5),
		"disableAfterDaysInactive":   float64(0),
		"warnUserDaysBefore":         float64(7),
		"cloudSyncIntervalSeconds":   float64(600),
		"clusterSyncIntervalSeconds": float64(120),
		"usageRetainmentPeriod":      float64(180),
		"invoiceRetainmentPeriod":    float64(365),
		"incidentRetainmentPeriod":   float64(30),
		"reportsRetainmentPeriod":    float64(0),
		"defaultVdiGateway":          float64(3),
		"maxOptionListSize":          float64(2000),
	})
}

// The flags travel as booleans, false included, and the text settings as
// strings, an empty one included.
func TestUnitSettingsBodyFlagsAndText(t *testing.T) {
	t.Parallel()

	plan := unconfiguredPlan()
	plan.NoAgent = types.BoolValue(false)
	plan.AgentSslVerify = types.BoolValue(true)
	plan.DisableSshPasswordAuth = types.BoolValue(true)
	plan.EnableMurph = types.BoolValue(false)
	plan.TwilioAccountSid = types.StringValue("AC0123")
	plan.TwilioSmsFrom = types.StringValue("+15550100")
	plan.HttpBlacklistHosts = types.StringValue("169.254.169.254,blocked.test.local")
	plan.HttpApprovelistHosts = types.StringValue("")
	plan.DefaultLocale = types.StringValue("en")
	plan.DashboardsToDisplay = types.StringValue("1,2")

	body, send := settingsBody(&plan, nil)
	if !send {
		t.Fatal("configured settings must be sent")
	}

	expectWire(t, body, map[string]any{
		"noAgent":                false,
		"agentSSLVerify":         true,
		"disableSSHPasswordAuth": true,
		"enableMurph":            false,
		"twilioAccountSid":       "AC0123",
		"twilioSmsFrom":          "+15550100",
		"httpBlacklistHosts":     "169.254.169.254,blocked.test.local",
		"httpApprovelistHosts":   "",
		"defaultLocale":          "en",
		"dashboardsToDisplay":    "1,2",
	})
}

// On an update, a setting is sent only when its planned value differs from
// the state: the fixture's password length is 8 and its cloud sync interval
// 300.
func TestUnitSettingsBodyUpdateSendsOnlyTheChangedSettings(t *testing.T) {
	t.Parallel()

	state := currentState(t)

	plan := state
	plan.PasswordMinLength = types.Int64Value(12)
	plan.CloudSyncIntervalSeconds = types.Int64Value(300)
	plan.EnableMurph = types.BoolValue(true)

	body, send := settingsBody(&plan, &state)
	if !send {
		t.Fatal("changed settings must be sent")
	}

	expectWire(t, body, map[string]any{
		"passwordMinLength": "12",
		"enableMurph":       true,
	})
}

func TestUnitAddCredentials(t *testing.T) {
	t.Parallel()

	config := unconfiguredPlan()
	config.SmtpPassword = types.StringValue("smtp-secret")
	config.CurrencyProviderApiKey = types.StringValue("currency-key")
	config.TwilioAuthToken = types.StringValue("twilio-token")

	var body settingsUpdate
	if !addCredentials(&body, &config) {
		t.Fatal("configured credentials must be added")
	}

	// The proxy password is not configured, so the appliance keeps its own.
	expectWire(t, body, map[string]any{ //nolint:gosec // G101: test fixture, not a real credential
		"smtpPassword":    "smtp-secret",
		"currencyKey":     "currency-key",
		"twilioAuthToken": "twilio-token",
	})

	none := unconfiguredPlan()

	var empty settingsUpdate
	if addCredentials(&empty, &none) {
		t.Errorf("no credential is configured, got %v", wire(t, empty))
	}
}

func TestUnitMapSettings(t *testing.T) {
	t.Parallel()

	var got SettingApplianceModel

	diags := mapSettings(&got, readFixture(t))
	if diags.HasError() || diags.WarningsCount() > 0 {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}

	want := SettingApplianceModel{
		ApplianceUrl:          types.StringValue("https://morpheus.test.local"),
		InternalApplianceUrl:  types.StringNull(),
		ApiAllowedOrigins:     types.StringValue(""),
		RegistrationEnabled:   types.BoolValue(false),
		DefaultRoleId:         types.Int64Value(5),
		DefaultUserRoleId:     types.Int64Null(),
		DockerPrivilegedMode:  types.BoolValue(true),
		SmtpFromAddress:       types.StringValue("info@gomorpheus.com"),
		SmtpServer:            types.StringValue("smtp01.test.local"),
		SmtpPort:              types.Int64Value(465),
		SmtpUseSsl:            types.BoolValue(true),
		SmtpUseTls:            types.BoolValue(false),
		SmtpUsername:          types.StringValue("jsmith"),
		ProxyHost:             types.StringValue(""),
		ProxyPort:             types.Int64Null(),
		ProxyUser:             types.StringNull(),
		ProxyDomain:           types.StringNull(),
		ProxyWorkstation:      types.StringNull(),
		NoProxy:               types.StringValue("localhost,127.0.0.1"),
		CurrencyProvider:      types.StringValue("fixer"),
		StatsRetainmentPeriod: types.Int64Value(30),
		// The numeric settings are strings on the wire.
		PasswordMinLength:            types.Int64Value(8),
		PasswordMinUpperCase:         types.Int64Value(1),
		PasswordMinNumbers:           types.Int64Value(1),
		PasswordMinSymbols:           types.Int64Value(1),
		UserBrowserSessionTimeout:    types.Int64Value(20),
		UserBrowserSessionWarning:    types.Int64Value(15),
		ExpirePasswordAfterDays:      types.Int64Value(90),
		DisableUserAfterAttempts:     types.Int64Value(0),
		DisableUserAfterDaysInactive: types.Int64Null(),
		WarnUserDaysBeforeDisable:    types.Int64Null(),
		CloudSyncIntervalSeconds:     types.Int64Value(300),
		ClusterSyncIntervalSeconds:   types.Int64Null(),
		UsageRetainmentPeriod:        types.Int64Null(),
		InvoiceRetainmentPeriod:      types.Int64Null(),
		IncidentRetainmentPeriod:     types.Int64Null(),
		ReportsRetainmentPeriod:      types.Int64Value(90),
		DefaultVdiGatewayId:          types.Int64Value(3),
		MaxOptionListSize:            types.Int64Value(1000),
		HttpBlacklistHosts:           types.StringValue("169.254.169.254,blocked.test.local"),
		HttpApprovelistHosts:         types.StringNull(),
		NoAgent:                      types.BoolValue(false),
		AgentSslVerify:               types.BoolValue(true),
		DisableSshPasswordAuth:       types.BoolValue(false),
		DefaultLocale:                types.StringValue("en"),
		DashboardsToDisplay:          types.StringNull(),
		EnableMurph:                  types.BoolValue(false),
		TwilioAccountSid:             types.StringValue("AC0123456789abcdef0123456789abcdef"),
		TwilioSmsFrom:                types.StringValue("+15550100"),
		// The masks, the hashes and the clear currency key are never read
		// into the write-only attributes.
		SmtpPassword:           types.StringNull(),
		ProxyPassword:          types.StringNull(),
		CurrencyProviderApiKey: types.StringNull(),
		TwilioAuthToken:        types.StringNull(),
		// Not mapped: set by Create and Update only.
		WriteOnlyRevision: types.Int64{},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("mapped settings:\n got  %+v\n want %+v", got, want)
	}
}

// A port that is not an integer, set so outside Terraform, is read as unset
// with a warning rather than making the settings unreadable.
func TestUnitMapSettingsNonNumeric(t *testing.T) {
	t.Parallel()

	settings := readFixture(t)
	port := "smtp"
	settings.SmtpPort.Set(&port)

	var got SettingApplianceModel

	diags := mapSettings(&got, settings)
	if diags.HasError() || diags.WarningsCount() != 1 {
		t.Fatalf("want exactly one warning, got %v", diags)
	}

	if !got.SmtpPort.IsNull() {
		t.Errorf("smtp_port: got %v, want null", got.SmtpPort)
	}
}

// Morpheus stores a setting emptied with "" as none and reports it null: the ""
// the model holds is kept, so "" clears a setting without an inconsistent
// result. A setting the appliance does report, or one the model holds another
// value for, is read as the appliance reports it.
func TestUnitMapSettingsKeepsAnEmptiedSetting(t *testing.T) {
	t.Parallel()

	got := SettingApplianceModel{
		ProxyUser:            types.StringValue(""),       // reported null
		HttpApprovelistHosts: types.StringValue(""),       // reported null
		ApiAllowedOrigins:    types.StringValue(""),       // reported ""
		TwilioSmsFrom:        types.StringValue(""),       // reported "+15550100"
		ProxyDomain:          types.StringValue("corp"),   // reported null
		DashboardsToDisplay:  types.StringUnknown(),       // reported null
		HttpBlacklistHosts:   types.StringValue("a.test"), // reported a list
	}

	if diags := mapSettings(&got, readFixture(t)); diags.HasError() {
		t.Fatal(diags)
	}

	for name, c := range map[string]struct{ got, want types.String }{
		"proxy_user":             {got.ProxyUser, types.StringValue("")},
		"http_approvelist_hosts": {got.HttpApprovelistHosts, types.StringValue("")},
		"api_allowed_origins":    {got.ApiAllowedOrigins, types.StringValue("")},
		"twilio_sms_from":        {got.TwilioSmsFrom, types.StringValue("+15550100")},
		"proxy_domain":           {got.ProxyDomain, types.StringNull()},
		"dashboards_to_display":  {got.DashboardsToDisplay, types.StringNull()},
		"http_blacklist_hosts":   {got.HttpBlacklistHosts, types.StringValue("169.254.169.254,blocked.test.local")},
	} {
		if !c.got.Equal(c.want) {
			t.Errorf("%s: got %v, want %v", name, c.got, c.want)
		}
	}
}

func TestUnitSettingsBodyStatsRetainmentPeriod(t *testing.T) {
	t.Parallel()

	state := currentState(t)

	plan := state
	plan.StatsRetainmentPeriod = types.Int64Value(60)

	body, send := settingsBody(&plan, &state)
	if !send {
		t.Fatal("a changed stats retention must be sent")
	}

	expectWire(t, body, map[string]any{"statsRetainmentPeriod": float64(60)})
}

// The appliance renders the stats retention as a string; the SDK may hand it
// over as either variant.
func TestUnitStatsRetainmentPeriod(t *testing.T) {
	t.Parallel()

	type period = sdk.ListApplianceSettings200ResponseApplianceSettingsStatsRetainmentPeriod

	thirty, sixty, monthly := int32(30), "60", "monthly"

	cases := []struct {
		name     string
		value    *period
		want     types.Int64
		warnings int
	}{
		{"absent", nil, types.Int64Null(), 0},
		{"an integer", &period{Int32: &thirty}, types.Int64Value(30), 0},
		{"a string", &period{String: &sixty}, types.Int64Value(60), 0},
		{"not a number", &period{String: &monthly}, types.Int64Null(), 1},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var diags diag.Diagnostics

			if got := statsRetainmentPeriod(tc.value, &diags); !got.Equal(tc.want) {
				t.Errorf("got %v, want %v", got, tc.want)
			}

			if diags.HasError() || diags.WarningsCount() != tc.warnings {
				t.Errorf("diagnostics: %v", diags)
			}
		})
	}
}
