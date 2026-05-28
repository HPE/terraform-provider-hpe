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
)

const sweeperName = "hpe_morpheus_network_dhcp_server"

func init() {
	testsweep.RegisterTypedAPISweeper(
		sweeperName,
		// List DHCP server resources by iterating network servers.
		func(ctx context.Context, client *sdk.APIClient) (
			[]sdk.GetNetworkDhcpServers200ResponseAllOfNetworkDhcpServersInner,
			*http.Response,
			error,
		) {
			serversResp, hresp, err := client.NetworksAPI.ListNetworkServers(ctx).Execute()
			if err != nil {
				return nil, hresp, err
			}

			if serversResp == nil {
				return nil, hresp, err
			}

			var allDhcpServers []sdk.GetNetworkDhcpServers200ResponseAllOfNetworkDhcpServersInner

			for _, ns := range serversResp.GetNetworkServers() {
				nsID, ok := ns.GetIdOk()
				if !ok || nsID == nil {
					continue
				}

				dhcpResp, dhcpHresp, err := client.NetworksAPI.
					GetNetworkDhcpServers(ctx, *nsID).Execute()
				if err != nil {
					log.Printf(
						"[WARN] Failed to list DHCP servers for network server %d: %v",
						*nsID, err,
					)

					continue
				}

				if dhcpHresp == nil || dhcpHresp.StatusCode != http.StatusOK {
					continue
				}

				allDhcpServers = append(allDhcpServers, dhcpResp.GetNetworkDhcpServers()...)
			}

			return allDhcpServers, &http.Response{StatusCode: http.StatusOK}, nil
		},
		// Is this a test DHCP server?
		func(item sdk.GetNetworkDhcpServers200ResponseAllOfNetworkDhcpServersInner) bool {
			name, ok := item.GetNameOk()
			if !ok || name == nil {
				return false
			}

			return strings.HasPrefix(*name, testsweep.TestResourcePrefix)
		},
		// Delete the test DHCP server.
		func(
			ctx context.Context,
			client *sdk.APIClient,
			item sdk.GetNetworkDhcpServers200ResponseAllOfNetworkDhcpServersInner,
		) (*http.Response, error) {
			id, ok := item.GetIdOk()
			if !ok || id == nil {
				return nil, fmt.Errorf("could not get ID")
			}

			ns, ok := item.GetNetworkServerOk()
			if !ok || ns == nil {
				return nil, fmt.Errorf("could not get network server")
			}

			nsID, ok := ns.GetIdOk()
			if !ok || nsID == nil {
				return nil, fmt.Errorf("could not get network server ID")
			}

			_, hresp, err := client.NetworksAPI.
				DeleteNetworkDhcpServer(ctx, *id, *nsID).Execute()

			return hresp, err
		},
	)
}
