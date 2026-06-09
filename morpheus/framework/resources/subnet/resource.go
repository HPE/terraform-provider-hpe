// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package subnet

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
	"github.com/HPE/terraform-provider-hpe/utils/cleanup"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

var (
	_ resource.Resource                = &subnetResource{}
	_ resource.ResourceWithConfigure   = &subnetResource{}
	_ resource.ResourceWithImportState = &subnetResource{}
)

type subnetResource struct {
	configure.ResourceWithMorpheusConfigure
}

func NewResource() resource.Resource {
	return &subnetResource{}
}

func (r *subnetResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_morpheus_subnet"
}

func (r *subnetResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = SubnetResourceSchema(ctx)
}

func (r *subnetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var plan SubnetModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := &sdk.CreateSubnetRequestSubnet{}
	body.Type = &sdk.CreateSubnetRequestSubnetType{
		Id: plan.TypeId.ValueInt64Pointer(),
	}
	body.NetworkId = plan.NetworkId.ValueInt64Pointer()

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		body.Description = plan.Description.ValueStringPointer()
	}
	if !plan.Cidr.IsNull() && !plan.Cidr.IsUnknown() {
		body.Cidr = plan.Cidr.ValueStringPointer()
	}
	if !plan.Active.IsNull() && !plan.Active.IsUnknown() {
		body.Active = plan.Active.ValueBoolPointer()
	}
	if !plan.DhcpServer.IsNull() && !plan.DhcpServer.IsUnknown() {
		body.DhcpServer = plan.DhcpServer.ValueBoolPointer()
	}
	if !plan.AllowStaticOverride.IsNull() && !plan.AllowStaticOverride.IsUnknown() {
		body.AllowStaticOverride = plan.AllowStaticOverride.ValueBoolPointer()
	}
	if !plan.PoolId.IsNull() && !plan.PoolId.IsUnknown() {
		body.Pool = &sdk.CreateSubnetRequestSubnetPool{
			Id: plan.PoolId.ValueInt64Pointer(),
		}
	}
	if !plan.Visibility.IsNull() && !plan.Visibility.IsUnknown() {
		body.Visibility = plan.Visibility.ValueStringPointer()
	}
	if !plan.Labels.IsNull() && !plan.Labels.IsUnknown() {
		var labels []string
		resp.Diagnostics.Append(plan.Labels.ElementsAs(ctx, &labels, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		body.Labels = labels
	}
	if !plan.TenantIds.IsNull() && !plan.TenantIds.IsUnknown() {
		var tenantIDs []int64
		resp.Diagnostics.Append(plan.TenantIds.ElementsAs(ctx, &tenantIDs, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		tenants := make([]sdk.CreateSubnetRequestSubnetTenantsInner, len(tenantIDs))
		for i, tid := range tenantIDs {
			id := tid
			tenants[i] = sdk.CreateSubnetRequestSubnetTenantsInner{Id: &id}
		}
		body.Tenants = tenants
	}
	if !plan.Config.IsNull() && !plan.Config.IsUnknown() {
		configValue := plan.Config.UnderlyingValue()
		configAny, err := convert.ValueToAny(ctx, configValue)
		if err != nil {
			resp.Diagnostics.AddError(
				"create subnet resource",
				"failed to convert config: "+err.Error(),
			)

			return
		}

		configMap, ok := configAny.(map[string]any)
		if !ok {
			resp.Diagnostics.AddError(
				"create subnet resource",
				"config must be a valid object/map",
			)

			return
		}
		body.Config = configMap
	}

	createReq := sdk.CreateSubnetRequest{
		Subnet: body,
	}

	// Resource permissions
	if !plan.ResourcePermissionGroupsAll.IsNull() && !plan.ResourcePermissionGroupsAll.IsUnknown() {
		rp := &sdk.CreateSubnetRequestResourcePermission{
			All: plan.ResourcePermissionGroupsAll.ValueBoolPointer(),
		}
		if !plan.ResourcePermissionGroupIds.IsNull() && !plan.ResourcePermissionGroupIds.IsUnknown() {
			var groupIDs []int64
			resp.Diagnostics.Append(plan.ResourcePermissionGroupIds.ElementsAs(ctx, &groupIDs, false)...)
			if resp.Diagnostics.HasError() {
				return
			}
			sites := make([]sdk.CreateSubnetRequestResourcePermissionSitesInner, len(groupIDs))
			for i, gid := range groupIDs {
				id := gid
				sites[i] = sdk.CreateSubnetRequestResourcePermissionSitesInner{Id: &id}
			}
			rp.Sites = sites
		}
		createReq.ResourcePermission = rp
	}

	result, httpResp, err := client.NetworksAPI.CreateSubnet(ctx).CreateSubnetRequest(createReq).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpCreate, "subnet", "", err, httpResp)

		return
	}

	subnet := result.Subnet
	if subnet == nil {
		resp.Diagnostics.AddError("API returned nil", "Subnet is nil in the response")

		return
	}
	if subnet.Id == nil {
		resp.Diagnostics.AddError("create subnet resource", "API returned a subnet with no ID")

		return
	}

	id := *subnet.Id
	plan.Id = convert.Int64ToType(&id)

	// Re-read via GET to populate full state
	readResult, readHTTPResp, readErr := client.NetworksAPI.GetSubnet(ctx, id).Execute()
	if readErr := errfmt.CheckResponse(readErr, readHTTPResp); readErr != nil {
		resp.Diagnostics.AddError(
			"create subnet resource",
			fmt.Sprintf("Subnet %d was created but could not be read: %s", id, readErr.Error()),
		)
		cleanup.TaintResourceState(ctx, cleanup.TaintResourceStateConfig{
			ResourceType: "subnet",
			ResourceID:   id,
			StateWriter:  &resp.State,
			Diagnostics:  &resp.Diagnostics,
		})

		return
	}

	readSubnet := readResult.Subnet
	if readSubnet == nil {
		resp.Diagnostics.AddError("API returned nil", "Subnet is nil in the response")

		return
	}

	mapResponseToModel(&plan, readSubnet)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *subnetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var state SubnetModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Id.ValueInt64()

	result, httpResp, err := client.NetworksAPI.GetSubnet(ctx, id).Execute()
	if errfmt.IsNotFound(httpResp) {
		resp.State.RemoveResource(ctx)

		return
	}
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "subnet", "", err, httpResp)

		return
	}

	subnet := result.Subnet
	if subnet == nil {
		resp.Diagnostics.AddError("API returned nil", "Subnet is nil in the response")

		return
	}
	mapResponseToModel(&state, subnet)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *subnetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var plan SubnetModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.Id.ValueInt64()

	body := &sdk.UpdateSubnetRequestSubnet{}

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		body.Description = plan.Description.ValueStringPointer()
	}
	if !plan.Cidr.IsNull() && !plan.Cidr.IsUnknown() {
		body.Cidr = plan.Cidr.ValueStringPointer()
	}
	if !plan.Active.IsNull() && !plan.Active.IsUnknown() {
		body.Active = plan.Active.ValueBoolPointer()
	}
	if !plan.DhcpServer.IsNull() && !plan.DhcpServer.IsUnknown() {
		body.DhcpServer = plan.DhcpServer.ValueBoolPointer()
	}
	if !plan.AllowStaticOverride.IsNull() && !plan.AllowStaticOverride.IsUnknown() {
		body.AllowStaticOverride = plan.AllowStaticOverride.ValueBoolPointer()
	}
	if !plan.PoolId.IsNull() && !plan.PoolId.IsUnknown() {
		body.Pool = &sdk.UpdateSubnetRequestSubnetPool{
			Id: plan.PoolId.ValueInt64Pointer(),
		}
	}
	if !plan.Visibility.IsNull() && !plan.Visibility.IsUnknown() {
		body.Visibility = plan.Visibility.ValueStringPointer()
	}
	if !plan.Labels.IsNull() && !plan.Labels.IsUnknown() {
		var labels []string
		resp.Diagnostics.Append(plan.Labels.ElementsAs(ctx, &labels, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		body.Labels = labels
	}
	// Tenants — always send to avoid perpetual diff when user removes tenant_ids from config.
	if !plan.TenantIds.IsNull() && !plan.TenantIds.IsUnknown() {
		var tenantIDs []int64
		resp.Diagnostics.Append(plan.TenantIds.ElementsAs(ctx, &tenantIDs, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		tenants := make([]sdk.UpdateSubnetRequestSubnetTenantsInner, len(tenantIDs))
		for i, tid := range tenantIDs {
			id := tid
			tenants[i] = sdk.UpdateSubnetRequestSubnetTenantsInner{Id: &id}
		}
		body.Tenants = tenants
	} else {
		body.Tenants = []sdk.UpdateSubnetRequestSubnetTenantsInner{}
	}

	updateReq := sdk.UpdateSubnetRequest{
		Subnet: body,
	}

	// Resource permissions
	if !plan.ResourcePermissionGroupsAll.IsNull() && !plan.ResourcePermissionGroupsAll.IsUnknown() {
		rp := &sdk.UpdateSubnetRequestResourcePermission{
			All: plan.ResourcePermissionGroupsAll.ValueBoolPointer(),
		}
		if !plan.ResourcePermissionGroupIds.IsNull() && !plan.ResourcePermissionGroupIds.IsUnknown() {
			var groupIDs []int64
			resp.Diagnostics.Append(plan.ResourcePermissionGroupIds.ElementsAs(ctx, &groupIDs, false)...)
			if resp.Diagnostics.HasError() {
				return
			}
			sites := make([]sdk.UpdateSubnetRequestResourcePermissionSitesInner, len(groupIDs))
			for i, gid := range groupIDs {
				id := gid
				sites[i] = sdk.UpdateSubnetRequestResourcePermissionSitesInner{Id: &id}
			}
			rp.Sites = sites
		}
		updateReq.ResourcePermission = rp
	}

	_, httpResp, err := client.NetworksAPI.UpdateSubnet(ctx, id).UpdateSubnetRequest(updateReq).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpUpdate, "subnet", "", err, httpResp)

		return
	}

	// Re-read via GET to populate full state
	readResult, readHTTPResp, readErr := client.NetworksAPI.GetSubnet(ctx, id).Execute()
	if readErr := errfmt.CheckResponse(readErr, readHTTPResp); readErr != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "subnet", "", readErr, readHTTPResp)

		return
	}

	readSubnet := readResult.Subnet
	if readSubnet == nil {
		resp.Diagnostics.AddError("API returned nil", "Subnet is nil in the response")

		return
	}

	mapResponseToModel(&plan, readSubnet)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *subnetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var state SubnetModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Id.ValueInt64()

	_, httpResp, err := client.NetworksAPI.DeleteSubnet(ctx, id).Execute()
	if errfmt.IsNotFound(httpResp) {
		return
	}
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpDelete, "subnet", "", err, httpResp)

		return
	}
}

func (r *subnetResource) ImportState(
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

func mapResponseToModel(model *SubnetModel, subnet *sdk.GetSubnet200ResponseSubnet) {
	if subnet.Id != nil {
		model.Id = convert.Int64ToType(subnet.Id)
	}
	if subnet.Name != nil {
		model.Name = convert.StrToType(subnet.Name)
	}
	if subnet.Description.IsSet() && subnet.Description.Get() != nil {
		model.Description = convert.StrToType(subnet.Description.Get())
	} else {
		model.Description = types.StringNull()
	}
	if subnet.Cidr != nil {
		model.Cidr = convert.StrToType(subnet.Cidr)
	}
	if subnet.Gateway.IsSet() && subnet.Gateway.Get() != nil {
		model.Gateway = convert.StrToType(subnet.Gateway.Get())
	} else {
		model.Gateway = types.StringNull()
	}
	if subnet.Netmask != nil {
		model.Netmask = convert.StrToType(subnet.Netmask)
	}
	if subnet.SubnetAddress != nil {
		model.SubnetAddress = convert.StrToType(subnet.SubnetAddress)
	}
	if subnet.Active != nil {
		model.Active = convert.BoolToType(subnet.Active)
	}
	if subnet.DhcpServer != nil {
		model.DhcpServer = convert.BoolToType(subnet.DhcpServer)
	}
	if subnet.Visibility != nil {
		model.Visibility = convert.StrToType(subnet.Visibility)
	}
	if subnet.Type != nil && subnet.Type.Id != nil {
		model.TypeId = convert.Int64ToType(subnet.Type.Id)
	}
	if subnet.Network != nil && subnet.Network.Id != nil {
		model.NetworkId = convert.Int64ToType(subnet.Network.Id)
	}
	if subnet.Pool != nil && subnet.Pool.Id != nil {
		model.PoolId = convert.Int64ToType(subnet.Pool.Id)
	} else {
		model.PoolId = types.Int64Null()
	}
	if subnet.Zone != nil && subnet.Zone.Id != nil {
		model.CloudId = convert.Int64ToType(subnet.Zone.Id)
	} else {
		model.CloudId = types.Int64Null()
	}
	if subnet.Labels != nil {
		labels, _ := types.SetValueFrom(context.Background(), types.StringType, subnet.Labels)
		model.Labels = labels
	}

	// Tenants — intentionally not updated from the API response.
	// The API always includes the master tenant in the returned list even if it
	// was not explicitly set in the configuration, which would cause a perpetual
	// diff. We preserve the plan/state value instead.

	// Resource permissions
	if subnet.ResourcePermission != nil {
		model.ResourcePermissionGroupsAll = convert.BoolToType(subnet.ResourcePermission.All)
		model.ResourcePermissionGroupIds = extractGroupIDsFromSites(subnet.ResourcePermission.Sites)
	} else {
		model.ResourcePermissionGroupsAll = types.BoolNull()
		model.ResourcePermissionGroupIds = types.SetNull(types.Int64Type)
	}
}

// extractGroupIDsFromSites converts the untyped sites slice ([]map[string]interface{})
// into a types.Set of Int64 group IDs.
func extractGroupIDsFromSites(sites []map[string]interface{}) types.Set {
	if len(sites) == 0 {
		return types.SetNull(types.Int64Type)
	}

	groupValues := make([]attr.Value, 0, len(sites))
	for _, site := range sites {
		if idVal, ok := site["id"]; ok && idVal != nil {
			switch v := idVal.(type) {
			case float64:
				groupValues = append(groupValues, types.Int64Value(int64(v)))
			case int64:
				groupValues = append(groupValues, types.Int64Value(v))
			}
		}
	}

	if len(groupValues) == 0 {
		return types.SetNull(types.Int64Type)
	}

	result, _ := types.SetValue(types.Int64Type, groupValues)

	return result
}
