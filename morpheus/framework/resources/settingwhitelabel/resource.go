package settingwhitelabel

import (
	"context"
	"fmt"
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

	// Upload any logo/favicon image files via the multipart images endpoint.
	if err := r.uploadImages(ctx, client, &plan); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpCreate, "setting_whitelabel", "images", err, nil)

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

	var plan SettingWhitelabelModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
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

	// Upload any logo/favicon image files via the multipart images endpoint.
	if err := r.uploadImages(ctx, client, &plan); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpUpdate, "setting_whitelabel", "images", err, nil)

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

// uploadImages uploads the logo/favicon image files referenced by the plan to
// the multipart whitelabel images endpoint.
//
// The four image attributes (header_logo, footer_logo, login_logo, favicon) are
// local file paths. The API stores the uploaded bytes and, on read, returns a
// server-generated storage URL that never matches the supplied path, so these
// values are carried through from the plan and are never reconciled from the
// read response. Clearing a path does not reset the image on the appliance; the
// resource must be destroyed (which resets all images) to remove one.
func (r *settingWhitelabelResource) uploadImages(
	ctx context.Context,
	client *sdk.APIClient,
	plan *SettingWhitelabelModel,
) error {
	// The whitelabel images endpoint returns an error for an empty upload, so
	// there is nothing to do when no image paths are set.
	if plan.HeaderLogo.IsNull() && plan.FooterLogo.IsNull() &&
		plan.LoginLogo.IsNull() && plan.Favicon.IsNull() {
		return nil
	}

	imgReq := client.WhitelabelSettingsAPI.UpdateWhitelabelImages(ctx)

	if !plan.HeaderLogo.IsNull() && !plan.HeaderLogo.IsUnknown() {
		headerLogo, err := openImageFile(plan.HeaderLogo)
		if err != nil {
			return err
		}
		if headerLogo != nil {
			defer headerLogo.Close()
			imgReq = imgReq.HeaderLogoFile(headerLogo)
		}
	}

	if !plan.FooterLogo.IsNull() && !plan.FooterLogo.IsUnknown() {
		footerLogo, err := openImageFile(plan.FooterLogo)
		if err != nil {
			return err
		}
		if footerLogo != nil {
			defer footerLogo.Close()
			imgReq = imgReq.FooterLogoFile(footerLogo)
		}
	}

	if !plan.LoginLogo.IsNull() && !plan.LoginLogo.IsUnknown() {
		loginLogo, err := openImageFile(plan.LoginLogo)
		if err != nil {
			return err
		}
		if loginLogo != nil {
			defer loginLogo.Close()
			imgReq = imgReq.LoginLogoFile(loginLogo)
		}
	}

	if !plan.Favicon.IsNull() && !plan.Favicon.IsUnknown() {
		favicon, err := openImageFile(plan.Favicon)
		if err != nil {
			return err
		}
		if favicon != nil {
			defer favicon.Close()
			imgReq = imgReq.FaviconFile(favicon)
		}
	}

	_, httpResp, err := imgReq.Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		return fmt.Errorf("uploading whitelabel images: %w", err)
	}

	return nil
}

// openImageFile opens the image file at path for upload. The caller is
// responsible for closing the returned file.
func openImageFile(path types.String) (*os.File, error) {
	f, err := os.Open(path.ValueString())
	if err != nil {
		return nil, fmt.Errorf("could not read whitelabel image file: %w", err)
	}

	return f, nil
}

func buildUpdateRequest(plan *SettingWhitelabelModel) sdk.UpdateWhitelabelSettingsRequest {
	settings := sdk.UpdateWhitelabelSettingsRequestWhitelabelSettings{}
	if !plan.Enabled.IsNull() && !plan.Enabled.IsUnknown() {
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
