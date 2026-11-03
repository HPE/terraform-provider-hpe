// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package clusteraffinitygroupmember

import "github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"

// clearableAttributes lists the optional attributes a user can remove from
// the configuration to clear on the membership. There are none.
var clearableAttributes = clearable.Attributes{}

// notClearable records why each remaining optional attribute is not cleared
// when it is removed from the configuration.
var notClearable = clearable.NotClearable{}
