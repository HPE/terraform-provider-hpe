// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkrouterfirewallrulegroup

import "github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"

// clearableAttributes are the optional attributes a user can remove from the
// configuration to clear on the rule group: the description, sent as null.
// It is Optional+Computed, so ModifyPlan plans the clear once an apply has
// recorded it as configured. tenant_ids is Optional-only and is already sent
// as an empty list when null.
var clearableAttributes = clearable.Attributes{
	{Name: "description", Mode: clearable.Clear},
	{Name: "tenant_ids", Mode: clearable.Clear},
}

// notClearable records why each remaining optional attribute is not cleared
// when it is removed from the configuration.
var notClearable = clearable.NotClearable{
	"group_layer": "a change replaces the rule group; the current value is kept when it is removed",
	"priority":    "Morpheus assigns a priority; the current one is kept when it is removed",
	"visibility":  "has a default; Terraform plans the default when it is removed",
}
