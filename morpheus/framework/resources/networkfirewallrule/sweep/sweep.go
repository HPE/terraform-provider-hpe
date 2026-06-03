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

			for _, server := range serverResp.NetworkServers {
				serverID, ok := getsafe.GetSafeOk(server.Id)
				if !ok || serverID == nil {
					continue
				}

				ruleResp, _, err := client.NetworksAPI.
					GetNetworkFirewallRules(ctx, *serverID).Execute()
				if err != nil || ruleResp == nil {
					continue
				}

				for _, rule := range ruleResp.Rules {
					id, ok := getsafe.GetSafeOk(rule.Id)
					if !ok || id == nil {
						continue
					}

					name, ok := getsafe.GetSafeOk(rule.Name)
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
