// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

// Package ostypeimage is the os_type_image resource
package ostypeimage

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

var (
	_ resource.Resource                = &Resource{}
	_ resource.ResourceWithImportState = &Resource{}
)

func NewResource() resource.Resource {
	return &Resource{}
}

type Resource struct {
	configure.ResourceWithMorpheusConfigure
	resource.Resource
}

func (r *Resource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_morpheus_os_type_image"
}

func (r *Resource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = OsTypeImageResourceSchema(ctx)
}

// getOsTypeImageAsState reads the remote object and maps it to the Terraform model.
// It fetches the os_type_image record, then looks up the associated virtual image
// to extract the OsTypeId from virtualImage.osType.id.
func getOsTypeImageAsState(
	ctx context.Context,
	id int64,
	client *sdk.APIClient,
) (OsTypeImageModel, diag.Diagnostics) {
	var state OsTypeImageModel
	var diags diag.Diagnostics

	resp, hresp, err := client.LibraryAPI.GetOsTypeImage(ctx, id).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		diags.AddError("read ostypeimage",
			fmt.Sprintf("GET osTypeImage %d failed: %s", id, errfmt.ErrMsg(err, hresp)))

		return state, diags
	}

	img := resp.OsTypeImage
	if img == nil {
		diags.AddError("read ostypeimage", fmt.Sprintf("GET osTypeImage %d returned no osTypeImage", id))

		return state, diags
	}

	state.Id = convert.Int64ToType(img.Id)
	state.VirtualImageId = convert.Int64ToType(img.VirtualImageId)

	if zone := img.Zone.Get(); zone != nil {
		state.CloudId = types.Int64Value(*zone)
	}

	if provisionType := img.ProvisionType.Get(); provisionType != nil {
		state.ProvisionTypeId = types.Int64Value(*provisionType)
	}

	if img.VirtualImageId == nil {
		return state, diags
	}

	// Resolve OsTypeId by fetching the virtual image and reading its osType.id.
	viResp, viHresp, viErr := client.LibraryAPI.GetVirtualImage(ctx, *img.VirtualImageId).Execute()
	if viErr != nil || viHresp.StatusCode != http.StatusOK {
		diags.AddError("read ostypeimage",
			fmt.Sprintf("GET virtualImage %d (for osTypeImage %d) failed: %s",
				*img.VirtualImageId, id, errfmt.ErrMsg(viErr, viHresp)))

		return state, diags
	}

	vi := viResp.VirtualImage
	if vi != nil && vi.OsType != nil && vi.OsType.Id != nil {
		state.OsTypeId = types.Int64Value(*vi.OsType.Id)
	}

	return state, diags
}

func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan OsTypeImageModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("create ostypeimage", "failed to create client: "+err.Error())

		return
	}

	osTypeImage := &sdk.AddOsTypeImageRequestOsTypeImage{
		OsType:       plan.OsTypeId.ValueInt64(),
		VirtualImage: plan.VirtualImageId.ValueInt64(),
	}

	if !plan.CloudId.IsNull() && !plan.CloudId.IsUnknown() {
		osTypeImage.Zone.Set(plan.CloudId.ValueInt64Pointer())
	}
	if !plan.ProvisionTypeId.IsNull() && !plan.ProvisionTypeId.IsUnknown() {
		osTypeImage.ProvisionType.Set(plan.ProvisionTypeId.ValueInt64Pointer())
	}

	addReq := &sdk.AddOsTypeImageRequest{OsTypeImage: osTypeImage}

	// AddOsTypeImage returns (*http.Response, error) with no parsed body.
	hresp, err := client.LibraryAPI.AddOsTypeImage(ctx).
		AddOsTypeImageRequest(*addReq).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError("create ostypeimage",
			"POST failed: "+errfmt.ErrMsg(err, hresp))

		return
	}

	defer hresp.Body.Close()
	body, err := io.ReadAll(hresp.Body)
	if err != nil {
		resp.Diagnostics.AddError("create ostypeimage", "failed to read response body: "+err.Error())

		return
	}

	var createResp struct {
		ID      int64 `json:"id"`
		Success bool  `json:"success"`
	}
	if err := json.Unmarshal(body, &createResp); err != nil {
		resp.Diagnostics.AddError("create ostypeimage", "failed to parse response: "+err.Error())

		return
	}

	if createResp.ID == 0 {
		resp.Diagnostics.AddError("create ostypeimage",
			"POST failed: "+errfmt.ErrMsg(err, hresp))

		return
	}

	plan.Id = types.Int64Value(createResp.ID)

	// Save state with ID immediately in case the subsequent read fails.
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, diags := getOsTypeImageAsState(ctx, createResp.ID, client)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var current OsTypeImageModel
	resp.Diagnostics.Append(req.State.Get(ctx, &current)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("read ostypeimage", "failed to create client: "+err.Error())

		return
	}

	id := current.Id.ValueInt64()
	state, diags := getOsTypeImageAsState(ctx, id, client)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// All attributes use RequiresReplace — Terraform will never call Update.
	resp.Diagnostics.AddError("update ostypeimage",
		"update is not supported; all changes require replacement")
}

func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data OsTypeImageModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("delete ostypeimage", "failed to create client: "+err.Error())

		return
	}

	id := data.Id.ValueInt64()
	_, hresp, err := client.LibraryAPI.DeleteOsTypeImage(ctx, id).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError("delete ostypeimage",
			fmt.Sprintf("DELETE %d failed: %s", id, errfmt.ErrMsg(err, hresp)))
	}
}

func (r *Resource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	id, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("import ostypeimage",
			"provided import ID '"+req.ID+"' is invalid (non-number)")

		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}
