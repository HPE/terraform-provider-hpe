// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkrouterbgpneighbor

import "github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"

// clearableAttributes are the optional attributes a user can remove from the
// configuration to clear on the BGP neighbour: the description, the NSX-V
// forwarding and protocol addresses and the route filter names, all sent as
// null. All are Optional+Computed, so ModifyPlan plans the clear once an apply
// has recorded them as configured.
var clearableAttributes = clearable.Attributes{
	{Name: "description", Mode: clearable.Clear},
	{Name: "forwarding_address", Mode: clearable.Clear},
	{Name: "protocol_address", Mode: clearable.Clear},
	{Name: "route_filtering_in", Mode: clearable.Clear},
	{Name: "route_filtering_out", Mode: clearable.Clear},
}

// notClearable records why each remaining optional attribute is not cleared
// when it is removed from the configuration.
var notClearable = clearable.NotClearable{
	"remote_as":            "required by the neighbour; the current value is kept when it is removed",
	"bfd_interval":         "a BFD timer Morpheus keeps a value for; the current value is kept when it is removed",
	"bfd_multiple":         "a BFD timer Morpheus keeps a value for; the current value is kept when it is removed",
	"config":               "free-form configuration; removing the whole map is not yet supported",
	"config_nsxt":          "nested block; removing a whole block is not yet supported",
	"config_nsxv":          "nested block; removing a whole block is not yet supported",
	"hold_down":            "has a default; Terraform plans the default when it is removed",
	"hop_limit":            "has a default; Terraform plans the default when it is removed",
	"keep_alive":           "has a default; Terraform plans the default when it is removed",
	"restart_mode":         "has a default; Terraform plans the default when it is removed",
	"route_filtering_type": "has a default; Terraform plans the default when it is removed",
	"weight":               "has a default; Terraform plans the default when it is removed",
	"allow_as_in":          "boolean",
	"bfd_enabled":          "boolean",
	"password_wo_version":  "version companion of a write-only value",
}
