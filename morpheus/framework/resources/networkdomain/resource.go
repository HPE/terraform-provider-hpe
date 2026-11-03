// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkdomain

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
)

// schemaVersion is the current framework schema version for the network domain
// resource. It is 1 (not 0) because this resource replaces the retired
// terraform-plugin-sdk/v2 implementation, whose state stored "id" as a string
// (the framework model uses int64). UpgradeState also carries the deprecated
// auto_join_domain and domain_password attributes, which remain accepted.
// UpgradeState (version 0) converts that legacy state to this schema. See
// upgrade_state.go.
const schemaVersion = 1

// passwordMask is the sentinel Morpheus uses for "password unchanged".
// NetworkDomainsController.update() does
//
//	if (opts.networkDomain?.domainPassword != '************')
//	    domain.domainPassword = opts.networkDomain?.domainPassword
//
// so OMITTING domainPassword (a missing key) assigns null and CLEARS the stored
// password. The framework therefore always sends domainPassword on update: the
// configured value when it is being re-sent (see wodigest.Resend), otherwise
// this mask, which the API reads as "leave it as it is".
const passwordMask = "************"

type Resource struct {
	configure.ResourceWithMorpheusConfigure
}

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

func (r *Resource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_" + "network_domain"
}

func (r *Resource) Schema(
	ctx context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	s := NetworkDomainResourceSchema(ctx)
	// The generated schema defaults to version 0. Bump it so Terraform invokes
	// UpgradeState when reading state written by the legacy SDKv2 resource.
	s.Version = schemaVersion
	resp.Schema = s
}
