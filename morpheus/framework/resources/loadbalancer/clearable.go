// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package loadbalancer

import "github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"

// clearableAttributes are the optional attributes a user can remove from the
// configuration to clear on the load balancer: the description, sent as null.
// It is Optional+Computed, so ModifyPlan plans the clear once an apply has
// recorded it as configured.
var clearableAttributes = clearable.Attributes{
	{Name: "description", Mode: clearable.Clear},
}

// notClearable records why each remaining optional attribute is not cleared
// when it is removed from the configuration.
var notClearable = clearable.NotClearable{
	"cloud_id":                  "a change replaces the load balancer; the current value is kept when it is removed",
	"type_code":                 "a change replaces the load balancer; the current value is kept when it is removed",
	"config":                    "free-form configuration; removing the whole map is not yet supported",
	"config_haproxy":            "nested block; removing a whole block is not yet supported",
	"config_nsxt":               "nested block; removing a whole block is not yet supported",
	"permissions":               "nested block; removing a whole block is not yet supported",
	"tenants":                   "nested block; removing a whole block is not yet supported",
	"group_id_version":          "version companion of a write-only value",
	"network_server_id_version": "version companion of a write-only value",
	"visibility":                "has a default; Terraform plans the default when it is removed",
	"enabled":                   "boolean",
}
