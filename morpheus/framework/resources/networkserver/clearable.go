// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkserver

import "github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"

// clearableAttributes are the optional attributes a user can remove from the
// configuration to clear on the network server. All are Optional-only (not
// Computed), so the plan carries the null and Update sends the clear; no
// private-state record is needed.
var clearableAttributes = clearable.Attributes{
	{Name: "service_url", Mode: clearable.Clear},
	{Name: "service_username", Mode: clearable.Clear},
	{Name: "network_filter", Mode: clearable.Clear},
	{Name: "tenant_match", Mode: clearable.Clear},
	{Name: "cloud_id", Mode: clearable.Clear},
	{Name: "tenant_ids", Mode: clearable.Clear},
}

// notClearable records why each remaining optional attribute is not cleared
// when it is removed from the configuration.
var notClearable = clearable.NotClearable{ //nolint:gosec // attribute names and reasons, not credentials
	"config":                      "free-form type-specific option map; the stored values are kept when it is removed",
	"config_juniper_sdn":          "a nested type block; its controller address is RequiresReplace, so removing it is not a clear",
	"credential_id":               "a credential reference; switch back to local credentials by setting service_username and service_password_wo",
	"service_host":                "has a server-computed default; the current host is kept when it is removed",
	"service_mode":                "has a server-computed default; the current mode is kept when it is removed",
	"service_port":                "has a server-computed default; the current port is kept when it is removed",
	"service_password_wo":         "a write-only value; it is sent only when set",
	"service_password_wo_version": "version companion of a write-only value",
	"service_token_wo":            "a write-only value; it is sent only when set",
	"visibility":                  "has a default; Terraform plans the default when it is removed",
}
