// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

//go:build sweep

package sweep

import (
	"context"
	"net/http"
	"strings"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	testsweep "github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/sweep"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/getsafe"
)

const sweeperName = "hpe_morpheus_network_router_firewall_rule"

type firewallRuleSweeperItem struct {
	routerID int64
	rule     sdk.GetNetworkRoutersFirewallRules200ResponseRulesInner
}

func init() {
	testsweep.RegisterTypedAPISweeper(
		sweeperName,
		// List network router firewall rule resources by iterating routers.
		func(ctx context.Context, client *sdk.APIClient) ([]firewallRuleSweeperItem, *http.Response, error) {
			routersResp, routersHTTPResp, err := client.NetworksAPI.GetNetworkRouters(ctx).Execute()
			if err != nil || routersResp == nil {
				return nil, routersHTTPResp, err
			}

			items := make([]firewallRuleSweeperItem, 0)

			for _, router := range routersResp.NetworkRouters {
				routerID, ok := getsafe.GetOk(router.Id)
				if !ok || routerID == nil {
					continue
				}

				rulesResp, _, listErr := client.NetworksAPI.GetNetworkRoutersFirewallRules(ctx, *routerID).Execute()
				if listErr != nil || rulesResp == nil {
					continue
				}

				for _, rule := range rulesResp.Rules {
					items = append(items, firewallRuleSweeperItem{routerID: *routerID, rule: rule})
				}
			}

			return items, routersHTTPResp, nil
		},
		// Is this a test network router firewall rule?
		func(item firewallRuleSweeperItem) bool {
			name, ok := getsafe.GetOk(item.rule.Name)
			if !ok || name == nil {
				return false
			}

			return strings.HasPrefix(*name, testsweep.TestResourcePrefix)
		},
		// Delete the test network router firewall rule.
		func(ctx context.Context, client *sdk.APIClient, item firewallRuleSweeperItem) (*http.Response, error) {
			id, ok := getsafe.GetOk(item.rule.Id)
			if !ok || id == nil {
				return &http.Response{StatusCode: http.StatusOK}, nil
			}

			_, hresp, err := client.NetworksAPI.DeleteNetworkRouterFirewallRule(ctx, *id, item.routerID).Execute()
			if err != nil && hresp != nil && hresp.StatusCode == http.StatusNotFound {
				return &http.Response{StatusCode: http.StatusOK}, nil
			}

			return hresp, err
		},
	)
}
