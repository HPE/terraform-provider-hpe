// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package security_group_test

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	testsweep "github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/sweep"
)

const sweeperName = "hpe_morpheus_security_group"

func init() {
	testsweep.RegisterTypedAPISweeper(
		sweeperName,
		// List security group resources.
		func(ctx context.Context, client *sdk.APIClient) (
			[]sdk.ListSecurityGroups200ResponseAllOfSecurityGroupsInner,
			*http.Response,
			error,
		) {
			resp, hresp, err := client.SecurityGroupsAPI.ListSecurityGroups(ctx).Execute()
			if resp == nil {
				return nil, hresp, err
			}

			return resp.GetSecurityGroups(), hresp, err
		},
		// Is this a test security group?
		func(item sdk.ListSecurityGroups200ResponseAllOfSecurityGroupsInner) bool {
			name, ok := item.GetNameOk()
			if !ok || name == nil {
				return false
			}

			return strings.HasPrefix(*name, testsweep.TestResourcePrefix)
		},
		// Delete the test security group.
		func(
			ctx context.Context,
			client *sdk.APIClient,
			item sdk.ListSecurityGroups200ResponseAllOfSecurityGroupsInner,
		) (*http.Response, error) {
			id, ok := item.GetIdOk()
			if !ok || id == nil {
				return nil, fmt.Errorf("could not get ID")
			}

			_, hresp, err := client.SecurityGroupsAPI.RemoveSecurityGroups(ctx, *id).Execute()

			return hresp, err
		},
		testsweep.WithIgnoreListStatuses[sdk.ListSecurityGroups200ResponseAllOfSecurityGroupsInner](
			http.StatusNotFound,
			http.StatusForbidden,
		),
	)
}
