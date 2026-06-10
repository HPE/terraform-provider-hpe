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

const sweeperName = "hpe_morpheus_vdi_pool"

func init() {
	testsweep.RegisterTypedAPISweeper(
		sweeperName,
		// List VDI pool resources.
		func(ctx context.Context, client *sdk.APIClient) (
			[]sdk.ListVDIPools200ResponseAllOfVdiPoolsInner,
			*http.Response,
			error,
		) {
			resp, hresp, err := client.VDIAPI.ListVDIPools(ctx).Execute()
			if resp == nil {
				return nil, hresp, err
			}

			return resp.VdiPools, hresp, err
		},
		// Is this a test VDI pool?
		func(item sdk.ListVDIPools200ResponseAllOfVdiPoolsInner) bool {
			if item.Name == nil {
				return false
			}

			return strings.HasPrefix(*item.Name, testsweep.TestResourcePrefix)
		},
		// Delete the test VDI pool.
		func(
			ctx context.Context,
			client *sdk.APIClient,
			item sdk.ListVDIPools200ResponseAllOfVdiPoolsInner,
		) (*http.Response, error) {
			if item.Id == nil {
				return nil, fmt.Errorf("could not get ID")
			}

			_, hresp, err := client.VDIAPI.RemoveVDIPools(ctx, *item.Id).Execute()

			return hresp, err
		},
		testsweep.WithIgnoreListStatuses[sdk.ListVDIPools200ResponseAllOfVdiPoolsInner](
			http.StatusNotFound,
			http.StatusForbidden,
		),
	)
}
