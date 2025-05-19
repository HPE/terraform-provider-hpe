package user_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/HPE/terraform-provider-hpe/internal/provider"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
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

resource "hpe_morpheus_user" "foo" {
	username = "testacc-TestAccMorpheusUserRequiredAttrsOk"
	email = "foo@hpe.com"
	password = "Secret123!"
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
			"password",
			"Secret123!",
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
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:             providerConfig,
				ExpectNonEmptyPlan: false,
				Check:              checkFn,
				PlanOnly:           false,
			},
			{
				ImportState: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					// Read ID from the pre-import state
					rs := s.RootModule().
						Resources["hpe_morpheus_user.foo"]

					return rs.Primary.ID + "," + "Secret123!", nil
				},
				ImportStateVerify: true, // Check state post import
				ResourceName:      "hpe_morpheus_user.foo",
				Check:             checkFn,
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

resource "hpe_morpheus_user" "foo" {
	username = "testacc-TestAccMorpheusUserAllAttrsOk"
	email = "foo@hpe.com"
	password = "Secret123!"
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

	checks := []resource.TestCheckFunc{
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
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"password",
			"Secret123!",
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

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:             providerConfig,
				ExpectNonEmptyPlan: false,
				Check:              checkFn,
				PlanOnly:           false,
			},
			{
				ImportState: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					// Read ID from the pre-import state
					rs := s.RootModule().
						Resources["hpe_morpheus_user.foo"]

					return rs.Primary.ID + "," + "Secret123!", nil
				},
				ImportStateVerify: true, // Check state post import
				ResourceName:      "hpe_morpheus_user.foo",
				Check:             checkFn,
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
