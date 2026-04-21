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
	testhelpers.RegisterTypedAPISweeper(
		"hpe_morpheus_policy",
		testPolicyPrefix,
		func(ctx context.Context, client *sdk.APIClient, prefix string) ([]sdk.ListPolicies200ResponseAllOfPoliciesInner, *http.Response, error) {
			resp, hresp, err := client.PoliciesAPI.ListPolicies(ctx).Phrase(prefix).Execute()
			if resp == nil {
				return nil, hresp, err
			}

			return resp.GetPolicies(), hresp, err
		},
		func(item sdk.ListPolicies200ResponseAllOfPoliciesInner) (string, bool) {
			name, ok := item.GetNameOk()
			if !ok || name == nil {
				return "", false
			}

			return *name, true
		},
		func(item sdk.ListPolicies200ResponseAllOfPoliciesInner) (int64, bool) {
			id, ok := item.GetIdOk()
			if !ok || id == nil {
				return 0, false
			}

			return *id, true
		},
		func(ctx context.Context, client *sdk.APIClient, id int64, _ sdk.ListPolicies200ResponseAllOfPoliciesInner) (*http.Response, error) {
			_, hresp, err := client.PoliciesAPI.RemovePolicies(ctx, id).Execute()
			return hresp, err
		},
		testhelpers.WithIgnoreListStatuses[sdk.ListPolicies200ResponseAllOfPoliciesInner](
			http.StatusNotFound,
			http.StatusForbidden,
		),
	)
}
