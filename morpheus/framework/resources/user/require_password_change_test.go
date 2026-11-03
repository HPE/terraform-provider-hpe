// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package user_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

// userID captures the user's id from state so a later step can act on the
// user behind Terraform's back.
type userID struct{ id int64 }

func (u *userID) record(addr string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[addr]
		if !ok {
			return fmt.Errorf("%s not in state", addr)
		}

		id, err := strconv.ParseInt(rs.Primary.ID, 10, 64)
		if err != nil {
			return fmt.Errorf("%s: id %q: %w", addr, rs.Primary.ID, err)
		}

		u.id = id

		return nil
	}
}

// clearPasswordExpired stands in for the user changing their password in the
// Morpheus UI, which clears the flag: it clears it through the API. The
// recorder is reset afterwards so the step's checks see only what the
// provider sends.
func (u *userID) clearPasswordExpired(t *testing.T, rec *testhelpers.RequestRecorder) func() {
	t.Helper()

	return func() {
		defer rec.Reset()

		ctx := context.Background()

		client, err := testhelpers.NewClientForServer(ctx, "")
		if err != nil {
			t.Fatal(err)
		}

		_, _, err = client.UsersAPI.UpdateUser(ctx, u.id).UpdateUserRequest(sdk.UpdateUserRequest{
			User: sdk.UpdateUserRequestUser{PasswordExpired: sdk.PtrBool(false)},
		}).Execute()
		if err != nil {
			t.Fatalf("clearing passwordExpired on user %d: %v", u.id, err)
		}
	}
}

// TestAccMorpheusUserResourceRequirePasswordChangeOk covers the acceptance
// criteria of MORPH-18374:
//
//  1. require_password_change = true on create -> the user is created with
//     passwordExpired: true and password_expired reads back true (the flag
//     redirects the next UI sign-in; API sign-in with the password still
//     works, which the password checks confirm);
//  2. once the user has changed their password (the flag is false again) the
//     plan is empty -- the attribute is not read back, so it does not fight
//     the user's change;
//  3. a new password_wo while require_password_change is true re-sends
//     passwordExpired: true, and the flag is set again;
//  4. with require_password_change false (or unset) passwordExpired is never
//     sent, so a password set by the administrator is not a temporary one;
//  5. turning require_password_change on for an existing user, without a
//     password change, expires their current password (the flag alone is
//     sent); the user's own change is then left alone as in 2; and turning
//     it off sends nothing (the review finding on the original change).
//
// The request recorder is process-wide, so the test does not run in parallel.
func TestAccMorpheusUserResourceRequirePasswordChangeOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	name := acctest.RandomWithPrefix(t.Name())
	email := name + "@testacc.com"
	addr := "hpe_morpheus_user.wo"
	rec := testhelpers.RecordRequests(t)
	uid := &userID{}

	update := resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{
		plancheck.ExpectResourceAction(addr, plancheck.ResourceActionUpdate),
	}}

	const (
		required    = "\trequire_password_change = true\n"
		notRequired = "\trequire_password_change = false\n"
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				// 1. created with the password marked expired.
				Config: userConfig(name, email, "Secret123!", required),
				Check: resource.ComposeAggregateTestCheckFunc(
					uid.record(addr),
					resource.TestCheckResourceAttr(addr, "require_password_change", "true"),
					resource.TestCheckResourceAttr(addr, "password_expired", "true"),
					rec.ExpectSent("/api/users", "user.passwordExpired", true),
					// The flag only redirects the UI sign-in; the password is live.
					checkPassword(t, name, "Secret123!", true),
				),
			},
			{
				Config:   userConfig(name, email, "Secret123!", required),
				PlanOnly: true,
			},
			{
				// 2. the user changed their password: the flag is false, the
				// plan is empty, nothing is sent.
				PreConfig: uid.clearPasswordExpired(t, rec),
				Config:    userConfig(name, email, "Secret123!", required),
				ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(addr, plancheck.ResourceActionNoop),
				}},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(addr, "password_expired", "false"),
					rec.ExpectSent("/api/users", "user.passwordExpired", false),
					checkPassword(t, name, "Secret123!", true),
				),
			},
			{
				// 3. a new password from the administrator is temporary again.
				Config:           userConfig(name, email, "Secret456!", required),
				ConfigPlanChecks: update,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(addr, "password_expired", "true"),
					rec.ExpectSent("/api/users", "user.passwordExpired", true),
					checkPassword(t, name, "Secret456!", true),
				),
			},
			{
				// 4a. require_password_change = false: a new password is sent
				// without passwordExpired, and the user can sign in with it.
				PreConfig:        uid.clearPasswordExpired(t, rec),
				Config:           userConfig(name, email, "Secret789!", notRequired),
				ConfigPlanChecks: update,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(addr, "password_expired", "false"),
					rec.ExpectSentKeys("/api/users", map[string]bool{
						"user.password":        true,
						"user.passwordExpired": false,
					}),
					checkPassword(t, name, "Secret789!", true),
				),
			},
			{
				// 4b. unset: an unrelated update does not send it either.
				Config:           userConfig(name, name+"-2@testacc.com", "Secret789!", ""),
				ConfigPlanChecks: update,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckNoResourceAttr(addr, "require_password_change"),
					resource.TestCheckResourceAttr(addr, "password_expired", "false"),
					rec.ExpectSent("/api/users", "user.passwordExpired", false),
				),
			},
			{
				Config:   userConfig(name, name+"-2@testacc.com", "Secret789!", ""),
				PlanOnly: true,
			},
			{
				// 5a. turned on for an existing user, nothing else changed: the
				// flag alone is sent and the current password is now expired.
				Config:           userConfig(name, name+"-2@testacc.com", "Secret789!", required),
				ConfigPlanChecks: update,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(addr, "require_password_change", "true"),
					resource.TestCheckResourceAttr(addr, "password_expired", "true"),
					rec.ExpectSentKeys("/api/users", map[string]bool{
						"user.passwordExpired": true,
						"user.password":        false,
					}),
					checkPassword(t, name, "Secret789!", true),
				),
			},
			{
				// 5b. the user changed their password: left alone, as in 2.
				PreConfig: uid.clearPasswordExpired(t, rec),
				Config:    userConfig(name, name+"-2@testacc.com", "Secret789!", required),
				ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(addr, plancheck.ResourceActionNoop),
				}},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(addr, "password_expired", "false"),
					rec.ExpectSent("/api/users", "user.passwordExpired", false),
				),
			},
			{
				// 5c. turned off: an update that sends nothing about the flag.
				Config:           userConfig(name, name+"-2@testacc.com", "Secret789!", notRequired),
				ConfigPlanChecks: update,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(addr, "require_password_change", "false"),
					resource.TestCheckResourceAttr(addr, "password_expired", "false"),
					rec.ExpectSent("/api/users", "user.passwordExpired", false),
					checkPassword(t, name, "Secret789!", true),
				),
			},
			{
				Config:   userConfig(name, name+"-2@testacc.com", "Secret789!", notRequired),
				PlanOnly: true,
			},
		},
	})
}
