package provisioning_license

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
	_ resource.Resource                = &provisioningLicenseResource{}
	_ resource.ResourceWithConfigure   = &provisioningLicenseResource{}
	_ resource.ResourceWithImportState = &provisioningLicenseResource{}
)

type provisioningLicenseResource struct {
	configure.ResourceWithMorpheusConfigure
}

func NewResource() resource.Resource {
	return &provisioningLicenseResource{}
}

func (r *provisioningLicenseResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_morpheus_provisioning_license"
}

func (r *provisioningLicenseResource) Schema(
	ctx context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = ProvisioningLicenseSchema(ctx)
}

func (r *provisioningLicenseResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var plan provisioningLicenseModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var config provisioningLicenseModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := sdk.AddProvisioningLicenseRequestLicense{
		Name:        plan.Name.ValueString(),
		LicenseType: plan.LicenseType.ValueString(),
		LicenseKey:  config.LicenseKeyWo.ValueString(),
	}
	if !plan.Description.IsNull() {
		body.Description = plan.Description.ValueStringPointer()
	}
	if !plan.VirtualImages.IsNull() {
		var images []int64
		resp.Diagnostics.Append(plan.VirtualImages.ElementsAs(ctx, &images, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		body.VirtualImages = images
	}
	if !plan.Tenants.IsNull() {
		var tenants []int64
		resp.Diagnostics.Append(plan.Tenants.ElementsAs(ctx, &tenants, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		body.Tenants = tenants
	}

	result, httpResp, err := client.ProvisioningLicensesAPI.AddProvisioningLicense(ctx).
		AddProvisioningLicenseRequest(sdk.AddProvisioningLicenseRequest{
			License: &body,
		}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpCreate, "provisioning_license", plan.Name.ValueString(), err, httpResp)

		return
	}

	license := result.License
	if license == nil {
		resp.Diagnostics.AddError("API returned nil", "License is nil in the response")
		return
	}
	if license.Id != nil {
		plan.ID = types.Int64Value(*license.Id)
	}
	plan.LicenseKeyWoVersion = config.LicenseKeyWoVersion

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *provisioningLicenseResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var state provisioningLicenseModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	result, httpResp, err := client.ProvisioningLicensesAPI.GetProvisioningLicense(ctx, id).Execute()
	if errfmt.IsNotFound(httpResp) {
		resp.State.RemoveResource(ctx)

		return
	}
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "provisioning_license", "", err, httpResp)

		return
	}

	license := result.License
	if license == nil {
		resp.Diagnostics.AddError("API returned nil", "License is nil in the response")
		return
	}
	if license.Id != nil {
		state.ID = types.Int64Value(*license.Id)
	}
	if license.Name != nil {
		state.Name = types.StringValue(*license.Name)
	}
	if license.Description != nil {
		state.Description = types.StringValue(*license.Description)
	} else {
		state.Description = types.StringNull()
	}
	if license.LicenseType != nil && license.LicenseType.Code != nil {
		state.LicenseType = types.StringValue(*license.LicenseType.Code)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *provisioningLicenseResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var plan provisioningLicenseModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var config provisioningLicenseModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state provisioningLicenseModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.ID.ValueInt64()

	body := sdk.UpdateProvisioningLicenseRequestLicense{
		Name: plan.Name.ValueStringPointer(),
	}
	if !plan.Description.IsNull() {
		body.Description = plan.Description.ValueStringPointer()
	}
	if !plan.VirtualImages.IsNull() {
		var images []int64
		resp.Diagnostics.Append(plan.VirtualImages.ElementsAs(ctx, &images, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		body.VirtualImages = images
	}
	if !plan.Tenants.IsNull() {
		var tenants []int64
		resp.Diagnostics.Append(plan.Tenants.ElementsAs(ctx, &tenants, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		body.Tenants = tenants
	}
	if !plan.LicenseKeyWoVersion.Equal(state.LicenseKeyWoVersion) {
		body.AdditionalProperties = map[string]interface{}{
			"licenseKey": config.LicenseKeyWo.ValueString(),
		}
	}

	_, httpResp, err := client.ProvisioningLicensesAPI.UpdateProvisioningLicense(ctx, id).
		UpdateProvisioningLicenseRequest(sdk.UpdateProvisioningLicenseRequest{
			License: &body,
		}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpUpdate, "provisioning_license", plan.Name.ValueString(), err, httpResp)

		return
	}

	plan.LicenseKeyWoVersion = config.LicenseKeyWoVersion

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *provisioningLicenseResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var state provisioningLicenseModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	_, httpResp, err := client.ProvisioningLicensesAPI.RemoveProvisioningLicense(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpDelete, "provisioning_license", "", err, httpResp)

		return
	}
}

func (r *provisioningLicenseResource) ImportState(
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
