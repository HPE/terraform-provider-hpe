// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkrouter

import "github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"

// clearableAttributes are the optional attributes a user can remove from the
// configuration to clear on the router. There are none: the router's
// optional attributes are its type, cloud and integration (replace on
// change), booleans, an enumeration, nested blocks and the tenant
// permissions, which the API does not return and the configured set of
// which is kept. The classification is enforced by the unit test so a new
// attribute has to be decided.
var clearableAttributes = clearable.Attributes{}

// notClearable records why each remaining optional attribute is not cleared
// when it is removed from the configuration.
var notClearable = clearable.NotClearable{
	"tenant_ids":                "tenant permissions are not returned by the API; the configured set is kept",
	"visibility":                "an enumeration with no empty value; the current one is kept when it is removed",
	"type_id":                   "a change replaces the router; the current value is kept when it is removed",
	"cloud_id":                  "a change replaces the router; the current value is kept when it is removed",
	"network_integration_id":    "a change replaces the router; the current value is kept when it is removed",
	"config":                    "free-form configuration; removing the whole map is not yet supported",
	"config_nsxt_gateway_tier0": "nested block; removing a whole block is not yet supported",
	"config_nsxt_gateway_tier1": "nested block; removing a whole block is not yet supported",
	"enable_bgp":                "boolean",
	"enabled":                   "boolean",
	"shared_group_access":       "boolean",
}
