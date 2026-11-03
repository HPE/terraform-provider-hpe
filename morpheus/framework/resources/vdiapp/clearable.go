// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package vdiapp

import "github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"

// clearableAttributes are the optional attributes a user can remove from the
// configuration to clear on the VDI app. description is Optional without
// Computed: the plan carries the null and Update sends it as null; with no
// clearable Optional+Computed attribute the app keeps no record of what was
// configured.
var clearableAttributes = clearable.Attributes{
	{Name: "description", Mode: clearable.Clear},
}

// notClearable records why each remaining optional attribute is not cleared
// when it is removed from the configuration.
var notClearable = clearable.NotClearable{}
