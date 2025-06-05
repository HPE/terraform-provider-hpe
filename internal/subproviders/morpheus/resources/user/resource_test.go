// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package user_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/HPE/terraform-provider-hpe/internal/provider"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus"
)

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

// Check that we can create a user with only
// required attributes specified
func TestAccMorpheusUserRequiredAttrsOk(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := `
variable "testacc_morpheus_url" {}
variable "testacc_morpheus_username" {}
variable "testacc_morpheus_password" {}
variable "testacc_morpheus_insecure" {}

provider "hpe" {
	morpheus {
		url = var.testacc_morpheus_url
		username = var.testacc_morpheus_username
		password = var.testacc_morpheus_password
		insecure = var.testacc_morpheus_insecure
	}
}
`
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

func TestAccMorpheusUserAllAttrsOk(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := `
variable "testacc_morpheus_url" {}
variable "testacc_morpheus_username" {}
variable "testacc_morpheus_password" {}
variable "testacc_morpheus_insecure" {}

provider "hpe" {
	morpheus {
		url = var.testacc_morpheus_url
		username = var.testacc_morpheus_username
		password = var.testacc_morpheus_password
		insecure = var.testacc_morpheus_insecure
	}
}

# Role id 0 causes a test failure because it is ignored by
# the server and only the other two roles are created
#resource "hpe_morpheus_user" "bar" {
#username = "test101"
#email = "foo@hpe.com"
#password = "Secret123!"
#roles = [3,0,1]
#}
`

	resourceCfg := `
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
	linux_key_pair_id = 100
	receive_notifications = false
	windows_username = "bill"
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

	passwordWoImportCheck := resource.TestCheckNoResourceAttr(
		"hpe_morpheus_user.foo",
		"password_wo_version",
	)

	checkFn := resource.ComposeAggregateTestCheckFunc(
		append(baseChecks, passwordWoCheck)...,
	)

	checkImportFn := resource.ComposeAggregateTestCheckFunc(
		append(baseChecks, passwordWoImportCheck)...,
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
				ImportStateVerify:       true, // Check state post import (in memory)
				ImportStateVerifyIgnore: []string{"password_wo_version"},
				ResourceName:            "hpe_morpheus_user.foo",
				Check:                   checkImportFn,
			},
		},
	})
}

func TestAccMorpheusUserMissingRoles(t *testing.T) {
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
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := `
variable "testacc_morpheus_url" {}
variable "testacc_morpheus_username" {}
variable "testacc_morpheus_password" {}
variable "testacc_morpheus_insecure" {}

provider "hpe" {
	morpheus {
		url = var.testacc_morpheus_url
		username = var.testacc_morpheus_username
		password = var.testacc_morpheus_password
		insecure = var.testacc_morpheus_insecure
	}
}
`
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
