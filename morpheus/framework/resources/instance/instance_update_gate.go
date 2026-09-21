// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instance

import (
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// The instance update path makes up to two appliance calls, each gated on its
// own set of attributes:
//
//   - the instance update call carries the attributes in putAttributes;
//   - the resize call carries the attributes in resizeAttributes.
//
// Everything else on the resource is either provider-only (timeouts,
// wait_for_ip_address), applied through a dedicated operation (evars -- see
// applyEvarsUpdate, which Update runs unconditionally and which issues its own
// add/update and delete calls only for variables that differ), or
// create-only and handled by replacement (ports, the config_* blocks and the
// other RequiresReplace attributes). None of those may trigger either call, so
// a change to nothing but `timeouts` must make no appliance call at all. MORPH-17751.
//
// The comparison is done on the raw plan and state values with attributeChanged,
// which fills unknown planned values from prior state before comparing. That is
// what makes the gate correct: an optional-and-computed attribute the
// practitioner did not set is unknown in the plan but known in state, and a
// model-level equality check reads that as a change. Almost every configuration
// has such attributes (a volume's controller_mount_point, for instance), which
// is why the previous model-level gate fired on every update.

// putAttributes are the root attributes the instance update call carries.
// Keep in step with makeUpdateAPIcalls.
var putAttributes = []string{
	"name",
	"description",
	"instance_context",
	"group_id",
	"tags",
	"labels",
	"config_azure", // the only config block applied in place; the rest require replacement
}

// resizeAttributes are the root attributes the resize call carries. Keep in
// step with makeResizeAPICalls and createResizeRequest.
var resizeAttributes = []string{
	"plan_id",
	"volumes",
	"network_interfaces",
	"service_plan_options",
}

// isUpdateCallNeeded reports whether the instance update call must be made:
// true when any attribute it carries was changed by the practitioner.
func isUpdateCallNeeded(plan, state tftypes.Value) bool {
	return anyAttributeChanged(plan, state, putAttributes)
}

// isResizeNeeded reports whether the resize call must be made: true when any
// attribute it carries was changed by the practitioner.
func isResizeNeeded(plan, state tftypes.Value) bool {
	return anyAttributeChanged(plan, state, resizeAttributes)
}

// UpdateCallObserver, when non-nil, is invoked with the name of each appliance
// call the update path issues ("update" or "resize"). It exists so acceptance
// tests -- which run the provider in-process -- can assert which calls an
// update made, which is the whole contract of the gates above. It is nil in
// production and costs a nil check.
var UpdateCallObserver func(call string)

func observeUpdateCall(call string) {
	if UpdateCallObserver != nil {
		UpdateCallObserver(call)
	}
}
