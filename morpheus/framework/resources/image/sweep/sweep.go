// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

//go:build sweep

package sweep

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	testsweep "github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/sweep"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/getsafe"
	"github.com/HPE/terraform-provider-hpe/utils/paging"
)

const sweeperName = "hpe_morpheus_image"

func init() {
	testsweep.RegisterTypedAPISweeper(
		sweeperName,
		// List image resources.
		func(ctx context.Context, client *sdk.APIClient) (
			[]sdk.ListVirtualImages200ResponseAllOfVirtualImagesInner,
			*http.Response,
			error,
		) {
			// The caller distinguishes "nothing to sweep" from "the sweep
			// failed" by the list status, so one response has to come back
			// with the items.
			//
			// Only the first page's response is kept. Later pages are fetched
			// concurrently, so assigning from every one of them would be a
			// data race — and it would gain nothing: a 403 or 404 arrives on
			// the first request, before any wave is launched, and every
			// response after it on a successful walk is a 200.
			var (
				firstResp *http.Response
				once      sync.Once
			)

			items, err := paging.Collect(
				ctx,
				func(ctx context.Context, offset, max int64) (
					[]sdk.ListVirtualImages200ResponseAllOfVirtualImagesInner,
					int64,
					error,
				) {
					resp, hresp, err := client.LibraryAPI.ListVirtualImages(ctx).
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

					return getsafe.Get(&resp.VirtualImages), total, err
				},
			)

			return items, firstResp, err
		},
		// Is this a test image?
		func(item sdk.ListVirtualImages200ResponseAllOfVirtualImagesInner) bool {
			name, ok := getsafe.GetOk(item.Name)
			if !ok || name == nil {
				return false
			}

			return strings.HasPrefix(*name, testsweep.TestResourcePrefix)
		},
		// Delete the test image.
		func(
			ctx context.Context,
			client *sdk.APIClient,
			item sdk.ListVirtualImages200ResponseAllOfVirtualImagesInner,
		) (*http.Response, error) {
			id, ok := getsafe.GetOk(item.Id)
			if !ok || id == nil {
				return nil, fmt.Errorf("could not get ID")
			}

			_, hresp, err := client.LibraryAPI.RemoveVirtualImage(ctx, *id).Execute()

			return hresp, err
		},
		testsweep.WithIgnoreListStatuses[sdk.ListVirtualImages200ResponseAllOfVirtualImagesInner](
			http.StatusNotFound,
			http.StatusForbidden,
		),
	)
}
