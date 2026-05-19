package account

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
	_ resource.Resource                = &accountResource{}
	_ resource.ResourceWithConfigure   = &accountResource{}
	_ resource.ResourceWithImportState = &accountResource{}
)

type accountResource struct {
	configure.ResourceWithMorpheusConfigure
}

func NewResource() resource.Resource {
	return &accountResource{}
}

func (r *accountResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_morpheus_account"
}

func (r *accountResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = AccountSchema(ctx)
}

func (r *accountResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var plan accountModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := sdk.AddTenantRequestAccount{
		Name: plan.Name.ValueString(),
	}
	if !plan.Description.IsNull() {
		body.Description = *sdk.NewNullableString(plan.Description.ValueStringPointer())
	}
	if !plan.Subdomain.IsNull() {
		body.Subdomain = *sdk.NewNullableString(plan.Subdomain.ValueStringPointer())
	}
	if !plan.Currency.IsNull() {
		body.Currency = plan.Currency.ValueStringPointer()
	}
	if !plan.RoleID.IsNull() {
		roleID := plan.RoleID.ValueInt64()
		body.Role = &sdk.AddTenantRequestAccountRole{
			Id: &roleID,
		}
	}

	result, httpResp, err := client.TenantsAPI.AddTenant(ctx).AddTenantRequest(sdk.AddTenantRequest{
		Account: body,
	}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpCreate, "account", plan.Name.ValueString(), err, httpResp)
		return
	}

	acct := result.GetAccount()
	mapAddResponseToModel(&plan, &acct)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *accountResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var state accountModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	result, httpResp, err := client.TenantsAPI.GetTenant(ctx, id).Execute()
	if errfmt.IsNotFound(httpResp) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "account", "", err, httpResp)
		return
	}

	acct := result.GetAccount()
	mapGetResponseToModel(&state, &acct)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *accountResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var plan accountModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.ID.ValueInt64()

	body := sdk.UpdateTenantRequestAccount{
		Name: plan.Name.ValueStringPointer(),
	}
	if !plan.Description.IsNull() {
		body.Description = *sdk.NewNullableString(plan.Description.ValueStringPointer())
	}
	if !plan.Subdomain.IsNull() {
		body.Subdomain = *sdk.NewNullableString(plan.Subdomain.ValueStringPointer())
	}
	if !plan.Currency.IsNull() {
		body.Currency = plan.Currency.ValueStringPointer()
	}
	if !plan.RoleID.IsNull() {
		roleID := plan.RoleID.ValueInt64()
		body.Role = &sdk.UpdateTenantRequestAccountRole{
			Id: &roleID,
		}
	}

	_, httpResp, err := client.TenantsAPI.UpdateTenant(ctx, id).UpdateTenantRequest(sdk.UpdateTenantRequest{
		Account: body,
	}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpUpdate, "account", plan.Name.ValueString(), err, httpResp)
		return
	}

	// Re-read to get updated state
	result, httpResp, err := client.TenantsAPI.GetTenant(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "account", plan.Name.ValueString(), err, httpResp)
		return
	}

	acct := result.GetAccount()
	mapGetResponseToModel(&plan, &acct)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *accountResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var state accountModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	_, httpResp, err := client.TenantsAPI.RemoveTenant(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpDelete, "account", "", err, httpResp)
		return
	}
}

func (r *accountResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Could not parse ID %q as integer: %s", req.ID, err))
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func mapAddResponseToModel(model *accountModel, acct *sdk.AddTenant200ResponseAllOfAccount) {
	if acct.Id != nil {
		model.ID = types.Int64Value(*acct.Id)
	}
	if acct.Name != nil {
		model.Name = types.StringValue(*acct.Name)
	}
	if v := acct.Description.Get(); v != nil {
		model.Description = types.StringValue(*v)
	} else {
		model.Description = types.StringNull()
	}
	if acct.Subdomain != nil {
		model.Subdomain = types.StringValue(*acct.Subdomain)
	} else {
		model.Subdomain = types.StringNull()
	}
	if acct.Active != nil {
		model.Active = types.BoolValue(*acct.Active)
	}
	if acct.Currency != nil {
		model.Currency = types.StringValue(*acct.Currency)
	}
}

func mapGetResponseToModel(model *accountModel, acct *sdk.GetTenant200ResponseAccount) {
	if acct.Id != nil {
		model.ID = types.Int64Value(*acct.Id)
	}
	if acct.Name != nil {
		model.Name = types.StringValue(*acct.Name)
	}
	if v := acct.Description.Get(); v != nil {
		model.Description = types.StringValue(*v)
	} else {
		model.Description = types.StringNull()
	}
	if acct.Subdomain != nil {
		model.Subdomain = types.StringValue(*acct.Subdomain)
	} else {
		model.Subdomain = types.StringNull()
	}
	if acct.Active != nil {
		model.Active = types.BoolValue(*acct.Active)
	}
	if acct.Currency != nil {
		model.Currency = types.StringValue(*acct.Currency)
	}
	if acct.Role != nil && acct.Role.Id != nil {
		model.RoleID = types.Int64Value(*acct.Role.Id)
	} else {
		model.RoleID = types.Int64Null()
	}
}
