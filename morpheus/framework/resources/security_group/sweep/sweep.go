// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

//go:build sweep

package sweep

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	testsweep "github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/sweep"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/getsafe"
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

			return getsafe.GetSafe(&resp.SecurityGroups), hresp, err
		},
		// Is this a test security group?
		func(item sdk.ListSecurityGroups200ResponseAllOfSecurityGroupsInner) bool {
			name, ok := getsafe.GetSafeOk(item.Name)
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
			id, ok := getsafe.GetSafeOk(item.Id)
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
		testsweep.WithDependencies[sdk.ListSecurityGroups200ResponseAllOfSecurityGroupsInner](
			"hpe_morpheus_security_group_rule",
		),
	)
}
