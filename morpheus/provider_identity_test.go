// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package morpheus_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	providermod "github.com/HPE/terraform-provider-hpe/provider"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

// newHpeServer returns the muxed "hpe" provider's framework server, built the
// same way main.go builds it.
func newHpeServer(t *testing.T) tfprotov6.ProviderServer {
	t.Helper()

	server, err := providerserver.NewProtocol6WithError(
		providermod.New("test", adapter.NewMorpheus())(),
	)()
	if err != nil {
		t.Fatalf("could not create provider server: %v", err)
	}

	return server
}

// providerObjectType returns the tftypes.Object of the "hpe" provider schema.
func providerObjectType(t *testing.T, server tfprotov6.ProviderServer) tftypes.Object {
	t.Helper()

	resp, err := server.GetProviderSchema(
		context.Background(),
		&tfprotov6.GetProviderSchemaRequest{},
	)
	if err != nil {
		t.Fatalf("GetProviderSchema() unexpected error: %v", err)
	}

	obj, ok := resp.Provider.ValueType().(tftypes.Object)
	if !ok {
		t.Fatalf("provider value type = %T, want tftypes.Object", resp.Provider.ValueType())
	}

	return obj
}

// absentBlocks describes how an unconfigured block is represented on the wire.
//
// Terraform's HCL decoder returns cty.ListValEmpty for a block that is not
// present (hcldec BlockListSpec), so emptyBlocks is the faithful case. Both are
// exercised because the distinction is exactly what makes the framework's
// ConflictsWith validators unsuitable here: they treat only null as absent.
type absentBlocks int

const (
	nullBlocks absentBlocks = iota
	emptyBlocks
)

func (a absentBlocks) String() string {
	if a == emptyBlocks {
		return "empty list"
	}

	return "null list"
}

// object builds a value of the given object type, defaulting every attribute to
// null and applying the supplied overrides. Every attribute has to be present
// for tftypes to accept the value.
//
// Unset list attributes follow absent, so that a test can reproduce either wire
// representation of a block that was not written.
func object(
	t *testing.T,
	typ tftypes.Object,
	absent absentBlocks,
	set map[string]tftypes.Value,
) tftypes.Value {
	t.Helper()

	attrs := make(map[string]tftypes.Value, len(typ.AttributeTypes))

	for name, attrType := range typ.AttributeTypes {
		if v, ok := set[name]; ok {
			attrs[name] = v

			continue
		}

		if listType, ok := attrType.(tftypes.List); ok && absent == emptyBlocks {
			attrs[name] = tftypes.NewValue(listType, []tftypes.Value{})

			continue
		}

		attrs[name] = tftypes.NewValue(attrType, nil)
	}

	for name := range set {
		if _, ok := typ.AttributeTypes[name]; !ok {
			t.Fatalf("attribute %q is not in the schema", name)
		}
	}

	return tftypes.NewValue(typ, attrs)
}

// blockTypes returns the list type of a named block and its element object
// type.
func blockTypes(t *testing.T, parent tftypes.Object, name string) (tftypes.List, tftypes.Object) {
	t.Helper()

	listType, ok := parent.AttributeTypes[name].(tftypes.List)
	if !ok {
		t.Fatalf("%s type = %T, want tftypes.List", name, parent.AttributeTypes[name])
	}

	objType, ok := listType.ElementType.(tftypes.Object)
	if !ok {
		t.Fatalf("%s element type = %T, want tftypes.Object", name, listType.ElementType)
	}

	return listType, objType
}

// validateProviderConfig runs the provider's ValidateProviderConfig with a
// morpheus block built from the supplied attributes, and returns the
// diagnostics summaries and details.
func validateProviderConfig(
	t *testing.T,
	absent absentBlocks,
	morpheusAttrs func(t *testing.T, obj tftypes.Object, absent absentBlocks) map[string]tftypes.Value,
) []string {
	t.Helper()

	server := newHpeServer(t)
	providerType := providerObjectType(t, server)

	morpheusList, morpheusObj := blockTypes(t, providerType, "morpheus")

	config := object(t, providerType, absent, map[string]tftypes.Value{
		"morpheus": tftypes.NewValue(morpheusList, []tftypes.Value{
			object(t, morpheusObj, absent, morpheusAttrs(t, morpheusObj, absent)),
		}),
	})

	dv, err := tfprotov6.NewDynamicValue(providerType, config)
	if err != nil {
		t.Fatalf("could not build a dynamic value: %v", err)
	}

	resp, err := server.ValidateProviderConfig(
		context.Background(),
		&tfprotov6.ValidateProviderConfigRequest{Config: &dv},
	)
	if err != nil {
		t.Fatalf("ValidateProviderConfig() unexpected error: %v", err)
	}

	var out []string
	for _, d := range resp.Diagnostics {
		out = append(out, d.Summary+": "+d.Detail)
	}

	return out
}

// eachBlockRepresentation runs fn for both wire representations of an
// unconfigured block.
func eachBlockRepresentation(t *testing.T, fn func(t *testing.T, absent absentBlocks)) {
	t.Helper()

	for _, absent := range []absentBlocks{nullBlocks, emptyBlocks} {
		t.Run(absent.String(), func(t *testing.T) {
			fn(t, absent)
		})
	}
}

// identityBlock returns a single element list for a named identity block, with
// only the supplied attributes set.
func identityBlockValue(
	t *testing.T,
	morpheusObj tftypes.Object,
	absent absentBlocks,
	name string,
	set map[string]string,
) tftypes.Value {
	t.Helper()

	listType, objType := blockTypes(t, morpheusObj, name)

	attrs := make(map[string]tftypes.Value, len(set))
	for k, v := range set {
		attrs[k] = tftypes.NewValue(tftypes.String, v)
	}

	return tftypes.NewValue(listType, []tftypes.Value{object(t, objType, absent, attrs)})
}

// A morpheus block with connection details and no identity block is the
// ordinary case and must validate cleanly. If an absent block were treated as
// configured, this is what would break.
func TestValidateProviderConfigAcceptsDirectConnectionDetails(t *testing.T) {
	eachBlockRepresentation(t, func(t *testing.T, absent absentBlocks) {
		diags := validateProviderConfig(t, absent, func(_ *testing.T, _ tftypes.Object, _ absentBlocks) map[string]tftypes.Value {
			return map[string]tftypes.Value{
				"url":          tftypes.NewValue(tftypes.String, "https://morpheus.example.invalid"),
				"access_token": tftypes.NewValue(tftypes.String, "token"),
			}
		})

		if len(diags) != 0 {
			t.Errorf("ValidateProviderConfig() diagnostics = %v, want none", diags)
		}
	})
}

// An identity block on its own must validate cleanly too.
func TestValidateProviderConfigAcceptsIdentityBlockAlone(t *testing.T) {
	eachBlockRepresentation(t, func(t *testing.T, absent absentBlocks) {
		diags := validateProviderConfig(t, absent, func(t *testing.T, obj tftypes.Object, absent absentBlocks) map[string]tftypes.Value {
			return map[string]tftypes.Value{
				"pce_disconnected_identity": identityBlockValue(t, obj, absent, "pce_disconnected_identity",
					map[string]string{
						"iam_token":    "token",
						"location":     "site-a",
						"workspace_id": "workspace",
						"broker_url":   "https://broker.example.invalid",
					}),
			}
		})

		if len(diags) != 0 {
			t.Errorf("ValidateProviderConfig() diagnostics = %v, want none", diags)
		}
	})
}

// Declaring the conflict in the schema means it is reported during validation,
// so "terraform validate" catches it without configuring the provider.
func TestValidateProviderConfigRejectsIdentityBlockWithDirectURL(t *testing.T) {
	eachBlockRepresentation(t, func(t *testing.T, absent absentBlocks) {
		diags := validateProviderConfig(t, absent, func(t *testing.T, obj tftypes.Object, absent absentBlocks) map[string]tftypes.Value {
			return map[string]tftypes.Value{
				"url": tftypes.NewValue(tftypes.String, "https://morpheus.example.invalid"),
				"pce_identity": identityBlockValue(t, obj, absent, "pce_identity", map[string]string{
					"iam_token": "token",
					"location":  "site-a",
				}),
			}
		})

		want := `Attribute "morpheus[0].url" cannot be specified when ` +
			`"morpheus[0].pce_identity" is specified`

		if !containsDiag(diags, want) {
			t.Errorf("ValidateProviderConfig() diagnostics = %v, want one containing %q", diags, want)
		}
	})
}

// Both identity blocks together resolve different Morpheus instances.
//
// The conflict is declared on pce_identity only. Declaring it on both blocks
// would work, but each block's validator would report it, so the same mistake
// would be described twice. This asserts the single diagnostic to keep that
// decision from being undone by accident.
func TestValidateProviderConfigRejectsBothIdentityBlocks(t *testing.T) {
	eachBlockRepresentation(t, func(t *testing.T, absent absentBlocks) {
		diags := validateProviderConfig(t, absent, func(t *testing.T, obj tftypes.Object, absent absentBlocks) map[string]tftypes.Value {
			return map[string]tftypes.Value{
				"pce_identity": identityBlockValue(t, obj, absent, "pce_identity", map[string]string{
					"iam_token": "token",
					"location":  "site-a",
				}),
				"pce_disconnected_identity": identityBlockValue(t, obj, absent, "pce_disconnected_identity",
					map[string]string{
						"iam_token":    "token",
						"location":     "site-a",
						"workspace_id": "workspace",
						"broker_url":   "https://broker.example.invalid",
					}),
			}
		})

		want := `Attribute "morpheus[0].pce_disconnected_identity" cannot be ` +
			`specified when "morpheus[0].pce_identity" is specified`

		if !containsDiag(diags, want) {
			t.Errorf("ValidateProviderConfig() diagnostics = %v, want one containing %q", diags, want)
		}

		if len(diags) != 1 {
			t.Errorf("ValidateProviderConfig() reported %d diagnostics, want 1: %v", len(diags), diags)
		}
	})
}

// A value that is not known during validation, such as one taken from a
// variable, must not be reported as a conflict. Terraform validates again once
// the value is known.
func TestValidateProviderConfigIgnoresUnknownConnectionDetails(t *testing.T) {
	eachBlockRepresentation(t, func(t *testing.T, absent absentBlocks) {
		diags := validateProviderConfig(t, absent, func(t *testing.T, obj tftypes.Object, absent absentBlocks) map[string]tftypes.Value {
			return map[string]tftypes.Value{
				"url": tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
				"pce_identity": identityBlockValue(t, obj, absent, "pce_identity", map[string]string{
					"iam_token": "token",
					"location":  "site-a",
				}),
			}
		})

		if len(diags) != 0 {
			t.Errorf("ValidateProviderConfig() diagnostics = %v, want none for an unknown url", diags)
		}
	})
}

// identityBlockAttrs returns the attributes that make a named identity block
// valid on its own, so that a test can add to them to isolate one rule.
func identityBlockAttrs(block string) map[string]string {
	attrs := map[string]string{
		"location": "site-a",
	}

	if block == "pce_disconnected_identity" {
		attrs["workspace_id"] = "workspace"
		attrs["broker_url"] = "https://broker.example.invalid"
	}

	return attrs
}

// A pre-generated token and the credentials it would be generated from are two
// ways of obtaining the same thing, so configuring both is rejected. Each
// conflicting attribute is reported, rather than only the first.
func TestValidateProviderConfigRejectsIamTokenWithCredentials(t *testing.T) {
	for _, block := range []string{"pce_identity", "pce_disconnected_identity"} {
		t.Run(block, func(t *testing.T) {
			eachBlockRepresentation(t, func(t *testing.T, absent absentBlocks) {
				attrs := identityBlockAttrs(block)
				attrs["iam_token"] = "token"
				attrs["client_id"] = "client-id"
				attrs["client_secret"] = "client-secret"
				attrs["issuer_url"] = "https://issuer.example.invalid"

				diags := validateProviderConfig(t, absent,
					func(t *testing.T, obj tftypes.Object, absent absentBlocks) map[string]tftypes.Value {
						return map[string]tftypes.Value{
							block: identityBlockValue(t, obj, absent, block, attrs),
						}
					})

				for _, conflicting := range []string{"client_id", "client_secret", "issuer_url"} {
					want := fmt.Sprintf(
						`Attribute "morpheus[0].%s[0].%s" cannot be specified when `+
							`"morpheus[0].%s[0].iam_token" is specified`,
						block, conflicting, block,
					)

					if !containsDiag(diags, want) {
						t.Errorf("ValidateProviderConfig() diagnostics = %v, want one containing %q",
							diags, want)
					}
				}
			})
		})
	}
}

// The credentials are the ordinary way of obtaining a token, so using them
// without iam_token must not trip the conflict.
func TestValidateProviderConfigAcceptsCredentialsWithoutIamToken(t *testing.T) {
	for _, block := range []string{"pce_identity", "pce_disconnected_identity"} {
		t.Run(block, func(t *testing.T) {
			eachBlockRepresentation(t, func(t *testing.T, absent absentBlocks) {
				attrs := identityBlockAttrs(block)
				attrs["client_id"] = "client-id"
				attrs["client_secret"] = "client-secret"
				attrs["issuer_url"] = "https://issuer.example.invalid"

				diags := validateProviderConfig(t, absent,
					func(t *testing.T, obj tftypes.Object, absent absentBlocks) map[string]tftypes.Value {
						return map[string]tftypes.Value{
							block: identityBlockValue(t, obj, absent, block, attrs),
						}
					})

				if len(diags) != 0 {
					t.Errorf("ValidateProviderConfig() diagnostics = %v, want none", diags)
				}
			})
		})
	}
}

// Every field that supplies credentials is optional, so a block can satisfy the
// schema with only its required fields and no way to authenticate at all. That
// has to be caught during validation rather than surfacing later as an opaque
// token exchange failure.
func TestValidateProviderConfigRejectsIdentityBlockWithoutCredentials(t *testing.T) {
	for _, block := range []string{"pce_identity", "pce_disconnected_identity"} {
		t.Run(block, func(t *testing.T) {
			eachBlockRepresentation(t, func(t *testing.T, absent absentBlocks) {
				// Only the required fields: no credentials and no token.
				attrs := identityBlockAttrs(block)

				diags := validateProviderConfig(t, absent,
					func(t *testing.T, obj tftypes.Object, absent absentBlocks) map[string]tftypes.Value {
						return map[string]tftypes.Value{
							block: identityBlockValue(t, obj, absent, block, attrs),
						}
					})

				want := "Configure either the GreenLake API client credentials " +
					"(client_id, client_secret and issuer_url) or a pre-generated " +
					"iam_token."

				if !containsDiag(diags, want) {
					t.Errorf("ValidateProviderConfig() diagnostics = %v, want one containing %q",
						diags, want)
				}

				// The rule is declared once, on iam_token, so that it cannot
				// fire alongside the conflict rule on the same attribute.
				if len(diags) != 1 {
					t.Errorf("ValidateProviderConfig() reported %d diagnostics, want 1: %v",
						len(diags), diags)
				}
			})
		})
	}
}

func containsDiag(diags []string, want string) bool {
	for _, d := range diags {
		if strings.Contains(d, want) {
			return true
		}
	}

	return false
}
