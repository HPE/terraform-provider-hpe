// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package subnet

import "github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"

// clearableAttributes are the optional attributes a user can remove from the
// configuration to clear on the subnet. All of them are Optional without
// Computed, so the plan carries the null and Update sends the clear — an
// empty string for the description (Morpheus stores it as null), a null pool
// id, an empty list for the labels and for the group permissions; with no
// clearable Optional+Computed attribute the subnet keeps no record of what
// was configured.
var clearableAttributes = clearable.Attributes{
	{Name: "description", Mode: clearable.Clear},
	{Name: "pool_id", Mode: clearable.Clear},
	{Name: "labels", Mode: clearable.Clear},
	{Name: "resource_permission_group_ids", Mode: clearable.Clear},
}

// notClearable records why each remaining optional attribute is not cleared
// when it is removed from the configuration.
var notClearable = clearable.NotClearable{
	"active":                         "boolean",
	"allow_static_override":          "boolean",
	"dhcp_server":                    "boolean",
	"resource_permission_groups_all": "boolean",
	"visibility":                     "has a default; Terraform plans the default when it is removed",
	"config_version":                 "version companion of a write-only value",
	"tenant_ids":                     "a change or removal replaces the subnet",
	"name":                           "Morpheus requires a name; the current one is kept when it is removed",
	"cidr": "Morpheus derives it from the subnet's configuration; the current value is kept when it is " +
		"removed",
}
