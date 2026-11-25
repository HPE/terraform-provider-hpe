// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

// Package settingbackup implements hpe_morpheus_setting_backup, the Morpheus
// backup settings: whether backups are scheduled for provisioned instances,
// whether Morpheus configures instances for backups, whether the appliance
// database is backed up, the default backup storage bucket and schedule, and
// how many successful backups to retain.
//
// The settings are a singleton the API identifies by nothing: GET and PUT
// /api/backup-settings, no DELETE. Create and Update PUT the settings that
// change -- Morpheus binds just the keys present under backupSettings -- and
// read the result back, since the PUT answers with a success envelope only; a
// setting Morpheus accepts but does not apply fails the apply with its name
// (checkApplied). Delete forgets the resource and leaves the appliance as it
// is. The id is always "1", as it was in the SDKv2 resource, and import
// accepts any ID.
//
// Morpheus keeps the scheduled-backups and create-backups flags, the
// retention and the default schedule once per appliance (shared by every
// tenant with the permission), the default storage bucket once per tenant,
// and one appliance backup for the whole appliance that only the master
// tenant can manage; ModifyPlan refuses backup_appliance = true from any
// other tenant.
//
// retention_count and retention_days are two attributes over one setting
// (retentionCount): retention_days is the deprecated name the SDKv2 resource
// used, and it is a number of backups, not days (a misnomer). The two must not
// both be set (a ConfigValidator), and ModifyPlan mirrors whichever is
// configured into the other so the single setting has one source of truth.
//
// The resource replaces an SDKv2 implementation of the same type;
// upgrade_state.go converts the state that one wrote.
package settingbackup

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
// whose state (version 0) had integer ids stored as 0 when unset and the
// retention under retention_days only. UpgradeState converts it; see
// upgrade_state.go.
const schemaVersion = 1

// resourceType names the resource in diagnostics.
const resourceType = "setting_backup"

// singletonID is the id of the settings: always "1", as in the SDKv2
// resource. The settings are a singleton, so it identifies nothing.
const singletonID = "1"

// resourceModel is the resource's model: the generated one, and the id that
// resourceSchema adds to the generated schema.
type resourceModel struct {
	SettingBackupModel
	Id types.String `tfsdk:"id"`
}

// resourceSchema is the generated schema with the id of the SDKv2 resource
// added, and its version. The id is added here rather than in the codegen
// spec: the spec also generates the Ansible module, and its generator manages
// a resource whose spec has no id as a singleton (no state argument).
func resourceSchema(ctx context.Context) schema.Schema {
	s := SettingBackupResourceSchema(ctx)

	// The generated schema defaults to version 0. Bump it so Terraform invokes
	// UpgradeState when reading state written by the SDKv2 resource.
	s.Version = schemaVersion

	s.Attributes["id"] = schema.StringAttribute{
		Computed:            true,
		Description:         `The ID of the backup settings: always "1", as they are a singleton.`,
		MarkdownDescription: "The ID of the backup settings: always `\"1\"`, as they are a singleton.",
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

// ConfigValidators keeps retention_count and retention_days from both being
// set: they are two attributes over one setting (retentionCount), and sending
// one value is unambiguous. The rule lives here rather than in the codegen
// spec because the spec also generates the Ansible module, whose generator
// drops the deprecated retention_days attribute.
func (r *Resource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.Conflicting(
			path.MatchRoot("retention_count"),
			path.MatchRoot("retention_days"),
		),
	}
}
