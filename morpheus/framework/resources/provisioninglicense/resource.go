// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package provisioninglicense

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
	resp.TypeName = req.ProviderTypeName + "_" + "provisioning_license"
}

func (r *provisioningLicenseResource) Schema(
	ctx context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = ProvisioningLicenseResourceSchema(ctx)
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

	var plan ProvisioningLicenseModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var config ProvisioningLicenseModel
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
	if !plan.Tenants.IsNull() && !plan.Tenants.IsUnknown() {
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

	if result.License == nil || result.License.Id == nil {
		resp.Diagnostics.AddError("API returned nil ID", "License ID is nil in the create response")

		return
	}

	id := *result.License.Id

	readResult, httpResp, err := client.ProvisioningLicensesAPI.GetProvisioningLicense(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "provisioning_license", plan.Name.ValueString(), err, httpResp)
		cleanup.TaintResourceState(ctx, cleanup.TaintResourceStateConfig{
			ResourceType: "provisioning_license",
			ResourceID:   id,
			StateWriter:  &resp.State,
			Diagnostics:  &resp.Diagnostics,
		})

		return
	}

	readLicense := readResult.License
	if readLicense == nil {
		resp.Diagnostics.AddError("API returned nil", "License is nil in the response")

		return
	}

	// The API may normalise the tenants list on GET (e.g. replacing submitted IDs
	// with the master tenant). Preserve the plan value so Terraform's post-apply
	// consistency check passes. Read() will surface any real divergence on the
	// next plan.
	savedTenants := plan.Tenants
	resp.Diagnostics.Append(mapGetResponseToModel(&plan, readLicense)...)
	plan.Tenants = savedTenants
	plan.LicenseKeyWoVersion = config.LicenseKeyWoVersion

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *provisioningLicenseResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var state ProvisioningLicenseModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Detect import: ImportState sets only id; name is null.
	// On normal refresh, name is always a known string from prior state.
	isImport := state.Name.IsNull()
	priorTenants := state.Tenants

	id := state.Id.ValueInt64()

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
	resp.Diagnostics.Append(mapGetResponseToModel(&state, license)...)

	// On normal refresh, preserve tenants from prior state. The API may
	// silently drop IDs that don't exist in the environment, which would
	// cause a spurious diff. On import there is no prior state, so we use
	// the API values that mapGetResponseToModel just populated.
	if !isImport {
		state.Tenants = priorTenants
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

	var plan ProvisioningLicenseModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var config ProvisioningLicenseModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state ProvisioningLicenseModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.Id.ValueInt64()

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
	if !plan.Tenants.IsNull() && !plan.Tenants.IsUnknown() {
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

	readResult, httpResp, err := client.ProvisioningLicensesAPI.GetProvisioningLicense(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "provisioning_license", plan.Name.ValueString(), err, httpResp)

		return
	}

	readLicense := readResult.License
	if readLicense == nil {
		resp.Diagnostics.AddError("API returned nil", "License is nil in the response")

		return
	}

	// Same as Create: preserve plan value for tenants.
	savedTenants := plan.Tenants
	resp.Diagnostics.Append(mapGetResponseToModel(&plan, readLicense)...)
	plan.Tenants = savedTenants
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

	var state ProvisioningLicenseModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Id.ValueInt64()

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

func mapGetResponseToModel(
	model *ProvisioningLicenseModel,
	license *sdk.GetProvisioningLicense200ResponseLicense,
) diag.Diagnostics {
	var diags diag.Diagnostics

	if license.Id != nil {
		model.Id = types.Int64Value(*license.Id)
	}
	if license.Name != nil {
		model.Name = types.StringValue(*license.Name)
	}
	if license.Description != nil {
		model.Description = types.StringValue(*license.Description)
	} else {
		model.Description = types.StringNull()
	}
	if license.LicenseType != nil && license.LicenseType.Code != nil {
		model.LicenseType = types.StringValue(*license.LicenseType.Code)
	}
	tenantVals := make([]attr.Value, 0, len(license.Tenants))
	for _, t := range license.Tenants {
		if t.Id != nil {
			tenantVals = append(tenantVals, types.Int64Value(*t.Id))
		}
	}
	list, listDiags := types.ListValue(types.Int64Type, tenantVals)
	diags.Append(listDiags...)
	model.Tenants = list

	return diags
}
