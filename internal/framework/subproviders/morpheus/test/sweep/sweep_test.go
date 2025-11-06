// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package sweep

import (
	"context"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/resources/policy"
)

func init() {
	// Skip sweep client creation and sweeper registration in short mode
	if testing.Short() {
		log.Printf("[INFO] Skipping sweeper registration in short mode")
		return
	}

	ctx := context.Background()
	client, err := NewSweepClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create sweep client: %v", err)
	}

	// Register sweepers that use the centralized client
	policy.NewPolicySweeper(client)

	// Register sweepers that still create their own clients (to be converted later)
	Datastores()
	Instances()
	Networks()
	Users()
}

func TestMain(m *testing.M) {
	resource.TestMain(m)
}
