package resource

import (
	"context"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

type Resource interface {
	resource.Resource
	NewClient(ctx context.Context) *sdk.APIClient
}
