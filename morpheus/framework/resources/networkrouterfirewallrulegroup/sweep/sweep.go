// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

//go:build sweep

package sweep

import (
	"context"
	"net/http"
	"strings"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	testsweep "github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/sweep"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/getsafe"
)

const sweeperName = "hpe_morpheus_network_router_firewall_rule_group"

type ruleGroupSweeperItem struct {
	routerID  int64
	ruleGroup sdk.GetNetworkRouterFirewallRuleGroups200ResponseRuleGroupsInner
}

func init() {
	testsweep.RegisterTypedAPISweeper(
		sweeperName,
		// List network router firewall rule groups by iterating routers.
		func(ctx context.Context, client *sdk.APIClient) ([]ruleGroupSweeperItem, *http.Response, error) {
			routersResp, routersHTTPResp, err := client.NetworksAPI.GetNetworkRouters(ctx).Execute()
			if err != nil || routersResp == nil {
				return nil, routersHTTPResp, err
			}

			items := make([]ruleGroupSweeperItem, 0)

			for _, router := range routersResp.NetworkRouters {
				routerID, ok := getsafe.GetOk(router.Id)
				if !ok || routerID == nil {
					continue
				}

				groupsResp, _, listErr := client.NetworksAPI.
					GetNetworkRouterFirewallRuleGroups(ctx, *routerID).Execute()
				if listErr != nil || groupsResp == nil {
					continue
				}

				for _, rg := range groupsResp.RuleGroups {
					items = append(items, ruleGroupSweeperItem{routerID: *routerID, ruleGroup: rg})
				}
			}

			return items, routersHTTPResp, nil
		},
		// Is this a test network router firewall rule group?
		func(item ruleGroupSweeperItem) bool {
			name, ok := getsafe.GetOk(item.ruleGroup.Name)
			if !ok || name == nil {
				return false
			}

			return strings.HasPrefix(*name, testsweep.TestResourcePrefix)
		},
		// Delete the test network router firewall rule group.
		func(ctx context.Context, client *sdk.APIClient, item ruleGroupSweeperItem) (*http.Response, error) {
			id, ok := getsafe.GetOk(item.ruleGroup.Id)
			if !ok || id == nil {
				return &http.Response{StatusCode: http.StatusOK}, nil
			}

			_, hresp, err := client.NetworksAPI.
				DeleteNetworkRouterFirewallRuleGroup(ctx, *id, item.routerID).Execute()
			if err != nil && hresp != nil && hresp.StatusCode == http.StatusNotFound {
				return &http.Response{StatusCode: http.StatusOK}, nil
			}

			return hresp, err
		},
	)
}
