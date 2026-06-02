package storage_bucket

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
	_ resource.Resource                = &storageBucketResource{}
	_ resource.ResourceWithConfigure   = &storageBucketResource{}
	_ resource.ResourceWithImportState = &storageBucketResource{}
)

type storageBucketResource struct {
	configure.ResourceWithMorpheusConfigure
}

func NewResource() resource.Resource {
	return &storageBucketResource{}
}

func (r *storageBucketResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_morpheus_storage_bucket"
}

func (r *storageBucketResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = StorageBucketSchema(ctx)
}

func (r *storageBucketResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var plan storageBucketModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := sdk.AddStorageBucketsRequestStorageBucket{
		Name:         plan.Name.ValueString(),
		ProviderType: plan.ProviderType.ValueString(),
	}
	if !plan.BucketName.IsNull() {
		body.BucketName = plan.BucketName.ValueString()
	}
	if !plan.DefaultBackupTarget.IsNull() {
		body.DefaultBackupTarget = plan.DefaultBackupTarget.ValueBoolPointer()
	}
	if !plan.RetentionDays.IsNull() {
		body.RetentionPolicyDays = plan.RetentionDays.ValueInt64Pointer()
		retType := "delete"
		body.RetentionPolicyType = &retType
	}
	if !plan.Endpoint.IsNull() || !plan.AccessKey.IsNull() || !plan.SecretKey.IsNull() {
		// These are typically passed via config; set as additional properties
		body.AdditionalProperties = map[string]interface{}{}
		if !plan.Endpoint.IsNull() {
			body.AdditionalProperties["endpoint"] = plan.Endpoint.ValueString()
		}
		if !plan.AccessKey.IsNull() {
			body.AdditionalProperties["accessKey"] = plan.AccessKey.ValueString()
		}
		if !plan.SecretKey.IsNull() {
			body.AdditionalProperties["secretKey"] = plan.SecretKey.ValueString()
		}
	}

	result, httpResp, err := client.StorageAPI.AddStorageBuckets(ctx).
		AddStorageBucketsRequest(sdk.AddStorageBucketsRequest{
			StorageBucket: body,
		}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpCreate, "storage_bucket", plan.Name.ValueString(), err, httpResp)

		return
	}

	sb := result.StorageBucket
	mapCreateResponseToModel(&plan, sb)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *storageBucketResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var state storageBucketModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	result, httpResp, err := client.StorageAPI.GetStorageBuckets(ctx, id).Execute()
	if errfmt.IsNotFound(httpResp) {
		resp.State.RemoveResource(ctx)

		return
	}
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "storage_bucket", "", err, httpResp)

		return
	}

	sb := result.StorageBucket
	mapGetResponseToModel(&state, sb)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *storageBucketResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var plan storageBucketModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.ID.ValueInt64()

	body := sdk.UpdateStorageBucketsRequestStorageBucket{
		Name:         plan.Name.ValueStringPointer(),
		ProviderType: plan.ProviderType.ValueStringPointer(),
	}
	if !plan.BucketName.IsNull() {
		body.BucketName = plan.BucketName.ValueStringPointer()
	}
	if !plan.DefaultBackupTarget.IsNull() {
		body.DefaultBackupTarget = plan.DefaultBackupTarget.ValueBoolPointer()
	}
	if !plan.RetentionDays.IsNull() {
		body.RetentionPolicyDays = plan.RetentionDays.ValueInt64Pointer()
		retType := "delete"
		body.RetentionPolicyType = &retType
	}
	if !plan.Endpoint.IsNull() || !plan.AccessKey.IsNull() || !plan.SecretKey.IsNull() {
		body.AdditionalProperties = map[string]interface{}{}
		if !plan.Endpoint.IsNull() {
			body.AdditionalProperties["endpoint"] = plan.Endpoint.ValueString()
		}
		if !plan.AccessKey.IsNull() {
			body.AdditionalProperties["accessKey"] = plan.AccessKey.ValueString()
		}
		if !plan.SecretKey.IsNull() {
			body.AdditionalProperties["secretKey"] = plan.SecretKey.ValueString()
		}
	}

	_, httpResp, err := client.StorageAPI.UpdateStorageBuckets(ctx, id).
		UpdateStorageBucketsRequest(sdk.UpdateStorageBucketsRequest{
			StorageBucket: body,
		}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpUpdate, "storage_bucket", plan.Name.ValueString(), err, httpResp)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *storageBucketResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var state storageBucketModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	_, httpResp, err := client.StorageAPI.RemoveStorageBuckets(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpDelete, "storage_bucket", "", err, httpResp)

		return
	}
}

func (r *storageBucketResource) ImportState(
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

func mapCreateResponseToModel(model *storageBucketModel, sb *sdk.AddStorageBuckets200ResponseAllOfStorageBucket) {
	if sb.Id != nil {
		model.ID = types.Int64Value(*sb.Id)
	}
	if sb.Name != nil {
		model.Name = types.StringValue(*sb.Name)
	}
	if sb.ProviderType != nil {
		model.ProviderType = types.StringValue(*sb.ProviderType)
	}
	if sb.BucketName != nil {
		model.BucketName = types.StringValue(*sb.BucketName)
	} else {
		model.BucketName = types.StringNull()
	}
	if sb.DefaultBackupTarget != nil {
		model.DefaultBackupTarget = types.BoolValue(*sb.DefaultBackupTarget)
	}
}

func mapGetResponseToModel(model *storageBucketModel, sb *sdk.GetStorageBuckets200ResponseStorageBucket) {
	if sb.Id != nil {
		model.ID = types.Int64Value(*sb.Id)
	}
	if sb.Name != nil {
		model.Name = types.StringValue(*sb.Name)
	}
	if sb.ProviderType != nil {
		model.ProviderType = types.StringValue(*sb.ProviderType)
	}
	if sb.BucketName != nil {
		model.BucketName = types.StringValue(*sb.BucketName)
	} else {
		model.BucketName = types.StringNull()
	}
	if sb.DefaultBackupTarget != nil {
		model.DefaultBackupTarget = types.BoolValue(*sb.DefaultBackupTarget)
	}
}
