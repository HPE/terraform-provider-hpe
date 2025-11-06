// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

// Package sweep allows deletion of dangling test resources
package sweep

import (
	"context"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/resources/policy"
)

func Policies() {
	resource.AddTestSweepers(
		"hpe_morpheus_policy",
		&resource.Sweeper{
			Name: "hpe_morpheus_policy",
			F: func(region string) error {
				ctx := context.Background()
				client, err := NewSweepClient(ctx)
				if err != nil {
					return err
				}
				return policy.SweepPolicies(client)
			},
		})
}
