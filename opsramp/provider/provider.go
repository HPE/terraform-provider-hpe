// SPDX-FileCopyrightText: Copyright Hewlett Packard Enterprise Development LP
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"context"
	"os"

	"github.com/HPE/terraform-provider-opsramp/internal/client"
	"github.com/HPE/terraform-provider-opsramp/internal/data"
	"github.com/HPE/terraform-provider-opsramp/internal/resources"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure OpsRampProvider satisfies various provider interfaces.
var _ provider.Provider = &OpsRampProvider{}
var _ provider.ProviderWithFunctions = &OpsRampProvider{}
var _ provider.ProviderWithEphemeralResources = &OpsRampProvider{}

// OpsRampProvider defines the provider implementation.
type OpsRampProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// OpsRampProviderModel describes the provider data model.
type OpsRampProviderModel struct {
	ClientId     types.String `tfsdk:"client_id"`
	ClientSecret types.String `tfsdk:"client_secret"`
	Endpoint     types.String `tfsdk:"endpoint"`
	Tenant       types.String `tfsdk:"tenant"`
}

func (p *OpsRampProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "opsramp"
	resp.Version = p.version
}

func (p *OpsRampProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"client_id": schema.StringAttribute{
				Optional:    true,
				Description: "OAuth client ID for OpsRamp API",
			},
			"client_secret": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "OAuth client secret for OpsRamp API",
			},
			"endpoint": schema.StringAttribute{
				Optional:    true,
				Description: "OpsRamp Endpoint",
			},
			"tenant": schema.StringAttribute{
				Optional:    true,
				Description: "OpsRamp Tenant ID",
			},
		},
	}
}

func (p *OpsRampProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Retrieve provider data from configuration
	var config OpsRampProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	endpoint := os.Getenv("OPSRAMP_ENDPOINT")
	tenant := os.Getenv("OPSRAMP_TENANT")
	client_id := os.Getenv("OPSRAMP_CLIENT_ID")
	client_secret := os.Getenv("OPSRAMP_CLIENT_SECRET")

	if config.ClientId.ValueString() != "" && config.ClientId.ValueString() != "*****" {
		client_id = config.ClientId.ValueString()
	}

	if config.ClientSecret.ValueString() != "" && config.ClientSecret.ValueString() != "*****" {
		client_secret = config.ClientSecret.ValueString()
	}

	if config.Tenant.ValueString() != "" && config.Tenant.ValueString() != "*****" {
		tenant = config.Tenant.ValueString()
	}

	if config.Endpoint.ValueString() != "" && config.Endpoint.ValueString() != "*****" {
		endpoint = config.Endpoint.ValueString()
	}

	// Create a new OpsRamp client using the configuration values
	client, err := client.NewOpsRampClient(client_id, client_secret, endpoint, tenant)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create OpsRamp API Client",
			"An unexpected error occurred when creating the OpsRamp API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"OpsRamp Client Error: "+err.Error(),
		)
		return
	}

	// Make the HashiCups client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *OpsRampProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		resources.NewResource,
		resources.NewServicemap,
		resources.NewServicemapLink,
		resources.NewDeviceGroup,
		resources.NewSite,
		resources.NewServiceDeskCategory,
		resources.NewServiceDeskUrgency,
		resources.NewServiceDeskBusinessImpact,
		resources.NewClient,
		resources.NewUser,
		resources.NewRole,
		resources.NewPermissionSet,
		resources.NewUserGroup,
		resources.NewCustomIntegration,
		resources.NewIntegration,
		resources.NewIntegrationApp,
		resources.NewIntegrationConfig,
		resources.NewIntegrationEvent,
		resources.NewScript,
		resources.NewScriptCategory,
		resources.NewKBCategory,
		resources.NewKBArticle,
		resources.NewAlertCorrelationPolicy,
		resources.NewAlertEscalationPolicy,
		resources.NewAlertPredictionPolicy,
		resources.NewFirstResponsePolicy,
		resources.NewCredentialSet,
		resources.NewManagementProfile,
		resources.NewScheduleMaintenance,
		resources.NewMetricAlertDefinition,
		resources.NewLogAlertDefinition,
	}
}

func (p *OpsRampProvider) EphemeralResources(ctx context.Context) []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{}
}

func (p *OpsRampProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		data.NewResourceLookupDataSource,
		data.NewResourceTenantDataSource,
		data.NewDataRoleSource,
		data.NewCustomEventAlertSourceDataSource,
		data.NewServiceDeskUrgencyDataSource,
		data.NewServiceDeskBusinessImpactDataSource,
		data.NewServiceDeskCategoryDataSource,
	}
}

func (p *OpsRampProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &OpsRampProvider{
			version: version,
		}
	}
}
