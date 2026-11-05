// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

//go:build sweep

package sweep

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	testsweep "github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/sweep"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
)

const sweeperName = "hpe_morpheus_network_pool_ip"

// networkPoolIpSweepItem pairs a reservation with its parent pool ID.
type networkPoolIpSweepItem struct {
	networkPoolID int64
	id            int64
	hostname      string
}

func init() {
	testsweep.RegisterTypedAPISweeper(
		sweeperName,
		// List network pool IP reservations across every pool.
		func(ctx context.Context, client *sdk.APIClient) (
			[]networkPoolIpSweepItem,
			*http.Response,
			error,
		) {
			poolsResp, hresp, err := client.NetworksAPI.
				GetNetworkPools(ctx).Execute()
			if err != nil {
				return nil, hresp, err
			}

			poolIDs, err := extractPoolIDs(poolsResp)
			if err != nil {
				return nil, hresp, err
			}

			var (
				items  []networkPoolIpSweepItem
				failed []int64
			)

			// A pool whose reservations cannot be listed is not skipped
			// quietly: its test reservations would stay behind while the sweep
			// reported a clean run. The failure is logged per pool and summed
			// up at the end, but NOT returned as an error: the shared runner
			// treats a list error as fatal and deletes nothing, which would
			// turn one pool's 403 or 5xx into every other pool's reservations
			// surviving too. What was listed is swept; what was not is said.
			for _, poolID := range poolIDs {
				ipsResp, hresp, err := client.NetworksAPI.
					GetNetworkPoolIps(ctx, poolID).Execute()
				if err != nil {
					log.Printf("[ERROR] network pool IP sweeper: cannot list IPs of pool %d, its test reservations are NOT swept: %s",
						poolID, errfmt.ErrMsg(err, hresp))
					failed = append(failed, poolID)

					continue
				}
				if ipsResp == nil {
					continue
				}

				for _, ip := range ipsResp.NetworkPoolIps {
					if ip.Id == nil || ip.Hostname == nil {
						continue
					}

					items = append(items, networkPoolIpSweepItem{
						networkPoolID: poolID,
						id:            *ip.Id,
						hostname:      *ip.Hostname,
					})
				}
			}

			if len(failed) > 0 {
				log.Printf("[ERROR] network pool IP sweeper: partial sweep, %d of %d pool(s) could not be listed: %v",
					len(failed), len(poolIDs), failed)
			}

			return items, hresp, nil
		},
		// Is this a test reservation (named after its acceptance test)?
		func(item networkPoolIpSweepItem) bool {
			return strings.HasPrefix(item.hostname, testsweep.TestResourcePrefix)
		},
		// Delete the test reservation.
		func(
			ctx context.Context,
			client *sdk.APIClient,
			item networkPoolIpSweepItem,
		) (*http.Response, error) {
			if item.id == 0 {
				return nil, fmt.Errorf("could not get ID")
			}

			_, hresp, err := client.NetworksAPI.
				DeleteNetworkPoolIp(ctx, item.networkPoolID, item.id).Execute()

			return hresp, err
		},
		testsweep.WithIgnoreListStatuses[networkPoolIpSweepItem](
			http.StatusNotFound,
			http.StatusForbidden,
		),
	)
}

// extractPoolIDs pulls pool IDs out of GetNetworkPools' loosely-typed
// networkPools payload (the generated model types it as interface{}). It is
// re-marshalled through JSON into a minimal typed shape so the sweeper does
// not depend on the concrete runtime type.
func extractPoolIDs(resp *sdk.GetNetworkPools200Response) ([]int64, error) {
	if resp == nil || resp.NetworkPools == nil {
		return nil, nil
	}

	raw, err := json.Marshal(resp.NetworkPools)
	if err != nil {
		return nil, fmt.Errorf("marshalling network pools: %w", err)
	}

	var pools []struct {
		ID *int64 `json:"id"`
	}
	if err := json.Unmarshal(raw, &pools); err != nil {
		return nil, fmt.Errorf("unmarshalling network pools: %w", err)
	}

	var ids []int64
	for _, p := range pools {
		if p.ID != nil {
			ids = append(ids, *p.ID)
		}
	}

	return ids, nil
}
