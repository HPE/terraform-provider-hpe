// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package adapter

import (
	"github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/model"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/clientfactory"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
	"github.com/hashicorp/terraform-plugin-framework/provider"
)

// Returns a Morpheus provider via the Adapter Layer
// This saves having to directly import the "adapter" package where an
// adapted Morpheus provider is needed, e.g. Acceptance Tests.
func NewAdaptedMorpheus(opts ...morpheus.Option) provider.Provider {

	f := func(m model.MorpheusProviderModel) *clientfactory.ClientFactory {
		return clientfactory.New(m)
	}

	p := &morpheus.MorpheusProvider{
		NewClientFactory: f,
	}

	return adapter.NewAdaptedProvider(p)
}
