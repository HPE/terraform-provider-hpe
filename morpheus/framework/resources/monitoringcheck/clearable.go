// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package monitoringcheck

import "github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"

// clearableAttributes are the optional attributes a user can remove from the
// configuration to clear on the check. description is Optional without
// Computed: the plan carries the null and Update sends it as null; with no
// clearable Optional+Computed attribute the check keeps no record of what
// was configured.
var clearableAttributes = clearable.Attributes{
	{Name: "description", Mode: clearable.Clear},
}

// notClearable records why each remaining optional attribute is not cleared
// when it is removed from the configuration.
var notClearable = clearable.NotClearable{
	"active":         "boolean",
	"check_interval": "has a default; Terraform plans the default when it is removed",
	"in_uptime":      "boolean",
	"severity":       "has a default; Terraform plans the default when it is removed",
}
