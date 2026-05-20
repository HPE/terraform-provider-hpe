package credential

import (
	"context"
	"fmt"
	"strconv"

	sdk "github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
)

var (
	_ resource.Resource                = &credentialResource{}
	_ resource.ResourceWithConfigure   = &credentialResource{}
	_ resource.ResourceWithImportState = &credentialResource{}
)

type credentialResource struct {
	configure.ResourceWithMorpheusConfigure
}

func NewResource() resource.Resource {
	return &credentialResource{}
}

func (r *credentialResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_morpheus_credential"
}

func (r *credentialResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = CredentialSchema(ctx)
}

func (r *credentialResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var plan credentialModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	inner := sdk.AddCredentialsRequestCredentialOneOf{
		Type: plan.Type.ValueString(),
		Name: plan.Name.ValueString(),
	}
	if !plan.Description.IsNull() {
		inner.Description = plan.Description.ValueStringPointer()
	}
	if !plan.Enabled.IsNull() {
		inner.Enabled = plan.Enabled.ValueBoolPointer()
	}
	if !plan.Username.IsNull() {
		inner.Username = plan.Username.ValueString()
	}
	if !plan.Password.IsNull() {
		inner.Password = plan.Password.ValueString()
	}
	if !plan.IntegrationID.IsNull() {
		integrationID := plan.IntegrationID.ValueInt64()
		integrationIdWrapper := sdk.Int64AsAddCredentialsRequestCredentialOneOfIntegrationId(&integrationID)
		inner.Integration = &sdk.AddCredentialsRequestCredentialOneOfIntegration{
			Id: &integrationIdWrapper,
		}
	}

	body := sdk.AddCredentialsRequestCredentialOneOfAsAddCredentialsRequestCredential(&inner)

	result, httpResp, err := client.CredentialsAPI.AddCredentials(ctx).AddCredentialsRequest(sdk.AddCredentialsRequest{
		Credential: body,
	}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpCreate, "credential", plan.Name.ValueString(), err, httpResp)

		return
	}

	cred := result.GetCredential()
	mapAddResponseToModel(&plan, &cred)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *credentialResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var state credentialModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	result, httpResp, err := client.CredentialsAPI.GetCredentials(ctx, id).Execute()
	if errfmt.IsNotFound(httpResp) {
		resp.State.RemoveResource(ctx)

		return
	}
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "credential", "", err, httpResp)

		return
	}

	cred := result.GetCredential()
	mapGetResponseToModel(&state, &cred)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *credentialResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var plan credentialModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.ID.ValueInt64()

	inner := sdk.UpdateCredentialsRequestCredentialOneOf{
		Type: plan.Type.ValueString(),
		Name: plan.Name.ValueString(),
	}
	if !plan.Description.IsNull() {
		inner.Description = plan.Description.ValueStringPointer()
	}
	if !plan.Enabled.IsNull() {
		inner.Enabled = plan.Enabled.ValueBoolPointer()
	}
	if !plan.Username.IsNull() {
		inner.Username = plan.Username.ValueString()
	}
	if !plan.Password.IsNull() {
		inner.Password = plan.Password.ValueString()
	}
	if !plan.IntegrationID.IsNull() {
		integrationID := plan.IntegrationID.ValueInt64()
		integrationIdWrapper := sdk.Int64AsUpdateCredentialsRequestCredentialOneOfIntegrationId(&integrationID)
		inner.Integration = &sdk.UpdateCredentialsRequestCredentialOneOfIntegration{
			Id: &integrationIdWrapper,
		}
	}

	body := sdk.UpdateCredentialsRequestCredentialOneOfAsUpdateCredentialsRequestCredential(&inner)

	_, httpResp, err := client.CredentialsAPI.UpdateCredentials(ctx, id).
		UpdateCredentialsRequest(sdk.UpdateCredentialsRequest{
			Credential: body,
		}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpUpdate, "credential", plan.Name.ValueString(), err, httpResp)

		return
	}

	// Re-read
	getResult, httpResp, err := client.CredentialsAPI.GetCredentials(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "credential", "", err, httpResp)

		return
	}

	cred := getResult.GetCredential()
	mapGetResponseToModel(&plan, &cred)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *credentialResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var state credentialModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	_, httpResp, err := client.CredentialsAPI.RemoveCredentials(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpDelete, "credential", "", err, httpResp)

		return
	}
}

func (r *credentialResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Could not parse ID %q as integer: %s", req.ID, err))

		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func mapAddResponseToModel(model *credentialModel, cred *sdk.AddCredentials200ResponseAllOfCredential) {
	if cred.Id != nil {
		model.ID = types.Int64Value(*cred.Id)
	}
	if cred.Name != nil {
		model.Name = types.StringValue(*cred.Name)
	}
	if cred.Description.IsSet() && cred.Description.Get() != nil {
		model.Description = types.StringValue(*cred.Description.Get())
	} else {
		model.Description = types.StringNull()
	}
	if cred.Enabled != nil {
		model.Enabled = types.BoolValue(*cred.Enabled)
	}
}

func mapGetResponseToModel(model *credentialModel, cred *sdk.GetCredentials200ResponseCredential) {
	if cred.Id != nil {
		model.ID = types.Int64Value(*cred.Id)
	}
	if cred.Name != nil {
		model.Name = types.StringValue(*cred.Name)
	}
	if cred.Description.IsSet() && cred.Description.Get() != nil {
		model.Description = types.StringValue(*cred.Description.Get())
	} else {
		model.Description = types.StringNull()
	}
	if cred.Enabled != nil {
		model.Enabled = types.BoolValue(*cred.Enabled)
	}
	if cred.Type != nil && cred.Type.Code != nil {
		model.Type = types.StringValue(*cred.Type.Code)
	}
}
