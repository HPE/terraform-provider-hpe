// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package datastore

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
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

	// storage_server and resource_pool are Optional+Computed with an
	// Optional+Computed nested id, so an empty block (e.g. storage_server = {})
	// satisfies !IsNull() while leaving a null id that the API would reject.
	//
	// Unknown is different from null: an unknown block or id is a reference
	// Terraform has not resolved yet (a data source during validate, say, or
	// a plan that depends on another change). It cannot be judged here --
	// the check runs again once the value is known, and the API has the final
	// say on apply -- so it is let through rather than reported as missing,
	// which would fail `terraform validate` for every configuration that
	// resolves its storage server or pool by name.
	if nestedIDMissing(config.StorageServer.IsNull(), config.StorageServer.IsUnknown(), config.StorageServer.Id) {
		diags.AddAttributeError(
			path.Root("storage_server"),
			"Missing storage_server for BMaaS datastore",
			"storage_server (with a populated id) is required when config_alletramp_bmaas "+
				"is set: HPE Alletra MP Bare Metal datastores must be created against a "+
				"storage server.",
		)
	}

	if nestedIDMissing(config.ResourcePool.IsNull(), config.ResourcePool.IsUnknown(), config.ResourcePool.Id) {
		diags.AddAttributeError(
			path.Root("resource_pool"),
			"Missing resource_pool for BMaaS datastore",
			"resource_pool (with a populated id) is required when config_alletramp_bmaas "+
				"is set: HPE Alletra MP Bare Metal datastores must be created in a "+
				"resource pool.",
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

// nestedIDMissing reports whether a single-nested block with an id is
// decidably missing: the block is null, or it is known and its id is null. An
// unknown block or id is not missing -- it is not yet known -- and reports
// false. The block's own id is only consulted when the block is known and
// non-null, since an unknown or null block carries a zero-value id.
func nestedIDMissing(blockNull, blockUnknown bool, id types.Int64) bool {
	switch {
	case blockUnknown:
		return false
	case blockNull:
		return true
	default:
		return id.IsNull()
	}
}
