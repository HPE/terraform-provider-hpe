// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package security_group_rule_test

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	testsweep "github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/sweep"
)

const sweeperName = "hpe_morpheus_security_group_rule"

// securityGroupRuleSweepItem pairs a rule with its parent security group ID.
// RemoveSecurityGroupRules takes the security group ID as float32 (SDK quirk).
type securityGroupRuleSweepItem struct {
	securityGroupID float32
	id              int64
	name            string
}

func init() {
	testsweep.RegisterTypedAPISweeper(
		sweeperName,
		// List security group rule resources.
		func(ctx context.Context, client *sdk.APIClient) (
			[]securityGroupRuleSweepItem,
			*http.Response,
			error,
		) {
			sgResp, hresp, err := client.SecurityGroupsAPI.ListSecurityGroups(ctx).Execute()
			if err != nil {
				return nil, hresp, err
			}

			var items []securityGroupRuleSweepItem

			for _, sg := range sgResp.GetSecurityGroups() {
				sgID, ok := sg.GetIdOk()
				if !ok || sgID == nil {
					continue
				}

				ruleResp, _, err := client.SecurityGroupsAPI.
					ListSecurityGroupRules(ctx, *sgID).Execute()
				if err != nil || ruleResp == nil {
					continue
				}

				for _, rule := range ruleResp.GetRules() {
					id, ok := rule.GetIdOk()
					if !ok || id == nil {
						continue
					}

					name, ok := rule.GetNameOk()
					if !ok || name == nil {
						continue
					}

					items = append(items, securityGroupRuleSweepItem{
						securityGroupID: float32(*sgID),
						id:              *id,
						name:            *name,
					})
				}
			}

			return items, hresp, nil
		},
		// Is this a test security group rule?
		func(item securityGroupRuleSweepItem) bool {
			return strings.HasPrefix(item.name, testsweep.TestResourcePrefix)
		},
		// Delete the test security group rule.
		func(
			ctx context.Context,
			client *sdk.APIClient,
			item securityGroupRuleSweepItem,
		) (*http.Response, error) {
			if item.id == 0 {
				return nil, fmt.Errorf("could not get ID")
			}

			_, hresp, err := client.SecurityGroupsAPI.
				RemoveSecurityGroupRules(ctx, item.id, item.securityGroupID).Execute()

			return hresp, err
		},
		testsweep.WithIgnoreListStatuses[securityGroupRuleSweepItem](
			http.StatusNotFound,
			http.StatusForbidden,
		),
	)
}
