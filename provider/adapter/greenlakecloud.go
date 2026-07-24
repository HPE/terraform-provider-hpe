// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package adapter

import (
	"github.com/hashicorp/terraform-plugin-framework/provider"

	glcs "github.com/HPE/terraform-provider-hpe/greenlake/cloud"
)

// NewGreenLakeCloud returns a GreenLake Cloud Services (GLCS) provider via the
// Adapter Layer. This saves having to manually wrap a GreenLakeCloud provider
// where an adapted provider is needed, e.g. Acceptance Tests.
func NewGreenLakeCloud(opts ...glcs.Option) provider.Provider {
	p := glcs.New(opts...)

	return NewAdaptedProvider(p)
}
