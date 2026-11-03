// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

//go:build sweep

package sweep

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	testsweep "github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/sweep"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/getsafe"
	"github.com/HPE/terraform-provider-hpe/utils/paging"
)

const sweeperName = "hpe_morpheus_plugin"

// testPluginCodes are the plugin codes the acceptance tests install. The
// sweeper removes only these, and never a prepackaged plugin.
var testPluginCodes = map[string]bool{
	"hpe-terraform-skeleton-plugin":   true,
	"morpheus-hashicorp-vault-plugin": true,
}

func init() {
	testsweep.RegisterTypedAPISweeper(
		sweeperName,
		// List every plugin. The endpoint pages (25 by default), and a
		// leaked test plugin is as likely to be on the second page as the
		// first.
		func(ctx context.Context, client *sdk.APIClient) (
			[]sdk.ListPlugins200ResponseAllOfPluginsInner,
			*http.Response,
			error,
		) {
			// Only the first page's response is kept for the list status:
			// later pages are fetched concurrently, and a 403 or 404 arrives
			// on the first request anyway.
			var (
				firstResp *http.Response
				once      sync.Once
			)

			items, err := paging.Collect(
				ctx,
				func(ctx context.Context, offset, max int64) (
					[]sdk.ListPlugins200ResponseAllOfPluginsInner,
					int64,
					error,
				) {
					resp, hresp, err := client.PluginsAPI.ListPlugins(ctx).
						Max(max).
						Offset(offset).
						Execute()

					once.Do(func() { firstResp = hresp })

					if resp == nil {
						return nil, 0, err
					}

					var total int64
					if resp.Meta != nil && resp.Meta.Total != nil {
						total = *resp.Meta.Total
					}

					return getsafe.Get(&resp.Plugins), total, err
				},
			)

			return items, firstResp, err
		},
		// Is this a test plugin?
		func(item sdk.ListPlugins200ResponseAllOfPluginsInner) bool {
			if prepackaged, ok := getsafe.GetOk(item.IsPrepackaged); ok && prepackaged != nil && *prepackaged {
				return false
			}

			code, ok := getsafe.GetOk(item.Code)
			if !ok || code == nil {
				return false
			}

			return testPluginCodes[*code]
		},
		// Delete the test plugin.
		func(
			ctx context.Context,
			client *sdk.APIClient,
			item sdk.ListPlugins200ResponseAllOfPluginsInner,
		) (*http.Response, error) {
			id, ok := getsafe.GetOk(item.Id)
			if !ok || id == nil {
				return nil, fmt.Errorf("could not get ID")
			}

			_, hresp, err := client.PluginsAPI.RemovePlugin(ctx, *id).Execute()

			return hresp, err
		},
		testsweep.WithIgnoreListStatuses[sdk.ListPlugins200ResponseAllOfPluginsInner](
			http.StatusNotFound,
			http.StatusForbidden,
		),
	)
}
