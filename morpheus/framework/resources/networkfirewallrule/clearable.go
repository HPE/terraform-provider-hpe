// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkfirewallrule

import "github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"

// clearableAttributes are the optional attributes a user can remove from the
// configuration to clear on the firewall rule, all Optional+Computed and
// cleared once an apply has recorded them: the description (Update sends
// null) and the application and profile lists (an empty list matches any).
var clearableAttributes = clearable.Attributes{
	{Name: "description", Mode: clearable.Clear},
	{Name: "application", Mode: clearable.Clear},
	{Name: "profile", Mode: clearable.Clear},
}

// notClearable records why each remaining optional attribute is not cleared
// when it is removed from the configuration.
var notClearable = clearable.NotClearable{
	"destinations":  "a nested configuration block, sent as configured and not cleared as a whole",
	"rule_group_id": "a nested configuration block, sent as configured and not cleared as a whole",
	"scopes":        "a nested configuration block, sent as configured and not cleared as a whole",
	"sources":       "a nested configuration block, sent as configured and not cleared as a whole",
	"direction":     "Morpheus requires a direction; the current one is kept when it is removed",
	"enabled":       "boolean",
	"policy":        "Morpheus requires a policy; the current one is kept when it is removed",
	"priority":      "Morpheus assigns a priority when none is given; the current one is kept when it is removed",
}
