package catalog_item_type

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
	_ resource.Resource                = &catalogItemTypeResource{}
	_ resource.ResourceWithConfigure   = &catalogItemTypeResource{}
	_ resource.ResourceWithImportState = &catalogItemTypeResource{}
)

type catalogItemTypeResource struct {
	configure.ResourceWithMorpheusConfigure
}

func NewResource() resource.Resource {
	return &catalogItemTypeResource{}
}

func (r *catalogItemTypeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_morpheus_catalog_item_type"
}

func (r *catalogItemTypeResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = CatalogItemTypeSchema(ctx)
}

func (r *catalogItemTypeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var plan catalogItemTypeModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	itemType := &sdk.InstanceCatalogItemType{
		Name: plan.Name.ValueStringPointer(),
	}
	if !plan.Description.IsNull() {
		itemType.Description = plan.Description.ValueStringPointer()
	}
	if !plan.Type.IsNull() {
		itemType.Type = plan.Type.ValueStringPointer()
	} else {
		defaultType := "instance"
		itemType.Type = &defaultType
	}
	if !plan.Enabled.IsNull() {
		itemType.Enabled = plan.Enabled.ValueBoolPointer()
	}
	if !plan.Visibility.IsNull() {
		itemType.Visibility = plan.Visibility.ValueStringPointer()
	}
	if !plan.Featured.IsNull() {
		itemType.Featured = plan.Featured.ValueBoolPointer()
	}

	// Config is required by the SDK; provide an empty JSON config string
	emptyConfig := "{}"
	itemType.Config = sdk.StringAsInstanceCatalogItemTypeConfig(&emptyConfig)

	body := sdk.AddCatalogItemTypeRequest{
		CatalogItemType: &sdk.AddCatalogItemTypeRequestCatalogItemType{
			InstanceCatalogItemType: itemType,
		},
	}

	result, httpResp, err := client.CatalogItemsAPI.AddCatalogItemType(ctx).AddCatalogItemTypeRequest(body).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpCreate, "catalog_item_type", plan.Name.ValueString(), err, httpResp)
		return
	}

	if cat := result.GetCatalogItemType(); cat.Id != nil {
		plan.ID = types.Int64Value(*cat.Id)
	}

	// Re-read to get full state
	readResult, httpResp, err := client.CatalogItemsAPI.GetCatalogItemType(ctx, plan.ID.ValueInt64()).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "catalog_item_type", plan.Name.ValueString(), err, httpResp)
		return
	}

	mapGetResponseToModel(&plan, readResult)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *catalogItemTypeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var state catalogItemTypeModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	result, httpResp, err := client.CatalogItemsAPI.GetCatalogItemType(ctx, id).Execute()
	if errfmt.IsNotFound(httpResp) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "catalog_item_type", "", err, httpResp)
		return
	}

	mapGetResponseToModel(&state, result)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *catalogItemTypeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var plan catalogItemTypeModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.ID.ValueInt64()

	itemType := &sdk.UpdateCatalogItemTypeRequestCatalogItemTypeAnyOf{
		Name: plan.Name.ValueStringPointer(),
	}
	if !plan.Description.IsNull() {
		itemType.Description = plan.Description.ValueStringPointer()
	}
	if !plan.Enabled.IsNull() {
		itemType.Enabled = plan.Enabled.ValueBoolPointer()
	}
	if !plan.Visibility.IsNull() {
		itemType.Visibility = plan.Visibility.ValueStringPointer()
	}
	if !plan.Featured.IsNull() {
		itemType.Featured = plan.Featured.ValueBoolPointer()
	}

	body := sdk.UpdateCatalogItemTypeRequest{
		CatalogItemType: &sdk.UpdateCatalogItemTypeRequestCatalogItemType{
			UpdateCatalogItemTypeRequestCatalogItemTypeAnyOf: itemType,
		},
	}

	_, httpResp, err := client.CatalogItemsAPI.UpdateCatalogItemType(ctx, id).UpdateCatalogItemTypeRequest(body).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpUpdate, "catalog_item_type", plan.Name.ValueString(), err, httpResp)
		return
	}

	// Re-read to get current state
	result, httpResp, err := client.CatalogItemsAPI.GetCatalogItemType(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "catalog_item_type", plan.Name.ValueString(), err, httpResp)
		return
	}

	mapGetResponseToModel(&plan, result)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *catalogItemTypeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var state catalogItemTypeModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	_, httpResp, err := client.CatalogItemsAPI.RemoveCatalogItemType(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpDelete, "catalog_item_type", "", err, httpResp)
		return
	}
}

func (r *catalogItemTypeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Could not parse ID %q as integer: %s", req.ID, err))
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func mapGetResponseToModel(model *catalogItemTypeModel, result *sdk.GetCatalogItemType200Response) {
	cat := result.GetCatalogItemType()
	if cat.Id != nil {
		model.ID = types.Int64Value(*cat.Id)
	}
	if cat.Name != nil {
		model.Name = types.StringValue(*cat.Name)
	}
	if cat.Description.IsSet() && cat.Description.Get() != nil {
		model.Description = types.StringValue(*cat.Description.Get())
	} else {
		model.Description = types.StringNull()
	}
	if cat.Type != nil {
		model.Type = types.StringValue(*cat.Type)
	}
	if cat.Enabled != nil {
		model.Enabled = types.BoolValue(*cat.Enabled)
	}
	if cat.Visibility != nil {
		model.Visibility = types.StringValue(*cat.Visibility)
	}
	if cat.Featured != nil {
		model.Featured = types.BoolValue(*cat.Featured)
	}
}
