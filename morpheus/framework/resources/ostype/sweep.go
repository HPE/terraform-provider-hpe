// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package ostype

import (
	"context"
	"net/http"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

const testOsTypePrefix = "TestAccMorpheusOsType"

func init() {
	testhelpers.RegisterTypedAPISweeper(
		"hpe_morpheus_os_type",
		testOsTypePrefix,
		func(ctx context.Context, client *sdk.APIClient, prefix string) ([]sdk.ListOsTypes200ResponseAllOfOsTypesInner, *http.Response, error) {
			resp, hresp, err := client.LibraryAPI.ListOsTypes(ctx).Phrase(prefix).Execute()
			if resp == nil {
				return nil, hresp, err
			}

			return resp.GetOsTypes(), hresp, err
		},
		func(item sdk.ListOsTypes200ResponseAllOfOsTypesInner) (string, bool) {
			name, ok := item.GetNameOk()
			if !ok || name == nil {
				return "", false
			}

			return *name, true
		},
		func(item sdk.ListOsTypes200ResponseAllOfOsTypesInner) (int64, bool) {
			id, ok := item.GetIdOk()
			if !ok || id == nil {
				return 0, false
			}

			return *id, true
		},
		func(ctx context.Context, client *sdk.APIClient, id int64, _ sdk.ListOsTypes200ResponseAllOfOsTypesInner) (*http.Response, error) {
			_, hresp, err := client.LibraryAPI.DeleteOsType(ctx, id).Execute()
			return hresp, err
		},
	)
}
