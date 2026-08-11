// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

// Package validators provides schema validators shared by the Morpheus
// provider.
package validators

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// identityCredentialAttributes are the attributes of an identity block that can
// supply a way of authenticating with GreenLake. At least one has to be set.
var identityCredentialAttributes = []string{
	"client_id",
	"client_secret",
	"issuer_url",
	"iam_token",
}

// IdentityCredentialsOrTokenValidator checks that an identity block configures
// some way of authenticating with GreenLake, being either API client
// credentials or a pre-generated token.
//
// Every credential attribute is optional, because either alternative is
// allowed. Without this rule a block satisfies the schema with only its
// required fields and fails much later, in Configure, with an opaque token
// exchange error.
//
// stringvalidator.AtLeastOneOf expresses the same rule, but builds its message
// from unresolved path expressions, so a relative sibling renders as
// "pce_identity[0].iam_token.<.client_id". This reports the block itself and
// names the two alternatives in prose instead.
func IdentityCredentialsOrTokenValidator() validator.Object {
	return identityCredentialsOrTokenValidator{}
}

type identityCredentialsOrTokenValidator struct{}

var _ validator.Object = identityCredentialsOrTokenValidator{}

func (v identityCredentialsOrTokenValidator) Description(_ context.Context) string {
	return "must configure either GreenLake API client credentials or a pre-generated IAM token"
}

func (v identityCredentialsOrTokenValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v identityCredentialsOrTokenValidator) ValidateObject(
	_ context.Context,
	req validator.ObjectRequest,
	resp *validator.ObjectResponse,
) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	attributes := req.ConfigValue.Attributes()

	for _, name := range identityCredentialAttributes {
		value, ok := attributes[name]
		if !ok {
			continue
		}

		// A configured value satisfies the rule. An unknown value is not null
		// either, which defers the decision until Terraform has resolved it
		// rather than reporting a configuration the user may well have
		// supplied.
		if !value.IsNull() {
			return
		}
	}

	resp.Diagnostics.AddAttributeError(
		req.Path,
		"Missing GreenLake credentials",
		"Configure either the GreenLake API client credentials (client_id, "+
			"client_secret and issuer_url) or a pre-generated iam_token.",
	)
}
