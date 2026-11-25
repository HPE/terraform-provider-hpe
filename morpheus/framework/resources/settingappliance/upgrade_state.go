// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package settingappliance

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// settingApplianceModelV0 is the state written by the retired SDKv2 resource
// (schema version 0). It differs from the current model in that it has string
// ports and role ids, the password attributes held the API's password hashes
// (the SDKv2 resource read them back into state), the currency key was kept in
// clear, and stats_retainment_period is always null: the SDKv2 resource never
// read or sent it.
type settingApplianceModelV0 struct {
	Id                     types.String `tfsdk:"id"`
	ApplianceUrl           types.String `tfsdk:"appliance_url"`
	InternalApplianceUrl   types.String `tfsdk:"internal_appliance_url"`
	ApiAllowedOrigins      types.String `tfsdk:"api_allowed_origins"`
	StatsRetainmentPeriod  types.Int64  `tfsdk:"stats_retainment_period"`
	RegistrationEnabled    types.Bool   `tfsdk:"registration_enabled"`
	DefaultRoleId          types.String `tfsdk:"default_role_id"`
	DefaultUserRoleId      types.String `tfsdk:"default_user_role_id"`
	DockerPrivilegedMode   types.Bool   `tfsdk:"docker_privileged_mode"`
	SMTPFromAddress        types.String `tfsdk:"smtp_from_address"`
	SMTPServer             types.String `tfsdk:"smtp_server"`
	SMTPPort               types.String `tfsdk:"smtp_port"`
	SMTPUseSSL             types.Bool   `tfsdk:"smtp_use_ssl"`
	SMTPUseTLS             types.Bool   `tfsdk:"smtp_use_tls"`
	SMTPUsername           types.String `tfsdk:"smtp_username"`
	SMTPPassword           types.String `tfsdk:"smtp_password"`
	ProxyHost              types.String `tfsdk:"proxy_host"`
	ProxyPort              types.String `tfsdk:"proxy_port"`
	ProxyUser              types.String `tfsdk:"proxy_user"`
	ProxyPassword          types.String `tfsdk:"proxy_password"`
	ProxyDomain            types.String `tfsdk:"proxy_domain"`
	ProxyWorkstation       types.String `tfsdk:"proxy_workstation"`
	NoProxy                types.String `tfsdk:"no_proxy"`
	CurrencyProvider       types.String `tfsdk:"currency_provider"`
	CurrencyProviderApiKey types.String `tfsdk:"currency_provider_api_key"`
}

// settingApplianceSchemaV0 describes the SDKv2 state so the framework can
// decode it during an upgrade. Only the attribute types matter for decoding;
// the flags mirror the SDKv2 declaration. It is pinned by
// testdata/state_v0.json: the unit test fails if the two disagree.
func settingApplianceSchemaV0() schema.Schema {
	optionalComputed := func() schema.StringAttribute {
		return schema.StringAttribute{Optional: true, Computed: true}
	}

	return schema.Schema{
		Version: 0,
		Attributes: map[string]schema.Attribute{
			"id":                        schema.StringAttribute{Computed: true},
			"appliance_url":             optionalComputed(),
			"internal_appliance_url":    optionalComputed(),
			"api_allowed_origins":       optionalComputed(),
			"stats_retainment_period":   schema.Int64Attribute{Optional: true, Computed: true},
			"registration_enabled":      schema.BoolAttribute{Optional: true, Computed: true},
			"default_role_id":           optionalComputed(),
			"default_user_role_id":      optionalComputed(),
			"docker_privileged_mode":    schema.BoolAttribute{Optional: true, Computed: true},
			"smtp_from_address":         optionalComputed(),
			"smtp_server":               optionalComputed(),
			"smtp_port":                 optionalComputed(),
			"smtp_use_ssl":              schema.BoolAttribute{Optional: true, Computed: true},
			"smtp_use_tls":              schema.BoolAttribute{Optional: true, Computed: true},
			"smtp_username":             optionalComputed(),
			"smtp_password":             schema.StringAttribute{Optional: true, Computed: true, Sensitive: true},
			"proxy_host":                optionalComputed(),
			"proxy_port":                optionalComputed(),
			"proxy_user":                optionalComputed(),
			"proxy_password":            schema.StringAttribute{Optional: true, Computed: true, Sensitive: true},
			"proxy_domain":              optionalComputed(),
			"proxy_workstation":         optionalComputed(),
			"no_proxy":                  optionalComputed(),
			"currency_provider":         optionalComputed(),
			"currency_provider_api_key": optionalComputed(),
		},
	}
}

func (r *Resource) UpgradeState(_ context.Context) map[int64]resource.StateUpgrader {
	priorSchema := settingApplianceSchemaV0()

	return map[int64]resource.StateUpgrader{
		0: {
			PriorSchema:   &priorSchema,
			StateUpgrader: upgradeStateV0toV1,
		},
	}
}

// upgradeStateV0toV1 converts the SDKv2 state to the current model. It keeps
// id ("1"), parses the ports and role ids, carries stats_retainment_period
// (null: the refresh reads the appliance's value), and drops the three
// credentials: write-only values are never in state, and the passwords held
// the API's hashes anyway. The settings the SDKv2 resource did not have (the
// password policy, sessions, user deactivation, retention periods, sync
// intervals and the rest) are left null, their zero value, and the refresh
// reads them too.
//
// write_only_revision is set null. An upgraded resource has no digests in
// private state, so wodigest.Tracker.ModifyPlan takes its missing-digest path
// and marks write_only_revision unknown: the first plan after the upgrade is
// one "~ write_only_revision" in-place update that re-sends the configured
// credentials and records their digests, and every plan after it is clean.
func upgradeStateV0toV1(
	ctx context.Context,
	req resource.UpgradeStateRequest,
	resp *resource.UpgradeStateResponse,
) {
	var old settingApplianceModelV0

	resp.Diagnostics.Append(req.State.Get(ctx, &old)...)
	if resp.Diagnostics.HasError() {
		return
	}

	upgraded := resourceModel{
		SettingApplianceModel: SettingApplianceModel{
			ApplianceUrl:          old.ApplianceUrl,
			InternalApplianceUrl:  old.InternalApplianceUrl,
			ApiAllowedOrigins:     old.ApiAllowedOrigins,
			RegistrationEnabled:   old.RegistrationEnabled,
			DefaultRoleId:         upgradeInt64("default_role_id", old.DefaultRoleId, &resp.Diagnostics),
			DefaultUserRoleId:     upgradeInt64("default_user_role_id", old.DefaultUserRoleId, &resp.Diagnostics),
			DockerPrivilegedMode:  old.DockerPrivilegedMode,
			SmtpFromAddress:       old.SMTPFromAddress,
			SmtpServer:            old.SMTPServer,
			SmtpPort:              upgradeInt64("smtp_port", old.SMTPPort, &resp.Diagnostics),
			SmtpUseSsl:            old.SMTPUseSSL,
			SmtpUseTls:            old.SMTPUseTLS,
			SmtpUsername:          old.SMTPUsername,
			ProxyHost:             old.ProxyHost,
			ProxyPort:             upgradeInt64("proxy_port", old.ProxyPort, &resp.Diagnostics),
			ProxyUser:             old.ProxyUser,
			ProxyDomain:           old.ProxyDomain,
			ProxyWorkstation:      old.ProxyWorkstation,
			NoProxy:               old.NoProxy,
			CurrencyProvider:      old.CurrencyProvider,
			StatsRetainmentPeriod: old.StatsRetainmentPeriod,
			// Write-only: never in state.
			SmtpPassword:           types.StringNull(),
			ProxyPassword:          types.StringNull(),
			CurrencyProviderApiKey: types.StringNull(),
			// Null so the first plan after the upgrade records the digests; see
			// above.
			WriteOnlyRevision: types.Int64Null(),
		},
		Id: types.StringValue(singletonID),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, upgraded)...)
}

// upgradeInt64 converts a port or role id the SDKv2 resource kept as a
// string. An empty value is unset. One that is not an integer is dropped with
// a warning, as Read does: the refresh that follows the upgrade reads the
// appliance's value again.
func upgradeInt64(name string, v types.String, diags *diag.Diagnostics) types.Int64 {
	if v.IsNull() || v.IsUnknown() || strings.TrimSpace(v.ValueString()) == "" {
		return types.Int64Null()
	}

	n, err := strconv.ParseInt(strings.TrimSpace(v.ValueString()), 10, 64)
	if err != nil {
		diags.AddAttributeWarning(
			path.Root(name),
			"Legacy state value is not an integer",
			fmt.Sprintf("The state written by the previous version of the provider has %s = %q, which is "+
				"not an integer; it is dropped and read from the appliance again.", name, v.ValueString()),
		)

		return types.Int64Null()
	}

	return types.Int64Value(n)
}
