// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

//go:build sweep

package sweep

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	testsweep "github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/sweep"
)

const sweeperName = "hpe_morpheus_virtual_switch"

// virtualSwitchSweepItem is a switch to sweep; the cluster it belongs to is
// the one whose list it came from (the delete needs only the switch id).
type virtualSwitchSweepItem struct {
	id   int64
	name string
}

func init() {
	testsweep.RegisterTypedAPISweeper(
		sweeperName,
		// List virtual switches across every cluster.
		func(ctx context.Context, client *sdk.APIClient) (
			[]virtualSwitchSweepItem,
			*http.Response,
			error,
		) {
			clusterResp, hresp, err := client.ClustersAPI.ListClusters(ctx).
				Max(testsweep.ListPageSize).Execute()
			if err != nil {
				return nil, hresp, err
			}

			var items []virtualSwitchSweepItem

			for _, cluster := range clusterResp.Clusters {
				if cluster.Id == nil {
					continue
				}

				swResp, swHresp, err := client.VirtualSwitchesAPI.
					ListVirtualSwitches(ctx).ClusterId(*cluster.Id).Execute()
				if err != nil || swResp == nil {
					// A cluster that cannot carry virtual switches answers 404;
					// any other failure must not silently shrink the sweep.
					if swHresp != nil && swHresp.StatusCode == http.StatusNotFound {
						continue
					}

					// A nil response: the harness would read a 403 on it as
					// "nothing to sweep" and end the whole sweep as a success.
					return nil, nil, fmt.Errorf("listing virtual switches on cluster %d: %s", *cluster.Id, errfmt.ErrMsg(err, swHresp))
				}

				for _, sw := range swResp.VirtualSwitches {
					if sw.Id == nil || sw.Name == nil {
						continue
					}

					items = append(items, virtualSwitchSweepItem{
						id:   *sw.Id,
						name: *sw.Name,
					})
				}
			}

			return items, hresp, nil
		},
		// Is this a test virtual switch?
		func(item virtualSwitchSweepItem) bool {
			// Switch names are limited to 16 characters, so the tests use
			// their own short prefix rather than the suite-wide one.
			return strings.HasPrefix(item.name, TestNamePrefix) ||
				strings.HasPrefix(item.name, testsweep.TestResourcePrefix)
		},
		// Delete the test virtual switch (force past any in-flight sync).
		func(
			ctx context.Context,
			client *sdk.APIClient,
			item virtualSwitchSweepItem,
		) (*http.Response, error) {
			if item.id == 0 {
				return nil, fmt.Errorf("could not get ID")
			}

			// Morpheus renders a refused delete as HTTP 200 with success:false,
			// so inspect the envelope: a refusal must not count as swept.
			result, hresp, err := client.VirtualSwitchesAPI.
				DeleteVirtualSwitch(ctx, item.id).Force(true).Execute()

			return hresp, errfmt.CheckDeleted(result, err, hresp)
		},
		testsweep.WithIgnoreListStatuses[virtualSwitchSweepItem](
			http.StatusNotFound,
			http.StatusForbidden,
		),
		// A switch with networks on it cannot be deleted: sweep them first.
		testsweep.WithDependencies[virtualSwitchSweepItem]("hpe_morpheus_network"),
	)
}

// TestNamePrefix is the name prefix of virtual switches created by the
// acceptance tests. Switch names are limited to 16 characters, so the
// suite-wide sweep prefix cannot be used; the sweeper matches this one. It
// lives here rather than in the resource package so the sweeper does not
// import the resource it sweeps (an import cycle under the sweep build tag).
const TestNamePrefix = "tfacc"
