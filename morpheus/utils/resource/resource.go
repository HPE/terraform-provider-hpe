// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package resource

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

type Resource interface {
	resource.Resource
	NewClient(ctx context.Context) *sdk.APIClient
}
