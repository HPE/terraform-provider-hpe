// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package morpheus

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/morpheus/model"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/clientfactory"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/constants"
)

var _ provider.Provider = &MorpheusProvider{}

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

	// Apply any options
	for _, opt := range opts {
		opt(p)
	}

	return p
}

func (p *MorpheusProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = constants.ProviderName
	resp.Version = "dev"
}

func (p *MorpheusProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var m model.MorpheusProviderModel

	diags := req.Config.Get(ctx, &m)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)

		return
	}

	// A greenlake_connected block supplies the Morpheus url and access token
	// by exchanging GreenLake credentials, rather than having them configured
	// directly. The schema limits the block to at most one element.
	if len(m.GreenLakeConnected) > 0 {
		url, token, err := greenlakeConnectedTokenExchange(ctx, &m.GreenLakeConnected[0])
		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to obtain Morpheus connection details from GreenLake",
				err.Error(),
			)

			return
		}

		m.URL = types.StringValue(url)
		m.AccessToken = types.StringValue(token)
	}

	cf := p.NewClientFactory(m)
	resp.ResourceData = cf
	resp.DataSourceData = cf
}

func (p *MorpheusProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"url": schema.StringAttribute{
				Description: "Morpheus instance URL. May be omitted when it is " +
					"supplied by a greenlake_connected block.",
				Optional: true,
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
		Blocks: map[string]schema.Block{
			"greenlake_connected": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"client_id": schema.StringAttribute{
							Description: "GreenLake API client ID used for authentication.",
							Optional:    true,
							Validators: []validator.String{
								stringvalidator.AlsoRequires(
									path.MatchRelative().AtParent().AtName("client_secret"),
									path.MatchRelative().AtParent().AtName("issuer_url"),
								),
							},
						},
						"client_secret": schema.StringAttribute{
							Description: "GreenLake API client secret used for authentication.",
							Optional:    true,
							Sensitive:   true,
							Validators: []validator.String{
								stringvalidator.AlsoRequires(
									path.MatchRelative().AtParent().AtName("client_id"),
									path.MatchRelative().AtParent().AtName("issuer_url"),
								),
							},
						},
						"location": schema.StringAttribute{
							Description: "Location of the GreenLake VMaaS service.",
							Optional:    true,
						},
						"space": schema.StringAttribute{
							Description: "GreenLake VMaaS space name (IAM Space) used for the broker exchange.",
							Optional:    true,
						},
						"issuer_url": schema.StringAttribute{
							Description: `GreenLake IAM Issuer URL used to generate access tokens. ` +
								`This should be set to the "Issuer" URL of the API client.`,
							Optional: true,
							Validators: []validator.String{
								stringvalidator.AlsoRequires(
									path.MatchRelative().AtParent().AtName("client_id"),
									path.MatchRelative().AtParent().AtName("client_secret"),
								),
							},
						},
						"iam_token": schema.StringAttribute{
							Description: "Pre-generated GreenLake IAM token. If set, token " +
								"generation from credentials is skipped.",
							Optional:  true,
							Sensitive: true,
							Validators: []validator.String{
								stringvalidator.ConflictsWith(
									path.MatchRelative().AtParent().AtName("client_id"),
									path.MatchRelative().AtParent().AtName("client_secret"),
									path.MatchRelative().AtParent().AtName("issuer_url"),
								),
							},
						},
						"broker_url": schema.StringAttribute{
							Description: "URL of the VMaaS broker used for the CMP details exchange. " +
								"Defaults to the US1 production broker if not set.",
							Optional: true,
						},
					},
				},
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
				// Only Description is set, never MarkdownDescription: the
				// SDKv2 provider this is muxed with can only report plain text
				// descriptions, and the schemas must match exactly.
				Description: "Configuration block for using Morpheus with GreenLake Connected",
			},
		},
	}
}
