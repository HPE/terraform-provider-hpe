// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package cloud_test

import (
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	morpheus "github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/cloud/consts"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
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

func TestAccMorpheusFindCloudById(t *testing.T) {
	if capabilities.Missing(t, capabilities.All) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)
	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	name := acctest.RandomWithPrefix(t.Name())

	providerConfig := testhelpers.ProviderBlock()

	cloudResourceConfig := `
resource "hpe_morpheus_cloud" "test_cloud" {
  # Required fields
  name      = "` + name + `"
  tenant_id = 1
  group_id  = 1


  # General configuration
  code             = "` + strings.ToLower(name) + `"
  labels           = ["aLabel1", "aLabel2"]
  data_center_name = "aDatacenter"
  enabled          = true
  location         = "somewhere"
  visibility       = "public"

  # Agent and provisioning settings
  agent_install_mode       = "ssh"
  appliance_url            = "https://somewhere.com"
  auto_recover_power_state = true
  import_existing_vms      = "off"

  # Cost and guidance settings
  costing_mode  = "costing"
  guidance_mode = "off"

  # Security settings
  security_mode = "off"

  # Console settings
  keyboard_layout = "us"

  # HVM-specific configuration
  config_hvm = {
    certificate_provider          = "internal"
    enable_network_type_selection = false
  }
}
`

	dataSourceConfig, err := testhelpers.RenderExample(t,
		"example-id.tf.tmpl", "Id", "hpe_morpheus_cloud.test_cloud.id")
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_cloud.example",
			"name",
			name,
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + cloudResourceConfig + dataSourceConfig,
				Check:  checkFn,
			},
		},
	})
}

func TestAccMorpheusFindCloudByName(t *testing.T) {
	if capabilities.Missing(t, capabilities.All) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)
	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	name := acctest.RandomWithPrefix(t.Name())

	providerConfig := testhelpers.ProviderBlock()

	cloudResourceConfig := `
resource "hpe_morpheus_cloud" "test_cloud" {
  # Required fields
  name      = "` + name + `"
  tenant_id = 1
  group_id  = 1


  # General configuration
  code             = "` + strings.ToLower(name) + `"
  labels           = ["aLabel1", "aLabel2"]
  data_center_name = "aDatacenter"
  enabled          = true
  location         = "somewhere"
  visibility       = "public"

  # Agent and provisioning settings
  agent_install_mode       = "ssh"
  appliance_url            = "https://somewhere.com"
  auto_recover_power_state = true
  import_existing_vms      = "off"

  # Cost and guidance settings
  costing_mode  = "costing"
  guidance_mode = "off"

  # Security settings
  security_mode = "off"

  # Console settings
  keyboard_layout = "us"

  # HVM-specific configuration
  config_hvm = {
    certificate_provider          = "internal"
    enable_network_type_selection = false
  }
}
`

	dataSourceConfig, err := testhelpers.RenderExample(t,
		"example-name.tf.tmpl", "Name", "hpe_morpheus_cloud.test_cloud.name")
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_cloud.example",
			"name",
			name,
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + cloudResourceConfig + dataSourceConfig,
				Check:  checkFn,
			},
		},
	})
}

func TestAccMorpheusFindCloudNotFound(t *testing.T) {
	if capabilities.Missing(t, capabilities.All) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)
	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	config := providerConfig + `
      data "hpe_morpheus_cloud" "test" {
        name = "______" 
      }`

	checks := []resource.TestCheckFunc{
		resource.TestCheckNoResourceAttr(
			"data.hpe_morpheus_cloud.test",
			"id",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	expected := consts.ErrorNoCloudFound

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

func TestAccMorpheusFindCloudNoSearchAttrs(t *testing.T) {
	if capabilities.Missing(t, capabilities.All) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)
	t.Parallel()

	config := providerConfigOffline + `
      data "hpe_morpheus_cloud" "test" {
      }`

	checks := []resource.TestCheckFunc{
		resource.TestCheckNoResourceAttr(
			"data.hpe_morpheus_cloud.test",
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

func TestAccMorpheusFindCloudBothSearchAttrs(t *testing.T) {
	if capabilities.Missing(t, capabilities.All) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)
	t.Parallel()

	config := providerConfigOffline + `
      data "hpe_morpheus_cloud" "test" {
        id = 1
        name = "______" 
      }`

	checks := []resource.TestCheckFunc{
		resource.TestCheckNoResourceAttr(
			"data.hpe_morpheus_cloud.test",
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
