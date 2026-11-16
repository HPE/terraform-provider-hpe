// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package appliancesetup

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/auth"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
)

// Create implements resource.Resource.
//
// It waits for the appliance's API to answer, then initialises the appliance
// when it reports that setup is needed and adopts it otherwise. The password
// and licence key are write-only, so they are read from the configuration.
func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, config ApplianceSetupModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createTimeout, diags := plan.Timeouts.Create(ctx, defaultCreateTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, createTimeout)
	defer cancel()

	unauth, err := r.NewUnauthenticated()
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	setup := newSetupClient(unauth)

	fail := func(err error) {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpCreate, resourceName, unauth.BaseURL, err, nil)
	}

	check, err := setup.awaitAPI(ctx)
	if err != nil {
		fail(err)

		return
	}

	plan.Id = types.StringValue(unauth.BaseURL)
	if plan.ApplianceURL.IsUnknown() || plan.ApplianceURL.IsNull() {
		plan.ApplianceURL = types.StringValue(unauth.BaseURL)
	}

	plan.AlreadyInitialised = types.BoolValue(!check.SetupNeeded)
	plan.BuildVersion = types.StringNull()

	if check.SetupNeeded {
		result, err := setup.init(ctx, payloadFor(plan, config))

		switch {
		case errors.Is(err, errAlreadySetUp):
			plan.AlreadyInitialised = types.BoolValue(true)
		case err != nil:
			fail(err)

			return
		default:
			if result.BuildVersion != "" {
				plan.BuildVersion = types.StringValue(result.BuildVersion)
			}

			tflog.Info(ctx, "appliance initialised", map[string]any{"appliance_url": plan.ApplianceURL.ValueString()})
		}
	}

	if plan.AlreadyInitialised.ValueBool() {
		resp.Diagnostics.AddWarning(
			"Appliance already set up",
			fmt.Sprintf("The appliance at %s had already been set up, so nothing was sent: the configured tenant, "+
				"user and licence key may not be the ones on the appliance. The resource adopts it as it is.",
				unauth.BaseURL),
		)
	}

	if plan.BuildVersion.IsNull() {
		plan.BuildVersion = r.buildVersion(ctx, plan.BuildVersion)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// payloadFor builds the setup request from the planned inputs and the
// write-only values read from the configuration.
func payloadFor(plan, config ApplianceSetupModel) initPayload {
	hubMode := plan.HubMode.ValueString()
	if hubMode == "" {
		hubMode = hubModeSkip
	}

	return initPayload{
		ApplianceURL:  plan.ApplianceURL.ValueString(),
		ApplianceName: plan.ApplianceName.ValueString(),
		AccountName:   plan.AccountName.ValueString(),
		Username:      plan.Username.ValueString(),
		Password:      config.PasswordWo.ValueString(),
		Email:         plan.Email.ValueString(),
		FirstName:     plan.FirstName.ValueString(),
		LastName:      plan.LastName.ValueString(),
		HubMode:       hubMode,
		LicenseKey:    config.LicenseKeyWo.ValueString(),
	}
}

// Read implements resource.Resource.
//
// The setup check is the oracle: an appliance that again reports setupNeeded
// has been reinstalled behind the same URL, and the resource is removed from state so
// that the next apply sets it up again. An appliance whose API does not
// answer is an error, not a removal: setup must not be re-run blind.
func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	var state ApplianceSetupModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	unauth, err := r.NewUnauthenticated()
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	check, err := newSetupClient(unauth).check(ctx)
	if err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, resourceName, state.Id.ValueString(), err, nil)

		return
	}

	if check.SetupNeeded {
		tflog.Warn(ctx, "appliance reports that setup is needed again; removing the resource from state",
			map[string]any{"appliance_url": state.Id.ValueString()})
		resp.State.RemoveResource(ctx)

		return
	}

	state.BuildVersion = r.buildVersion(ctx, state.BuildVersion)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// buildVersion reads the appliance's build version with the provider's
// credentials. The unauthenticated check blanks the version, so this is the
// only way to see it; when the call fails (the credentials are not the ones
// setup created, say) the previous value is kept.
func (r *Resource) buildVersion(ctx context.Context, previous types.String) types.String {
	client, err := r.NewClient(ctx)
	if err != nil {
		return previous
	}

	whoami, httpResp, err := client.AuthenticationAPI.Whoami(ctx).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		tflog.Debug(ctx, "build version not read: whoami failed", map[string]any{"error": err.Error()})

		return previous
	}

	if whoami == nil || whoami.Appliance == nil || whoami.Appliance.BuildVersion == nil ||
		*whoami.Appliance.BuildVersion == "" {
		return previous
	}

	return types.StringValue(*whoami.Appliance.BuildVersion)
}

// Update implements resource.Resource.
func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state ApplianceSetupModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(validateSetupChanges(plan, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.Id = state.Id
	plan.AlreadyInitialised = state.AlreadyInitialised
	plan.BuildVersion = state.BuildVersion

	if plan.ApplianceURL.IsUnknown() {
		plan.ApplianceURL = state.ApplianceURL
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete implements resource.Resource. Setup cannot be undone; the resource
// is forgotten and the appliance is left as it is.
func (r *Resource) Delete(ctx context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
	tflog.Info(ctx, "hpe_morpheus_appliance_setup destroyed: the appliance is left set up")
}

// ImportState implements resource.ResourceWithImportState. The import ID is
// the appliance URL, which must be the one the provider is configured with;
// the appliance must already be set up. The required inputs are not known to
// the appliance and are taken from the configuration on the next apply.
func (r *Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	unauth, err := r.NewUnauthenticated()
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	id := auth.NormalizeBaseURL(req.ID)
	if id != unauth.BaseURL {
		resp.Diagnostics.AddError("invalid import ID",
			fmt.Sprintf("expected the provider's appliance URL %q, got %q", unauth.BaseURL, req.ID))

		return
	}

	check, err := newSetupClient(unauth).check(ctx)
	if err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpImport, resourceName, id, err, nil)

		return
	}

	if check.SetupNeeded {
		errfmt.DiagErrorf(&resp.Diagnostics, errfmt.OpImport, resourceName,
			"the appliance at %s has not been set up yet; apply the resource instead of importing it", id)

		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("appliance_url"), id)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("hub_mode"), hubModeSkip)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("already_initialised"), true)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("build_version"),
		r.buildVersion(ctx, types.StringNull()))...)
}
