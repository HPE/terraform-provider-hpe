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

const sweeperName = "hpe_morpheus_subnet"

func init() {
	testsweep.RegisterTypedAPISweeper(
		sweeperName,
		// List subnet resources.
		func(ctx context.Context, client *sdk.APIClient) (
			[]sdk.ListSubnets200ResponseAllOfSubnetsInner,
			*http.Response,
			error,
		) {
			resp, hresp, err := client.NetworksAPI.ListSubnets(ctx).Execute()
			if resp == nil {
				return nil, hresp, err
			}

			return getsafe.Get(&resp.Subnets), hresp, err
		},
		// Is this a test subnet?
		func(item sdk.ListSubnets200ResponseAllOfSubnetsInner) bool {
			name, ok := getsafe.GetOk(item.Name)
			if !ok || name == nil {
				return false
			}

			return strings.HasPrefix(*name, testsweep.TestResourcePrefix)
		},
		// Delete the test subnet.
		func(
			ctx context.Context,
			client *sdk.APIClient,
			item sdk.ListSubnets200ResponseAllOfSubnetsInner,
		) (*http.Response, error) {
			id, ok := getsafe.GetOk(item.Id)
			if !ok || id == nil {
				return nil, fmt.Errorf("could not get ID")
			}

			_, hresp, err := client.NetworksAPI.DeleteSubnet(ctx, *id).Execute()

			return hresp, err
		},
		testsweep.WithIgnoreListStatuses[sdk.ListSubnets200ResponseAllOfSubnetsInner](
			http.StatusNotFound,
			http.StatusForbidden,
		),
	)
}
