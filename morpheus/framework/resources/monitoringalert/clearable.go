// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package monitoringalert

import "github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"

// clearableAttributes lists the optional attributes a user can remove from
// the configuration to clear on the alert. There are none: every optional
// attribute is a boolean or has a default.
var clearableAttributes = clearable.Attributes{}

// notClearable records why each remaining optional attribute is not cleared
// when it is removed from the configuration.
var notClearable = clearable.NotClearable{
	"active":       "boolean",
	"all_checks":   "boolean",
	"all_groups":   "boolean",
	"min_duration": "has a default; Terraform plans the default when it is removed",
	"min_severity": "has a default; Terraform plans the default when it is removed",
}
