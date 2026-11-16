// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

// Package appliancesetup implements the hpe_morpheus_appliance_setup
// resource: the first-time setup of a freshly installed Morpheus appliance,
// through the unauthenticated /api/setup endpoints behind the setup wizard.
package appliancesetup

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
)

// Interface compliance assertions.
var (
	_ resource.Resource                = &Resource{}
	_ resource.ResourceWithImportState = &Resource{}
)

// No sweeper is needed for this resource: it creates nothing that can be
// removed, and an appliance can only be set up once.

const (
	resourceName = "appliance_setup"

	// defaultCreateTimeout bounds the wait for a starting appliance. A first
	// boot seeds the database before the API answers, which takes a good
	// twenty minutes on a small VM.
	defaultCreateTimeout = 40 * time.Minute
)

// ApplianceSetupModel is the Terraform state model.
type ApplianceSetupModel struct {
	Id                 types.String   `tfsdk:"id"`
	ApplianceURL       types.String   `tfsdk:"appliance_url"`
	ApplianceName      types.String   `tfsdk:"appliance_name"`
	AccountName        types.String   `tfsdk:"account_name"`
	Username           types.String   `tfsdk:"username"`
	PasswordWo         types.String   `tfsdk:"password_wo"`
	Email              types.String   `tfsdk:"email"`
	FirstName          types.String   `tfsdk:"first_name"`
	LastName           types.String   `tfsdk:"last_name"`
	HubMode            types.String   `tfsdk:"hub_mode"`
	LicenseKeyWo       types.String   `tfsdk:"license_key_wo"`
	AlreadyInitialised types.Bool     `tfsdk:"already_initialised"`
	BuildVersion       types.String   `tfsdk:"build_version"`
	Timeouts           timeouts.Value `tfsdk:"timeouts"`
}

// Resource implements the hpe_morpheus_appliance_setup resource.
type Resource struct {
	configure.ResourceWithMorpheusConfigure
}

// NewResource returns a new appliance setup resource.
func NewResource() resource.Resource {
	return &Resource{}
}

// Metadata implements resource.Resource.
func (r *Resource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_" + resourceName
}

// Schema implements resource.Resource.
func (r *Resource) Schema(
	ctx context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Performs the first-time setup of a freshly installed Morpheus appliance: waits for its " +
			"API to come up, then creates the master tenant and its first administrator through the " +
			"unauthenticated setup endpoint the setup wizard uses, optionally installing a licence key. " +
			"An appliance already set up is adopted as is (already_initialised = true). Setup happens " +
			"once: later changes to identity inputs are refused, and destroying the resource " +
			"changes nothing on the appliance.",
		MarkdownDescription: "Performs the first-time setup of a freshly installed Morpheus appliance: waits for its " +
			"API to come up, then creates the master tenant and its first administrator through the " +
			"unauthenticated setup endpoint the setup wizard uses, optionally installing a licence key. " +
			"An appliance already set up is adopted as is (`already_initialised = true`). Setup happens " +
			"once: later changes to identity inputs are refused, and destroying the resource " +
			"changes nothing on the appliance.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "Resource identifier: the appliance URL the provider is configured with.",
				MarkdownDescription: "Resource identifier: the appliance URL the provider is configured with.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"appliance_url": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Description: "The URL the appliance records as its own (Administration > Settings > Appliance URL); " +
					"hosts and agents call back to it. Defaults to the provider's url.",
				MarkdownDescription: "The URL the appliance records as its own (*Administration > Settings > Appliance URL*); " +
					"hosts and agents call back to it. Defaults to the provider's `url`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"appliance_name": schema.StringAttribute{
				Optional: true,
				Description: "A name for the appliance instance. When omitted the appliance generates a " +
					"random identifier.",
				MarkdownDescription: "A name for the appliance instance. When omitted the appliance generates a " +
					"random identifier.",
				Validators: []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"account_name": schema.StringAttribute{
				Required:            true,
				Description:         "Name of the master tenant to create.",
				MarkdownDescription: "Name of the master tenant to create.",
				Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"username": schema.StringAttribute{
				Required:            true,
				Description:         "Login of the master tenant's first administrator.",
				MarkdownDescription: "Login of the master tenant's first administrator.",
				Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"password_wo": schema.StringAttribute{
				Required:  true,
				Sensitive: true,
				WriteOnly: true,
				Description: "Password of the first administrator; it must satisfy the appliance's password policy. " +
					"Write-only: sent once at setup and never stored in state or plan. Changing it later has no " +
					"effect on the appliance.",
				MarkdownDescription: "Password of the first administrator; it must satisfy the appliance's password policy. " +
					"Write-only: sent once at setup and never stored in state or plan. Changing it later has no " +
					"effect on the appliance.",
				Validators: []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"email": schema.StringAttribute{
				Required:            true,
				Description:         "Email address of the first administrator.",
				MarkdownDescription: "Email address of the first administrator.",
				Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"first_name": schema.StringAttribute{
				Optional:            true,
				Description:         "First name of the first administrator.",
				MarkdownDescription: "First name of the first administrator.",
			},
			"last_name": schema.StringAttribute{
				Optional:            true,
				Description:         "Last name of the first administrator.",
				MarkdownDescription: "Last name of the first administrator.",
			},
			"hub_mode": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(hubModeSkip),
				Description: "How setup deals with Morpheus Hub: skip (the default) registers nothing; login and " +
					"register use the hub credentials already entered on the appliance. A licence key implies skip.",
				MarkdownDescription: "How setup deals with Morpheus Hub: `skip` (the default) registers nothing; `login` and " +
					"`register` use the hub credentials already entered on the appliance. A licence key implies `skip`.",
				Validators: []validator.String{
					stringvalidator.OneOf(hubModeSkip, hubModeLogin, hubModeRegister),
				},
			},
			"license_key_wo": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
				WriteOnly: true,
				Description: "Licence key to install during setup. The appliance validates it before creating anything, " +
					"so a bad key fails the setup and leaves the appliance untouched. Without a key the appliance " +
					"installs its default trial licence. Write-only; to manage the licence afterwards use " +
					"hpe_morpheus_license.",
				MarkdownDescription: "Licence key to install during setup. The appliance validates it before creating anything, " +
					"so a bad key fails the setup and leaves the appliance untouched. Without a key the appliance " +
					"installs its default trial licence. Write-only; to manage the licence afterwards use " +
					"`hpe_morpheus_license`.",
				Validators: []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"already_initialised": schema.BoolAttribute{
				Computed: true,
				Description: "Whether the appliance had already been set up when the resource was created, in which " +
					"case nothing was sent and the configured tenant and user may not be the ones on the appliance.",
				MarkdownDescription: "Whether the appliance had already been set up when the resource was created, in which " +
					"case nothing was sent and the configured tenant and user may not be the ones on the appliance.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"build_version": schema.StringAttribute{
				Computed: true,
				Description: "The appliance's build version, as reported after setup and on each refresh that can " +
					"authenticate with the provider's credentials; null until then.",
				MarkdownDescription: "The appliance's build version, as reported after setup and on each refresh that can " +
					"authenticate with the provider's credentials; null until then.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
			}),
		},
	}
}

const (
	hubModeSkip     = "skip"
	hubModeLogin    = "login"
	hubModeRegister = "register"
)
