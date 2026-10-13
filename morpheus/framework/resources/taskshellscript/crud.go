// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package taskshellscript

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HPE/terraform-provider-hpe/utils/schemadefaults"
)

func (r *Resource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan, config TaskShellScriptModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("creating client failed", err.Error())

		return
	}

	kind.Create(ctx, client, req, resp, toFlat(plan), toFlat(config), setState(&resp.State))
}

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var data TaskShellScriptModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("creating client failed", err.Error())

		return
	}

	kind.Read(ctx, client, resp, toFlat(data), setState(&resp.State))
	if resp.Diagnostics.HasError() || resp.State.Raw.IsNull() {
		return
	}

	// Fill any schema-declared default the API omitted (null in state) so an
	// imported resource does not plan a change nobody made. MORPH-16192.
	resp.Diagnostics.Append(
		schemadefaults.Apply(ctx, TaskShellScriptResourceSchema(ctx), &resp.State)...,
	)
}

func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan, state, config TaskShellScriptModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("creating client failed", err.Error())

		return
	}

	kind.Update(ctx, client, req, resp, toFlat(plan), toFlat(state), toFlat(config), setState(&resp.State))
}

func (r *Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var data TaskShellScriptModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("creating client failed", err.Error())

		return
	}

	kind.Delete(ctx, client, resp, toFlat(data))
}
