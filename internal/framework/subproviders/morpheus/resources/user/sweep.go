// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package user

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/errors"
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

// userSweeper handles cleanup of test users
type userSweeper struct {
	client *sdk.APIClient
}

// NewUserSweeper creates and registers a user sweeper
func NewUserSweeper(client *sdk.APIClient) {
	s := &userSweeper{
		client: client,
	}

	resource.AddTestSweepers(
		"hpe_morpheus_user",
		&resource.Sweeper{
			Name: "hpe_morpheus_user",
			F:    s.Sweep,
		})
}

// Sweep cleans up test users
func (s *userSweeper) Sweep(_ string) error {
	ctx := context.Background()

	if s.client == nil {
		log.Printf("[INFO] No client provided, skipping user sweep")

		return nil
	}

	users, hresp, err := s.client.UsersAPI.ListUsers(ctx).
		Phrase(testUserPrefix).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to list users: %s", errors.ErrMsg(err, hresp))
	}

	userList := users.GetUsers()
	var sweptCount int
	var sweepErrors []string

	for _, user := range userList {
		username, ok := user.GetUsernameOk()
		if !ok || username == nil {
			continue
		}

		if !strings.HasPrefix(*username, testUserPrefix) {
			log.Printf("[INFO] Skipping user (name): %s", *username)

			continue
		}

		email, ok := user.GetEmailOk()
		if !ok || email == nil || !isSweepableEmail(*email) {
			log.Printf("[INFO] Skipping user (email): %s", *username)

			continue
		}

		id, ok := user.GetIdOk()
		if !ok || id == nil {
			log.Printf("[INFO] Skipping user (id): %s", *username)

			continue
		}

		log.Printf("[INFO] Sweeping user: %s (id: %d)", *username, *id)

		// Delete the user
		_, hresp, err := s.client.UsersAPI.DeleteUser(ctx, *id).Execute()
		if err != nil || hresp.StatusCode != http.StatusOK {
			errMsg := fmt.Sprintf(
				"failed to delete user %s (id: %d): %s",
				*username, *id, errors.ErrMsg(err, hresp),
			)
			log.Printf("[ERROR] %s", errMsg)
			sweepErrors = append(sweepErrors, errMsg)

			continue
		}

		sweptCount++
	}

	log.Printf(
		"[INFO] User sweep completed. Users swept: %d, errors: %d",
		sweptCount, len(sweepErrors),
	)

	if len(sweepErrors) > 0 {
		return fmt.Errorf("%s", strings.Join(sweepErrors, "\n"))
	}

	return nil
}
