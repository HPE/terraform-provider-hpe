// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package optionlist

import "github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"

// clearableAttributes are the optional attributes a user can remove from the
// configuration to clear on the option list: the description, source URL
// and API type, all Optional-only, so the plan carries the null and Update
// sends it.
var clearableAttributes = clearable.Attributes{
	{Name: "description", Mode: clearable.Clear},
	{Name: "source_url", Mode: clearable.Clear},
	{Name: "api_type", Mode: clearable.Clear},
}

// notClearable records why each remaining optional attribute is not cleared
// when it is removed from the configuration.
var notClearable = clearable.NotClearable{
	"real_time":  "has a default; Terraform plans the default when it is removed",
	"type":       "has a default; Terraform plans the default when it is removed",
	"visibility": "has a default; Terraform plans the default when it is removed",
}
