// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package morpheus

import (
	"context"
	"errors"

	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/constants"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/model"
	"github.com/HPE/terraform-provider-hpe/subprovider"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ subprovider.SubProvider = (*SubProvider)(nil)

type SubProvider struct {
	model *model.SubModel
}

func New() subprovider.SubProvider {
	return &SubProvider{}
}

func (s *SubProvider) Configure(_ context.Context, f func(any)) error {
	var m []model.SubModel

	f(&m)

	switch len(m) {
	case 0:
		// no morpheus provider block
		return nil
	case 1:
		s.model = &m[0]

		return nil
	default:
		msg := "invalid morpheus provider block length"

		return errors.New(msg)
	}
}

func (SubProvider) GetName(_ context.Context) string {
	return constants.SubProviderName
}

func (SubProvider) GetSchema(_ context.Context) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"url": schema.StringAttribute{
			Required: true,
		},
	}
}

func (SubProvider) GetDataSources(
	_ context.Context,
) []func() datasource.DataSource {
	return nil
}

func (s SubProvider) GetResources(
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
