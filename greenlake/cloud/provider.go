// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package cloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/greenlake/cloud/datasources/morpheusdetails"
	vmaascmpclient "github.com/HPE/terraform-provider-hpe/greenlake/cloud/sdk/vmaascmp/client"
	"github.com/HPE/terraform-provider-hpe/greenlake/sdk/token/iamversion"
	"github.com/HPE/terraform-provider-hpe/greenlake/sdk/token/retrieve"
	"github.com/HPE/terraform-provider-hpe/greenlake/sdk/token/serviceclient"
)

var _ provider.Provider = &GreenLakeCloudProvider{}

// ProviderName is the HCL block name this child provider exposes inside the
// parent "hpe" provider, i.e.:
//
//	provider "hpe" {
//	  greenlake_cloud {
//	    # ...
//	  }
//	}
const ProviderName = "greenlake_cloud"

const defaultBrokerURL = "https://vmaas-broker.us1.greenlake-hpe.com"

// GreenLakeCloudProviderModel maps the greenlake_cloud provider schema to a Go
// struct for use in Configure.
type GreenLakeCloudProviderModel struct {
	UserID     types.String `tfsdk:"user_id"`
	UserSecret types.String `tfsdk:"user_secret"`
	Location   types.String `tfsdk:"location"`
	Space      types.String `tfsdk:"space"`
	IssuerURL  types.String `tfsdk:"issuer_url"`
	IAMToken   types.String `tfsdk:"iam_token"`
	BrokerURL  types.String `tfsdk:"broker_url"`
}

// GreenLakeCloudProvider is the GreenLake Cloud Services (GLCS) child provider.
type GreenLakeCloudProvider struct{}

// Option configures a GreenLakeCloudProvider.
type Option func(*GreenLakeCloudProvider)

// New constructs a GreenLakeCloudProvider, applying any options.
func New(opts ...Option) *GreenLakeCloudProvider {
	p := &GreenLakeCloudProvider{}

	for _, opt := range opts {
		opt(p)
	}

	return p
}

func (p *GreenLakeCloudProvider) Metadata(
	_ context.Context,
	_ provider.MetadataRequest,
	resp *provider.MetadataResponse,
) {
	resp.TypeName = ProviderName
	resp.Version = "dev"
}

func (p *GreenLakeCloudProvider) Schema(
	_ context.Context,
	_ provider.SchemaRequest,
	resp *provider.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"user_id": schema.StringAttribute{
				Description: "GreenLake API client ID used for authentication.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.AlsoRequires(
						path.MatchRelative().AtParent().AtName("user_secret"),
						path.MatchRelative().AtParent().AtName("issuer_url"),
					),
				},
			},
			"user_secret": schema.StringAttribute{
				Description: "GreenLake API client secret used for authentication.",
				Optional:    true,
				Sensitive:   true,
				Validators: []validator.String{
					stringvalidator.AlsoRequires(
						path.MatchRelative().AtParent().AtName("user_id"),
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
				Description: `GreenLake IAM Issuer URL used to generate access tokens. This should be set to the "Issuer" URL of the API client.`,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.AlsoRequires(
						path.MatchRelative().AtParent().AtName("user_id"),
						path.MatchRelative().AtParent().AtName("user_secret"),
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
						path.MatchRelative().AtParent().AtName("user_id"),
						path.MatchRelative().AtParent().AtName("user_secret"),
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
	}
}

func (p *GreenLakeCloudProvider) Configure(
	ctx context.Context,
	req provider.ConfigureRequest,
	resp *provider.ConfigureResponse,
) {
	var m GreenLakeCloudProviderModel

	diags := req.Config.Get(ctx, &m)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)

		return
	}

	// Obtain a GreenLake IAM token via the serviceclient handler.
	// The handler supports auto-refresh and retry; WithIAMToken short-circuits
	// to a pre-supplied token when set.
	handlerOpts := []serviceclient.CreateOpt{
		serviceclient.WithIAMVersion(iamversion.GLCS),
	}

	if !m.IAMToken.IsNull() && m.IAMToken.ValueString() != "" {
		handlerOpts = append(handlerOpts, serviceclient.WithIAMToken(m.IAMToken.ValueString()))
	} else {
		handlerOpts = append(handlerOpts,
			serviceclient.WithIAMServiceURL(m.IssuerURL.ValueString()),
			serviceclient.WithClientID(m.UserID.ValueString()),
			serviceclient.WithClientSecret(m.UserSecret.ValueString()),
		)
	}

	handler, err := serviceclient.NewHandler(handlerOpts...)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create GreenLake IAM token handler",
			err.Error(),
		)

		return
	}

	getToken := retrieve.NewTokenRetrieveFunc(handler)

	iamToken, err := getToken(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to obtain GreenLake IAM token",
			err.Error(),
		)

		return
	}

	// Build the VMaaS broker client. The morpheus_details data source calls
	// GetCMPDetails on this client in its Read to perform the actual token
	// exchange, so Configure only needs to construct and authenticate it.
	brokerURL := defaultBrokerURL
	if !m.BrokerURL.IsNull() && m.BrokerURL.ValueString() != "" {
		brokerURL = m.BrokerURL.ValueString()
	}

	brokerCfg := vmaascmpclient.NewConfiguration()
	brokerCfg.Host = brokerURL

	if !m.Location.IsNull() {
		brokerCfg.DefaultQueryParams["location"] = m.Location.ValueString()
	}

	if !m.Space.IsNull() {
		brokerCfg.DefaultQueryParams["space"] = m.Space.ValueString()
	}

	brokerClient := vmaascmpclient.NewAPIClient(brokerCfg)

	// Inject the GreenLake IAM token on every request's context, which
	// prepareRequest reads for Bearer auth.
	brokerClient.SetMetaFnAndVersion(nil, 0, func(ctx *context.Context, _ interface{}) {
		*ctx = context.WithValue(*ctx, vmaascmpclient.ContextAccessToken, iamToken)
	})

	resp.DataSourceData = brokerClient
}

func (p *GreenLakeCloudProvider) Resources(
	_ context.Context,
) []func() resource.Resource {
	return nil
}

func (p *GreenLakeCloudProvider) DataSources(
	_ context.Context,
) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		morpheusdetails.NewDataSource,
	}
}
