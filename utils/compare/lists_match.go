// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package compare

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// ListsMatch is a generic that compares lists from plan and state to see if they are the same.
// Returns true if they are, false otherwise.
func ListsMatch[S attr.Value](
	ctx context.Context,
	planList, stateList basetypes.ListValue,
) (bool, diag.Diagnostics) {
	var planVals, stateVals []S

	diags := planList.ElementsAs(ctx, &planVals, false)
	if diags.HasError() {
		tflog.Error(ctx, fmt.Sprintf("cannot convert plan list values to type %T", planVals))

		return false, diags
	}

	diags = stateList.ElementsAs(ctx, &stateVals, false)
	if diags.HasError() {
		tflog.Error(ctx, fmt.Sprintf("cannot convert state list values to type %T", stateVals))

		return false, diags
	}

	// Check length of lists first
	if len(planVals) != len(stateVals) {
		return false, nil
	}

	// Compare each element in the lists to see if they are the same
	for i, planVal := range planVals {
		stateVal := stateVals[i]

		if !planVal.Equal(stateVal) {
			return false, nil
		}
	}

	return true, nil
}