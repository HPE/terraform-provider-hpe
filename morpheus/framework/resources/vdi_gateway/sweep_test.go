// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package vdi_gateway_test

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	testsweep "github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/sweep"
)

const sweeperName = "hpe_morpheus_vdi_gateway"

func init() {
	testsweep.RegisterTypedAPISweeper(
		sweeperName,
		// List VDI gateway resources.
		func(ctx context.Context, client *sdk.APIClient) (
			[]sdk.ListVDIGateways200ResponseAllOfVdiGatewaysInner,
			*http.Response,
			error,
		) {
			resp, hresp, err := client.VDIAPI.ListVDIGateways(ctx).Execute()
			if resp == nil {
				return nil, hresp, err
			}

			return resp.GetVdiGateways(), hresp, err
		},
		// Is this a test VDI gateway?
		func(item sdk.ListVDIGateways200ResponseAllOfVdiGatewaysInner) bool {
			name, ok := item.GetNameOk()
			if !ok || name == nil {
				return false
			}

			return strings.HasPrefix(*name, testsweep.TestResourcePrefix)
		},
		// Delete the test VDI gateway.
		func(
			ctx context.Context,
			client *sdk.APIClient,
			item sdk.ListVDIGateways200ResponseAllOfVdiGatewaysInner,
		) (*http.Response, error) {
			id, ok := item.GetIdOk()
			if !ok || id == nil {
				return nil, fmt.Errorf("could not get ID")
			}

			_, hresp, err := client.VDIAPI.RemoveVDIGateways(ctx, *id).Execute()

			return hresp, err
		},
		testsweep.WithIgnoreListStatuses[sdk.ListVDIGateways200ResponseAllOfVdiGatewaysInner](
			http.StatusNotFound,
			http.StatusForbidden,
		),
	)
}
