package resource

import (
	"context"

	"github.com/HewlettPackard/hpe-morpheus-client/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

type DataSource interface {
	datasource.DataSource
	NewClient(ctx context.Context) *client.APIClient
}
