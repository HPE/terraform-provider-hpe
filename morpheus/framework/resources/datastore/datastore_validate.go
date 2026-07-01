// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package datastore

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// ValidateConfig enforces the cross-attribute requirements for BMaaS datastores
// that the generated schema cannot express.
func (r *Resource) ValidateConfig(
	ctx context.Context,
	req resource.ValidateConfigRequest,
	resp *resource.ValidateConfigResponse,
) {
	var config DatastoreModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(validateBmaasConfig(config)...)
}

// validateBmaasConfig returns the plan-time diagnostics for a datastore config.
//
// The HPE Alletra MP Bare Metal (BMaaS) datastore type is Cloud-scoped and
// requires both a storage server and a resource pool - the storage plugin
// rejects the create otherwise (validateDatastore). protocol_type lives in the
// config blob, while storage_server and resource_pool are top-level datastore
// fields, so this ties them together at plan time rather than failing the apply
// with a generic API error. It is a pure function (no framework request) so it
// can be unit tested directly.
func validateBmaasConfig(config DatastoreModel) diag.Diagnostics {
	var diags diag.Diagnostics

	// Only BMaaS datastores carry these extra requirements.
	if config.ConfigAlletrampBmaas.IsNull() || config.ConfigAlletrampBmaas.IsUnknown() {
		return diags
	}

	if config.StorageServer.IsNull() || config.StorageServer.IsUnknown() {
		diags.AddAttributeError(
			path.Root("storage_server"),
			"Missing storage_server for BMaaS datastore",
			"storage_server is required when config_alletramp_bmaas is set: HPE Alletra "+
				"MP Bare Metal datastores must be created against a storage server.",
		)
	}

	if config.ResourcePool.IsNull() || config.ResourcePool.IsUnknown() {
		diags.AddAttributeError(
			path.Root("resource_pool"),
			"Missing resource_pool for BMaaS datastore",
			"resource_pool is required when config_alletramp_bmaas is set: HPE Alletra "+
				"MP Bare Metal datastores must be created in a resource pool.",
		)
	}

	if !config.AssociatedResourceType.IsNull() && !config.AssociatedResourceType.IsUnknown() &&
		config.AssociatedResourceType.ValueString() != associatedResourceTypeCloud {
		diags.AddAttributeError(
			path.Root("associated_resource_type"),
			"Invalid associated_resource_type for BMaaS datastore",
			"config_alletramp_bmaas is only supported for Cloud datastores; set "+
				"associated_resource_type = \"Cloud\".",
		)
	}

	return diags
}
