// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package group

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/convert"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/errors"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &Resource{}
	_ resource.ResourceWithImportState = &Resource{}
)

func NewResource() resource.Resource {
	return &Resource{}
}

// Resource defines the resource implementation.
type Resource struct {
	configure.ResourceWithMorpheusConfigure
	resource.Resource
}

func (r *Resource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_morpheus_group"
}

func (r *Resource) Schema(
	ctx context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = GroupResourceSchema(ctx)
}

// populate group resource model with current API values
func getGroupAsState(
	ctx context.Context,
	id int64,
	client *sdk.APIClient,
) (GroupModel, diag.Diagnostics) {
	var state GroupModel
	var diags diag.Diagnostics

	g, hresp, err := client.GroupsAPI.GetGroups(ctx, id).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		diags.AddError(
			"populate group resource",
			fmt.Sprintf("group %d GET failed: ", id)+errors.ErrMsg(err, hresp),
		)

		return state, diags
	}

	state.Id = convert.Int64ToType(g.Group.Id)
	state.Name = convert.StrToType(g.Group.Name)
	state.Code = convert.StrToType(g.Group.Code.Get())
	state.Location = convert.StrToType(g.Group.Location.Get())
	state.Labels = convert.StrSliceToSet(g.Group.Labels)

	return state, diags
}

func (r *Resource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan GroupModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := plan.Name.ValueString()
	addGroup := sdk.NewAddGroupsRequestGroup(name)

	var config GroupModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !plan.Code.IsUnknown() {
		addGroup.SetCode(plan.Code.ValueString())
	}

	if !plan.Location.IsUnknown() {
		addGroup.SetLocation(plan.Location.ValueString())
	}

	if !plan.Labels.IsUnknown() {
		var labels []string

		for _, l := range plan.Labels.Elements() {
			v, err := convert.ValueToAny(ctx, l)
			if err != nil {
				resp.Diagnostics.AddError(
					"create group resource",
					"group "+name+": failed to parse label: "+err.Error(),
				)

				return
			}

			labels = append(labels, v.(string))
		}

		addGroup.SetLabels(labels)
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"create group resource",
			"group "+name+": failed to create client: "+err.Error(),
		)

		return
	}

	addGroupReq := sdk.NewAddGroupsRequest(*addGroup)

	group, hresp, err := client.GroupsAPI.AddGroups(ctx).AddGroupsRequest(*addGroupReq).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"create group resource",
			"group "+name+" POST failed: "+errors.ErrMsg(err, hresp),
		)

		return
	}

	if group.GetGroup().Id == nil {
		resp.Diagnostics.AddError(
			"create group resource",
			"group "+name+": id is nil",
		)

		return
	}

	id := *group.GetGroup().Id
	plan.Id = types.Int64Value(id)

	// write id as soon as possible
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, pdiags := getGroupAsState(ctx, id, client)
	if pdiags.HasError() {
		resp.Diagnostics.Append(pdiags...)
		resp.Diagnostics.AddError(
			"create group resource",
			fmt.Sprintf("group %d: failed to read from api", id),
		)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan, state, config GroupModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateGroup := sdk.NewAddGroupsRequestGroupWithDefaults()

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := plan.Name.ValueString()

	updateGroup.SetName(name)

	if !plan.Code.IsNull() {
		updateGroup.SetCode(plan.Code.ValueString())
	}

	if !plan.Location.IsNull() {
		updateGroup.SetLocation(plan.Location.ValueString())
	}

	if plan.Labels.IsNull() {
		updateGroup.SetLabels([]string{})
	} else {
		var labels []string

		for _, l := range plan.Labels.Elements() {
			v, err := convert.ValueToAny(ctx, l)
			if err != nil {
				resp.Diagnostics.AddError(
					"update group resource",
					"group "+name+": failed to parse label: "+err.Error(),
				)

				return
			}

			labels = append(labels, v.(string))
		}

		updateGroup.SetLabels(labels)
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"update group resource",
			"group "+name+": failed to create client: "+err.Error(),
		)

		return
	}

	id := plan.Id.ValueInt64()

	updateGroupReq := sdk.NewAddGroupsRequest(*updateGroup)

	group, hresp, err := client.GroupsAPI.UpdateGroups(ctx, id).
		AddGroupsRequest(*updateGroupReq).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"update group resource",
			"group "+name+" PUT failed: "+errors.ErrMsg(err, hresp),
		)

		return
	}

	if group.GetGroup().Id == nil {
		resp.Diagnostics.AddError(
			"update group resource",
			"group "+name+": id is nil",
		)

		return
	}

	newid := *group.GetGroup().Id
	if newid != id {
		resp.Diagnostics.AddError(
			"update group resource",
			"group "+name+": id mismatch "+fmt.Sprintf("%d != %d", id, newid),
		)

		return
	}

	state, pdiags := getGroupAsState(ctx, newid, client)
	if pdiags.HasError() {
		resp.Diagnostics.Append(pdiags...)
		resp.Diagnostics.AddError(
			"update group resource",
			fmt.Sprintf("group %d: failed to read from api", id),
		)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var plan GroupModel

	diags := req.State.Get(ctx, &plan)
	if diags.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"read group resource",
			"new client call failed with "+err.Error(),
		)

		return
	}

	id := plan.Id.ValueInt64()
	state, pdiags := getGroupAsState(ctx, id, client)
	if pdiags.HasError() {
		resp.Diagnostics.Append(pdiags...)
		resp.Diagnostics.AddError(
			"read group resource",
			fmt.Sprintf("group %d: failed to read from api", id),
		)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var data GroupModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	id := data.Id.ValueInt64()

	client, _ := r.NewClient(ctx)

	_, hresp, err := client.GroupsAPI.RemoveGroups(ctx, id).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"delete group resource",
			fmt.Sprintf("group %d: DELETE failed ", id)+errors.ErrMsg(err, hresp),
		)

		return
	}
}

func (r *Resource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	id, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"import group resource",
			"provided import ID '"+req.ID+"' is invalid (non-number)",
		)

		return
	}

	diags := resp.State.SetAttribute(
		ctx, path.Root("id"), id,
	)
	resp.Diagnostics.Append(diags...)
}
