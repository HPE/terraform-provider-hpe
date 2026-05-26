// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package monitoring_group_test

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	testsweep "github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/sweep"
)

const sweeperName = "hpe_morpheus_monitoring_group"

func init() {
	testsweep.RegisterTypedAPISweeper(
		sweeperName,
		// List monitoring group resources.
		func(ctx context.Context, client *sdk.APIClient) (
			[]sdk.ListCheckGroups200ResponseAllOfCheckGroupsInner,
			*http.Response,
			error,
		) {
			resp, hresp, err := client.ChecksAPI.ListCheckGroups(ctx).Execute()
			if resp == nil {
				return nil, hresp, err
			}

			return resp.GetCheckGroups(), hresp, err
		},
		// Is this a test monitoring group?
		func(item sdk.ListCheckGroups200ResponseAllOfCheckGroupsInner) bool {
			name, ok := item.GetNameOk()
			if !ok || name == nil {
				return false
			}

			return strings.HasPrefix(*name, testsweep.TestResourcePrefix)
		},
		// Delete the test monitoring group.
		func(
			ctx context.Context,
			client *sdk.APIClient,
			item sdk.ListCheckGroups200ResponseAllOfCheckGroupsInner,
		) (*http.Response, error) {
			id, ok := item.GetIdOk()
			if !ok || id == nil {
				return nil, fmt.Errorf("could not get ID")
			}

			_, hresp, err := client.ChecksAPI.DeleteCheckGroups(ctx, *id).Execute()

			return hresp, err
		},
		testsweep.WithIgnoreListStatuses[sdk.ListCheckGroups200ResponseAllOfCheckGroupsInner](
			http.StatusNotFound,
			http.StatusForbidden,
		),
	)
}
