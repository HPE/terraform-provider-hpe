// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package user

import (
	"context"
	"net/http"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
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
	testhelpers.RegisterTypedAPISweeper(
		"hpe_morpheus_user",
		testUserPrefix,
		func(ctx context.Context, client *sdk.APIClient, prefix string) ([]sdk.ListUsers200ResponseAllOfUsersInner, *http.Response, error) {
			resp, hresp, err := client.UsersAPI.ListUsers(ctx).Phrase(prefix).Execute()
			if resp == nil {
				return nil, hresp, err
			}

			return resp.GetUsers(), hresp, err
		},
		func(item sdk.ListUsers200ResponseAllOfUsersInner) (string, bool) {
			name, ok := item.GetUsernameOk()
			if !ok || name == nil {
				return "", false
			}

			return *name, true
		},
		func(item sdk.ListUsers200ResponseAllOfUsersInner) (int64, bool) {
			id, ok := item.GetIdOk()
			if !ok || id == nil {
				return 0, false
			}

			return *id, true
		},
		func(ctx context.Context, client *sdk.APIClient, id int64, _ sdk.ListUsers200ResponseAllOfUsersInner) (*http.Response, error) {
			_, hresp, err := client.UsersAPI.DeleteUser(ctx, id).Execute()
			return hresp, err
		},
		testhelpers.WithFilter(func(
			_ context.Context,
			_ *sdk.APIClient,
			user sdk.ListUsers200ResponseAllOfUsersInner,
		) (bool, string, error) {
			email, ok := user.GetEmailOk()
			if !ok || email == nil || !isSweepableEmail(*email) {
				return false, "email", nil
			}

			return true, "", nil
		}),
	)
}
