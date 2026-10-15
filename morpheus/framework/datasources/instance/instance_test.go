// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

//go:generate ../../../../bin/render example-id.tf.tmpl Id 99
//go:generate ../../../../bin/render example-name.tf.tmpl Name "HVM Instance"
package instance_test

import (
	"os"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/instance"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

func TestMain(m *testing.M) {
	code := m.Run()
	testhelpers.WriteMergedResults()
	os.Exit(code)
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

// Combining tests for both ID and Name under the same function as instance
// creation can be time consuming, so we only want to create one instance.
func TestAccMorpheusInstanceDatasourceByIdAndName(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := instance.RenderInstanceConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	dataSourceConfigWithID, err := testhelpers.RenderExample(t,
		"example-id.tf.tmpl", "Id", "hpe_morpheus_instance.example.id")
	if err != nil {
		t.Fatal(err)
	}

	dataSourceConfigWithName, err := testhelpers.RenderExample(t,
		"example-name.tf.tmpl", "Name", name)
	if err != nil {
		t.Fatal(err)
	}

	// check some random fields to make sure we get the expected data from the
	// datasource
	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_instance.example",
			"name",
			name,
		),
		// Whether the VM has the agent follows the HVM fixture.
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_instance.example",
			"config.noAgent",
			strconv.FormatBool(testhelpers.HVMNoAgent()),
		),
		// The pool is the cluster's, resolved by the configuration's data
		// source; the instance type is the layout's.
		testhelpers.CheckResourceAttrEqual(
			"data.hpe_morpheus_instance.example", "config.resourcePoolId",
			"hpe_morpheus_instance.example", "config_hvm.resource_pool_id",
		),
		testhelpers.CheckResourceAttrEqual(
			"data.hpe_morpheus_instance.example", "instance_type.id",
			"hpe_morpheus_instance.example", "instance_type_id",
		),
		resource.TestCheckResourceAttrPair(
			"data.hpe_morpheus_instance.example", "instance_type.id",
			"data.hpe_morpheus_instance_type_layout.vme", "instance_type.id",
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_instance.example",
			"container_details.#",
			"1",
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_instance.example",
			"container_details.0.server.agent_installed",
			strconv.FormatBool(!testhelpers.HVMNoAgent()),
		),
		// The server's pool is the bare id of the cluster's pool.
		resource.TestCheckResourceAttrPair(
			"data.hpe_morpheus_instance.example", "container_details.0.server.resource_pool_id",
			"data.hpe_morpheus_cluster.vme", "permissions.resource_pool.id",
		),
		testhelpers.CheckResourceAttrEqual(
			"data.hpe_morpheus_instance.example", "container_details.0.server.container_plan.id",
			"hpe_morpheus_instance.example", "plan_id",
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_instance.example",
			"volumes.#",
			"3",
		),
		testhelpers.CheckResourceAttrEqual(
			"data.hpe_morpheus_instance.example", "volumes.0.datastore_id",
			"hpe_morpheus_instance.example", "volumes.0.datastore_id",
		),
		testhelpers.CheckResourceAttrEqual(
			"data.hpe_morpheus_instance.example", "interfaces.0.network.id",
			"hpe_morpheus_instance.example", "network_interfaces.0.network_id",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + resourceConfig,
			},
			{
				Config: providerConfig + resourceConfig + dataSourceConfigWithID,
				Check:  checkFn,
			},
			{
				Config: providerConfig + resourceConfig + dataSourceConfigWithName,
				Check:  checkFn,
			},
		},
	})
}

// this should fail due to a conflict between id and name
func TestAccMorpheusInstanceDatasourceBothAttrs(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	dataSourceConfig := `
data "hpe_morpheus_instance" "test" {
  name = "Example Instance"
  id = 5
}
`

	errMatch := regexp.MustCompile("Attribute \"(.*)\" cannot be specified when \"id\" is specified")
	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      providerConfig + dataSourceConfig,
				ExpectError: errMatch,
			},
		},
	})
}

// this should fail due to id or name being required
func TestAccMorpheusInstanceDatasourceNoAttrs(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	dataSourceConfig := `
data "hpe_morpheus_instance" "test" {
}
`

	errMatch := regexp.MustCompile("At least one attribute out of (.*) must be specified")
	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      providerConfig + dataSourceConfig,
				ExpectError: errMatch,
			},
		},
	})
}
