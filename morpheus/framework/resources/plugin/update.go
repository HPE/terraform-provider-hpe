// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package plugin

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state PluginModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateTimeout, diags := plan.Timeouts.Update(ctx, 10*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, updateTimeout)
	defer cancel()

	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	id := state.Id.ValueInt64()
	plan.Id = state.Id

	// Re-upload when the jar source or checksum changed.
	jarChanged := !plan.JarPath.Equal(state.JarPath) ||
		!plan.JarUrl.Equal(state.JarUrl) ||
		!plan.JarSha256.Equal(state.JarSha256)

	newSHA := state.JarSha256

	if jarChanged {
		jar, err := resolveJar(ctx, plan.JarPath.ValueString(), plan.JarUrl.ValueString(),
			configuredSHA(plan))
		if err != nil {
			errfmt.DiagError(&resp.Diagnostics, errfmt.OpUpdate, resourceName, fmt.Sprint(id), err, nil)

			return
		}
		defer jar.cleanup()

		uploadResp, httpResp, err := client.PluginsAPI.UploadPlugin(ctx).File(jar.file).Execute()
		if err := errfmt.CheckResponse(err, httpResp); err != nil {
			errfmt.DiagError(&resp.Diagnostics, errfmt.OpUpdate, resourceName, fmt.Sprint(id), err, httpResp)

			return
		}
		if uploadResp == nil || uploadResp.Plugin == nil {
			resp.Diagnostics.AddError("update plugin resource", "the upload response did not include a plugin")

			return
		}

		// A jar with a different code installs a separate plugin rather than
		// updating this one; refuse and leave state untouched.
		uploadedCode := ""
		if uploadResp.Plugin.Code != nil {
			uploadedCode = *uploadResp.Plugin.Code
		}
		wantCode := state.Code.ValueString()
		if uploadedCode != wantCode {
			resp.Diagnostics.AddError(
				"update plugin resource",
				fmt.Sprintf("the uploaded jar is plugin code %q but this resource manages plugin "+
					"code %q; a jar with a different code installs a separate plugin instead of "+
					"updating this one", uploadedCode, wantCode),
			)

			return
		}

		if _, err := r.pollLoaded(ctx, client, id, updateTimeout); err != nil {
			errfmt.DiagError(&resp.Diagnostics, errfmt.OpUpdate, resourceName, fmt.Sprint(id), err, nil)

			return
		}

		newSHA = types.StringValue(jar.sha256)
	}

	// Apply enabled/config when they changed.
	enabledChanged := !plan.Enabled.Equal(state.Enabled)
	configChanged := !plan.Config.Equal(state.Config)

	if enabledChanged || configChanged {
		if err := r.applyEnabledConfig(ctx, client, id, plan, enabledChanged); err != nil {
			errfmt.DiagError(&resp.Diagnostics, errfmt.OpUpdate, resourceName, fmt.Sprint(id), err, nil)

			return
		}
	}

	apiResp, httpResp, err := client.PluginsAPI.GetPlugin(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, resourceName, fmt.Sprint(id), err, httpResp)

		return
	}
	// CheckResponse passes a 200 with an empty body. The prior state still
	// describes the plugin, so leave it for the next plan to refresh.
	if apiResp == nil || apiResp.Plugin == nil {
		resp.Diagnostics.AddError("update plugin resource",
			fmt.Sprintf("plugin %d: API returned an empty plugin", id))

		return
	}

	plan.JarSha256 = newSHA
	newState, diags := getPluginAsState(ctx, apiResp.Plugin, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

// configToMap converts a plan config Dynamic value into the free-form map the
// API accepts.
func configToMap(ctx context.Context, cfg types.Dynamic) (map[string]any, error) {
	v, err := convert.ValueToAny(ctx, cfg.UnderlyingValue())
	if err != nil {
		return nil, fmt.Errorf("converting plugin config: %w", err)
	}

	m, ok := v.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("plugin config must be an object")
	}

	return m, nil
}
