// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package cluster

import "github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"

// clearableAttributes are the optional attributes a user can remove from the
// configuration to clear on the cluster. Both are Optional without Computed,
// so the plan carries the null and Update sends the clear: an empty string
// for the description (Morpheus stores it as null) and, for the labels, an
// empty labelArray — the cluster API rejects an empty labels list as an
// unrecognised field, because its handler only renames a non-empty one to
// the key it accepts. With no clearable Optional+Computed attribute the
// cluster keeps no record of what was configured.
var clearableAttributes = clearable.Attributes{
	{Name: "description", Mode: clearable.Clear},
	{Name: "labels", Mode: clearable.Clear},
}

// notClearable records why each remaining optional attribute is not cleared
// when it is removed from the configuration.
var notClearable = clearable.NotClearable{
	"cluster_type_code": "the cluster's type; the current one is kept when it is removed",
	"config":            "free-form configuration; removing the whole map is not yet supported",
	"config_hvm":        "nested block; removing a whole block is not yet supported",
	"timeouts":          "nested block; removing a whole block is not yet supported",
	"workflow_id":       "a change or removal replaces the cluster",
}
