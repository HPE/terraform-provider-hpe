// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package cypher_test

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	testsweep "github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/sweep"
)

const sweeperName = "hpe_morpheus_cypher"

// cypherSweepItem holds the path (ItemKey) needed to identify and delete a cypher entry.
// The SDK uses the path string as the resource identifier rather than a numeric ID.
type cypherSweepItem struct {
	path string
}

func init() {
	testsweep.RegisterTypedAPISweeper(
		sweeperName,
		// List cypher resources.
		func(ctx context.Context, client *sdk.APIClient) (
			[]cypherSweepItem,
			*http.Response,
			error,
		) {
			resp, hresp, err := client.CypherAPI.ListCypherKeys(ctx).Execute()
			if err != nil {
				return nil, hresp, err
			}

			cyphers := resp.GetCyphers()
			items := make([]cypherSweepItem, 0, len(cyphers))

			for _, c := range cyphers {
				path, ok := c.GetItemKeyOk()
				if !ok || path == nil {
					continue
				}

				items = append(items, cypherSweepItem{path: *path})
			}

			return items, hresp, nil
		},
		// Is this a test cypher?
		func(item cypherSweepItem) bool {
			return strings.Contains(item.path, testsweep.TestResourcePrefix)
		},
		// Delete the test cypher.
		func(
			ctx context.Context,
			client *sdk.APIClient,
			item cypherSweepItem,
		) (*http.Response, error) {
			if item.path == "" {
				return nil, fmt.Errorf("could not get path")
			}

			_, hresp, err := client.CypherAPI.RemoveCypher(ctx, item.path).Execute()

			return hresp, err
		},
		testsweep.WithIgnoreListStatuses[cypherSweepItem](
			http.StatusNotFound,
			http.StatusForbidden,
		),
	)
}
