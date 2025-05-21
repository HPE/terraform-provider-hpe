package group

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/datasources/group/consts"
)

const summary = "read group data source"

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
	resp.TypeName = req.ProviderTypeName + "_morpheus_group"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = GroupDataSourceSchema(ctx)
}

func strToType(s *string) types.String {
	if s == nil {
		return types.StringNull()
	}

	return types.StringValue(*s)
}

func int64ToType(i *int64) types.Int64 {
	if i == nil {
		return types.Int64Null()
	}

	return types.Int64Value(*i)
}

func doRead(
	ctx context.Context,
	data GroupModel,
	apiClient *sdk.APIClient,
) (*sdk.ListGroups200ResponseAllOfGroupsInner, error) {
	// Get by id
	if !data.Id.IsNull() {
		id := data.Id.ValueInt64()

		g, hresp, err := apiClient.GroupsAPI.GetGroups(ctx, id).Execute()
		if err != nil || hresp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("GET failed for group %d", id)
		}

		group := g.GetGroup()

		return &group, nil

		// Get by name
	} else if !data.Name.IsNull() {
		name := data.Name.ValueString()

		gs, hresp, err := apiClient.GroupsAPI.ListGroups(ctx).Name(name).Execute()
		if err != nil || hresp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("GET failed for group %s", name)
		}

		var groups []sdk.ListGroups200ResponseAllOfGroupsInner

		for _, g := range gs.Groups {
			if g.GetName() == name {
				groups = append(groups, g)
			}
		}

		if len(groups) == 1 {
			return &groups[0], nil
		} else if len(groups) > 1 {
			return nil, errors.New(consts.ErrorMultipleGroups)
		}

	} else {
		return nil, errors.New(consts.ErrorNoValidSearchTerms)
	}

	return nil, errors.New(consts.ErrorNoGroupFound)
}

// Read refreshes the Terraform state with the latest data.
func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var data GroupModel

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

	group, err := doRead(ctx, data, apiClient)
	if err != nil {
		resp.Diagnostics.AddError(
			summary,
			err.Error(),
		)

		return
	}

	data.Id = int64ToType(group.Id)
	data.Name = strToType(group.Name)
	data.Code = strToType(group.Code)
	data.Location = strToType(group.Location)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
