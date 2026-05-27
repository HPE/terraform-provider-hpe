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

const sweeperName = "hpe_morpheus_catalog_item_type"

func init() {
	testsweep.RegisterTypedAPISweeper(
		sweeperName,
		// List catalog item type resources.
		func(ctx context.Context, client *sdk.APIClient) (
			[]sdk.ListCatalogItemTypes200ResponseAllOfCatalogItemTypesInner,
			*http.Response,
			error,
		) {
			resp, hresp, err := client.CatalogItemsAPI.ListCatalogItemTypes(ctx).Execute()
			if resp == nil {
				return nil, hresp, err
			}

			return resp.GetCatalogItemTypes(), hresp, err
		},
		// Is this a test catalog item type?
		func(item sdk.ListCatalogItemTypes200ResponseAllOfCatalogItemTypesInner) bool {
			name, ok := item.GetNameOk()
			if !ok || name == nil {
				return false
			}

			return strings.HasPrefix(*name, testsweep.TestResourcePrefix)
		},
		// Delete the test catalog item type.
		func(
			ctx context.Context,
			client *sdk.APIClient,
			item sdk.ListCatalogItemTypes200ResponseAllOfCatalogItemTypesInner,
		) (*http.Response, error) {
			id, ok := item.GetIdOk()
			if !ok || id == nil {
				return nil, fmt.Errorf("could not get ID")
			}

			_, hresp, err := client.CatalogItemsAPI.RemoveCatalogItemType(ctx, *id).Execute()

			return hresp, err
		},
		testsweep.WithIgnoreListStatuses[sdk.ListCatalogItemTypes200ResponseAllOfCatalogItemTypesInner](
			http.StatusNotFound,
			http.StatusForbidden,
		),
	)
}
