// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkrouternat

import "github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"

// clearableAttributes are the optional attributes a user can remove from the
// configuration to clear on the NAT rule: the description and the optional
// match and translation fields, all sent as null (an empty source or
// destination matches any address). All are Optional+Computed, so ModifyPlan
// plans the clear once an apply has recorded them as configured.
var clearableAttributes = clearable.Attributes{
	{Name: "description", Mode: clearable.Clear},
	{Name: "source_network", Mode: clearable.Clear},
	{Name: "destination_network", Mode: clearable.Clear},
	{Name: "translated_network", Mode: clearable.Clear},
	{Name: "translated_ports", Mode: clearable.Clear},
}

// notClearable records why each remaining optional attribute is not cleared
// when it is removed from the configuration.
var notClearable = clearable.NotClearable{
	"protocol":         "deprecated, superseded by the service block; the current value is kept when it is removed",
	"priority":         "has a default; Terraform plans the default when it is removed",
	"enabled":          "boolean",
	"action_version":   "version companion of a write-only value",
	"firewall_version": "version companion of a write-only value",
	"service_version":  "version companion of a write-only value",
}
