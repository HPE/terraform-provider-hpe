// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instance

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// newControllerID is the id component of a controller_mount_point that asks
// Morpheus to create the controller rather than use an existing one. It is what
// hpe_morpheus_instance_storage_controller composes, and what the documented
// "id:busNumber:typeId:unitNumber" format calls out for new storage controllers.
const newControllerID = "-1"

// reconcileControllerMountPoint decides what state records for a volume's
// controller_mount_point given what was configured and what Morpheus reports
// (MORPH-18424).
//
// Morpheus resolves a mount point on bus number and controller type only, never
// on the id (AbstractBoxProvisionService.assignStorageVolumeController), so a
// configured "-1:1:6:0" and a reported "1559905:1:6:0" are the same request,
// satisfied. Terraform, however, requires the applied value to equal the
// planned one, so the configured form is kept whenever the bus:type:unit tail
// agrees. Any other case — nothing configured, a configured real controller id,
// or a tail that differs — records the API value.
func reconcileControllerMountPoint(configured types.String, api *string) *string {
	if configured.IsNull() || configured.IsUnknown() || api == nil {
		return api
	}
	want := configured.ValueString()
	id, wantTail, ok := strings.Cut(want, ":")
	if !ok || id != newControllerID {
		return api
	}
	if mountPointTail(*api) != wantTail {
		return api
	}

	return &want
}

// mountPointTail returns the "busNumber:typeId:unitNumber" part of a
// controller_mount_point, i.e. everything after the controller id.
func mountPointTail(mountPoint string) string {
	_, tail, ok := strings.Cut(mountPoint, ":")
	if !ok {
		return ""
	}

	return tail
}
