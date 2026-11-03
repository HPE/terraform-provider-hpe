// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package loadbalancerpool

import "github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"

// clearableAttributes are the optional attributes a user can remove from the
// configuration to clear on the pool: the description, sent as null. It is
// Optional+Computed, so ModifyPlan plans the clear once an apply has recorded
// it as configured.
var clearableAttributes = clearable.Attributes{
	{Name: "description", Mode: clearable.Clear},
}

// notClearable records why each remaining optional attribute is not cleared
// when it is removed from the configuration.
var notClearable = clearable.NotClearable{
	"config":             "free-form configuration; removing the whole map is not yet supported",
	"config_nsxt":        "nested block; removing a whole block is not yet supported",
	"min_active":         "a pool setting Morpheus keeps a value for; the current value is kept when it is removed",
	"port":               "a pool setting Morpheus keeps a value for; the current value is kept when it is removed",
	"vip_balance":        "an enumeration with no empty value; the current one is kept when it is removed",
	"vip_client_ip_mode": "an enumeration with no empty value; the current one is kept when it is removed",
	"vip_sticky":         "an enumeration with no empty value; the current one is kept when it is removed",
}
