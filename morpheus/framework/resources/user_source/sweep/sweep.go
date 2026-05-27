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

const sweeperName = "hpe_morpheus_user_source"

// userSourceSweepItem holds the parsed fields needed to sweep a user source.
// The SDK returns user sources as a polymorphic anyOf type, so a JSON
// round-trip is used to extract the common name and ID fields.
type userSourceSweepItem struct {
	id   int64
	name string
}

func init() {
	testsweep.RegisterTypedAPISweeper(
		sweeperName,
		// List user source resources.
		func(ctx context.Context, client *sdk.APIClient) (
			[]userSourceSweepItem,
			*http.Response,
			error,
		) {
			resp, hresp, err := client.IdentitySourcesAPI.ListIdentitySources(ctx).Execute()
			if resp == nil {
				return nil, hresp, err
			}

			sources := resp.GetUserSources()
			items := make([]userSourceSweepItem, 0, len(sources))

			for _, src := range sources {
				// Use JSON round-trip to extract common fields from the polymorphic type.
				payload, err := json.Marshal(src)
				if err != nil {
					continue
				}

				var parsed struct {
					ID   *int64  `json:"id"`
					Name *string `json:"name"`
				}

				if err := json.Unmarshal(payload, &parsed); err != nil || parsed.ID == nil || parsed.Name == nil {
					continue
				}

				items = append(items, userSourceSweepItem{id: *parsed.ID, name: *parsed.Name})
			}

			return items, hresp, nil
		},
		// Is this a test user source?
		func(item userSourceSweepItem) bool {
			return strings.HasPrefix(item.name, testsweep.TestResourcePrefix)
		},
		// Delete the test user source.
		func(
			ctx context.Context,
			client *sdk.APIClient,
			item userSourceSweepItem,
		) (*http.Response, error) {
			if item.id == 0 {
				return nil, fmt.Errorf("could not get ID")
			}

			_, hresp, err := client.IdentitySourcesAPI.RemoveIdentitySources(ctx, item.id).Execute()

			return hresp, err
		},
		testsweep.WithIgnoreListStatuses[userSourceSweepItem](
			http.StatusNotFound,
			http.StatusForbidden,
		),
	)
}
