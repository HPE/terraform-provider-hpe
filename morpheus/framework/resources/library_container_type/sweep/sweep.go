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

const sweeperName = "hpe_morpheus_library_container_type"

func init() {
	testsweep.RegisterTypedAPISweeper(
		sweeperName,
		// List library container type resources.
		func(ctx context.Context, client *sdk.APIClient) (
			[]sdk.ListNodeTypes200ResponseAllOfContainerTypesInner,
			*http.Response,
			error,
		) {
			resp, hresp, err := client.LibraryAPI.ListNodeTypes(ctx).Execute()
			if resp == nil {
				return nil, hresp, err
			}

			return resp.GetContainerTypes(), hresp, err
		},
		// Is this a test library container type?
		func(item sdk.ListNodeTypes200ResponseAllOfContainerTypesInner) bool {
			name, ok := item.GetNameOk()
			if !ok || name == nil {
				return false
			}

			return strings.HasPrefix(*name, testsweep.TestResourcePrefix)
		},
		// Delete the test library container type.
		func(
			ctx context.Context,
			client *sdk.APIClient,
			item sdk.ListNodeTypes200ResponseAllOfContainerTypesInner,
		) (*http.Response, error) {
			id, ok := item.GetIdOk()
			if !ok || id == nil {
				return nil, fmt.Errorf("could not get ID")
			}

			_, hresp, err := client.LibraryAPI.DeleteNodeType(ctx, int64(*id)).Execute()

			return hresp, err
		},
		testsweep.WithIgnoreListStatuses[sdk.ListNodeTypes200ResponseAllOfContainerTypesInner](
			http.StatusNotFound,
			http.StatusForbidden,
		),
	)
}
