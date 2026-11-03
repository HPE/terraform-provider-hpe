// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package loadbalancermonitor

import "github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"

// clearableAttributes are the optional attributes a user can remove from the
// configuration to clear on the monitor: the description and the free-text
// probe settings, all sent as null. All are Optional+Computed, so ModifyPlan
// plans the clear once an apply has recorded them as configured; the
// intervals and counts Morpheus keeps a value for are not clearable.
var clearableAttributes = clearable.Attributes{
	{Name: "description", Mode: clearable.Clear},
	{Name: "send_data", Mode: clearable.Clear},
	{Name: "send_type", Mode: clearable.Clear},
	{Name: "receive_data", Mode: clearable.Clear},
	{Name: "receive_code", Mode: clearable.Clear},
	{Name: "monitor_destination", Mode: clearable.Clear},
	{Name: "extra_config", Mode: clearable.Clear},
}

// notClearable records why each remaining optional attribute is not cleared
// when it is removed from the configuration.
var notClearable = clearable.NotClearable{
	"monitor_username":            "set only when the monitor is created; Morpheus ignores it on update, and a change replaces the monitor",
	"monitor_type":                "the monitor's type; the current one is kept when it is removed",
	"monitor_interval":            "a monitor setting Morpheus keeps a value for; the current value is kept when it is removed",
	"monitor_timeout":             "a monitor setting Morpheus keeps a value for; the current value is kept when it is removed",
	"fall_count":                  "a monitor setting Morpheus keeps a value for; the current value is kept when it is removed",
	"rise_count":                  "a monitor setting Morpheus keeps a value for; the current value is kept when it is removed",
	"alias_port":                  "a monitor setting Morpheus keeps a value for; the current value is kept when it is removed",
	"data_length":                 "a monitor setting Morpheus keeps a value for; the current value is kept when it is removed",
	"max_retry":                   "a monitor setting Morpheus keeps a value for; the current value is kept when it is removed",
	"send_version":                "version companion of a write-only value",
	"monitor_password_wo_version": "version companion of a write-only value",
	"config":                      "nested block; removing a whole block is not yet supported",
}
