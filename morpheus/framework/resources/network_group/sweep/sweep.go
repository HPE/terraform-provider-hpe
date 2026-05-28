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

const sweeperName = "hpe_morpheus_network_group"

func init() {
	testsweep.RegisterTypedAPISweeper(
		sweeperName,
		// List network group resources.
		func(ctx context.Context, client *sdk.APIClient) (
			[]sdk.GetNetworkGroups200ResponseNetworkGroupsInner,
			*http.Response,
			error,
		) {
			resp, hresp, err := client.NetworksAPI.GetNetworkGroups(ctx).Execute()
			if resp == nil {
				return nil, hresp, err
			}

			return resp.GetNetworkGroups(), hresp, err
		},
		// Is this a test network group?
		func(item sdk.GetNetworkGroups200ResponseNetworkGroupsInner) bool {
			name, ok := item.GetNameOk()
			if !ok || name == nil {
				return false
			}

			return strings.HasPrefix(*name, testsweep.TestResourcePrefix)
		},
		// Delete the test network group.
		func(
			ctx context.Context,
			client *sdk.APIClient,
			item sdk.GetNetworkGroups200ResponseNetworkGroupsInner,
		) (*http.Response, error) {
			id, ok := item.GetIdOk()
			if !ok || id == nil {
				return nil, fmt.Errorf("could not get ID")
			}

			_, hresp, err := client.NetworksAPI.DeleteNetworkGroup(ctx, *id).Execute()

			return hresp, err
		},
		testsweep.WithIgnoreListStatuses[sdk.GetNetworkGroups200ResponseNetworkGroupsInner](
			http.StatusNotFound,
			http.StatusForbidden,
		),
	)
}
