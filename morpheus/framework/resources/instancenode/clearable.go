// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instancenode

import "github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"

// clearableAttributes lists the optional attributes a user can remove from
// the configuration to clear on the node. There are none: every optional
// attribute is set when the node is created.
var clearableAttributes = clearable.Attributes{}

// notClearable records why each remaining optional attribute is not cleared
// when it is removed from the configuration.
var notClearable = clearable.NotClearable{
	"pre_provisioned":     "boolean",
	"timeouts":            "nested block; removing a whole block is not yet supported",
	"wait_for_ip_address": "boolean",
	"resource_pool_id":    "set when the node is created; a change replaces it, and the current value is kept when it is removed",
	"selected_server_id":  "set when the node is created; a change replaces it, and the current value is kept when it is removed",
	"server_uuid":         "set when the node is created; a change replaces it, and the current value is kept when it is removed",
	"ssh_host":            "set when the node is created; a change replaces it, and the current value is kept when it is removed",
	"ssh_key_pair_id":     "set when the node is created; a change replaces it, and the current value is kept when it is removed",
	"ssh_username":        "set when the node is created; a change replaces it, and the current value is kept when it is removed",
}
