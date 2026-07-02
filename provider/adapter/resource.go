// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package adapter

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// ResourceAdapter wraps a Terraform Plugin Framework resource and adapts it
// to work as a child resource within a parent provider architecture. It allows
// child provider resources to be properly namespaced and integrated into a
// parent provider.
//
// The adapter implements all optional resource interfaces
// (ResourceWithConfigure, ResourceWithConfigValidators,
// ResourceWithImportState, ResourceWithModifyPlan, ResourceWithMoveState,
// ResourceWithUpgradeState, and ResourceWithValidateConfig) by delegating to
// the wrapped resource if it implements them. For optional interfaces not
// implemented by the wrapped resource, the adapter returns nil or no-op
// responses.
//
// Key transformations performed by the adapter:
//
// 1. Metadata: Transforms the resource TypeName by prepending the child
// provider's TypeName. For example, a "network" resource from a "morpheus"
// provider becomes "morpheus_network", and when used in a parent "hpe"
// provider becomes "hpe_morpheus_network".
//
// 2. Configure: Extracts the child provider's configuration data from the
// parent provider's ConfigureRequest.ProviderData map, ensuring the resource
// receives only its own provider's data.
//
// The adapter does not support ResourceWithIdentity and
// ResourceWithUpgradeIdentity interfaces because Terraform Plugin Framework
// checks for Resource Identity at gRPC server creation time, unlike other
// optional interfaces which are handled during RPC request processing.
type ResourceAdapter struct {
	in       resource.Resource
	provider provider.Provider // we need the provider so we can access its name from metadata

	withConfigure        resource.ResourceWithConfigure
	withConfigValidators resource.ResourceWithConfigValidators
	withImportState      resource.ResourceWithImportState
	withModifyPlan       resource.ResourceWithModifyPlan
	withMoveState        resource.ResourceWithMoveState
	withUpgradeState     resource.ResourceWithUpgradeState
	withValidateConfig   resource.ResourceWithValidateConfig

	// We don't support ResourceWithIdentity and ResourceWithUpgradeIdentity in the Resource Adapter.
	// This is because Terraform Plugin Framework checks for Resource Identity at gRPC server create,
	// unlike the other interfaces which are handled at RPC request time.
}

var (
	_ resource.Resource                     = &ResourceAdapter{}
	_ resource.ResourceWithConfigure        = &ResourceAdapter{}
	_ resource.ResourceWithConfigValidators = &ResourceAdapter{}
	_ resource.ResourceWithImportState      = &ResourceAdapter{}
	_ resource.ResourceWithModifyPlan       = &ResourceAdapter{}
	_ resource.ResourceWithMoveState        = &ResourceAdapter{}
	_ resource.ResourceWithUpgradeState     = &ResourceAdapter{}
	_ resource.ResourceWithValidateConfig   = &ResourceAdapter{}
)

func NewResourceAdapter(in resource.Resource, p provider.Provider) *ResourceAdapter {
	r := &ResourceAdapter{in: in, provider: p}

	r.withConfigure, _ = in.(resource.ResourceWithConfigure)
	r.withConfigValidators, _ = in.(resource.ResourceWithConfigValidators)
	r.withImportState, _ = in.(resource.ResourceWithImportState)
	r.withModifyPlan, _ = in.(resource.ResourceWithModifyPlan)
	r.withMoveState, _ = in.(resource.ResourceWithMoveState)
	r.withUpgradeState, _ = in.(resource.ResourceWithUpgradeState)
	r.withValidateConfig, _ = in.(resource.ResourceWithValidateConfig)

	return r
}

func NewAdaptedResource(in resource.Resource, p provider.Provider) resource.Resource {
	return NewResourceAdapter(in, p)
}

// We use the Provider Adapter's name to the Metadata request.
// This will transform the resource name from e.g.:
// resource -> {child_provider}_resource
// When a parent provider is introduced, the resource name will
// then be registered as e.g.:
// {parent_provider}_{child_provider}_resource
func (r *ResourceAdapter) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	// Get the provider's name from its metadata
	providerMetaResp := &provider.MetadataResponse{}
	r.provider.Metadata(ctx, provider.MetadataRequest{}, providerMetaResp)

	req.ProviderTypeName = req.ProviderTypeName + "_" + providerMetaResp.TypeName
	r.in.Metadata(ctx, req, resp)
}

func (r *ResourceAdapter) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	r.in.Schema(ctx, req, resp)
}

func (r *ResourceAdapter) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	r.in.Create(ctx, req, resp)
}

func (r *ResourceAdapter) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	r.in.Read(ctx, req, resp)
}

func (r *ResourceAdapter) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	r.in.Update(ctx, req, resp)
}

func (r *ResourceAdapter) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	r.in.Delete(ctx, req, resp)
}

func (r *ResourceAdapter) Configure(
	ctx context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	if r.withConfigure == nil {
		return
	}

	// Extract child provider configure data for ConfigureRequest.ProviderData
	if providerData, ok := req.ProviderData.(map[string]any); ok {
		metaResp := &provider.MetadataResponse{}
		r.provider.Metadata(ctx, provider.MetadataRequest{}, metaResp)

		childData, exists := providerData[metaResp.TypeName]
		if !exists {
			resp.Diagnostics.AddError(
				"Missing provider configuration",
				fmt.Sprintf(
					"The %q provider block is required but was not found in the provider configuration.",
					metaResp.TypeName,
				),
			)

			return
		}

		req.ProviderData = childData
	}

	r.withConfigure.Configure(ctx, req, resp)
}

func (r *ResourceAdapter) ConfigValidators(
	ctx context.Context,
) []resource.ConfigValidator {
	if r.withConfigValidators == nil {
		return nil
	}

	return r.withConfigValidators.ConfigValidators(ctx)
}

func (r *ResourceAdapter) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	if r.withImportState == nil {
		resp.Diagnostics.AddError(
			"Import Not Supported",
			"This resource does not support import.",
		)

		return
	}

	r.withImportState.ImportState(ctx, req, resp)
}

func (r *ResourceAdapter) ModifyPlan(
	ctx context.Context,
	req resource.ModifyPlanRequest,
	resp *resource.ModifyPlanResponse,
) {
	if r.withModifyPlan == nil {
		return
	}

	r.withModifyPlan.ModifyPlan(ctx, req, resp)
}

func (r *ResourceAdapter) MoveState(
	ctx context.Context,
) []resource.StateMover {
	if r.withMoveState == nil {
		return nil
	}

	return r.withMoveState.MoveState(ctx)
}

func (r *ResourceAdapter) UpgradeState(
	ctx context.Context,
) map[int64]resource.StateUpgrader {
	if r.withUpgradeState == nil {
		return nil
	}

	return r.withUpgradeState.UpgradeState(ctx)
}

func (r *ResourceAdapter) ValidateConfig(
	ctx context.Context,
	req resource.ValidateConfigRequest,
	resp *resource.ValidateConfigResponse,
) {
	if r.withValidateConfig == nil {
		return
	}

	r.withValidateConfig.ValidateConfig(ctx, req, resp)
}
