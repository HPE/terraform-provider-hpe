// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkpoolip

import "github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"

// clearableAttributes are the optional attributes a user can remove from the
// configuration to clear on the reservation: ref_type and ref_id, the
// association with a server or other resource. Both are Optional (not
// Computed), so the framework plans the null and Update sends it: ref_id as an
// explicit null (NullableInt64), ref_type as refType:null in the request body
// (see update.go for why the plain *string pointer cannot express that).
var clearableAttributes = clearable.Attributes{
	{Name: "ref_type", Mode: clearable.Clear},
	{Name: "ref_id", Mode: clearable.Clear},
}

// notClearable records why each remaining optional attribute is not cleared
// when it is removed from the configuration.
var notClearable = clearable.NotClearable{
	"description": "create-only; changing or removing it replaces the reservation (RequiresReplace)",
	"domain_id":   "create-only; changing or removing it replaces the reservation (RequiresReplace)",
	"create_dns":  "create-only; changing or removing it replaces the reservation (RequiresReplace)",
	"ip_address":  "create-only; changing or removing it replaces the reservation (RequiresReplace)",
}
