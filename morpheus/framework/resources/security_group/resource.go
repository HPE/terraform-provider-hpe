package security_group

import (
	"context"
	"fmt"
	"strconv"

	sdk "github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
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

func (r *securityGroupResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_morpheus_security_group"
}

func (r *securityGroupResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = SecurityGroupResourceSchema(ctx)
}

func (r *securityGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var plan SecurityGroupModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := sdk.NewAddSecurityGroupsRequestSecurityGroupWithDefaults()
	body.Name = plan.Name.ValueString()
	if !plan.CloudId.IsNull() && !plan.CloudId.IsUnknown() {
		body.ZoneId = plan.CloudId.ValueInt64()
	}
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		body.Description = plan.Description.ValueStringPointer()
	}
	if !plan.Active.IsNull() && !plan.Active.IsUnknown() {
		body.Active = plan.Active.ValueBoolPointer()
	}
	if !plan.Visibility.IsNull() && !plan.Visibility.IsUnknown() {
		body.Visibility = plan.Visibility.ValueStringPointer()
	}

	// TODO: Add network_server_id (Optional, create-only) to schema and set body.NetworkServerId here.
	// Use case: HVM/Standard clouds with multiple network integrations (e.g. both NSX-T and another
	// network server). When cloud_id alone is insufficient to disambiguate which network server should
	// own the security group, network_server_id lets the user target a specific one. Not needed for
	// NSX-T clouds (where cloud_id automatically resolves to the single network server) or Azure.
	// The field is create-only (not updatable) and not returned in the GET response, making it a
	// WriteOnly attribute candidate. SDK field: AddSecurityGroupsRequestSecurityGroup.NetworkServerId.

	// Tenant permissions
	if !plan.TenantIds.IsNull() && !plan.TenantIds.IsUnknown() {
		var tenantIDs []int64
		resp.Diagnostics.Append(plan.TenantIds.ElementsAs(ctx, &tenantIDs, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		body.TenantPermissions = &sdk.AddSecurityGroupsRequestSecurityGroupTenantPermissions{
			Accounts: tenantIDs,
		}
	}

	// Resource permissions
	if !plan.ResourcePermissionGroupsAll.IsNull() && !plan.ResourcePermissionGroupsAll.IsUnknown() {
		rp := &sdk.AddSecurityGroupsRequestSecurityGroupResourcePermissions{
			All: plan.ResourcePermissionGroupsAll.ValueBoolPointer(),
		}
		if !plan.ResourcePermissionGroupIds.IsNull() && !plan.ResourcePermissionGroupIds.IsUnknown() {
			var groupIDs []int64
			resp.Diagnostics.Append(plan.ResourcePermissionGroupIds.ElementsAs(ctx, &groupIDs, false)...)
			if resp.Diagnostics.HasError() {
				return
			}
			sites := make([]sdk.UpdateCloudFoldersRequestFolderResourcePermissionsSitesInner, len(groupIDs))
			for i, gid := range groupIDs {
				id := gid
				sites[i] = sdk.UpdateCloudFoldersRequestFolderResourcePermissionsSitesInner{Id: &id}
			}
			rp.Sites = sites
		}
		body.ResourcePermissions = rp
	}

	result, httpResp, err := client.SecurityGroupsAPI.AddSecurityGroups(ctx).
		AddSecurityGroupsRequest(sdk.AddSecurityGroupsRequest{
			SecurityGroup: *body,
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

	var state SecurityGroupModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Id.ValueInt64()

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

	var plan SecurityGroupModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.Id.ValueInt64()

	body := sdk.NewUpdateSecurityGroupsRequestSecurityGroupWithDefaults()
	body.Name = plan.Name.ValueStringPointer()
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		body.Description = plan.Description.ValueStringPointer()
	}
	if !plan.Active.IsNull() && !plan.Active.IsUnknown() {
		body.Active = plan.Active.ValueBoolPointer()
	}
	if !plan.Visibility.IsNull() && !plan.Visibility.IsUnknown() {
		body.Visibility = plan.Visibility.ValueStringPointer()
	}

	// Tenant permissions — always send to avoid perpetual diff when user removes tenant_ids from config.
	if !plan.TenantIds.IsNull() && !plan.TenantIds.IsUnknown() {
		var tenantIDs []int64
		resp.Diagnostics.Append(plan.TenantIds.ElementsAs(ctx, &tenantIDs, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		body.TenantPermissions = &sdk.UpdateSecurityGroupsRequestSecurityGroupTenantPermissions{
			Accounts: tenantIDs,
		}
	} else {
		body.TenantPermissions = &sdk.UpdateSecurityGroupsRequestSecurityGroupTenantPermissions{
			Accounts: []int64{},
		}
	}

	// Resource permissions
	if !plan.ResourcePermissionGroupsAll.IsNull() && !plan.ResourcePermissionGroupsAll.IsUnknown() {
		rp := &sdk.UpdateSecurityGroupsRequestSecurityGroupResourcePermissions{
			All: plan.ResourcePermissionGroupsAll.ValueBoolPointer(),
		}
		if !plan.ResourcePermissionGroupIds.IsNull() && !plan.ResourcePermissionGroupIds.IsUnknown() {
			var groupIDs []int64
			resp.Diagnostics.Append(plan.ResourcePermissionGroupIds.ElementsAs(ctx, &groupIDs, false)...)
			if resp.Diagnostics.HasError() {
				return
			}
			sites := make([]sdk.UpdateCloudFoldersRequestFolderResourcePermissionsSitesInner, len(groupIDs))
			for i, gid := range groupIDs {
				id := gid
				sites[i] = sdk.UpdateCloudFoldersRequestFolderResourcePermissionsSitesInner{Id: &id}
			}
			rp.Sites = sites
		}
		body.ResourcePermissions = rp
	}

	result, httpResp, err := client.SecurityGroupsAPI.UpdateSecurityGroups(ctx, id).
		UpdateSecurityGroupsRequest(sdk.UpdateSecurityGroupsRequest{
			SecurityGroup: *body,
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

	var state SecurityGroupModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Id.ValueInt64()

	_, httpResp, err := client.SecurityGroupsAPI.RemoveSecurityGroups(ctx, id).Execute()
	if errfmt.IsNotFound(httpResp) {
		return
	}
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpDelete, "security_group", "", err, httpResp)

		return
	}
}

func (r *securityGroupResource) ImportState(
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

func mapCreateResponseToModel(
	model *SecurityGroupModel,
	sg *sdk.AddSecurityGroups200ResponseSecurityGroup,
) {
	model.Id = convert.Int64ToType(sg.Id)
	model.Name = convert.StrToType(sg.Name)
	if sg.Description.IsSet() {
		model.Description = convert.StrToType(sg.Description.Get())
	} else {
		model.Description = types.StringNull()
	}
	model.Active = convert.BoolToType(sg.Active)
	model.Visibility = convert.StrToType(sg.Visibility)
	zone := sg.GetZone()
	model.CloudId = convert.Int64ToType(zone.Id)

	// Tenants
	if len(sg.Tenants) > 0 {
		tenantValues := make([]attr.Value, 0, len(sg.Tenants))
		for _, t := range sg.Tenants {
			if t.Id != nil {
				tenantValues = append(tenantValues, types.Int64Value(*t.Id))
			}
		}
		model.TenantIds, _ = types.SetValue(types.Int64Type, tenantValues)
	} else {
		model.TenantIds = types.SetNull(types.Int64Type)
	}

	// Resource permissions
	if sg.ResourcePermission != nil {
		model.ResourcePermissionGroupsAll = convert.BoolToType(sg.ResourcePermission.All)
		model.ResourcePermissionGroupIds = extractGroupIDsFromCreateSites(sg.ResourcePermission.Sites)
	} else {
		model.ResourcePermissionGroupsAll = types.BoolNull()
		model.ResourcePermissionGroupIds = types.SetNull(types.Int64Type)
	}
}

func mapResponseToModel(
	model *SecurityGroupModel,
	sg *sdk.GetSecurityGroups200ResponseSecurityGroup,
) {
	model.Id = convert.Int64ToType(sg.Id)
	model.Name = convert.StrToType(sg.Name)
	if sg.Description.IsSet() {
		model.Description = convert.StrToType(sg.Description.Get())
	} else {
		model.Description = types.StringNull()
	}
	model.Active = convert.BoolToType(sg.Active)
	model.Visibility = convert.StrToType(sg.Visibility)
	zone := sg.GetZone()
	model.CloudId = convert.Int64ToType(zone.Id)

	// Tenants
	if len(sg.Tenants) > 0 {
		tenantValues := make([]attr.Value, 0, len(sg.Tenants))
		for _, t := range sg.Tenants {
			if t.Id != nil {
				tenantValues = append(tenantValues, types.Int64Value(*t.Id))
			}
		}
		model.TenantIds, _ = types.SetValue(types.Int64Type, tenantValues)
	} else {
		model.TenantIds = types.SetNull(types.Int64Type)
	}

	// Resource permissions
	if sg.ResourcePermission != nil {
		model.ResourcePermissionGroupsAll = convert.BoolToType(sg.ResourcePermission.All)
		model.ResourcePermissionGroupIds = extractGroupIDsFromCreateSites(sg.ResourcePermission.Sites)
	} else {
		model.ResourcePermissionGroupsAll = types.BoolNull()
		model.ResourcePermissionGroupIds = types.SetNull(types.Int64Type)
	}
}

func mapUpdateResponseToModel(
	model *SecurityGroupModel,
	sg *sdk.UpdateSecurityGroups200ResponseSecurityGroup,
) {
	model.Id = convert.Int64ToType(sg.Id)
	model.Name = convert.StrToType(sg.Name)
	if sg.Description.IsSet() {
		model.Description = convert.StrToType(sg.Description.Get())
	} else {
		model.Description = types.StringNull()
	}
	model.Active = convert.BoolToType(sg.Active)
	model.Visibility = convert.StrToType(sg.Visibility)
	zone := sg.GetZone()
	model.CloudId = convert.Int64ToType(zone.Id)

	// Tenants
	if len(sg.Tenants) > 0 {
		tenantValues := make([]attr.Value, 0, len(sg.Tenants))
		for _, t := range sg.Tenants {
			if t.Id != nil {
				tenantValues = append(tenantValues, types.Int64Value(*t.Id))
			}
		}
		model.TenantIds, _ = types.SetValue(types.Int64Type, tenantValues)
	} else {
		model.TenantIds = types.SetNull(types.Int64Type)
	}

	// Resource permissions
	if sg.ResourcePermission != nil {
		model.ResourcePermissionGroupsAll = convert.BoolToType(sg.ResourcePermission.All)
		model.ResourcePermissionGroupIds = extractGroupIDsFromUpdateSites(sg.ResourcePermission.Sites)
	} else {
		model.ResourcePermissionGroupsAll = types.BoolNull()
		model.ResourcePermissionGroupIds = types.SetNull(types.Int64Type)
	}
}

func extractGroupIDsFromCreateSites(
	sites []sdk.AddSecurityGroups200ResponseSecurityGroupAllOfResourcePermissionSitesInner,
) types.Set {
	if len(sites) == 0 {
		return types.SetNull(types.Int64Type)
	}

	groupValues := make([]attr.Value, 0, len(sites))
	for _, site := range sites {
		if site.Id != nil {
			groupValues = append(groupValues, types.Int64Value(*site.Id))
		}
	}

	if len(groupValues) == 0 {
		return types.SetNull(types.Int64Type)
	}

	result, _ := types.SetValue(types.Int64Type, groupValues)

	return result
}

func extractGroupIDsFromUpdateSites(
	sites []sdk.UpdateSecurityGroups200ResponseSecurityGroupAllOfResourcePermissionSitesInner,
) types.Set {
	if len(sites) == 0 {
		return types.SetNull(types.Int64Type)
	}

	groupValues := make([]attr.Value, 0, len(sites))
	for _, site := range sites {
		if site.Id != nil {
			groupValues = append(groupValues, types.Int64Value(*site.Id))
		}
	}

	if len(groupValues) == 0 {
		return types.SetNull(types.Int64Type)
	}

	result, _ := types.SetValue(types.Int64Type, groupValues)

	return result
}
