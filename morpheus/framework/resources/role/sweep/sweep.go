// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

//go:build sweep

package sweep

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

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
		// Morpheus refuses to delete a role that is any tenant's base role
		// ("Role is already in use"), and the tenant tests create tenants
		// whose base role is a test role. Sweep tenants first so their roles
		// can go in the same run.
		testsweep.WithDependencies[sdk.ListRoles200ResponseAllOfRolesInner](
			"hpe_morpheus_tenant",
		),
	)
}
