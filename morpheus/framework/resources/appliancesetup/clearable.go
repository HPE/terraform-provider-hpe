// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package appliancesetup

import "github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"

// clearableAttributes lists the optional attributes a user can remove from
// the configuration to clear on the appliance. There are none: setup is a
// one-time operation and nothing is sent after it.
var clearableAttributes = clearable.Attributes{}

// notClearable records why each optional attribute is not cleared when it is
// removed from the configuration.
var notClearable = clearable.NotClearable{
	"appliance_url":  "recorded URL changes are refused; omission retains the computed URL",
	"appliance_name": "changing or removing a recorded name is refused; setup cannot rename the appliance",
	"first_name":     "changing or removing a recorded value is refused; manage the user through user resources",
	"last_name":      "changing or removing a recorded value is refused; manage the user through user resources",
	"hub_mode":       "omission restores the default; a resulting change to recorded Hub mode is refused",
	"timeouts":       "nested block; removing a whole block is not yet supported",
}
