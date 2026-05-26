// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkfirewallrule_test

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	testsweep "github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/sweep"
)

const sweeperName = "hpe_morpheus_network_firewall_rule"

// networkFirewallRuleSweepItem pairs a firewall rule with its parent network server ID.
type networkFirewallRuleSweepItem struct {
	serverID int64
	id       int64
	name     string
}

func init() {
	testsweep.RegisterTypedAPISweeper(
		sweeperName,
		// List network firewall rule resources.
		func(ctx context.Context, client *sdk.APIClient) (
			[]networkFirewallRuleSweepItem,
			*http.Response,
			error,
		) {
			serverResp, hresp, err := client.NetworksAPI.ListNetworkServers(ctx).Execute()
			if err != nil {
				return nil, hresp, err
			}

			var items []networkFirewallRuleSweepItem

			for _, server := range serverResp.GetNetworkServers() {
				serverID, ok := server.GetIdOk()
				if !ok || serverID == nil {
					continue
				}

				ruleResp, _, err := client.NetworksAPI.
					GetNetworkFirewallRules(ctx, *serverID).Execute()
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

					items = append(items, networkFirewallRuleSweepItem{
						serverID: *serverID,
						id:       *id,
						name:     *name,
					})
				}
			}

			return items, hresp, nil
		},
		// Is this a test network firewall rule?
		func(item networkFirewallRuleSweepItem) bool {
			return strings.HasPrefix(item.name, testsweep.TestResourcePrefix)
		},
		// Delete the test network firewall rule.
		func(
			ctx context.Context,
			client *sdk.APIClient,
			item networkFirewallRuleSweepItem,
		) (*http.Response, error) {
			if item.id == 0 {
				return nil, fmt.Errorf("could not get ID")
			}

			_, hresp, err := client.NetworksAPI.
				DeleteNetworkFirewallRule(ctx, item.id, item.serverID).Execute()

			return hresp, err
		},
		testsweep.WithIgnoreListStatuses[networkFirewallRuleSweepItem](
			http.StatusNotFound,
			http.StatusForbidden,
		),
	)
}
