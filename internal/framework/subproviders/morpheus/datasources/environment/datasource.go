// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package environment

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/convert"
)

const (
	summary                   = "read environment data source"
	ErrorNoValidSearchTerms   = `no valid search terms - an id or name is required`
	ErrorRunningPreApply      = `Error running pre-apply plan: exit status 1`
	ErrorNoEnvironmentFound   = `no environment found`
	ErrorMultipleEnvironments = `multiple environments were returned`
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

func getEnvironmentByID(
	ctx context.Context,
	id int64,
	data *EnvironmentModel,
	apiClient *sdk.APIClient,
) error {
	e, hresp, err := apiClient.EnvironmentsAPI.GetEnvironments(ctx, id).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		return fmt.Errorf("GET failed for environment %d", id)
	}

	environment := e.GetEnvironment()

	data.Active = convert.BoolToType(environment.Active)
	data.Code = convert.StrToType(environment.Code)
	data.Description = convert.StrToType(environment.Description)
	data.Id = convert.Int64ToType(environment.Id)
	data.Name = convert.StrToType(environment.Name)
	data.Visibility = convert.StrToType(environment.Visibility)

	return nil
}

func getEnvironmentByName(
	ctx context.Context,
	name string,
	data *EnvironmentModel,
	apiClient *sdk.APIClient,
) error {
	es, hresp, err := apiClient.EnvironmentsAPI.ListEnvironments(ctx).Name(name).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		return fmt.Errorf("GET failed for environment %s", name)
	}

	environments := sdk.NewListEnvironments200Response().Environments

	for _, e := range es.Environments {
		if e.GetName() == name {
			environments = append(environments, e)
		}
	}

	if len(environments) == 0 {
		return errors.New(ErrorNoEnvironmentFound)
	} else if len(environments) > 1 {
		return errors.New(ErrorMultipleEnvironments)
	}

	environment := environments[0]

	data.Active = convert.BoolToType(environment.Active)
	data.Code = convert.StrToType(environment.Code)
	data.Description = convert.StrToType(environment.Description)
	data.Id = convert.Int64ToType(environment.Id)
	data.Name = convert.StrToType(environment.Name)
	data.Visibility = convert.StrToType(environment.Visibility)

	return nil
}

func getEnvironment(
	ctx context.Context,
	data *EnvironmentModel,
	apiClient *sdk.APIClient,
) error {
	if !data.Id.IsNull() {
		return getEnvironmentByID(ctx, data.Id.ValueInt64(), data, apiClient)
	} else if !data.Name.IsNull() {
		return getEnvironmentByName(ctx, data.Name.ValueString(), data, apiClient)
	}

	return errors.New(ErrorNoValidSearchTerms)
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

	if err := getEnvironment(ctx, &data, apiClient); err != nil {
		resp.Diagnostics.AddError(
			summary,
			err.Error(),
		)

		return
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
