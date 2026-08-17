// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkservergroup

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	providererrors "github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

const (
	summary                          = "read network server group data source"
	ErrorNoValidSearchTerms          = `name is required`
	ErrorNoNetworkServerGroupFound   = `no network server group found`
	ErrorMultipleNetworkServerGroups = `multiple network server groups matched`
	ErrorNoNSXTServerFound           = `no NSX-T network server found`

	nsxtTypeCode = "nsx-t"
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
	resp.TypeName = req.ProviderTypeName + "_" + "network_server_group"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = NetworkServerGroupDataSourceSchema(ctx)
}

// discoverNSXTServerID lists all network servers and returns the ID of the
// first one whose type.code == "nsx-t".
func discoverNSXTServerID(
	ctx context.Context,
	apiClient *sdk.APIClient,
) (int64, error) {
	rs, hresp, err := apiClient.NetworksAPI.ListNetworkServers(ctx).Max(10000).Execute()
	if rs == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("failed to list network servers: %s", providererrors.ErrMsg(err, hresp))
	}

	for i := range rs.NetworkServers {
		ns := &rs.NetworkServers[i]
		if ns.Type != nil && ns.Type.Code != nil && *ns.Type.Code == nsxtTypeCode {
			if ns.Id != nil {
				return *ns.Id, nil
			}
		}
	}

	return 0, errors.New(ErrorNoNSXTServerFound)
}

// findGroupByName lists groups under the given serverId and returns the first
// exact name match.
func findGroupByName(
	ctx context.Context,
	apiClient *sdk.APIClient,
	serverID int64,
	name string,
) (*sdk.ListNetworkServerGroups200ResponseAllOfGroupsInner, error) {
	rs, hresp, err := apiClient.NetworksAPI.ListNetworkServerGroups(ctx, serverID).
		Max(10000).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to list network server groups: %s", providererrors.ErrMsg(err, hresp))
	}

	if rs == nil {
		return nil, errors.New(ErrorNoNetworkServerGroupFound)
	}

	// The SDK types the slice, but round-trip through JSON for robustness
	// (same pattern as networkserver).
	raw, marshalErr := json.Marshal(rs.Groups)
	if marshalErr != nil {
		return nil, fmt.Errorf("error marshaling groups: %w", marshalErr)
	}

	var groups []sdk.ListNetworkServerGroups200ResponseAllOfGroupsInner

	if unmarshalErr := json.Unmarshal(raw, &groups); unmarshalErr != nil {
		return nil, fmt.Errorf("error decoding groups: %w", unmarshalErr)
	}

	var matched *sdk.ListNetworkServerGroups200ResponseAllOfGroupsInner

	var matchCount int

	for i := range groups {
		if groups[i].Name != nil && *groups[i].Name == name {
			matched = &groups[i]
			matchCount++
		}
	}

	if matchCount == 0 {
		return nil, errors.New(ErrorNoNetworkServerGroupFound)
	} else if matchCount > 1 {
		return nil, errors.New(ErrorMultipleNetworkServerGroups)
	}

	return matched, nil
}

func mapTag(t sdk.ListNetworkServerGroups200ResponseAllOfGroupsInnerTagsInner) TagsValue {
	return TagsValue{
		Id:    convert.Int64ToType(t.Id),
		Name:  convert.StrToType(t.Name),
		Value: convert.StrToType(t.Value),
		state: attr.ValueStateKnown,
	}
}

func mapMember(m sdk.ListNetworkServerGroups200ResponseAllOfGroupsInnerMembersInner) MembersValue {
	return MembersValue{
		Id:               convert.Int64ToType(m.Id),
		Category:         convert.StrToType(m.Category),
		MembersType:      convert.StrToType(m.Type),
		MemberName:       convert.StrToType(m.MemberName.Get()),
		MemberType:       convert.StrToType(m.MemberType.Get()),
		MemberValue:      convert.StrToType(m.MemberValue.Get()),
		MemberExpression: convert.StrToType(m.MemberExpression.Get()),
		DisplayOrder:     convert.Int64ToType(m.DisplayOrder),
		InternalId:       convert.StrToType(m.InternalId.Get()),
		ExternalId:       convert.StrToType(m.ExternalId),
		state:            attr.ValueStateKnown,
	}
}

func idOnlyAccount(a *sdk.ListNetworkServerGroups200ResponseAllOfGroupsInnerAccount) AccountValue {
	if a == nil {
		return NewAccountValueNull()
	}

	return AccountValue{
		Id:    convert.Int64ToType(a.Id),
		state: attr.ValueStateKnown,
	}
}

func idOnlyOwner(o *sdk.ListNetworkServerGroups200ResponseAllOfGroupsInnerOwner) OwnerValue {
	if o == nil {
		return NewOwnerValueNull()
	}

	return OwnerValue{
		Id:    convert.Int64ToType(o.Id),
		state: attr.ValueStateKnown,
	}
}

func groupAsState(
	ctx context.Context,
	g *sdk.ListNetworkServerGroups200ResponseAllOfGroupsInner,
	serverID int64,
) (NetworkServerGroupModel, error) {
	tagsSet, diags := convert.ToSetType(ctx, g.Tags, mapTag)
	if diags.HasError() {
		return NetworkServerGroupModel{}, fmt.Errorf("error mapping tags: %s", diags.Errors())
	}

	membersSet, diags := convert.ToSetType(ctx, g.Members, mapMember)
	if diags.HasError() {
		return NetworkServerGroupModel{}, fmt.Errorf("error mapping members: %s", diags.Errors())
	}

	state := NetworkServerGroupModel{
		Id:              convert.Int64ToType(g.Id),
		Name:            convert.StrToType(g.Name),
		Description:     convert.StrToType(g.Description.Get()),
		ExternalId:      convert.StrToType(g.ExternalId),
		InternalId:      convert.StrToType(g.InternalId),
		Visibility:      convert.StrToType(g.Visibility),
		NetworkServerId: convert.Int64ToType(&serverID),
		Account:         idOnlyAccount(g.Account),
		Owner:           idOnlyOwner(g.Owner),
		Tags:            tagsSet,
		Members:         membersSet,
	}

	return state, nil
}

// Read refreshes the Terraform state with the latest data.
func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config NetworkServerGroupModel

	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	if config.Name.IsNull() || config.Name.ValueString() == "" {
		resp.Diagnostics.AddError(summary, ErrorNoValidSearchTerms)

		return
	}

	apiClient, err := d.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(summary, "could not create sdk client")

		return
	}

	// Determine network server ID: explicit or NSX-T fallback.
	var serverID int64

	if !config.NetworkServerId.IsNull() {
		serverID = config.NetworkServerId.ValueInt64()
	} else {
		serverID, err = discoverNSXTServerID(ctx, apiClient)
		if err != nil {
			resp.Diagnostics.AddError(summary, err.Error())

			return
		}
	}

	group, err := findGroupByName(ctx, apiClient, serverID, config.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(summary, err.Error())

		return
	}

	state, err := groupAsState(ctx, group, serverID)
	if err != nil {
		resp.Diagnostics.AddError(summary, err.Error())

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
