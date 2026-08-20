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
	"github.com/HPE/terraform-provider-hpe/morpheus/pce"
	"github.com/HPE/terraform-provider-hpe/morpheus/pce/sdk/token/iamversion"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/clientfactory"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/constants"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/validators"
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

	// An identity block supplies the Morpheus url and access token by
	// exchanging GreenLake credentials, rather than having them configured
	// directly. The schema limits each block to at most one element and
	// rejects combining them with each other or with direct connection
	// details, so at most one of these cases can apply.
	var (
		url, token string
		err        error
	)

	switch {
	case len(m.PCEIdentity) > 0:
		url, token, err = pceIdentityTokenExchange(ctx, &m.PCEIdentity[0])
	case len(m.PCEDisconnectedIdentity) > 0:
		url, token, err = pceDisconnectedIdentityTokenExchange(ctx, &m.PCEDisconnectedIdentity[0])
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to obtain Morpheus connection details from GreenLake",
			err.Error(),
		)

		return
	}

	if url != "" {
		m.URL = types.StringValue(url)
		m.AccessToken = types.StringValue(token)
	}

	// "url" is optional in the schema because an identity block can supply it,
	// so nothing has rejected a configuration that sets neither. Report it here
	// rather than building a client that cannot reach anything and failing
	// later with a less obvious error.
	if m.URL.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Missing Morpheus connection details",
			missingConnectionDetails,
		)

		return
	}

	cf := p.NewClientFactory(m)
	resp.ResourceData = cf
	resp.DataSourceData = cf
}

// missingConnectionDetails is reported when neither the connection details nor
// a usable identity block were configured. It mirrors the equivalent message in
// the SDKv2 provider so that the two do not tell users different things.
const missingConnectionDetails = `The morpheus provider block does not set "url", and no usable identity block
was found.

Set the connection details explicitly, or configure a pce_identity or
pce_disconnected_identity block so that they can be obtained from GreenLake:

 provider "hpe" {
   morpheus {
     url          = "https://example.com"
     access_token = "..."
   }
 }`

// identityBlockValidators returns the validators for an identity block, which
// may not be combined with the named sibling paths. Both identity blocks use
// it so that they cannot drift apart.
//
// The raw listvalidator.SizeAtMost value has to appear in the returned slice:
// utils/convert derives the SDKv2 MaxItems by reflecting on the validator's
// concrete type name, so wrapping or replacing it would silently drop the
// constraint and leave the muxed providers disagreeing.
//
// "insecure" is deliberately not conflicted with: it is a transport setting
// that applies however the Morpheus URL was obtained.
func identityBlockValidators(conflictsWith ...string) []validator.List {
	expressions := make([]path.Expression, 0, len(conflictsWith))
	for _, name := range conflictsWith {
		expressions = append(expressions, path.MatchRelative().AtParent().AtName(name))
	}

	return []validator.List{
		listvalidator.SizeAtMost(1),
		listvalidator.ConflictsWith(expressions...),
	}
}

func (p *MorpheusProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"url": schema.StringAttribute{
				Description: "Morpheus instance URL. May be omitted when it is " +
					"supplied by a pce_identity or pce_disconnected_identity block.",
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
			"pce_identity": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Validators: []validator.Object{
						validators.IdentityCredentialsOrTokenValidator(),
					},
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
							Description: "The PCE instance's Location.",
							Required:    true,
						},
						"space": schema.StringAttribute{
							Description: "The name of the GreenLake Space that the PCE instance is in.",
							Required:    true,
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
							Description: "GreenLake IAM access token. If set, token " +
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
							Description: "URL of the PCE broker. Defaults to the " +
								"HPE-hosted broker if not set.",
							Optional: true,
						},
					},
				},
				Validators: identityBlockValidators(
					// The mutual conflict is declared here only, and not on
					// the other block as well, so that configuring both blocks
					// is reported once rather than once per block.
					"pce_disconnected_identity",
					"url", "username", "password", "access_token", "tenant_subdomain",
				),
				// Only Description is set, never MarkdownDescription: the
				// SDKv2 provider this is muxed with can only report plain text
				// descriptions, and the schemas must match exactly.
				Description: "Configuration block for using Morpheus with PCE (Private Cloud Enterprise) Identity",
			},
			"pce_disconnected_identity": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Validators: []validator.Object{
						validators.IdentityCredentialsOrTokenValidator(),
					},
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
							Description: "GreenLake IAM access token. If set, token " +
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
						"location": schema.StringAttribute{
							Description: "The PCE instance's Location.",
							Required:    true,
						},
						"workspace_id": schema.StringAttribute{
							Description: "The GreenLake Workspace ID that the PCE instance is in.",
							Required:    true,
						},
						"broker_url": schema.StringAttribute{
							Description: "URL of the PCE broker for this deployment. There is " +
								"no default: a Disconnected deployment has no hosted " +
								"broker to fall back to.",
							Required: true,
						},
					},
				},
				Validators: identityBlockValidators(
					"url", "username", "password", "access_token", "tenant_subdomain",
				),
				// Only Description is set, never MarkdownDescription: the
				// SDKv2 provider this is muxed with can only report plain text
				// descriptions, and the schemas must match exactly.
				Description: "Configuration block for using Morpheus with Disconnected PCE " +
					"(Private Cloud Enterprise) Identity",
			},
		},
	}
}

// pceIdentityTokenExchange resolves the Morpheus url and access token
// for a pce_identity block.
//
// The SDKv2 Morpheus provider performs the same exchange from the same
// configuration, so a configuration using this block exchanges once per muxed
// provider.
func pceIdentityTokenExchange(
	ctx context.Context,
	m *model.PCEIdentityModel,
) (string, string, error) {
	// ValueString returns "" for null values, which matches how the SDKv2
	// provider reads the same block. The two must agree so that both resolve
	// the same Morpheus instance.
	return pce.TokenExchange(ctx, pce.Config{
		ClientID:     m.ClientID.ValueString(),
		ClientSecret: m.ClientSecret.ValueString(),
		Location:     m.Location.ValueString(),
		Space:        m.Space.ValueString(),
		IssuerURL:    m.IssuerURL.ValueString(),
		IAMToken:     m.IAMToken.ValueString(),
		BrokerURL:    m.BrokerURL.ValueString(),
		Version:      iamversion.GLCS,
	})
}

// pceDisconnectedIdentityTokenExchange resolves the Morpheus url and access
// token for a pce_disconnected_identity block.
//
// It differs from the Connected exchange only in the IAM version and in how the
// broker request is scoped; see pce.Config.
func pceDisconnectedIdentityTokenExchange(
	ctx context.Context,
	m *model.PCEDisconnectedIdentityModel,
) (string, string, error) {
	return pce.TokenExchange(ctx, pce.Config{
		ClientID:     m.ClientID.ValueString(),
		ClientSecret: m.ClientSecret.ValueString(),
		IssuerURL:    m.IssuerURL.ValueString(),
		IAMToken:     m.IAMToken.ValueString(),
		Location:     m.Location.ValueString(),
		WorkspaceID:  m.WorkspaceID.ValueString(),
		BrokerURL:    m.BrokerURL.ValueString(),
		Version:      iamversion.GLP,
	})
}
