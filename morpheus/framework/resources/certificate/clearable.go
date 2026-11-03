// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package certificate

import "github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"

// clearableAttributes are the optional attributes a user can remove from the
// configuration to clear on the certificate: the description, Optional-only,
// so the plan carries the null and Update sends it.
var clearableAttributes = clearable.Attributes{
	{Name: "description", Mode: clearable.Clear},
}

// notClearable records why each remaining optional attribute is not cleared
// when it is removed from the configuration.
var notClearable = clearable.NotClearable{
	"domain_name": "derived from the certificate when it is not set; the current value is kept when it is removed",
}
