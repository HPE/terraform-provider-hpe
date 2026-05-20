package security_group

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
	_ resource.Resource                = &securityGroupResource{}
	_ resource.ResourceWithConfigure   = &securityGroupResource{}
	_ resource.ResourceWithImportState = &securityGroupResource{}
)

type securityGroupResource struct {
	configure.ResourceWithMorpheusConfigure
}

func NewResource() resource.Resource {
	return &securityGroupResource{}
}

func (r *securityGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_morpheus_security_group"
}

func (r *securityGroupResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = SecurityGroupSchema(ctx)
}

func (r *securityGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var plan securityGroupModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := sdk.AddSecurityGroupsRequestSecurityGroup{
		Name: plan.Name.ValueString(),
	}
	if !plan.ZoneID.IsNull() && !plan.ZoneID.IsUnknown() {
		body.ZoneId = plan.ZoneID.ValueInt64()
	}
	if !plan.Description.IsNull() {
		body.Description = plan.Description.ValueStringPointer()
	}
	if !plan.Active.IsNull() {
		body.Active = plan.Active.ValueBoolPointer()
	}

	result, httpResp, err := client.SecurityGroupsAPI.AddSecurityGroups(ctx).AddSecurityGroupsRequest(sdk.AddSecurityGroupsRequest{
		SecurityGroup: body,
	}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpCreate, "security_group", plan.Name.ValueString(), err, httpResp)
		return
	}

	sg := result.GetSecurityGroup()
	mapCreateResponseToModel(&plan, &sg)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *securityGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var state securityGroupModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	result, httpResp, err := client.SecurityGroupsAPI.GetSecurityGroups(ctx, id).Execute()
	if errfmt.IsNotFound(httpResp) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "security_group", "", err, httpResp)
		return
	}

	sg := result.GetSecurityGroup()
	mapResponseToModel(&state, &sg)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *securityGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var plan securityGroupModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.ID.ValueInt64()

	body := sdk.UpdateSecurityGroupsRequestSecurityGroup{
		Name: plan.Name.ValueStringPointer(),
	}
	if !plan.Description.IsNull() {
		body.Description = plan.Description.ValueStringPointer()
	}
	if !plan.Active.IsNull() {
		body.Active = plan.Active.ValueBoolPointer()
	}

	result, httpResp, err := client.SecurityGroupsAPI.UpdateSecurityGroups(ctx, id).UpdateSecurityGroupsRequest(sdk.UpdateSecurityGroupsRequest{
		SecurityGroup: body,
	}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpUpdate, "security_group", plan.Name.ValueString(), err, httpResp)
		return
	}

	sg := result.GetSecurityGroup()
	mapUpdateResponseToModel(&plan, &sg)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *securityGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var state securityGroupModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	_, httpResp, err := client.SecurityGroupsAPI.RemoveSecurityGroups(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpDelete, "security_group", "", err, httpResp)
		return
	}
}

func (r *securityGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Could not parse ID %q as integer: %s", req.ID, err))
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func mapCreateResponseToModel(model *securityGroupModel, sg *sdk.AddSecurityGroups200ResponseSecurityGroup) {
	if sg.Id != nil {
		model.ID = types.Int64Value(*sg.Id)
	}
	if sg.Name != nil {
		model.Name = types.StringValue(*sg.Name)
	}
	if sg.Description.IsSet() && sg.Description.Get() != nil {
		model.Description = types.StringValue(*sg.Description.Get())
	} else {
		model.Description = types.StringNull()
	}
	if sg.Active != nil {
		model.Active = types.BoolValue(*sg.Active)
	}
	if sg.Visibility != nil {
		model.Visibility = types.StringValue(*sg.Visibility)
	}
	zone := sg.GetZone()
	if zone.Id != nil {
		model.ZoneID = types.Int64Value(*zone.Id)
	}
}

func mapResponseToModel(model *securityGroupModel, sg *sdk.GetSecurityGroups200ResponseSecurityGroup) {
	if sg.Id != nil {
		model.ID = types.Int64Value(*sg.Id)
	}
	if sg.Name != nil {
		model.Name = types.StringValue(*sg.Name)
	}
	if sg.Description.IsSet() && sg.Description.Get() != nil {
		model.Description = types.StringValue(*sg.Description.Get())
	} else {
		model.Description = types.StringNull()
	}
	if sg.Active != nil {
		model.Active = types.BoolValue(*sg.Active)
	}
	if sg.Visibility != nil {
		model.Visibility = types.StringValue(*sg.Visibility)
	}
	zone := sg.GetZone()
	if zone.Id != nil {
		model.ZoneID = types.Int64Value(*zone.Id)
	}
}

func mapUpdateResponseToModel(model *securityGroupModel, sg *sdk.UpdateSecurityGroups200ResponseSecurityGroup) {
	if sg.Id != nil {
		model.ID = types.Int64Value(*sg.Id)
	}
	if sg.Name != nil {
		model.Name = types.StringValue(*sg.Name)
	}
	if sg.Description.IsSet() && sg.Description.Get() != nil {
		model.Description = types.StringValue(*sg.Description.Get())
	} else {
		model.Description = types.StringNull()
	}
	if sg.Active != nil {
		model.Active = types.BoolValue(*sg.Active)
	}
	if sg.Visibility != nil {
		model.Visibility = types.StringValue(*sg.Visibility)
	}
	zone := sg.GetZone()
	if zone.Id != nil {
		model.ZoneID = types.Int64Value(*zone.Id)
	}
}
