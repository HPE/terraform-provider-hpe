// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package plugin

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/cenkalti/backoff/v5"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
	"github.com/HPE/terraform-provider-hpe/utils/schemadefaults"
)

func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state PluginModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	id := state.Id.ValueInt64()
	apiResp, httpResp, err := client.PluginsAPI.GetPlugin(ctx, id).Execute()
	if errfmt.IsNotFound(httpResp) {
		resp.State.RemoveResource(ctx)

		return
	}
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, resourceName, fmt.Sprint(id), err, httpResp)

		return
	}
	if apiResp == nil || apiResp.Plugin == nil {
		resp.Diagnostics.AddError("read plugin resource",
			fmt.Sprintf("plugin %d: API returned an empty plugin", id))

		return
	}

	newState, diags := getPluginAsState(ctx, apiResp.Plugin, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Fill any schema-declared default the API omitted (null in state) so an
	// imported resource does not plan a change nobody made. MORPH-16192.
	resp.Diagnostics.Append(
		schemadefaults.Apply(ctx, PluginResourceSchema(ctx), &resp.State)...,
	)
}

// getPluginAsState maps an API plugin onto a PluginModel, preserving from prior
// the fields the API does not return (jar_path/jar_url/jar_sha256) and the
// user's configured config value to avoid perpetual diffs.
func getPluginAsState(
	ctx context.Context,
	p *sdk.GetPlugin200ResponseAllOfPlugin,
	prior PluginModel,
) (PluginModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	state := PluginModel{
		Id:                    convert.Int64ToType(p.Id),
		Name:                  convert.StrToType(p.Name),
		Code:                  convert.StrToType(p.Code),
		Description:           convert.StrToType(p.Description),
		Version:               convert.StrToType(p.Version),
		Status:                convert.StrToType(p.Status),
		StatusMessage:         convert.StrToType(p.StatusMessage.Get()),
		Author:                convert.StrToType(p.Author.Get()),
		WebsiteUrl:            convert.StrToType(p.WebsiteUrl.Get()),
		SourceCodeLocationUrl: convert.StrToType(p.SourceCodeLocationUrl.Get()),
		IssueTrackerUrl:       convert.StrToType(p.IssueTrackerUrl.Get()),
		Enabled:               convert.BoolToType(p.Enabled),
		Valid:                 convert.BoolToType(p.Valid),
		HasValidUpdate:        convert.BoolToType(p.HasValidUpdate),
		IsPrepackaged:         convert.BoolToType(p.IsPrepackaged),
		DateCreated:           convert.TimeToType(p.DateCreated),
		LastUpdated:           convert.TimeToType(p.LastUpdated),

		// Not returned by the API — carried through from prior state/plan.
		JarPath:   prior.JarPath,
		JarUrl:    prior.JarUrl,
		JarSha256: prior.JarSha256,
		Timeouts:  prior.Timeouts,
	}

	// providers: the API returns them in a different order on every read, and
	// {type, name} is not unique (a plugin may contribute several UI tabs with
	// the same name), so keep the list but give it a stable order.
	providers := make([]ProvidersValue, 0, len(p.Providers))
	for _, pr := range p.Providers {
		providers = append(providers, ProvidersValue{
			Name:          convert.StrToType(pr.Name),
			ProvidersType: convert.StrToType(pr.Type),
			state:         attr.ValueStateKnown,
		})
	}
	sortProviders(providers)
	providersList, d := types.ListValueFrom(ctx, ProvidersValue{}.Type(ctx), providers)
	diags.Append(d...)
	state.Providers = providersList

	// config: keep the user's configured value when present; when they set
	// nothing and the API returns {} (or nothing), keep null to avoid a diff.
	switch {
	case !prior.Config.IsNull() && !prior.Config.IsUnknown():
		state.Config = prior.Config
	case len(p.Config) == 0:
		state.Config = types.DynamicNull()
	default:
		cfg, err := convert.MapToDynamic(ctx, p.Config)
		if err != nil {
			diags.AddError("read plugin resource", "could not convert plugin config: "+err.Error())

			return state, diags
		}
		state.Config = cfg
	}

	return state, diags
}

// pollLoaded polls GetPlugin until the plugin reaches status "loaded", returning
// the loaded plugin. A status of "error" fails with the status message;
// "installing"/"updating"/empty are treated as in-progress.
func (r *Resource) pollLoaded(
	ctx context.Context,
	client *sdk.APIClient,
	id int64,
	timeout time.Duration,
) (*sdk.GetPlugin200ResponseAllOfPlugin, error) {
	what := fmt.Sprintf("plugin %d", id)

	op := func() (*sdk.GetPlugin200ResponseAllOfPlugin, error) {
		apiResp, httpResp, err := client.PluginsAPI.GetPlugin(ctx, id).Execute()
		if err := errfmt.CheckResponse(err, httpResp); err != nil {
			return nil, backoff.Permanent(fmt.Errorf("%s: GET failed: %w", what, err))
		}
		if apiResp == nil || apiResp.Plugin == nil {
			return nil, fmt.Errorf("%s: GET returned an empty plugin", what)
		}

		status := ""
		if apiResp.Plugin.Status != nil {
			status = *apiResp.Plugin.Status
		}

		switch status {
		case statusLoaded:
			return apiResp.Plugin, nil
		case statusError:
			msg := apiResp.Plugin.StatusMessage.Get()
			detail := ""
			if msg != nil {
				detail = ": " + *msg
			}

			return nil, backoff.Permanent(fmt.Errorf("%s: failed to load%s", what, detail))
		default:
			return nil, fmt.Errorf("%s: still loading (status %q)", what, status)
		}
	}

	return backoff.Retry(
		ctx,
		op,
		backoff.WithBackOff(backoff.NewConstantBackOff(3*time.Second)),
		backoff.WithMaxElapsedTime(timeout),
	)
}

// sortProviders orders providers by type, then name, so the list is the same
// on every read regardless of the order the API returns it in.
func sortProviders(providers []ProvidersValue) {
	sort.SliceStable(providers, func(i, j int) bool {
		ti, tj := providers[i].ProvidersType.ValueString(), providers[j].ProvidersType.ValueString()
		if ti != tj {
			return ti < tj
		}

		return providers[i].Name.ValueString() < providers[j].Name.ValueString()
	})
}
