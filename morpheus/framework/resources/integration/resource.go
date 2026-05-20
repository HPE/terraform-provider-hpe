package integration

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
	_ resource.Resource                = &integrationResource{}
	_ resource.ResourceWithConfigure   = &integrationResource{}
	_ resource.ResourceWithImportState = &integrationResource{}
)

type integrationResource struct {
	configure.ResourceWithMorpheusConfigure
}

func NewResource() resource.Resource {
	return &integrationResource{}
}

func (r *integrationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_morpheus_integration"
}

func (r *integrationResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = IntegrationSchema(ctx)
}

func (r *integrationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var plan integrationModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	integration := sdk.AddIntegrationsRequestOneOfIntegration{
		Name: plan.Name.ValueString(),
		Type: plan.Type.ValueString(),
	}
	if !plan.Enabled.IsNull() {
		integration.Enabled = plan.Enabled.ValueBoolPointer()
	}
	if !plan.URL.IsNull() {
		integration.AdditionalProperties = make(map[string]interface{})
		integration.AdditionalProperties["url"] = plan.URL.ValueString()
	}
	if !plan.Username.IsNull() {
		if integration.AdditionalProperties == nil {
			integration.AdditionalProperties = make(map[string]interface{})
		}
		integration.AdditionalProperties["username"] = plan.Username.ValueString()
	}
	if !plan.Password.IsNull() {
		if integration.AdditionalProperties == nil {
			integration.AdditionalProperties = make(map[string]interface{})
		}
		integration.AdditionalProperties["password"] = plan.Password.ValueString()
	}

	body := sdk.AddIntegrationsRequestOneOfAsAddIntegrationsRequest(&sdk.AddIntegrationsRequestOneOf{
		Integration: integration,
	})

	result, httpResp, err := client.IntegrationsAPI.AddIntegrations(ctx).AddIntegrationsRequest(body).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpCreate, "integration", plan.Name.ValueString(), err, httpResp)
		return
	}

	mapAddResponseToModel(&plan, result)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *integrationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var state integrationModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	result, httpResp, err := client.IntegrationsAPI.GetIntegrations(ctx, id).Execute()
	if errfmt.IsNotFound(httpResp) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "integration", "", err, httpResp)
		return
	}

	mapGetResponseToModel(&state, result)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *integrationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var plan integrationModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.ID.ValueInt64()

	integration := sdk.UpdateIntegrationsRequestOneOfIntegration{
		Name: plan.Name.ValueString(),
		Type: plan.Type.ValueString(),
	}
	if !plan.Enabled.IsNull() {
		integration.Enabled = plan.Enabled.ValueBoolPointer()
	}
	if !plan.URL.IsNull() {
		integration.AdditionalProperties = make(map[string]interface{})
		integration.AdditionalProperties["url"] = plan.URL.ValueString()
	}
	if !plan.Username.IsNull() {
		if integration.AdditionalProperties == nil {
			integration.AdditionalProperties = make(map[string]interface{})
		}
		integration.AdditionalProperties["username"] = plan.Username.ValueString()
	}
	if !plan.Password.IsNull() {
		if integration.AdditionalProperties == nil {
			integration.AdditionalProperties = make(map[string]interface{})
		}
		integration.AdditionalProperties["password"] = plan.Password.ValueString()
	}

	body := sdk.UpdateIntegrationsRequestOneOfAsUpdateIntegrationsRequest(&sdk.UpdateIntegrationsRequestOneOf{
		Integration: integration,
	})

	_, httpResp, err := client.IntegrationsAPI.UpdateIntegrations(ctx, id).UpdateIntegrationsRequest(body).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpUpdate, "integration", plan.Name.ValueString(), err, httpResp)
		return
	}

	// Re-read to get current state
	result, httpResp, err := client.IntegrationsAPI.GetIntegrations(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "integration", plan.Name.ValueString(), err, httpResp)
		return
	}

	mapGetResponseToModel(&plan, result)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *integrationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var state integrationModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	_, httpResp, err := client.IntegrationsAPI.RemoveIntegrations(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpDelete, "integration", "", err, httpResp)
		return
	}
}

func (r *integrationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Could not parse ID %q as integer: %s", req.ID, err))
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func mapAddResponseToModel(model *integrationModel, result *sdk.AddIntegrations200Response) {
	integration := result.GetIntegration()
	if oneOf := integration.AddIntegrations200ResponseAllOfIntegrationOneOf; oneOf != nil {
		if oneOf.Id != nil {
			model.ID = types.Int64Value(*oneOf.Id)
		}
		if oneOf.Name != nil {
			model.Name = types.StringValue(*oneOf.Name)
		}
		if oneOf.Type != nil {
			model.Type = types.StringValue(*oneOf.Type)
		}
		if oneOf.Enabled != nil {
			model.Enabled = types.BoolValue(*oneOf.Enabled)
		}
		if oneOf.Url != nil {
			model.URL = types.StringValue(*oneOf.Url)
		} else {
			model.URL = types.StringNull()
		}
	}
}

func mapGetResponseToModel(model *integrationModel, result *sdk.GetIntegrations200Response) {
	integration := result.GetIntegration()
	if oneOf := integration.GetIntegrations200ResponseAllOfIntegrationOneOf; oneOf != nil {
		if oneOf.Id != nil {
			model.ID = types.Int64Value(*oneOf.Id)
		}
		if oneOf.Name != nil {
			model.Name = types.StringValue(*oneOf.Name)
		}
		if oneOf.Type != nil {
			model.Type = types.StringValue(*oneOf.Type)
		}
		if oneOf.Enabled != nil {
			model.Enabled = types.BoolValue(*oneOf.Enabled)
		}
		if oneOf.Url != nil {
			model.URL = types.StringValue(*oneOf.Url)
		} else {
			model.URL = types.StringNull()
		}
	}
}
