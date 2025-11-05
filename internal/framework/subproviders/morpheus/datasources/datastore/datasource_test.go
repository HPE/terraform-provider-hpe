// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package datastore_test

//go:generate go run ../../../../../../cmd/render example-id.tf.tmpl Id 99
//go:generate go run ../../../../../../cmd/render example-name.tf.tmpl Name "\"Example name\""

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/internal/framework/provider"
	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus"
	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
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

func testCheckFns(name string) []resource.TestCheckFunc {
	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_datastore.test",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_datastore.test",
			"associated_resource_type",
			"Cluster",
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_datastore.test",
			"associated_resource_id",
			"6032",
		),
		resource.TestCheckTypeSetElemNestedAttrs(
			"data.hpe_morpheus_datastore.test",
			"resource_permissions.groups.*",
			map[string]string{
				"id": "17830",
			},
		),
		resource.TestCheckTypeSetElemNestedAttrs(
			"data.hpe_morpheus_datastore.test",
			"tenants.*",
			map[string]string{
				"id": "466",
			},
		),
		resource.TestCheckTypeSetElemNestedAttrs(
			"data.hpe_morpheus_datastore.test",
			"tenants.*",
			map[string]string{
				"id": "1",
			},
		),
	}

	return checks
}

func TestAccMorpheusFindDatastoreById(t *testing.T) {
	defer testhelpers.RecordResult(t)
	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	name := acctest.RandomWithPrefix(t.Name())

	providerConfig := testhelpers.ProviderBlock()

	datastoreResouceConfig, err := testhelpers.RenderExample(t, "example_alletramp_hvm.tf.tmpl",
		"Name", name,
		"AssociatedResourceID", "6032",
		"StorageServerID", "1489",
		"GroupID", "17830",
		"TenantID", "466",
	)
	if err != nil {
		t.Fatal(err)
	}

	dataSourceConfig, err := testhelpers.RenderExample(t,
		"example-id.tf.tmpl", "Id", "hpe_morpheus_datastore.example.id")
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + datastoreResouceConfig + dataSourceConfig,
				Check:  resource.ComposeAggregateTestCheckFunc(testCheckFns(name)...),
			},
		},
	})
}

func TestAccMorpheusFindDatastoreByName(t *testing.T) {
	defer testhelpers.RecordResult(t)
	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	name := acctest.RandomWithPrefix(t.Name())

	providerConfig := testhelpers.ProviderBlock()

	datastoreResouceConfig, err := testhelpers.RenderExample(t, "example_alletramp_hvm.tf.tmpl",
		"Name", name,
		"AssociatedResourceID", "6032",
		"StorageServerID", "1489",
		"GroupID", "17830",
		"TenantID", "466",
	)
	if err != nil {
		t.Fatal(err)
	}

	dataSourceConfig, err := testhelpers.RenderExample(t,
		"example-name.tf.tmpl", "Name", "hpe_morpheus_datastore.example.name")
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + datastoreResouceConfig + dataSourceConfig,
				Check:  resource.ComposeAggregateTestCheckFunc(testCheckFns(name)...),
			},
		},
	})
}

func TestAccMorpheusFindDatastoreNotFound(t *testing.T) {
	defer testhelpers.RecordResult(t)
	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	config := providerConfig + `
      data "hpe_morpheus_datastore" "test" {
        name = "blah" 
      }`

	checks := []resource.TestCheckFunc{
		resource.TestCheckNoResourceAttr(
			"data.hpe_morpheus_datastore.test",
			"id",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	expected := fmt.Sprint("datastore blah list failed")

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

func TestAccMorpheusFindDatastoreNoSearchAttrs(t *testing.T) {
	defer testhelpers.RecordResult(t)
	t.Parallel()

	config := providerConfigOffline + `
      data "hpe_morpheus_datastore" "test" {
      }`

	checks := []resource.TestCheckFunc{
		resource.TestCheckNoResourceAttr(
			"data.hpe_morpheus_datastore.test",
			"id",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	expected := fmt.Sprintf("either id or name must be specified")

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

func TestAccMorpheusFindDatastoreBothSearchAttrs(t *testing.T) {
	defer testhelpers.RecordResult(t)
	t.Parallel()

	config := providerConfigOffline + `
      data "hpe_morpheus_datastore" "test" {
        id = 1
        name = "blah" 
      }`

	checks := []resource.TestCheckFunc{
		resource.TestCheckNoResourceAttr(
			"data.hpe_morpheus_datastore.test",
			"id",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	expected := fmt.Sprint("Error running pre-apply plan: exit status 1")

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
