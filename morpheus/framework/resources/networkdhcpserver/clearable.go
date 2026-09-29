// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkdhcpserver

import "github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"

// clearableAttributes lists the optional attributes a user can remove from
// the configuration to clear on the DHCP server. There are none.
var clearableAttributes = clearable.Attributes{}

// notClearable records why each remaining optional attribute is not cleared
// when it is removed from the configuration.
var notClearable = clearable.NotClearable{
	"config_nsxt": "a nested configuration block, sent as configured and not cleared as a whole",
	"lease_time":  "has a default; Terraform plans the default when it is removed",
	"config":      "free-form configuration; its keys are sent as configured and are not cleared individually",
}
