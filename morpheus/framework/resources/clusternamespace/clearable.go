// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package clusternamespace

import "github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"

// clearableAttributes are the optional attributes a user can remove from the
// configuration to clear on the namespace: the description and tenants, both
// Optional-only, so the plan carries the null and Update sends null and an
// empty list.
var clearableAttributes = clearable.Attributes{
	{Name: "description", Mode: clearable.Clear},
	{Name: "tenant_ids", Mode: clearable.Clear},
}

// notClearable records why each remaining optional attribute is not cleared
// when it is removed from the configuration.
var notClearable = clearable.NotClearable{
	"resource_permissions": "a nested configuration block, sent as configured and not cleared as a whole",
	"active":               "has a default; Terraform plans the default when it is removed",
	"visibility":           "has a default; Terraform plans the default when it is removed",
}
