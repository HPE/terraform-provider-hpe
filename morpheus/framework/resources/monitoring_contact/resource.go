package monitoring_contact

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
	_ resource.Resource                = &monitoringContactResource{}
	_ resource.ResourceWithConfigure   = &monitoringContactResource{}
	_ resource.ResourceWithImportState = &monitoringContactResource{}
)

type monitoringContactResource struct {
	configure.ResourceWithMorpheusConfigure
}

func NewResource() resource.Resource {
	return &monitoringContactResource{}
}

func (r *monitoringContactResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_morpheus_monitoring_contact"
}

func (r *monitoringContactResource) Schema(
	ctx context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = MonitoringContactSchema(ctx)
}

func (r *monitoringContactResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var plan monitoringContactModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := sdk.AddContactsRequestContact{
		Name: plan.Name.ValueString(),
	}
	if !plan.EmailAddress.IsNull() {
		body.EmailAddress = plan.EmailAddress.ValueStringPointer()
	}
	if !plan.SmsAddress.IsNull() {
		body.SmsAddress = plan.SmsAddress.ValueStringPointer()
	}
	if !plan.SlackHook.IsNull() {
		body.SlackHook = plan.SlackHook.ValueStringPointer()
	}

	result, httpResp, err := client.ContactsAPI.AddContacts(ctx).AddContactsRequest(sdk.AddContactsRequest{
		Contact: body,
	}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpCreate, "monitoring_contact", plan.Name.ValueString(), err, httpResp)

		return
	}

	contact := result.GetContact()
	mapAddContactResponseToModel(&plan, &contact)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *monitoringContactResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var state monitoringContactModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	result, httpResp, err := client.ContactsAPI.GetContacts(ctx, id).Execute()
	if errfmt.IsNotFound(httpResp) {
		resp.State.RemoveResource(ctx)

		return
	}
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "monitoring_contact", "", err, httpResp)

		return
	}

	contact := result.GetContact()
	mapGetContactResponseToModel(&state, &contact)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *monitoringContactResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var plan monitoringContactModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.ID.ValueInt64()

	body := sdk.UpdateContactsRequestContact{
		Name: plan.Name.ValueStringPointer(),
	}
	if !plan.EmailAddress.IsNull() {
		body.EmailAddress = plan.EmailAddress.ValueStringPointer()
	}
	if !plan.SmsAddress.IsNull() {
		body.SmsAddress = plan.SmsAddress.ValueStringPointer()
	}
	if !plan.SlackHook.IsNull() {
		body.SlackHook = plan.SlackHook.ValueStringPointer()
	}

	result, httpResp, err := client.ContactsAPI.UpdateContacts(ctx, id).UpdateContactsRequest(sdk.UpdateContactsRequest{
		Contact: body,
	}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpUpdate, "monitoring_contact", plan.Name.ValueString(), err, httpResp)

		return
	}

	contact := result.GetContact()
	mapUpdateContactResponseToModel(&plan, &contact)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *monitoringContactResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var state monitoringContactModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueInt64()

	_, httpResp, err := client.ContactsAPI.DeleteContacts(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpDelete, "monitoring_contact", "", err, httpResp)

		return
	}
}

func (r *monitoringContactResource) ImportState(
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

func mapAddContactResponseToModel(model *monitoringContactModel, contact *sdk.AddContacts200ResponseAllOfContact) {
	if contact.Id != nil {
		model.ID = types.Int64Value(*contact.Id)
	}
	if contact.Name != nil {
		model.Name = types.StringValue(*contact.Name)
	}
	if contact.EmailAddress != nil {
		model.EmailAddress = types.StringValue(*contact.EmailAddress)
	} else {
		model.EmailAddress = types.StringNull()
	}
	if contact.SmsAddress != nil {
		model.SmsAddress = types.StringValue(*contact.SmsAddress)
	} else {
		model.SmsAddress = types.StringNull()
	}
	if contact.SlackHook.IsSet() && contact.SlackHook.Get() != nil {
		model.SlackHook = types.StringValue(*contact.SlackHook.Get())
	} else {
		model.SlackHook = types.StringNull()
	}
}

func mapGetContactResponseToModel(model *monitoringContactModel, contact *sdk.GetContacts200ResponseContact) {
	if contact.Id != nil {
		model.ID = types.Int64Value(*contact.Id)
	}
	if contact.Name != nil {
		model.Name = types.StringValue(*contact.Name)
	}
	if contact.EmailAddress != nil {
		model.EmailAddress = types.StringValue(*contact.EmailAddress)
	} else {
		model.EmailAddress = types.StringNull()
	}
	if contact.SmsAddress != nil {
		model.SmsAddress = types.StringValue(*contact.SmsAddress)
	} else {
		model.SmsAddress = types.StringNull()
	}
	if contact.SlackHook.IsSet() && contact.SlackHook.Get() != nil {
		model.SlackHook = types.StringValue(*contact.SlackHook.Get())
	} else {
		model.SlackHook = types.StringNull()
	}
}

func mapUpdateContactResponseToModel(
	model *monitoringContactModel,
	contact *sdk.UpdateContacts200ResponseAllOfContact,
) {
	if contact.Id != nil {
		model.ID = types.Int64Value(*contact.Id)
	}
	if contact.Name != nil {
		model.Name = types.StringValue(*contact.Name)
	}
	if contact.EmailAddress != nil {
		model.EmailAddress = types.StringValue(*contact.EmailAddress)
	} else {
		model.EmailAddress = types.StringNull()
	}
	if contact.SmsAddress != nil {
		model.SmsAddress = types.StringValue(*contact.SmsAddress)
	} else {
		model.SmsAddress = types.StringNull()
	}
	if contact.SlackHook.IsSet() && contact.SlackHook.Get() != nil {
		model.SlackHook = types.StringValue(*contact.SlackHook.Get())
	} else {
		model.SlackHook = types.StringNull()
	}
}
