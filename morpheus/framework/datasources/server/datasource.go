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
	resp.TypeName = req.ProviderTypeName + "_morpheus_server"
}

func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = ServerDataSourceSchema(ctx)
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
	response, hresp, err := client.HostsAPI.GetHost(ctx, sdk.Int64AsGetHostIdParameter(&id)).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server %d GET failed: %s", id, errfmt.ErrMsg(err, hresp))
	}

	server, ok := response.GetServerOk()
	if !ok {
		return nil, fmt.Errorf("server %d is nil", id)
	}

	state := &ServerModel{}
	state.Id = types.Int64Value(id)
	state.Name = convert.StrToType(server.Name)
	state.Hostname = convert.StrToType(server.Hostname)
	state.ExternalIp = convert.StrToType(server.ExternalIp.Get())
	state.InternalIp = convert.StrToType(server.InternalIp.Get())
	state.Status = convert.StrToType(server.Status)
	state.PowerState = convert.StrToType(server.PowerState)
	state.Platform = convert.StrToType(server.Platform.Get())

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
	for _, server := range servers.GetServers() {
		if serverName, ok := server.GetNameOk(); ok && *serverName == name {
			matchingServers = append(matchingServers, server)
		}
	}

	if len(matchingServers) == 0 {
		return nil, fmt.Errorf("server %s not found", name)
	}

	if len(matchingServers) > 1 {
		var serverIDs []string
		for _, s := range matchingServers {
			if id, ok := s.GetIdOk(); ok {
				serverIDs = append(serverIDs, fmt.Sprintf("%d", *id))
			}
		}

		return nil, fmt.Errorf(
			"multiple servers found with name %s. Server IDs: %s. "+
				"Please specify an ID instead",
			name,
			strings.Join(serverIDs, ", "),
		)
	}

	id, ok := matchingServers[0].GetIdOk()
	if !ok {
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
