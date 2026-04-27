// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package sweep

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	testsweep "github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/sweep"
)

func init() {
	testsweep.RegisterTypedAPISweeper(
		"hpe_morpheus_role",
		// List all role resources.
		func(ctx context.Context, client *sdk.APIClient) (
			[]sdk.ListRoles200ResponseAllOfRolesInner,
			*http.Response,
			error,
		) {
			resp, hresp, err := client.RolesAPI.ListRoles(ctx).Execute()
			if resp == nil {
				return nil, hresp, err
			}

			return resp.GetRoles(), hresp, err
		},
		// Is this a test role?
		func(item sdk.ListRoles200ResponseAllOfRolesInner) bool {
			name, ok := item.GetNameOk()
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
			id, ok := item.GetIdOk()
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
