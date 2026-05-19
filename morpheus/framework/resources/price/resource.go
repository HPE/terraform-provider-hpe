package price

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
	_ resource.Resource                = &priceResource{}
	_ resource.ResourceWithConfigure   = &priceResource{}
	_ resource.ResourceWithImportState = &priceResource{}
)

type priceResource struct {
	configure.ResourceWithMorpheusConfigure
}

func NewResource() resource.Resource {
	return &priceResource{}
}

func (r *priceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_morpheus_price"
}

func (r *priceResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = PriceSchema(ctx)
}

func (r *priceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var plan priceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := sdk.AddPricesRequestPrice{
		Name:         plan.Name.ValueString(),
		Code:         plan.Code.ValueString(),
		PriceType:    plan.PriceType.ValueString(),
		PriceUnit:    plan.PriceUnit.ValueString(),
		Cost:         float32(plan.Cost.ValueFloat64()),
		Currency:     plan.Currency.ValueString(),
		IncurCharges: "always",
	}
	if !plan.MarkupType.IsNull() {
		body.MarkupType = plan.MarkupType.ValueStringPointer()
	}
	if !plan.Markup.IsNull() {
		v := float32(plan.Markup.ValueFloat64())
		body.Markup = &v
	}

	result, httpResp, err := client.PricesAPI.AddPrices(ctx).AddPricesRequest(sdk.AddPricesRequest{
		Price: body,
	}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpCreate, "price", plan.Name.ValueString(), err, httpResp)
		return
	}

	// The API returns {"success": true, "id": <int>} without a nested price object.
	// Extract the ID from AdditionalProperties and do a follow-up read.
	var id int64
	if p := result.GetPrice(); p.Id != nil {
		id = *p.Id
	} else if raw, ok := result.AdditionalProperties["id"]; ok {
		switch v := raw.(type) {
		case float64:
			id = int64(v)
		case int64:
			id = v
		}
	}
	if id == 0 {
		resp.Diagnostics.AddError("Error creating price", "Could not determine ID from create response")
		return
	}
	plan.ID = types.Int64Value(id)

	// Read back the full object
	readResult, readHttpResp, readErr := client.PricesAPI.GetPrices(ctx, id).Execute()
	if err := errfmt.CheckResponse(readErr, readHttpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "price", plan.Name.ValueString(), err, readHttpResp)
		return
	}
	rp := readResult.GetPrice()
	mapGetResponseToModel(&plan, &rp)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *priceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var state priceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	result, httpResp, err := client.PricesAPI.GetPrices(ctx, id).Execute()
	if errfmt.IsNotFound(httpResp) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "price", "", err, httpResp)
		return
	}

	p := result.GetPrice()
	mapGetResponseToModel(&state, &p)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *priceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var plan priceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.ID.ValueInt64()

	body := sdk.UpdatePricesRequestPrice{
		Name:      plan.Name.ValueStringPointer(),
		Code:      plan.Code.ValueStringPointer(),
		PriceType: plan.PriceType.ValueStringPointer(),
		PriceUnit: plan.PriceUnit.ValueStringPointer(),
		Currency:  plan.Currency.ValueStringPointer(),
	}
	cost := float32(plan.Cost.ValueFloat64())
	body.Cost = &cost
	if !plan.MarkupType.IsNull() {
		body.MarkupType = plan.MarkupType.ValueStringPointer()
	}
	if !plan.Markup.IsNull() {
		v := float32(plan.Markup.ValueFloat64())
		body.Markup = &v
	}

	result, httpResp, err := client.PricesAPI.UpdatePrices(ctx, id).UpdatePricesRequest(sdk.UpdatePricesRequest{
		Price: body,
	}).Execute()
	_ = result
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpUpdate, "price", plan.Name.ValueString(), err, httpResp)
		return
	}

	// Read back the full object since update only returns {"success": true}
	readResult, readHttpResp, readErr := client.PricesAPI.GetPrices(ctx, id).Execute()
	if err := errfmt.CheckResponse(readErr, readHttpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "price", plan.Name.ValueString(), err, readHttpResp)
		return
	}
	rp := readResult.GetPrice()
	mapGetResponseToModel(&plan, &rp)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *priceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var state priceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	_, httpResp, err := client.PricesAPI.DeactivatePrices(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpDelete, "price", "", err, httpResp)
		return
	}
}

func (r *priceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Could not parse ID %q as integer: %s", req.ID, err))
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func mapAddResponseToModel(model *priceModel, p *sdk.AddPrices200ResponseAllOfPrice) {
	if p.Id != nil {
		model.ID = types.Int64Value(*p.Id)
	}
	if p.Name != nil {
		model.Name = types.StringValue(*p.Name)
	}
	if p.Code != nil {
		model.Code = types.StringValue(*p.Code)
	}
	if p.PriceType != nil {
		model.PriceType = types.StringValue(*p.PriceType)
	}
	if p.PriceUnit != nil {
		model.PriceUnit = types.StringValue(*p.PriceUnit)
	}
	if v := p.Cost.Get(); v != nil {
		model.Cost = types.Float64Value(float64(*v))
	}
	if v := p.MarkupType.Get(); v != nil {
		model.MarkupType = types.StringValue(*v)
	} else {
		model.MarkupType = types.StringNull()
	}
	if v := p.Markup.Get(); v != nil {
		model.Markup = types.Float64Value(float64(*v))
	} else {
		model.Markup = types.Float64Null()
	}
	if p.Currency != nil {
		model.Currency = types.StringValue(*p.Currency)
	}
}

func mapGetResponseToModel(model *priceModel, p *sdk.GetPrices200ResponsePrice) {
	if p.Id != nil {
		model.ID = types.Int64Value(*p.Id)
	}
	if p.Name != nil {
		model.Name = types.StringValue(*p.Name)
	}
	if p.Code != nil {
		model.Code = types.StringValue(*p.Code)
	}
	if p.PriceType != nil {
		model.PriceType = types.StringValue(*p.PriceType)
	}
	if p.PriceUnit != nil {
		model.PriceUnit = types.StringValue(*p.PriceUnit)
	}
	if v := p.Cost.Get(); v != nil {
		model.Cost = types.Float64Value(float64(*v))
	}
	if v := p.MarkupType.Get(); v != nil {
		model.MarkupType = types.StringValue(*v)
	} else {
		model.MarkupType = types.StringNull()
	}
	if v := p.Markup.Get(); v != nil {
		model.Markup = types.Float64Value(float64(*v))
	} else {
		model.Markup = types.Float64Null()
	}
	if p.Currency != nil {
		model.Currency = types.StringValue(*p.Currency)
	}
}

func mapUpdateResponseToModel(model *priceModel, p *sdk.UpdatePrices200ResponseAllOfPrice) {
	if p.Id != nil {
		model.ID = types.Int64Value(*p.Id)
	}
	if p.Name != nil {
		model.Name = types.StringValue(*p.Name)
	}
	if p.Code != nil {
		model.Code = types.StringValue(*p.Code)
	}
	if p.PriceType != nil {
		model.PriceType = types.StringValue(*p.PriceType)
	}
	if p.PriceUnit != nil {
		model.PriceUnit = types.StringValue(*p.PriceUnit)
	}
	if v := p.Cost.Get(); v != nil {
		model.Cost = types.Float64Value(float64(*v))
	}
	if v := p.MarkupType.Get(); v != nil {
		model.MarkupType = types.StringValue(*v)
	} else {
		model.MarkupType = types.StringNull()
	}
	if v := p.Markup.Get(); v != nil {
		model.Markup = types.Float64Value(float64(*v))
	} else {
		model.Markup = types.Float64Null()
	}
	if p.Currency != nil {
		model.Currency = types.StringValue(*p.Currency)
	}
}
