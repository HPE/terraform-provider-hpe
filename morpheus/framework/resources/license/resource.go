// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package license

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
)

// schemaVersion is the current framework schema version for the license
// resource. It is 1 (not 0) because this resource replaces the retired
// terraform-plugin-sdk/v2 implementation, whose state version 0 stored the
// licence key in state as a plain (Sensitive) string. UpgradeState (version 0)
// drops that plaintext key: the key is now write-only and never stored. See
// upgrade_state.go.
const schemaVersion = 1

// Operation names for diagnostics.
const (
	createOperation = "create license"
	readOperation   = "read license"
	updateOperation = "update license"
	deleteOperation = "delete license"
)

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
	resp.TypeName = req.ProviderTypeName + "_" + "license"
}

func (r *Resource) Schema(
	ctx context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	s := LicenseResourceSchema(ctx)
	// The generated schema defaults to version 0. Bump it so Terraform invokes
	// UpgradeState when reading state written by the legacy SDKv2 resource.
	s.Version = schemaVersion
	resp.Schema = s
}
