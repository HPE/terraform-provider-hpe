// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

func timeToType(t *time.Time) types.String {
	if t == nil {
		return types.StringNull()
	}

	return types.StringValue(t.Format(time.RFC3339))
}

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
	state.OsDevice = convert.StrToType(server.OsDevice)
	state.MaxCores = convert.Int64ToType(server.MaxCores)
	state.MaxMemory = convert.Int64ToType(server.MaxMemory)
	state.MaxStorage = convert.Int64ToType(server.MaxStorage)
	state.GroupId = convert.Int64ToType(server.SiteId)
	state.AccountId = convert.Int64ToType(server.AccountId)
	state.Visibility = convert.StrToType(server.Visibility)
	state.ExternalName = convert.StrToType(server.ExternalName)
	state.AgentInstalled = convert.BoolToType(server.AgentInstalled)
	state.Enabled = convert.BoolToType(server.Enabled)
	state.ManageInternalFirewall = convert.BoolToType(server.ManageInternalFirewall)
	state.EnableLogs = convert.BoolToType(server.EnableLogs)
	state.SshPort = convert.Int64ToType(server.SshPort)
	state.DateCreated = timeToType(server.DateCreated)
	state.LastUpdated = timeToType(server.LastUpdated)

	// NullableString fields
	state.ExternalIp = convert.StrToType(server.ExternalIp.Get())
	state.InternalIp = convert.StrToType(server.InternalIp.Get())
	state.Platform = convert.StrToType(server.Platform.Get())
	state.PlatformVersion = convert.StrToType(server.PlatformVersion.Get())
	state.Description = convert.StrToType(server.Description.Get())
	state.ExternalId = convert.StrToType(server.ExternalId.Get())
	state.InternalId = convert.StrToType(server.InternalId.Get())
	state.SshHost = convert.StrToType(server.SshHost.Get())
	state.StatusMessage = convert.StrToType(server.StatusMessage.Get())
	state.ErrorMessage = convert.StrToType(server.ErrorMessage.Get())
	state.AgentVersion = convert.StrToType(server.AgentVersion.Get())

	// NullableInt64 fields
	state.CoresPerSocket = convert.Int64ToType(server.CoresPerSocket.Get())
	state.ResourcePoolId = convert.Int64ToType(server.ResourcePoolId.Get())
	state.FolderId = convert.Int64ToType(server.FolderId.Get())
	state.MaxCpu = convert.Int64ToType(server.MaxCpu.Get())

	// NullableBool fields
	state.TagCompliant = convert.BoolToType(server.TagCompliant.Get())

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
