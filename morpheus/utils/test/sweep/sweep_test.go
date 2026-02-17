// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package sweep

import (
	"context"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/datastore"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/instance"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/network"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/policy"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/user"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/systemoverride"
)

func init() {
	ctx := context.Background()
	client, err := NewSweepClient(ctx)
	if err != nil {
		// Log the error but don't fail - this is expected in short mode or when env vars aren't set
		log.Printf("[WARN] Cannot create sweep client (likely running in short mode or env vars not set): %v", err)

		return
	}

	// Register sweepers that use the centralized client
	policy.NewPolicySweeper(client)
	datastore.NewDatastoreSweeper(client)
	instance.NewInstanceSweeper(client)
	network.NewNetworkSweeper(client)
	user.NewUserSweeper(client)
}

func TestMain(m *testing.M) {
	systemoverride.ParseFlags()
	resource.TestMain(m)
}
