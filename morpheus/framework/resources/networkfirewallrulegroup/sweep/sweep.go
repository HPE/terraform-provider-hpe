// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

//go:build sweep

package sweep

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	testsweep "github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/sweep"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/getsafe"
)

const sweeperName = "hpe_morpheus_network_firewall_rule_group"

// ruleGroupItem represents a firewall rule group item for sweep operations.
type ruleGroupItem struct {
	Id              int64
	Name            string
	NetworkServerID int64
}

func init() {
	testsweep.RegisterTypedAPISweeper(
		sweeperName,
		// List network firewall rule group resources by iterating network servers.
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

			for _, ns := range serversResp.NetworkServers {
				nsID, ok := getsafe.GetOk(ns.Id)
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

				for _, rg := range listResp.RuleGroups {
					id, idOk := getsafe.GetOk(rg.Id)
					name, nameOk := getsafe.GetOk(rg.Name)

					if !idOk || id == nil || !nameOk || name == nil {
						continue
					}

					allGroups = append(allGroups, ruleGroupItem{
						Id:              *id,
						Name:            *name,
						NetworkServerID: *nsID,
					})
				}
			}

			return allGroups, &http.Response{StatusCode: http.StatusOK}, nil
		},
		// Is this a test network firewall rule group?
		func(item ruleGroupItem) bool {
			return strings.HasPrefix(item.Name, testsweep.TestResourcePrefix)
		},
		// Delete the test network firewall rule group.
		func(
			ctx context.Context,
			client *sdk.APIClient,
			item ruleGroupItem,
		) (*http.Response, error) {
			if item.NetworkServerID == 0 {
				return nil, fmt.Errorf("could not get network server ID")
			}

			_, hresp, err := client.NetworksAPI.
				DeleteNetworkFirewallRuleGroup(ctx, item.Id, item.NetworkServerID).
				Execute()

			return hresp, err
		},
	)
}
