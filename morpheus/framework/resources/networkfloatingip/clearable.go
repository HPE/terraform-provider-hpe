// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkfloatingip

import "github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"

// clearableAttributes is empty: this resource has no optional attributes to
// clear. The two user inputs are Required and RequiresReplace, and every other
// attribute is computed.
var clearableAttributes = clearable.Attributes{}

// notClearable is empty for the same reason. Both are kept, with the
// completeness test, so a future optional attribute must be classified.
var notClearable = clearable.NotClearable{}
