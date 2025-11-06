// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package policy

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/clientfactory"
	morpheuserrors "github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/errors"
)

// Policies whose name begins with this string will be eligible for deletion
const TestPolicyPrefix = "TestAccMorpheusPolicy"

// newSweepClient creates a client for sweeping test resources
func newSweepClient(ctx context.Context) (*sdk.APIClient, error) {
	var username, password string

	url, ok := os.LookupEnv("TF_VAR_testacc_morpheus_url")
	if !ok {
		return nil, errors.New("TF_VAR_testacc_morpheus_url not set")
	}

	token, ok := os.LookupEnv("TF_VAR_testacc_morpheus_access_token")
	if !ok {
		username, ok = os.LookupEnv("TF_VAR_testacc_morpheus_username")
		if !ok {
			return nil, errors.New(
				"one of TF_VAR_testacc_morpheus_access_token or " +
					"TF_VAR_testacc_morpheus_username must be set",
			)
		}

		password, ok = os.LookupEnv("TF_VAR_testacc_morpheus_password")
		if !ok {
			return nil, errors.New(
				"one of TF_VAR_testacc_morpheus_access_token or " +
					"TF_VAR_testacc_morpheus_password must be set",
			)
		}
	}

	// If set to any value, use insecure
	_, insecure := os.LookupEnv("TF_VAR_testacc_morpheus_insecure")
	var opts []clientfactory.ClientOption
	if insecure {
		opts = append(opts, clientfactory.WithInsecureTLS())
	}

	client := clientfactory.NewAPIClient(
		ctx,
		url,
		username,
		password,
		token,
		opts...,
	)

	return client, nil
}

// SweepPolicies cleans up test policies - exported for use by global sweep
func SweepPolicies(_ string) error {
	ctx := context.Background()

	client, err := newSweepClient(ctx)
	if err != nil {
		// If we can't create a client (e.g., env vars not set), just log and continue
		log.Printf("[INFO] Cannot create sweep client, skipping policy sweep: %v", err)

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
