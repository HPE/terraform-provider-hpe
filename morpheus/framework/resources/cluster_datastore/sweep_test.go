// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package cluster_datastore_test

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	testsweep "github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/sweep"
)

const sweeperName = "hpe_morpheus_cluster_datastore"

// clusterDatastoreSweepItem pairs a datastore with its parent cluster ID.
type clusterDatastoreSweepItem struct {
	clusterID int64
	id        int64
	name      string
}

func init() {
	testsweep.RegisterTypedAPISweeper(
		sweeperName,
		// List cluster datastore resources.
		func(ctx context.Context, client *sdk.APIClient) (
			[]clusterDatastoreSweepItem,
			*http.Response,
			error,
		) {
			clusterResp, hresp, err := client.ClustersAPI.ListClusters(ctx).Execute()
			if err != nil {
				return nil, hresp, err
			}

			var items []clusterDatastoreSweepItem

			for _, cluster := range clusterResp.GetClusters() {
				clusterID, ok := cluster.GetIdOk()
				if !ok || clusterID == nil {
					continue
				}

				dsResp, _, err := client.ClustersAPI.
					ListClusterDatastores(ctx, *clusterID).Execute()
				if err != nil || dsResp == nil {
					continue
				}

				for _, ds := range dsResp.GetDatastores() {
					id, ok := ds.GetIdOk()
					if !ok || id == nil {
						continue
					}

					name, ok := ds.GetNameOk()
					if !ok || name == nil {
						continue
					}

					items = append(items, clusterDatastoreSweepItem{
						clusterID: *clusterID,
						id:        *id,
						name:      *name,
					})
				}
			}

			return items, hresp, nil
		},
		// Is this a test cluster datastore?
		func(item clusterDatastoreSweepItem) bool {
			return strings.HasPrefix(item.name, testsweep.TestResourcePrefix)
		},
		// Delete the test cluster datastore.
		func(
			ctx context.Context,
			client *sdk.APIClient,
			item clusterDatastoreSweepItem,
		) (*http.Response, error) {
			if item.id == 0 {
				return nil, fmt.Errorf("could not get ID")
			}

			_, hresp, err := client.ClustersAPI.
				DeleteClusterDatastore(ctx, item.clusterID, item.id).Execute()

			return hresp, err
		},
		testsweep.WithIgnoreListStatuses[clusterDatastoreSweepItem](
			http.StatusNotFound,
			http.StatusForbidden,
		),
	)
}
