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

const sweeperName = "hpe_morpheus_user"

func init() {
	testsweep.RegisterTypedAPISweeper(
		sweeperName,
		// List user resources.
		func(ctx context.Context, client *sdk.APIClient) (
			[]sdk.ListUsers200ResponseAllOfUsersInner,
			*http.Response,
			error,
		) {
			resp, hresp, err := client.UsersAPI.ListUsers(ctx).Execute()
			if resp == nil {
				return nil, hresp, err
			}

			return getsafe.GetSafe(&resp.Users), hresp, err
		},
		// Is this a test user?
		func(item sdk.ListUsers200ResponseAllOfUsersInner) bool {
			username, ok := getsafe.GetSafeOk(item.Username)
			if !ok || username == nil {
				return false
			}

			return strings.HasPrefix(*username, testsweep.TestResourcePrefix)
		},
		// Delete the test user.
		func(
			ctx context.Context,
			client *sdk.APIClient,
			item sdk.ListUsers200ResponseAllOfUsersInner,
		) (*http.Response, error) {
			id, ok := getsafe.GetSafeOk(item.Id)
			if !ok || id == nil {
				return nil, fmt.Errorf("could not get ID")
			}

			_, hresp, err := client.UsersAPI.DeleteUser(ctx, *id).Execute()

			return hresp, err
		},
	)
}
