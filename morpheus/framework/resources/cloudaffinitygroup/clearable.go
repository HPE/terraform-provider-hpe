// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package cloudaffinitygroup

import "github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"

// clearableAttributes lists the optional attributes a user can remove from
// the configuration to clear on the affinity group. There are none: the
// type and pool are set when the group is created, and the API manages its
// tenants (see create.go).
var clearableAttributes = clearable.Attributes{}

// notClearable records why each remaining optional attribute is not cleared
// when it is removed from the configuration.
var notClearable = clearable.NotClearable{
	"active":               "boolean",
	"resource_permissions": "a nested configuration block, sent as configured and not cleared as a whole",
	"visibility":           "has a default; Terraform plans the default when it is removed",
	"affinity_type":        "set when the affinity group is created; a change replaces it, and the current value is kept when it is removed",
	"pool_id":              "set when the affinity group is created; a change replaces it, and the current value is kept when it is removed",
	"tenant_ids":           "not sent to the API, which manages the tenants (MORPH-15806)",
}
