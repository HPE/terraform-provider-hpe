// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package tenant

import "github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"

// clearableAttributes are the optional attributes a user can remove from the
// configuration to clear on the tenant. description is Optional without
// Computed and has always been sent as an explicit null; the billing fields
// and the subdomain are Optional+Computed, so ModifyPlan plans the clear once
// an apply has recorded them as configured, and Update sends null.
var clearableAttributes = clearable.Attributes{
	{Name: "description", Mode: clearable.Clear},
	{Name: "subdomain", Mode: clearable.Clear},
	{Name: "account_number", Mode: clearable.Clear},
	{Name: "account_name", Mode: clearable.Clear},
	{Name: "customer_number", Mode: clearable.Clear},
}

// notClearable records why each remaining optional attribute is not cleared
// when it is removed from the configuration.
var notClearable = clearable.NotClearable{
	"enabled":          "boolean",
	"remove_resources": "boolean",
	"currency":         "has a default; Terraform plans the default when it is removed",
	"base_role_id":     "the tenant's role; Morpheus requires one, so the current role is kept when it is removed",
	"parent_id":        "a change replaces the tenant; the current parent is kept when it is removed",
}
