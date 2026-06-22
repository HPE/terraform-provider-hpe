package morpheus

import (
	"context"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/clientfactory"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/model"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ provider.Provider = &MorpheusProvider{}

const providerName = "morpheus"

type MorpheusProvider struct {
	// model.SubModel = morpheus SubModel
	// we can tidy this up later.
	newClientFactory func(model.SubModel) *clientfactory.ClientFactory
}

type MorpheusProviderModel struct {
	URL             types.String `tfsdk:"url"`
	Username        types.String `tfsdk:"username"`
	Password        types.String `tfsdk:"password"`
	AccessToken     types.String `tfsdk:"access_token"`
	TenantSubdomain types.String `tfsdk:"tenant_subdomain"`
	Insecure        types.Bool   `tfsdk:"insecure"`
}

func NewMorpheusProvider() *MorpheusProvider {
	f := func(m model.SubModel) *clientfactory.ClientFactory {
		return clientfactory.New(m)
	}

	p := &MorpheusProvider{
		newClientFactory: f,
	}

	return p
}

func (p *MorpheusProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = providerName
	resp.Version = "dev"
}

func (p *MorpheusProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var m model.SubModel // Morpheus Provider Model

	diags := req.Config.Get(ctx, &m)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)

		return
	}

	cf := p.newClientFactory(m)
	resp.ResourceData = cf
	resp.DataSourceData = cf
	// // diags := req.Config.GetAttribute(ctx, path.Root(providerName), &m)
	// resp.Diagnostics.Append(diags...)
	// if resp.Diagnostics.HasError() {
	// 	return
	// }
	//
	// switch len(m) {
	// case 0:
	// 	return
	// case 1:
	// 	cf := p.newClientFactory(m[0])
	// 	resp.ResourceData = cf
	// 	resp.DataSourceData = cf
	// default:
	// 	resp.Diagnostics.AddError("Failed to configure morpheus", "invalid morpheus provider block length")
	// }
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
