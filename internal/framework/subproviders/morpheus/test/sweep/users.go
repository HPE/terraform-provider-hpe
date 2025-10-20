// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

// Package sweep allows deletion of dangling test resources
package sweep

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

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

func Users() {
	resource.AddTestSweepers(
		"hpe_morpheus_user",
		&resource.Sweeper{
			Name: "hpe_morpheus_user",
			F:    testSweepMorpheusUsers,
		})
}

func testSweepMorpheusUsers(_ string) error {
	ctx := context.Background()

	client, err := NewSweepClient(ctx)
	if err != nil {
		return err
	}

	users, hresp, err := client.UsersAPI.ListUsers(ctx).
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
		_, hresp, err := client.UsersAPI.DeleteUser(ctx, *id).Execute()
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
