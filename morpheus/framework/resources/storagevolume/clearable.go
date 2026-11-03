// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package storagevolume

import "github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"

// clearableAttributes lists the optional attributes a user can remove from
// the configuration to clear on the volume. There are none: every optional
// attribute is set when the volume is created, or is its size, which
// Morpheus keeps.
var clearableAttributes = clearable.Attributes{}

// notClearable records why each remaining optional attribute is not cleared
// when it is removed from the configuration.
var notClearable = clearable.NotClearable{
	"config_alletramp_bmaas_wo_version": "version companion of a write-only value",
	"config_wo_version":                 "version companion of a write-only value",
	"max_storage":                       "the volume's size; Morpheus keeps the current size when it is removed",
	"provision_type":                    "set when the volume is created; a change replaces it, and the current value is kept when it is removed",
	"storage_group_id":                  "set when the volume is created; a change replaces it, and the current value is kept when it is removed",
	"storage_server_id":                 "set when the volume is created; a change replaces it, and the current value is kept when it is removed",
	"type_code":                         "set when the volume is created; a change replaces it, and the current value is kept when it is removed",
	"type_id":                           "set when the volume is created; a change replaces it, and the current value is kept when it is removed",
}
