// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package morpheus

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"

	"github.com/HPE/terraform-provider-hpe/morpheus/functions"
)

var _ provider.ProviderWithFunctions = &MorpheusProvider{}

func (p *MorpheusProvider) Functions(_ context.Context) []func() function.Function {
	return []func() function.Function{functions.NewSHA512Crypt}
}
