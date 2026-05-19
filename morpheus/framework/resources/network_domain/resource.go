package network_domain

import (
	"context"
	"fmt"
	"strconv"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	sdk "github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &networkDomainResource{}
	_ resource.ResourceWithConfigure   = &networkDomainResource{}
	_ resource.ResourceWithImportState = &networkDomainResource{}
)

type networkDomainResource struct {
	configure.ResourceWithMorpheusConfigure
}

func NewResource() resource.Resource {
	return &networkDomainResource{}
}

func (r *networkDomainResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_morpheus_network_domain"
}

func (r *networkDomainResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = NetworkDomainSchema(ctx)
}

func (r *networkDomainResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var plan networkDomainModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := sdk.CreateNetworkDomainRequestNetworkDomain{
		Name: plan.Name.ValueStringPointer(),
	}
	if !plan.Description.IsNull() {
		body.Description = plan.Description.ValueStringPointer()
	}
	if !plan.PublicZone.IsNull() {
		body.PublicZone = plan.PublicZone.ValueBoolPointer()
	}
	if !plan.Active.IsNull() {
		body.Active = plan.Active.ValueBoolPointer()
	}
	if !plan.DomainController.IsNull() {
		body.DomainController = plan.DomainController.ValueBoolPointer()
	}
	if !plan.DomainUsername.IsNull() {
		body.DomainUsername = plan.DomainUsername.ValueStringPointer()
	}
	if !plan.DomainPassword.IsNull() {
		body.DomainPassword = plan.DomainPassword.ValueStringPointer()
	}
	if !plan.DcServer.IsNull() {
		body.DcServer = plan.DcServer.ValueStringPointer()
	}

	result, httpResp, err := client.NetworksAPI.CreateNetworkDomain(ctx).CreateNetworkDomainRequest(sdk.CreateNetworkDomainRequest{
		NetworkDomain: &body,
	}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpCreate, "network_domain", plan.Name.ValueString(), err, httpResp)
		return
	}

	domain := result.GetNetworkDomain()
	mapCreateResponseToModel(&plan, &domain)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *networkDomainResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var state networkDomainModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	result, httpResp, err := client.NetworksAPI.GetNetworkDomain(ctx, id).Execute()
	if errfmt.IsNotFound(httpResp) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "network_domain", "", err, httpResp)
		return
	}

	domain := result.GetNetworkDomain()
	mapResponseToModel(&state, &domain)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *networkDomainResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var plan networkDomainModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.ID.ValueInt64()

	body := sdk.UpdateNetworkDomainRequestNetworkDomain{
		Name: plan.Name.ValueStringPointer(),
	}
	if !plan.Description.IsNull() {
		body.Description = plan.Description.ValueStringPointer()
	}
	if !plan.PublicZone.IsNull() {
		body.PublicZone = plan.PublicZone.ValueBoolPointer()
	}
	if !plan.Active.IsNull() {
		body.Active = plan.Active.ValueBoolPointer()
	}
	if !plan.DomainController.IsNull() {
		body.DomainController = plan.DomainController.ValueBoolPointer()
	}
	if !plan.DomainUsername.IsNull() {
		body.DomainUsername = plan.DomainUsername.ValueStringPointer()
	}
	if !plan.DomainPassword.IsNull() {
		body.DomainPassword = plan.DomainPassword.ValueStringPointer()
	}
	if !plan.DcServer.IsNull() {
		body.DcServer = plan.DcServer.ValueStringPointer()
	}

	result, httpResp, err := client.NetworksAPI.UpdateNetworkDomain(ctx, id).UpdateNetworkDomainRequest(sdk.UpdateNetworkDomainRequest{
		NetworkDomain: &body,
	}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpUpdate, "network_domain", plan.Name.ValueString(), err, httpResp)
		return
	}

	domain := result.GetNetworkDomain()
	mapResponseToModel(&plan, &domain)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *networkDomainResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var state networkDomainModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	_, httpResp, err := client.NetworksAPI.DeleteNetworkDomain(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpDelete, "network_domain", "", err, httpResp)
		return
	}
}

func (r *networkDomainResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Could not parse ID %q as integer: %s", req.ID, err))
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func mapCreateResponseToModel(model *networkDomainModel, domain *sdk.CreateNetworkDomain200ResponseNetworkDomain) {
	if domain.Id != nil {
		model.ID = types.Int64Value(*domain.Id)
	}
	if domain.Name != nil {
		model.Name = types.StringValue(*domain.Name)
	}
	if domain.Description.IsSet() && domain.Description.Get() != nil {
		model.Description = types.StringValue(*domain.Description.Get())
	} else {
		model.Description = types.StringNull()
	}
	if domain.Active != nil {
		model.Active = types.BoolValue(*domain.Active)
	}
	if domain.PublicZone != nil {
		model.PublicZone = types.BoolValue(*domain.PublicZone)
	}
	if domain.DomainController != nil {
		model.DomainController = types.BoolValue(*domain.DomainController)
	}
	if domain.DomainUsername.IsSet() && domain.DomainUsername.Get() != nil {
		model.DomainUsername = types.StringValue(*domain.DomainUsername.Get())
	}
	if domain.Fqdn.IsSet() && domain.Fqdn.Get() != nil {
		model.Fqdn = types.StringValue(*domain.Fqdn.Get())
	} else {
		model.Fqdn = types.StringNull()
	}
	if domain.Visibility != nil {
		model.Visibility = types.StringValue(*domain.Visibility)
	}
}

func mapResponseToModel(model *networkDomainModel, domain *sdk.GetNetworkDomain200ResponseNetworkDomain) {
	if domain.Id != nil {
		model.ID = types.Int64Value(*domain.Id)
	}
	if domain.Name != nil {
		model.Name = types.StringValue(*domain.Name)
	}
	if domain.Description.IsSet() && domain.Description.Get() != nil {
		model.Description = types.StringValue(*domain.Description.Get())
	} else {
		model.Description = types.StringNull()
	}
	if domain.Active != nil {
		model.Active = types.BoolValue(*domain.Active)
	}
	if domain.PublicZone != nil {
		model.PublicZone = types.BoolValue(*domain.PublicZone)
	}
	if domain.DomainController != nil {
		model.DomainController = types.BoolValue(*domain.DomainController)
	}
	if domain.DomainUsername.IsSet() && domain.DomainUsername.Get() != nil {
		model.DomainUsername = types.StringValue(*domain.DomainUsername.Get())
	}
	if domain.DcServer.IsSet() && domain.DcServer.Get() != nil {
		model.DcServer = types.StringValue(*domain.DcServer.Get())
	}
	if domain.Fqdn.IsSet() && domain.Fqdn.Get() != nil {
		model.Fqdn = types.StringValue(*domain.Fqdn.Get())
	} else {
		model.Fqdn = types.StringNull()
	}
	if domain.Visibility != nil {
		model.Visibility = types.StringValue(*domain.Visibility)
	}
}
