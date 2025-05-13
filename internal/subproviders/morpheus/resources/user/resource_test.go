package user_test

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/HPE/terraform-provider-hpe/internal/provider"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func getExpectedFn(expected map[string]map[string]string) func(string, string) string {
	return func(id string, key string) string {
		return expected[id][key]
	}
}

func getExpected(index string, key string) string {
	expectedRoles := map[string]map[string]string{
		"3": {
			"id":          "3",
			"name":        "User Admin",
			"description": "Sub Tenant User Template",
			"authority":   "User Admin",
		},
		"1": {
			"id":          "1",
			"name":        "System Admin",
			"description": "Super User",
			"authority":   "System Admin",
		},
	}

	return expectedRoles[index][key]
}

func checkRole(
	resourceName string,
	listkey string,
	index int,
	keys []string,
	refkey string,
	f func(string, string) string,
) func(*terraform.State) error {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		indexStr := strconv.Itoa(index)
		refattr := fmt.Sprintf("%s.%s.%s", listkey, indexStr, refkey)
		refvalue := rs.Primary.Attributes[refattr]

		for _, key := range keys {
			attr := fmt.Sprintf("%s.%s.%s", listkey, indexStr, key)
			value := rs.Primary.Attributes[attr]
			expected := f(refvalue, key)
			if value != expected {
				msg := fmt.Sprintf(
					"expected '%s' got '%s'",
					expected,
					value,
				)

				return errors.New(msg)
			}
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

func TestAccMorpheusUserOk(t *testing.T) {
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
#email = "foo@bar.com"
#password = "Secret123!"
#roles = [
#	{
#		id = 3
#	},
#	{
#		id = 0
#	},
#	{
#		id = 1
#	}
#]
#}

resource "hpe_morpheus_user" "foo" {
	username = "testacc-1"
	email = "foo@bar.com"
	password = "Secret123!"
	roles = [
		{
			id = 3
		},
		{
			id = 1
		}
	]
}
`
	expectedRoles := map[string]map[string]string{
		"3": {
			"id":          "3",
			"name":        "User Admin",
			"description": "Sub Tenant User Template",
			"authority":   "User Admin",
		},
		"1": {
			"id":          "1",
			"name":        "System Admin",
			"description": "Super User",
			"authority":   "System Admin",
		},
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"username",
			"test2",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"email",
			"foo@bar.com",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"password",
			"Secret123!",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_user.foo",
			"roles.#",
			"2",
		),
		checkRole(
			"hpe_morpheus_user.foo",
			"roles",
			0,
			[]string{"authority", "description", "name", "id"},
			"id",
			getExpectedFn(expectedRoles),
		),
		checkRole(
			"hpe_morpheus_user.foo",
			"roles",
			1,
			[]string{"authority", "description", "name", "id"},
			"id",
			getExpected,
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

func TestAccMorpheusUserMissingRoles(t *testing.T) {
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
	username = "test2"
	email = "bar@bar.com"
	password = "Secret123!"
}
`
	expected := `The argument "roles" is required`

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
	#username = "test2"
	email = "bar@bar.com"
	password = "Secret123!"
	roles = [
		{
			id = 3
		},
		{
			id = 1
		}
	]
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
	username = "test2"
	#email = "bar@bar.com"
	password = "Secret123!"
	roles = [
		{
			id = 3
		},
		{
			id = 1
		}
	]
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
