// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package deployment

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	providererrors "github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

const (
	summary                  = "read deployment data source"
	ErrorNoValidSearchTerms  = `no valid search terms - an id or name is required`
	ErrorRunningPreApply     = `Error running pre-apply plan: exit status 1`
	ErrorNoDeploymentFound   = `no deployment found`
	ErrorMultipleDeployments = `multiple deployments were returned`
)

// Ensure the implementation satisfies the expected interfaces.
var _ datasource.DataSource = &DataSource{}

// NewDataSource is a helper function to simplify the provider implementation.
func NewDataSource() datasource.DataSource {
	return &DataSource{}
}

// DataSource is the data source implementation.
type DataSource struct {
	configure.DataSourceWithMorpheusConfigure
	datasource.DataSource
}

// Metadata returns the data source type name.
func (d *DataSource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_" + "deployment"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = DeploymentDataSourceSchema(ctx)
}

func deploymentAsState(
	dep *sdk.GetDeployment200ResponseDeployment,
) DeploymentModel {
	return DeploymentModel{
		Id:          convert.Int64ToType(dep.Id),
		Name:        convert.StrToType(dep.Name),
		Description: convert.StrToType(dep.Description),
	}
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

	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiClient, err := d.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			summary,
			"could not create sdk client",
		)

		return
	}

	dep, err := getDeployment(ctx, &config, apiClient)
	if err != nil {
		resp.Diagnostics.AddError(
			summary,
			err.Error(),
		)

		return
	}

	state := deploymentAsState(dep)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
