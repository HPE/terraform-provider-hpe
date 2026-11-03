// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package securitygrouprule

import "github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"

// clearableAttributes are the optional attributes a user can remove from the
// configuration to clear on the rule. name is Optional without Computed: the
// plan carries the null and Update sends it as null; with no clearable
// Optional+Computed attribute the rule keeps no record of what was
// configured.
var clearableAttributes = clearable.Attributes{
	{Name: "name", Mode: clearable.Clear},
}

// notClearable records why each remaining optional attribute is not cleared
// when it is removed from the configuration.
var notClearable = clearable.NotClearable{
	"source":           "required by the rule's source type (a CIDR, or a group); removing it is not a valid rule",
	"destination":      "required by the rule's destination type (a CIDR, or a group); removing it is not a valid rule",
	"port_range":       "required by the rule's protocol; removing it is not a valid rule",
	"destination_type": "has a default; Terraform plans the default when it is removed",
	"direction":        "has a default; Terraform plans the default when it is removed",
	"policy":           "has a default; Terraform plans the default when it is removed",
	"source_type":      "has a default; Terraform plans the default when it is removed",
}
