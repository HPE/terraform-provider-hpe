package securitygroup

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/cleanup"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
	"github.com/HPE/terraform-provider-hpe/utils/schemadefaults"
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
	resp.TypeName = req.ProviderTypeName + "_" + "security_group"
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

	body := &sdk.AddSecurityGroupsRequestSecurityGroup{}
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
	//
	// These are also sent here for completeness, but the create endpoint does not
	// act on them (see applyPlanAfterCreate below), so the authoritative write
	// happens in the follow-up update after the create.
	rp, diags := resourcePermissionsFromPlan(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if rp != nil {
		body.ResourcePermissions = &sdk.AddSecurityGroupsRequestSecurityGroupResourcePermissions{
			All:   rp.All,
			Sites: rp.Sites,
		}
	}

	result, httpResp, err := client.SecurityGroupsAPI.AddSecurityGroups(ctx).
		AddSecurityGroupsRequest(sdk.AddSecurityGroupsRequest{
			SecurityGroup: *body,
		}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpCreate, "security_group", plan.Name.ValueString(), err, httpResp)

		return
	}

	sg := result.SecurityGroup
	if sg == nil {
		resp.Diagnostics.AddError("API returned nil", "SecurityGroup is nil in the response")

		return
	}
	if sg.Id == nil {
		resp.Diagnostics.AddError("Create Error", "Security group ID not returned")

		return
	}

	// The create endpoint ignores resourcePermissions (it parses them from the
	// request and then never applies them), so a group created with
	// resource_permission_groups_all = false and a list of group ids would come
	// back as "all groups" and fail the apply with an inconsistent result. The
	// update endpoint does honour them, so apply them with a follow-up update.
	//
	// The follow-up carries the FULL planned body, not just the permissions:
	// the update endpoint coerces an absent `active` to false and an absent
	// `visibility` to "private", so a permissions-only update would silently
	// deactivate the group. MORPH-16355.
	if rp != nil {
		if err := applyPlanAfterCreate(ctx, client, *sg.Id, &plan); err != nil {
			errfmt.DiagError(&resp.Diagnostics, errfmt.OpCreate, "security_group", plan.Name.ValueString(), err, nil)
			cleanup.TaintResourceState(ctx, cleanup.TaintResourceStateConfig{
				ResourceType: "security_group",
				ResourceID:   *sg.Id,
				StateWriter:  &resp.State,
				Diagnostics:  &resp.Diagnostics,
			})

			return
		}
	}

	state, diags := r.getSecurityGroupAsState(ctx, *sg.Id)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		cleanup.TaintResourceState(ctx, cleanup.TaintResourceStateConfig{
			ResourceType: "security_group",
			ResourceID:   *sg.Id,
			StateWriter:  &resp.State,
			Diagnostics:  &resp.Diagnostics,
		})

		return
	}

	preservePlannedPermissions(state, &plan)

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *securityGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state SecurityGroupModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	model, diags := r.getSecurityGroupAsState(ctx, state.Id.ValueInt64())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if model == nil {
		resp.State.RemoveResource(ctx)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)

	// Fill any schema-declared default the API omitted (null in state) so an
	// imported resource does not plan a change nobody made. MORPH-16192.
	resp.Diagnostics.Append(
		schemadefaults.Apply(ctx, SecurityGroupResourceSchema(ctx), &resp.State)...,
	)
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

	body, diags := updateBodyFromPlan(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, httpResp, err := client.SecurityGroupsAPI.UpdateSecurityGroups(ctx, id).
		UpdateSecurityGroupsRequest(sdk.UpdateSecurityGroupsRequest{
			SecurityGroup: *body,
		}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpUpdate, "security_group", plan.Name.ValueString(), err, httpResp)

		return
	}

	state, diags := r.getSecurityGroupAsState(ctx, id)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	preservePlannedPermissions(state, &plan)

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
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

// getSecurityGroupAsState fetches a security group by ID and returns it as a model.
// Returns nil model if the resource is not found (404).
func (r *securityGroupResource) getSecurityGroupAsState(
	ctx context.Context,
	id int64,
) (*SecurityGroupModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&diags, err)

		return nil, diags
	}

	result, httpResp, err := client.SecurityGroupsAPI.GetSecurityGroups(ctx, id).Execute()
	if errfmt.IsNotFound(httpResp) {
		return nil, diags
	}
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&diags, errfmt.OpRead, "security_group", "", err, httpResp)

		return nil, diags
	}

	sg := result.SecurityGroup
	if sg == nil {
		diags.AddError("API returned nil", "SecurityGroup is nil in the response")

		return nil, diags
	}
	var model SecurityGroupModel
	mapResponseToModel(&model, sg)

	return &model, diags
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

	if sg.Zone != nil {
		model.CloudId = convert.Int64ToType(sg.Zone.Id)
	} else {
		model.CloudId = types.Int64Null()
	}

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

// preservePlannedPermissions carries planned permission values into the state
// read back after a create or update when the API omits them.
//
// The GET may not return resourcePermission or tenants -- a known quirk of the
// endpoint. Without this, a value the practitioner configured would land in
// state as null, and every subsequent plan would propose to set it again. Only
// null state values are filled, so a value the API does return always wins.
// Shared by Create and Update so the two paths cannot drift.
func preservePlannedPermissions(state, plan *SecurityGroupModel) {
	if state.ResourcePermissionGroupsAll.IsNull() && !plan.ResourcePermissionGroupsAll.IsNull() {
		state.ResourcePermissionGroupsAll = plan.ResourcePermissionGroupsAll
	}
	if state.ResourcePermissionGroupIds.IsNull() && !plan.ResourcePermissionGroupIds.IsNull() {
		state.ResourcePermissionGroupIds = plan.ResourcePermissionGroupIds
	}
	if state.TenantIds.IsNull() && !plan.TenantIds.IsNull() {
		state.TenantIds = plan.TenantIds
	}
}

// resourcePermissionsFromPlan builds the resourcePermissions payload from the
// plan, or returns nil when resource_permission_groups_all is not configured.
//
// The two attributes are designed to be used together: `all = false` plus a
// list of group ids is the only way to express "specific groups have access",
// and the ids are only meaningful inside the permissions object. The update
// request type is returned because it is the one the API actually honours (see
// applyPlanAfterCreate); the create body adapts it.
func resourcePermissionsFromPlan(
	ctx context.Context,
	plan *SecurityGroupModel,
) (*sdk.UpdateSecurityGroupsRequestSecurityGroupResourcePermissions, diag.Diagnostics) {
	var diags diag.Diagnostics

	if plan.ResourcePermissionGroupsAll.IsNull() || plan.ResourcePermissionGroupsAll.IsUnknown() {
		return nil, diags
	}

	rp := &sdk.UpdateSecurityGroupsRequestSecurityGroupResourcePermissions{
		All: plan.ResourcePermissionGroupsAll.ValueBoolPointer(),
	}

	if !plan.ResourcePermissionGroupIds.IsNull() && !plan.ResourcePermissionGroupIds.IsUnknown() {
		var groupIDs []int64
		diags.Append(plan.ResourcePermissionGroupIds.ElementsAs(ctx, &groupIDs, false)...)
		if diags.HasError() {
			return nil, diags
		}
		sites := make([]sdk.UpdateCloudFoldersRequestFolderResourcePermissionsSitesInner, len(groupIDs))
		for i, gid := range groupIDs {
			id := gid
			sites[i] = sdk.UpdateCloudFoldersRequestFolderResourcePermissionsSitesInner{Id: &id}
		}
		rp.Sites = sites
	}

	return rp, diags
}

// updateBodyFromPlan builds the full update request body from the plan. It is
// shared by Update and by the follow-up update Create issues.
//
// Every field is sent explicitly, because the update endpoint does not treat
// absent fields as "leave unchanged" for all attributes: it coerces an absent
// `active` to false and an absent `visibility` to "private". Sending the whole
// planned body makes the request idempotent with respect to the plan.
func updateBodyFromPlan(
	ctx context.Context,
	plan *SecurityGroupModel,
) (*sdk.UpdateSecurityGroupsRequestSecurityGroup, diag.Diagnostics) {
	var diags diag.Diagnostics

	body := &sdk.UpdateSecurityGroupsRequestSecurityGroup{}
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

	// Tenant permissions: send only when tenant_ids is explicitly configured.
	//
	// Do NOT send an empty tenantPermissions to "clear" tenants. Tenant and
	// group (resource) permissions share one server-side table keyed by
	// account, and an empty accounts list makes the API delete every row it
	// holds -- including the owner account's own row that carries
	// resource_permission_groups_all / resource_permission_group_ids. The GET
	// then stops returning resourcePermission at all. Omitting the field leaves
	// permissions untouched, which is the intended "unchanged" semantics;
	// tenant_ids is Computed, so the API's own value is carried in state.
	if !plan.TenantIds.IsNull() && !plan.TenantIds.IsUnknown() {
		var tenantIDs []int64
		diags.Append(plan.TenantIds.ElementsAs(ctx, &tenantIDs, false)...)
		if diags.HasError() {
			return nil, diags
		}
		body.TenantPermissions = &sdk.UpdateSecurityGroupsRequestSecurityGroupTenantPermissions{
			Accounts: tenantIDs,
		}
	}

	rp, rpDiags := resourcePermissionsFromPlan(ctx, plan)
	diags.Append(rpDiags...)
	if diags.HasError() {
		return nil, diags
	}
	body.ResourcePermissions = rp

	return body, diags
}

// applyPlanAfterCreate re-sends the full planned body as an update immediately
// after a create.
//
// The create endpoint accepts a resourcePermissions object but does not act on
// it -- the API parses it from the request and then never applies it, so a
// group created with `all = false` and a list of group ids is stored as "all
// groups" with no sites. The update endpoint does apply it, so the plan is
// applied a second time through that path. See updateBodyFromPlan for why the
// whole body, rather than just the permissions, is sent.
func applyPlanAfterCreate(
	ctx context.Context,
	client *sdk.APIClient,
	id int64,
	plan *SecurityGroupModel,
) error {
	body, diags := updateBodyFromPlan(ctx, plan)
	if diags.HasError() {
		return fmt.Errorf("building follow-up update after create: %s", diags.Errors()[0].Detail())
	}

	_, httpResp, err := client.SecurityGroupsAPI.UpdateSecurityGroups(ctx, id).
		UpdateSecurityGroupsRequest(sdk.UpdateSecurityGroupsRequest{
			SecurityGroup: *body,
		}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		return fmt.Errorf("applying resource permissions after create: %w", err)
	}

	return nil
}
