package adapter

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

type AdapterResource struct {
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

type AdapterResourceWithConfigure struct {
	in       resource.ResourceWithConfigure
	provider provider.Provider
	resource.Resource
}

func (r *AdapterResourceWithConfigure) Configure(
	ctx context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	r.in.Configure(ctx, req, resp)
}

var _ resource.Resource = &AdapterResource{}
var _ resource.ResourceWithConfigure = &AdapterResource{}
var _ resource.ResourceWithConfigValidators = &AdapterResource{}
var _ resource.ResourceWithImportState = &AdapterResource{}
var _ resource.ResourceWithModifyPlan = &AdapterResource{}
var _ resource.ResourceWithMoveState = &AdapterResource{}
var _ resource.ResourceWithUpgradeState = &AdapterResource{}
var _ resource.ResourceWithValidateConfig = &AdapterResource{}

func NewAdapterResource(in resource.Resource, p provider.Provider) *AdapterResource {
	r := &AdapterResource{in: in, provider: p}

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
	return NewAdapterResource(in, p)
}

// Metadata is the only method implementation that varies from `in`
// We use the Provider Adapter's name to the Metadata request.
// This will transform the resource name from e.g.:
// resource -> {child_provider}_resource
// When a parent provier is introduced, the resource name will
// then be registered as e.g.:
// {parent_provider}_{child_provider}_resource
func (r *AdapterResource) Metadata(
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

func (r *AdapterResource) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	r.in.Schema(ctx, req, resp)
}

func (r *AdapterResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	r.in.Create(ctx, req, resp)
}

func (r *AdapterResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	r.in.Read(ctx, req, resp)
}

func (r *AdapterResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	r.in.Update(ctx, req, resp)
}

func (r *AdapterResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	r.in.Delete(ctx, req, resp)
}

func (r *AdapterResource) Configure(
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

		req.ProviderData = providerData[metaResp.TypeName]
	}
	r.withConfigure.Configure(ctx, req, resp)
}

func (r *AdapterResource) ConfigValidators(
	ctx context.Context,
) []resource.ConfigValidator {
	if r.withConfigValidators == nil {
		return nil
	}

	return r.withConfigValidators.ConfigValidators(ctx)
}

func (r *AdapterResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	if r.withImportState == nil {
		return
	}

	r.withImportState.ImportState(ctx, req, resp)
}

func (r *AdapterResource) ModifyPlan(
	ctx context.Context,
	req resource.ModifyPlanRequest,
	resp *resource.ModifyPlanResponse,
) {
	if r.withModifyPlan == nil {
		return
	}

	r.withModifyPlan.ModifyPlan(ctx, req, resp)
}

func (r *AdapterResource) MoveState(
	ctx context.Context,
) []resource.StateMover {
	if r.withMoveState == nil {
		return nil
	}

	return r.withMoveState.MoveState(ctx)
}

func (r *AdapterResource) UpgradeState(
	ctx context.Context,
) map[int64]resource.StateUpgrader {
	if r.withUpgradeState == nil {
		return nil
	}

	return r.withUpgradeState.UpgradeState(ctx)
}

func (r *AdapterResource) ValidateConfig(
	ctx context.Context,
	req resource.ValidateConfigRequest,
	resp *resource.ValidateConfigResponse,
) {
	if r.withValidateConfig == nil {
		return
	}

	r.withValidateConfig.ValidateConfig(ctx, req, resp)
}
