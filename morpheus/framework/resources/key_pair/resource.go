package key_pair

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
	_ resource.Resource                = &keyPairResource{}
	_ resource.ResourceWithConfigure   = &keyPairResource{}
	_ resource.ResourceWithImportState = &keyPairResource{}
)

type keyPairResource struct {
	configure.ResourceWithMorpheusConfigure
}

func NewResource() resource.Resource {
	return &keyPairResource{}
}

func (r *keyPairResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_morpheus_key_pair"
}

func (r *keyPairResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = KeyPairSchema(ctx)
}

func (r *keyPairResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var plan keyPairModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := sdk.AddKeyPairsRequestKeyPair{
		Name:      plan.Name.ValueString(),
		PublicKey: plan.PublicKey.ValueString(),
	}
	if !plan.PrivateKey.IsNull() {
		body.PrivateKey = plan.PrivateKey.ValueStringPointer()
	}
	if !plan.Passphrase.IsNull() {
		body.Passphrase = plan.Passphrase.ValueStringPointer()
	}

	result, httpResp, err := client.KeyPairsAPI.AddKeyPairs(ctx).AddKeyPairsRequest(sdk.AddKeyPairsRequest{
		KeyPair: body,
	}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpCreate, "key_pair", plan.Name.ValueString(), err, httpResp)
		return
	}

	// SDK maps response field as "account" but API returns "keyPair" — extract from AdditionalProperties
	if kpData, ok := result.AdditionalProperties["keyPair"]; ok {
		if kpMap, ok := kpData.(map[string]interface{}); ok {
			mapGenericResponseToModel(&plan, kpMap)
		}
	} else {
		// Fallback: try the SDK accessor
		kp := result.GetAccount()
		mapAddResponseToModel(&plan, &kp)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *keyPairResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var state keyPairModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	result, httpResp, err := client.KeyPairsAPI.GetKeyPairs(ctx, id).Execute()
	if errfmt.IsNotFound(httpResp) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "key_pair", "", err, httpResp)
		return
	}

	// SDK maps response field as "account" but API returns "keyPair" — extract from AdditionalProperties
	if kpData, ok := result.AdditionalProperties["keyPair"]; ok {
		if kpMap, ok := kpData.(map[string]interface{}); ok {
			mapGenericResponseToModel(&state, kpMap)
		}
	} else {
		kp := result.GetAccount()
		mapGetResponseToModel(&state, &kp)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *keyPairResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Key pairs cannot be updated in place. All mutable attributes use RequiresReplace,
	// so this method should never be called.
	resp.Diagnostics.AddError(
		"Update Not Supported",
		"Key pairs cannot be updated in place. Changes require resource replacement.",
	)
}

func (r *keyPairResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var state keyPairModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	_, httpResp, err := client.KeyPairsAPI.RemoveKeyPairs(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpDelete, "key_pair", "", err, httpResp)
		return
	}
}

func (r *keyPairResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Could not parse ID %q as integer: %s", req.ID, err))
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func mapAddResponseToModel(model *keyPairModel, kp *sdk.AddKeyPairs200ResponseAllOfAccount) {
	if kp.Id != nil {
		model.ID = types.Int64Value(*kp.Id)
	}
	if kp.Name != nil {
		model.Name = types.StringValue(*kp.Name)
	}
	if kp.PublicKey.IsSet() && kp.PublicKey.Get() != nil {
		model.PublicKey = types.StringValue(*kp.PublicKey.Get())
	}
	if kp.HasPrivateKey != nil {
		model.HasPrivateKey = types.BoolValue(*kp.HasPrivateKey)
	}
	if kp.Fingerprint.IsSet() && kp.Fingerprint.Get() != nil {
		model.Fingerprint = types.StringValue(*kp.Fingerprint.Get())
	} else {
		model.Fingerprint = types.StringNull()
	}
}

func mapGetResponseToModel(model *keyPairModel, kp *sdk.GetKeyPairs200ResponseAccount) {
	if kp.Id != nil {
		model.ID = types.Int64Value(*kp.Id)
	}
	if kp.Name != nil {
		model.Name = types.StringValue(*kp.Name)
	}
	if kp.PublicKey.IsSet() {
		model.PublicKey = types.StringValue(*kp.PublicKey.Get())
	}
	if kp.HasPrivateKey != nil {
		model.HasPrivateKey = types.BoolValue(*kp.HasPrivateKey)
	}
	if kp.Fingerprint.IsSet() {
		model.Fingerprint = types.StringValue(*kp.Fingerprint.Get())
	} else {
		model.Fingerprint = types.StringNull()
	}
}

func mapGenericResponseToModel(model *keyPairModel, m map[string]interface{}) {
	if id, ok := m["id"].(float64); ok {
		model.ID = types.Int64Value(int64(id))
	}
	if name, ok := m["name"].(string); ok {
		model.Name = types.StringValue(name)
	}
	if pk, ok := m["publicKey"].(string); ok {
		model.PublicKey = types.StringValue(pk)
	}
	if hp, ok := m["hasPrivateKey"].(bool); ok {
		model.HasPrivateKey = types.BoolValue(hp)
	}
	if fp, ok := m["fingerprint"].(string); ok && fp != "" {
		model.Fingerprint = types.StringValue(fp)
	} else {
		model.Fingerprint = types.StringNull()
	}
}
