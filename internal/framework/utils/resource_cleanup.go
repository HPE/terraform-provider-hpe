package utils

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/cenkalti/backoff/v5"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// DeleteFunc defines a function that deletes a resource by ID.
// It should return the HTTP response and any error encountered.
type DeleteFunc func(ctx context.Context, id int64) (*http.Response, error)

// GetFunc defines a function that retrieves a resource by ID.
// It should return the HTTP response and any error encountered.
type GetFunc func(ctx context.Context, id int64) (resp *http.Response, err error)

// CleanupConfig holds configuration for resource cleanup on error.
type CleanupConfig struct {
	// ResourceType is the human-readable name of the resource (e.g., "instance", "network")
	ResourceType string
	// ResourceID is the ID of the resource to clean up
	ResourceID int64
	// DeleteFunc is the function that performs the deletion
	DeleteFunc DeleteFunc
	// GetFunc is the function that checks if the resource still exists
	GetFunc GetFunc
	// Timeout is how long to wait for cleanup (default: 10 minutes)
	Timeout time.Duration
	// Diagnostics is where to add warnings/errors during cleanup
	Diagnostics *diag.Diagnostics
}

// CleanupResourceOnError attempts to delete a resource that was created but couldn't be fully configured.
// This implements "Option 2" - clean up after ourselves by deleting partially-created resources.
//
// Usage example:
//
//	cleanup := utils.CleanupConfig{
//	    ResourceType: "instance",
//	    ResourceID: instanceId,
//	    DeleteFunc: func(ctx context.Context, id int64) (*http.Response, error) {
//	        _, resp, err := client.InstancesAPI.DeleteInstance(ctx, id).Execute()
//	        return resp, err
//	    },
//	    GetFunc: func(ctx context.Context, id int64) (string, *http.Response, error) {
//	        instance, resp, err := client.InstancesAPI.GetInstance(ctx, id).Execute()
//	        if err != nil {
//	            return "", resp, err
//	        }
//	        return instance.Instance.GetStatus(), resp, nil
//	    },
//	    Diagnostics: &resp.Diagnostics,
//	}
//	utils.CleanupResourceOnError(ctx, cleanup)
func CleanupResourceOnError(ctx context.Context, config CleanupConfig) {
	if config.Timeout == 0 {
		config.Timeout = 10 * time.Minute
	}

	tflog.Warn(ctx, fmt.Sprintf("%s %d: creation failed, attempting cleanup by deleting",
		config.ResourceType, config.ResourceID))

	deleteCtx, deleteCancel := context.WithTimeout(context.Background(), config.Timeout)
	defer deleteCancel()

	// Attempt to delete the resource
	delResp, delErr := config.DeleteFunc(deleteCtx, config.ResourceID)
	if delErr != nil || (delResp != nil && delResp.StatusCode != http.StatusOK) {
		tflog.Error(ctx, fmt.Sprintf("%s %d: failed to delete during error cleanup: %v",
			config.ResourceType, config.ResourceID, delErr))
		config.Diagnostics.AddWarning(
			fmt.Sprintf("failed to clean up %s after error", config.ResourceType),
			fmt.Sprintf("%s %d was created but subsequent operations failed. "+
				"Attempted to delete the %s but deletion also failed. "+
				"You may need to manually delete %s %d. Error: %v",
				capitalize(config.ResourceType), config.ResourceID,
				config.ResourceType,
				config.ResourceType, config.ResourceID, delErr),
		)
		return
	}

	// Wait for deletion to complete
	waitForDelete := func() (string, error) {
		getHttpResp, getErr := config.GetFunc(deleteCtx, config.ResourceID)
		if getHttpResp != nil && getHttpResp.StatusCode == http.StatusNotFound {
			// Resource is gone, success!
			return "deleted", nil
		}
		if getErr != nil {
			// Some other error, keep trying
			return "", nil
		}
		// Resource still exists
		return "exists", nil
	}

	if _, waitErr := backoff.Retry(
		deleteCtx,
		waitForDelete,
		backoff.WithBackOff(backoff.NewConstantBackOff(5*time.Second)),
		backoff.WithMaxElapsedTime(config.Timeout),
	); waitErr != nil {
		tflog.Error(ctx, fmt.Sprintf("%s %d: timed out waiting for deletion during cleanup",
			config.ResourceType, config.ResourceID))
		config.Diagnostics.AddWarning(
			fmt.Sprintf("failed to confirm %s deletion after error", config.ResourceType),
			fmt.Sprintf("%s %d was created but subsequent operations failed. "+
				"Deletion was initiated but we couldn't confirm completion. "+
				"You may need to manually verify/delete %s %d.",
				capitalize(config.ResourceType), config.ResourceID,
				config.ResourceType, config.ResourceID),
		)
	} else {
		tflog.Info(ctx, fmt.Sprintf("%s %d: successfully deleted during error cleanup",
			config.ResourceType, config.ResourceID))
	}
}

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
// This implements "Option 3" - let Terraform handle cleanup by tainting the resource.
// Uses SetAttribute to set only the ID field, which works even when diagnostics has errors.
//
// Call this AFTER adding your error to diagnostics. This function only sets the ID and adds
// a warning about the tainted resource - it does not add the error itself.
//
// Usage example:
//
//	resp.Diagnostics.AddError(
//	    "instance provisioning failed",
//	    fmt.Sprintf("Instance %d failed to reach running status", instanceId),
//	)
//	utils.SetPartialState(ctx, utils.SetPartialStateConfig{
//	    ResourceType: "instance",
//	    ResourceID:   instanceId,
//	    StateWriter:  resp.State,
//	    Diagnostics:  &resp.Diagnostics,
//	})
//	return
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
		fmt.Sprintf("%s partially created", capitalize(config.ResourceType)),
		fmt.Sprintf("%s %d was created but could not be fully configured. "+
			"The resource has been marked as tainted. "+
			"On the next 'terraform apply', Terraform will destroy and recreate this %s. "+
			"If you want to keep this %s, you can import it manually: "+
			"'terraform import <resource_type>.<name> %d'",
			capitalize(config.ResourceType), config.ResourceID,
			config.ResourceType,
			config.ResourceType, config.ResourceID),
	)
}

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return string(s[0]-32) + s[1:]
}
