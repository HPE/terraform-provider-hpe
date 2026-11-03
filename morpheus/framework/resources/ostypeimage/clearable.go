// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package ostypeimage

import "github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"

// clearableAttributes lists the optional attributes a user can remove from
// the configuration to clear. There are none: both are set when the
// association is created.
var clearableAttributes = clearable.Attributes{}

// notClearable records why each remaining optional attribute is not cleared
// when it is removed from the configuration.
var notClearable = clearable.NotClearable{
	"cloud_id":          "set when the association is created; a change replaces it, and the current value is kept when it is removed",
	"provision_type_id": "set when the association is created; a change replaces it, and the current value is kept when it is removed",
}
