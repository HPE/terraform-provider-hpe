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

	var plan, config SettingWhitelabelModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
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

	// Upload logo/favicon image files (write-only paths, read from config) via
	// the multipart images endpoint. On create there is no prior state, so every
	// provided image is uploaded.
	r.applyImageChanges(ctx, client, &config, &plan, nil, &resp.Diagnostics)
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

	var plan, state, config SettingWhitelabelModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
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

	// Upload logo/favicon files whose *_wo_version changed, and reset any whose
	// version changed with no file supplied.
	r.applyImageChanges(ctx, client, &config, &plan, &state, &resp.Diagnostics)
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

// applyImageChanges uploads whitelabel logo/favicon image files and resets any
// that were cleared. The image paths are write-only attributes
// (header_logo_wo, footer_logo_wo, login_logo_wo, favicon_wo), so they are read
// from the config, not the plan/state (write-only values are always null in
// plan and state). A given image is acted on only when its *_wo_version changed
// relative to the prior state; on Create (state == nil) every provided image is
// uploaded.
//
// The image paths are local files. The API stores the uploaded bytes and, on
// read, returns a server-generated storage URL that never matches the supplied
// path, so nothing is reconciled back onto these attributes.
func (r *settingWhitelabelResource) applyImageChanges(
	ctx context.Context,
	client *sdk.APIClient,
	config *SettingWhitelabelModel,
	plan *SettingWhitelabelModel,
	state *SettingWhitelabelModel,
	diagnostics *diag.Diagnostics,
) {
	imgReq := client.WhitelabelSettingsAPI.UpdateWhitelabelImages(ctx)
	resets := sdk.UpdateWhitelabelSettingsRequestWhitelabelSettings{}

	var openFiles []*os.File
	defer func() {
		for _, f := range openFiles {
			_ = f.Close()
		}
	}()

	upload := false
	needReset := false

	// apply decides what to do for a single image. It returns false only when a
	// fatal error was recorded in diagnostics.
	apply := func(
		name string,
		cfgPath types.String,
		planVersion types.Int64,
		stateVersion types.Int64,
		attach func(*os.File),
		setReset func(),
	) bool {
		// On update, act only when the version changed; on create (state == nil)
		// every provided image is applied.
		triggered := state == nil || !planVersion.Equal(stateVersion)
		if !triggered {
			return true
		}

		switch {
		case cfgPath.IsUnknown():
			// On update a bumped version with no resolvable path is a user error;
			// on create an unset path simply means "no image provided".
			if state != nil {
				diagnostics.AddError(
					"update setting_whitelabel",
					name+"_wo_version changed but "+name+"_wo is not set",
				)

				return false
			}

			return true
		case cfgPath.IsNull():
			// No file supplied. On update, a bumped version resets the image; on
			// create there is nothing to reset.
			if state != nil {
				setReset()
				needReset = true
			}

			return true
		default:
			f, err := os.Open(cfgPath.ValueString())
			if err != nil {
				diagnostics.AddError("could not read whitelabel image file", err.Error())

				return false
			}
			openFiles = append(openFiles, f)
			attach(f)
			upload = true

			return true
		}
	}

	// stateVersionOf returns the prior version for an image, or null on create.
	stateVersionOf := func(get func(*SettingWhitelabelModel) types.Int64) types.Int64 {
		if state == nil {
			return types.Int64Null()
		}

		return get(state)
	}

	if !apply(
		"header_logo", config.HeaderLogoWo,
		plan.HeaderLogoWoVersion,
		stateVersionOf(func(m *SettingWhitelabelModel) types.Int64 { return m.HeaderLogoWoVersion }),
		func(f *os.File) { imgReq = imgReq.HeaderLogoFile(f) },
		func() { resets.ResetHeaderLogo = boolPtr(true) },
	) {
		return
	}
	if !apply(
		"footer_logo", config.FooterLogoWo,
		plan.FooterLogoWoVersion,
		stateVersionOf(func(m *SettingWhitelabelModel) types.Int64 { return m.FooterLogoWoVersion }),
		func(f *os.File) { imgReq = imgReq.FooterLogoFile(f) },
		func() { resets.ResetFooterLogo = boolPtr(true) },
	) {
		return
	}
	if !apply(
		"login_logo", config.LoginLogoWo,
		plan.LoginLogoWoVersion,
		stateVersionOf(func(m *SettingWhitelabelModel) types.Int64 { return m.LoginLogoWoVersion }),
		func(f *os.File) { imgReq = imgReq.LoginLogoFile(f) },
		func() { resets.ResetLoginLogo = boolPtr(true) },
	) {
		return
	}
	if !apply(
		"favicon", config.FaviconWo,
		plan.FaviconWoVersion,
		stateVersionOf(func(m *SettingWhitelabelModel) types.Int64 { return m.FaviconWoVersion }),
		func(f *os.File) { imgReq = imgReq.FaviconFile(f) },
		func() { resets.ResetFavicon = boolPtr(true) },
	) {
		return
	}

	if upload {
		_, httpResp, err := imgReq.Execute()
		if err := errfmt.CheckResponse(err, httpResp); err != nil {
			errfmt.DiagError(diagnostics, errfmt.OpUpdate, "setting_whitelabel", "images", err, httpResp)

			return
		}
	}

	if !needReset {
		return
	}

	// Route resets through the JSON settings endpoint, which honours the reset
	// flags (the multipart images endpoint's favicon reset field is misnamed in
	// the spec).
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
	// header_logo, footer_logo, login_logo and favicon are write-only local file
	// paths uploaded via the multipart images endpoint. The API returns a
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
