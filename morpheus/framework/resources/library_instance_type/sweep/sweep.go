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
)

const sweeperName = "hpe_morpheus_library_instance_type"

func init() {
	testsweep.RegisterTypedAPISweeper(
		sweeperName,
		// List library instance type resources.
		func(ctx context.Context, client *sdk.APIClient) (
			[]sdk.ListInstanceTypes200ResponseAllOfInstanceTypesInner,
			*http.Response,
			error,
		) {
			resp, hresp, err := client.LibraryAPI.ListInstanceTypes(ctx).Execute()
			if resp == nil {
				return nil, hresp, err
			}

			return resp.GetInstanceTypes(), hresp, err
		},
		// Is this a test library instance type?
		func(item sdk.ListInstanceTypes200ResponseAllOfInstanceTypesInner) bool {
			name, ok := item.GetNameOk()
			if !ok || name == nil {
				return false
			}

			return strings.HasPrefix(*name, testsweep.TestResourcePrefix)
		},
		// Delete the test library instance type.
		func(
			ctx context.Context,
			client *sdk.APIClient,
			item sdk.ListInstanceTypes200ResponseAllOfInstanceTypesInner,
		) (*http.Response, error) {
			id, ok := item.GetIdOk()
			if !ok || id == nil {
				return nil, fmt.Errorf("could not get ID")
			}

			_, hresp, err := client.LibraryAPI.DeleteInstanceType(ctx, *id).Execute()

			return hresp, err
		},
		testsweep.WithIgnoreListStatuses[sdk.ListInstanceTypes200ResponseAllOfInstanceTypesInner](
			http.StatusNotFound,
			http.StatusForbidden,
		),
	)
}
