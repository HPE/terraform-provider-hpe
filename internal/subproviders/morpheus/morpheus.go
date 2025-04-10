// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package morpheus

import (
	"errors"

	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/clientfactory"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/constants"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/model"
	"github.com/HPE/terraform-provider-hpe/subprovider"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ subprovider.SubProvider = (*SubProvider)(nil)

type SubProvider struct{}

func New() subprovider.SubProvider {
	return SubProvider{}
}

func (SubProvider) Configure(f func(any)) error {
	var m []model.SubModel

	f(&m)

	switch len(m) {
	case 0:
		// no morpheus provider block
		return nil
	case 1:
		clientfactory.SetClientFactory(m[0])

		return nil
	default:
		msg := "invalid morpheus provider block length"

		return errors.New(msg)
	}
}

func (SubProvider) GetName() string {
	return constants.SubProviderName
}

func (SubProvider) GetSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{}
}

func (SubProvider) GetDataSources() []func() datasource.DataSource {
	return nil
}

func (SubProvider) GetResources() []func() resource.Resource {
	return nil
}
