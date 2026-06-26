// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package datasource

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

type DataSource interface {
	datasource.DataSource
	NewClient(ctx context.Context) *sdk.APIClient
}
