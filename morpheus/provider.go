// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package morpheus

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/HPE/terraform-provider-hpe/morpheus/model"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/clientfactory"
)

var _ provider.Provider = &MorpheusProvider{}

const ProviderName = "morpheus"

type MorpheusProvider struct {
	NewClientFactory func(model.MorpheusProviderModel) *clientfactory.ClientFactory
}

type Option func(*MorpheusProvider)

func MorpheusWithClientFactory(f func(model.MorpheusProviderModel) *clientfactory.ClientFactory) Option {
	return func(p *MorpheusProvider) {
		p.NewClientFactory = f
	}
}

func New(opts ...Option) *MorpheusProvider {
	f := func(m model.MorpheusProviderModel) *clientfactory.ClientFactory {
		return clientfactory.New(m)
	}

	p := &MorpheusProvider{
		NewClientFactory: f,
	}

	return p
}

func (p *MorpheusProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = ProviderName
	resp.Version = "dev"
}

func (p *MorpheusProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var m model.MorpheusProviderModel

	diags := req.Config.Get(ctx, &m)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)

		return
	}

	cf := p.NewClientFactory(m)
	resp.ResourceData = cf
	resp.DataSourceData = cf
}

func (p *MorpheusProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"url": schema.StringAttribute{
				Description: "Morpheus instance URL",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.Any(
						stringvalidator.AlsoRequires(path.MatchRelative().AtParent().AtName("username")),
						stringvalidator.AlsoRequires(path.MatchRelative().AtParent().AtName("access_token")),
					),
				},
			},
			"username": schema.StringAttribute{
				Description: "Morpheus username for authentication, required if password is set",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.AlsoRequires(path.MatchRelative().AtParent().AtName("password")),
				},
			},
			"password": schema.StringAttribute{
				Description: "Morpheus password for authentication, required if username is set",
				Optional:    true,
				Sensitive:   true,
				Validators: []validator.String{
					stringvalidator.AlsoRequires(path.MatchRelative().AtParent().AtName("username")),
				},
			},
			"access_token": schema.StringAttribute{
				Description: "Morpheus access token for authentication",
				Optional:    true,
				Sensitive:   true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("username")),
					stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("password")),
					stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("tenant_subdomain")),
				},
			},
			"tenant_subdomain": schema.StringAttribute{
				Description: "Morpheus tenant subdomain used for authentication",
				Optional:    true,
			},
			"insecure": schema.BoolAttribute{
				Description: "Explicitly allow the provider to perform " +
					"\"insecure\" SSL requests. If omitted, " +
					"default value is `false`",
				Optional: true,
			},
		},
	}
}
