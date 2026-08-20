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
	// PCEIdentity is a list to stay compatible with the SDKv2 provider
	// that the "hpe" provider is muxed with: SDKv2 can only express nested
	// blocks as list, set or map nesting, never single nesting. The schema
	// limits it to at most one element.
	PCEIdentity []PCEIdentityModel `tfsdk:"pce_identity"`
	// PCEDisconnectedIdentity is the Disconnected PCE counterpart of
	// PCEIdentity, and is a list for the same reason. The two blocks are
	// mutually exclusive.
	PCEDisconnectedIdentity []PCEDisconnectedIdentityModel `tfsdk:"pce_disconnected_identity"`
}

type PCEIdentityModel struct {
	ClientID     types.String `tfsdk:"client_id"`
	ClientSecret types.String `tfsdk:"client_secret"`
	Location     types.String `tfsdk:"location"`
	Space        types.String `tfsdk:"space"`
	IssuerURL    types.String `tfsdk:"issuer_url"`
	IAMToken     types.String `tfsdk:"iam_token"`
	BrokerURL    types.String `tfsdk:"broker_url"`
}

// PCEDisconnectedIdentityModel is the Disconnected PCE identity block.
//
// It scopes by workspace_id where the Connected block uses space, as GLP
// identifies a workspace where GLCS identifies an IAM space. Location is
// required for both: the broker uses it to identify which deployment's Morpheus
// instance to return.
type PCEDisconnectedIdentityModel struct {
	ClientID     types.String `tfsdk:"client_id"`
	ClientSecret types.String `tfsdk:"client_secret"`
	// TokenIssuerURL is this block's spelling of the Connected block's
	// issuer_url: GLP names the same value differently. The two blocks are
	// configured independently, so the Connected spelling is left as it is.
	TokenIssuerURL types.String `tfsdk:"token_issuer_url"`
	IAMToken       types.String `tfsdk:"iam_token"`
	Location       types.String `tfsdk:"location"`
	WorkspaceID    types.String `tfsdk:"workspace_id"`
	BrokerURL      types.String `tfsdk:"broker_url"`
}
