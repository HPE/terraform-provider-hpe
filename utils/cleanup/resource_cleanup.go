package cleanup

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// TaintResourceStateConfig holds configuration for tainting resource state on error.
type TaintResourceStateConfig struct {
	// ResourceType is the human-readable name of the resource (e.g. "instance", "cloud")
	ResourceType string
	// ResourceID is the ID of the resource that was created
	ResourceID int64
	// StateWriter is the response state to write to (resource.CreateResponse.State)
	StateWriter State
	// Diagnostics is where to add warnings (errors should already be added before calling this)
	Diagnostics *diag.Diagnostics
}

// State is an interface for state operations (implemented by resource.CreateResponse.State)
type State interface {
	Set(ctx context.Context, val any) diag.Diagnostics
	SetAttribute(ctx context.Context, path path.Path, val any) diag.Diagnostics
}

// TaintResourceState sets just the ID in state to mark the resource as tainted.
//
// A Create that has already written the plan to state (to record the id before
// a long poll, say) leaves every computed attribute unknown. Terraform rejects
// unknown values in post-apply state even when the provider returns an error,
// and would add a "Provider returned invalid result object after apply …
// attribute X is unknown" diagnostic per attribute on top of the real failure.
// So, when the writer is the framework's state, every unknown is replaced with
// null first (MORPH-18935).
func TaintResourceState(ctx context.Context, config TaintResourceStateConfig) {
	tflog.Warn(ctx, fmt.Sprintf("%s %d created but encountered error - setting partial state with ID only",
		config.ResourceType, config.ResourceID))

	if st, ok := config.StateWriter.(*tfsdk.State); ok {
		nullUnknowns(st)
	}

	// Set ONLY the ID using SetAttribute - this works even when returning errors
	setDiags := config.StateWriter.SetAttribute(ctx, path.Root("id"), types.Int64Value(config.ResourceID))
	config.Diagnostics.Append(setDiags...)
	if setDiags.HasError() {
		tflog.Error(ctx, fmt.Sprintf("Failed to set ID attribute: %v", setDiags))
	} else {
		tflog.Info(ctx, fmt.Sprintf("Successfully set ID attribute to %d", config.ResourceID))
	}

	// Add helpful guidance to the user
	config.Diagnostics.AddWarning(
		fmt.Sprintf("%s partially created", config.ResourceType),
		fmt.Sprintf("%s %d was created but could not be fully configured. "+
			"The resource has been marked as tainted. "+
			"On the next 'terraform apply', Terraform will destroy and recreate this %s. "+
			"If you want to keep this %s, you can import it manually: "+
			"'terraform import <resource_type>.<name> %d'",
			config.ResourceType, config.ResourceID,
			config.ResourceType,
			config.ResourceType, config.ResourceID),
	)
}

// nullUnknowns replaces every unknown value in the state with null, in place.
// A null or already-known state is left untouched.
func nullUnknowns(st *tfsdk.State) {
	if st == nil || st.Raw.IsNull() || st.Raw.IsFullyKnown() {
		return
	}

	transformed, err := tftypes.Transform(st.Raw, func(_ *tftypes.AttributePath, v tftypes.Value) (tftypes.Value, error) {
		if !v.IsKnown() {
			return tftypes.NewValue(v.Type(), nil), nil
		}

		return v, nil
	})
	if err != nil {
		// Transform only fails on an internal inconsistency; leave the state as
		// it was rather than lose the id write that follows.
		return
	}

	st.Raw = transformed
}
