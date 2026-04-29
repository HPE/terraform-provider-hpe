// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package sweep

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/clientfactory"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
)

const testPrefix = "TestAccMorpheusNetworkFirewallRule"

func init() {
	resource.AddTestSweepers(
		"hpe_morpheus_network_firewall_rule",
		&resource.Sweeper{
			Name: "hpe_morpheus_network_firewall_rule",
			F: func(_ string) error {
				return sweepFirewallRules()
			},
		},
	)
}

func sweepFirewallRules() error {
	serverIDStr := os.Getenv("MORPHEUS_TEST_NETWORK_SERVER_ID")
	if serverIDStr == "" {
		log.Println("[WARN] MORPHEUS_TEST_NETWORK_SERVER_ID not set; skipping network firewall rule sweep")

		return nil
	}

	serverID, err := strconv.ParseInt(serverIDStr, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid MORPHEUS_TEST_NETWORK_SERVER_ID: %w", err)
	}

	ctx := context.Background()

	client, err := newSweepClient(ctx)
	if err != nil {
		log.Printf("[WARN] Cannot create sweep client: %v", err)

		return nil
	}

	resp, httpResp, err := client.NetworksAPI.
		GetNetworkFirewallRules(ctx, serverID).Execute()
	if err != nil || httpResp.StatusCode != http.StatusOK {
		return fmt.Errorf("listing firewall rules failed: %s", errfmt.ErrMsg(err, httpResp))
	}

	rules := resp.GetRules()

	rulesJSON, err := json.Marshal(rules)
	if err != nil {
		return fmt.Errorf("marshaling rules response: %w", err)
	}

	var rulesList []map[string]interface{}
	if err := json.Unmarshal(rulesJSON, &rulesList); err != nil {
		return fmt.Errorf("decoding rules list: %w", err)
	}

	for _, rule := range rulesList {
		name, _ := rule["name"].(string)
		if !strings.HasPrefix(name, testPrefix) {
			continue
		}

		idRaw, ok := rule["id"]
		if !ok {
			continue
		}

		var id int64

		switch v := idRaw.(type) {
		case float64:
			id = int64(v)
		case json.Number:
			id, _ = v.Int64()
		default:
			continue
		}

		log.Printf("[INFO] Sweeping network firewall rule %q (id=%d)", name, id)

		_, httpResp, err := client.NetworksAPI.
			DeleteNetworkFirewallRule(ctx, id, serverID).Execute()
		if err != nil {
			log.Printf("[WARN] Failed to delete firewall rule %d: %s", id, errfmt.ErrMsg(err, httpResp))
		}
	}

	return nil
}

func newSweepClient(ctx context.Context) (*sdk.APIClient, error) {
	var username, password string

	url, ok := os.LookupEnv("TF_VAR_testacc_morpheus_url")
	if !ok {
		return nil, errors.New("TF_VAR_testacc_morpheus_url not set")
	}

	token, ok := os.LookupEnv("TF_VAR_testacc_morpheus_access_token")
	if !ok {
		username, ok = os.LookupEnv("TF_VAR_testacc_morpheus_username")
		if !ok {
			return nil, errors.New(
				"one of TF_VAR_testacc_morpheus_access_token or " +
					"TF_VAR_testacc_morpheus_username must be set",
			)
		}

		password, ok = os.LookupEnv("TF_VAR_testacc_morpheus_password")
		if !ok {
			return nil, errors.New(
				"one of TF_VAR_testacc_morpheus_access_token or " +
					"TF_VAR_testacc_morpheus_password must be set",
			)
		}
	}

	_, insecure := os.LookupEnv("TF_VAR_testacc_morpheus_insecure")

	var opts []clientfactory.ClientOption
	if insecure {
		opts = append(opts, clientfactory.WithInsecureTLS())
	}

	client := clientfactory.NewAPIClient(
		ctx,
		url,
		username,
		password,
		"",
		token,
		opts...,
	)

	return client, nil
}
