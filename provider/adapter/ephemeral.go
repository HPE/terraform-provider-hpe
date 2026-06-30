// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package adapter

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/provider"
)

// EphemeralResourceAdapter wraps a Terraform Plugin Framework ephemeral
// resource and adapts it to work as a child ephemeral resource within a parent
// provider architecture. It allows child provider ephemeral resources to be
// properly namespaced and integrated into a parent provider.
//
// Key transformations performed by the adapter:
//
// 1. Metadata: Transforms the ephemeral resource TypeName by prepending the
// child provider's TypeName.
//
// 2. Configure: Extracts the child provider's configuration data from the
// parent provider's ConfigureRequest.ProviderData map.
type EphemeralResourceAdapter struct {
	in       ephemeral.EphemeralResource
	provider provider.Provider

	withConfigure        ephemeral.EphemeralResourceWithConfigure
	withConfigValidators ephemeral.EphemeralResourceWithConfigValidators
	withRenew            ephemeral.EphemeralResourceWithRenew
	withClose            ephemeral.EphemeralResourceWithClose
	withValidateConfig   ephemeral.EphemeralResourceWithValidateConfig
}

var (
	_ ephemeral.EphemeralResource                 = &EphemeralResourceAdapter{}
	_ ephemeral.EphemeralResourceWithConfigure    = &EphemeralResourceAdapter{}
	_ ephemeral.EphemeralResourceWithConfigValidators = &EphemeralResourceAdapter{}
	_ ephemeral.EphemeralResourceWithRenew        = &EphemeralResourceAdapter{}
	_ ephemeral.EphemeralResourceWithClose        = &EphemeralResourceAdapter{}
	_ ephemeral.EphemeralResourceWithValidateConfig = &EphemeralResourceAdapter{}
)

func NewEphemeralResourceAdapter(in ephemeral.EphemeralResource, p provider.Provider) *EphemeralResourceAdapter {
	e := &EphemeralResourceAdapter{in: in, provider: p}

	e.withConfigure, _ = in.(ephemeral.EphemeralResourceWithConfigure)
	e.withConfigValidators, _ = in.(ephemeral.EphemeralResourceWithConfigValidators)
	e.withRenew, _ = in.(ephemeral.EphemeralResourceWithRenew)
	e.withClose, _ = in.(ephemeral.EphemeralResourceWithClose)
	e.withValidateConfig, _ = in.(ephemeral.EphemeralResourceWithValidateConfig)

	return e
}

func NewAdaptedEphemeralResource(in ephemeral.EphemeralResource, p provider.Provider) ephemeral.EphemeralResource {
	return NewEphemeralResourceAdapter(in, p)
}

// Metadata transforms the ephemeral resource name by prepending the child
// provider's TypeName.
func (e *EphemeralResourceAdapter) Metadata(
	ctx context.Context,
	req ephemeral.MetadataRequest,
	resp *ephemeral.MetadataResponse,
) {
	providerMetaResp := &provider.MetadataResponse{}
	e.provider.Metadata(ctx, provider.MetadataRequest{}, providerMetaResp)

	req.ProviderTypeName = req.ProviderTypeName + "_" + providerMetaResp.TypeName
	e.in.Metadata(ctx, req, resp)
}

func (e *EphemeralResourceAdapter) Schema(
	ctx context.Context,
	req ephemeral.SchemaRequest,
	resp *ephemeral.SchemaResponse,
) {
	e.in.Schema(ctx, req, resp)
}

func (e *EphemeralResourceAdapter) Open(
	ctx context.Context,
	req ephemeral.OpenRequest,
	resp *ephemeral.OpenResponse,
) {
	e.in.Open(ctx, req, resp)
}

func (e *EphemeralResourceAdapter) Configure(
	ctx context.Context,
	req ephemeral.ConfigureRequest,
	resp *ephemeral.ConfigureResponse,
) {
	if e.withConfigure == nil {
		return
	}

	// Extract child provider configure data for ConfigureRequest.ProviderData
	if providerData, ok := req.ProviderData.(map[string]any); ok {
		metaResp := &provider.MetadataResponse{}
		e.provider.Metadata(ctx, provider.MetadataRequest{}, metaResp)

		childData, exists := providerData[metaResp.TypeName]
		if !exists {
			resp.Diagnostics.AddError(
				"Missing provider configuration",
				fmt.Sprintf(
					"The %q provider block is required but was not found in the provider configuration.",
					metaResp.TypeName,
				),
			)

			return
		}

		req.ProviderData = childData
	}

	e.withConfigure.Configure(ctx, req, resp)
}

func (e *EphemeralResourceAdapter) ConfigValidators(
	ctx context.Context,
) []ephemeral.ConfigValidator {
	if e.withConfigValidators == nil {
		return nil
	}

	return e.withConfigValidators.ConfigValidators(ctx)
}

func (e *EphemeralResourceAdapter) Renew(
	ctx context.Context,
	req ephemeral.RenewRequest,
	resp *ephemeral.RenewResponse,
) {
	if e.withRenew == nil {
		return
	}

	e.withRenew.Renew(ctx, req, resp)
}

func (e *EphemeralResourceAdapter) Close(
	ctx context.Context,
	req ephemeral.CloseRequest,
	resp *ephemeral.CloseResponse,
) {
	if e.withClose == nil {
		return
	}

	e.withClose.Close(ctx, req, resp)
}

func (e *EphemeralResourceAdapter) ValidateConfig(
	ctx context.Context,
	req ephemeral.ValidateConfigRequest,
	resp *ephemeral.ValidateConfigResponse,
) {
	if e.withValidateConfig == nil {
		return
	}

	e.withValidateConfig.ValidateConfig(ctx, req, resp)
}
