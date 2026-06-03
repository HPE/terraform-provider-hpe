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

const sweeperName = "hpe_morpheus_group"

func init() {
	testsweep.RegisterTypedAPISweeper(
		sweeperName,
		// List group resources.
		func(ctx context.Context, client *sdk.APIClient) (
			[]sdk.ListGroups200ResponseAllOfGroupsInner,
			*http.Response,
			error,
		) {
			resp, hresp, err := client.GroupsAPI.ListGroups(ctx).Execute()
			if resp == nil {
				return nil, hresp, err
			}

			return getsafe.GetSafe(&resp.Groups), hresp, err
		},
		// Is this a test group?
		func(item sdk.ListGroups200ResponseAllOfGroupsInner) bool {
			name, ok := getsafe.GetSafeOk(item.Name)
			if !ok || name == nil {
				return false
			}

			return strings.HasPrefix(*name, testsweep.TestResourcePrefix)
		},
		// Delete the test group.
		func(
			ctx context.Context,
			client *sdk.APIClient,
			item sdk.ListGroups200ResponseAllOfGroupsInner,
		) (*http.Response, error) {
			id, ok := getsafe.GetSafeOk(item.Id)
			if !ok || id == nil {
				return nil, fmt.Errorf("could not get ID")
			}

			_, hresp, err := client.GroupsAPI.RemoveGroups(ctx, *id).Execute()

			return hresp, err
		},
		testsweep.WithIgnoreListStatuses[sdk.ListGroups200ResponseAllOfGroupsInner](
			http.StatusNotFound,
			http.StatusForbidden,
		),
	)
}
