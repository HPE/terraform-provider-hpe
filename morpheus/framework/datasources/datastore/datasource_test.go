// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package datastore_test

//go:generate ../../../../bin/render example-id.tf.tmpl Id 99
//go:generate ../../../../bin/render example-name.tf.tmpl Name "\"Example name\""

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/provider"
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
	checks := testCheckFnsNoResourcePermissions(name)
	checksResourcePermissions := []resource.TestCheckFunc{
		resource.TestCheckTypeSetElemNestedAttrs(
			"data.hpe_morpheus_datastore.test",
			"resource_permissions.groups.*",
			map[string]string{
				"id": "17830",
			},
		),
		// it seems that adding resource permissions results in the root tenant being added
		// for our tests with TerraformTester on feature
		resource.TestCheckTypeSetElemNestedAttrs(
			"data.hpe_morpheus_datastore.test",
			"tenants.*",
			map[string]string{
				"id": "1",
			},
		),
	}
	checks = append(checks, checksResourcePermissions...)

	return checks
}

func testCheckFnsNoResourcePermissions(name string) []resource.TestCheckFunc {
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
			"tenants.*",
			map[string]string{
				"id": "466",
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

func TestAccMorpheusFindDatastoreNoResourcePermissionsById(t *testing.T) {
	defer testhelpers.RecordResult(t)
	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	name := acctest.RandomWithPrefix(t.Name())

	providerConfig := testhelpers.ProviderBlock()

	datastoreResouceConfig, err := testhelpers.RenderExample(t, "example_alletramp_hvm_no_resource_permissions.tf.tmpl",
		"Name", name,
		"AssociatedResourceID", "6032",
		"StorageServerID", "1489",
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
				Check:  resource.ComposeAggregateTestCheckFunc(testCheckFnsNoResourcePermissions(name)...),
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

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      config,
				Check:       checkFn,
				ExpectError: regexp.MustCompile("datastore blah list failed"),
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

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      config,
				Check:       checkFn,
				ExpectError: regexp.MustCompile("either id or name must be specified"),
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

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      config,
				Check:       checkFn,
				ExpectError: regexp.MustCompile("Error running pre-apply plan: exit status 1"),
			},
		},
	})
}
