// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instancesnapshot

import "github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"

// clearableAttributes lists the optional attributes a user can remove from
// the configuration to clear on the snapshot. There are none: a snapshot is
// not edited after it is taken.
var clearableAttributes = clearable.Attributes{}

// notClearable records why each remaining optional attribute is not cleared
// when it is removed from the configuration.
var notClearable = clearable.NotClearable{
	"for_export":       "boolean",
	"memory_snapshot":  "boolean",
	"retain_on_delete": "boolean",
	"timeouts":         "nested block; removing a whole block is not yet supported",
	"description":      "set when the snapshot is created; a change replaces it, and the current value is kept when it is removed",
}
