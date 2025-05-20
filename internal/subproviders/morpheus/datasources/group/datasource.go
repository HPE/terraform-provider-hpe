package group

import (
	"context"
	"fmt"
	"net/http"

	"github.com/HewlettPackard/hpe-morpheus-client/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/configure"
)

const subject = "read group resource"

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource = &DataSource{}
)

// NewCoffeesDataSource is a helper function to simplify the provider implementation.
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

	apiClient, _ := d.NewClient(ctx)
	var group client.ListGroups200ResponseAllOfGroupsInner

	// Get by id
	if !data.Id.IsNull() {
		id := data.Id.ValueInt64()

		g, hresp, err := apiClient.GroupsAPI.GetGroups(ctx, id).Execute()
		if err != nil || hresp.StatusCode != http.StatusOK {
			resp.Diagnostics.AddError(
				subject,
				fmt.Sprintf("group %d GET failed: ", id),
			)

			return
		}

		group = g.GetGroup()

		// Get by name
	} else if data.Name.ValueString() != "" {
		name := data.Name.ValueString()

		gs, hresp, err := apiClient.GroupsAPI.ListGroups(ctx).Name(name).Execute()
		if err != nil || hresp.StatusCode != http.StatusOK {
			resp.Diagnostics.AddError(
				subject,
				fmt.Sprintf("group %s GET failed: ", name),
			)

			return
		}

		for _, g := range gs.Groups {
			if g.Name != nil && *g.Name == name {
				group = g
				break
			}
		}

	} else {
		resp.Diagnostics.AddError(
			subject,
			"no valid search terms",
		)

		return
	}

	if group.Id == nil {
		resp.Diagnostics.AddError(
			subject,
			"no group found",
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
