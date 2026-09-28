// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package backuphost

import "github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"

// clearableAttributes are the optional attributes a user can remove from the
// configuration to clear on the appliance. ModifyPlan plans the clear, Update
// sends it in the shape each field takes on the wire, and Create and Update
// record which of them the configuration set.
var clearableAttributes = clearable.Attributes{
	{Name: "path", Mode: clearable.Clear},
	{Name: "ssh_username", Mode: clearable.Clear},
	{Name: "storage_provider_id", Mode: clearable.Clear},
}

// notClearable records why each remaining optional attribute is not cleared
// when it is removed from the configuration.
var notClearable = clearable.NotClearable{
	"ssh_password_wo_version": "version companion of a write-only value",
	"enabled":                 "boolean",
}
