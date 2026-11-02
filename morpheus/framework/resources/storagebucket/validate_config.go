// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package storagebucket

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ resource.ResourceWithValidateConfig = &storageBucketResource{}

// ValidateConfig ties base_path to the local provider type without touching the
// appliance, so `terraform validate` catches it: a local bucket is a directory
// on the appliance and has to say which, and no other type has one (Morpheus
// would silently keep the key in the bucket's config). Unknown values -- a
// provider type or path computed from another resource -- wait for apply.
func (r *storageBucketResource) ValidateConfig(
	ctx context.Context,
	req resource.ValidateConfigRequest,
	resp *resource.ValidateConfigResponse,
) {
	var config StorageBucketModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.ProviderType.IsNull() || config.ProviderType.IsUnknown() || config.BasePath.IsUnknown() {
		return
	}

	switch {
	case isLocal(config.ProviderType) && config.BasePath.IsNull():
		resp.Diagnostics.AddAttributeError(path.Root("base_path"),
			"base_path is required for a local storage bucket",
			"A bucket of provider_type \"local\" stores into a directory on the appliance; set base_path to it "+
				"(for example \"/var/opt/morpheus/morpheus-ui/images\").")
	case !isLocal(config.ProviderType) && !config.BasePath.IsNull():
		resp.Diagnostics.AddAttributeError(path.Root("base_path"),
			"base_path applies to local storage buckets only",
			"base_path is the directory a bucket of provider_type \"local\" stores into; it has no meaning for "+
				"provider_type \""+config.ProviderType.ValueString()+"\". Remove it, or change the provider type.")
	}
}
