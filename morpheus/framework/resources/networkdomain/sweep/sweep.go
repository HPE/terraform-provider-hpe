// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

//go:build sweep

package sweep

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	testsweep "github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/sweep"
)

const sweeperName = "hpe_morpheus_network_domain"

// networkDomainSweepItem is one network domain from the list.
type networkDomainSweepItem struct {
	id   int64
	name string
}

// networkDomainListItem decodes the untyped networkDomains list entries.
type networkDomainListItem struct {
	Id   *int64  `json:"id"`
	Name *string `json:"name"`
}

func init() {
	testsweep.RegisterTypedAPISweeper(
		sweeperName,
		// List network domains.
		func(ctx context.Context, client *sdk.APIClient) (
			[]networkDomainSweepItem,
			*http.Response,
			error,
		) {
			rs, hresp, err := client.NetworksAPI.GetNetworkDomains(ctx).Execute()
			if err != nil {
				return nil, hresp, err
			}

			raw, err := json.Marshal(rs.NetworkDomains)
			if err != nil {
				return nil, hresp, err
			}

			var entries []networkDomainListItem
			if err := json.Unmarshal(raw, &entries); err != nil {
				return nil, hresp, err
			}

			var items []networkDomainSweepItem

			for _, e := range entries {
				if e.Id == nil || e.Name == nil {
					continue
				}

				items = append(items, networkDomainSweepItem{id: *e.Id, name: *e.Name})
			}

			return items, hresp, nil
		},
		// Is this a test network domain?
		func(item networkDomainSweepItem) bool {
			return strings.HasPrefix(item.name, testsweep.TestResourcePrefix)
		},
		// Delete the test network domain.
		func(
			ctx context.Context,
			client *sdk.APIClient,
			item networkDomainSweepItem,
		) (*http.Response, error) {
			if item.id == 0 {
				return nil, fmt.Errorf("could not get ID")
			}

			_, hresp, err := client.NetworksAPI.
				DeleteNetworkDomain(ctx, item.id).Execute()

			return hresp, err
		},
		testsweep.WithIgnoreListStatuses[networkDomainSweepItem](
			http.StatusNotFound,
			http.StatusForbidden,
		),
	)
}
