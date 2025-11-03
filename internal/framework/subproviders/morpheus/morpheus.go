// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package morpheus

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/clientfactory"
	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/constants"
	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/model"
	"github.com/HPE/terraform-provider-hpe/subprovider"
)

var _ subprovider.SubProvider = (*SubProvider)(nil)

type Option func(*SubProvider)

// Option to override newClientFactory
func WithClientFactory(f func(model.SubModel) *clientfactory.ClientFactory) Option {
	return func(sp *SubProvider) {
		sp.newClientFactory = f
	}
}

type SubProvider struct {
	newClientFactory func(model.SubModel) *clientfactory.ClientFactory
}

func New(opts ...Option) subprovider.SubProvider {
	f := func(m model.SubModel) *clientfactory.ClientFactory {
		return clientfactory.New(m)
	}

	sp := &SubProvider{
		newClientFactory: f,
	}

	// Apply any options
	for _, opt := range opts {
		opt(sp)
	}

	return sp
}

func (s SubProvider) Configure(_ context.Context, f func(any)) (any, error) {
	var m []model.SubModel

	f(&m)

	switch len(m) {
	case 0:
		// no morpheus provider block
		return nil, nil
	case 1:

		return s.newClientFactory(m[0]), nil
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
			Description: "Morpheus instance URL",
			Required:    true,
			Validators: []validator.String{
				stringvalidator.Any(
					stringvalidator.AlsoRequires(parentBlock.AtName("username")),
					stringvalidator.AlsoRequires(parentBlock.AtName("access_token")),
				),
			},
		},
		"username": schema.StringAttribute{
			Description: "Morpheus username for authentication, required if password is set",
			Optional:    true,
			Validators: []validator.String{
				stringvalidator.AlsoRequires(parentBlock.AtName("password")),
			},
		},
		"password": schema.StringAttribute{
			Description: "Morpheus password for authentication, required if username is set",
			Optional:    true,
			Sensitive:   true,
			Validators: []validator.String{
				stringvalidator.AlsoRequires(parentBlock.AtName("username")),
			},
		},
		"access_token": schema.StringAttribute{
			Description: "Morpheus access token for authentication",
			Optional:    true,
			Sensitive:   true,
			Validators: []validator.String{
				stringvalidator.ConflictsWith(parentBlock.AtName("username")),
				stringvalidator.ConflictsWith(parentBlock.AtName("password")),
			},
		},
		"insecure": schema.BoolAttribute{
			Description: "Explicitly allow the provider to perform " +
				"\"insecure\" SSL requests. If omitted, " +
				"default value is `false`",
			Optional: true,
		},
	}
}
