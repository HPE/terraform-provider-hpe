// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package securitygroup

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

var _ datasource.DataSource = &DataSource{}

func NewDataSource() datasource.DataSource {
	return &DataSource{}
}

type DataSource struct {
	configure.DataSourceWithMorpheusConfigure
	datasource.DataSource
}

func (d *DataSource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_morpheus_security_group"
}

func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = SecurityGroupDataSourceSchema(ctx)
}

func getSecurityGroup(
	ctx context.Context,
	config SecurityGroupModel,
	client *sdk.APIClient,
) (*SecurityGroupModel, error) {
	if !config.Id.IsNull() {
		return getSecurityGroupByID(ctx, config.Id.ValueInt64(), client)
	}
	if !config.Name.IsNull() {
		return getSecurityGroupByName(ctx, config.Name.ValueString(), client)
	}

	return nil, fmt.Errorf("either id or name must be specified")
}

func getSecurityGroupByID(
	ctx context.Context,
	id int64,
	client *sdk.APIClient,
) (*SecurityGroupModel, error) {
	response, hresp, err := client.SecurityGroupsAPI.GetSecurityGroups(ctx, id).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("security group %d GET failed: %s", id, errfmt.ErrMsg(err, hresp))
	}

	state := &SecurityGroupModel{}

	sg, ok := response.GetSecurityGroupOk()
	if !ok {
		return nil, fmt.Errorf("security group %d is nil", id)
	}

	state.Id = types.Int64Value(id)
	state.Name = convert.StrToType(sg.Name)
	state.Description = convert.StrToType(sg.Description.Get())
	state.Active = convert.BoolToType(sg.Active)
	state.Visibility = convert.StrToType(sg.Visibility)
	state.ExternalId = convert.StrToType(sg.ExternalId.Get())

	if zone, ok := sg.GetZoneOk(); ok && zone != nil {
		if zoneId, ok := zone.GetIdOk(); ok {
			state.CloudId = types.Int64Value(*zoneId)
		}
	}

	return state, nil
}

func getSecurityGroupByName(
	ctx context.Context,
	name string,
	client *sdk.APIClient,
) (*SecurityGroupModel, error) {
	response, hresp, err := client.SecurityGroupsAPI.ListSecurityGroups(ctx).Name(name).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("security group %s list failed: %s", name, errfmt.ErrMsg(err, hresp))
	}

	var matchingGroups []sdk.ListSecurityGroups200ResponseAllOfSecurityGroupsInner
	for _, sg := range response.GetSecurityGroups() {
		if sgName, ok := sg.GetNameOk(); ok && *sgName == name {
			matchingGroups = append(matchingGroups, sg)
		}
	}

	if len(matchingGroups) == 0 {
		return nil, fmt.Errorf("security group %s not found", name)
	}

	if len(matchingGroups) > 1 {
		var sgIDs []string
		for _, sg := range matchingGroups {
			if id, ok := sg.GetIdOk(); ok {
				sgIDs = append(sgIDs, fmt.Sprintf("%d", *id))
			}
		}

		return nil, fmt.Errorf(
			"multiple security groups found with name %s. Security group IDs: %s. "+
				"Please specify an ID instead",
			name,
			strings.Join(sgIDs, ", "),
		)
	}

	id, ok := matchingGroups[0].GetIdOk()
	if !ok {
		return nil, fmt.Errorf("security group %s has missing ID", name)
	}

	return getSecurityGroupByID(ctx, *id, client)
}

func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config SecurityGroupModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := d.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"read security group data source",
			fmt.Sprintf("failed to create client: %s", err.Error()),
		)

		return
	}

	state, err := getSecurityGroup(ctx, config, client)
	if err != nil {
		resp.Diagnostics.AddError(
			"read security group data source",
			err.Error(),
		)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
