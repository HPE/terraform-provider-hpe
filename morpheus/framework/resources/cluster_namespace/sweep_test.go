// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package cluster_namespace_test

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	testsweep "github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/sweep"
)

const sweeperName = "hpe_morpheus_cluster_namespace"

// clusterNamespaceSweepItem pairs a namespace with its parent cluster ID.
type clusterNamespaceSweepItem struct {
	clusterID int64
	id        int64
	name      string
}

func init() {
	testsweep.RegisterTypedAPISweeper(
		sweeperName,
		// List cluster namespace resources.
		func(ctx context.Context, client *sdk.APIClient) (
			[]clusterNamespaceSweepItem,
			*http.Response,
			error,
		) {
			clusterResp, hresp, err := client.ClustersAPI.ListClusters(ctx).Execute()
			if err != nil {
				return nil, hresp, err
			}

			var items []clusterNamespaceSweepItem

			for _, cluster := range clusterResp.GetClusters() {
				clusterID, ok := cluster.GetIdOk()
				if !ok || clusterID == nil {
					continue
				}

				nsResp, _, err := client.ClustersAPI.
					GetClusterNamespaces(ctx, *clusterID).Execute()
				if err != nil || nsResp == nil {
					continue
				}

				for _, ns := range nsResp.GetNamespaces() {
					id, ok := ns.GetIdOk()
					if !ok || id == nil {
						continue
					}

					name, ok := ns.GetNameOk()
					if !ok || name == nil {
						continue
					}

					items = append(items, clusterNamespaceSweepItem{
						clusterID: *clusterID,
						id:        *id,
						name:      *name,
					})
				}
			}

			return items, hresp, nil
		},
		// Is this a test cluster namespace?
		func(item clusterNamespaceSweepItem) bool {
			return strings.HasPrefix(item.name, testsweep.TestResourcePrefix)
		},
		// Delete the test cluster namespace.
		func(
			ctx context.Context,
			client *sdk.APIClient,
			item clusterNamespaceSweepItem,
		) (*http.Response, error) {
			if item.id == 0 {
				return nil, fmt.Errorf("could not get ID")
			}

			_, hresp, err := client.ClustersAPI.
				DeleteClusterNamespace(ctx, item.clusterID, item.id).Execute()

			return hresp, err
		},
		testsweep.WithIgnoreListStatuses[clusterNamespaceSweepItem](
			http.StatusNotFound,
			http.StatusForbidden,
		),
	)
}
