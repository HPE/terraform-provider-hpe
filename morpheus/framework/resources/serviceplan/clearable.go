// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package serviceplan

import "github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"

// clearableAttributes are the optional attributes a user can remove from the
// configuration to clear on the service plan: the description (sent as
// null) and the price sets (sent as an empty list). price_set_ids is
// Optional+Computed, so ModifyPlan plans the clear once an apply has recorded
// it as configured; description is Optional-only and needs no record.
var clearableAttributes = clearable.Attributes{
	{Name: "description", Mode: clearable.Clear},
	{Name: "price_set_ids", Mode: clearable.Clear},
}

// notClearable records why each remaining optional attribute is not cleared
// when it is removed from the configuration.
var notClearable = clearable.NotClearable{
	"add_volumes":        "boolean",
	"custom_cores":       "boolean",
	"custom_cpu":         "boolean",
	"custom_max_memory":  "boolean",
	"custom_max_storage": "boolean",
	"memory_size_type":   "has a default; Terraform plans the default when it is removed",
	"storage_size_type":  "has a default; Terraform plans the default when it is removed",
	"config_ranges":      "nested block; removing a whole block is not yet supported",
	"id":                 "the plan's identifier, configured only to adopt an existing plan",
	"max_cores":          "a plan limit Morpheus keeps a value for; the current value is kept when it is removed",
	"max_cpu":            "a plan limit Morpheus keeps a value for; the current value is kept when it is removed",
	"max_disks":          "a plan limit Morpheus keeps a value for; the current value is kept when it is removed",
	"sort_order":         "Morpheus assigns an order; the current value is kept when it is removed",
}
