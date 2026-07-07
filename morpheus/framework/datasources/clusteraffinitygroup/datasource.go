// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

// Package clusteraffinitygroup implements a data source for cluster_affinity_group
package clusteraffinitygroup

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	providererrors "github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

const (
	summary                            = "read cluster affinity group data source"
	ErrorNoValidSearchTerms            = `no valid search terms - an id or name is required`
	ErrorNoClusterAffinityGroupFound   = `no cluster affinity group found`
	ErrorMultipleClusterAffinityGroups = `multiple cluster affinity groups were returned`
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
	resp.TypeName = req.ProviderTypeName + "_" + "cluster_affinity_group"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = ClusterAffinityGroupDataSourceSchema(ctx)
}

func affinityGroupAsState(
	ag *sdk.GetClusterAffinityGroup200ResponseAffinityGroup,
	clusterID int64,
) ClusterAffinityGroupModel {
	return ClusterAffinityGroupModel{
		Id:         convert.Int64ToType(ag.Id),
		ClusterId:  types.Int64Value(clusterID),
		Name:       convert.StrToType(ag.Name),
		Active:     convert.BoolToType(ag.Active),
		Visibility: convert.StrToType(ag.Visibility),
	}
}

func getAffinityGroupByID(
	ctx context.Context,
	id int64,
	clusterID int64,
	apiClient *sdk.APIClient,
) (*sdk.GetClusterAffinityGroup200ResponseAffinityGroup, error) {
	// GetClusterAffinityGroup(ctx, clusterId, id)
	r, hresp, err := apiClient.ClustersAPI.GetClusterAffinityGroup(
		ctx, clusterID, id,
	).Execute()
	if r == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"GET failed for cluster affinity group %d: %s",
			id, providererrors.ErrMsg(err, hresp),
		)
	}

	return r.AffinityGroup, nil
}

func getAffinityGroupByName(
	ctx context.Context,
	name string,
	clusterID int64,
	apiClient *sdk.APIClient,
) (*sdk.GetClusterAffinityGroup200ResponseAffinityGroup, error) {
	// ListClusterAffinityGroups(ctx, clusterId)
	rs, hresp, err := apiClient.ClustersAPI.ListClusterAffinityGroups(
		ctx, clusterID,
	).Execute()
	if rs == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"GET failed for cluster affinity groups with name %s: %s",
			name, providererrors.ErrMsg(err, hresp),
		)
	}

	items := rs.AffinityGroups
	if len(items) == 0 {
		return nil, errors.New(ErrorNoClusterAffinityGroupFound)
	}

	var matchedIDs []int64

	for i := range items {
		if items[i].Name == nil || *items[i].Name != name {
			continue
		}
		if items[i].Id == nil {
			continue
		}

		matchedIDs = append(matchedIDs, *items[i].Id)
	}

	if len(matchedIDs) == 0 {
		return nil, errors.New(ErrorNoClusterAffinityGroupFound)
	} else if len(matchedIDs) > 1 {
		return nil, errors.New(ErrorMultipleClusterAffinityGroups)
	}

	return getAffinityGroupByID(ctx, matchedIDs[0], clusterID, apiClient)
}

func getAffinityGroup(
	ctx context.Context,
	config *ClusterAffinityGroupModel,
	apiClient *sdk.APIClient,
) (*sdk.GetClusterAffinityGroup200ResponseAffinityGroup, error) {
	clusterID := config.ClusterId.ValueInt64()

	if !config.Id.IsNull() {
		return getAffinityGroupByID(ctx, config.Id.ValueInt64(), clusterID, apiClient)
	} else if !config.Name.IsNull() {
		return getAffinityGroupByName(ctx, config.Name.ValueString(), clusterID, apiClient)
	}

	return nil, errors.New(ErrorNoValidSearchTerms)
}

// Read refreshes the Terraform state with the latest data.
func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config ClusterAffinityGroupModel

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

	ag, err := getAffinityGroup(ctx, &config, apiClient)
	if err != nil {
		resp.Diagnostics.AddError(
			summary,
			err.Error(),
		)

		return
	}

	clusterID := config.ClusterId.ValueInt64()
	state := affinityGroupAsState(ag, clusterID)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
