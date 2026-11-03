// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package backupjob

import "github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"

// clearableAttributes are the optional attributes a user can remove from the
// configuration to clear on the backup job. Both are Optional without
// Computed: the plan carries the null and Update sends it as null; with no
// clearable Optional+Computed attribute the job keeps no record of what was
// configured.
var clearableAttributes = clearable.Attributes{
	{Name: "retention_count", Mode: clearable.Clear},
	{Name: "schedule_id", Mode: clearable.Clear},
}

// notClearable records why each remaining optional attribute is not cleared
// when it is removed from the configuration.
var notClearable = clearable.NotClearable{
	"code":    "a change or removal replaces the resource",
	"enabled": "boolean",
}
