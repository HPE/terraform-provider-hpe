package utils

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// SetPartialStateConfig holds configuration for setting partial state on error.
type SetPartialStateConfig struct {
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

// SetPartialState sets just the ID in state to mark the resource as tainted.
func SetPartialState(ctx context.Context, config SetPartialStateConfig) {
	tflog.Warn(ctx, fmt.Sprintf("%s %d created but encountered error - setting partial state with ID only",
		config.ResourceType, config.ResourceID))

	// Set ONLY the ID using SetAttribute - this works even when returning errors
	// Unlike State.Set() which is ignored on error, SetAttribute can set individual fields
	// IMPORTANT: Must convert int64 to types.Int64 - raw primitives don't work!
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
