// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package provisioninglicense

import "github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"

// clearableAttributes are the optional attributes a user can remove from the
// configuration to clear on the license: the description and virtual images
// (Optional-only; Update sends null and an empty list) and the tenants
// (Optional+Computed, cleared once an apply has recorded them; an empty
// list).
var clearableAttributes = clearable.Attributes{
	{Name: "description", Mode: clearable.Clear},
	{Name: "virtual_images", Mode: clearable.Clear},
	{Name: "tenants", Mode: clearable.Clear},
}

// notClearable records why each remaining optional attribute is not cleared
// when it is removed from the configuration.
var notClearable = clearable.NotClearable{
	"license_key_wo":         "write-only",
	"license_key_wo_version": "version companion of a write-only value",
}
