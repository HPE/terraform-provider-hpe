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
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/getsafe"
)

const sweeperName = "hpe_morpheus_cluster_affinity_group"

// clusterAffinityGroupSweepItem pairs an affinity group with its parent cluster ID.
type clusterAffinityGroupSweepItem struct {
	clusterID int64
	id        int64
	name      string
}

func init() {
	testsweep.RegisterTypedAPISweeper(
		sweeperName,
		// List cluster affinity group resources.
		func(ctx context.Context, client *sdk.APIClient) (
			[]clusterAffinityGroupSweepItem,
			*http.Response,
			error,
		) {
			clusterResp, hresp, err := client.ClustersAPI.ListClusters(ctx).Execute()
			if err != nil {
				return nil, hresp, err
			}

			var items []clusterAffinityGroupSweepItem

			for _, cluster := range clusterResp.Clusters {
				clusterID, ok := getsafe.GetOk(cluster.Id)
				if !ok || clusterID == nil {
					continue
				}

				agResp, _, err := client.ClustersAPI.
					ListClusterAffinityGroups(ctx, *clusterID).Execute()
				if err != nil || agResp == nil {
					continue
				}

				for _, ag := range agResp.AffinityGroups {
					id, ok := getsafe.GetOk(ag.Id)
					if !ok || id == nil {
						continue
					}

					name, ok := getsafe.GetOk(ag.Name)
					if !ok || name == nil {
						continue
					}

					items = append(items, clusterAffinityGroupSweepItem{
						clusterID: *clusterID,
						id:        *id,
						name:      *name,
					})
				}
			}

			return items, hresp, nil
		},
		// Is this a test cluster affinity group?
		func(item clusterAffinityGroupSweepItem) bool {
			return strings.HasPrefix(item.name, testsweep.TestResourcePrefix)
		},
		// Delete the test cluster affinity group.
		func(
			ctx context.Context,
			client *sdk.APIClient,
			item clusterAffinityGroupSweepItem,
		) (*http.Response, error) {
			if item.id == 0 {
				return nil, fmt.Errorf("could not get ID")
			}

			_, hresp, err := client.ClustersAPI.
				DeleteClusterAffinityGroup(ctx, item.clusterID, item.id).Execute()

			return hresp, err
		},
		testsweep.WithIgnoreListStatuses[clusterAffinityGroupSweepItem](
			http.StatusNotFound,
			http.StatusForbidden,
		),
	)
}
