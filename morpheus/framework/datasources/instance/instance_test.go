// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

//go:generate ../../../../bin/render example-id.tf.tmpl Id 99
//go:generate ../../../../bin/render example-name.tf.tmpl Name "HVM Instance"
package instance_test

import (
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/provider"
)

func TestMain(m *testing.M) {
	code := m.Run()
	testhelpers.WriteMergedResults()
	os.Exit(code)
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

// Combining tests for both ID and Name under the same function as instance
// creation can be time consuming, so we only want to create one instance.
func TestAccMorpheusInstanceDatasourceByIdAndName(t *testing.T) {
	defer testhelpers.RecordResult(t)
	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())
	instanceTypeID := "9"
	resourcePool := "pool-62299"
	resourceConfig, err := testhelpers.RenderExample(t, "../../resources/instance/example.tf.tmpl",
		"Name", name,
		"InstanceType", instanceTypeID,
		"ResourcePool", resourcePool,
	)
	if err != nil {
		t.Fatal(err)
	}

	dataSourceConfigWithId, err := testhelpers.RenderExample(t,
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
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_instance.example",
			"config.noAgent",
			"true",
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_instance.example",
			"config.resourcePoolId",
			resourcePool,
		),
		testhelpers.CheckResourceAttrEqual(
			"data.hpe_morpheus_instance.example", "instance_type.id",
			"hpe_morpheus_instance.example", "instance_type_id",
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_instance.example",
			"instance_type.id",
			instanceTypeID,
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_instance.example",
			"container_details.#",
			"1",
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_instance.example",
			"container_details.0.server.agent_installed",
			"false",
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_instance.example",
			"container_details.0.server.resource_pool_id",
			strings.TrimPrefix(resourcePool, "pool-"),
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
				Config: providerConfig + resourceConfig + dataSourceConfigWithId,
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
