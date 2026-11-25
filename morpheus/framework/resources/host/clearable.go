// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package host

import "github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"

// clearableAttributes are the optional attributes a user can remove from the
// configuration to clear on the host in place. Update sends the cleared value.
var clearableAttributes = clearable.Attributes{
	{Name: "description", Mode: clearable.Clear},
	{Name: "labels", Mode: clearable.Clear},
	{Name: "ssh_username", Mode: clearable.Clear},
	{Name: "user_data", Mode: clearable.Clear},
	{Name: "ssh_key_pair_id", Mode: clearable.Clear},
}

// notClearable records why each remaining optional attribute is not cleared
// when it is removed from the configuration.
var notClearable = clearable.NotClearable{
	"cloud_id":            "create-only; changing or removing replaces the host",
	"group_id":            "create-only; changing or removing replaces the host",
	"compute_server_type": "create-only; changing or removing replaces the host",
	"plan_id":             "create-only; changing or removing replaces the host",
	"config":              "create-only; nested/dynamic, removing replaces the host",
	"config_bmaas":        "create-only; nested block, removing replaces the host",
	"config_omega_bmaas":  "create-only; nested block, removing replaces the host",
	"network_interfaces":  "create-only; nested list, removing replaces the host",
	"volumes":             "create-only; nested list, removing replaces the host",
	"delete_options":      "nested block that only affects destroy; not updatable",
	"enabled":             "boolean with a computed value; Terraform keeps the current value",
	"visibility":          "has a default; Terraform plans the default when it is removed",
	"ssh_password_wo":     "write-only value; not stored in state",
	"write_only_revision": "computed: the time a write-only value was last sent",
	"timeouts":            "nested block; removing a whole block is not yet supported",
}
