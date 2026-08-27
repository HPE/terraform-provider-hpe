// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package adapter

import (
	"github.com/hashicorp/terraform-plugin-framework/provider"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/model"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/clientfactory"
)

// Returns a Morpheus provider via the Adapter Layer
// This saves having to manually wrap a Morpheus provider where an
// adapted Morpheus provider is needed, e.g. Acceptance Tests.
func NewMorpheus(opts ...morpheus.Option) provider.Provider {
	f := func(m model.MorpheusProviderModel) *clientfactory.ClientFactory {
		return clientfactory.New(m)
	}

	p := &morpheus.MorpheusProvider{
		NewClientFactory: f,
	}

	// Apply any options
	for _, opt := range opts {
		opt(p)
	}

	return NewAdaptedProvider(p)
}
