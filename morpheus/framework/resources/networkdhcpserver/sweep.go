// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package networkdhcpserver

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
)

const testDhcpServerPrefix = "TestAccMorpheusNetworkDhcpServer"

type dhcpServerSweeper struct {
	client   *sdk.APIClient
	serverID float32
}

// NewNetworkDhcpServerSweeper creates and registers a DHCP server sweeper.
// The serverID identifies the network server whose DHCP servers will be swept.
func NewNetworkDhcpServerSweeper(
	client *sdk.APIClient,
	serverID float32,
) {
	s := &dhcpServerSweeper{
		client:   client,
		serverID: serverID,
	}

	resource.AddTestSweepers(
		"hpe_morpheus_network_dhcp_server",
		&resource.Sweeper{
			Name: "hpe_morpheus_network_dhcp_server",
			F:    s.Sweep,
		})
}

// Sweep deletes DHCP servers whose name starts with the test prefix.
// The list endpoint returns an untyped interface{}, so we decode manually.
func (s *dhcpServerSweeper) Sweep(_ string) error {
	ctx := context.Background()

	if s.client == nil {
		log.Printf(
			"[INFO] No client provided, skipping network dhcp server sweep",
		)

		return nil
	}

	resp, hresp, err := s.client.NetworksAPI.
		GetNetworkDhcpServers(ctx, s.serverID).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		return fmt.Errorf(
			"failed to list network dhcp servers: %s",
			errfmt.ErrMsg(err, hresp),
		)
	}

	items, ok := resp.GetNetworkDhcpServers().([]interface{})
	if !ok {
		log.Printf(
			"[INFO] Unexpected dhcp servers list type, skipping sweep",
		)

		return nil
	}

	var sweptCount int
	var sweepErrors []string

	for _, item := range items {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		name, _ := m["name"].(string)
		if !strings.HasPrefix(name, testDhcpServerPrefix) {
			log.Printf(
				"[INFO] Skipping network dhcp server (name): %s", name,
			)

			continue
		}

		idFloat, ok := m["id"].(float64)
		if !ok {
			log.Printf(
				"[INFO] Skipping network dhcp server (id): %s", name,
			)

			continue
		}

		id := int64(idFloat)
		log.Printf(
			"[INFO] Sweeping network dhcp server: %s (id: %d)",
			name, id,
		)

		_, hresp, err := s.client.NetworksAPI.
			DeleteNetworkDhcpServer(ctx, id, s.serverID).Execute()
		if err != nil || hresp.StatusCode != http.StatusOK {
			errMsg := fmt.Sprintf(
				"failed to delete network dhcp server %s (id: %d): %s",
				name, id, errfmt.ErrMsg(err, hresp),
			)
			log.Printf("[ERROR] %s", errMsg)
			sweepErrors = append(sweepErrors, errMsg)

			continue
		}

		sweptCount++
	}

	log.Printf(
		"[INFO] Network DHCP server sweep completed. Swept: %d, errors: %d",
		sweptCount, len(sweepErrors),
	)

	if len(sweepErrors) > 0 {
		return fmt.Errorf("%s", strings.Join(sweepErrors, "\n"))
	}

	return nil
}
