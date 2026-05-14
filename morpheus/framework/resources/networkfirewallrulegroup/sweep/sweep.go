// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package sweep

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	testsweep "github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/sweep"
)

// ruleGroupItem represents a firewall rule group item returned by the list API.
type ruleGroupItem struct {
	Id            *int64  `json:"id"`
	Name          *string `json:"name"`
	NetworkServer *struct {
		Id *float64 `json:"id"`
	} `json:"networkServer"`
}

func init() {
	testsweep.RegisterTypedAPISweeper(
		"hpe_morpheus_network_firewall_rule_group",
		// List all firewall rule groups across all network servers.
		func(ctx context.Context, client *sdk.APIClient) (
			[]ruleGroupItem,
			*http.Response,
			error,
		) {
			serversResp, hresp, err := client.NetworksAPI.
				ListNetworkServers(ctx).Execute()
			if err != nil {
				return nil, hresp, err
			}

			if serversResp == nil {
				return nil, hresp, err
			}

			var allGroups []ruleGroupItem

			for _, ns := range serversResp.GetNetworkServers() {
				nsID, ok := ns.GetIdOk()
				if !ok || nsID == nil {
					continue
				}

				listResp, listHresp, err := client.NetworksAPI.
					GetNetworkFirewallRuleGroups(ctx, *nsID).Execute()
				if err != nil {
					log.Printf(
						"[WARN] Failed to list firewall rule groups for server %d: %v",
						*nsID, err,
					)

					continue
				}

				if listHresp == nil || listHresp.StatusCode != http.StatusOK {
					continue
				}

				raw := listResp.GetRuleGroups()
				if raw == nil {
					continue
				}

				encoded, err := json.Marshal(raw)
				if err != nil {
					continue
				}

				var items []ruleGroupItem
				if err := json.Unmarshal(encoded, &items); err != nil {
					continue
				}

				allGroups = append(allGroups, items...)
			}

			return allGroups, &http.Response{StatusCode: http.StatusOK}, nil
		},
		// Is this a test firewall rule group?
		func(item ruleGroupItem) bool {
			if item.Name == nil {
				return false
			}

			return strings.HasPrefix(*item.Name, testsweep.TestResourcePrefix)
		},
		// Delete the test firewall rule group.
		func(
			ctx context.Context,
			client *sdk.APIClient,
			item ruleGroupItem,
		) (*http.Response, error) {
			if item.Id == nil {
				return nil, fmt.Errorf("could not get ID")
			}

			if item.NetworkServer == nil || item.NetworkServer.Id == nil {
				return nil, fmt.Errorf("could not get network server ID")
			}

			serverID := int64(*item.NetworkServer.Id)

			_, hresp, err := client.NetworksAPI.
				DeleteNetworkFirewallRuleGroup(ctx, *item.Id, serverID).
				Execute()

			return hresp, err
		},
	)
}
