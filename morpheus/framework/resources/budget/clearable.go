// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package budget

import "github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"

// clearableAttributes are the optional attributes a user can remove from the
// configuration to clear on the budget: the description, Optional-only, so
// the plan carries the null and Update sends it.
var clearableAttributes = clearable.Attributes{
	{Name: "description", Mode: clearable.Clear},
}

// notClearable records why each remaining optional attribute is not cleared
// when it is removed from the configuration.
var notClearable = clearable.NotClearable{
	"associated_resource_id": "the resource the budget is scoped to; the current scope is kept when it is removed",
	"enabled":                "has a default; Terraform plans the default when it is removed",
	"end_date":               "defines the budget period; the current period is kept when it is removed",
	"interval":               "has a default; Terraform plans the default when it is removed",
	"scope":                  "has a default; Terraform plans the default when it is removed",
	"start_date":             "defines the budget period; the current period is kept when it is removed",
	"year":                   "defines the budget period; the current period is kept when it is removed",
}
