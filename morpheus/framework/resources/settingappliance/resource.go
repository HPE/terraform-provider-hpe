// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

// Package settingappliance implements hpe_morpheus_setting_appliance, the
// appliance-wide settings: every setting PUT /api/appliance-settings takes --
// appliance URLs, API allowed origins, tenant registration, Docker privileged
// mode, password policy, browser sessions, user deactivation, SMTP, proxy,
// currency provider, Twilio, sync intervals, data retention, HTTP task hosts,
// Agent and SSH flags, locale, VDI gateway, option lists, dashboards and Murph
// -- bar the cloud-type commands and the plugin exchange URL.
//
// The settings are a singleton the API identifies by nothing: GET and PUT
// /api/appliance-settings, no DELETE. Create and Update PUT the settings that
// change -- Morpheus binds just the keys present -- and read the result back,
// since the PUT answers with a success envelope only. Delete forgets the
// resource and leaves the appliance as it is. The id is always "1", as it was
// in the SDKv2 resource, and import accepts any ID.
//
// Every PUT also resets the plugin exchange URL to its default: Morpheus does
// so whenever the body lacks exchangeUrl, and refuses a body that carries it.
// The provider never sends it; the documentation says so.
//
// The resource replaces an SDKv2 implementation of the same type;
// upgrade_state.go converts the state that one wrote.
package settingappliance

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
)

// schemaVersion is 1 because the resource replaces the SDKv2 implementation,
// whose state (version 0) had string ports and role ids, and the API's
// password hashes in the password attributes. UpgradeState converts it; see
// upgrade_state.go.
const schemaVersion = 1

// resourceType names the resource in diagnostics and in the private-state key
// of its write-only digests.
const resourceType = "setting_appliance"

// singletonID is the id of the settings: always "1", as in the SDKv2
// resource. The settings are a singleton, so it identifies nothing.
const singletonID = "1"

// resourceModel is the resource's model: the generated one, and the id that
// resourceSchema adds to the generated schema.
type resourceModel struct {
	SettingApplianceModel
	Id types.String `tfsdk:"id"`
}

// resourceSchema is the generated schema with the id of the SDKv2 resource
// added, and its version. The id is added here rather than in the codegen
// spec: the spec also generates the Ansible module, and its generator manages
// a resource whose spec has no id as a singleton (no state argument).
func resourceSchema(ctx context.Context) schema.Schema {
	s := SettingApplianceResourceSchema(ctx)

	// The generated schema defaults to version 0. Bump it so Terraform invokes
	// UpgradeState when reading state written by the SDKv2 resource.
	s.Version = schemaVersion

	s.Attributes["id"] = schema.StringAttribute{
		Computed:            true,
		Description:         `The ID of the appliance settings: always "1", as they are a singleton.`,
		MarkdownDescription: "The ID of the appliance settings: always `\"1\"`, as they are a singleton.",
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	}

	return s
}

type Resource struct {
	configure.ResourceWithMorpheusConfigure
}

var (
	_ resource.Resource                     = &Resource{}
	_ resource.ResourceWithConfigure        = &Resource{}
	_ resource.ResourceWithConfigValidators = &Resource{}
	_ resource.ResourceWithImportState      = &Resource{}
	_ resource.ResourceWithModifyPlan       = &Resource{}
	_ resource.ResourceWithUpgradeState     = &Resource{}
)

func NewResource() resource.Resource {
	return &Resource{}
}

func (r *Resource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_" + resourceType
}

func (r *Resource) Schema(
	ctx context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = resourceSchema(ctx)
}

// ConfigValidators keeps the SDKv2 resource's rule that the registration flag
// and the two default role ids are configured together (each was RequiredWith
// the other two). They travel together on the wire as well; see settingsBody.
// It also rejects an empty smtp_from_address; see smtpFromAddressValidator.
// The rules live here rather than in the codegen spec because the spec also
// generates the Ansible module, which lets each of the three travel alone.
func (r *Resource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.RequiredTogether(
			path.MatchRoot("registration_enabled"),
			path.MatchRoot("default_role_id"),
			path.MatchRoot("default_user_role_id"),
		),
		smtpFromAddressValidator{},
	}
}
