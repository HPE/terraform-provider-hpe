// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package loadbalancerprofile

import "github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"

// clearableAttributes are the optional attributes a user can remove from the
// configuration to clear on the profile: the description, Optional without
// Computed, so the plan carries the null and Update sends it as null; with no
// clearable Optional+Computed attribute the profile keeps no record of what
// was configured.
var clearableAttributes = clearable.Attributes{
	{Name: "description", Mode: clearable.Clear},
}

// notClearable records why each remaining optional attribute is not cleared
// when it is removed from the configuration.
var notClearable = clearable.NotClearable{
	"config_client_ssl":            "nested block; removing a whole block is not yet supported",
	"config_cookie_persistence":    "nested block; removing a whole block is not yet supported",
	"config_fast_tcp":              "nested block; removing a whole block is not yet supported",
	"config_fast_udp":              "nested block; removing a whole block is not yet supported",
	"config_generic_persistence":   "nested block; removing a whole block is not yet supported",
	"config_http":                  "nested block; removing a whole block is not yet supported",
	"config_server_ssl":            "nested block; removing a whole block is not yet supported",
	"config_source_ip_persistence": "nested block; removing a whole block is not yet supported",
	"tags":                         "nested block; removing a whole block is not yet supported",
}
