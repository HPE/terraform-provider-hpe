package resource

import (
	"context"

	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

type Resource interface {
	resource.Resource
	NewClient(ctx context.Context) client.Client
}
