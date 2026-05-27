// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package option_list_test

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	testsweep "github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/sweep"
)

const sweeperName = "hpe_morpheus_option_list"

func init() {
	testsweep.RegisterTypedAPISweeper(
		sweeperName,
		// List library option type list resources.
		func(ctx context.Context, client *sdk.APIClient) (
			[]sdk.ListOptionLists200ResponseAllOfOptionTypesInner,
			*http.Response,
			error,
		) {
			resp, hresp, err := client.LibraryAPI.ListOptionLists(ctx).Execute()
			if resp == nil {
				return nil, hresp, err
			}

			return resp.GetOptionTypes(), hresp, err
		},
		// Is this a test library option type list?
		func(item sdk.ListOptionLists200ResponseAllOfOptionTypesInner) bool {
			name, ok := item.GetNameOk()
			if !ok || name == nil {
				return false
			}

			return strings.HasPrefix(*name, testsweep.TestResourcePrefix)
		},
		// Delete the test library option type list.
		func(
			ctx context.Context,
			client *sdk.APIClient,
			item sdk.ListOptionLists200ResponseAllOfOptionTypesInner,
		) (*http.Response, error) {
			id, ok := item.GetIdOk()
			if !ok || id == nil {
				return nil, fmt.Errorf("could not get ID")
			}

			_, hresp, err := client.LibraryAPI.DeleteOptionList(ctx, *id).Execute()

			return hresp, err
		},
		testsweep.WithIgnoreListStatuses[sdk.ListOptionLists200ResponseAllOfOptionTypesInner](
			http.StatusNotFound,
			http.StatusForbidden,
		),
	)
}
