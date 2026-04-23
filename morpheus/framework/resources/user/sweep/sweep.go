// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package sweep

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	testsweep "github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/sweep"
)

// Users whose name begins with this string will be eligible for deletion
const testUserPrefix = "TestAccMorpheusUser"

// Additionally, the user email must match one of these in order to be deleted
var sweepableEmails = map[string]bool{
	"foo@testacc.com": true,
	"bar@testacc.com": true,
}

func isSweepableEmail(email string) bool {
	return sweepableEmails[email]
}

func init() {
	testsweep.RegisterTypedAPISweeper(
		"hpe_morpheus_user",
		// List all user resources.
		func(ctx context.Context, client *sdk.APIClient) (
			[]sdk.ListUsers200ResponseAllOfUsersInner,
			*http.Response,
			error,
		) {
			resp, hresp, err := client.UsersAPI.ListUsers(ctx).Execute()
			if resp == nil {
				return nil, hresp, err
			}

			return resp.GetUsers(), hresp, err
		},
		// Is this a test user?
		func(item sdk.ListUsers200ResponseAllOfUsersInner) bool {
			username, ok := item.GetUsernameOk()
			if !ok || username == nil || !strings.HasPrefix(*username, testUserPrefix) {
				return false
			}

			email, ok := item.GetEmailOk()
			if !ok || email == nil || !isSweepableEmail(*email) {
				return false
			}

			return true
		},
		// Delete the test user.
		func(
			ctx context.Context,
			client *sdk.APIClient,
			item sdk.ListUsers200ResponseAllOfUsersInner,
		) (*http.Response, error) {
			id, ok := item.GetIdOk()
			if !ok || id == nil {
				return nil, fmt.Errorf("could not get ID")
			}

			_, hresp, err := client.UsersAPI.DeleteUser(ctx, *id).Execute()

			return hresp, err
		},
	)
}
