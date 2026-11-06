// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package identitysource_test

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/resources/identitysource"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

// userRoleConfig is an hpe_morpheus_role of type user, for an identity source
// to map to. The identity source tests create the roles they name rather than
// depending on roles an appliance happens to seed: the default roles' ids
// differ between appliances, and a missing one fails the apply.
func userRoleConfig(label, name string) string {
	return fmt.Sprintf(`
resource "hpe_morpheus_role" %q {
  name        = %q
  multitenant = false
  description = "Role for the identity source acceptance test"
  role_type   = "user"
}
`, label, name)
}

// checkRoleMappingAttrSet passes when every role_mapping element of the
// resource has a non-empty value for attribute. TestCheckTypeSetElemNestedAttrs
// needs the values, which for an attribute Morpheus computes are not known in
// advance.
func checkRoleMappingAttrSet(addr, attribute string) resource.TestCheckFunc {
	key := regexp.MustCompile(`^role_mapping\.\d+\.` + regexp.QuoteMeta(attribute) + `$`)

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[addr]
		if !ok {
			return fmt.Errorf("not found: %s", addr)
		}
		found := 0
		for k, v := range rs.Primary.Attributes {
			if !key.MatchString(k) {
				continue
			}
			found++
			if v == "" {
				return fmt.Errorf("%s: %s is empty", addr, k)
			}
		}
		if found == 0 {
			return fmt.Errorf("%s: no role_mapping has %s", addr, attribute)
		}

		return nil
	}
}

// activeDirectoryVariables declares the Terraform variables the test
// configuration takes the directory's details from. Terraform reads their
// values from the TF_VAR_testacc_active_directory_* environment itself — the
// convention the TF_VAR_ prefix signals (testhelpers/capabilities/README.md;
// testhelpers.ProviderBlockPceIdentity does the same) — so the values are
// never spliced into HCL text: a DOMAIN\user binding name, a password with a
// quote, backslash or ${, or a distinguished name with an escaped comma all
// reach the provider intact, and the password is not written into the
// rendered configuration. The password variable is sensitive so that plan
// output redacts it. The group FQN is declared only when the test has a
// value for it; without one the configuration leaves the attribute to
// Morpheus.
func activeDirectoryVariables(withGroupFqn bool) string {
	var b strings.Builder
	for _, name := range []string{"server", "domain", "binding_username", "group_name"} {
		fmt.Fprintf(&b, "variable \"testacc_active_directory_%s\" {}\n\n", name)
	}
	b.WriteString("variable \"testacc_active_directory_binding_password\" {\n  sensitive = true\n}\n\n")
	if withGroupFqn {
		b.WriteString("variable \"testacc_active_directory_group_fqn\" {}\n\n")
	}

	return b.String()
}

// Morpheus validates an Active Directory identity source against the
// directory when it is saved (ActiveDirectoryUserService.validateUserSource):
// it binds as the binding user, then looks the required group and every mapped
// group up, and rejects the save if any of that fails. It also rewrites each
// mapping's group name and FQN to the group's sAMAccountName and distinguished
// name as the directory reports them. So the documented example's placeholder
// directory cannot be applied, and a mapped group cannot be made up: the test
// takes a reachable directory and a group in it from the environment:
//
//	TF_VAR_testacc_active_directory_server            e.g. dc01.example.com
//	TF_VAR_testacc_active_directory_domain            e.g. example.com
//	TF_VAR_testacc_active_directory_binding_username  a user that can bind
//	TF_VAR_testacc_active_directory_binding_password  that user's password
//	TF_VAR_testacc_active_directory_group_name        a group to map a role to, spelt as its sAMAccountName
//	TF_VAR_testacc_active_directory_group_fqn         optional: that group's distinguished name
//
// and skips when any of the first five is unset. The configuration reads
// them as Terraform variables (activeDirectoryVariables); the test reads the
// same environment to decide whether to run and to know what to expect in
// state. When the FQN is given it is configured and asserted; when it is
// not, the configuration leaves the (Computed) attribute to Morpheus and the
// test only asserts that it was filled in. The required group is the
// directory's built-in Administrators group, which every Active Directory
// has. The Morpheus roles the identity source names are created by the test.
func TestAccMorpheusIdentitySourceActiveDirectoryExampleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.LDAP, capabilities.MasterTenant)

	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	directory := map[string]string{
		"AdServer":             os.Getenv("TF_VAR_testacc_active_directory_server"),
		"Domain":               os.Getenv("TF_VAR_testacc_active_directory_domain"),
		"BindingUsername":      os.Getenv("TF_VAR_testacc_active_directory_binding_username"),
		"BindingPassword":      os.Getenv("TF_VAR_testacc_active_directory_binding_password"),
		"RoleMappingGroupName": os.Getenv("TF_VAR_testacc_active_directory_group_name"),
	}
	for key, value := range directory {
		if value == "" {
			t.Skipf("Skipping: set TF_VAR_testacc_active_directory_{server,domain,binding_username,binding_password,group_name} "+
				"to a reachable Active Directory and a group in it to run this test (%s is unset)", key)
		}
	}
	groupFqn := os.Getenv("TF_VAR_testacc_active_directory_group_fqn")

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())
	defaultRoleName := name + "-default"
	mappedRoleName := name + "-mapped"

	// The directory's details are passed as variable references, not values:
	// see activeDirectoryVariables.
	overrides := map[string]string{ //nolint:gosec // G101: the password entry is a Terraform variable reference, not a credential
		"Name":                 name,
		"AdServer":             "var.testacc_active_directory_server",
		"Domain":               "var.testacc_active_directory_domain",
		"BindingUsername":      "var.testacc_active_directory_binding_username",
		"BindingPassword":      "var.testacc_active_directory_binding_password",
		"RoleMappingGroupName": "var.testacc_active_directory_group_name",
		"RoleMappingGroupFqn":  "",
		"DefaultAccountRoleId": "hpe_morpheus_role.default_role.id",
		"RoleMappingRoleId":    "hpe_morpheus_role.mapped_role.id",
		"RoleMappingRoleName":  mappedRoleName,
	}
	if groupFqn != "" {
		overrides["RoleMappingGroupFqn"] = "var.testacc_active_directory_group_fqn"
	}
	resourceConfig, err := identitysource.RenderIdentitySourceActiveDirectoryConfig(t, overrides)
	if err != nil {
		t.Fatal(err)
	}
	config := providerConfig +
		activeDirectoryVariables(groupFqn != "") +
		userRoleConfig("default_role", defaultRoleName) +
		userRoleConfig("mapped_role", mappedRoleName) +
		resourceConfig

	const addr = "hpe_morpheus_identity_source_active_directory.tf_example_active_directory_identity_source"
	mapping := map[string]string{
		"role_name":                   mappedRoleName,
		"active_directory_group_name": directory["RoleMappingGroupName"],
	}
	fqnCheck := checkRoleMappingAttrSet(addr, "active_directory_group_fqn")
	if groupFqn != "" {
		mapping["active_directory_group_fqn"] = groupFqn
	}
	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(addr, "tenant_id", "1"),
		resource.TestCheckResourceAttr(addr, "name", name),
		resource.TestCheckResourceAttr(addr, "description", "TF example AD identity source"),
		resource.TestCheckResourceAttr(addr, "ad_server", directory["AdServer"]),
		resource.TestCheckResourceAttr(addr, "domain", directory["Domain"]),
		resource.TestCheckResourceAttr(addr, "use_ssl", "false"),
		resource.TestCheckResourceAttr(addr, "binding_username", directory["BindingUsername"]),
		resource.TestCheckResourceAttr(addr, "required_group", "administrators"),
		resource.TestCheckResourceAttr(addr, "search_member_groups", "true"),
		resource.TestCheckResourceAttrPair(addr, "default_account_role_id", "hpe_morpheus_role.default_role", "id"),
		resource.TestCheckResourceAttr(addr, "role_mapping.#", "1"),
		resource.TestCheckTypeSetElemAttrPair(addr, "role_mapping.*.role_id", "hpe_morpheus_role.mapped_role", "id"),
		resource.TestCheckTypeSetElemNestedAttrs(addr, "role_mapping.*", mapping),
		fqnCheck,
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			// Plan
			{
				Config:             config,
				ExpectNonEmptyPlan: true,
				Check:              checkFn,
				PlanOnly:           true,
			},
			// Apply
			{
				Config: config,
				Check:  checkFn,
			},
			// Plan after apply
			{
				Config:             config,
				ExpectNonEmptyPlan: false,
				Check:              checkFn,
				PlanOnly:           true,
			},
		},
	})
}
