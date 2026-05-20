package whitelabel_settings

import (
	"context"

	sdk "github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
)

var (
	_ resource.Resource              = &whitelabelSettingsResource{}
	_ resource.ResourceWithConfigure = &whitelabelSettingsResource{}
)

type whitelabelSettingsResource struct {
	configure.ResourceWithMorpheusConfigure
}

func NewResource() resource.Resource {
	return &whitelabelSettingsResource{}
}

func (r *whitelabelSettingsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_morpheus_whitelabel_settings"
}

func (r *whitelabelSettingsResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = WhitelabelSettingsSchema(ctx)
}

func (r *whitelabelSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Singleton resource: Create applies the settings via Update.
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var plan whitelabelSettingsModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := buildUpdateRequest(&plan)

	_, httpResp, err := client.WhitelabelSettingsAPI.UpdateWhitelabelSettings(ctx).UpdateWhitelabelSettingsRequest(body).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpCreate, "whitelabel_settings", "settings", err, httpResp)
		return
	}

	plan.ID = types.StringValue("settings")

	// Re-read to capture server state
	r.readIntoModel(ctx, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *whitelabelSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var state whitelabelSettingsModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, httpResp, err := client.WhitelabelSettingsAPI.ListWhitelabelSettings(ctx).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "whitelabel_settings", "", err, httpResp)
		return
	}

	settings := result.GetWhitelabelSettings()
	mapResponseToModel(&state, &settings)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *whitelabelSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	var plan whitelabelSettingsModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := buildUpdateRequest(&plan)

	_, httpResp, err := client.WhitelabelSettingsAPI.UpdateWhitelabelSettings(ctx).UpdateWhitelabelSettingsRequest(body).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpUpdate, "whitelabel_settings", "settings", err, httpResp)
		return
	}

	plan.ID = types.StringValue("settings")

	// Re-read to capture server state
	r.readIntoModel(ctx, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *whitelabelSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Singleton resource: Delete resets settings to defaults.
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)
		return
	}

	// Reset by disabling whitelabel and clearing fields.
	enabled := false
	body := sdk.UpdateWhitelabelSettingsRequest{
		WhitelabelSettings: &sdk.UpdateWhitelabelSettingsRequestWhitelabelSettings{
			Enabled:         &enabled,
			ResetHeaderLogo: boolPtr(true),
			ResetFooterLogo: boolPtr(true),
			ResetLoginLogo:  boolPtr(true),
			ResetFavicon:    boolPtr(true),
		},
	}

	_, httpResp, err := client.WhitelabelSettingsAPI.UpdateWhitelabelSettings(ctx).UpdateWhitelabelSettingsRequest(body).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpDelete, "whitelabel_settings", "", err, httpResp)
		return
	}
}

func (r *whitelabelSettingsResource) readIntoModel(ctx context.Context, model *whitelabelSettingsModel, diagnostics *diag.Diagnostics) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(diagnostics, err)
		return
	}

	result, httpResp, err := client.WhitelabelSettingsAPI.ListWhitelabelSettings(ctx).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(diagnostics, errfmt.OpRead, "whitelabel_settings", "", err, httpResp)
		return
	}

	settings := result.GetWhitelabelSettings()
	mapResponseToModel(model, &settings)
}

func buildUpdateRequest(plan *whitelabelSettingsModel) sdk.UpdateWhitelabelSettingsRequest {
	settings := sdk.UpdateWhitelabelSettingsRequestWhitelabelSettings{}
	if !plan.Enabled.IsNull() {
		settings.Enabled = plan.Enabled.ValueBoolPointer()
	}
	if !plan.ApplianceName.IsNull() {
		settings.ApplianceName = plan.ApplianceName.ValueStringPointer()
	}
	if !plan.PrimaryColor.IsNull() {
		settings.HeaderBgColor = plan.PrimaryColor.ValueStringPointer()
	}
	if !plan.SecondaryColor.IsNull() {
		settings.HeaderFgColor = plan.SecondaryColor.ValueStringPointer()
	}
	return sdk.UpdateWhitelabelSettingsRequest{
		WhitelabelSettings: &settings,
	}
}

func mapResponseToModel(model *whitelabelSettingsModel, settings *sdk.ListWhitelabelSettings200ResponseWhitelabelSettings) {
	model.ID = types.StringValue("settings")
	if settings.Enabled != nil {
		model.Enabled = types.BoolValue(*settings.Enabled)
	} else {
		model.Enabled = types.BoolNull()
	}
	if settings.ApplianceName != nil {
		model.ApplianceName = types.StringValue(*settings.ApplianceName)
	} else {
		model.ApplianceName = types.StringNull()
	}
	if settings.HeaderLogo != nil {
		model.HeaderLogo = types.StringValue(*settings.HeaderLogo)
	} else {
		model.HeaderLogo = types.StringNull()
	}
	if settings.FooterLogo != nil {
		model.FooterLogo = types.StringValue(*settings.FooterLogo)
	} else {
		model.FooterLogo = types.StringNull()
	}
	if settings.LoginLogo != nil {
		model.LoginLogo = types.StringValue(*settings.LoginLogo)
	} else {
		model.LoginLogo = types.StringNull()
	}
	if settings.Favicon != nil {
		model.Favicon = types.StringValue(*settings.Favicon)
	} else {
		model.Favicon = types.StringNull()
	}
	if settings.HeaderBgColor != nil {
		model.PrimaryColor = types.StringValue(*settings.HeaderBgColor)
	} else {
		model.PrimaryColor = types.StringNull()
	}
	if settings.HeaderFgColor != nil {
		model.SecondaryColor = types.StringValue(*settings.HeaderFgColor)
	} else {
		model.SecondaryColor = types.StringNull()
	}
	// SupportMenuLinks is not directly mapped from the response as it requires JSON serialization
}

func boolPtr(v bool) *bool {
	return &v
}
