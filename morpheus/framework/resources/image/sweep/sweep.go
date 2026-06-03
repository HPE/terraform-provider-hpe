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
			resp, hresp, err := client.LibraryAPI.ListVirtualImages(ctx).Execute()
			if resp == nil {
				return nil, hresp, err
			}

			return getsafe.GetSafe(&resp.VirtualImages), hresp, err
		},
		// Is this a test image?
		func(item sdk.ListVirtualImages200ResponseAllOfVirtualImagesInner) bool {
			name, ok := getsafe.GetSafeOk(item.Name)
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
			id, ok := getsafe.GetSafeOk(item.Id)
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
