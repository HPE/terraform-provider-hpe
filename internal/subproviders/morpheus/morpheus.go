// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package morpheus

import (
	"context"
	"errors"

	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/clientfactory"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/constants"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/model"
	"github.com/HPE/terraform-provider-hpe/subprovider"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ subprovider.SubProvider = (*SubProvider)(nil)

type SubProvider struct{}

func New() subprovider.SubProvider {
	return &SubProvider{}
}

func (s SubProvider) Configure(_ context.Context, f func(any)) (any, error) {
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

func (SubProvider) GetName(_ context.Context) string {
	return constants.SubProviderName
}

func (SubProvider) GetSchema(_ context.Context) map[string]schema.Attribute {
	parentBlock := path.MatchRelative().AtParent()
	return map[string]schema.Attribute{
		"url": schema.StringAttribute{
			MarkdownDescription: "Morpheus instance URL",
			Required:            true,
		},
		"username": schema.StringAttribute{
			MarkdownDescription: "Morpheus username for authentication, required if password is set",
			Optional:            true,
			Validators: []validator.String{
				stringvalidator.AlsoRequires(parentBlock.AtName("password")),
			},
		},
		"password": schema.StringAttribute{
			MarkdownDescription: "Morpheus password for authentication, required if username is set",
			Optional:            true,
			Sensitive:           true,
			Validators: []validator.String{
				stringvalidator.AlsoRequires(parentBlock.AtName("username")),
			},
		},
		"access_token": schema.StringAttribute{
			MarkdownDescription: "Morpheus access token for authentication",
			Optional:            true,
			Sensitive:           true,
			Validators: []validator.String{
				stringvalidator.ConflictsWith(parentBlock.AtName("username")),
				stringvalidator.ConflictsWith(parentBlock.AtName("password")),
			},
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
