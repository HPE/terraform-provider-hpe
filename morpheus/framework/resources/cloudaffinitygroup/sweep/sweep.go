// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

//go:build sweep

package sweep

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	testsweep "github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/sweep"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/getsafe"
)

const sweeperName = "hpe_morpheus_cloud_affinity_group"

// cloudAffinityGroupSweepItem pairs an affinity group with its parent cloud ID.
type cloudAffinityGroupSweepItem struct {
	cloudID int64
	id      int64
	name    string
}

func init() {
	testsweep.RegisterTypedAPISweeper(
		sweeperName,
		// List cloud affinity group resources across all clouds.
		func(ctx context.Context, client *sdk.APIClient) (
			[]cloudAffinityGroupSweepItem,
			*http.Response,
			error,
		) {
			cloudResp, hresp, err := client.CloudsAPI.ListClouds(ctx).Execute()
			if err != nil {
				return nil, hresp, err
			}

			var items []cloudAffinityGroupSweepItem

			for _, cloud := range cloudResp.Zones {
				cloudID, ok := getsafe.GetOk(cloud.Id)
				if !ok || cloudID == nil {
					continue
				}

				agResp, _, err := client.CloudsAPI.
					ListCloudAffinityGroups(ctx, *cloudID).Execute()
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

					items = append(items, cloudAffinityGroupSweepItem{
						cloudID: *cloudID,
						id:      *id,
						name:    *name,
					})
				}
			}

			return items, hresp, nil
		},
		// Is this a test cloud affinity group?
		func(item cloudAffinityGroupSweepItem) bool {
			return strings.HasPrefix(item.name, testsweep.TestResourcePrefix)
		},
		// Delete the test cloud affinity group.
		func(
			ctx context.Context,
			client *sdk.APIClient,
			item cloudAffinityGroupSweepItem,
		) (*http.Response, error) {
			if item.id == 0 {
				return nil, fmt.Errorf("could not get ID")
			}

			_, hresp, err := client.CloudsAPI.
				DeleteCloudAffinityGroup(ctx, item.cloudID, item.id).Execute()

			return hresp, err
		},
		testsweep.WithIgnoreListStatuses[cloudAffinityGroupSweepItem](
			http.StatusNotFound,
			http.StatusForbidden,
		),
	)
}
