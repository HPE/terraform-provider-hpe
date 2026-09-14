// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package tenant

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
)

// schemaVersion is the current framework schema version for the tenant
// resource. It is 1 (not 0) because this resource replaces the retired
// terraform-plugin-sdk/v2 implementation, whose state stored "id" as a string.
// UpgradeState (version 0) converts that legacy state to the int64 "id" this
// schema uses. See upgrade_state.go.
const schemaVersion = 1

// Master-tenant guard messages. Morpheus lets a master-tenant caller update
// the master account's name, description, currency, subdomain and billing
// fields, but rejects disabling it (AccountsService.validateUpdate) with an
// opaque 400, ignores any role assigned to it (the master account has no base
// role), and rejects deleting it. Surface each as a named diagnostic instead
// (see Update and Delete).
const (
	masterTenantDisableMsg = "The master tenant cannot be disabled. Morpheus rejects setting " +
		"enabled = false on the master account; leave enabled unset or true."
	masterTenantRoleMsg = "The master tenant has no base role. Morpheus ignores a role assigned " +
		"to the master account, so base_role_id cannot be set on it; remove base_role_id from " +
		"the configuration."
	masterTenantDeleteMsg = "The master tenant cannot be deleted. Morpheus rejects deletion of the " +
		"master account; remove it from Terraform state (terraform state rm) if it was imported."
)

// parentIdFeature names the parent_id attribute in the appliance version gate
// diagnostic, phrased as a plural noun so it reads as "... require a Morpheus
// appliance version >= 8.1.0". See constants.TenantParentMinVersion.
const parentIdFeature = "Nominated parent tenants (parent_id)"

var (
	_ resource.Resource                 = &Resource{}
	_ resource.ResourceWithConfigure    = &Resource{}
	_ resource.ResourceWithImportState  = &Resource{}
	_ resource.ResourceWithModifyPlan   = &Resource{}
	_ resource.ResourceWithUpgradeState = &Resource{}
)

func NewResource() resource.Resource {
	return &Resource{}
}

type Resource struct {
	configure.ResourceWithMorpheusConfigure
}

func (r *Resource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_" + "tenant"
}

func (r *Resource) Schema(
	ctx context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	s := TenantResourceSchema(ctx)
	// The generated schema defaults to version 0. Bump it so Terraform invokes
	// UpgradeState when reading state written by the legacy SDKv2 resource.
	s.Version = schemaVersion
	resp.Schema = s
}
