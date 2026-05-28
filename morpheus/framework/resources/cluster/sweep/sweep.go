// (C) Copyright 2026 Hewlett Packard Enterprise Development LP


//go:build sweep

package sweep

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	testsweep "github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/sweep"
)

const sweeperName = "hpe_morpheus_cluster"

func init() {
	testsweep.RegisterTypedAPISweeper(
		sweeperName,
		// List cluster resources.
		func(ctx context.Context, client *sdk.APIClient) (
			[]sdk.ListClusters200ResponseAllOfClustersInner,
			*http.Response,
			error,
		) {
			resp, hresp, err := client.ClustersAPI.ListClusters(ctx).Execute()
			if resp == nil {
				return nil, hresp, err
			}

			return resp.GetClusters(), hresp, err
		},
		// Is this a test cluster?
		func(item sdk.ListClusters200ResponseAllOfClustersInner) bool {
			name, ok := item.GetNameOk()
			if !ok || name == nil {
				return false
			}

			return strings.HasPrefix(*name, testsweep.TestResourcePrefix)
		},
		// Delete the test cluster.
		func(
			ctx context.Context,
			client *sdk.APIClient,
			item sdk.ListClusters200ResponseAllOfClustersInner,
		) (*http.Response, error) {
			id, ok := item.GetIdOk()
			if !ok || id == nil {
				return nil, fmt.Errorf("could not get ID")
			}

			_, hresp, err := client.ClustersAPI.DeleteCluster(ctx, *id).Force("on").Execute()

			return hresp, err
		},
		// Ignore (i.e. just log) "not found" and "forbidden" errors.
		testsweep.WithIgnoreListStatuses[sdk.ListClusters200ResponseAllOfClustersInner](
			http.StatusNotFound,
			http.StatusForbidden,
		),
		testsweep.WithDependencies[sdk.ListClusters200ResponseAllOfClustersInner](
			"hpe_morpheus_cluster_affinity_group",
			"hpe_morpheus_cluster_datastore",
			"hpe_morpheus_cluster_namespace",
		),
	)
}
