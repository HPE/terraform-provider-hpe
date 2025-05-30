// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package role

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/convert"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/errors"
	"github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource = &Resource{}
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
	resp.TypeName = req.ProviderTypeName + "_morpheus_role"
}

func (r *Resource) Schema(
	ctx context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = RoleResourceSchema(ctx)
}

// populate role resource model with current API values
func getRoleAsState(
	ctx context.Context,
	id int64,
	client *sdk.APIClient,
) (RoleModel, diag.Diagnostics) {
	var state RoleModel
	var diags diag.Diagnostics

	r, hresp, err := client.RolesAPI.GetRole(ctx, id).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		diags.AddError(
			"populate role resource",
			fmt.Sprintf("role %d GET failed: ", id)+errors.ErrMsg(err, hresp),
		)

		return state, diags
	}

	state.Id = convert.Int64ToType(r.Role.Id)
	state.Name = convert.StrToType(r.Role.Name)
	state.Description = convert.StrToType(r.Role.Description)
	state.Multitenant = convert.BoolToType(r.Role.Multitenant)
	state.RoleType = convert.StrToType(r.Role.RoleType)

	return state, diags
}

func (r *Resource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan RoleModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRole := sdk.NewAddRolesRequestRoleWithDefaults()

	// required
	name := plan.Name.ValueString()
	addRole.SetAuthority(name)

	// optional
	if !plan.Description.IsUnknown() {
		addRole.SetDescription(plan.Description.ValueString())
	}
	if !plan.Multitenant.IsUnknown() {
		addRole.SetMultitenant(plan.Multitenant.ValueBool())
	}
	if !plan.RoleType.IsUnknown() {
		if plan.RoleType.ValueString() != "user" {
			resp.Diagnostics.AddError(
				"create role resource",
				"role "+name+": currently only 'user' role_type is supported",
			)

			return
		}
		addRole.SetRoleType(plan.RoleType.ValueString())
	}

	addRoleReq := sdk.NewAddRolesRequest(*addRole)

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"create role resource",
			"role "+name+": failed to create client: "+err.Error(),
		)

		return
	}

	role, hresp, err := client.RolesAPI.AddRoles(ctx).
		AddRolesRequest(*addRoleReq).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"create role resource",
			"role "+name+" POST failed: "+errors.ErrMsg(err, hresp),
		)

		return
	}

	if role.GetRole().Id == nil {
		resp.Diagnostics.AddError(
			"create role resource",
			"role "+name+": id is nil",
		)

		return
	}

	id := *role.GetRole().Id
	plan.Id = types.Int64Value(id)

	// write id as soon as possible
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, pdiags := getRoleAsState(ctx, id, client)
	if pdiags.HasError() {
		resp.Diagnostics.Append(pdiags...)
		resp.Diagnostics.AddError(
			"create role resource",
			fmt.Sprintf("role %d: failed to read from api", id),
		)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var plan RoleModel

	diags := req.State.Get(ctx, &plan)
	if diags.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"read role resource",
			"new client call failed with "+err.Error(),
		)

		return
	}

	id := plan.Id.ValueInt64()
	state, pdiags := getRoleAsState(ctx, id, client)
	if pdiags.HasError() {
		resp.Diagnostics.Append(pdiags...)
		resp.Diagnostics.AddError(
			"read role resource",
			fmt.Sprintf("role %d: failed to read from api", id),
		)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *Resource) Update(
	_ context.Context,
	_ resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	resp.Diagnostics.AddError(
		"update role resource",
		"update of 'role' resources has not been implemented",
	)
}

func (r *Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var data RoleModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	id := data.Id.ValueInt64()
	client, _ := r.NewClient(ctx)
	_, hresp, err := client.RolesAPI.DeleteRole(ctx, id).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"delete role resource",
			fmt.Sprintf("role %d: DELETE failed ", id)+errors.ErrMsg(err, hresp),
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
			"import role resource",
			"provided import ID '"+req.ID+"' is invalid (non-number)",
		)

		return
	}

	diags := resp.State.SetAttribute(ctx, path.Root("id"), id)

	resp.Diagnostics.Append(diags...)
}
