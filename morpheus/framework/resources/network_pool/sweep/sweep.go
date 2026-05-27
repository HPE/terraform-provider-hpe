// (C) Copyright 2026 Hewlett Packard Enterprise Development LP


//go:build sweep

package sweep

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	testsweep "github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/sweep"
)

const sweeperName = "hpe_morpheus_network_pool"

// networkPoolSweepItem holds the parsed fields needed to sweep a network pool.
// The SDK returns network pools as interface{} (untyped), so a JSON round-trip
// is used to extract the name and ID.
type networkPoolSweepItem struct {
	id   int64
	name string
}

func init() {
	testsweep.RegisterTypedAPISweeper(
		sweeperName,
		// List network pool resources.
		func(ctx context.Context, client *sdk.APIClient) (
			[]networkPoolSweepItem,
			*http.Response,
			error,
		) {
			resp, hresp, err := client.NetworksAPI.GetNetworkPools(ctx).Execute()
			if resp == nil {
				return nil, hresp, err
			}

			// The SDK returns network pools as interface{} — use a JSON round-trip
			// to extract typed name/ID fields.
			raw, err := json.Marshal(resp.NetworkPools)
			if err != nil {
				return nil, hresp, fmt.Errorf("marshaling network pools: %w", err)
			}

			var pools []map[string]interface{}
			if err := json.Unmarshal(raw, &pools); err != nil {
				return nil, hresp, fmt.Errorf("decoding network pools: %w", err)
			}

			items := make([]networkPoolSweepItem, 0, len(pools))

			for _, p := range pools {
				name, _ := p["name"].(string)

				var id int64

				switch v := p["id"].(type) {
				case float64:
					id = int64(v)
				case json.Number:
					id, _ = v.Int64()
				default:
					continue
				}

				items = append(items, networkPoolSweepItem{id: id, name: name})
			}

			return items, hresp, nil
		},
		// Is this a test network pool?
		func(item networkPoolSweepItem) bool {
			return strings.HasPrefix(item.name, testsweep.TestResourcePrefix)
		},
		// Delete the test network pool.
		func(
			ctx context.Context,
			client *sdk.APIClient,
			item networkPoolSweepItem,
		) (*http.Response, error) {
			if item.id == 0 {
				return nil, fmt.Errorf("could not get ID")
			}

			_, hresp, err := client.NetworksAPI.DeleteNetworkPool(ctx, item.id).Execute()

			return hresp, err
		},
		testsweep.WithIgnoreListStatuses[networkPoolSweepItem](
			http.StatusNotFound,
			http.StatusForbidden,
		),
	)
}
