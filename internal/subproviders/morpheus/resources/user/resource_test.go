// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package user_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/HPE/terraform-provider-hpe/internal/provider"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/testhelpers"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestMain(m *testing.M) {
	code := m.Run()
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func checkRole(
	resourceName string,
	roleIDAttr string,
	expectedRoles map[string]struct{},
) func(*terraform.State) error {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		roleID := rs.Primary.Attributes[roleIDAttr]
		if _, ok := expectedRoles[roleID]; !ok {
			return fmt.Errorf("role ID %s not found ", roleID)
		}

		delete(expectedRoles, roleID)

		return nil
	}
}

func checkStrayRoles(
	expectedRoles map[string]struct{},
) func(*terraform.State) error {
	return func(_ *terraform.State) error {
		if len(expectedRoles) != 0 {
			return fmt.Errorf("not all role_ids found %s", expectedRoles)
		}

		return nil
	}
}

func newProviderWithError() (tfprotov6.ProviderServer, error) {
	providerInstance := provider.New("test", morpheus.New())()

	return providerserver.NewProtocol6WithError(providerInstance)()
}

var testAccProtoV6ProviderFactories = map[string]func() (
	tfprotov6.ProviderServer, error,
){
	"hpe": newProviderWithError,
}

func TestAccMorpheusUserExample(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	// nolint: goconst
	providerConfig := testhelpers.ProviderBlock()

	path := "../../../../../examples/resources/hpe_morpheus_user/resource.tf"
	exampleConfig, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Error reading example config: %v", err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.example",
			"username",
			"testacc-example",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.example",
			"email",
			"user@example.com",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.example",
			"role_ids.#",
			"1",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.example",
			"role_ids.0",
			"1",
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
				Config:   providerConfig + string(exampleConfig),
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
func TestAccMorpheusUserUpdateTestIdOk(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	// nolint: goconst
	providerConfig := testhelpers.ProviderBlock()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "hpe_morpheus_user" "foo" {
	username = "testacc-TestAccMorpheusUserUpdateTestIdOk"
	email = "foo@hpe.com"
	password_wo = "Secret123!"
	role_ids = [3]
	tenant_id = 1
}`,
				Check: resource.TestCheckResourceAttr(
					"hpe_morpheus_user.foo",
					"tenant_id",
					"1",
				),
			},
			{
				Config: providerConfig + `
resource "hpe_morpheus_user" "foo" {
	username = "testacc-TestAccMorpheusUserUpdateTestIdOk"
	email = "foo@hpe.com"
	password_wo = "Secret123!"
	role_ids = [3]
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
func TestAccMorpheusUserRequiredAttrsOk(t *testing.T) {
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	// nolint: goconst
	resourceConfig := `
resource "hpe_morpheus_user" "foo" {
	username = "testacc-TestAccMorpheusUserRequiredAttrsOk"
	email = "foo@hpe.com"
	password_wo = "Secret123!"
	role_ids = [3]
}
`
	resourceConfigPostImport := `
resource "hpe_morpheus_user" "foo" {
	username = "testacc-TestAccMorpheusUserRequiredAttrsOk"
	email = "foo@hpe.com"
	# password_wo = "Secret123!"
	role_ids = [3]
}
`
	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"username",
			"testacc-TestAccMorpheusUserRequiredAttrsOk",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"email",
			"foo@hpe.com",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"role_ids.#",
			"1",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"role_ids.0",
			"3",
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
				ResourceName:      "hpe_morpheus_user.foo",
				Check:             checkFn,
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

func TestAccMorpheusUserUpdateOk(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()
	expectedRoles := map[string]struct{}{"3": {}, "1": {}}

	baseChecks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"tenant_id",
			"1",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"username",
			"testacc-TestAccMorpheusUserUpdateOk",
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
			"foo@hpe.com",
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
			expectedRoles,
		),
		checkRole(
			"hpe_morpheus_user.foo",
			"role_ids.1",
			expectedRoles,
		),
		checkStrayRoles(expectedRoles),
	}

	passwordWoCheck := resource.TestCheckResourceAttr(
		"hpe_morpheus_user.foo",
		"password_wo_version",
		"1",
	)

	checkFn := resource.ComposeAggregateTestCheckFunc(
		append(baseChecks, passwordWoCheck)...,
	)

	expectedUpdateRoles := map[string]struct{}{"1": {}}
	updateChecks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"tenant_id",
			"1",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"username",
			"testacc-TestAccMorpheusUserUpdateOkChanged",
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
			"bar@hpe.com",
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
			"1",
		),
		checkRole(
			"hpe_morpheus_user.foo",
			"role_ids.0",
			expectedUpdateRoles,
		),
		checkStrayRoles(expectedUpdateRoles),
	}

	checkUpdateFn := resource.ComposeAggregateTestCheckFunc(
		updateChecks...,
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "hpe_morpheus_user" "foo" {
	# Assumes tenant_id 1 pre-exists
	tenant_id = 1
	username = "testacc-TestAccMorpheusUserUpdateOk"
	email = "foo@hpe.com"
	password_wo = "Secret123!"
	password_wo_version = 1
	role_ids = [3,1]
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
				Config: providerConfig + `
# checks plan has no effect
resource "hpe_morpheus_user" "foo" {
	tenant_id = 1
	username = "testacc-TestAccMorpheusUserUpdateOk"
	email = "foo@hpe.com"
	password_wo = "Secret123!"
	password_wo_version = 1
	role_ids = [3,1]
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
				Config: providerConfig + `
# checks plan detects first_name change to null
resource "hpe_morpheus_user" "foo" {
	tenant_id = 1
	username = "testacc-TestAccMorpheusUserUpdateOk"
	email = "foo@hpe.com"
	password_wo = "Secret123!"
	password_wo_version = 1
	role_ids = [3,1]
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
				Config: providerConfig + `
# checks plan detects first_name change
resource "hpe_morpheus_user" "foo" {
	tenant_id = 1
	username = "testacc-TestAccMorpheusUserUpdateOk"
	email = "foo@hpe.com"
	password_wo = "Secret123!"
	password_wo_version = 1
	role_ids = [3,1]
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
				Config: providerConfig + `
# checks plan detects last_name change to null
resource "hpe_morpheus_user" "foo" {
	tenant_id = 1
	username = "testacc-TestAccMorpheusUserUpdateOk"
	email = "foo@hpe.com"
	password_wo = "Secret123!"
	password_wo_version = 1
	role_ids = [3,1]
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
				Config: providerConfig + `
# checks plan detects last_name change
resource "hpe_morpheus_user" "foo" {
	tenant_id = 1
	username = "testacc-TestAccMorpheusUserUpdateOk"
	email = "foo@hpe.com"
	password_wo = "Secret123!"
	password_wo_version = 1
	role_ids = [3,1]
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
				Config: providerConfig + `
# checks plan detects password_wo_version to null
resource "hpe_morpheus_user" "foo" {
	tenant_id = 1
	username = "testacc-TestAccMorpheusUserUpdateOk"
	email = "foo@hpe.com"
	password_wo = "Secret123!"
	# changed
	# password_wo_version = 1
	role_ids = [3,1]
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
				Config: providerConfig + `
# checks plan detects no change if only password_wo is changed
resource "hpe_morpheus_user" "foo" {
	tenant_id = 1
	username = "testacc-TestAccMorpheusUserUpdateOk"
	email = "foo@hpe.com"
	# changed
	# password_wo = "Secret123!"
	password_wo_version = 1
	role_ids = [3,1]
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
				Config: providerConfig + `
# checks plan detects changed role_ids
resource "hpe_morpheus_user" "foo" {
	tenant_id = 1
	username = "testacc-TestAccMorpheusUserUpdateOk"
	email = "foo@hpe.com"
	password_wo = "Secret123!"
	password_wo_version = 1
	# changed
	role_ids = [3]
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
				Config: providerConfig + `
# checks plan detects changed username
resource "hpe_morpheus_user" "foo" {
	tenant_id = 1
	# changed
	username = "testacc-TestAccMorpheusUserUpdateOkNew"
	email = "foo@hpe.com"
	password_wo = "Secret123!"
	password_wo_version = 1
	role_ids = [3,1]
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
				Config: providerConfig + `
# checks plan detects changed windows username
resource "hpe_morpheus_user" "foo" {
	tenant_id = 1
	username = "testacc-TestAccMorpheusUserUpdateOk"
	email = "foo@hpe.com"
	password_wo = "Secret123!"
	password_wo_version = 1
	role_ids = [3,1]
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
				Config: providerConfig + `
# checks plan detects changed linux username
resource "hpe_morpheus_user" "foo" {
	tenant_id = 1
	username = "testacc-TestAccMorpheusUserUpdateOk"
	email = "foo@hpe.com"
	password_wo = "Secret123!"
	password_wo_version = 1
	role_ids = [3,1]
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
				Config: providerConfig + `
# checks plan detects changed linux password version
resource "hpe_morpheus_user" "foo" {
	tenant_id = 1
	username = "testacc-TestAccMorpheusUserUpdateOk"
	email = "foo@hpe.com"
	password_wo = "Secret123!"
	password_wo_version = 1
	role_ids = [3,1]
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
				Config: providerConfig + `
# checks plan detects changed windows password version
resource "hpe_morpheus_user" "foo" {
	tenant_id = 1
	username = "testacc-TestAccMorpheusUserUpdateOk"
	email = "foo@hpe.com"
	password_wo = "Secret123!"
	password_wo_version = 1
	role_ids = [3,1]
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
				Config: providerConfig + `
# checks plan detects changed linux key pair id
resource "hpe_morpheus_user" "foo" {
	tenant_id = 1
	username = "testacc-TestAccMorpheusUserUpdateOk"
	email = "foo@hpe.com"
	password_wo = "Secret123!"
	password_wo_version = 1
	role_ids = [3,1]
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
				Config: providerConfig + `
# checks apply of changes to all changeable fields
resource "hpe_morpheus_user" "foo" {
	tenant_id = 1
	# changed
	username = "testacc-TestAccMorpheusUserUpdateOkChanged"
	# changed
	email = "bar@hpe.com"
	# changed
	password_wo = "Secret456!"
	# changed
	password_wo_version = 2
	# changed
	role_ids = [1]
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
	receive_notifications = false
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
				Config: providerConfig + `
# checks plan has no effect
resource "hpe_morpheus_user" "foo" {
	tenant_id = 1
	# changed
	username = "testacc-TestAccMorpheusUserUpdateOkChanged"
	# changed
	email = "bar@hpe.com"
	# changed
	password_wo = "Secret456!"
	# changed
	password_wo_version = 2
	# changed
	role_ids = [1]
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
	receive_notifications = false
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

func TestAccMorpheusUserAllAttrsOk(t *testing.T) {
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	resourceCfg := `
# Role id 0 causes a test failure because it is ignored by
# the server and only the other two roles are created
#resource "hpe_morpheus_user" "bar" {
#username = "test101"
#email = "foo@hpe.com"
#password = "Secret123!"
#roles = [3,0,1]
#}
resource "hpe_morpheus_user" "foo" {
	# Assumes tenant_id 1 pre-exists
	tenant_id = 1
	username = "testacc-TestAccMorpheusUserAllAttrsOk"
	email = "foo@hpe.com"
	password_wo = "Secret123!"
	password_wo_version = 1
	role_ids = [3,1]
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
	expectedRoles := map[string]struct{}{"3": {}, "1": {}}

	baseChecks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"tenant_id",
			"1",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"username",
			"testacc-TestAccMorpheusUserAllAttrsOk",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"email",
			"foo@hpe.com",
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
			expectedRoles,
		),
		checkRole(
			"hpe_morpheus_user.foo",
			"role_ids.1",
			expectedRoles,
		),
		checkStrayRoles(expectedRoles),
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
				Config:   providerConfig + resourceCfg,
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
				},
				ResourceName: "hpe_morpheus_user.foo",
				Check:        checkImportFn,
			},
		},
	})
}

func TestAccMorpheusUserMissingRoles(t *testing.T) {
	defer testhelpers.RecordResult(t)
	providerConfig := `
provider "hpe" {
	morpheus {
		url = ""
		username = ""
		password = ""
	}
}

resource "hpe_morpheus_user" "foo" {
	username = "test2"
	email = "bar@hpe.com"
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

func TestAccMorpheusUserMissingUsername(t *testing.T) {
	defer testhelpers.RecordResult(t)
	providerConfig := `
provider "hpe" {
	morpheus {
		url = ""
		username = ""
		password = ""
	}
}

resource "hpe_morpheus_user" "foo" {
	#username = "test2"
	email = "bar@hpe.com"
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

func TestAccMorpheusUserMissingEmail(t *testing.T) {
	defer testhelpers.RecordResult(t)
	providerConfig := `
provider "hpe" {
	morpheus {
		url = ""
		username = ""
		password = ""
	}
}

resource "hpe_morpheus_user" "foo" {
	username = "test2"
	#email = "bar@hpe.com"
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
func TestAccMorpheusUserMissingPasswordWo(t *testing.T) {
	defer testhelpers.RecordResult(t)
	providerConfig := `
provider "hpe" {
	morpheus {
		url = ""
		username = ""
		password = ""
	}
}

resource "hpe_morpheus_user" "foo" {
	username = "test2"
	email = "bar@hpe.com"
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
func TestAccMorpheusUserImportOk(t *testing.T) {
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	// nolint: gosec
	resourceCfgWithPassword := `
resource "hpe_morpheus_user" "foo" {
	# Assumes tenant_id 1 pre-exists
	tenant_id = 1
	username = "testacc-TestAccMorpheusUserImportOk"
	email = "foo@hpe.com"
	password_wo = "Secret123!"
	password_wo_version = 1
	role_ids = [3,1]
	first_name = "foo"
	last_name = "bar"
	linux_username = "linus"
	linux_key_pair_id = 100
	receive_notifications = false
	windows_username = "bill"
}
`
	// nolint: gosec
	resourceCfgNoPassword := `
resource "hpe_morpheus_user" "foo" {
	# Assumes tenant_id 1 pre-exists
	tenant_id = 1
	username = "testacc-TestAccMorpheusUserImportOk"
	email = "foo@hpe.com"
        #password_wo = "Secret123!"
        #password_wo_version = 1
	role_ids = [3,1]
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
removed {
from = hpe_morpheus_user.foo

lifecycle {
destroy = false
}
}
`
	baseChecks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"tenant_id",
			"1",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"username",
			"testacc-TestAccMorpheusUserImportOk",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"email",
			"foo@hpe.com",
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

	expectedCreateRoles := map[string]struct{}{"3": {}, "1": {}}
	createChecks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"password_wo_version",
			"1",
		),
		resource.TestCheckNoResourceAttr(
			"hpe_morpheus_user.foo",
			"password_wo",
		),
		checkRole("hpe_morpheus_user.foo", "role_ids.0", expectedCreateRoles),
		checkRole("hpe_morpheus_user.foo", "role_ids.1", expectedCreateRoles),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(
		append(baseChecks, createChecks...)...,
	)

	var cachedID string

	// This is a new TestCase - we know for sure
	// we inherit no state from the TestCase above
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + resourceCfgWithPassword,
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						// Cache ID for use later
						rs := s.RootModule().Resources["hpe_morpheus_user.foo"]
						if rs == nil {
							return fmt.Errorf("resource not found")
						}
						cachedID = rs.Primary.ID

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

	importCfg := providerConfig + resourceCfgNoPassword + `
	import {
	  to = hpe_morpheus_user.foo
	  id = ` + cachedID + `
	}
	`
	expectedImportRoles := map[string]struct{}{"3": {}, "1": {}}
	importChecks := []resource.TestCheckFunc{
		resource.TestCheckNoResourceAttr(
			"hpe_morpheus_user.foo",
			"password_wo_version",
		),
		resource.TestCheckNoResourceAttr(
			"hpe_morpheus_user.foo",
			"password_wo",
		),
		checkRole("hpe_morpheus_user.foo", "role_ids.0", expectedImportRoles),
		checkRole("hpe_morpheus_user.foo", "role_ids.1", expectedImportRoles),
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
				Config:             providerConfig + resourceCfgNoPassword,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}
