// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package policy_test

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	morpheuserrors "github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/errors"
	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/test/sweep"
)

// Policies whose name begins with this string will be eligible for deletion
const testPolicyPrefix = "TestAccMorpheusPolicy"

// sweepPolicies cleans up test policies after tests complete
func sweepPolicies() error {
	ctx := context.Background()

	client, err := sweep.NewSweepClient(ctx)
	if err != nil {
		// If we can't create a client (e.g., env vars not set), just log and continue
		log.Printf("[INFO] Cannot create sweep client, skipping policy sweep: %v", err)

		return nil
	}

	policies, hresp, err := client.PoliciesAPI.ListPolicies(ctx).
		Phrase(testPolicyPrefix).Execute()
	if err != nil {
		// Handle 404 and 403 as "no matches found" rather than an error
		if hresp != nil && (hresp.StatusCode == http.StatusNotFound || hresp.StatusCode == http.StatusForbidden) {
			log.Printf("[INFO] No policies found matching prefix (status %d): %s", hresp.StatusCode, testPolicyPrefix)

			return nil
		}

		return fmt.Errorf("failed to list policies: %s", morpheuserrors.ErrMsg(err, hresp))
	}
	if hresp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to list policies: %s", morpheuserrors.ErrMsg(err, hresp))
	}

	policyList := policies.GetPolicies()
	var sweptCount int
	var sweepErrors []string

	for _, policy := range policyList {
		name, ok := policy.GetNameOk()
		if !ok || name == nil {
			continue
		}

		if !strings.HasPrefix(*name, testPolicyPrefix) {
			log.Printf("[INFO] Skipping policy (name): %s", *name)

			continue
		}

		id, ok := policy.GetIdOk()
		if !ok || id == nil {
			log.Printf("[INFO] Skipping policy (id): %s", *name)

			continue
		}

		log.Printf("[INFO] Sweeping policy: %s (id: %d)", *name, *id)

		// Delete the policy
		_, hresp, err := client.PoliciesAPI.RemovePolicies(ctx, *id).Execute()
		if err != nil || hresp.StatusCode != http.StatusOK {
			errMsg := fmt.Sprintf(
				"failed to delete policy %s (id: %d): %s",
				*name, *id, morpheuserrors.ErrMsg(err, hresp),
			)
			log.Printf("[ERROR] %s", errMsg)
			sweepErrors = append(sweepErrors, errMsg)

			continue
		}

		sweptCount++
	}

	log.Printf(
		"[INFO] Policy sweep completed. Policies swept: %d, errors: %d",
		sweptCount, len(sweepErrors),
	)

	if len(sweepErrors) > 0 {
		return fmt.Errorf("%s", strings.Join(sweepErrors, "\n"))
	}

	return nil
}
