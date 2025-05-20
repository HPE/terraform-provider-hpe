package resource

import (
	"context"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

type DataSource interface {
	datasource.DataSource
	NewClient(ctx context.Context) *sdk.APIClient
}
