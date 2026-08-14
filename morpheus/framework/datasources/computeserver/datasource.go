// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

// Package computeserver implements a data source for compute_server
package computeserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	providererrors "github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

const (
	summary                     = "read compute server data source"
	ErrorNoValidSearchTerms     = `no valid search terms - an id or name is required`
	ErrorNoComputeServerFound   = `no compute server found`
	ErrorMultipleComputeServers = `multiple compute servers were returned`
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
	resp.TypeName = req.ProviderTypeName + "_" + "compute_server"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = ComputeServerDataSourceSchema(ctx)
}

func serverAsState(
	s *sdk.GetHost200ResponseServer,
) ComputeServerModel {
	state := ComputeServerModel{
		Id:             convert.Int64ToType(s.Id),
		Name:           convert.StrToType(s.Name),
		Description:    convert.StrToType(s.Description.Get()),
		Hostname:       convert.StrToType(s.Hostname),
		ExternalId:     convert.StrToType(s.ExternalId.Get()),
		ExternalIp:     convert.StrToType(s.ExternalIp.Get()),
		InternalId:     convert.StrToType(s.InternalId.Get()),
		InternalIp:     convert.StrToType(s.InternalIp.Get()),
		Platform:       convert.StrToType(s.Platform.Get()),
		PowerState:     convert.StrToType(s.PowerState),
		Status:         convert.StrToType(s.Status),
		Uuid:           convert.StrToType(s.Uuid),
		Visibility:     convert.StrToType(s.Visibility),
		AgentInstalled: convert.BoolToType(s.AgentInstalled),
		MaxMemory:      convert.Int64ToType(s.MaxMemory),
		MaxStorage:     convert.Int64ToType(s.MaxStorage),
		CloudId:        convert.Int64ToType(s.ZoneId),
		GroupId:        convert.Int64ToType(s.SiteId),
	}

	// Cloud name from Zone nested object
	if s.Zone != nil && s.Zone.Name != nil {
		state.CloudName = convert.StrToType(s.Zone.Name)
	} else {
		state.CloudName = types.StringNull()
	}

	// ComputeServerType nested object
	if s.ComputeServerType != nil {
		state.ComputeServerTypeId = convert.Int64ToType(s.ComputeServerType.Id)
		state.ComputeServerTypeCode = convert.StrToType(s.ComputeServerType.Code)
		state.ComputeServerTypeName = convert.StrToType(s.ComputeServerType.Name)
		state.ComputeServerTypeManaged = convert.BoolToType(s.ComputeServerType.Managed)
	} else {
		state.ComputeServerTypeId = types.Int64Null()
		state.ComputeServerTypeCode = types.StringNull()
		state.ComputeServerTypeName = types.StringNull()
		state.ComputeServerTypeManaged = types.BoolNull()
	}

	// Plan nested object
	if s.Plan != nil {
		state.PlanId = convert.Int64ToType(s.Plan.Id.Get())
		state.PlanCode = convert.StrToType(s.Plan.Code.Get())
		state.PlanName = convert.StrToType(s.Plan.Name.Get())
	} else {
		state.PlanId = types.Int64Null()
		state.PlanCode = types.StringNull()
		state.PlanName = types.StringNull()
	}

	// Instance
	if s.Instance != nil && s.Instance.Id != nil {
		state.InstanceId = convert.Int64ToType(s.Instance.Id)
	} else {
		state.InstanceId = types.Int64Null()
	}

	// Parent host: the hypervisor this compute server runs on. Null for a host
	// itself, which has no parent.
	if s.ParentServer != nil {
		state.ParentHostId = convert.Int64ToType(s.ParentServer.Id)
		state.ParentHostName = convert.StrToType(s.ParentServer.Name)
	} else {
		state.ParentHostId = types.Int64Null()
		state.ParentHostName = types.StringNull()
	}

	// Labels
	if len(s.Labels) > 0 {
		vals := make([]attr.Value, 0, len(s.Labels))
		for _, l := range s.Labels {
			vals = append(vals, types.StringValue(l))
		}
		state.Labels, _ = types.SetValue(types.StringType, vals)
	} else {
		state.Labels = types.SetNull(types.StringType)
	}

	return state
}

func getComputeServerByID(
	ctx context.Context,
	id int64,
	apiClient *sdk.APIClient,
) (*sdk.GetHost200ResponseServer, error) {
	idParam := sdk.GetHostIdParameter{Int64: &id}
	r, hresp, err := apiClient.HostsAPI.GetHost(ctx, idParam).Execute()
	if r == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"GET failed for compute server %d: %s",
			id, providererrors.ErrMsg(err, hresp),
		)
	}

	return r.Server, nil
}

func getComputeServerByName(
	ctx context.Context,
	name string,
	apiClient *sdk.APIClient,
) (*sdk.GetHost200ResponseServer, error) {
	rs, hresp, err := apiClient.HostsAPI.ListHosts(ctx).
		Name(name).Max(1000).Execute()
	if rs == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"GET failed for compute servers with name %s: %s",
			name, providererrors.ErrMsg(err, hresp),
		)
	}

	var matchedIDs []int64

	for i := range rs.Servers {
		if rs.Servers[i].Name == nil || *rs.Servers[i].Name != name {
			continue
		}
		if rs.Servers[i].Id == nil {
			continue
		}

		matchedIDs = append(matchedIDs, *rs.Servers[i].Id)
	}

	if len(matchedIDs) == 0 {
		return nil, errors.New(ErrorNoComputeServerFound)
	} else if len(matchedIDs) > 1 {
		return nil, errors.New(ErrorMultipleComputeServers)
	}

	return getComputeServerByID(ctx, matchedIDs[0], apiClient)
}

func getComputeServer(
	ctx context.Context,
	config *ComputeServerModel,
	apiClient *sdk.APIClient,
) (*sdk.GetHost200ResponseServer, error) {
	if !config.Id.IsNull() {
		return getComputeServerByID(ctx, config.Id.ValueInt64(), apiClient)
	} else if !config.Name.IsNull() {
		return getComputeServerByName(ctx, config.Name.ValueString(), apiClient)
	}

	return nil, errors.New(ErrorNoValidSearchTerms)
}

// Read refreshes the Terraform state with the latest data.
func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config ComputeServerModel

	// Read config
	diags := req.Config.Get(ctx, &config)
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

	s, err := getComputeServer(ctx, &config, apiClient)
	if err != nil {
		resp.Diagnostics.AddError(
			summary,
			err.Error(),
		)

		return
	}

	state := serverAsState(s)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
