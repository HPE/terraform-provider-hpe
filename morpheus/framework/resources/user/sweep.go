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
	testhelpers.RegisterAPISweeper(
		"hpe_morpheus_user",
		testUserPrefix,
		func(ctx context.Context, client *sdk.APIClient, prefix string) (any, *http.Response, error) {
			return client.UsersAPI.ListUsers(ctx).Phrase(prefix).Execute()
		},
		"GetUsers",
		func(ctx context.Context, client *sdk.APIClient, id int64, _ any) (*http.Response, error) {
			_, hresp, err := client.UsersAPI.DeleteUser(ctx, id).Execute()
			return hresp, err
		},
		testhelpers.WithNameMethod("GetUsernameOk"),
		testhelpers.WithFilter(func(_ context.Context, _ *sdk.APIClient, item any) (bool, string, error) {
			user := item.(sdk.User)
			email, ok := user.GetEmailOk()
			if !ok || email == nil || !isSweepableEmail(*email) {
				return false, "email", nil
			}

			return true, "", nil
		}),
	)
}
