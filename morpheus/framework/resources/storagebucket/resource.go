// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

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
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/wodigest"
	"github.com/HPE/terraform-provider-hpe/utils/cleanup"
	"github.com/HPE/terraform-provider-hpe/utils/schemadefaults"
)

var (
	_ resource.Resource                = &storageBucketResource{}
	_ resource.ResourceWithConfigure   = &storageBucketResource{}
	_ resource.ResourceWithImportState = &storageBucketResource{}
	_ resource.ResourceWithModifyPlan  = &storageBucketResource{}
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

// providerTypeLocal is the provider type of a bucket on the appliance's own
// disk. Its config is a directory (basePath) rather than credentials, and
// Morpheus sets its bucketName to "." itself.
const providerTypeLocal = "local"

// isLocal reports whether the configured provider type is the local one. The
// value is required, so it is known at plan time.
func isLocal(providerType types.String) bool {
	return !providerType.IsNull() && !providerType.IsUnknown() && providerType.ValueString() == providerTypeLocal
}

// knownBasePath returns the configured base_path, or nil when it is unset.
func knownBasePath(basePath types.String) *string {
	if basePath.IsNull() || basePath.IsUnknown() {
		return nil
	}

	return basePath.ValueStringPointer()
}

// addBucketConfig builds the create request's config object for the provider
// type: the local variant carries basePath, every other type the S3-shaped
// credentials and endpoint.
//
// The oneOf wrapper's generated MarshalJSON returns (nil, nil) when no variant
// is set -- which encoding/json rejects with "unexpected end of JSON input". A
// variant must therefore always be selected, even when every field inside it is
// empty, so that config is always sent as an object.
func addBucketConfig(
	providerType, basePath, endpoint, accessKey, secretKey types.String,
) *sdk.AddStorageBucketsRequestStorageBucketConfig {
	if isLocal(providerType) {
		return &sdk.AddStorageBucketsRequestStorageBucketConfig{
			AddStorageBucketsRequestStorageBucketConfigOneOf3: &sdk.AddStorageBucketsRequestStorageBucketConfigOneOf3{
				BasePath: knownBasePath(basePath),
			},
		}
	}

	ak, sk, ep := bucketConfigFields(endpoint, accessKey, secretKey)

	return &sdk.AddStorageBucketsRequestStorageBucketConfig{
		AddStorageBucketsRequestStorageBucketConfigOneOf: &sdk.AddStorageBucketsRequestStorageBucketConfigOneOf{
			AccessKey: ak,
			SecretKey: sk,
			Endpoint:  ep,
		},
	}
}

// updateBucketConfig builds the update request's config object, by provider
// type as addBucketConfig does. The credentials are passed in already resolved
// (nil when they are not being re-sent).
func updateBucketConfig(
	providerType, basePath types.String,
	accessKey, secretKey, endpoint *string,
) *sdk.UpdateStorageBucketsRequestStorageBucketConfig {
	if isLocal(providerType) {
		return &sdk.UpdateStorageBucketsRequestStorageBucketConfig{
			UpdateStorageBucketsRequestStorageBucketConfigOneOf3: &sdk.UpdateStorageBucketsRequestStorageBucketConfigOneOf3{
				BasePath: knownBasePath(basePath),
			},
		}
	}

	return &sdk.UpdateStorageBucketsRequestStorageBucketConfig{
		UpdateStorageBucketsRequestStorageBucketConfigOneOf: &sdk.UpdateStorageBucketsRequestStorageBucketConfigOneOf{
			AccessKey: accessKey,
			SecretKey: secretKey,
			Endpoint:  endpoint,
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
		Config:       addBucketConfig(plan.ProviderType, plan.BasePath, plan.Endpoint, config.AccessKey, config.SecretKey),
	}
	// bucket_name is computed when unset: Morpheus names a local bucket "." and
	// an nfs one "/" itself.
	if !plan.BucketName.IsNull() && !plan.BucketName.IsUnknown() {
		body.BucketName = plan.BucketName.ValueStringPointer()
	}
	if !plan.DefaultBackupTarget.IsNull() {
		body.DefaultBackupTarget = plan.DefaultBackupTarget.ValueBoolPointer()
	}
	// default_virtual_image_target has no default: unset, it is unknown on
	// create and carries the state value on update, so a bucket that is the
	// default stays so. Unknown is skipped explicitly -- an unknown bool's
	// pointer is to false, which would take the default away.
	if !plan.DefaultVirtualImageTarget.IsNull() && !plan.DefaultVirtualImageTarget.IsUnknown() {
		body.DefaultVirtualImageTarget = plan.DefaultVirtualImageTarget.ValueBoolPointer()
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

	plan.WriteOnlyRevision = wodigest.NewRevision()
	resp.Diagnostics.Append(wo.Store(ctx, resp.Private, req.Config)...)

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

	// Fill any schema-declared default the API omitted (null in state) so an
	// imported resource does not plan a change nobody made. MORPH-16192.
	resp.Diagnostics.Append(
		schemadefaults.Apply(ctx, StorageBucketResourceSchema(ctx), &resp.State)...,
	)
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

	// access_key and secret_key are write-only. They are re-sent when ModifyPlan
	// marked write_only_revision unknown (a digest changed, or a deprecated
	// marker did); otherwise the update carries only the endpoint. See
	// write_only.go.
	resend := wodigest.Resend(ctx, req.Plan)
	accessKey, secretKey, endpoint := bucketConfigFields(plan.Endpoint, config.AccessKey, config.SecretKey)
	if !resend {
		accessKey, secretKey = nil, nil
	}

	// Attributes removed from the configuration since the last apply, to be
	// cleared. Read before Record, below, overwrites what it reads.
	pending, pdiags := clearable.Compute(ctx, req.Config, req.State, req.Private, clearableAttributes)
	resp.Diagnostics.Append(pdiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := sdk.UpdateStorageBucketsRequestStorageBucket{
		Name:         plan.Name.ValueStringPointer(),
		ProviderType: plan.ProviderType.ValueStringPointer(),
		Config:       updateBucketConfig(plan.ProviderType, plan.BasePath, accessKey, secretKey, endpoint),
	}
	if !plan.BucketName.IsNull() && !plan.BucketName.IsUnknown() {
		body.BucketName = plan.BucketName.ValueStringPointer()
	}
	if !plan.DefaultBackupTarget.IsNull() {
		body.DefaultBackupTarget = plan.DefaultBackupTarget.ValueBoolPointer()
	}
	// default_virtual_image_target has no default: unset, it is unknown on
	// create and carries the state value on update, so a bucket that is the
	// default stays so. Unknown is skipped explicitly -- an unknown bool's
	// pointer is to false, which would take the default away.
	if !plan.DefaultVirtualImageTarget.IsNull() && !plan.DefaultVirtualImageTarget.IsUnknown() {
		body.DefaultVirtualImageTarget = plan.DefaultVirtualImageTarget.ValueBoolPointer()
	}
	sent := clearable.NullableInt64(pending, "retention_days", plan.RetentionDays, &body.RetentionPolicyDays)
	if sent && !plan.RetentionDays.IsNull() {
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

	var priorState StorageBucketModel
	resp.Diagnostics.Append(req.State.Get(ctx, &priorState)...)
	plan.WriteOnlyRevision = priorState.WriteOnlyRevision
	if resend {
		plan.WriteOnlyRevision = wodigest.NewRevision()
	}
	// Always record the digests of the configured values, sent or not.
	resp.Diagnostics.Append(wo.Store(ctx, resp.Private, req.Config)...)

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

	result, httpResp, err := client.StorageAPI.RemoveStorageBuckets(ctx, id).Execute()
	if err := errfmt.CheckDeleted(result, err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpDelete, "storage_bucket", "", err, httpResp)

		return
	}
}

func (r *storageBucketResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	// "<id>" or "<id>#access_key=<sha256 hex>&secret_key=<sha256 hex>": see wodigest.
	rawID, ok := wo.Import(ctx, req, resp)
	if !ok {
		return
	}

	id, err := strconv.ParseInt(rawID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Could not parse ID %q as integer: %s", rawID, err))

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
	if sb.DefaultVirtualImageTarget != nil {
		model.DefaultVirtualImageTarget = types.BoolValue(*sb.DefaultVirtualImageTarget)
	}
	model.BasePath = basePathFromResponse(sb)
}

// basePathFromResponse reads a local bucket's directory back from the
// response's config. The response's config model is typed for the S3 shape,
// so basePath arrives among its additional properties; any other provider type
// has none, and the attribute is null.
func basePathFromResponse(sb *sdk.GetStorageBuckets200ResponseStorageBucket) types.String {
	if sb.ProviderType == nil || *sb.ProviderType != providerTypeLocal || sb.Config == nil {
		return types.StringNull()
	}

	if v, ok := sb.Config.AdditionalProperties["basePath"].(string); ok && v != "" {
		return types.StringValue(v)
	}

	return types.StringNull()
}
