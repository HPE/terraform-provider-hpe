// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkfirewallrulegroup

import "github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"

// clearableAttributes are the optional attributes a user can remove from the
// configuration to clear on the rule group, both Optional+Computed and
// cleared once an apply has recorded them: the description (Update sends
// null) and the tenants (an empty list).
var clearableAttributes = clearable.Attributes{
	{Name: "description", Mode: clearable.Clear},
	{Name: "tenant_ids", Mode: clearable.Clear},
}

// notClearable records why each remaining optional attribute is not cleared
// when it is removed from the configuration.
var notClearable = clearable.NotClearable{
	"visibility":  "has a default; Terraform plans the default when it is removed",
	"group_layer": "set when the rule group is created; a change replaces it, and the current value is kept when it is removed",
	"priority":    "Morpheus orders the groups by priority; the current one is kept when it is removed",
}
