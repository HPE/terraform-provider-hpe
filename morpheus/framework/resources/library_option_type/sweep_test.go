// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package library_option_type_test

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	testsweep "github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/sweep"
)

const sweeperName = "hpe_morpheus_library_option_type"

func init() {
	testsweep.RegisterTypedAPISweeper(
		sweeperName,
		// List library option type resources.
		func(ctx context.Context, client *sdk.APIClient) (
			[]sdk.ListInputs200ResponseAllOfOptionTypesInner,
			*http.Response,
			error,
		) {
			resp, hresp, err := client.LibraryAPI.ListInputs(ctx).Execute()
			if resp == nil {
				return nil, hresp, err
			}

			return resp.GetOptionTypes(), hresp, err
		},
		// Is this a test library option type?
		func(item sdk.ListInputs200ResponseAllOfOptionTypesInner) bool {
			name, ok := item.GetNameOk()
			if !ok || name == nil {
				return false
			}

			return strings.HasPrefix(*name, testsweep.TestResourcePrefix)
		},
		// Delete the test library option type.
		func(
			ctx context.Context,
			client *sdk.APIClient,
			item sdk.ListInputs200ResponseAllOfOptionTypesInner,
		) (*http.Response, error) {
			id, ok := item.GetIdOk()
			if !ok || id == nil {
				return nil, fmt.Errorf("could not get ID")
			}

			_, hresp, err := client.LibraryAPI.DeleteOptionType(ctx, *id).Execute()

			return hresp, err
		},
		testsweep.WithIgnoreListStatuses[sdk.ListInputs200ResponseAllOfOptionTypesInner](
			http.StatusNotFound,
			http.StatusForbidden,
		),
	)
}
