// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package cluster

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

const (
	summary = "read cluster data source"

	ErrorNoValidSearchTerms = "no valid search terms - an id or name is required"
	ErrorRunningPreApply    = "Error running pre-apply plan: exit status 1"
	ErrorNoClusterFound     = "no cluster found"
	ErrorMultipleClusters   = "multiple clusters were returned"
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
	resp.TypeName = req.ProviderTypeName + "_morpheus_cluster"
}

func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = ClusterDataSourceSchema(ctx)
}

func populateClusterData(
	data *ClusterModel,
	cluster *sdk.GetCluster200ResponseCluster,
) {
	data.CloudId = convert.Int64ToType(cluster.Zone.Id)
	data.Description = convert.StrToType(cluster.Description.Get())
	data.GroupId = convert.Int64ToType(cluster.Site.Id)
	data.Id = convert.Int64ToType(cluster.Id)
	data.Labels = convert.StrSliceToSet(cluster.GetLabels())
	data.LayoutId = convert.Int64ToType(cluster.Layout.Id)
	data.Name = convert.StrToType(cluster.Name)
	data.ServiceUrl = convert.StrToType(cluster.ServiceUrl.Get())
	data.Uuid = convert.StrToType(cluster.Uuid)
}

func getClusterByID(
	ctx context.Context,
	id int64,
	data *ClusterModel,
	apiClient *sdk.APIClient,
) error {
	clusterResp, hresp, err := apiClient.ClustersAPI.GetCluster(ctx, id).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		return fmt.Errorf("cluster %d GET failed: %s", id, errfmt.ErrMsg(err, hresp))
	}

	cluster, ok := clusterResp.GetClusterOk()
	if !ok || cluster == nil {
		return errors.New(ErrorNoClusterFound)
	}

	populateClusterData(data, cluster)

	return nil
}

func getClusterByName(
	ctx context.Context,
	name string,
	data *ClusterModel,
	apiClient *sdk.APIClient,
) error {
	clustersResp, hresp, err := apiClient.ClustersAPI.ListClusters(ctx).Name(name).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		return fmt.Errorf("cluster %s GET failed: %s", name, errfmt.ErrMsg(err, hresp))
	}

	var matchingClusters []sdk.ListClusters200ResponseAllOfClustersInner
	for _, c := range clustersResp.GetClusters() {
		if c.GetName() == name {
			matchingClusters = append(matchingClusters, c)
		}
	}

	if len(matchingClusters) == 0 {
		return errors.New(ErrorNoClusterFound)
	}

	if len(matchingClusters) > 1 {
		return errors.New(ErrorMultipleClusters)
	}

	id, ok := matchingClusters[0].GetIdOk()
	if !ok || id == nil {
		return errors.New(ErrorNoClusterFound)
	}

	return getClusterByID(ctx, *id, data, apiClient)
}

func getCluster(
	ctx context.Context,
	data *ClusterModel,
	apiClient *sdk.APIClient,
) error {
	if !data.Id.IsNull() {
		return getClusterByID(ctx, data.Id.ValueInt64(), data, apiClient)
	}

	if !data.Name.IsNull() {
		return getClusterByName(ctx, data.Name.ValueString(), data, apiClient)
	}

	return errors.New(ErrorNoValidSearchTerms)
}

func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var data ClusterModel

	diags := req.Config.Get(ctx, &data)
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

	if err := getCluster(ctx, &data, apiClient); err != nil {
		resp.Diagnostics.AddError(
			summary,
			err.Error(),
		)

		return
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
