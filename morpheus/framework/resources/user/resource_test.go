// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package user_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/role"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/user"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func checkRole(
	resourceName string,
	roleIDAttr string,
	expectedRoles *map[string]struct{},
) func(*terraform.State) error {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		roleID := rs.Primary.Attributes[roleIDAttr]
		if _, ok := (*expectedRoles)[roleID]; !ok {
			return fmt.Errorf("role ID %s not found ", roleID)
		}

		delete(*expectedRoles, roleID)

		return nil
	}
}

func checkStrayRoles(
	expectedRoles *map[string]struct{},
) func(*terraform.State) error {
	return func(_ *terraform.State) error {
		if len(*expectedRoles) != 0 {
			return fmt.Errorf("not all role_ids found %s", *expectedRoles)
		}

		return nil
	}
}

// MORPH-16460: roleResourceConfig returns HCL for a per-run user role with the
// given resource label and name, replacing the previously hardcoded
// master-tenant role ids (1/3). Two-role fixtures use two labels so the
// created role ids can be correlated against the user's role_ids set.
func roleResourceConfig(label, name string) string {
	return `
resource "hpe_morpheus_role" "` + label + `" {
  name        = "` + name + `"
  multitenant = false
  description = "An example user role"
  role_type   = "user"
}
`
}

// captureRoleIDs populates the expected-role set with the actual ids of the
// created role resources so the set-correlation checks (checkRole /
// checkStrayRoles) can run without hardcoded ids. It preserves the
// post-apply role_ids set-correlation behaviour (MORPH-13857).
func captureRoleIDs(
	expectedRoles *map[string]struct{},
	roleResources ...string,
) func(*terraform.State) error {
	return func(s *terraform.State) error {
		m := map[string]struct{}{}
		for _, r := range roleResources {
			rs := s.RootModule().Resources[r]
			if rs == nil {
				return fmt.Errorf("role resource not found: %s", r)
			}
			m[rs.Primary.ID] = struct{}{}
		}
		*expectedRoles = m

		return nil
	}
}

func newProviderWithError() (tfprotov6.ProviderServer, error) {
	providerInstance := provider.New("test", adapter.NewMorpheus())()

	return providerserver.NewProtocol6WithError(providerInstance)()
}

var testAccProtoV6ProviderFactories = map[string]func() (
	tfprotov6.ProviderServer, error,
){
	"hpe": newProviderWithError,
}

func TestAccMorpheusUserResourceExample(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	// nolint: goconst
	providerConfig := testhelpers.ProviderBlock()

	dependenciesConfig := testhelpers.WhoamiBlock()

	name := acctest.RandomWithPrefix(t.Name())

	// MORPH-16460: the user fixtures previously hardcoded master-tenant role ids
	// (1/3), which fail as a sub-tenant. Create the role per-run and reference
	// it so the user is self-contained.
	roleConfig, err := role.RenderRoleUserConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}
	dependenciesConfig += roleConfig

	resourceConfig, err := user.RenderUserConfig(t, map[string]string{
		"Username": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.example",
			"username",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.example",
			"email",
			"foo@testacc.com",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.example",
			"role_ids.#",
			"1",
		),
		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_user.example",
			"role_ids.0",
			"hpe_morpheus_role.example",
			"id",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.example",
			"linux_key_pair_id",
			"100",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.example",
			"first_name",
			"Joe",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.example",
			"last_name",
			"User",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.example",
			"windows_username",
			"winuser",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.example",
			"receive_notifications",
			"false",
		),
		resource.TestCheckNoResourceAttr(
			"hpe_morpheus_user.example",
			"password_wo",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.example",
			"password_wo_version",
			"1",
		),
		resource.TestCheckNoResourceAttr(
			"hpe_morpheus_user.example",
			"password_wo",
		),
		resource.TestCheckNoResourceAttr(
			"hpe_morpheus_user.example",
			"windows_password_wo",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.example",
			"windows_password_wo_version",
			"1",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.example",
			"linux_username",
			"linuser",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.example",
			"windows_username",
			"winuser",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.example",
			"linux_password_wo_version",
			"1",
		),
		resource.TestCheckNoResourceAttr(
			"hpe_morpheus_user.example",
			"linux_password_wo",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:   providerConfig + dependenciesConfig + resourceConfig,
				Check:    checkFn,
				PlanOnly: false,
			},
		},
	})
}

// Test update of tenant_id attribute separately, as it
// requires delete/recreate.
// We may update this test once we can create a second tenant using
// the provider.
func TestAccMorpheusUserResourceUpdateTestIdOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	// nolint: goconst
	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	// MORPH-16460: create the role per-run instead of using master-tenant id 3.
	roleConfig, err := role.RenderRoleUserConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + testhelpers.WhoamiBlock() + roleConfig + `
resource "hpe_morpheus_user" "foo" {
	username = "` + name + `"
	email = "foo@testacc.com"
	password_wo = "Secret123!"
	role_ids = [hpe_morpheus_role.example.id]
	tenant_id = data.hpe_morpheus_whoami.current.tenant_id
}`,
				Check: resource.TestCheckResourceAttrPair(
					"hpe_morpheus_user.foo", "tenant_id",
					"data.hpe_morpheus_whoami.current", "tenant_id",
				),
			},
			{
				Config: providerConfig + testhelpers.WhoamiBlock() + roleConfig + `
resource "hpe_morpheus_user" "foo" {
	username = "` + name + `"
	email = "foo@testacc.com"
	password_wo = "Secret123!"
	role_ids = [hpe_morpheus_role.example.id]
	# changed
	tenant_id = 2
}`,
				ExpectNonEmptyPlan: true, // implicit delete/recreate
				PlanOnly:           true,
			},
		},
	})
}

// Check that we can create a user with only
// required attributes specified
func TestAccMorpheusUserResourceRequiredAttrsOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	// MORPH-16460: create the role per-run instead of using master-tenant id 3.
	roleConfig, err := role.RenderRoleUserConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	// nolint: goconst
	resourceConfig := roleConfig + `
resource "hpe_morpheus_user" "foo" {
	username = "` + name + `"
	email = "foo@testacc.com"
	password_wo = "Secret123!"
	role_ids = [hpe_morpheus_role.example.id]
}
`
	resourceConfigPostImport := roleConfig + `
resource "hpe_morpheus_user" "foo" {
	username = "` + name + `"
	email = "foo@testacc.com"
	# password_wo = "Secret123!"
	role_ids = [hpe_morpheus_role.example.id]
}
`
	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"username",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"email",
			"foo@testacc.com",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"role_ids.#",
			"1",
		),
		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_user.foo",
			"role_ids.0",
			"hpe_morpheus_role.example",
			"id",
		),
		resource.TestCheckNoResourceAttr(
			"hpe_morpheus_user.foo",
			"linux_username",
		),
		resource.TestCheckNoResourceAttr(
			"hpe_morpheus_user.foo",
			"linux_key_pair_id",
		),
		resource.TestCheckNoResourceAttr(
			"hpe_morpheus_user.foo",
			"windows_username",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"receive_notifications",
			"true",
		),
		resource.TestCheckNoResourceAttr(
			"hpe_morpheus_user.foo",
			"password_wo",
		),
		resource.TestCheckNoResourceAttr(
			"hpe_morpheus_user.foo",
			"password_wo_version",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:   providerConfig + resourceConfig,
				Check:    checkFn,
				PlanOnly: false,
			},
			{
				// Check that a post-apply plan detects no changes
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				Check:              checkFn,
				PlanOnly:           true,
			},
			{
				ImportState: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					// Read ID from the pre-import state
					rs := s.RootModule().
						Resources["hpe_morpheus_user.foo"]

					return rs.Primary.ID, nil
				},
				ImportStateVerify: true, // Check state post import
				// write_only_revision is the time of the last apply that sent a
				// write-only value; an import sets its own.
				ImportStateVerifyIgnore: []string{"write_only_revision"},
				ResourceName:            "hpe_morpheus_user.foo",
				Check:                   checkFn,
			},
			{
				// Check that a post-import plan detects no changes
				// if write-only fields are omitted
				Config:             providerConfig + resourceConfigPostImport,
				ExpectNonEmptyPlan: false,
				Check:              checkFn,
				PlanOnly:           true,
			},
		},
	})
}

func TestAccMorpheusUserResourceUpdateOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()
	expectedRoles := map[string]struct{}{}

	name := acctest.RandomWithPrefix(t.Name())

	// MORPH-16460: create two roles per-run instead of master-tenant ids 3 and
	// 1. hpe_morpheus_role.example stands in for the former id 3 and
	// hpe_morpheus_role.example2 for the former id 1.
	rolesConfig := roleResourceConfig("example", name+"-a") +
		roleResourceConfig("example2", name+"-b")

	baseChecks := []resource.TestCheckFunc{
		captureRoleIDs(
			&expectedRoles,
			"hpe_morpheus_role.example",
			"hpe_morpheus_role.example2",
		),
		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_user.foo", "tenant_id",
			"data.hpe_morpheus_whoami.current", "tenant_id",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"username",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"first_name",
			"foo",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"last_name",
			"bar",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"email",
			"foo@testacc.com",
		),
		resource.TestCheckNoResourceAttr(
			"hpe_morpheus_user.foo",
			"password_wo",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"password_wo_version",
			"1",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"linux_username",
			"linus",
		),
		resource.TestCheckNoResourceAttr(
			"hpe_morpheus_user.foo",
			"linux_password_wo",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"linux_password_wo_version",
			"1",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"linux_key_pair_id",
			"100",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"windows_username",
			"bill",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"windows_password_wo_version",
			"1",
		),
		resource.TestCheckNoResourceAttr(
			"hpe_morpheus_user.foo",
			"windows_password_wo",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"receive_notifications",
			"false",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"receive_notifications",
			"false",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"role_ids.#",
			"2",
		),
		checkRole(
			"hpe_morpheus_user.foo",
			"role_ids.0",
			&expectedRoles,
		),
		checkRole(
			"hpe_morpheus_user.foo",
			"role_ids.1",
			&expectedRoles,
		),
		checkStrayRoles(&expectedRoles),
	}

	passwordWoCheck := resource.TestCheckResourceAttr(
		"hpe_morpheus_user.foo",
		"password_wo_version",
		"1",
	)

	checkFn := resource.ComposeAggregateTestCheckFunc(
		append(baseChecks, passwordWoCheck)...,
	)

	expectedUpdateRoles := map[string]struct{}{}
	updateChecks := []resource.TestCheckFunc{
		captureRoleIDs(
			&expectedUpdateRoles,
			"hpe_morpheus_role.example2",
		),
		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_user.foo", "tenant_id",
			"data.hpe_morpheus_whoami.current", "tenant_id",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"username",
			name+"Updated",
		),
		resource.TestCheckNoResourceAttr(
			"hpe_morpheus_user.foo",
			"first_name",
		),
		resource.TestCheckNoResourceAttr(
			"hpe_morpheus_user.foo",
			"first_name",
		),
		resource.TestCheckNoResourceAttr(
			"hpe_morpheus_user.foo",
			"last_name",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"email",
			"bar@testacc.com",
		),
		resource.TestCheckNoResourceAttr(
			"hpe_morpheus_user.foo",
			"password_wo",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"password_wo_version",
			"2",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"linux_username",
			"torvalds",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"linux_password_wo_version",
			"2",
		),
		resource.TestCheckNoResourceAttr(
			"hpe_morpheus_user.foo",
			"linux_password_wo",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"linux_key_pair_id",
			"101",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"windows_username",
			"gates",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"windows_password_wo_version",
			"2",
		),
		resource.TestCheckNoResourceAttr(
			"hpe_morpheus_user.foo",
			"windows_password_wo",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"receive_notifications",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"role_ids.#",
			"1",
		),
		checkRole(
			"hpe_morpheus_user.foo",
			"role_ids.0",
			&expectedUpdateRoles,
		),
		checkStrayRoles(&expectedUpdateRoles),
	}

	checkUpdateFn := resource.ComposeAggregateTestCheckFunc(
		updateChecks...,
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + testhelpers.WhoamiBlock() + rolesConfig + `
resource "hpe_morpheus_user" "foo" {
	# Assumes tenant_id 1 pre-exists
	tenant_id = data.hpe_morpheus_whoami.current.tenant_id
	username = "` + name + `"
	email = "foo@testacc.com"
	password_wo = "Secret123!"
	password_wo_version = 1
	role_ids = [hpe_morpheus_role.example.id, hpe_morpheus_role.example2.id]
	first_name = "foo"
	last_name = "bar"
	linux_username = "linus"
	linux_password_wo = "Linux123!"
	linux_password_wo_version = 1
	linux_key_pair_id = 100
	receive_notifications = false
	windows_username = "bill"
	windows_password_wo = "Windows123!"
	windows_password_wo_version = 1
}`,
				Check:    checkFn,
				PlanOnly: false,
			},
			{
				Config: providerConfig + testhelpers.WhoamiBlock() + rolesConfig + `
# checks plan has no effect
resource "hpe_morpheus_user" "foo" {
	tenant_id = data.hpe_morpheus_whoami.current.tenant_id
	username = "` + name + `"
	email = "foo@testacc.com"
	password_wo = "Secret123!"
	password_wo_version = 1
	role_ids = [hpe_morpheus_role.example.id, hpe_morpheus_role.example2.id]
	first_name = "foo"
	last_name = "bar"
	linux_username = "linus"
	linux_password_wo = "Linux123!"
	linux_password_wo_version = 1
	linux_key_pair_id = 100
	receive_notifications = false
	windows_username = "bill"
	windows_password_wo = "Windows123!"
	windows_password_wo_version = 1
}`,
				Check:              checkFn,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
			{
				Config: providerConfig + testhelpers.WhoamiBlock() + rolesConfig + `
# checks plan detects first_name change to null
resource "hpe_morpheus_user" "foo" {
	tenant_id = data.hpe_morpheus_whoami.current.tenant_id
	username = "` + name + `"
	email = "foo@testacc.com"
	password_wo = "Secret123!"
	password_wo_version = 1
	role_ids = [hpe_morpheus_role.example.id, hpe_morpheus_role.example2.id]
	# changed
	# first_name = "foo"
	last_name = "bar"
	linux_username = "linus"
	linux_password_wo = "Linux123!"
	linux_password_wo_version = 1
	linux_key_pair_id = 100
	receive_notifications = false
	windows_username = "bill"
	windows_password_wo = "Windows123!"
	windows_password_wo_version = 1
}`,
				ExpectNonEmptyPlan: true,
				PlanOnly:           true,
			},
			{
				Config: providerConfig + testhelpers.WhoamiBlock() + rolesConfig + `
# checks plan detects first_name change
resource "hpe_morpheus_user" "foo" {
	tenant_id = data.hpe_morpheus_whoami.current.tenant_id
	username = "` + name + `"
	email = "foo@testacc.com"
	password_wo = "Secret123!"
	password_wo_version = 1
	role_ids = [hpe_morpheus_role.example.id, hpe_morpheus_role.example2.id]
	# changed
	first_name = "newfoo"
	last_name = "bar"
	linux_username = "linus"
	linux_password_wo = "Linux123!"
	linux_password_wo_version = 1
	linux_key_pair_id = 100
	receive_notifications = false
	windows_username = "bill"
	windows_password_wo = "Windows123!"
	windows_password_wo_version = 1
}`,
				ExpectNonEmptyPlan: true,
				PlanOnly:           true,
			},
			{
				Config: providerConfig + testhelpers.WhoamiBlock() + rolesConfig + `
# checks plan detects last_name change to null
resource "hpe_morpheus_user" "foo" {
	tenant_id = data.hpe_morpheus_whoami.current.tenant_id
	username = "` + name + `"
	email = "foo@testacc.com"
	password_wo = "Secret123!"
	password_wo_version = 1
	role_ids = [hpe_morpheus_role.example.id, hpe_morpheus_role.example2.id]
	first_name = "foo"
	# changed
	# last_name = "bar"
	linux_username = "linus"
	linux_password_wo = "Linux123!"
	linux_password_wo_version = 1
	linux_key_pair_id = 100
	receive_notifications = false
	windows_username = "bill"
	windows_password_wo = "Windows123!"
	windows_password_wo_version = 1
}`,
				ExpectNonEmptyPlan: true,
				PlanOnly:           true,
			},
			{
				Config: providerConfig + testhelpers.WhoamiBlock() + rolesConfig + `
# checks plan detects last_name change
resource "hpe_morpheus_user" "foo" {
	tenant_id = data.hpe_morpheus_whoami.current.tenant_id
	username = "` + name + `"
	email = "foo@testacc.com"
	password_wo = "Secret123!"
	password_wo_version = 1
	role_ids = [hpe_morpheus_role.example.id, hpe_morpheus_role.example2.id]
	first_name = "foo"
	# changed
	last_name = "newbar"
	linux_username = "linus"
	linux_password_wo = "Linux123!"
	linux_password_wo_version = 1
	linux_key_pair_id = 100
	receive_notifications = false
	windows_username = "bill"
	windows_password_wo = "Windows123!"
	windows_password_wo_version = 1
}`,
				ExpectNonEmptyPlan: true,
				PlanOnly:           true,
			},
			{
				Config: providerConfig + testhelpers.WhoamiBlock() + rolesConfig + `
# checks plan detects password_wo_version to null
resource "hpe_morpheus_user" "foo" {
	tenant_id = data.hpe_morpheus_whoami.current.tenant_id
	username = "` + name + `"
	email = "foo@testacc.com"
	password_wo = "Secret123!"
	# changed
	# password_wo_version = 1
	role_ids = [hpe_morpheus_role.example.id, hpe_morpheus_role.example2.id]
	first_name = "foo"
	last_name = "bar"
	linux_username = "linus"
	linux_password_wo = "Linux123!"
	linux_password_wo_version = 1
	linux_key_pair_id = 100
	receive_notifications = false
	windows_username = "bill"
	windows_password_wo = "Windows123!"
	windows_password_wo_version = 1
}`,
				ExpectNonEmptyPlan: true,
				PlanOnly:           true,
			},
			{
				Config: providerConfig + testhelpers.WhoamiBlock() + rolesConfig + `
# checks plan detects no change if only password_wo is removed: an unset
# write-only attribute is not managed, so it is not a change
resource "hpe_morpheus_user" "foo" {
	tenant_id = data.hpe_morpheus_whoami.current.tenant_id
	username = "` + name + `"
	email = "foo@testacc.com"
	# changed
	# password_wo = "Secret123!"
	password_wo_version = 1
	role_ids = [hpe_morpheus_role.example.id, hpe_morpheus_role.example2.id]
	first_name = "foo"
	last_name = "bar"
	linux_username = "linus"
	linux_password_wo = "Linux123!"
	linux_password_wo_version = 1
	linux_key_pair_id = 100
	receive_notifications = false
	windows_username = "bill"
	windows_password_wo = "Windows123!"
	windows_password_wo_version = 1
}`,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
			{
				Config: providerConfig + testhelpers.WhoamiBlock() + rolesConfig + `
# checks plan detects changed role_ids
resource "hpe_morpheus_user" "foo" {
	tenant_id = data.hpe_morpheus_whoami.current.tenant_id
	username = "` + name + `"
	email = "foo@testacc.com"
	password_wo = "Secret123!"
	password_wo_version = 1
	# changed
	role_ids = [hpe_morpheus_role.example.id]
	first_name = "foo"
	last_name = "bar"
	linux_username = "linus"
	linux_password_wo = "Linux123!"
	linux_password_wo_version = 1
	linux_key_pair_id = 100
	receive_notifications = false
	windows_username = "bill"
	windows_password_wo = "Windows123!"
	windows_password_wo_version = 1
}`,
				ExpectNonEmptyPlan: true,
				PlanOnly:           true,
			},
			{
				Config: providerConfig + testhelpers.WhoamiBlock() + rolesConfig + `
# checks plan detects changed username
resource "hpe_morpheus_user" "foo" {
	tenant_id = data.hpe_morpheus_whoami.current.tenant_id
	# changed
	username = "` + name + `Updated"
	email = "foo@testacc.com"
	password_wo = "Secret123!"
	password_wo_version = 1
	role_ids = [hpe_morpheus_role.example.id, hpe_morpheus_role.example2.id]
	first_name = "foo"
	last_name = "bar"
	linux_username = "linus"
	linux_password_wo = "Linux123!"
	linux_password_wo_version = 1
	linux_key_pair_id = 100
	receive_notifications = false
	windows_username = "bill"
	windows_password_wo = "Windows123!"
	windows_password_wo_version = 1
}`,
				ExpectNonEmptyPlan: true,
				PlanOnly:           true,
			},
			{
				Config: providerConfig + testhelpers.WhoamiBlock() + rolesConfig + `
# checks plan detects changed windows username
resource "hpe_morpheus_user" "foo" {
	tenant_id = data.hpe_morpheus_whoami.current.tenant_id
	username = "` + name + `"
	email = "foo@testacc.com"
	password_wo = "Secret123!"
	password_wo_version = 1
	role_ids = [hpe_morpheus_role.example.id, hpe_morpheus_role.example2.id]
	first_name = "foo"
	last_name = "bar"
	linux_username = "linus"
	linux_password_wo = "Linux123!"
	linux_password_wo_version = 1
	linux_key_pair_id = 100
	receive_notifications = false
	# changed
	windows_username = "melinda"
	windows_password_wo = "Windows123!"
	windows_password_wo_version = 1
}`,
				ExpectNonEmptyPlan: true,
				PlanOnly:           true,
			},
			{
				Config: providerConfig + testhelpers.WhoamiBlock() + rolesConfig + `
# checks plan detects changed linux username
resource "hpe_morpheus_user" "foo" {
	tenant_id = data.hpe_morpheus_whoami.current.tenant_id
	username = "` + name + `"
	email = "foo@testacc.com"
	password_wo = "Secret123!"
	password_wo_version = 1
	role_ids = [hpe_morpheus_role.example.id, hpe_morpheus_role.example2.id]
	first_name = "foo"
	last_name = "bar"
	# changed
	linux_username = "bsd"
	linux_password_wo = "Linux123!"
	linux_password_wo_version = 1
	linux_key_pair_id = 100
	receive_notifications = false
	windows_username = "bill"
	windows_password_wo = "Windows123!"
	windows_password_wo_version = 1
}`,
				ExpectNonEmptyPlan: true,
				PlanOnly:           true,
			},
			{
				Config: providerConfig + testhelpers.WhoamiBlock() + rolesConfig + `
# checks plan detects changed linux password version
resource "hpe_morpheus_user" "foo" {
	tenant_id = data.hpe_morpheus_whoami.current.tenant_id
	username = "` + name + `"
	email = "foo@testacc.com"
	password_wo = "Secret123!"
	password_wo_version = 1
	role_ids = [hpe_morpheus_role.example.id, hpe_morpheus_role.example2.id]
	first_name = "foo"
	last_name = "bar"
	linux_username = "linus"
	linux_password_wo = "Linux123!"
	# changed
	linux_password_wo_version = 2
	linux_key_pair_id = 100
	receive_notifications = false
	windows_username = "bill"
	windows_password_wo = "Windows123!"
	windows_password_wo_version = 1
}`,
				ExpectNonEmptyPlan: true,
				PlanOnly:           true,
			},
			{
				Config: providerConfig + testhelpers.WhoamiBlock() + rolesConfig + `
# checks plan detects changed windows password version
resource "hpe_morpheus_user" "foo" {
	tenant_id = data.hpe_morpheus_whoami.current.tenant_id
	username = "` + name + `"
	email = "foo@testacc.com"
	password_wo = "Secret123!"
	password_wo_version = 1
	role_ids = [hpe_morpheus_role.example.id, hpe_morpheus_role.example2.id]
	first_name = "foo"
	last_name = "bar"
	linux_username = "linus"
	linux_password_wo = "Linux123!"
	linux_password_wo_version = 1
	linux_key_pair_id = 100
	receive_notifications = false
	windows_username = "bill"
	windows_password_wo = "Windows123!"
	# changed
	windows_password_wo_version = 2
}`,
				ExpectNonEmptyPlan: true,
				PlanOnly:           true,
			},
			{
				Config: providerConfig + testhelpers.WhoamiBlock() + rolesConfig + `
# checks plan detects changed linux key pair id
resource "hpe_morpheus_user" "foo" {
	tenant_id = data.hpe_morpheus_whoami.current.tenant_id
	username = "` + name + `"
	email = "foo@testacc.com"
	password_wo = "Secret123!"
	password_wo_version = 1
	role_ids = [hpe_morpheus_role.example.id, hpe_morpheus_role.example2.id]
	first_name = "foo"
	last_name = "bar"
	linux_username = "linus"
	linux_password_wo = "Linux123!"
	linux_password_wo_version = 1
	# changed
	linux_key_pair_id = 101
	receive_notifications = false
	windows_username = "bill"
	windows_password_wo = "Windows123!"
	windows_password_wo_version = 1
}`,
				ExpectNonEmptyPlan: true,
				PlanOnly:           true,
			},
			{
				Config: providerConfig + testhelpers.WhoamiBlock() + rolesConfig + `
# checks plan detects changed receive_notifications
resource "hpe_morpheus_user" "foo" {
	tenant_id = data.hpe_morpheus_whoami.current.tenant_id
	username = "` + name + `"
	email = "foo@testacc.com"
	password_wo = "Secret123!"
	password_wo_version = 1
	role_ids = [hpe_morpheus_role.example.id, hpe_morpheus_role.example2.id]
	first_name = "foo"
	last_name = "bar"
	linux_username = "linus"
	linux_password_wo = "Linux123!"
	linux_password_wo_version = 1
	linux_key_pair_id = 100
	# changed
	receive_notifications = true
	windows_username = "bill"
	windows_password_wo = "Windows123!"
	windows_password_wo_version = 1
}`,
				ExpectNonEmptyPlan: true,
				PlanOnly:           true,
			},
			{
				Config: providerConfig + testhelpers.WhoamiBlock() + rolesConfig + `
# checks apply of changes to all changeable fields
resource "hpe_morpheus_user" "foo" {
	tenant_id = data.hpe_morpheus_whoami.current.tenant_id
	# changed
	username = "` + name + `Updated"
	# changed
	email = "bar@testacc.com"
	# changed
	password_wo = "Secret456!"
	# changed
	password_wo_version = 2
	# changed
	role_ids = [hpe_morpheus_role.example2.id]
	# changed
	# first_name = ""
	# changed - explicitly null
	last_name = null
	# changed
	linux_username = "torvalds"
	# changed
	linux_password_wo = "Linux1.0!"
	# changed
	linux_password_wo_version = 2
	# changed
	linux_key_pair_id = 101
	# changed
	receive_notifications = true
	# changed
	windows_username = "gates"
	# changed
	windows_password_wo = "Windows95!"
	# changed
	windows_password_wo_version = 2
}`,
				Check:    checkUpdateFn,
				PlanOnly: false,
			},
			{
				Config: providerConfig + testhelpers.WhoamiBlock() + rolesConfig + `
# checks plan has no effect
resource "hpe_morpheus_user" "foo" {
	tenant_id = data.hpe_morpheus_whoami.current.tenant_id
	# changed
	username = "` + name + `Updated"
	# changed
	email = "bar@testacc.com"
	# changed
	password_wo = "Secret456!"
	# changed
	password_wo_version = 2
	# changed
	role_ids = [hpe_morpheus_role.example2.id]
	# changed
	# first_name = ""
	# changed - explicitly null
	last_name = null
	# changed
	linux_username = "torvalds"
	# changed
	linux_password_wo = "Linux1.0!"
	# changed
	linux_password_wo_version = 2
	# changed
	linux_key_pair_id = 101
	# changed
	receive_notifications = true
	# changed
	windows_username = "gates"
	# changed
	windows_password_wo = "Windows95!"
	# changed
	windows_password_wo_version = 2
}`,
				Check:              checkUpdateFn,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func TestAccMorpheusUserResourceUpdateNoTenantIdOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()
	expectedRoles := map[string]struct{}{}

	name := acctest.RandomWithPrefix(t.Name())

	// MORPH-16460: create two roles per-run instead of master-tenant ids 3 and 1.
	rolesConfig := roleResourceConfig("example", name+"-a") +
		roleResourceConfig("example2", name+"-b")

	baseChecks := []resource.TestCheckFunc{
		captureRoleIDs(
			&expectedRoles,
			"hpe_morpheus_role.example",
			"hpe_morpheus_role.example2",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"username",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"first_name",
			"foo",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"last_name",
			"bar",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"email",
			"foo@testacc.com",
		),
		resource.TestCheckNoResourceAttr(
			"hpe_morpheus_user.foo",
			"password_wo",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"password_wo_version",
			"1",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"linux_username",
			"linus",
		),
		resource.TestCheckNoResourceAttr(
			"hpe_morpheus_user.foo",
			"linux_password_wo",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"linux_password_wo_version",
			"1",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"linux_key_pair_id",
			"100",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"windows_username",
			"bill",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"windows_password_wo_version",
			"1",
		),
		resource.TestCheckNoResourceAttr(
			"hpe_morpheus_user.foo",
			"windows_password_wo",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"receive_notifications",
			"false",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"receive_notifications",
			"false",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"role_ids.#",
			"2",
		),
		checkRole(
			"hpe_morpheus_user.foo",
			"role_ids.0",
			&expectedRoles,
		),
		checkRole(
			"hpe_morpheus_user.foo",
			"role_ids.1",
			&expectedRoles,
		),
		checkStrayRoles(&expectedRoles),
	}

	passwordWoCheck := resource.TestCheckResourceAttr(
		"hpe_morpheus_user.foo",
		"password_wo_version",
		"1",
	)

	checkFn := resource.ComposeAggregateTestCheckFunc(
		append(baseChecks, passwordWoCheck)...,
	)

	expectedUpdateRoles := map[string]struct{}{}
	updateChecks := []resource.TestCheckFunc{
		captureRoleIDs(
			&expectedUpdateRoles,
			"hpe_morpheus_role.example",
			"hpe_morpheus_role.example2",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"username",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"first_name",
			"foo",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"last_name",
			"bar2",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"email",
			"foo@testacc.com",
		),
		resource.TestCheckNoResourceAttr(
			"hpe_morpheus_user.foo",
			"password_wo",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"password_wo_version",
			"1",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"linux_username",
			"linus",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"linux_password_wo_version",
			"1",
		),
		resource.TestCheckNoResourceAttr(
			"hpe_morpheus_user.foo",
			"linux_password_wo",
		),
		resource.TestCheckNoResourceAttr(
			"hpe_morpheus_user.foo",
			"linux_key_pair_id",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"windows_username",
			"bill",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"windows_password_wo_version",
			"1",
		),
		resource.TestCheckNoResourceAttr(
			"hpe_morpheus_user.foo",
			"windows_password_wo",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"receive_notifications",
			"false",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"role_ids.#",
			"2",
		),
		checkRole(
			"hpe_morpheus_user.foo",
			"role_ids.0",
			&expectedUpdateRoles,
		),
		checkRole(
			"hpe_morpheus_user.foo",
			"role_ids.1",
			&expectedUpdateRoles,
		),
		checkStrayRoles(&expectedUpdateRoles),
	}

	checkUpdateFn := resource.ComposeAggregateTestCheckFunc(
		updateChecks...,
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + testhelpers.WhoamiBlock() + rolesConfig + `
resource "hpe_morpheus_user" "foo" {
	# create
	username = "` + name + `"
	email = "foo@testacc.com"
	password_wo = "Secret123!"
	password_wo_version = 1
	role_ids = [hpe_morpheus_role.example.id, hpe_morpheus_role.example2.id]
	first_name = "foo"
	last_name = "bar"
	linux_username = "linus"
	linux_password_wo = "Linux123!"
	linux_password_wo_version = 1
	linux_key_pair_id = 100
	receive_notifications = false
	windows_username = "bill"
	windows_password_wo = "Windows123!"
	windows_password_wo_version = 1
}`,
				Check:    checkFn,
				PlanOnly: false,
			},
			{
				Config: providerConfig + testhelpers.WhoamiBlock() + rolesConfig + `
resource "hpe_morpheus_user" "foo" {
	# post create plan
	username = "` + name + `"
	email = "foo@testacc.com"
	password_wo = "Secret123!"
	password_wo_version = 1
	role_ids = [hpe_morpheus_role.example.id, hpe_morpheus_role.example2.id]
	first_name = "foo"
	last_name = "bar"
	linux_username = "linus"
	linux_password_wo = "Linux123!"
	linux_password_wo_version = 1
	linux_key_pair_id = 100
	receive_notifications = false
	windows_username = "bill"
	windows_password_wo = "Windows123!"
	windows_password_wo_version = 1
}`,
				Check:              checkFn,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			{
				Config: providerConfig + testhelpers.WhoamiBlock() + rolesConfig + `
# checks plan has no effect
resource "hpe_morpheus_user" "foo" {
	# update 1
	username = "` + name + `"
	email = "foo@testacc.com"
	password_wo = "Secret123!"
	password_wo_version = 1
	role_ids = [hpe_morpheus_role.example.id, hpe_morpheus_role.example2.id]
	first_name = "foo"
	# changed
	last_name = "bar2"
	linux_username = "linus"
	linux_password_wo = "Linux123!"
	linux_password_wo_version = 1
	#linux_key_pair_id = 100
	receive_notifications = false
	windows_username = "bill"
	windows_password_wo = "Windows123!"
	windows_password_wo_version = 1
}`,
				Check: checkUpdateFn,
			},
			{
				Config: providerConfig + testhelpers.WhoamiBlock() + rolesConfig + `
# checks plan detects first_name change to null
resource "hpe_morpheus_user" "foo" {
	# post update 1 plan
	username = "` + name + `"
	email = "foo@testacc.com"
	password_wo = "Secret123!"
	password_wo_version = 1
	role_ids = [hpe_morpheus_role.example.id, hpe_morpheus_role.example2.id]
	first_name = "foo"
	last_name = "bar2"
	linux_username = "linus"
	linux_password_wo = "Linux123!"
	linux_password_wo_version = 1
	#linux_key_pair_id = 100
	receive_notifications = false
	windows_username = "bill"
	windows_password_wo = "Windows123!"
	windows_password_wo_version = 1
}`,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
			{
				Config: providerConfig + testhelpers.WhoamiBlock() + rolesConfig + `
resource "hpe_morpheus_user" "foo" {
	# update 2
	username = "` + name + `"
	email = "foo@testacc.com"
	password_wo = "Secret123!"
	password_wo_version = 1
	role_ids = [hpe_morpheus_role.example.id, hpe_morpheus_role.example2.id]
	first_name = "foo"
	last_name = "bar"
	linux_username = "linus"
	linux_password_wo = "Linux123!"
	linux_password_wo_version = 1
	linux_key_pair_id = 100
	receive_notifications = false
	windows_username = "bill"
	windows_password_wo = "Windows123!"
	windows_password_wo_version = 1
}`,
				Check:    checkFn,
				PlanOnly: false,
			},
			{
				Config: providerConfig + testhelpers.WhoamiBlock() + rolesConfig + `
resource "hpe_morpheus_user" "foo" {
	# post update 2 plan
	username = "` + name + `"
	email = "foo@testacc.com"
	password_wo = "Secret123!"
	password_wo_version = 1
	role_ids = [hpe_morpheus_role.example.id, hpe_morpheus_role.example2.id]
	first_name = "foo"
	last_name = "bar"
	linux_username = "linus"
	linux_password_wo = "Linux123!"
	linux_password_wo_version = 1
	linux_key_pair_id = 100
	receive_notifications = false
	windows_username = "bill"
	windows_password_wo = "Windows123!"
	windows_password_wo_version = 1
}`,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func TestAccMorpheusUserResourceAllAttrsOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	// MORPH-16460: create two roles per-run instead of master-tenant ids 3 and 1.
	rolesConfig := roleResourceConfig("example", name+"-a") +
		roleResourceConfig("example2", name+"-b")

	resourceCfg := rolesConfig + `
# Role id 0 causes a test failure because it is ignored by
# the server and only the other two roles are created
#resource "hpe_morpheus_user" "bar" {
#username = "test101"
#email = "foo@testacc.com"
#password = "Secret123!"
#roles = [3,0,1]
#}
resource "hpe_morpheus_user" "foo" {
	# Assumes tenant_id 1 pre-exists
	tenant_id = data.hpe_morpheus_whoami.current.tenant_id
	username = "` + name + `"
	email = "foo@testacc.com"
	password_wo = "Secret123!"
	password_wo_version = 1
	role_ids = [hpe_morpheus_role.example.id, hpe_morpheus_role.example2.id]
	first_name = "foo"
	last_name = "bar"
	linux_username = "linus"
	linux_password_wo = "Linux123!"
	linux_password_wo_version = 1
	linux_key_pair_id = 100
	receive_notifications = false
	windows_username = "bill"
	windows_password_wo = "Windows123!"
	windows_password_wo_version = 1
}
`
	expectedRoles := map[string]struct{}{}

	baseChecks := []resource.TestCheckFunc{
		captureRoleIDs(
			&expectedRoles,
			"hpe_morpheus_role.example",
			"hpe_morpheus_role.example2",
		),
		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_user.foo", "tenant_id",
			"data.hpe_morpheus_whoami.current", "tenant_id",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"username",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"email",
			"foo@testacc.com",
		),
		resource.TestCheckNoResourceAttr(
			"hpe_morpheus_user.foo",
			"password_wo",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"linux_username",
			"linus",
		),
		resource.TestCheckNoResourceAttr(
			"hpe_morpheus_user.foo",
			"linux_password_wo",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"linux_key_pair_id",
			"100",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"windows_username",
			"bill",
		),
		resource.TestCheckNoResourceAttr(
			"hpe_morpheus_user.foo",
			"windows_password_wo",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"windows_password_wo_version",
			"1",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"receive_notifications",
			"false",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"receive_notifications",
			"false",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"role_ids.#",
			"2",
		),
		checkRole(
			"hpe_morpheus_user.foo",
			"role_ids.0",
			&expectedRoles,
		),
		checkRole(
			"hpe_morpheus_user.foo",
			"role_ids.1",
			&expectedRoles,
		),
		checkStrayRoles(&expectedRoles),
	}

	passwordWoCheck := resource.TestCheckResourceAttr(
		"hpe_morpheus_user.foo",
		"password_wo_version",
		"1",
	)
	linuxPasswordWoVersionCheck := resource.TestCheckResourceAttr(
		"hpe_morpheus_user.foo",
		"linux_password_wo_version",
		"1",
	)
	windowsPasswordWoVersionCheck := resource.TestCheckResourceAttr(
		"hpe_morpheus_user.foo",
		"windows_password_wo_version",
		"1",
	)

	checkFn := resource.ComposeAggregateTestCheckFunc(
		append(
			baseChecks,
			passwordWoCheck,
			linuxPasswordWoVersionCheck,
			windowsPasswordWoVersionCheck,
		)...,
	)

	linuxPasswordWoVersionImportCheck := resource.TestCheckNoResourceAttr(
		"hpe_morpheus_user.foo",
		"linux_password_wo_version",
	)
	windowsPasswordWoVersionImportCheck := resource.TestCheckNoResourceAttr(
		"hpe_morpheus_user.foo",
		"windows_password_wo_version",
	)
	passwordWoImportCheck := resource.TestCheckNoResourceAttr(
		"hpe_morpheus_user.foo",
		"password_wo_version",
	)

	checkImportFn := resource.ComposeAggregateTestCheckFunc(
		append(
			baseChecks,
			passwordWoImportCheck,
			linuxPasswordWoVersionImportCheck,
			windowsPasswordWoVersionImportCheck,
		)...,
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:   providerConfig + testhelpers.WhoamiBlock() + resourceCfg,
				Check:    checkFn,
				PlanOnly: false,
			},
			{
				// state from import test exists in memory (not written to disk)
				ImportState: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					// Read ID from the pre-import state
					rs := s.RootModule().
						Resources["hpe_morpheus_user.foo"]

					return rs.Primary.ID, nil
				},
				ImportStateVerify: true, // Check state post import (in memory)
				ImportStateVerifyIgnore: []string{
					"password_wo_version",
					"linux_password_wo_version",
					"windows_password_wo_version",
					"write_only_revision",
				},
				ResourceName: "hpe_morpheus_user.foo",
				Check:        checkImportFn,
			},
		},
	})
}

func TestAccMorpheusUserResourceMissingRoles(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	name := acctest.RandomWithPrefix(t.Name())

	providerConfig := `
provider "hpe" {
	morpheus {
		url = ""
		username = ""
		password = ""
	}
}

resource "hpe_morpheus_user" "foo" {
	username = "` + name + `"
	email = "bar@testacc.com"
	password = "Secret123!"
	# role_ids = [3,1]
}
`
	expected := `The argument "role_ids" is required`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:             providerConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
				ExpectError:        regexp.MustCompile(expected),
			},
		},
	})
}

func TestAccMorpheusUserResourceMissingUsername(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	name := acctest.RandomWithPrefix(t.Name())

	providerConfig := `
provider "hpe" {
	morpheus {
		url = ""
		username = ""
		password = ""
	}
}

resource "hpe_morpheus_user" "foo" {
	#username = "` + name + `"
	email = "bar@testacc.com"
	password = "Secret123!"
	role_ids = [3,1]
}
`
	expected := `The argument "username" is required`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:             providerConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
				ExpectError:        regexp.MustCompile(expected),
			},
		},
	})
}

func TestAccMorpheusUserResourceMissingEmail(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	name := acctest.RandomWithPrefix(t.Name())

	providerConfig := `
provider "hpe" {
	morpheus {
		url = ""
		username = ""
		password = ""
	}
}

resource "hpe_morpheus_user" "foo" {
	username = "` + name + `"
	#email = "bar@testacc.com"
	password = "Secret123!"
	role_ids = [3,1]
}
`
	expected := `The argument "email" is required`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:             providerConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
				ExpectError:        regexp.MustCompile(expected),
			},
		},
	})
}

// password_wo is required for create (but not import) here we check that it is
// correctly identified as missing during plan (i.e. before Create is called)
func TestAccMorpheusUserResourceMissingPasswordWo(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	name := acctest.RandomWithPrefix(t.Name())

	// A real provider block is required here. The "required on create" check for
	// password_wo is an attribute plan modifier
	// (modifiers.RequireOnCreateModifier, see the password_wo attribute in
	// schema_gen.go), so it only runs during PlanResourceChange, which Terraform
	// performs after ConfigureProvider. With an unconfigured provider the
	// provider fails first with "Missing Morpheus connection details" and the
	// plan modifier is never reached.
	//
	// This differs from the sibling MissingEmail/MissingRoles/MissingUsername
	// tests, which omit genuinely Required attributes and are therefore rejected
	// by schema validation before the provider is ever configured.
	providerConfig := testhelpers.ProviderBlock() + `
resource "hpe_morpheus_user" "foo" {
	username = "` + name + `"
	email = "bar@testacc.com"
	#password_wo = "Secret123!"
	role_ids = [3,1]
}
`
	expected := `'password_wo' not set`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:             providerConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
				ExpectError:        regexp.MustCompile(expected),
			},
		},
	})
}

// Here we use a two phase approach to import that
// allows creating a resource using terraform
// while preserving the import state for follow
// on tests.
//
// The testing here is similar to other import
// related tests in this file, but here we
// are able to run plan after import, having
// inherited the import state.
func TestAccMorpheusUserResourceImportOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	// MORPH-16460: create two roles per-run instead of master-tenant ids 3 and
	// 1. The roles must persist across both TestCases (create then import), so
	// TestCase 1 removes them from state without deleting them and TestCase 2
	// re-imports them alongside the user.
	rolesConfig := roleResourceConfig("example", name+"-a") +
		roleResourceConfig("example2", name+"-b")

	// nolint: gosec
	resourceCfgWithPassword := rolesConfig + `
resource "hpe_morpheus_user" "foo" {
	# Assumes tenant_id 1 pre-exists
	tenant_id = data.hpe_morpheus_whoami.current.tenant_id
	username = "` + name + `"
	email = "foo@testacc.com"
	password_wo = "Secret123!"
	password_wo_version = 1
	role_ids = [hpe_morpheus_role.example.id, hpe_morpheus_role.example2.id]
	first_name = "foo"
	last_name = "bar"
	linux_username = "linus"
	linux_key_pair_id = 100
	receive_notifications = false
	windows_username = "bill"
}
`
	// nolint: gosec
	resourceCfgNoPassword := rolesConfig + `
resource "hpe_morpheus_user" "foo" {
	# Assumes tenant_id 1 pre-exists
	tenant_id = data.hpe_morpheus_whoami.current.tenant_id
	username = "` + name + `"
	email = "foo@testacc.com"
        #password_wo = "Secret123!"
        #password_wo_version = 1
	role_ids = [hpe_morpheus_role.example.id, hpe_morpheus_role.example2.id]
	first_name = "foo"
	last_name = "bar"
	linux_username = "linus"
	linux_key_pair_id = 100
	receive_notifications = false
	windows_username = "bill"
}
`
	resourceCfgRemove := `
# This allows us to create a resource using the provider
# and then import it in a separate resource.Test.
#
# A regular 'Config:' test step (with import block) can be used for the
# import test (rather than an 'ImportState:' style test step)
#
# The state is preserved after import and available to
# subsequent tests.
#
# The 'removed' block means the resource is removed from state
# (and terraform control) but not deleted.
#
# This avoids both triggering two deletes for the same resource,
# and the dreaded "resource is already under terraform control"
# error when running the follow on import test.
#
# MORPH-16460: the per-run roles are removed from state (destroy = false)
# as well, so they survive for the import TestCase below, which re-imports
# them alongside the user.
removed {
from = hpe_morpheus_user.foo

lifecycle {
destroy = false
}
}

removed {
from = hpe_morpheus_role.example

lifecycle {
destroy = false
}
}

removed {
from = hpe_morpheus_role.example2

lifecycle {
destroy = false
}
}
`
	baseChecks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_user.foo", "tenant_id",
			"data.hpe_morpheus_whoami.current", "tenant_id",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"username",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"email",
			"foo@testacc.com",
		),
		resource.TestCheckNoResourceAttr(
			"hpe_morpheus_user.foo",
			"password_wo",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"linux_username",
			"linus",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"linux_key_pair_id",
			"100",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"windows_username",
			"bill",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"receive_notifications",
			"false",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"receive_notifications",
			"false",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"role_ids.#",
			"2",
		),
	}

	expectedCreateRoles := map[string]struct{}{}
	createChecks := []resource.TestCheckFunc{
		captureRoleIDs(
			&expectedCreateRoles,
			"hpe_morpheus_role.example",
			"hpe_morpheus_role.example2",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"password_wo_version",
			"1",
		),
		resource.TestCheckNoResourceAttr(
			"hpe_morpheus_user.foo",
			"password_wo",
		),
		checkRole("hpe_morpheus_user.foo", "role_ids.0", &expectedCreateRoles),
		checkRole("hpe_morpheus_user.foo", "role_ids.1", &expectedCreateRoles),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(
		append(baseChecks, createChecks...)...,
	)

	var cachedID string
	var cachedRoleID0, cachedRoleID1 string

	// This is a new TestCase - we know for sure
	// we inherit no state from the TestCase above
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + testhelpers.WhoamiBlock() + resourceCfgWithPassword,
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						// Cache ID for use later
						rs := s.RootModule().Resources["hpe_morpheus_user.foo"]
						if rs == nil {
							return fmt.Errorf("resource not found")
						}
						cachedID = rs.Primary.ID

						// MORPH-16460: cache the per-run role ids so the import
						// TestCase can re-import the same roles.
						r0 := s.RootModule().Resources["hpe_morpheus_role.example"]
						r1 := s.RootModule().Resources["hpe_morpheus_role.example2"]
						if r0 == nil || r1 == nil {
							return fmt.Errorf("role resource not found")
						}
						cachedRoleID0 = r0.Primary.ID
						cachedRoleID1 = r1.Primary.ID

						return nil
					},
					checkFn,
				),
				PlanOnly: false,
			},
			{
				// remove resource from terraform state (without deleting it)
				Config:   providerConfig + resourceCfgRemove,
				PlanOnly: false,
			},
		},
	})

	importCfg := providerConfig + testhelpers.WhoamiBlock() + resourceCfgNoPassword + `
	import {
	  to = hpe_morpheus_user.foo
	  id = ` + cachedID + `
	}
	import {
	  to = hpe_morpheus_role.example
	  id = ` + cachedRoleID0 + `
	}
	import {
	  to = hpe_morpheus_role.example2
	  id = ` + cachedRoleID1 + `
	}
	`
	expectedImportRoles := map[string]struct{}{}
	importChecks := []resource.TestCheckFunc{
		captureRoleIDs(
			&expectedImportRoles,
			"hpe_morpheus_role.example",
			"hpe_morpheus_role.example2",
		),
		resource.TestCheckNoResourceAttr(
			"hpe_morpheus_user.foo",
			"password_wo_version",
		),
		resource.TestCheckNoResourceAttr(
			"hpe_morpheus_user.foo",
			"password_wo",
		),
		checkRole("hpe_morpheus_user.foo", "role_ids.0", &expectedImportRoles),
		checkRole("hpe_morpheus_user.foo", "role_ids.1", &expectedImportRoles),
	}

	checkImportFn := resource.ComposeAggregateTestCheckFunc(
		append(baseChecks, importChecks...)...,
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:   importCfg,
				PlanOnly: false,
				Check: resource.ComposeTestCheckFunc(
					checkImportFn,
				),
			},
			{
				// check that a plan after import detects no changes
				Config:             providerConfig + testhelpers.WhoamiBlock() + resourceCfgNoPassword,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

// Test that when tenant_id is not set (use state for unknown)
// updating last_name results in update operation, not delete/recreate
func TestAccMorpheusUserResourceUpdateLastNameWithoutTenantIdOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())
	var userID string

	// MORPH-16460: create the role per-run instead of using master-tenant id 3.
	rolesConfig := roleResourceConfig("example", name)

	resource.Test(
		t,
		resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + testhelpers.WhoamiBlock() + rolesConfig + `
resource "hpe_morpheus_user" "foo" {
	# Note: tenant_id is NOT set - uses "state for unknown"
	username = "` + name + `"
	email = "foo@testacc.com"
	password_wo = "Secret123!"
	role_ids = [hpe_morpheus_role.example.id]
	first_name = "John"
	last_name = "Doe"
	linux_key_pair_id = 100
}`,
					Check: resource.ComposeAggregateTestCheckFunc(
						// Capture the user ID for comparison in subsequent steps
						func(s *terraform.State) error {
							rs := s.RootModule().Resources["hpe_morpheus_user.foo"]
							if rs == nil {
								return fmt.Errorf("resource not found")
							}
							userID = rs.Primary.ID

							return nil
						},
						resource.TestCheckResourceAttr(
							"hpe_morpheus_user.foo",
							"username",
							name,
						),
						resource.TestCheckResourceAttr(
							"hpe_morpheus_user.foo",
							"email",
							"foo@testacc.com",
						),
						resource.TestCheckResourceAttr(
							"hpe_morpheus_user.foo",
							"first_name",
							"John",
						),
						resource.TestCheckResourceAttr(
							"hpe_morpheus_user.foo",
							"last_name",
							"Doe",
						),
						resource.TestCheckResourceAttr(
							"hpe_morpheus_user.foo",
							"role_ids.#",
							"1",
						),
						resource.TestCheckResourceAttrPair(
							"hpe_morpheus_user.foo",
							"role_ids.0",
							"hpe_morpheus_role.example",
							"id",
						),
						// tenant_id should be computed from state (not set in config)
						resource.TestCheckResourceAttrSet(
							"hpe_morpheus_user.foo",
							"tenant_id",
						),
					),
				},
				{
					Config: providerConfig + testhelpers.WhoamiBlock() + rolesConfig + `
resource "hpe_morpheus_user" "foo" {
	# Note: tenant_id is still NOT set - uses "state for unknown"
	username = "` + name + `"
	email = "foo@testacc.com"
	password_wo = "Secret123!"
	role_ids = [hpe_morpheus_role.example.id]
	first_name = "John"
	# Changed: update last_name
	last_name = "Smith"
	linux_key_pair_id = 100
}`,
					Check: resource.ComposeAggregateTestCheckFunc(
						// Verify the user ID hasn't changed (proving it's an update, not recreate)
						func(s *terraform.State) error {
							rs := s.RootModule().Resources["hpe_morpheus_user.foo"]
							if rs == nil {
								return fmt.Errorf("resource not found")
							}
							if rs.Primary.ID != userID {
								return fmt.Errorf(
									"user ID changed from %s to %s, indicating resource was recreated instead of updated",
									userID,
									rs.Primary.ID,
								)
							}

							return nil
						},
						resource.TestCheckResourceAttr(
							"hpe_morpheus_user.foo",
							"username",
							name,
						),
						resource.TestCheckResourceAttr(
							"hpe_morpheus_user.foo",
							"email",
							"foo@testacc.com",
						),
						resource.TestCheckResourceAttr(
							"hpe_morpheus_user.foo",
							"first_name",
							"John",
						),
						// Verify last_name was updated
						resource.TestCheckResourceAttr(
							"hpe_morpheus_user.foo",
							"last_name",
							"Smith",
						),
						resource.TestCheckResourceAttr(
							"hpe_morpheus_user.foo",
							"role_ids.#",
							"1",
						),
						resource.TestCheckResourceAttrPair(
							"hpe_morpheus_user.foo",
							"role_ids.0",
							"hpe_morpheus_role.example",
							"id",
						),
						// tenant_id should still be computed from state
						resource.TestCheckResourceAttrSet(
							"hpe_morpheus_user.foo",
							"tenant_id",
						),
					),
					// This is the key assertion: we expect this to be an update, not a recreate
					// If the resource gets deleted/recreated, the test will fail
					ExpectNonEmptyPlan: false, // After update, plan should be empty
				},
				{
					// Verify the plan shows no changes after the update
					Config: providerConfig + testhelpers.WhoamiBlock() + rolesConfig + `
resource "hpe_morpheus_user" "foo" {
	# Note: tenant_id is still NOT set - uses "state for unknown"
	username = "` + name + `"
	email = "foo@testacc.com"
	password_wo = "Secret123!"
	role_ids = [hpe_morpheus_role.example.id]
	first_name = "John"
	last_name = "Smith"
	linux_key_pair_id = 100
}`,
					PlanOnly:           true,
					ExpectNonEmptyPlan: false,
				},
			},
		},
	)
}

// TestAccMorpheusUserResourceRemoveOsUsernamesOk pins MORPH-18061: removing
// a configured linux_username / windows_username from the configuration
// clears them on the appliance and the plan then converges.
//
// Before the fix the update sent both usernames as explicit nulls, which the
// API rejects with HTTP 500 (MORPH-18060), reported as "resource was not
// updated". Only this path reaches a null: an attribute that was never set is
// planned as unknown, whose string pointer is "" rather than nil, so the API
// receives "" and the existing tests never saw the failure.
func TestAccMorpheusUserResourceRemoveOsUsernamesOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())
	rolesConfig := roleResourceConfig("example", name+"-a")

	config := func(osUsernames string) string {
		return providerConfig + testhelpers.WhoamiBlock() + rolesConfig + `
resource "hpe_morpheus_user" "example" {
	tenant_id           = data.hpe_morpheus_whoami.current.tenant_id
	username            = "` + name + `"
	email               = "` + name + `@testacc.com"
	password_wo         = "Secret123!"
	password_wo_version = 1
	role_ids            = [hpe_morpheus_role.example.id]
` + osUsernames + `
}`
	}

	withOsUsernames := `
	linux_username   = "linus"
	windows_username = "bill"
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: config(withOsUsernames),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_user.example", "linux_username", "linus"),
					resource.TestCheckResourceAttr("hpe_morpheus_user.example", "windows_username", "bill"),
				),
			},
			{
				// Remove both usernames: the update that used to fail with
				// "resource was not updated".
				Config: config(""),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("hpe_morpheus_user.example", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckNoResourceAttr("hpe_morpheus_user.example", "linux_username"),
					resource.TestCheckNoResourceAttr("hpe_morpheus_user.example", "windows_username"),
				),
			},
			{
				// The plan converges: nothing left to send.
				Config:             config(""),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}
