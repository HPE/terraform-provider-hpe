// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package sweep

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	testsweep "github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/sweep"
)

// Clusters whose name begins with this string will be eligible for deletion.
const testResourcePrefix = "TestAccMorpheus"

func init() {
	testsweep.RegisterTypedAPISweeper(
		"hpe_morpheus_cluster",
		// List all cluster resources.
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

			return strings.HasPrefix(*name, testResourcePrefix)
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
	)
}
