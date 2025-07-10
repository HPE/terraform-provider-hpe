// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

//go:build experimental

package role_test

//go:generate go run ../../../../../cmd/render example-id.tf.tmpl Id 99
//go:generate go run ../../../../../cmd/render example-name.tf.tmpl Name "\"Example name\""

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
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/datasources/cloud/consts"
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

func newProviderWithError() (tfprotov6.ProviderServer, error) {
	providerInstance := provider.New("test", morpheus.New())()

	return providerserver.NewProtocol6WithError(providerInstance)()
}

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"hpe": newProviderWithError,
}

func TestAccMorpheusFindRoleById(t *testing.T) {
	defer testhelpers.RecordResult(t)
	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	name := acctest.RandomWithPrefix(t.Name())

	providerConfig := testhelpers.ProviderBlock()

	resourceConfig := `
resource "hpe_morpheus_role" "test" {
  name = "` + name + `"
}
`
	dataSourceConfig, err := testhelpers.RenderExample(t,
		"example-id.tf.tmpl", "Id", "hpe_morpheus_role.test.id")
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_role.test",
			"name",
			name,
		),
		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_role.test",
			"id",
			"data.hpe_morpheus_role.test",
			"id",
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

func TestAccMorpheusFindRoleByName(t *testing.T) {
	defer testhelpers.RecordResult(t)
	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	name := acctest.RandomWithPrefix(t.Name())

	providerConfig := testhelpers.ProviderBlock()

	resourceConfig := `
resource "hpe_morpheus_role" "test" {
  name = "` + name + `"
}
`
	dataSourceConfig, err := testhelpers.RenderExample(t,
		"example-name.tf.tmpl", "Name", "hpe_morpheus_role.test.name")
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_role.test",
			"name",
			name,
		),
		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_role.test",
			"id",
			"data.hpe_morpheus_role.test",
			"id",
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

func TestAccMorpheusFindRoleVerifyAttributes(t *testing.T) {
	defer testhelpers.RecordResult(t)
	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	name := acctest.RandomWithPrefix(t.Name())

	providerConfig := testhelpers.ProviderBlock()

	resourceConfig := `
resource "hpe_morpheus_role" "test" {
  name = "` + name + `"
  description = "test"
  landing_url = "https://test.morpheus.com"
  multitenant = false
  multitenant_locked = false
  role_type = "account"
}
`
	dataSourceConfig := `
data "hpe_morpheus_role" "test" {
  name = "` + name + `"
}
`

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_role.test",
			"name",
			name,
		),
		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_role.test",
			"id",
			"data.hpe_morpheus_role.test",
			"id",
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_role.test",
			"description",
			"test",
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_role.test",
			"landing_url",
			"https://test.morpheus.com",
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_role.test",
			"multitenant",
			"false",
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_role.test",
			"multitenant_locked",
			"false",
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_role.test",
			"role_type",
			"account",
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

func TestAccMorpheusFindRoleNoSearchAttrs(t *testing.T) {
	defer testhelpers.RecordResult(t)
	t.Parallel()

	config := providerConfigOffline + `
      data "hpe_morpheus_role" "test" {
      }`

	checks := []resource.TestCheckFunc{
		resource.TestCheckNoResourceAttr(
			"data.hpe_morpheus_role.test",
			"id",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	expected := consts.ErrorNoValidSearchTerms

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

func TestAccMorpheusFindRoleBothSearchAttrs(t *testing.T) {
	defer testhelpers.RecordResult(t)
	t.Parallel()

	config := providerConfigOffline + `
      data "hpe_morpheus_role" "test" {
        id = 1
        name = "______"
      }`

	checks := []resource.TestCheckFunc{
		resource.TestCheckNoResourceAttr(
			"data.hpe_morpheus_role.test",
			"id",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	expected := consts.ErrorRunningPreApply

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
