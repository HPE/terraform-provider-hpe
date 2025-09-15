package user

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/convert"
)

const summary = "read user data source"
const (
	ErrorNoUserFound      = "no user found"
	ErrorNoValidUserTerms = "no valid search terms - id is required"
)

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
	resp.TypeName = req.ProviderTypeName + "_morpheus_user"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = UserDataSourceSchema(ctx)
}

func getUserByID(
	ctx context.Context,
	id int64,
	apiClient *sdk.APIClient,
) (*sdk.AddUserTenant200ResponseAllOfUser, error) {
	u, hresp, err := apiClient.UsersAPI.GetUser(ctx, id).Execute()
	if u == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET failed for user %d", id)
	}
	user, ok := u.GetUserOk()
	if !ok {
		return nil, errors.New(ErrorNoUserFound)
	}
	return user, nil
}

func getUser(
	ctx context.Context,
	data UserModel,
	apiClient *sdk.APIClient,
) (*sdk.AddUserTenant200ResponseAllOfUser, error) {
	if !data.Id.IsNull() {
		return getUserByID(ctx, data.Id.ValueInt64(), apiClient)
	}
	return nil, errors.New(ErrorNoValidUserTerms)
}

// Read refreshes the Terraform state with the latest data.
func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var data UserModel

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

	user, err := getUser(ctx, data, apiClient)
	if err != nil {
		resp.Diagnostics.AddError(
			summary,
			err.Error(),
		)
		return
	}

	// Map API response to state
	data.Id = convert.Int64ToType(user.Id)
	data.Username = convert.StrToType(user.Username)
	data.Email = convert.StrToType(user.Email)
	data.Enabled = convert.BoolToType(user.Enabled)
	data.FirstName = convert.StrToType(user.FirstName)
	data.LastName = convert.StrToType(user.LastName)
	// Add other fields as needed... TODO: find out what other fields are needed

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
