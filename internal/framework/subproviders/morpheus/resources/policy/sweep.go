// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package policy

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

// Policies whose name begins with this string will be eligible for deletion
const testPolicyPrefix = "TestAccMorpheusPolicy"

// policySweeper handles cleanup of test policies
type policySweeper struct {
	client *sdk.APIClient
}

// NewPolicySweeper creates and registers a policy sweeper
func NewPolicySweeper(client *sdk.APIClient) {
	s := &policySweeper{
		client: client,
	}

	resource.AddTestSweepers(
		"hpe_morpheus_policy",
		&resource.Sweeper{
			Name: "hpe_morpheus_policy",
			F:    s.Sweep,
		})
}

// Sweep cleans up test policies
func (s *policySweeper) Sweep(_ string) error {
	ctx := context.Background()

	if s.client == nil {
		log.Printf("[INFO] No client provided, skipping policy sweep")

		return nil
	}

	policies, hresp, err := s.client.PoliciesAPI.ListPolicies(ctx).
		Phrase(testPolicyPrefix).Execute()
	if err != nil {
		// Handle 404 and 403 as "no matches found" rather than an error
		if hresp != nil && (hresp.StatusCode == http.StatusNotFound || hresp.StatusCode == http.StatusForbidden) {
			log.Printf("[INFO] No policies found matching prefix (status %d): %s", hresp.StatusCode, testPolicyPrefix)

			return nil
		}

		return fmt.Errorf("failed to list policies: %s", errors.ErrMsg(err, hresp))
	}
	if hresp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to list policies: %s", errors.ErrMsg(err, hresp))
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
		_, hresp, err := s.client.PoliciesAPI.RemovePolicies(ctx, *id).Execute()
		if err != nil || hresp.StatusCode != http.StatusOK {
			errMsg := fmt.Sprintf(
				"failed to delete policy %s (id: %d): %s",
				*name, *id, errors.ErrMsg(err, hresp),
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
