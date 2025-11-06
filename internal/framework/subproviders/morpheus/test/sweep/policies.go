// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

// Package sweep allows deletion of dangling test resources
package sweep

import (
	"context"
	"log"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/resources/policy"
)

func Policies() {
	ctx := context.Background()
	client, err := NewSweepClient(ctx)
	if err != nil {
		log.Printf("[INFO] Cannot create policy sweep client, skipping policy sweeper registration: %v", err)

		return
	}

	policy.NewPolicySweeper(client)
}
