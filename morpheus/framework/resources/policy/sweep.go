// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package policy

import (
	"context"
	"net/http"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

// Policies whose name begins with this string will be eligible for deletion
const testPolicyPrefix = "TestAccMorpheusPolicy"

func init() {
	testhelpers.RegisterAPISweeper(
		"hpe_morpheus_policy",
		testPolicyPrefix,
		func(ctx context.Context, client *sdk.APIClient, prefix string) (any, *http.Response, error) {
			return client.PoliciesAPI.ListPolicies(ctx).Phrase(prefix).Execute()
		},
		"GetPolicies",
		func(ctx context.Context, client *sdk.APIClient, id int64, _ any) (*http.Response, error) {
			_, hresp, err := client.PoliciesAPI.RemovePolicies(ctx, id).Execute()
			return hresp, err
		},
		testhelpers.WithIgnoreListStatuses(http.StatusNotFound, http.StatusForbidden),
	)
}
