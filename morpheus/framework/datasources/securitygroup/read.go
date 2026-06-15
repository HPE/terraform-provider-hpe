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

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

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

	if response.SecurityGroup == nil {
		return nil, fmt.Errorf("security group %d is nil", id)
	}

	sg := response.SecurityGroup
	state := &SecurityGroupModel{}

	state.Id = types.Int64Value(id)
	state.Name = convert.StrToType(sg.Name)
	state.Active = convert.BoolToType(sg.Active)
	state.Visibility = convert.StrToType(sg.Visibility)
	state.SyncSource = convert.StrToType(sg.SyncSource)
	state.TenantId = convert.Int64ToType(sg.AccountId)

	// NullableString fields
	state.Description = convert.StrToType(sg.Description.Get())
	state.ExternalId = convert.StrToType(sg.ExternalId.Get())
	state.GroupSource = convert.StrToType(sg.GroupSource.Get())
	state.Enabled = convert.StrToType(sg.Enabled.Get())

	if sg.Zone != nil && sg.Zone.Id != nil {
		state.CloudId = types.Int64Value(*sg.Zone.Id)
	} else {
		state.CloudId = types.Int64Null()
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

	for _, sg := range response.SecurityGroups {
		if sg.Name != nil && *sg.Name == name {
			matchingGroups = append(matchingGroups, sg)
		}
	}

	if len(matchingGroups) == 0 {
		return nil, fmt.Errorf("security group %s not found", name)
	}

	if len(matchingGroups) > 1 {
		var sgIDs []string

		for _, sg := range matchingGroups {
			if sg.Id != nil {
				sgIDs = append(sgIDs, fmt.Sprintf("%d", *sg.Id))
			}
		}

		return nil, fmt.Errorf(
			"multiple security groups found with name %s. Security group IDs: %s. "+
				"Please specify an ID instead",
			name,
			strings.Join(sgIDs, ", "),
		)
	}

	id := matchingGroups[0].Id
	if id == nil {
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
