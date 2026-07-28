// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package cloud

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/specs"
	"github.com/HPE/terraform-provider-hpe/utils/cleanup"
	"github.com/HPE/terraform-provider-hpe/utils/translate"
)

const (
	resourceName    = "cloud"
	createOperation = "create cloud resource"
	readOperation   = "read cloud resource"
)

// newTranslateClient creates a translate.Client from the SDK client.
func (r *Resource) newTranslateClient(ctx context.Context) (*translate.Client, error) {
	sdkClient, err := r.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("creating SDK client: %w", err)
	}

	cfg := sdkClient.GetConfig()
	baseURL := cfg.Servers[0].URL
	tc := translate.FromSDKClient(baseURL, cfg.HTTPClient)

	data, err := specs.ResourceConfig(resourceName)
	if err != nil {
		return nil, fmt.Errorf("loading cloud config: %w", err)
	}

	if err := tc.RegisterResource(resourceName, data,
		translate.WithPostRead(cloudPostRead),
	); err != nil {
		return nil, fmt.Errorf("registering cloud config: %w", err)
	}

	return tc, nil
}

func (r *Resource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan CloudModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Write-only fields are stripped from the plan by Terraform.
	// Source them from req.Config instead.
	var config CloudModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	mergeWriteOnlyFields(&plan, &config)

	tc, err := r.newTranslateClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(createOperation, err.Error())

		return
	}

	result, err := tc.Execute(ctx, translate.Request{
		Operation: translate.Create,
		Resource:  resourceName,
		Model:     &plan,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			createOperation,
			"cloud "+plan.Name.ValueString()+" POST failed: "+err.Error(),
		)

		return
	}

	id, err := result.ID()
	if err != nil {
		resp.Diagnostics.AddError(createOperation, "failed to extract ID: "+err.Error())

		return
	}

	plan.Id = types.Int64Value(id)

	readResult, err := tc.Execute(ctx, translate.Request{
		Operation: translate.Read,
		Resource:  resourceName,
		ID:        id,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			createOperation,
			fmt.Sprintf("Cloud %d was created but could not be read: %s", id, err.Error()),
		)
		cleanup.TaintResourceState(ctx, cleanup.TaintResourceStateConfig{
			ResourceType: resourceName,
			ResourceID:   id,
			StateWriter:  &resp.State,
			Diagnostics:  &resp.Diagnostics,
		})

		return
	}

	var state CloudModel
	if err := readResult.IntoWithPlan(ctx, &state, &plan); err != nil {
		resp.Diagnostics.AddError(
			createOperation,
			fmt.Sprintf("Cloud %d was created but state could not be parsed: %s", id, err.Error()),
		)
		cleanup.TaintResourceState(ctx, cleanup.TaintResourceStateConfig{
			ResourceType: resourceName,
			ResourceID:   id,
			StateWriter:  &resp.State,
			Diagnostics:  &resp.Diagnostics,
		})

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var data CloudModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := data.Id.ValueInt64()

	tc, err := r.newTranslateClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(readOperation, err.Error())

		return
	}

	result, err := tc.Execute(ctx, translate.Request{
		Operation: translate.Read,
		Resource:  resourceName,
		ID:        id,
	})
	if err != nil {
		if result != nil && result.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)

			return
		}

		resp.Diagnostics.AddError(
			readOperation,
			fmt.Sprintf("cloud %d GET failed: %s", id, err.Error()),
		)

		return
	}

	var state CloudModel
	if err := result.IntoWithPlan(ctx, &state, &data); err != nil {
		resp.Diagnostics.AddError(readOperation, err.Error())

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan CloudModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Write-only fields are stripped from the plan by Terraform.
	var config CloudModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	mergeWriteOnlyFields(&plan, &config)

	id := plan.Id.ValueInt64()

	tc, err := r.newTranslateClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("update cloud resource", err.Error())

		return
	}

	_, err = tc.Execute(ctx, translate.Request{
		Operation: translate.Update,
		Resource:  resourceName,
		Model:     &plan,
		ID:        id,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"update cloud resource",
			fmt.Sprintf("cloud %d PUT failed: %s", id, err.Error()),
		)

		return
	}

	readResult, err := tc.Execute(ctx, translate.Request{
		Operation: translate.Read,
		Resource:  resourceName,
		ID:        id,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"update cloud resource",
			fmt.Sprintf("cloud %d was updated but could not be read: %s", id, err.Error()),
		)

		return
	}

	var state CloudModel
	if err := readResult.IntoWithPlan(ctx, &state, &plan); err != nil {
		resp.Diagnostics.AddError("update cloud resource", err.Error())

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var data CloudModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := data.Id.ValueInt64()

	tc, err := r.newTranslateClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("delete cloud resource", err.Error())

		return
	}

	_, err = tc.Execute(ctx, translate.Request{
		Operation:   translate.Delete,
		Resource:    resourceName,
		ID:          id,
		QueryParams: map[string]string{"force": "true"},
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"delete cloud resource",
			fmt.Sprintf("cloud %d DELETE failed: %s", id, err.Error()),
		)

		return
	}
}

// mergeWriteOnlyFields copies write-only field values from the config model
// (req.Config) into the plan model. Terraform strips write-only values from
// the plan, so they must be sourced from the raw config.
func mergeWriteOnlyFields(plan *CloudModel, config *CloudModel) {
	// VMware password
	if !config.ConfigVmware.IsNull() && !config.ConfigVmware.IsUnknown() {
		if !config.ConfigVmware.Password.IsNull() {
			plan.ConfigVmware.Password = config.ConfigVmware.Password
		}
	}

	// Azure client secret
	if !config.ConfigAzure.IsNull() && !config.ConfigAzure.IsUnknown() {
		if !config.ConfigAzure.ClientSecret.IsNull() {
			plan.ConfigAzure.ClientSecret = config.ConfigAzure.ClientSecret
		}
	}
}
