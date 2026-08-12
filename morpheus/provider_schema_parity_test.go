// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package morpheus_test

import (
	"context"
	"testing"

	fwprovider "github.com/hashicorp/terraform-plugin-framework/provider"
	sdkv2schema "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	sdkv2Morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	providermod "github.com/HPE/terraform-provider-hpe/provider"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

// frameworkDerivedSchema returns the SDKv2 schema that main.go derives from the
// framework provider schema. This is the schema the SDKv2 provider actually
// serves in production, so it is what has to agree with the framework side for
// the mux to accept both providers.
//
// It is built the same way main.go builds it, from the top level "hpe"
// provider, because that is what wraps the Morpheus sub-provider schema in the
// "morpheus" block.
func frameworkDerivedSchema(t *testing.T) map[string]*sdkv2schema.Schema {
	t.Helper()

	p := providermod.New("test", adapter.NewMorpheus())

	resp := &fwprovider.SchemaResponse{}
	p().Schema(context.Background(), fwprovider.SchemaRequest{}, resp)

	return convert.FwToSdkv2SchemaMap(resp.Schema)
}

// identityBlockSchema returns the element schema of a named identity block.
func identityBlockSchema(
	t *testing.T,
	in map[string]*sdkv2schema.Schema,
	name string,
) *sdkv2schema.Schema {
	t.Helper()

	block, ok := in[name]
	if !ok {
		t.Fatalf("%s block not found in schema", name)
	}

	return block
}

// The SDKv2 MaxItems is derived by reflecting on the concrete type name of the
// framework validators, so a validator that wraps or replaces
// listvalidator.SizeAtMost would silently drop the constraint and leave the two
// providers disagreeing about how many blocks are allowed.
func TestIdentityBlocksDeriveMaxItems(t *testing.T) {
	derived := frameworkDerivedSchema(t)

	morpheusBlock := identityBlockSchema(t, derived, "morpheus")

	elem, ok := morpheusBlock.Elem.(*sdkv2schema.Resource)
	if !ok {
		t.Fatalf("morpheus block Elem = %T, want *schema.Resource", morpheusBlock.Elem)
	}

	for _, name := range []string{"pce_identity", "pce_disconnected_identity"} {
		block := identityBlockSchema(t, elem.Schema, name)

		if got, want := block.Type, sdkv2schema.TypeList; got != want {
			t.Errorf("%s Type = %v, want %v", name, got, want)
		}

		if got, want := block.MaxItems, 1; got != want {
			t.Errorf("%s MaxItems = %d, want %d", name, got, want)
		}

		if !block.Optional {
			t.Errorf("%s Optional = false, want true: the block itself is never required", name)
		}
	}
}

// Both identity blocks need a location: the broker resolves it to a service
// instance and to the zone that the returned token's roles are granted against,
// and rejects a request without one. The disconnected block additionally has no
// HPE hosted broker to fall back to and no implicit workspace.
func TestIdentityRequiredAttributesSurviveConversion(t *testing.T) {
	tests := map[string]map[string]bool{
		"pce_identity": {
			"location":      true,
			"client_id":     false,
			"client_secret": false,
			"space":         false,
			"issuer_url":    false,
			"iam_token":     false,
			"broker_url":    false,
		},
		"pce_disconnected_identity": {
			"broker_url":    true,
			"location":      true,
			"workspace_id":  true,
			"client_id":     false,
			"client_secret": false,
			"issuer_url":    false,
			"iam_token":     false,
		},
	}

	derived := frameworkDerivedSchema(t)

	morpheusBlock := identityBlockSchema(t, derived, "morpheus")

	elem, ok := morpheusBlock.Elem.(*sdkv2schema.Resource)
	if !ok {
		t.Fatalf("morpheus block Elem = %T, want *schema.Resource", morpheusBlock.Elem)
	}

	for blockName, required := range tests {
		t.Run(blockName, func(t *testing.T) {
			block := identityBlockSchema(t, elem.Schema, blockName)

			blockElem, ok := block.Elem.(*sdkv2schema.Resource)
			if !ok {
				t.Fatalf("%s Elem = %T, want *schema.Resource", blockName, block.Elem)
			}

			for name, wantRequired := range required {
				attr, ok := blockElem.Schema[name]
				if !ok {
					t.Errorf("%s.%s not found", blockName, name)

					continue
				}

				if attr.Required != wantRequired {
					t.Errorf("%s.%s Required = %v, want %v",
						blockName, name, attr.Required, wantRequired)
				}

				// Required and Optional are mutually exclusive in SDKv2, and
				// InternalValidate rejects a schema that sets both.
				if attr.Required && attr.Optional {
					t.Errorf("%s.%s is both Required and Optional", blockName, name)
				}
			}
		})
	}
}

// The hand written SDKv2 schema is what tests exercise, while production serves
// the framework derived one. They must not drift apart.
func TestHandWrittenSdkv2SchemaMatchesDerived(t *testing.T) {
	derived := frameworkDerivedSchema(t)

	handWritten := sdkv2Morpheus.Provider().Schema

	derivedMorpheus, ok := derived["morpheus"].Elem.(*sdkv2schema.Resource)
	if !ok {
		t.Fatal("derived morpheus block Elem is not a *schema.Resource")
	}

	handMorpheus, ok := handWritten["morpheus"].Elem.(*sdkv2schema.Resource)
	if !ok {
		t.Fatal("hand written morpheus block Elem is not a *schema.Resource")
	}

	for _, name := range []string{"pce_identity", "pce_disconnected_identity"} {
		gotBlock := identityBlockSchema(t, handMorpheus.Schema, name)
		wantBlock := identityBlockSchema(t, derivedMorpheus.Schema, name)

		if gotBlock.MaxItems != wantBlock.MaxItems {
			t.Errorf("%s MaxItems: hand written = %d, derived = %d",
				name, gotBlock.MaxItems, wantBlock.MaxItems)
		}

		if gotBlock.Description != wantBlock.Description {
			t.Errorf("%s Description:\n hand written = %q\n derived      = %q",
				name, gotBlock.Description, wantBlock.Description)
		}

		gotElem, ok := gotBlock.Elem.(*sdkv2schema.Resource)
		if !ok {
			t.Fatalf("hand written %s Elem is not a *schema.Resource", name)
		}

		wantElem, ok := wantBlock.Elem.(*sdkv2schema.Resource)
		if !ok {
			t.Fatalf("derived %s Elem is not a *schema.Resource", name)
		}

		compareAttrs(t, name, gotElem.Schema, wantElem.Schema)
	}
}

// compareAttrs reports the differences that the mux cares about between a hand
// written and a derived block schema.
func compareAttrs(t *testing.T, block string, got, want map[string]*sdkv2schema.Schema) {
	t.Helper()

	for name, wantAttr := range want {
		gotAttr, ok := got[name]
		if !ok {
			t.Errorf("%s.%s missing from the hand written schema", block, name)

			continue
		}

		if gotAttr.Type != wantAttr.Type {
			t.Errorf("%s.%s Type: hand written = %v, derived = %v",
				block, name, gotAttr.Type, wantAttr.Type)
		}

		if gotAttr.Required != wantAttr.Required {
			t.Errorf("%s.%s Required: hand written = %v, derived = %v",
				block, name, gotAttr.Required, wantAttr.Required)
		}

		if gotAttr.Optional != wantAttr.Optional {
			t.Errorf("%s.%s Optional: hand written = %v, derived = %v",
				block, name, gotAttr.Optional, wantAttr.Optional)
		}

		if gotAttr.Sensitive != wantAttr.Sensitive {
			t.Errorf("%s.%s Sensitive: hand written = %v, derived = %v",
				block, name, gotAttr.Sensitive, wantAttr.Sensitive)
		}

		if gotAttr.Description != wantAttr.Description {
			t.Errorf("%s.%s Description:\n hand written = %q\n derived      = %q",
				block, name, gotAttr.Description, wantAttr.Description)
		}
	}

	for name := range got {
		if _, ok := want[name]; !ok {
			t.Errorf("%s.%s is in the hand written schema but not the derived one", block, name)
		}
	}
}

// InternalValidate is what SDKv2 runs over a provider schema at startup. The
// hand written schema is served by tests, so a schema it rejects would only
// surface as a confusing failure much later.
func TestHandWrittenSdkv2SchemaIsInternallyValid(t *testing.T) {
	if err := sdkv2Morpheus.Provider().InternalValidate(); err != nil {
		t.Fatalf("InternalValidate() = %v, want nil", err)
	}
}
