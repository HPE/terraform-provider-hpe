package environment

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/convert"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/datasources/environment/consts"
)

const summary = "read environment data source"

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource = &DataSource{}
)

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
	resp.TypeName = req.ProviderTypeName + "_morpheus_environment"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = EnvironmentDataSourceSchema(ctx)
}

func getEnvironmentById(
	ctx context.Context,
	id int64,
	apiClient *sdk.APIClient,
) (*sdk.ListEnvironments200ResponseAllOfEnvironmentsInner, error) {
	e, hresp, err := apiClient.EnvironmentsAPI.GetEnvironments(ctx, id).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET failed for environment %d", id)
	}

	environment := e.GetEnvironment()

	return &environment, nil
}

func getEnvironmentByName(
	ctx context.Context,
	name string,
	apiClient *sdk.APIClient,
) (*sdk.ListEnvironments200ResponseAllOfEnvironmentsInner, error) {
	es, hresp, err := apiClient.EnvironmentsAPI.ListEnvironments(ctx).Name(name).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET failed for environment %s", name)
	}

	var environments []sdk.ListEnvironments200ResponseAllOfEnvironmentsInner

	for _, e := range es.Environments {
		if e.GetName() == name {
			environments = append(environments, e)
		}
	}

	if len(environments) == 1 {
		return &environments[0], nil
	} else if len(environments) > 1 {
		return nil, errors.New(consts.ErrorMultipleEnvironments)
	}

	return nil, errors.New(consts.ErrorNoEnvironmentFound)
}

func getEnvironment(
	ctx context.Context,
	data EnvironmentModel,
	apiClient *sdk.APIClient,
) (*sdk.ListEnvironments200ResponseAllOfEnvironmentsInner, error) {
	if !data.Id.IsNull() {
		return getEnvironmentById(ctx, data.Id.ValueInt64(), apiClient)
	} else if !data.Name.IsNull() {
		return getEnvironmentByName(ctx, data.Name.ValueString(), apiClient)
	}

	return nil, errors.New(consts.ErrorNoValidSearchTerms)
}

// Read refreshes the Terraform state with the latest data.
func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var data EnvironmentModel

	// Read config
	diags := req.Config.Get(ctx, &data)
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

	environment, err := getEnvironment(ctx, data, apiClient)
	if err != nil {
		resp.Diagnostics.AddError(
			summary,
			err.Error(),
		)

		return
	}

	data.Active = convert.BoolToType(environment.Active)
	data.Code = convert.StrToType(environment.Code)
	data.Description = convert.StrToType(environment.Description)
	data.Id = convert.Int64ToType(environment.Id)
	data.Name = convert.StrToType(environment.Name)
	data.Visibility = convert.StrToType(environment.Visibility)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
