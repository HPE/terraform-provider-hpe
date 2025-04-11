// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package clientfactory

import (
	"context"

	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/client"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/model"
)

func New(m *model.SubModel) *ClientFactory {
	return &ClientFactory{model: m}
}

type ClientFactory struct {
	model *model.SubModel
}

func (c ClientFactory) NewClient(_ context.Context) client.Client {
	return client.Client{URL: c.model.URL}
}
