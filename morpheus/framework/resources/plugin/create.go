// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package plugin

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/cleanup"
)

func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan PluginModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createTimeout, diags := plan.Timeouts.Create(ctx, 10*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, createTimeout)
	defer cancel()

	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	taint := func(id int64) {
		cleanup.TaintResourceState(ctx, cleanup.TaintResourceStateConfig{
			ResourceType: resourceName,
			ResourceID:   id,
			StateWriter:  &resp.State,
			Diagnostics:  &resp.Diagnostics,
		})
	}

	// Resolve, verify and upload the jar.
	jar, err := resolveJar(ctx, plan.JarPath.ValueString(), plan.JarUrl.ValueString(),
		configuredSHA(plan))
	if err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpCreate, resourceName, "", err, nil)

		return
	}
	defer jar.cleanup()

	uploadResp, httpResp, err := client.PluginsAPI.UploadPlugin(ctx).File(jar.file).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpCreate, resourceName, "", err, httpResp)

		return
	}
	if uploadResp == nil || uploadResp.Plugin == nil || uploadResp.Plugin.Id == nil {
		resp.Diagnostics.AddError("create plugin resource",
			"the upload response did not include a plugin id")

		return
	}

	id := *uploadResp.Plugin.Id

	// Wait for the plugin to finish loading.
	if _, err := r.pollLoaded(ctx, client, id, createTimeout); err != nil {
		taint(id)
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpCreate, resourceName, fmt.Sprint(id), err, nil)

		return
	}

	// The upload leaves enabled=true, config={}. Apply enabled/config when the
	// plan asks for anything other than that default.
	if needsPostUploadUpdate(plan) {
		if err := r.applyEnabledConfig(ctx, client, id, plan, true); err != nil {
			taint(id)
			errfmt.DiagError(&resp.Diagnostics, errfmt.OpCreate, resourceName, fmt.Sprint(id), err, nil)

			return
		}
	}

	// Read back the final state.
	apiResp, httpResp, err := client.PluginsAPI.GetPlugin(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		taint(id)
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, resourceName, fmt.Sprint(id), err, httpResp)

		return
	}
	// CheckResponse passes a 200 with an empty body; the plugin is installed
	// either way, so keep its id in state for the next plan to read.
	if apiResp == nil || apiResp.Plugin == nil {
		taint(id)
		resp.Diagnostics.AddError("create plugin resource",
			fmt.Sprintf("plugin %d: API returned an empty plugin", id))

		return
	}

	plan.JarSha256 = types.StringValue(jar.sha256)
	state, diags := getPluginAsState(ctx, apiResp.Plugin, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// configuredSHA returns the user-set jar_sha256 (empty when not configured).
// A computed-only value is unknown at create time, so it is treated as unset.
func configuredSHA(plan PluginModel) string {
	if plan.JarSha256.IsNull() || plan.JarSha256.IsUnknown() {
		return ""
	}

	return plan.JarSha256.ValueString()
}

// needsPostUploadUpdate reports whether a follow-up UpdatePlugin is needed after
// the upload: enabled explicitly false, or a config value supplied.
func needsPostUploadUpdate(plan PluginModel) bool {
	enabledFalse := !plan.Enabled.IsNull() && !plan.Enabled.IsUnknown() && !plan.Enabled.ValueBool()
	hasConfig := !plan.Config.IsNull() && !plan.Config.IsUnknown()

	return enabledFalse || hasConfig
}

// applyEnabledConfig sends an UpdatePlugin with enabled and/or config. When
// sendEnabled is true the enabled flag is always sent; otherwise it is omitted.
func (r *Resource) applyEnabledConfig(
	ctx context.Context,
	client *sdk.APIClient,
	id int64,
	plan PluginModel,
	sendEnabled bool,
) error {
	body := sdk.UpdatePluginRequestPlugin{}

	if sendEnabled && !plan.Enabled.IsNull() && !plan.Enabled.IsUnknown() {
		body.Enabled = plan.Enabled.ValueBoolPointer()
	}

	if !plan.Config.IsNull() && !plan.Config.IsUnknown() {
		cfg, err := configToMap(ctx, plan.Config)
		if err != nil {
			return err
		}
		body.Config = cfg
	}

	_, httpResp, err := client.PluginsAPI.UpdatePlugin(ctx, id).
		UpdatePluginRequest(sdk.UpdatePluginRequest{Plugin: body}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		return fmt.Errorf("applying plugin enabled/config: %w", err)
	}

	// Enabling or configuring can briefly flip the status; wait for it to
	// settle back to loaded.
	_, err = r.pollLoaded(ctx, client, id, settleTimeout)

	return err
}
