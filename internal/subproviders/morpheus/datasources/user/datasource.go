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
	ErrorNoValidUserTerms = "no valid search terms - an id or username is required"
	ErrorMultipleUsers    = "multiple users were returned"
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

func getUserByUsername(
	ctx context.Context,
	data *UserModel,
	apiClient *sdk.APIClient,
) error {
	us, hresp, err := apiClient.UsersAPI.ListUsers(ctx).Username(data.Username.ValueString()).Execute()
	if us == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return fmt.Errorf("GET failed for user with username %s", data.Username.ValueString())
	}

	users := sdk.NewListUsers200Response().Users

	for _, u := range us.GetUsers() {
		if u.GetUsername() == data.Username.ValueString() {
			users = append(users, u)
		}
	}

	if len(users) > 1 {
		return errors.New(ErrorMultipleUsers)
	} else if len(users) == 0 {
		return errors.New(ErrorNoUserFound)
	}

	user := users[0]
	// Map API response to state
	data.Id = convert.Int64ToType(user.Id)
	data.Username = convert.StrToType(user.Username)
	data.Email = convert.StrToType(user.Email)
	data.Enabled = convert.BoolToType(user.Enabled)
	data.FirstName = convert.StrToType(user.FirstName)
	data.LastName = convert.StrToType(user.LastName)
	// Add other fields as needed... TODO: find out what other fields are needed

	return nil
}

func getUserByID(
	ctx context.Context,
	id int64,
	data *UserModel,
	apiClient *sdk.APIClient,
) error {
	u, hresp, err := apiClient.UsersAPI.GetUser(ctx, id).Execute()
	if u == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return fmt.Errorf("GET failed for user %d", id)
	}
	user, ok := u.GetUserOk()
	if !ok {
		return errors.New(ErrorNoUserFound)
	}
	// Map API response to state
	data.Id = convert.Int64ToType(user.Id)
	data.Username = convert.StrToType(user.Username)
	data.Email = convert.StrToType(user.Email)
	data.Enabled = convert.BoolToType(user.Enabled)
	data.FirstName = convert.StrToType(user.FirstName)
	data.LastName = convert.StrToType(user.LastName)
	// Add other fields as needed... TODO: find out what other fields are needed

	return nil
}

func getUser(
	ctx context.Context,
	data *UserModel,
	apiClient *sdk.APIClient,
) error {
	if !data.Id.IsNull() {
		return getUserByID(ctx, data.Id.ValueInt64(), data, apiClient)
	} else if !data.Username.IsNull() {
		return getUserByUsername(ctx, data, apiClient)
	}

	return errors.New(ErrorNoValidUserTerms)
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

	err = getUser(ctx, &data, apiClient)
	if err != nil {
		resp.Diagnostics.AddError(
			summary,
			err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
