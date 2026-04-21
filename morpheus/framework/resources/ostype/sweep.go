// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package ostype

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	testsweep "github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/sweep"
)

const testOsTypePrefix = "TestAccMorpheusOsType"

func init() {
	testsweep.RegisterTypedAPISweeper(
		"hpe_morpheus_os_type",
		// List all OS type resources.
		func(ctx context.Context, client *sdk.APIClient) (
			[]sdk.ListOsTypes200ResponseAllOfOsTypesInner,
			*http.Response,
			error,
		) {
			resp, hresp, err := client.LibraryAPI.ListOsTypes(ctx).Execute()
			if resp == nil {
				return nil, hresp, err
			}

			return resp.GetOsTypes(), hresp, err
		},
		// Is this a test OS type?
		func(item sdk.ListOsTypes200ResponseAllOfOsTypesInner) bool {
			name, ok := item.GetNameOk()
			if !ok || name == nil {
				return false
			}

			return strings.HasPrefix(*name, testOsTypePrefix)
		},
		// Delete the test OS type.
		func(
			ctx context.Context,
			client *sdk.APIClient,
			item sdk.ListOsTypes200ResponseAllOfOsTypesInner,
		) (*http.Response, error) {
			id, ok := item.GetIdOk()
			if !ok || id == nil {
				return nil, fmt.Errorf("could not get ID")
			}

			_, hresp, err := client.LibraryAPI.DeleteOsType(ctx, *id).Execute()

			return hresp, err
		},
	)
}
