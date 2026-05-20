package price_set

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
	_ resource.Resource                = &priceSetResource{}
	_ resource.ResourceWithConfigure   = &priceSetResource{}
	_ resource.ResourceWithImportState = &priceSetResource{}
)

type priceSetResource struct {
	configure.ResourceWithMorpheusConfigure
}

func NewResource() resource.Resource {
	return &priceSetResource{}
}

func (r *priceSetResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_morpheus_price_set"
}

func (r *priceSetResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = PriceSetSchema(ctx)
}

func (r *priceSetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var plan priceSetModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := sdk.AddPriceSetsRequestPriceSet{
		Name:      plan.Name.ValueString(),
		Code:      plan.Code.ValueString(),
		PriceUnit: plan.PriceUnit.ValueString(),
		Type:      plan.Type.ValueString(),
	}
	if !plan.RegionCode.IsNull() {
		body.RegionCode = plan.RegionCode.ValueStringPointer()
	}

	result, httpResp, err := client.PriceSetsAPI.AddPriceSets(ctx).AddPriceSetsRequest(sdk.AddPriceSetsRequest{
		PriceSet: body,
	}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpCreate, "price_set", plan.Name.ValueString(), err, httpResp)

		return
	}

	// The API returns {"success": true, "id": <int>} without a nested priceSet object.
	var id int64
	if b := result.GetBudget(); b.Id != nil {
		id = *b.Id
	} else if raw, ok := result.AdditionalProperties["id"]; ok {
		switch v := raw.(type) {
		case float64:
			id = int64(v)
		case int64:
			id = v
		}
	}
	if id == 0 {
		resp.Diagnostics.AddError("Error creating price_set", "Could not determine ID from create response")

		return
	}
	plan.ID = types.Int64Value(id)

	// Read back the full object
	readResult, readHTTPResp, readErr := client.PriceSetsAPI.GetPriceSets(ctx, id).Execute()
	// Tolerate decode errors when HTTP response is OK (SDK type mismatch on 'account' field)
	if readHTTPResp != nil && readHTTPResp.StatusCode >= 400 {
		if err := errfmt.CheckResponse(readErr, readHTTPResp); err != nil {
			errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "price_set", plan.Name.ValueString(), err, readHTTPResp)

			return
		}
	}
	ps := readResult.GetPriceSet()
	mapGetResponseToModel(&plan, &ps)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *priceSetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var state priceSetModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	result, httpResp, err := client.PriceSetsAPI.GetPriceSets(ctx, id).Execute()
	if errfmt.IsNotFound(httpResp) {
		resp.State.RemoveResource(ctx)

		return
	}
	// Tolerate decode errors when HTTP response is OK (SDK type mismatch on 'account' field)
	if httpResp != nil && httpResp.StatusCode >= 400 {
		if err := errfmt.CheckResponse(err, httpResp); err != nil {
			errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "price_set", "", err, httpResp)

			return
		}
	}
	if result == nil {
		resp.Diagnostics.AddError("Error reading price_set", fmt.Sprintf("nil response for id %d: %v", id, err))

		return
	}

	ps := result.GetPriceSet()
	mapGetResponseToModel(&state, &ps)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *priceSetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var plan priceSetModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.ID.ValueInt64()

	body := sdk.UpdatePriceSetsRequestPriceSet{
		Name:      plan.Name.ValueStringPointer(),
		Code:      plan.Code.ValueStringPointer(),
		PriceUnit: plan.PriceUnit.ValueStringPointer(),
		Type:      plan.Type.ValueStringPointer(),
	}
	if !plan.RegionCode.IsNull() {
		body.RegionCode = plan.RegionCode.ValueStringPointer()
	}

	result, httpResp, err := client.PriceSetsAPI.UpdatePriceSets(ctx, id).
		UpdatePriceSetsRequest(sdk.UpdatePriceSetsRequest{
			PriceSet: body,
		}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpUpdate, "price_set", plan.Name.ValueString(), err, httpResp)

		return
	}

	ps := result.GetBudget()
	mapUpdateResponseToModel(&plan, &ps)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *priceSetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var state priceSetModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	_, httpResp, err := client.PriceSetsAPI.DeactivatePriceSets(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpDelete, "price_set", "", err, httpResp)

		return
	}
}

func (r *priceSetResource) ImportState(
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

func mapGetResponseToModel(model *priceSetModel, ps *sdk.GetPriceSets200ResponsePriceSet) {
	if ps.Id != nil {
		model.ID = types.Int64Value(*ps.Id)
	}
	if ps.Name != nil {
		model.Name = types.StringValue(*ps.Name)
	}
	if ps.Code != nil {
		model.Code = types.StringValue(*ps.Code)
	}
	if ps.PriceUnit != nil {
		model.PriceUnit = types.StringValue(*ps.PriceUnit)
	}
	if ps.Type != nil {
		model.Type = types.StringValue(*ps.Type)
	}
	if ps.RegionCode != nil {
		model.RegionCode = types.StringValue(*ps.RegionCode)
	} else {
		model.RegionCode = types.StringNull()
	}
}

func mapUpdateResponseToModel(model *priceSetModel, ps *sdk.UpdatePriceSets200ResponseAllOfBudget) {
	if ps.Id != nil {
		model.ID = types.Int64Value(*ps.Id)
	}
	if ps.Name != nil {
		model.Name = types.StringValue(*ps.Name)
	}
	if ps.Code != nil {
		model.Code = types.StringValue(*ps.Code)
	}
	if ps.PriceUnit != nil {
		model.PriceUnit = types.StringValue(*ps.PriceUnit)
	}
	if ps.Type != nil {
		model.Type = types.StringValue(*ps.Type)
	}
	if ps.RegionCode != nil {
		model.RegionCode = types.StringValue(*ps.RegionCode)
	} else {
		model.RegionCode = types.StringNull()
	}
}
