// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instance

import "github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"

// clearableAttributes are the optional attributes a user can remove from the
// configuration to clear on the instance. description is Optional without
// Computed, so the plan carries the null and Update sends it as null; labels
// is Optional+Computed, so ModifyPlan plans the clear (as null — the update
// gate compares the plan with state, and an unknown would be filled from
// state and never fire) and Update sends an empty list. Create and Update
// record which of them the configuration set.
var clearableAttributes = clearable.Attributes{
	{Name: "description", Mode: clearable.Clear},
	{Name: "labels", Mode: clearable.Clear},
}

// notClearable records why each remaining optional attribute is not cleared
// when it is removed from the configuration.
var notClearable = clearable.NotClearable{
	"task_set_id":          "never returned by the API; kept from the configuration",
	"cloud_id":             "the instance's cloud; the current one is kept when it is removed",
	"config":               "free-form configuration; removing the whole map is not yet supported",
	"instance_context":     "the instance's environment; the current one is kept when it is removed",
	"host_name":            "a change replaces the instance; the current value is kept when it is removed",
	"server_uuid":          "a change replaces the instance; the current value is kept when it is removed",
	"server_uuids":         "a change replaces the instance; the current value is kept when it is removed",
	"config_aws":           "nested block; removing a whole block is not yet supported",
	"config_azure":         "nested block; removing a whole block is not yet supported",
	"config_bmaas":         "nested block; removing a whole block is not yet supported",
	"config_hvm":           "nested block; removing a whole block is not yet supported",
	"config_vmware":        "nested block; removing a whole block is not yet supported",
	"evars":                "nested block; removing a whole block is not yet supported",
	"network_domain_id":    "a change or removal replaces the resource",
	"ports":                "nested block; removing a whole block is not yet supported",
	"service_plan_options": "nested block; removing a whole block is not yet supported",
	"tags":                 "nested block; removing a whole block is not yet supported",
	"timeouts":             "nested block; removing a whole block is not yet supported",
	"user_group":           "a change or removal replaces the resource",
	"layout_size":          "has a default; Terraform plans the default when it is removed",
	"volumes":              "nested block; removing a whole block is not yet supported",
	"wait_for_ip_address":  "boolean",
}
