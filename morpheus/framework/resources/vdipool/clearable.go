// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package vdipool

import "github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"

// clearableAttributes are the optional attributes a user can remove from the
// configuration to clear on the VDI pool. description is Optional without
// Computed: the plan carries the null and Update sends it as null; with no
// clearable Optional+Computed attribute the pool keeps no record of what was
// configured.
var clearableAttributes = clearable.Attributes{
	{Name: "description", Mode: clearable.Clear},
}

// notClearable records why each remaining optional attribute is not cleared
// when it is removed from the configuration.
var notClearable = clearable.NotClearable{
	"max_session_timeout":                   "never returned by the API; kept from the configuration",
	"idle_timeout":                          "pool sizing Morpheus defaults; the current value is kept when it is removed",
	"initial_pool_size":                     "pool sizing Morpheus defaults; the current value is kept when it is removed",
	"max_idle":                              "pool sizing Morpheus defaults; the current value is kept when it is removed",
	"min_idle":                              "pool sizing Morpheus defaults; the current value is kept when it is removed",
	"allow_copy":                            "boolean",
	"allow_fileshare":                       "boolean",
	"allow_hypervisor_console":              "boolean",
	"allow_printer":                         "boolean",
	"auto_create_local_user_on_reservation": "boolean",
	"enabled":                               "boolean",
	"persistent_user":                       "boolean",
	"recyclable":                            "boolean",
}
