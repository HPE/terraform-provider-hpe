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

const sweeperName = "hpe_morpheus_os_type"

func init() {
	testsweep.RegisterTypedAPISweeper(
		sweeperName,
		// List OS type resources.
		func(ctx context.Context, client *sdk.APIClient) (
			[]sdk.ListOsTypes200ResponseAllOfOsTypesInner,
			*http.Response,
			error,
		) {
			resp, hresp, err := client.LibraryAPI.ListOsTypes(ctx).Execute()
			if resp == nil {
				return nil, hresp, err
			}

			return getsafe.GetSafe(&resp.OsTypes), hresp, err
		},
		// Is this a test OS type?
		func(item sdk.ListOsTypes200ResponseAllOfOsTypesInner) bool {
			name, ok := getsafe.GetSafeOk(item.Name)
			if !ok || name == nil {
				return false
			}

			return strings.HasPrefix(*name, testsweep.TestResourcePrefix)
		},
		// Delete the test OS type.
		func(
			ctx context.Context,
			client *sdk.APIClient,
			item sdk.ListOsTypes200ResponseAllOfOsTypesInner,
		) (*http.Response, error) {
			id, ok := getsafe.GetSafeOk(item.Id)
			if !ok || id == nil {
				return nil, fmt.Errorf("could not get ID")
			}

			_, hresp, err := client.LibraryAPI.DeleteOsType(ctx, *id).Execute()

			return hresp, err
		},
		testsweep.WithDependencies[sdk.ListOsTypes200ResponseAllOfOsTypesInner](
			"hpe_morpheus_os_type_image",
		),
	)
}
