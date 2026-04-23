// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package sweep

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	testsweep "github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/sweep"
)

func init() {
	testsweep.RegisterTypedAPISweeper(
		"hpe_morpheus_image",
		// List all image resources.
		func(ctx context.Context, client *sdk.APIClient) (
			[]sdk.ListVirtualImages200ResponseAllOfVirtualImagesInner,
			*http.Response,
			error,
		) {
			resp, hresp, err := client.LibraryAPI.ListVirtualImages(ctx).Execute()
			if resp == nil {
				return nil, hresp, err
			}

			return resp.GetVirtualImages(), hresp, err
		},
		// Is this a test image?
		func(item sdk.ListVirtualImages200ResponseAllOfVirtualImagesInner) bool {
			name, ok := item.GetNameOk()
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
			id, ok := item.GetIdOk()
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
