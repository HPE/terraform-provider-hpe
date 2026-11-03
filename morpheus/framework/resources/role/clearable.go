// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package role

import "github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"

// clearableAttributes are the optional attributes a user can remove from the
// configuration to clear on the role. All three are Optional without Computed,
// so the plan carries the null and Update sends it as null; with no clearable
// Optional+Computed attribute the role keeps no record of what was
// configured.
var clearableAttributes = clearable.Attributes{
	{Name: "default_persona_code", Mode: clearable.Clear},
	{Name: "description", Mode: clearable.Clear},
	{Name: "landing_url", Mode: clearable.Clear},
}

// notClearable records why each remaining optional attribute is not cleared
// when it is removed from the configuration.
var notClearable = clearable.NotClearable{
	"multitenant":        "boolean",
	"multitenant_locked": "boolean",
	"role_type":          "has a default; Terraform plans the default when it is removed",
	"permissions":        "nested block; removing a whole block is not yet supported",
}
