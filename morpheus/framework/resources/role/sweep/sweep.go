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

const sweeperName = "hpe_morpheus_role"

func init() {
	testsweep.RegisterTypedAPISweeper(
		sweeperName,
		// List role resources.
		func(ctx context.Context, client *sdk.APIClient) (
			[]sdk.ListRoles200ResponseAllOfRolesInner,
			*http.Response,
			error,
		) {
			resp, hresp, err := client.RolesAPI.ListRoles(ctx).Execute()
			if resp == nil {
				return nil, hresp, err
			}

			return getsafe.Get(&resp.Roles), hresp, err
		},
		// Is this a test role?
		func(item sdk.ListRoles200ResponseAllOfRolesInner) bool {
			name, ok := getsafe.GetOk(item.Name)
			if !ok || name == nil {
				return false
			}

			return strings.HasPrefix(*name, testsweep.TestResourcePrefix)
		},
		// Delete the test role.
		func(
			ctx context.Context,
			client *sdk.APIClient,
			item sdk.ListRoles200ResponseAllOfRolesInner,
		) (*http.Response, error) {
			id, ok := getsafe.GetOk(item.Id)
			if !ok || id == nil {
				return nil, fmt.Errorf("could not get ID")
			}

			_, hresp, err := client.RolesAPI.DeleteRole(ctx, *id).Execute()

			return hresp, err
		},
		testsweep.WithIgnoreListStatuses[sdk.ListRoles200ResponseAllOfRolesInner](
			http.StatusNotFound,
			http.StatusForbidden,
		),
	)
}
