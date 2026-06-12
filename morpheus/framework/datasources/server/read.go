// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package server

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

func getServer(
	ctx context.Context,
	config ServerModel,
	client *sdk.APIClient,
) (*ServerModel, error) {
	if !config.Id.IsNull() {
		return getServerByID(ctx, config.Id.ValueInt64(), client)
	}

	if !config.Name.IsNull() {
		return getServerByName(ctx, config.Name.ValueString(), client)
	}

	return nil, fmt.Errorf("either id or name must be specified")
}

func getServerByID(
	ctx context.Context,
	id int64,
	client *sdk.APIClient,
) (*ServerModel, error) {
	getId := sdk.GetHostIdParameter{Int64: &id}

	response, hresp, err := client.HostsAPI.GetHost(ctx, getId).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server %d GET failed: %s", id, errfmt.ErrMsg(err, hresp))
	}

	if response.Server == nil {
		return nil, fmt.Errorf("server %d is nil", id)
	}

	server := response.Server
	state := &ServerModel{}

	state.Id = types.Int64Value(id)
	state.Name = convert.StrToType(server.Name)
	state.Hostname = convert.StrToType(server.Hostname)
	state.Status = convert.StrToType(server.Status)
	state.PowerState = convert.StrToType(server.PowerState)
	state.Uuid = convert.StrToType(server.Uuid)
	state.OsType = convert.StrToType(server.OsType)
	state.MaxCores = convert.Int64ToType(server.MaxCores)
	state.MaxMemory = convert.Int64ToType(server.MaxMemory)
	state.MaxStorage = convert.Int64ToType(server.MaxStorage)
	state.GroupId = convert.Int64ToType(server.SiteId)

	if server.ExternalIp.IsSet() {
		state.ExternalIp = types.StringValue(*server.ExternalIp.Get())
	} else {
		state.ExternalIp = types.StringNull()
	}

	if server.InternalIp.IsSet() {
		state.InternalIp = types.StringValue(*server.InternalIp.Get())
	} else {
		state.InternalIp = types.StringNull()
	}

	if server.Platform.IsSet() {
		state.Platform = types.StringValue(*server.Platform.Get())
	} else {
		state.Platform = types.StringNull()
	}

	if server.PlatformVersion.IsSet() {
		state.PlatformVersion = types.StringValue(*server.PlatformVersion.Get())
	} else {
		state.PlatformVersion = types.StringNull()
	}

	if server.Description.IsSet() {
		state.Description = types.StringValue(*server.Description.Get())
	} else {
		state.Description = types.StringNull()
	}

	if server.ExternalId.IsSet() {
		state.ExternalId = types.StringValue(*server.ExternalId.Get())
	} else {
		state.ExternalId = types.StringNull()
	}

	if server.CoresPerSocket.IsSet() {
		state.CoresPerSocket = convert.Int64ToType(server.CoresPerSocket.Get())
	} else {
		state.CoresPerSocket = types.Int64Null()
	}

	if server.Zone != nil {
		state.CloudId = convert.Int64ToType(server.Zone.Id)
	} else {
		state.CloudId = types.Int64Null()
	}

	return state, nil
}

func getServerByName(
	ctx context.Context,
	name string,
	client *sdk.APIClient,
) (*ServerModel, error) {
	servers, hresp, err := client.HostsAPI.ListHosts(ctx).Name(name).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server %s list failed: %s", name, errfmt.ErrMsg(err, hresp))
	}

	var matchingServers []sdk.ListHosts200ResponseAllOfServersInner

	for _, server := range servers.Servers {
		if server.Name != nil && *server.Name == name {
			matchingServers = append(matchingServers, server)
		}
	}

	if len(matchingServers) == 0 {
		return nil, fmt.Errorf("server %s not found", name)
	}

	if len(matchingServers) > 1 {
		var serverIDs []string

		for _, s := range matchingServers {
			if s.Id != nil {
				serverIDs = append(serverIDs, fmt.Sprintf("%d", *s.Id))
			}
		}

		return nil, fmt.Errorf(
			"multiple servers found with name %s. Server IDs: %s. "+
				"Please specify an ID instead",
			name,
			strings.Join(serverIDs, ", "),
		)
	}

	id := matchingServers[0].Id
	if id == nil {
		return nil, fmt.Errorf("server %s has missing ID", name)
	}

	return getServerByID(ctx, *id, client)
}

func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config ServerModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := d.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"read server data source",
			fmt.Sprintf("failed to create client: %s", err.Error()),
		)

		return
	}

	state, err := getServer(ctx, config, client)
	if err != nil {
		resp.Diagnostics.AddError(
			"read server data source",
			err.Error(),
		)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
