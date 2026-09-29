// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package containerscript

import "github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"

// clearableAttributes are the optional attributes a user can remove from the
// configuration to clear on the script: the category and run-as user
// (Optional-only; Update sends null) and the labels (Optional+Computed,
// cleared once an apply has recorded them; Update sends an empty list).
var clearableAttributes = clearable.Attributes{
	{Name: "category", Mode: clearable.Clear},
	{Name: "run_as_user", Mode: clearable.Clear},
	{Name: "labels", Mode: clearable.Clear},
}

// notClearable records why each remaining optional attribute is not cleared
// when it is removed from the configuration.
var notClearable = clearable.NotClearable{
	"fail_on_error":  "boolean",
	"script_phase":   "has a default; Terraform plans the default when it is removed",
	"script_type":    "has a default; Terraform plans the default when it is removed",
	"script_version": "has a default; Terraform plans the default when it is removed",
	"sudo_user":      "boolean",
	"script":         "the script body; Morpheus keeps the current one when it is removed",
}
