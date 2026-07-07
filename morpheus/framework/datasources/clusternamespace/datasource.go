// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

// Package clusternamespace implements a data source for cluster_namespace
package clusternamespace

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
	summary                        = "read cluster namespace data source"
	ErrorNoValidSearchTerms        = `no valid search terms - an id or name is required`
	ErrorNoClusterNamespaceFound   = `no cluster namespace found`
	ErrorMultipleClusterNamespaces = `multiple cluster namespaces were returned`
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
	resp.TypeName = req.ProviderTypeName + "_" + "cluster_namespace"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = ClusterNamespaceDataSourceSchema(ctx)
}

func namespaceAsState(
	ns *sdk.GetClusterNamespace200ResponseNamespace,
	clusterID int64,
) ClusterNamespaceModel {
	return ClusterNamespaceModel{
		Id:          convert.Int64ToType(ns.Id),
		ClusterId:   types.Int64Value(clusterID),
		Name:        convert.StrToType(ns.Name),
		Description: convert.StrToType(ns.Description),
		ExternalId:  convert.StrToType(ns.ExternalId),
		Visibility:  convert.StrToType(ns.Visibility),
	}
}

func getNamespaceByID(
	ctx context.Context,
	id int64,
	clusterID int64,
	apiClient *sdk.APIClient,
) (*sdk.GetClusterNamespace200ResponseNamespace, error) {
	// GetClusterNamespace(ctx, clusterId, id)
	r, hresp, err := apiClient.ClustersAPI.GetClusterNamespace(
		ctx, clusterID, id,
	).Execute()
	if r == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"GET failed for cluster namespace %d: %s",
			id, providererrors.ErrMsg(err, hresp),
		)
	}

	return r.Namespace, nil
}

func getNamespaceByName(
	ctx context.Context,
	name string,
	clusterID int64,
	apiClient *sdk.APIClient,
) (*sdk.GetClusterNamespace200ResponseNamespace, error) {
	// GetClusterNamespaces(ctx, clusterId)
	rs, hresp, err := apiClient.ClustersAPI.GetClusterNamespaces(
		ctx, clusterID,
	).Execute()
	if rs == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"GET failed for cluster namespaces with name %s: %s",
			name, providererrors.ErrMsg(err, hresp),
		)
	}

	items := rs.Namespaces
	if len(items) == 0 {
		return nil, errors.New(ErrorNoClusterNamespaceFound)
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
		return nil, errors.New(ErrorNoClusterNamespaceFound)
	} else if len(matchedIDs) > 1 {
		return nil, errors.New(ErrorMultipleClusterNamespaces)
	}

	return getNamespaceByID(ctx, matchedIDs[0], clusterID, apiClient)
}

func getNamespace(
	ctx context.Context,
	config *ClusterNamespaceModel,
	apiClient *sdk.APIClient,
) (*sdk.GetClusterNamespace200ResponseNamespace, error) {
	clusterID := config.ClusterId.ValueInt64()

	if !config.Id.IsNull() {
		return getNamespaceByID(ctx, config.Id.ValueInt64(), clusterID, apiClient)
	} else if !config.Name.IsNull() {
		return getNamespaceByName(ctx, config.Name.ValueString(), clusterID, apiClient)
	}

	return nil, errors.New(ErrorNoValidSearchTerms)
}

// Read refreshes the Terraform state with the latest data.
func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config ClusterNamespaceModel

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

	ns, err := getNamespace(ctx, &config, apiClient)
	if err != nil {
		resp.Diagnostics.AddError(
			summary,
			err.Error(),
		)

		return
	}

	clusterID := config.ClusterId.ValueInt64()
	state := namespaceAsState(ns, clusterID)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
