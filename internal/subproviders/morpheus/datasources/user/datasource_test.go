// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package user_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/internal/provider"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/datasources/environment"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/datasources/user/consts"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/testhelpers"
)

const providerConfigOffline = `
provider "hpe" {
  morpheus {
    url          = ""
    username     = ""
    password     = ""
  }
}
`

func TestMain(m *testing.M) {
	code := m.Run()
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func TestAccMorpheusFindUserByUsername(t *testing.T) {
	defer testhelpers.RecordResult(t)
	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	username := acctest.RandomWithPrefix(t.Name())
	providerConfig := testhelpers.ProviderBlock()

	userResourceConfig := `
resource "hpe_morpheus_user" "test_user" {
	username = "` + username + `"
	role_ids = [1]
	email    = "` + username + `@example.com"
	password_wo = "Test123!!"
}
`

	dataSourceConfig, err := testhelpers.RenderExample(t,
		"example-username.tf.tmpl", "Username", username)
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_user.test",
			"username",
			username,
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + userResourceConfig,
			},
			{
				Config: providerConfig + userResourceConfig + dataSourceConfig,
				Check:  checkFn,
			},
		},
	})
}

func newProviderWithError() (tfprotov6.ProviderServer, error) {
	providerInstance := provider.New("test", morpheus.New())()

	return providerserver.NewProtocol6WithError(providerInstance)()
}

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"hpe": newProviderWithError,
}

func TestAccMorpheusFindUserById(t *testing.T) {
	defer testhelpers.RecordResult(t)
	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	username := acctest.RandomWithPrefix(t.Name())

	providerConfig := testhelpers.ProviderBlock()

	userResourceConfig := `
resource "hpe_morpheus_user" "test_user" {
	username = "` + username + `"
	role_ids = [1]
	email    = "` + username + `@example.com"
	password_wo = "Test123!!"
}
`

	dataSourceConfig := `
    data "hpe_morpheus_user" "test" {
        id = hpe_morpheus_user.test_user.id
    }
    `

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_user.test",
			"username",
			username,
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + userResourceConfig + dataSourceConfig,
				Check:  checkFn,
			},
		},
	})
}

func TestAccMorpheusFindUserNotFound(t *testing.T) {
	defer testhelpers.RecordResult(t)
	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	config := providerConfig + `
      data "hpe_morpheus_user" "test" {
        username = "______"
      }`

	checks := []resource.TestCheckFunc{
		resource.TestCheckNoResourceAttr(
			"data.hpe_morpheus_user.test",
			"id",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	expected := consts.ErrorNoUserFound

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      config,
				Check:       checkFn,
				ExpectError: regexp.MustCompile(expected),
			},
		},
	})
}

func TestAccMorpheusFindUserNoSearchAttrs(t *testing.T) {
	defer testhelpers.RecordResult(t)
	t.Parallel()

	config := providerConfigOffline + `
      data "hpe_morpheus_user" "test" {
      }`

	checks := []resource.TestCheckFunc{
		resource.TestCheckNoResourceAttr(
			"data.hpe_morpheus_user.test",
			"id",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	expected := consts.ErrorNoValidUserTerms

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      config,
				Check:       checkFn,
				ExpectError: regexp.MustCompile(expected),
			},
		},
	})
}

func TestAccMorpheusFindUserBothSearchAttrs(t *testing.T) {
	defer testhelpers.RecordResult(t)
	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	config := providerConfig + `
      data "hpe_morpheus_user" "test" {
        id = "1"
        username = "testuser"
      }`

	checks := []resource.TestCheckFunc{
		resource.TestCheckNoResourceAttr(
			"data.hpe_morpheus_user.test",
			"id",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	expected := environment.ErrorRunningPreApply

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      config,
				Check:       checkFn,
				ExpectError: regexp.MustCompile(expected),
			},
		},
	})
}

// Test to verify that all of the attributes from a created user can be read
func TestAccMorpheusFindUserVerifyAttributes(t *testing.T) {
	defer testhelpers.RecordResult(t)
	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	username := acctest.RandomWithPrefix(t.Name())
	email := username + "@example.com"
	firstName := "TestFirst"
	lastName := "TestLast"

	providerConfig := testhelpers.ProviderBlock()

	resourceConfig := `
resource "hpe_morpheus_user" "test_all" {
  username     = "` + username + `"
  email        = "` + email + `"
  first_name   = "` + firstName + `"
  last_name    = "` + lastName + `"
  role_ids     = [1]
  password_wo  = "Test123!!"
  receive_notifications = true
  linux_username = "` + username + `"
  windows_username = "` + username + `"
}
`

	dataSourceConfig := `
data "hpe_morpheus_user" "test_all" {
  username = "` + username + `"
}
`

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_user.test_all",
			"username",
			username,
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_user.test_all",
			"email",
			email,
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_user.test_all",
			"first_name",
			firstName,
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_user.test_all",
			"last_name",
			lastName,
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_user.test_all",
			"receive_notifications",
			"true",
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_user.test_all",
			"linux_username",
			username,
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_user.test_all",
			"windows_username",
			username,
		),
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_user.test_all",
			"id",
		),
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_user.test_all",
			"tenant_id",
		),
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_user.test_all",
			"roles.#",
		),
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_user.test_all",
			"account.id",
		),
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_user.test_all",
			"account.name",
		),
		/* Can't test this yet as we cannot assign a default persona via terraform yet
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_user.test_all",
			"default_persona.id",
		),
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_user.test_all",
			"default_persona.code",
		),
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_user.test_all",
			"default_persona.name",
		),
		*/
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_user.test_all",
			"access.blueprints.#",
		),
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_user.test_all",
			"access.catalog_item_types.#",
		),
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_user.test_all",
			"access.features.#",
		),
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_user.test_all",
			"access.instance_types.#",
		),
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_user.test_all",
			"access.personas.#",
		),
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_user.test_all",
			"access.report_types.#",
		),
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_user.test_all",
			"access.groups.#",
		),
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_user.test_all",
			"access.workflows.#",
		),
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_user.test_all",
			"access.tasks.#",
		),
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_user.test_all",
			"access.vdi_pools.#",
		),
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_user.test_all",
			"access.clouds.#",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				ExpectNonEmptyPlan: false,
				Config:             providerConfig + resourceConfig,
			},
			{
				ExpectNonEmptyPlan: false,
				Config:             providerConfig + resourceConfig + dataSourceConfig,
				Check:              checkFn,
			},
		},
	})
}
