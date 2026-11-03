// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package backupinstance

import "github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"

// clearableAttributes lists the optional attributes a user can remove from
// the configuration to clear on the backup. There are none: every optional
// attribute is set when the backup is created.
var clearableAttributes = clearable.Attributes{}

// notClearable records why each remaining optional attribute is not cleared
// when it is removed from the configuration.
var notClearable = clearable.NotClearable{
	"backup_type_code":    "set when the backup is created; a change replaces it, and the current value is kept when it is removed",
	"container_id":        "set when the backup is created; a change replaces it, and the current value is kept when it is removed",
	"enabled":             "has a default; Terraform plans the default when it is removed",
	"storage_provider_id": "set when the backup is created; a change replaces it, and the current value is kept when it is removed",
}
