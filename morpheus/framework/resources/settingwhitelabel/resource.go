package settingwhitelabel

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
)

var (
	_ resource.Resource              = &settingWhitelabelResource{}
	_ resource.ResourceWithConfigure = &settingWhitelabelResource{}
)

type settingWhitelabelResource struct {
	configure.ResourceWithMorpheusConfigure
}

func NewResource() resource.Resource {
	return &settingWhitelabelResource{}
}

func (r *settingWhitelabelResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_" + "setting_whitelabel"
}

func (r *settingWhitelabelResource) Schema(
	ctx context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = SettingWhitelabelResourceSchema(ctx)
}

func (r *settingWhitelabelResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	// Singleton resource: Create applies the settings via Update.
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var plan SettingWhitelabelModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := buildUpdateRequest(&plan)

	_, httpResp, err := client.WhitelabelSettingsAPI.UpdateWhitelabelSettings(ctx).
		UpdateWhitelabelSettingsRequest(body).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpCreate, "setting_whitelabel", "settings", err, httpResp)

		return
	}

	// Upload any logo/favicon image files via the multipart images endpoint. On
	// create there is no prior state, so no resets are performed.
	r.uploadImages(ctx, client, &plan, nil, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.Id = types.StringValue("settings")

	// Re-read to capture server state
	r.readIntoModel(ctx, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *settingWhitelabelResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var state SettingWhitelabelModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, httpResp, err := client.WhitelabelSettingsAPI.ListWhitelabelSettings(ctx).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "setting_whitelabel", "", err, httpResp)

		return
	}

	settings := result.WhitelabelSettings
	if settings == nil {
		resp.Diagnostics.AddError("API returned nil", "WhitelabelSettings is nil in the response")

		return
	}
	mapResponseToModel(&state, settings)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *settingWhitelabelResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var plan, state SettingWhitelabelModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := buildUpdateRequest(&plan)

	_, httpResp, err := client.WhitelabelSettingsAPI.UpdateWhitelabelSettings(ctx).
		UpdateWhitelabelSettingsRequest(body).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpUpdate, "setting_whitelabel", "settings", err, httpResp)

		return
	}

	// Upload changed logo/favicon files and reset any that were cleared.
	r.uploadImages(ctx, client, &plan, &state, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.Id = types.StringValue("settings")

	// Re-read to capture server state
	r.readIntoModel(ctx, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *settingWhitelabelResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	// Singleton resource: Delete resets settings to defaults.
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	// Reset by disabling whitelabel and clearing all fields to zero values.
	enabled := false
	emptyStr := ""
	body := sdk.UpdateWhitelabelSettingsRequest{
		WhitelabelSettings: &sdk.UpdateWhitelabelSettingsRequestWhitelabelSettings{
			Enabled:         &enabled,
			ApplianceName:   &emptyStr,
			HeaderBgColor:   &emptyStr,
			HeaderFgColor:   &emptyStr,
			ResetHeaderLogo: boolPtr(true),
			ResetFooterLogo: boolPtr(true),
			ResetLoginLogo:  boolPtr(true),
			ResetFavicon:    boolPtr(true),
		},
	}

	_, httpResp, err := client.WhitelabelSettingsAPI.UpdateWhitelabelSettings(ctx).
		UpdateWhitelabelSettingsRequest(body).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpDelete, "setting_whitelabel", "", err, httpResp)

		return
	}
}

func (r *settingWhitelabelResource) readIntoModel(
	ctx context.Context,
	model *SettingWhitelabelModel,
	diagnostics *diag.Diagnostics,
) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(diagnostics, err)

		return
	}

	result, httpResp, err := client.WhitelabelSettingsAPI.ListWhitelabelSettings(ctx).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(diagnostics, errfmt.OpRead, "setting_whitelabel", "", err, httpResp)

		return
	}

	settings := result.WhitelabelSettings
	if settings == nil {
		diagnostics.AddError("API returned nil", "WhitelabelSettings is nil in the response")

		return
	}
	mapResponseToModel(model, settings)
}

// uploadImages uploads any logo/favicon image files referenced by the plan to
// the multipart whitelabel images endpoint, and resets any images that were
// present in the prior state but cleared in the plan. prior may be nil (on
// Create), in which case no resets are performed.
//
// The four image attributes (header_logo, footer_logo, login_logo, favicon) are
// local file paths. The API stores the uploaded bytes and, on read, returns a
// server-generated storage URL that never matches the supplied path, so these
// values are carried through from the plan and are never reconciled from the
// read response.
func (r *settingWhitelabelResource) uploadImages(
	ctx context.Context,
	client *sdk.APIClient,
	plan *SettingWhitelabelModel,
	prior *SettingWhitelabelModel,
	diagnostics *diag.Diagnostics,
) {
	imgReq := client.WhitelabelSettingsAPI.UpdateWhitelabelImages(ctx)

	var openFiles []*os.File
	defer func() {
		for _, f := range openFiles {
			_ = f.Close()
		}
	}()

	upload := false
	addFile := func(pathVal types.String, attach func(*os.File)) bool {
		if pathVal.IsNull() || pathVal.IsUnknown() {
			return true
		}
		f, err := os.Open(pathVal.ValueString())
		if err != nil {
			diagnostics.AddError("could not read whitelabel image file", err.Error())

			return false
		}
		openFiles = append(openFiles, f)
		attach(f)
		upload = true

		return true
	}

	if !addFile(plan.HeaderLogo, func(f *os.File) { imgReq = imgReq.HeaderLogoFile(f) }) {
		return
	}
	if !addFile(plan.FooterLogo, func(f *os.File) { imgReq = imgReq.FooterLogoFile(f) }) {
		return
	}
	if !addFile(plan.LoginLogo, func(f *os.File) { imgReq = imgReq.LoginLogoFile(f) }) {
		return
	}
	if !addFile(plan.Favicon, func(f *os.File) { imgReq = imgReq.FaviconFile(f) }) {
		return
	}

	if upload {
		_, httpResp, err := imgReq.Execute()
		if err := errfmt.CheckResponse(err, httpResp); err != nil {
			errfmt.DiagError(diagnostics, errfmt.OpUpdate, "setting_whitelabel", "images", err, httpResp)

			return
		}
	}

	if prior == nil {
		return
	}

	// Reset images that were set previously but cleared in the new plan. The JSON
	// settings endpoint honours the reset flags (the multipart endpoint's favicon
	// reset field is misnamed in the spec), so route resets there.
	resets := sdk.UpdateWhitelabelSettingsRequestWhitelabelSettings{}
	needReset := false
	clearImage := func(planVal, priorVal types.String, setReset func()) {
		if planVal.IsNull() && !priorVal.IsNull() {
			setReset()
			needReset = true
		}
	}
	clearImage(plan.HeaderLogo, prior.HeaderLogo, func() { resets.ResetHeaderLogo = boolPtr(true) })
	clearImage(plan.FooterLogo, prior.FooterLogo, func() { resets.ResetFooterLogo = boolPtr(true) })
	clearImage(plan.LoginLogo, prior.LoginLogo, func() { resets.ResetLoginLogo = boolPtr(true) })
	clearImage(plan.Favicon, prior.Favicon, func() { resets.ResetFavicon = boolPtr(true) })

	if !needReset {
		return
	}

	body := sdk.UpdateWhitelabelSettingsRequest{
		WhitelabelSettings: &resets,
	}
	_, httpResp, err := client.WhitelabelSettingsAPI.UpdateWhitelabelSettings(ctx).
		UpdateWhitelabelSettingsRequest(body).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(diagnostics, errfmt.OpUpdate, "setting_whitelabel", "reset-images", err, httpResp)

		return
	}
}

func buildUpdateRequest(plan *SettingWhitelabelModel) sdk.UpdateWhitelabelSettingsRequest {
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

func mapResponseToModel(
	model *SettingWhitelabelModel,
	settings *sdk.ListWhitelabelSettings200ResponseWhitelabelSettings,
) {
	model.Id = types.StringValue("settings")
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
	// header_logo, footer_logo, login_logo and favicon are local file paths
	// uploaded via the multipart images endpoint. The API returns a
	// server-generated storage URL for them (never the path the user supplied),
	// so they are intentionally NOT mapped back from the response; doing so would
	// clobber the configured path and produce an inconsistent-state error. Their
	// values are carried through from the plan/prior state instead.
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
