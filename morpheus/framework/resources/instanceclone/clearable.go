// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instanceclone

import "github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"

// clearableAttributes lists the optional attributes a user can remove from
// the configuration to clear on the cloned instance. There are none: the
// clone's inputs are consumed when it is created.
var clearableAttributes = clearable.Attributes{}

// notClearable records why each remaining optional attribute is not cleared
// when it is removed from the configuration.
var notClearable = clearable.NotClearable{
	"config_aws":                 "has a default; Terraform plans the default when it is removed",
	"config_azure":               "has a default; Terraform plans the default when it is removed",
	"config_hvm":                 "has a default; Terraform plans the default when it is removed",
	"config_vmware":              "has a default; Terraform plans the default when it is removed",
	"timeouts":                   "nested block; removing a whole block is not yet supported",
	"cloud_id":                   "set when the clone is created; a change replaces it, and the current value is kept when it is removed",
	"config":                     "free-form configuration; its keys are sent as configured and are not cleared individually",
	"group_id":                   "set when the clone is created; a change replaces it, and the current value is kept when it is removed",
	"plan_id":                    "the clone's plan; the current plan is kept when it is removed",
	"source_instance_id_version": "version companion of the source instance id",
	"user_group":                 "set when the clone is created; a change replaces it, and the current value is kept when it is removed",
}
