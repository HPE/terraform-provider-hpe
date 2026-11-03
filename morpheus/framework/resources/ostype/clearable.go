// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package ostype

import "github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"

// clearableAttributes are the optional attributes a user can remove from the
// configuration to clear on the OS type: the description, family, version,
// cloud-init version (Optional+Computed, cleared once an apply has recorded
// them) and the category, codename, name and vendor (Optional-only). Update
// sends null for each.
var clearableAttributes = clearable.Attributes{
	{Name: "description", Mode: clearable.Clear},
	{Name: "os_family", Mode: clearable.Clear},
	{Name: "os_version", Mode: clearable.Clear},
	{Name: "cloud_init_version", Mode: clearable.Clear},
	{Name: "category", Mode: clearable.Clear},
	{Name: "os_codename", Mode: clearable.Clear},
	{Name: "os_name", Mode: clearable.Clear},
	{Name: "vendor", Mode: clearable.Clear},
}

// notClearable records why each remaining optional attribute is not cleared
// when it is removed from the configuration.
var notClearable = clearable.NotClearable{
	"install_agent": "boolean",
}
