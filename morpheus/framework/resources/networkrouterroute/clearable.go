// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkrouterroute

import "github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"

// clearableAttributes are the optional attributes a user can remove from the
// configuration on the route. The route has no update: every change replaces
// it, and so does a removal — ModifyPlan plans the recorded attribute as
// unknown and requires replacement, and the route is recreated without it,
// taking whatever value Morpheus assigns. network_mtu is Optional-only, so
// its own RequiresReplace modifier already does this.
var clearableAttributes = clearable.Attributes{
	{Name: "description", Mode: clearable.Release},
	{Name: "name", Mode: clearable.Release},
	{Name: "priority", Mode: clearable.Release},
}

// notClearable records why each remaining optional attribute is not cleared
// when it is removed from the configuration.
var notClearable = clearable.NotClearable{
	"network_mtu":   "a change or removal replaces the route through its own plan modifier",
	"default_route": "boolean",
	"enabled":       "boolean",
}
