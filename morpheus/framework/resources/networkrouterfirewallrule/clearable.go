// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkrouterfirewallrule

import "github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"

// clearableAttributes are the optional attributes a user can remove from the
// configuration to clear on the firewall rule. description is Optional-only,
// so the plan carries the null; it is never returned by the API and is
// cleared with an empty string, as before. application and port_range are
// Optional+Computed — an empty one matches any application or port — so
// ModifyPlan plans the clear once an apply has recorded them as configured,
// and Update sends null.
var clearableAttributes = clearable.Attributes{
	{Name: "description", Mode: clearable.Clear},
	{Name: "application", Mode: clearable.Clear},
	{Name: "port_range", Mode: clearable.Clear},
}

// notClearable records why each remaining optional attribute is not cleared
// when it is removed from the configuration.
var notClearable = clearable.NotClearable{
	"direction":        "an enumeration with no empty value; the current one is kept when it is removed",
	"protocol":         "an enumeration with no empty value; the current one is kept when it is removed",
	"source_type":      "an enumeration with no empty value; the current one is kept when it is removed",
	"destination_type": "an enumeration with no empty value; the current one is kept when it is removed",
	"priority":         "Morpheus assigns a priority; the current one is kept when it is removed",
	"policy":           "has a default; Terraform plans the default when it is removed",
	"enabled":          "boolean",
}
