// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package securitygroup

import "github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"

// clearableAttributes are the optional attributes a user can remove from the
// configuration to clear on the security group. Both are Optional without
// Computed: the plan carries the null and Update sends the clear — a null
// description, an empty group list under the permissions object. The
// resource therefore has no ModifyPlan for clears and records nothing in
// private state, deliberately; an Optional+Computed attribute added here
// needs both (TestUnitClearableAttributesNeedNoRecord fails until it has
// them).
var clearableAttributes = clearable.Attributes{
	{Name: "description", Mode: clearable.Clear},
	{Name: "resource_permission_group_ids", Mode: clearable.Clear},
}

// notClearable records why each remaining optional attribute is not cleared
// when it is removed from the configuration.
var notClearable = clearable.NotClearable{
	"cloud_id": "set on create and not updatable; the current cloud is kept when it is removed",
	"tenant_ids": "an empty tenant list makes the API drop the owner's permission row, and with it the group " +
		"permissions; the current tenants are kept",
	"active":                         "boolean",
	"resource_permission_groups_all": "boolean",
	"visibility":                     "has a default; Terraform plans the default when it is removed",
}
