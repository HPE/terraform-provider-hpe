// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package image

import "github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"

// clearableAttributes are the optional attributes a user can remove from the
// configuration to clear on the image. All are Optional+Computed, so ModifyPlan
// plans the clear once an apply has recorded them as configured, and Update
// sends null (the minimum disk, the OS type, the SSH username, the user data)
// or an empty list (labels, tenants).
var clearableAttributes = clearable.Attributes{
	{Name: "labels", Mode: clearable.Clear},
	{Name: "min_disk", Mode: clearable.Clear},
	{Name: "os_type_id", Mode: clearable.Clear},
	{Name: "ssh_username", Mode: clearable.Clear},
	{Name: "tenant_ids", Mode: clearable.Clear},
	{Name: "user_data", Mode: clearable.Clear},
}

// notClearable records why each remaining optional attribute is not cleared
// when it is removed from the configuration.
var notClearable = clearable.NotClearable{
	"min_ram":                 "Morpheus rejects a null minimum memory with a server error and ignores zero; the current value is kept",
	"description":             "not updatable through the API; a change replaces the image, and the current value is kept when it is removed",
	"storage_provider_id":     "where the image was uploaded; Morpheus ignores it on update, a change replaces the image, and the current value is kept when it is removed",
	"url":                     "a change replaces the image; the current value is kept when it is removed",
	"visibility":              "has a default; Terraform plans the default when it is removed",
	"ssh_password_wo_version": "version companion of a write-only value",
	"config_azure":            "nested block; removing a whole block is not yet supported",
	"tags":                    "nested block; removing a whole block is not yet supported",
	"auto_join_domain":        "boolean",
	"cloud_init":              "boolean",
	"fips_enabled":            "boolean",
	"force_customization":     "boolean",
	"install_agent":           "boolean",
	"sysprep":                 "boolean",
	"trial_version":           "boolean",
	"uefi":                    "boolean",
	"virtio_supported":        "boolean",
	"vm_tools_installed":      "boolean",
}
