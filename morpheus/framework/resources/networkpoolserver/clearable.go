// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkpoolserver

import "github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"

// clearableAttributes are the optional attributes a user can remove from the
// configuration to clear on the pool server: the service URL and username
// and the network, zone and tenant filters, all Optional-only, so the plan
// carries the null and Update sends it.
var clearableAttributes = clearable.Attributes{
	{Name: "service_url", Mode: clearable.Clear},
	{Name: "service_username", Mode: clearable.Clear},
	{Name: "network_filter", Mode: clearable.Clear},
	{Name: "zone_filter", Mode: clearable.Clear},
	{Name: "tenant_match", Mode: clearable.Clear},
}

// notClearable records why each remaining optional attribute is not cleared
// when it is removed from the configuration.
var notClearable = clearable.NotClearable{ //nolint:gosec // attribute names and reasons, not credentials
	"service_mode":                "has a default on the API side; the current mode is kept when it is removed",
	"credential_id":               "a credential reference; the current credential is kept when it is removed",
	"app_id":                      "required by phpIPAM and unused by other types; the current value is kept when it is removed",
	"enabled":                     "has a default; Terraform plans the default when it is removed",
	"ignore_ssl":                  "has a default; Terraform plans the default when it is removed",
	"inventory_existing":          "has a default; Terraform plans the default when it is removed",
	"service_password_wo":         "a write-only value; it is sent only when set",
	"service_password_wo_version": "version companion of a write-only value",
	"service_throttle_rate":       "has a default; Terraform plans the default when it is removed",
	"type_code":                   "set when the pool server is created; a change replaces it, and the current value is kept when it is removed",
	"type_id":                     "set when the pool server is created; a change replaces it, and the current value is kept when it is removed",
}
