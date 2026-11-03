// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package group

import "github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"

// clearableAttributes are the optional attributes a user can remove from the
// configuration to clear on the group. ModifyPlan plans the clear, Update
// sends it — an empty string for code and location (Morpheus stores it as
// null), an empty list for labels and for the clouds — and Create and Update
// record which of them the configuration set.
var clearableAttributes = clearable.Attributes{
	// Optional, not computed: the plan already carries the null; Update has
	// to send it or the apply fails.
	{Name: "code", Mode: clearable.Clear},
	{Name: "location", Mode: clearable.Clear},
	// Lists: sent empty.
	{Name: "labels", Mode: clearable.Clear},
	{Name: "cloud_ids", Mode: clearable.Clear},
}

// notClearable records why each remaining optional attribute is not cleared
// when it is removed from the configuration. The group has none.
var notClearable = clearable.NotClearable{}
