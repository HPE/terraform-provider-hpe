// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package powerschedule

import "github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"

// clearableAttributes are the optional attributes a user can remove from the
// configuration to clear on the schedule: the description, Optional-only, so
// the plan carries the null and Update sends it.
var clearableAttributes = clearable.Attributes{
	{Name: "description", Mode: clearable.Clear},
}

// notClearable records why each remaining optional attribute is not cleared
// when it is removed from the configuration.
var notClearable = clearable.NotClearable{
	"enabled":            "boolean",
	"friday_off_time":    "has a default; Terraform plans the default when it is removed",
	"friday_on_time":     "has a default; Terraform plans the default when it is removed",
	"monday_off_time":    "has a default; Terraform plans the default when it is removed",
	"monday_on_time":     "has a default; Terraform plans the default when it is removed",
	"saturday_off_time":  "has a default; Terraform plans the default when it is removed",
	"saturday_on_time":   "has a default; Terraform plans the default when it is removed",
	"schedule_timezone":  "has a default; Terraform plans the default when it is removed",
	"schedule_type":      "has a default; Terraform plans the default when it is removed",
	"sunday_off_time":    "has a default; Terraform plans the default when it is removed",
	"sunday_on_time":     "has a default; Terraform plans the default when it is removed",
	"thursday_off_time":  "has a default; Terraform plans the default when it is removed",
	"thursday_on_time":   "has a default; Terraform plans the default when it is removed",
	"tuesday_off_time":   "has a default; Terraform plans the default when it is removed",
	"tuesday_on_time":    "has a default; Terraform plans the default when it is removed",
	"visibility":         "has a default; Terraform plans the default when it is removed",
	"wednesday_off_time": "has a default; Terraform plans the default when it is removed",
	"wednesday_on_time":  "has a default; Terraform plans the default when it is removed",
}
