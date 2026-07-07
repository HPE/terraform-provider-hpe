// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package deployment

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	providererrors "github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

func deploymentAsState(
	dep *sdk.GetDeployment200ResponseDeployment,
) DeploymentModel {
	state := DeploymentModel{
		Id:           convert.Int64ToType(dep.Id),
		Name:         convert.StrToType(dep.Name),
		Description:  convert.StrToType(dep.Description),
		AccountId:    convert.Int64ToType(dep.AccountId),
		DateCreated:  convert.TimeToType(dep.DateCreated),
		LastUpdated:  convert.TimeToType(dep.LastUpdated),
		VersionCount: convert.Int64ToType(dep.VersionCount),
	}

	// Nullable fields
	if dep.ExternalId.IsSet() {
		state.ExternalId = convert.StrToType(dep.ExternalId.Get())
	} else {
		state.ExternalId = types.StringNull()
	}

	return state
}

func getByID(
	ctx context.Context,
	id int64,
	apiClient *sdk.APIClient,
) (*sdk.GetDeployment200ResponseDeployment, error) {
	r, hresp, err := apiClient.DeploymentsAPI.GetDeployment(ctx, id).Execute()
	if r == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET failed for deployment %d: %s", id, providererrors.ErrMsg(err, hresp))
	}

	if r.Deployment == nil {
		return nil, fmt.Errorf("GET failed for deployment %d: response missing deployment", id)
	}

	dep := *r.Deployment

	return &dep, nil
}

func getByName(
	ctx context.Context,
	name string,
	apiClient *sdk.APIClient,
) (*sdk.GetDeployment200ResponseDeployment, error) {
	rs, hresp, err := apiClient.DeploymentsAPI.ListDeployments(ctx).Name(name).Execute()
	if rs == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET failed for deployment %s: %s", name, providererrors.ErrMsg(err, hresp))
	}

	var matched []sdk.ListDeployments200ResponseAllOfDeploymentsInner

	for _, o := range rs.Deployments {
		if o.Name != nil && *o.Name == name {
			matched = append(matched, o)
		}
	}

	if len(matched) == 0 {
		return nil, errors.New(ErrorNoDeploymentFound)
	} else if len(matched) > 1 {
		return nil, errors.New(ErrorMultipleDeployments)
	}

	if matched[0].Id == nil {
		return nil, fmt.Errorf("GET failed for deployment %s: response missing id", name)
	}

	return getByID(ctx, *matched[0].Id, apiClient)
}

func getDeployment(
	ctx context.Context,
	config *DeploymentModel,
	apiClient *sdk.APIClient,
) (*sdk.GetDeployment200ResponseDeployment, error) {
	if !config.Id.IsNull() {
		return getByID(ctx, config.Id.ValueInt64(), apiClient)
	} else if !config.Name.IsNull() {
		return getByName(ctx, config.Name.ValueString(), apiClient)
	}

	return nil, errors.New(ErrorNoValidSearchTerms)
}

// Read refreshes the Terraform state with the latest data.
func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config DeploymentModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiClient, err := d.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(summary, "could not create sdk client")

		return
	}

	dep, err := getDeployment(ctx, &config, apiClient)
	if err != nil {
		resp.Diagnostics.AddError(summary, err.Error())

		return
	}

	state := deploymentAsState(dep)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
