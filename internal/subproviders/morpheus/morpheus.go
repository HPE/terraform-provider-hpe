// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package morpheus

import (
	"context"
	"errors"

	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/clientfactory"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/constants"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/model"
	"github.com/HPE/terraform-provider-hpe/subprovider"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ subprovider.SubProvider = (*MorpheusSubProvider)(nil)

type MorpheusSubProvider struct{}

func New() subprovider.SubProvider {
	return &MorpheusSubProvider{}
}

func (s MorpheusSubProvider) Configure(_ context.Context, f func(any)) (any, error) {
	var m []model.SubModel

	f(&m)

	switch len(m) {
	case 0:
		// no morpheus provider block
		return nil, nil
	case 1:

		return clientfactory.New(m[0]), nil
	default:
		msg := "invalid morpheus provider block length"

		return nil, errors.New(msg)
	}
}

func (MorpheusSubProvider) GetName(_ context.Context) string {
	return constants.SubProviderName
}

func (MorpheusSubProvider) GetSchema(_ context.Context) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"url": schema.StringAttribute{
			Required: true,
		},
	}
}

func (MorpheusSubProvider) GetDataSources(
	_ context.Context,
) []func() datasource.DataSource {
	return nil
}

func (s MorpheusSubProvider) GetResources(
	_ context.Context,
) []func() resource.Resource {
	// Can uncomment this once we have an actual resource
	// f := func(r resource.Resource) func() resource.Resource {
	//   return func() resource.Resource {
	//    return r
	//   }
	// }
	// // s.model contents not  populated yet
	// cf := clientfactory.New(s.model)
	// resources := []func() resource.Resource{
	//   f(xxx.NewResource(cf)),
	// }
	// return resources

	return []func() resource.Resource{}
}
