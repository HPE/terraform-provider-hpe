package resource

import (
	"context"

	"github.com/HewlettPackard/hpe-morpheus-client/client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

type Resource interface {
	resource.Resource
	NewClient(ctx context.Context) *client.APIClient
}
