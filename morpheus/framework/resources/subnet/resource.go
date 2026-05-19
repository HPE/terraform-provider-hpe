package subnet

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
	resp.Schema = SubnetSchema(ctx)
}

func (r *subnetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var plan subnetModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := sdk.CreateSubnetRequestSubnet{
		Type: &sdk.CreateSubnetRequestSubnetType{
			Id: plan.TypeID.ValueInt64Pointer(),
		},
	}
	if !plan.Visibility.IsNull() {
		body.Visibility = plan.Visibility.ValueStringPointer()
	}
	if !plan.Labels.IsNull() {
		var labels []string
		resp.Diagnostics.Append(plan.Labels.ElementsAs(ctx, &labels, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		body.Labels = labels
	}
	if !plan.Config.IsNull() {
		var configMap map[string]string
		resp.Diagnostics.Append(plan.Config.ElementsAs(ctx, &configMap, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		cfg := make(map[string]interface{}, len(configMap))
		for k, v := range configMap {
			cfg[k] = v
		}
		body.Config = cfg
	}

	result, httpResp, err := client.NetworksAPI.CreateSubnet(ctx).CreateSubnetRequest(sdk.CreateSubnetRequest{
		Subnet: &body,
	}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpCreate, "subnet", "", err, httpResp)
		return
	}

	subnet := result.GetSubnet()
	mapCreateResponseToModel(&plan, &subnet)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *subnetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var state subnetModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	result, httpResp, err := client.NetworksAPI.GetSubnet(ctx, id).Execute()
	if errfmt.IsNotFound(httpResp) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "subnet", "", err, httpResp)
		return
	}

	subnet := result.GetSubnet()
	mapResponseToModel(&state, &subnet)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *subnetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var plan subnetModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.ID.ValueInt64()

	body := sdk.UpdateSubnetRequest{}
	if !plan.Visibility.IsNull() {
		body.Subnet = &sdk.UpdateSubnetRequestSubnet{
			Visibility: plan.Visibility.ValueStringPointer(),
		}
	}

	result, httpResp, err := client.NetworksAPI.UpdateSubnet(ctx, id).UpdateSubnetRequest(body).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpUpdate, "subnet", "", err, httpResp)
		return
	}

	subnet := result.GetSubnet()
	mapResponseToModel(&plan, &subnet)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *subnetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var state subnetModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	_, httpResp, err := client.NetworksAPI.DeleteSubnet(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpDelete, "subnet", "", err, httpResp)
		return
	}
}

func (r *subnetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Could not parse ID %q as integer: %s", req.ID, err))
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func mapCreateResponseToModel(model *subnetModel, subnet *sdk.CreateSubnet200ResponseSubnet) {
	if subnet.Id != nil {
		model.ID = types.Int64Value(*subnet.Id)
	}
	if subnet.Name != nil {
		model.Name = types.StringValue(*subnet.Name)
	}
	if subnet.Cidr != nil {
		model.Cidr = types.StringValue(*subnet.Cidr)
	}
	if subnet.Gateway.IsSet() && subnet.Gateway.Get() != nil {
		model.Gateway = types.StringValue(*subnet.Gateway.Get())
	} else {
		model.Gateway = types.StringNull()
	}
	if subnet.Netmask != nil {
		model.Netmask = types.StringValue(*subnet.Netmask)
	}
	if subnet.SubnetAddress != nil {
		model.SubnetAddress = types.StringValue(*subnet.SubnetAddress)
	}
	if subnet.Active != nil {
		model.Active = types.BoolValue(*subnet.Active)
	}
	if subnet.DhcpServer != nil {
		model.DhcpServer = types.BoolValue(*subnet.DhcpServer)
	}
	if subnet.Visibility != nil {
		model.Visibility = types.StringValue(*subnet.Visibility)
	}
	if subnet.Labels != nil {
		labels, _ := types.ListValueFrom(context.Background(), types.StringType, subnet.Labels)
		model.Labels = labels
	}
}

func mapResponseToModel(model *subnetModel, subnet *sdk.GetSubnet200ResponseSubnet) {
	if subnet.Id != nil {
		model.ID = types.Int64Value(*subnet.Id)
	}
	if subnet.Name != nil {
		model.Name = types.StringValue(*subnet.Name)
	}
	if subnet.Cidr != nil {
		model.Cidr = types.StringValue(*subnet.Cidr)
	}
	if subnet.Gateway.IsSet() && subnet.Gateway.Get() != nil {
		model.Gateway = types.StringValue(*subnet.Gateway.Get())
	} else {
		model.Gateway = types.StringNull()
	}
	if subnet.Netmask != nil {
		model.Netmask = types.StringValue(*subnet.Netmask)
	}
	if subnet.SubnetAddress != nil {
		model.SubnetAddress = types.StringValue(*subnet.SubnetAddress)
	}
	if subnet.Active != nil {
		model.Active = types.BoolValue(*subnet.Active)
	}
	if subnet.DhcpServer != nil {
		model.DhcpServer = types.BoolValue(*subnet.DhcpServer)
	}
	if subnet.Visibility != nil {
		model.Visibility = types.StringValue(*subnet.Visibility)
	}
	if subnet.Labels != nil {
		labels, _ := types.ListValueFrom(context.Background(), types.StringType, subnet.Labels)
		model.Labels = labels
	}
}
