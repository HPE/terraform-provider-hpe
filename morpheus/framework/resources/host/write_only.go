// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package host

import (
	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/wodigest"
)

// Change detection for the write-only attributes: see wodigest. A changed
// ssh_password_wo marks write_only_revision unknown, which plans an in-place
// update in which Update re-sends the configured value; APIHash additionally
// detects a change made outside Terraform through the appliance's hash of the
// stored password. The two typed-variant passwords are sent only when the
// host is added (their blocks are create-only), so a change to either plans a
// replacement.
var wo = wodigest.Tracker{ //nolint:gosec // attribute names, not credentials
	Resource: "host",
	Attributes: map[string]string{
		"ssh_password_wo":                        "",
		"config_bmaas.ilo_password_wo":           "",
		"config_omega_bmaas.console_password_wo": "",
	},
	Immutable: []string{"config_bmaas.ilo_password_wo", "config_omega_bmaas.console_password_wo"},
	APIHash: map[string]wodigest.HashExtractor{
		"ssh_password_wo": sshPasswordHash,
	},
	OnDrift: wodigest.ResendPolicy,
}

// sshPasswordHash pulls the appliance's hash of the stored SSH password from
// the read-back server; see wodigest.HashExtractor for the omitted/null rule.
func sshPasswordHash(resp any) (string, bool) {
	srv, ok := resp.(*sdk.GetHost200ResponseServer)
	if !ok || srv == nil {
		return "", false
	}

	return wodigest.NullableStringHash(srv.SshPasswordHash)
}

// writeOnlyFromConfig extracts the write-only password values from the config
// model (read from req.Config). Write-only values are null in plan and state, so
// they must be read from configuration. Each returns nil when unset.
func writeOnlyFromConfig(config HostModel) (sshPwd, bmaasPwd, omegaPwd *string) {
	if !config.SshPasswordWo.IsNull() && !config.SshPasswordWo.IsUnknown() {
		sshPwd = config.SshPasswordWo.ValueStringPointer()
	}

	if !config.ConfigBmaas.IsNull() && !config.ConfigBmaas.IsUnknown() {
		p := config.ConfigBmaas.IloPasswordWo
		if !p.IsNull() && !p.IsUnknown() {
			bmaasPwd = p.ValueStringPointer()
		}
	}

	if !config.ConfigOmegaBmaas.IsNull() && !config.ConfigOmegaBmaas.IsUnknown() {
		p := config.ConfigOmegaBmaas.ConsolePasswordWo
		if !p.IsNull() && !p.IsUnknown() {
			omegaPwd = p.ValueStringPointer()
		}
	}

	return sshPwd, bmaasPwd, omegaPwd
}
