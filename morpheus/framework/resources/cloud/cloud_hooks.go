// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package cloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// cloudPostRead handles API quirks that can't be expressed declaratively.
func cloudPostRead(_ context.Context, _ map[string]any, state any, plan any) error {
	cloudState, ok := state.(*CloudModel)
	if !ok {
		return nil
	}

	// Quirk: the API returns code="standard" as a default when no code was set.
	// This causes drift because the plan has code=null. Treat "standard" as null
	// unless the plan explicitly set it.
	if cloudState.Code.ValueString() == "standard" {
		shouldNull := true

		if plan != nil {
			if planModel, ok := plan.(*CloudModel); ok {
				if !planModel.Code.IsNull() && !planModel.Code.IsUnknown() {
					shouldNull = false
				}
			}
		}

		if shouldNull {
			cloudState.Code = types.StringNull()
		}
	}

	// On import (plan has null name), inject known defaults for write-only config fields
	// that the API doesn't return but that were set during creation.
	isImport := plan == nil
	if plan != nil {
		if planModel, ok := plan.(*CloudModel); ok {
			isImport = planModel.Name.IsNull()
		}
	}

	if isImport && !cloudState.ConfigHvm.IsNull() && !cloudState.ConfigHvm.IsUnknown() {
		if cloudState.ConfigHvm.CertificateProvider.IsNull() {
			cloudState.ConfigHvm.CertificateProvider = types.StringValue("internal")
		}
	}

	if isImport && !cloudState.ConfigVmware.IsNull() && !cloudState.ConfigVmware.IsUnknown() {
		if cloudState.ConfigVmware.CertificateProvider.IsNull() {
			cloudState.ConfigVmware.CertificateProvider = types.StringValue("internal")
		}
	}

	return nil
}
