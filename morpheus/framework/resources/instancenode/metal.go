// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instancenode

import (
	"fmt"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

// The per-node resource pool (selectedResourcePoolId) is only honoured by the
// HPE bare-metal provision provider. Every other provider silently ignores it
// and places the node in the instance's own pool, so setting resource_pool_id
// on a non-metal instance would succeed while doing the wrong thing.
const metalProvisionTypeCode = "hpe-baremetal-plugin.provision"

const notMetalSummary = "resource_pool_id is only valid for bare-metal instances"

func notMetalDetail(instanceID int64, actualCode string) string {
	return fmt.Sprintf(
		"resource_pool_id is set, but instance %d has provision type %q "+
			"(not %q). Setting a per-node resource pool on a non-metal "+
			"instance would succeed but silently ignore the pool, placing "+
			"the node in the instance's own pool. Remove resource_pool_id "+
			"to add a node without pool placement.",
		instanceID, actualCode, metalProvisionTypeCode,
	)
}

// provisionTypeCode extracts the provision type code from an instance's layout.
// Returns ("", false) if the layout or code is absent — uncertainty, not a
// verdict.
func provisionTypeCode(
	inst *sdk.GetInstance200ResponseInstance,
) (string, bool) {
	if inst == nil || inst.Layout == nil || inst.Layout.ProvisionTypeCode == nil {
		return "", false
	}

	return *inst.Layout.ProvisionTypeCode, true
}

// provisionTypeCodeOrUnknown returns the provision type code for error messages,
// falling back to "<unknown>" when it cannot be determined.
func provisionTypeCodeOrUnknown(
	inst *sdk.GetInstance200ResponseInstance,
) string {
	if code, ok := provisionTypeCode(inst); ok {
		return code
	}

	return "<unknown>"
}

// IsMetal returns true if the instance is an HPE bare-metal instance.
func IsMetal(inst *sdk.GetInstance200ResponseInstance) bool {
	code, ok := provisionTypeCode(inst)

	return ok && code == metalProvisionTypeCode
}
