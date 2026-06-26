// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package model

import "github.com/hashicorp/terraform-plugin-framework/types"

// We require a separate package for model, as with the ClientFactory
// architecture we experience an import cycle if the Morpheus provider model
// exists in the "morpheus" provider package.
type MorpheusProviderModel struct {
	URL             types.String `tfsdk:"url"`
	Username        types.String `tfsdk:"username"`
	Password        types.String `tfsdk:"password"`
	AccessToken     types.String `tfsdk:"access_token"`
	TenantSubdomain types.String `tfsdk:"tenant_subdomain"`
	Insecure        types.Bool   `tfsdk:"insecure"`
}
