package user_source

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
	_ resource.Resource                = &userSourceResource{}
	_ resource.ResourceWithConfigure   = &userSourceResource{}
	_ resource.ResourceWithImportState = &userSourceResource{}
)

type userSourceResource struct {
	configure.ResourceWithMorpheusConfigure
}

func NewResource() resource.Resource {
	return &userSourceResource{}
}

func (r *userSourceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_morpheus_user_source"
}

func (r *userSourceResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = UserSourceSchema(ctx)
}

func (r *userSourceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var plan userSourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	userSource := *sdk.NewAddIdentitySourcesRequestUserSource(plan.Name.ValueString(), plan.Type.ValueString())
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		userSource.Description = plan.Description.ValueStringPointer()
	}
	if !plan.DefaultAccountRoleId.IsNull() && !plan.DefaultAccountRoleId.IsUnknown() {
		userSource.DefaultAccountRole = &sdk.AddIdentitySourcesRequestUserSourceDefaultAccountRole{
			Id: plan.DefaultAccountRoleId.ValueInt64(),
		}
	}

	createResp, httpResp, err := client.IdentitySourcesAPI.AddIdentitySources(ctx).
		AccountId(plan.AccountId.ValueInt64()).
		AddIdentitySourcesRequest(sdk.AddIdentitySourcesRequest{
			UserSource: userSource,
		}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpCreate, "user_source", plan.Name.ValueString(), err, httpResp)
		return
	}

	// Extract ID from response AdditionalProperties
	var id int64
	if createResp != nil && createResp.AdditionalProperties != nil {
		if idVal, ok := createResp.AdditionalProperties["id"]; ok {
			switch v := idVal.(type) {
			case float64:
				id = int64(v)
			case int64:
				id = v
			}
		}
	}
	if id == 0 {
		resp.Diagnostics.AddError("Create Error", "Could not determine ID from create response")
		return
	}

	// Read back
	readResp, httpResp, err := client.IdentitySourcesAPI.GetIdentitySources(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "user_source", "", err, httpResp)
		return
	}

	mapGetUserSourceToModel(&plan, readResp)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *userSourceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var state userSourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	readResp, httpResp, err := client.IdentitySourcesAPI.GetIdentitySources(ctx, id).Execute()
	if errfmt.IsNotFound(httpResp) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "user_source", "", err, httpResp)
		return
	}

	mapGetUserSourceToModel(&state, readResp)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *userSourceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var plan userSourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.ID.ValueInt64()

	userSource := *sdk.NewUpdateIdentitySourcesRequestUserSource(plan.Name.ValueString(), plan.Type.ValueString())
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		userSource.Description = plan.Description.ValueStringPointer()
	}
	if !plan.DefaultAccountRoleId.IsNull() && !plan.DefaultAccountRoleId.IsUnknown() {
		userSource.DefaultAccountRole = &sdk.UpdateIdentitySourcesRequestUserSourceDefaultAccountRole{
			Id: plan.DefaultAccountRoleId.ValueInt64(),
		}
	}

	_, httpResp, err := client.IdentitySourcesAPI.UpdateIdentitySources(ctx, id).
		UpdateIdentitySourcesRequest(sdk.UpdateIdentitySourcesRequest{
			UserSource: userSource,
		}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpUpdate, "user_source", plan.Name.ValueString(), err, httpResp)
		return
	}

	// Read back
	readResp, httpResp, err := client.IdentitySourcesAPI.GetIdentitySources(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "user_source", "", err, httpResp)
		return
	}

	mapGetUserSourceToModel(&plan, readResp)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *userSourceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var state userSourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	_, httpResp, err := client.IdentitySourcesAPI.RemoveIdentitySources(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpDelete, "user_source", "", err, httpResp)
		return
	}
}

func (r *userSourceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Could not parse ID %q as integer: %s", req.ID, err))
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func mapGetUserSourceToModel(model *userSourceModel, resp *sdk.GetIdentitySources200Response) {
	if resp == nil {
		return
	}
	// The response uses a polymorphic UserSource; use AdditionalProperties to extract fields
	us := resp.UserSource
	if us == nil {
		// Try AdditionalProperties on the outer response
		if resp.AdditionalProperties != nil {
			extractUserSourceFromMap(model, resp.AdditionalProperties)
		}
		return
	}

	// Try the first AnyOf variant which has the common fields
	if us.GetIdentitySources200ResponseUserSourceAnyOf != nil {
		anyOf := us.GetIdentitySources200ResponseUserSourceAnyOf
		if anyOf.Id != nil {
			model.ID = types.Int64Value(*anyOf.Id)
		}
		if anyOf.Name != nil {
			model.Name = types.StringValue(*anyOf.Name)
		}
		if anyOf.Type != nil {
			model.Type = types.StringValue(*anyOf.Type)
		}
		if anyOf.Description != nil {
			model.Description = types.StringValue(*anyOf.Description)
		} else if model.Description.IsNull() {
			model.Description = types.StringNull()
		}
		if anyOf.Account != nil && anyOf.Account.Id != nil {
			model.AccountId = types.Int64Value(*anyOf.Account.Id)
		}
		if anyOf.DefaultAccountRole != nil && anyOf.DefaultAccountRole.Id != nil {
			model.DefaultAccountRoleId = types.Int64Value(*anyOf.DefaultAccountRole.Id)
		} else if model.DefaultAccountRoleId.IsNull() || model.DefaultAccountRoleId.IsUnknown() {
			model.DefaultAccountRoleId = types.Int64Null()
		}
	}
}

func extractUserSourceFromMap(model *userSourceModel, m map[string]interface{}) {
	if usMap, ok := m["userSource"].(map[string]interface{}); ok {
		if id, ok := usMap["id"].(float64); ok {
			model.ID = types.Int64Value(int64(id))
		}
		if name, ok := usMap["name"].(string); ok {
			model.Name = types.StringValue(name)
		}
		if t, ok := usMap["type"].(string); ok {
			model.Type = types.StringValue(t)
		}
		if desc, ok := usMap["description"].(string); ok {
			model.Description = types.StringValue(desc)
		}
		if acct, ok := usMap["account"].(map[string]interface{}); ok {
			if id, ok := acct["id"].(float64); ok {
				model.AccountId = types.Int64Value(int64(id))
			}
		}
		if role, ok := usMap["defaultAccountRole"].(map[string]interface{}); ok {
			if id, ok := role["id"].(float64); ok {
				model.DefaultAccountRoleId = types.Int64Value(int64(id))
			}
		}
	}
}
