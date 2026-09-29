// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package cloud

import "github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"

// clearableAttributes are the optional attributes a user can remove from the
// configuration to clear on the cloud. code and labels are Optional+Computed,
// so ModifyPlan plans the clear once an apply has recorded them as
// configured; the rest are Optional-only, so the plan already carries the
// null and only Update had to send it: null for code and location, an empty
// list for labels, and an empty string for the configuration-map keys
// (appliance_url, data_center_name, external_id, import_existing_vms,
// keyboard_layout), which Morpheus stores as unset and the read maps back to
// null.
var clearableAttributes = clearable.Attributes{
	{Name: "code", Mode: clearable.Clear},
	{Name: "labels", Mode: clearable.Clear},
	{Name: "location", Mode: clearable.Clear},
	{Name: "appliance_url", Mode: clearable.Clear},
	{Name: "data_center_name", Mode: clearable.Clear},
	{Name: "external_id", Mode: clearable.Clear},
	{Name: "import_existing_vms", Mode: clearable.Clear},
	{Name: "keyboard_layout", Mode: clearable.Clear},
}

// notClearable records why each remaining optional attribute is not cleared
// when it is removed from the configuration.
var notClearable = clearable.NotClearable{
	"cloud_type_code":                    "the cloud's type; the current one is kept when it is removed",
	"group_id":                           "the cloud's group; the current one is kept when it is removed",
	"config":                             "free-form configuration; removing the whole map is not yet supported",
	"config_aws":                         "nested block; removing a whole block is not yet supported",
	"config_azure":                       "nested block; removing a whole block is not yet supported",
	"config_hvm":                         "nested block; removing a whole block is not yet supported",
	"config_vmware":                      "nested block; removing a whole block is not yet supported",
	"agent_install_mode":                 "has a default; Terraform plans the default when it is removed",
	"costing_mode":                       "has a default; Terraform plans the default when it is removed",
	"guidance_mode":                      "has a default; Terraform plans the default when it is removed",
	"security_mode":                      "has a default; Terraform plans the default when it is removed",
	"visibility":                         "has a default; Terraform plans the default when it is removed",
	"auto_recover_power_state":           "boolean",
	"default_datastore_sync_active":      "boolean",
	"default_folder_sync_active":         "boolean",
	"default_network_sync_active":        "boolean",
	"default_plan_sync_active":           "boolean",
	"default_pool_sync_active":           "boolean",
	"default_security_group_sync_active": "boolean",
	"enabled":                            "boolean",
}
