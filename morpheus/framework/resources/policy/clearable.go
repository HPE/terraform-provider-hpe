// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package policy

import "github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"

// clearableAttributes are the optional attributes a user can remove from the
// configuration to clear on the policy, both Optional+Computed and cleared
// once an apply has recorded them: the description (Update sends null) and
// the tenants (an empty list).
var clearableAttributes = clearable.Attributes{
	{Name: "description", Mode: clearable.Clear},
	{Name: "tenants", Mode: clearable.Clear},
}

// notClearable records why each remaining optional attribute is not cleared
// when it is removed from the configuration.
var notClearable = clearable.NotClearable{
	"config_approval":             "a nested configuration block, sent as configured and not cleared as a whole",
	"config_backup_storage":       "a nested configuration block, sent as configured and not cleared as a whole",
	"config_create_backup":        "a nested configuration block, sent as configured and not cleared as a whole",
	"config_create_user":          "a nested configuration block, sent as configured and not cleared as a whole",
	"config_create_user_group":    "a nested configuration block, sent as configured and not cleared as a whole",
	"config_cypher":               "a nested configuration block, sent as configured and not cleared as a whole",
	"config_delayed_removal":      "a nested configuration block, sent as configured and not cleared as a whole",
	"config_host_naming":          "a nested configuration block, sent as configured and not cleared as a whole",
	"config_lifecycle":            "a nested configuration block, sent as configured and not cleared as a whole",
	"config_max_containers":       "a nested configuration block, sent as configured and not cleared as a whole",
	"config_max_cores":            "a nested configuration block, sent as configured and not cleared as a whole",
	"config_max_hosts":            "a nested configuration block, sent as configured and not cleared as a whole",
	"config_max_memory":           "a nested configuration block, sent as configured and not cleared as a whole",
	"config_max_networks":         "a nested configuration block, sent as configured and not cleared as a whole",
	"config_max_pool_members":     "a nested configuration block, sent as configured and not cleared as a whole",
	"config_max_pools":            "a nested configuration block, sent as configured and not cleared as a whole",
	"config_max_price":            "a nested configuration block, sent as configured and not cleared as a whole",
	"config_max_routers":          "a nested configuration block, sent as configured and not cleared as a whole",
	"config_max_snapshots":        "a nested configuration block, sent as configured and not cleared as a whole",
	"config_max_storage":          "a nested configuration block, sent as configured and not cleared as a whole",
	"config_max_virtual_servers":  "a nested configuration block, sent as configured and not cleared as a whole",
	"config_max_vms":              "a nested configuration block, sent as configured and not cleared as a whole",
	"config_motd":                 "a nested configuration block, sent as configured and not cleared as a whole",
	"config_naming":               "a nested configuration block, sent as configured and not cleared as a whole",
	"config_power_schedule":       "a nested configuration block, sent as configured and not cleared as a whole",
	"config_required_network":     "a nested configuration block, sent as configured and not cleared as a whole",
	"config_server_naming":        "a nested configuration block, sent as configured and not cleared as a whole",
	"config_shutdown":             "a nested configuration block, sent as configured and not cleared as a whole",
	"config_storage_server_quota": "a nested configuration block, sent as configured and not cleared as a whole",
	"config_tags":                 "a nested configuration block, sent as configured and not cleared as a whole",
	"config_workflow":             "a nested configuration block, sent as configured and not cleared as a whole",
	"associated_resource_id":      "set when the policy is created; a change replaces it, and the current value is kept when it is removed",
	"config":                      "free-form configuration; its keys are sent as configured and are not cleared individually",
	"each_user":                   "boolean",
	"enabled":                     "boolean",
}
