package storagebucket

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/cleanup"
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
	resp.TypeName = req.ProviderTypeName + "_" + "storage_bucket"
}

func (r *storageBucketResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = StorageBucketResourceSchema(ctx)
}

// bucketConfigFields resolves the credential/endpoint values that belong under
// the request's nested config object.
//
// The API reads these only from the nested config object; values sent at the
// top level of the request body are silently discarded.
func bucketConfigFields(endpoint, accessKey, secretKey types.String) (ak, sk, ep *string) {
	if !accessKey.IsNull() && !accessKey.IsUnknown() {
		ak = accessKey.ValueStringPointer()
	}

	if !secretKey.IsNull() && !secretKey.IsUnknown() {
		sk = secretKey.ValueStringPointer()
	}

	if !endpoint.IsNull() && !endpoint.IsUnknown() {
		ep = endpoint.ValueStringPointer()
	}

	return ak, sk, ep
}

// addBucketConfig builds the create request's config object.
//
// The oneOf wrapper's generated MarshalJSON returns (nil, nil) when no variant
// is set -- which encoding/json rejects with "unexpected end of JSON input". A
// variant must therefore always be selected, even when every field inside it is
// empty, so that config is always sent as an object.
func addBucketConfig(endpoint, accessKey, secretKey types.String) *sdk.AddStorageBucketsRequestStorageBucketConfig {
	ak, sk, ep := bucketConfigFields(endpoint, accessKey, secretKey)

	return &sdk.AddStorageBucketsRequestStorageBucketConfig{
		AddStorageBucketsRequestStorageBucketConfigOneOf: &sdk.AddStorageBucketsRequestStorageBucketConfigOneOf{
			AccessKey: ak,
			SecretKey: sk,
			Endpoint:  ep,
		},
	}
}

func (r *storageBucketResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var plan StorageBucketModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// access_key and secret_key are write-only, so their values are only
	// available from config -- they are always null in the plan.
	var config StorageBucketModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := sdk.AddStorageBucketsRequestStorageBucket{
		Name:         plan.Name.ValueString(),
		ProviderType: plan.ProviderType.ValueString(),
		Config:       addBucketConfig(plan.Endpoint, config.AccessKey, config.SecretKey),
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

	result, httpResp, err := client.StorageAPI.AddStorageBuckets(ctx).
		AddStorageBucketsRequest(sdk.AddStorageBucketsRequest{
			StorageBucket: body,
		}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpCreate, "storage_bucket", plan.Name.ValueString(), err, httpResp)

		return
	}

	if result.StorageBucket == nil || result.StorageBucket.Id == nil {
		resp.Diagnostics.AddError("API returned nil ID", "StorageBucket ID is nil in the create response")

		return
	}

	id := *result.StorageBucket.Id

	readResult, httpResp, err := client.StorageAPI.GetStorageBuckets(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "storage_bucket", plan.Name.ValueString(), err, httpResp)
		cleanup.TaintResourceState(ctx, cleanup.TaintResourceStateConfig{
			ResourceType: "storage_bucket",
			ResourceID:   id,
			StateWriter:  &resp.State,
			Diagnostics:  &resp.Diagnostics,
		})

		return
	}

	if readResult.StorageBucket == nil {
		resp.Diagnostics.AddError("API returned nil", "StorageBucket is nil in the response")

		return
	}
	mapGetResponseToModel(&plan, readResult.StorageBucket)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *storageBucketResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var state StorageBucketModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Id.ValueInt64()

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
	if sb == nil {
		resp.Diagnostics.AddError("API returned nil", "StorageBucket is nil in the response")

		return
	}
	mapGetResponseToModel(&state, sb)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *storageBucketResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var plan StorageBucketModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// access_key and secret_key are write-only, so their values are only
	// available from config -- they are always null in the plan.
	var config StorageBucketModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.Id.ValueInt64()

	accessKey, secretKey, endpoint := bucketConfigFields(plan.Endpoint, config.AccessKey, config.SecretKey)

	body := sdk.UpdateStorageBucketsRequestStorageBucket{
		Name:         plan.Name.ValueStringPointer(),
		ProviderType: plan.ProviderType.ValueStringPointer(),
		Config: &sdk.UpdateStorageBucketsRequestStorageBucketConfig{
			UpdateStorageBucketsRequestStorageBucketConfigOneOf: &sdk.UpdateStorageBucketsRequestStorageBucketConfigOneOf{
				AccessKey: accessKey,
				SecretKey: secretKey,
				Endpoint:  endpoint,
			},
		},
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

	_, httpResp, err := client.StorageAPI.UpdateStorageBuckets(ctx, id).
		UpdateStorageBucketsRequest(sdk.UpdateStorageBucketsRequest{
			StorageBucket: body,
		}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpUpdate, "storage_bucket", plan.Name.ValueString(), err, httpResp)

		return
	}

	readResult, httpResp, err := client.StorageAPI.GetStorageBuckets(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "storage_bucket", plan.Name.ValueString(), err, httpResp)

		return
	}

	if readResult.StorageBucket == nil {
		resp.Diagnostics.AddError("API returned nil", "StorageBucket is nil in the response")

		return
	}
	mapGetResponseToModel(&plan, readResult.StorageBucket)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *storageBucketResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var state StorageBucketModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Id.ValueInt64()

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

func mapGetResponseToModel(model *StorageBucketModel, sb *sdk.GetStorageBuckets200ResponseStorageBucket) {
	if sb.Id != nil {
		model.Id = types.Int64Value(*sb.Id)
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
