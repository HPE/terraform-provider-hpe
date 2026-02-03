// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package datasource

import (
	"context"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

type DataSource interface {
	datasource.DataSource
	NewClient(ctx context.Context) *sdk.APIClient
}
