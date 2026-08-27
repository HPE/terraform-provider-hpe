// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package adapter

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
)

// FunctionAdapter wraps a Terraform Plugin Framework function and adapts it
// to work as a child function within a parent provider architecture. It allows
// child provider functions to be properly namespaced and integrated into a
// parent provider.
//
// Key transformation performed by the adapter:
//
// Metadata: Transforms the function name by prepending the child provider's
// TypeName. For example, a "parse_cidr" function from a "morpheus" provider
// becomes "morpheus_parse_cidr", and when used in a parent "hpe" provider is
// callable as provider::hpe::morpheus_parse_cidr.
type FunctionAdapter struct {
	in       function.Function
	provider provider.Provider
}

var _ function.Function = &FunctionAdapter{}

func NewFunctionAdapter(in function.Function, p provider.Provider) *FunctionAdapter {
	return &FunctionAdapter{in: in, provider: p}
}

func NewAdaptedFunction(in function.Function, p provider.Provider) function.Function {
	return NewFunctionAdapter(in, p)
}

// Metadata transforms the function name by prepending the child provider's
// TypeName.
// Due to function namespacing limitations within Terraform Core and HCL,
// we cannot namespace functions for child providers with double colons (::).
// So we use the same pattern as in resource and data source namespacing.
func (f *FunctionAdapter) Metadata(
	ctx context.Context,
	req function.MetadataRequest,
	resp *function.MetadataResponse,
) {
	providerMetaResp := &provider.MetadataResponse{}
	f.provider.Metadata(ctx, provider.MetadataRequest{}, providerMetaResp)

	f.in.Metadata(ctx, req, resp)
	resp.Name = providerMetaResp.TypeName + "_" + resp.Name
}

func (f *FunctionAdapter) Definition(
	ctx context.Context,
	req function.DefinitionRequest,
	resp *function.DefinitionResponse,
) {
	f.in.Definition(ctx, req, resp)
}

func (f *FunctionAdapter) Run(
	ctx context.Context,
	req function.RunRequest,
	resp *function.RunResponse,
) {
	f.in.Run(ctx, req, resp)
}
