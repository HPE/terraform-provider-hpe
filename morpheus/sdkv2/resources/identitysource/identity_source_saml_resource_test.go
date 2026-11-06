// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package identitysource_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/resources/identitysource"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

// Unlike Active Directory, a SAML identity source is validated field by
// field without contacting the identity provider, so the documented example
// applies against any appliance; only the tenant and the roles it names have
// to exist, and the test creates the roles.
func TestAccMorpheusIdentitySourceSAMLExampleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.MasterTenant)

	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())
	firstRoleName := name + "-first"
	secondRoleName := name + "-second"

	resourceConfig, err := identitysource.RenderIdentitySourceSAMLConfig(t, map[string]string{
		"Name":                 name,
		"DefaultAccountRoleId": "hpe_morpheus_role.first_role.id",
		"RoleMappingRoleId":    "hpe_morpheus_role.first_role.id",
		"RoleMappingRoleName":  firstRoleName,
		"RoleMapping2RoleId":   "hpe_morpheus_role.second_role.id",
		"RoleMapping2RoleName": secondRoleName,
	})
	if err != nil {
		t.Fatal(err)
	}
	config := providerConfig +
		userRoleConfig("first_role", firstRoleName) +
		userRoleConfig("second_role", secondRoleName) +
		resourceConfig

	const addr = "hpe_morpheus_identity_source_saml.tf_example_saml_identity_source"
	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(addr, "tenant_id", "1"),
		resource.TestCheckResourceAttr(addr, "name", name),
		resource.TestCheckResourceAttr(addr, "description", "TF example SAML identity source"),
		resource.TestCheckResourceAttr(addr, "login_redirect_url", "https://tfexamplesaml.test.local:8443/realms/master/protocol/saml"),
		resource.TestCheckResourceAttr(addr, "logout_redirect_url", "https://tfexamplesaml.test.local:8443/realms/master/protocol/saml"),
		resource.TestCheckResourceAttr(addr, "include_saml_request_parameter", "true"),
		resource.TestCheckResourceAttr(addr, "saml_request", "SelfSigned"),
		resource.TestCheckResourceAttr(addr, "validate_assertion_signature", "false"),
		resource.TestCheckResourceAttr(addr, "given_name_attribute", "givenName"),
		resource.TestCheckResourceAttr(addr, "surname_attribute", "surname"),
		resource.TestCheckResourceAttr(addr, "email_attribute", "email"),
		resource.TestCheckResourceAttrPair(addr, "default_account_role_id", "hpe_morpheus_role.first_role", "id"),
		resource.TestCheckResourceAttr(addr, "role_attribute_name", "memberOf"),
		resource.TestCheckResourceAttr(addr, "required_role_attribute_value", "test"),
		resource.TestCheckResourceAttr(addr, "enable_role_mapping_permission", "false"),
		resource.TestCheckResourceAttr(addr, "role_mapping.#", "2"),
		resource.TestCheckTypeSetElemAttrPair(addr, "role_mapping.*.role_id", "hpe_morpheus_role.first_role", "id"),
		resource.TestCheckTypeSetElemNestedAttrs(addr, "role_mapping.*", map[string]string{
			"role_name":           firstRoleName,
			"assertion_attribute": "developers",
		}),
		resource.TestCheckTypeSetElemAttrPair(addr, "role_mapping.*.role_id", "hpe_morpheus_role.second_role", "id"),
		resource.TestCheckTypeSetElemNestedAttrs(addr, "role_mapping.*", map[string]string{
			"role_name":           secondRoleName,
			"assertion_attribute": "operators",
		}),
	}

	// The Service Provider metadata Morpheus computes (MORPH-12634) is only
	// known after apply; the example surfaces it through outputs.
	applyChecks := append(checks,
		resource.TestCheckResourceAttrSet(addr, "entity_id"),
		resource.TestCheckResourceAttrSet(addr, "acs_url"),
		resource.TestCheckResourceAttrSet(addr, "sp_metadata"),
		resource.TestMatchOutput("saml_entity_id", regexp.MustCompile(`.+`)),
		resource.TestMatchOutput("saml_acs_url", regexp.MustCompile(`^https?://`)),
		resource.TestMatchOutput("saml_sp_metadata", regexp.MustCompile(`EntityDescriptor`)),
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			// Plan
			{
				Config:             config,
				ExpectNonEmptyPlan: true,
				Check:              resource.ComposeAggregateTestCheckFunc(checks...),
				PlanOnly:           true,
			},
			// Apply
			{
				Config: config,
				Check:  resource.ComposeAggregateTestCheckFunc(applyChecks...),
			},
			// Plan after apply
			{
				Config:             config,
				ExpectNonEmptyPlan: false,
				Check:              resource.ComposeAggregateTestCheckFunc(applyChecks...),
				PlanOnly:           true,
			},
		},
	})
}
