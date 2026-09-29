// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package task

import "github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"

// clearableAttributes are the optional attributes a user can remove from the
// configuration to clear on the task, all Optional+Computed and cleared once
// an apply has recorded them: the code, result type, retry count and delay
// (Update sends null) and the labels (an empty list).
var clearableAttributes = clearable.Attributes{
	{Name: "code", Mode: clearable.Clear},
	{Name: "result_type", Mode: clearable.Clear},
	{Name: "retry_count", Mode: clearable.Clear},
	{Name: "retry_delay_seconds", Mode: clearable.Clear},
	{Name: "labels", Mode: clearable.Clear},
}

// notClearable records why each remaining optional attribute is not cleared
// when it is removed from the configuration.
var notClearable = clearable.NotClearable{
	"config_conditional_workflow": "a nested configuration block, sent as configured and not cleared as a whole",
	"allow_custom_config":         "has a default; Terraform plans the default when it is removed",
	"config":                      "free-form configuration; its keys are sent as configured and are not cleared individually",
	"execute_target":              "has a default; Terraform plans the default when it is removed",
	"retryable":                   "has a default; Terraform plans the default when it is removed",
	"visibility":                  "has a default; Terraform plans the default when it is removed",
}
