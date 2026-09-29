// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package loadbalancervirtualserver

import "github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"

// clearableAttributes are the optional attributes a user can remove from the
// configuration to clear on the virtual server. All are Optional without
// Computed, so the plan carries the null and Update sends it: null for the
// description, the host name, the certificates and the pool reference; with
// no clearable Optional+Computed attribute the virtual server keeps no record
// of what was configured.
var clearableAttributes = clearable.Attributes{
	{Name: "description", Mode: clearable.Clear},
	{Name: "vip_hostname", Mode: clearable.Clear},
	{Name: "ssl_cert", Mode: clearable.Clear},
	{Name: "ssl_server_cert", Mode: clearable.Clear},
	{Name: "pool_id", Mode: clearable.Clear},
}

// notClearable records why each remaining optional attribute is not cleared
// when it is removed from the configuration.
var notClearable = clearable.NotClearable{
	"vip_pool":     "the network pool the VIP address was leased from when the virtual server was created; Morpheus ignores it on update, and the leased address stays when it is removed",
	"vip_address":  "required by the virtual server; the current value is kept when it is removed",
	"vip_name":     "required by the virtual server; the current value is kept when it is removed",
	"vip_port":     "required by the virtual server; the current value is kept when it is removed",
	"vip_protocol": "required by the virtual server; the current value is kept when it is removed",
	"vip_type":     "a change replaces the virtual server; the current value is kept when it is removed",
	"config":       "a change replaces the virtual server; removing the whole map is not yet supported",
	"config_nsxt":  "nested block; removing a whole block is not yet supported",
}
