// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package ostype

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

const testOsTypePrefix = "TestAccMorpheusOsType"

func init() {
	testhelpers.RegisterTypedAPISweeper(
		"hpe_morpheus_os_type",
		// List candidate OS types for sweeping.
		func(ctx context.Context, client *sdk.APIClient, _ string) (
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
		// Match only acceptance-test OS types.
		func(item sdk.ListOsTypes200ResponseAllOfOsTypesInner) bool {
			name, ok := item.GetNameOk()
			if !ok || name == nil {
				return false
			}

			return strings.HasPrefix(*name, testOsTypePrefix)
		},
		// Delete a matched OS type.
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
