// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package policy

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	morpheuserrors "github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/errors"
)

// Policies whose name begins with this string will be eligible for deletion
const TestPolicyPrefix = "TestAccMorpheusPolicy"

// SweepPolicies cleans up test policies - exported for use by global sweep
func SweepPolicies(client *sdk.APIClient) error {
	ctx := context.Background()

	if client == nil {
		log.Printf("[INFO] No client provided, skipping policy sweep")
		return nil
	}

	policies, hresp, err := client.PoliciesAPI.ListPolicies(ctx).
		Phrase(TestPolicyPrefix).Execute()
	if err != nil {
		// Handle 404 and 403 as "no matches found" rather than an error
		if hresp != nil && (hresp.StatusCode == http.StatusNotFound || hresp.StatusCode == http.StatusForbidden) {
			log.Printf("[INFO] No policies found matching prefix (status %d): %s", hresp.StatusCode, TestPolicyPrefix)

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

		if !strings.HasPrefix(*name, TestPolicyPrefix) {
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
