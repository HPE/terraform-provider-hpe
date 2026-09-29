// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkpool

import "github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"

// clearableAttributes are the optional attributes a user can remove from the
// configuration to clear on the pool, all Optional+Computed and cleared once
// an apply has recorded them: the DNS domain, gateway, netmask and subnet
// address, each sent as null.
var clearableAttributes = clearable.Attributes{
	{Name: "dns_domain", Mode: clearable.Clear},
	{Name: "gateway", Mode: clearable.Clear},
	{Name: "netmask", Mode: clearable.Clear},
	{Name: "subnet_address", Mode: clearable.Clear},
}

// notClearable records why each remaining optional attribute is not cleared
// when it is removed from the configuration.
var notClearable = clearable.NotClearable{
	"dhcp_server":  "boolean",
	"pool_enabled": "boolean",
}
