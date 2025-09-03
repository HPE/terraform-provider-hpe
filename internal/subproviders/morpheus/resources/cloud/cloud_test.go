// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

//go:generate go run ../../../../../cmd/render example.tf.tmpl Name "TestCloud" Code "aCode" Label "aLabel"

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

	"github.com/HPE/terraform-provider-hpe/internal/provider"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/testhelpers"
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

// Tests that our example file template used for docs is a valid config
func TestAccMorpheusCloudExampleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())
	code := strings.ToLower(name)

	resourceConfig, err := testhelpers.RenderExample(t, "example.tf.tmpl",
		"Name", name,
		"Code", code,
		"Label", "aLabel",
	)
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cloud.example",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cloud.example",
			"code",
			code,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cloud.example",
			"labels.#",
			"2",
		),
		resource.TestCheckTypeSetElemAttr(
			"hpe_morpheus_cloud.example",
			"labels.*",
			"aLabel1",
		),
		resource.TestCheckTypeSetElemAttr(
			"hpe_morpheus_cloud.example",
			"labels.*",
			"aLabel2",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cloud.example",
			"agent_install_mode",
			"ssh",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cloud.example",
			"appliance_url",
			"https://somewhere.com",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cloud.example",
			"auto_recover_power_state",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cloud.example",
			"costing_mode",
			"costing",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cloud.example",
			"data_center_name",
			"aDatacenter",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cloud.example",
			"enabled",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cloud.example",
			"external_id",
			code,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cloud.example",
			"guidance_mode",
			"off",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cloud.example",
			"import_existing_vms",
			"off",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cloud.example",
			"keyboard_layout",
			"us",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cloud.example",
			"location",
			"somewhere",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cloud.example",
			"security_mode",
			"off",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cloud.example",
			"visibility",
			"public",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cloud.example",
			"config_hvm.certificate_provider",
			"internal",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cloud.example",
			"config_hvm.enable_network_type_selection",
			"false",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				Check:              checkFn,
				PlanOnly:           false,
			},
			{
				ImportState:       true,
				ImportStateVerify: true, // Check state post import
				ResourceName:      "hpe_morpheus_cloud.example",
				Check:             checkFn,
			},
		},
	})
}

func TestAccMorpheusCloudUpdate(t *testing.T) {
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())
	code := strings.ToLower(name)

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cloud.example",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cloud.example",
			"code",
			code,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cloud.example",
			"labels.#",
			"2",
		),
		resource.TestCheckTypeSetElemAttr(
			"hpe_morpheus_cloud.example",
			"labels.*",
			"Label1",
		),
		resource.TestCheckTypeSetElemAttr(
			"hpe_morpheus_cloud.example",
			"labels.*",
			"Label2",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cloud.example",
			"agent_install_mode",
			"ssh",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cloud.example",
			"appliance_url",
			"https://somewhere.com",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cloud.example",
			"auto_recover_power_state",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cloud.example",
			"costing_mode",
			"costing",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cloud.example",
			"data_center_name",
			"aDatacenter",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cloud.example",
			"enabled",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cloud.example",
			"external_id",
			code,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cloud.example",
			"guidance_mode",
			"off",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cloud.example",
			"import_existing_vms",
			"off",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cloud.example",
			"keyboard_layout",
			"us",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cloud.example",
			"location",
			"somewhere",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cloud.example",
			"security_mode",
			"off",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cloud.example",
			"visibility",
			"public",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cloud.example",
			"config_hvm.certificate_provider",
			"internal",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cloud.example",
			"config_hvm.enable_network_type_selection",
			"false",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(
		checks...,
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
					resource "hpe_morpheus_cloud" "example" {
						name      = "` + name + `"
						tenant_id = 1
						group_id  = 1

						code                     = "` + code + `"
						external_id              = "` + code + `"
						labels                   = ["Label1", "Label2"]
						agent_install_mode       = "ssh"
						appliance_url            = "https://somewhere.com"
						auto_recover_power_state = true
						costing_mode             = "costing"
						data_center_name         = "aDatacenter"
						enabled                  = true
						guidance_mode            = "off"
						import_existing_vms      = "off"
						keyboard_layout          = "us"
						location                 = "somewhere"
						security_mode            = "off"
						visibility               = "public"

						config_hvm = {
							certificate_provider          = "internal"
							enable_network_type_selection = false
						}
					}`,
				Check:    checkFn,
				PlanOnly: false,
			},
			{
				// checks plan has no effect
				Config: providerConfig + `
					resource "hpe_morpheus_cloud" "example" {
						name      = "` + name + `"
						tenant_id = 1
						group_id  = 1

						code                     = "` + code + `"
						external_id              = "` + code + `"
						labels                   = ["Label1", "Label2"]
						agent_install_mode       = "ssh"
						appliance_url            = "https://somewhere.com"
						auto_recover_power_state = true
						costing_mode             = "costing"
						data_center_name         = "aDatacenter"
						enabled                  = true
						guidance_mode            = "off"
						import_existing_vms      = "off"
						keyboard_layout          = "us"
						location                 = "somewhere"
						security_mode            = "off"
						visibility               = "public"

						config_hvm = {
							certificate_provider          = "internal"
							enable_network_type_selection = false
						}
					}`,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
			{
				// checks plan detects name change
				Config: providerConfig + `
					resource "hpe_morpheus_cloud" "example" {
						name      = "changed"
						tenant_id = 1
						group_id  = 1

						code                     = "` + code + `"
						external_id              = "` + code + `"
						labels                   = ["Label1", "Label2"]
						agent_install_mode       = "ssh"
						appliance_url            = "https://somewhere.com"
						auto_recover_power_state = true
						costing_mode             = "costing"
						data_center_name         = "aDatacenter"
						enabled                  = true
						guidance_mode            = "off"
						import_existing_vms      = "off"
						keyboard_layout          = "us"
						location                 = "somewhere"
						security_mode            = "off"
						visibility               = "public"

						config_hvm = {
							certificate_provider          = "internal"
							enable_network_type_selection = false
						}
					}`,
				ExpectNonEmptyPlan: true,
				PlanOnly:           true,
			},
			{
				// checks plan detects code change
				Config: providerConfig + `
					resource "hpe_morpheus_cloud" "example" {
						name      = "` + name + `"
						tenant_id = 1
						group_id  = 1

						code                     = "changed"
						external_id              = "` + code + `"
						labels                   = ["Label1", "Label2"]
						agent_install_mode       = "ssh"
						appliance_url            = "https://somewhere.com"
						auto_recover_power_state = true
						costing_mode             = "costing"
						data_center_name         = "aDatacenter"
						enabled                  = true
						guidance_mode            = "off"
						import_existing_vms      = "off"
						keyboard_layout          = "us"
						location                 = "somewhere"
						security_mode            = "off"
						visibility               = "public"

						config_hvm = {
							certificate_provider          = "internal"
							enable_network_type_selection = false
						}
					}`,
				ExpectNonEmptyPlan: true,
				PlanOnly:           true,
			},
			{
				// checks plan detects external_id change
				Config: providerConfig + `
					resource "hpe_morpheus_cloud" "example" {
						name      = "` + name + `"
						tenant_id = 1
						group_id  = 1

						code                     = "` + code + `"
						external_id              = "changed"
						labels                   = ["Label1", "Label2"]
						agent_install_mode       = "ssh"
						appliance_url            = "https://somewhere.com"
						auto_recover_power_state = true
						costing_mode             = "costing"
						data_center_name         = "aDatacenter"
						enabled                  = true
						guidance_mode            = "off"
						import_existing_vms      = "off"
						keyboard_layout          = "us"
						location                 = "somewhere"
						security_mode            = "off"
						visibility               = "public"

						config_hvm = {
							certificate_provider          = "internal"
							enable_network_type_selection = false
						}
					}`,
				ExpectNonEmptyPlan: true,
				PlanOnly:           true,
			},
			{
				// checks plan detects agent_install_mode change
				Config: providerConfig + `
					resource "hpe_morpheus_cloud" "example" {
						name      = "` + name + `"
						tenant_id = 1
						group_id  = 1

						code                     = "` + code + `"
						external_id              = "` + code + `"
						labels                   = ["Label1", "Label2"]
						agent_install_mode       = "cloudInit"
						appliance_url            = "https://somewhere.com"
						auto_recover_power_state = true
						costing_mode             = "costing"
						data_center_name         = "aDatacenter"
						enabled                  = true
						guidance_mode            = "off"
						import_existing_vms      = "off"
						keyboard_layout          = "us"
						location                 = "somewhere"
						security_mode            = "off"
						visibility               = "public"

						config_hvm = {
							certificate_provider          = "internal"
							enable_network_type_selection = false
						}
					}`,
				ExpectNonEmptyPlan: true,
				PlanOnly:           true,
			},
			{
				// checks plan detects appliance_url change
				Config: providerConfig + `
					resource "hpe_morpheus_cloud" "example" {
						name      = "` + name + `"
						tenant_id = 1
						group_id  = 1

						code                     = "` + code + `"
						external_id              = "` + code + `"
						labels                   = ["Label1", "Label2"]
						agent_install_mode       = "ssh"
						appliance_url            = "https://changed.com"
						auto_recover_power_state = true
						costing_mode             = "costing"
						data_center_name         = "aDatacenter"
						enabled                  = true
						guidance_mode            = "off"
						import_existing_vms      = "off"
						keyboard_layout          = "us"
						location                 = "somewhere"
						security_mode            = "off"
						visibility               = "public"

						config_hvm = {
							certificate_provider          = "internal"
							enable_network_type_selection = false
						}
					}`,
				ExpectNonEmptyPlan: true,
				PlanOnly:           true,
			},
			{
				// checks plan detects auto_recover_power_state change
				Config: providerConfig + `
					resource "hpe_morpheus_cloud" "example" {
						name      = "` + name + `"
						tenant_id = 1
						group_id  = 1

						code                     = "` + code + `"
						external_id              = "` + code + `"
						labels                   = ["Label1", "Label2"]
						agent_install_mode       = "ssh"
						appliance_url            = "https://somewhere.com"
						auto_recover_power_state = false
						costing_mode             = "costing"
						data_center_name         = "aDatacenter"
						enabled                  = true
						guidance_mode            = "off"
						import_existing_vms      = "off"
						keyboard_layout          = "us"
						location                 = "somewhere"
						security_mode            = "off"
						visibility               = "public"

						config_hvm = {
							certificate_provider          = "internal"
							enable_network_type_selection = false
						}
					}`,
				ExpectNonEmptyPlan: true,
				PlanOnly:           true,
			},
			{
				// checks plan detects costing_mode change
				Config: providerConfig + `
					resource "hpe_morpheus_cloud" "example" {
						name      = "` + name + `"
						tenant_id = 1
						group_id  = 1

						code                     = "` + code + `"
						external_id              = "` + code + `"
						labels                   = ["Label1", "Label2"]
						agent_install_mode       = "ssh"
						appliance_url            = "https://somewhere.com"
						auto_recover_power_state = true
						costing_mode             = "off"
						data_center_name         = "aDatacenter"
						enabled                  = true
						guidance_mode            = "off"
						import_existing_vms      = "off"
						keyboard_layout          = "us"
						location                 = "somewhere"
						security_mode            = "off"
						visibility               = "public"

						config_hvm = {
							certificate_provider          = "internal"
							enable_network_type_selection = false
						}
					}`,
				ExpectNonEmptyPlan: true,
				PlanOnly:           true,
			},
			{
				// checks plan detects data_center_name change
				Config: providerConfig + `
					resource "hpe_morpheus_cloud" "example" {
						name      = "` + name + `"
						tenant_id = 1
						group_id  = 1

						code                     = "` + code + `"
						external_id              = "` + code + `"
						labels                   = ["Label1", "Label2"]
						agent_install_mode       = "ssh"
						appliance_url            = "https://somewhere.com"
						auto_recover_power_state = true
						costing_mode             = "costing"
						data_center_name         = "changedDatacenter"
						enabled                  = true
						guidance_mode            = "off"
						import_existing_vms      = "off"
						keyboard_layout          = "us"
						location                 = "somewhere"
						security_mode            = "off"
						visibility               = "public"

						config_hvm = {
							certificate_provider          = "internal"
							enable_network_type_selection = false
						}
					}`,
				ExpectNonEmptyPlan: true,
				PlanOnly:           true,
			},
			{
				// checks plan detects enabled change
				Config: providerConfig + `
					resource "hpe_morpheus_cloud" "example" {
						name      = "` + name + `"
						tenant_id = 1
						group_id  = 1

						code                     = "` + code + `"
						external_id              = "` + code + `"
						labels                   = ["Label1", "Label2"]
						agent_install_mode       = "ssh"
						appliance_url            = "https://somewhere.com"
						auto_recover_power_state = true
						costing_mode             = "costing"
						data_center_name         = "aDatacenter"
						enabled                  = false
						guidance_mode            = "off"
						import_existing_vms      = "off"
						keyboard_layout          = "us"
						location                 = "somewhere"
						security_mode            = "off"
						visibility               = "public"

						config_hvm = {
							certificate_provider          = "internal"
							enable_network_type_selection = false
						}
					}`,
				ExpectNonEmptyPlan: true,
				PlanOnly:           true,
			},
			{
				// checks plan detects guidance_mode change
				Config: providerConfig + `
					resource "hpe_morpheus_cloud" "example" {
						name      = "` + name + `"
						tenant_id = 1
						group_id  = 1

						code                     = "` + code + `"
						external_id              = "` + code + `"
						labels                   = ["Label1", "Label2"]
						agent_install_mode       = "ssh"
						appliance_url            = "https://somewhere.com"
						auto_recover_power_state = true
						costing_mode             = "costing"
						data_center_name         = "aDatacenter"
						enabled                  = true
						guidance_mode            = "manual"
						import_existing_vms      = "off"
						keyboard_layout          = "us"
						location                 = "somewhere"
						security_mode            = "off"
						visibility               = "public"

						config_hvm = {
							certificate_provider          = "internal"
							enable_network_type_selection = false
						}
					}`,
				ExpectNonEmptyPlan: true,
				PlanOnly:           true,
			},
			{
				// checks plan detects import_existing_vms change
				Config: providerConfig + `
					resource "hpe_morpheus_cloud" "example" {
						name      = "` + name + `"
						tenant_id = 1
						group_id  = 1

						code                     = "` + code + `"
						external_id              = "` + code + `"
						labels                   = ["Label1", "Label2"]
						agent_install_mode       = "ssh"
						appliance_url            = "https://somewhere.com"
						auto_recover_power_state = true
						costing_mode             = "costing"
						data_center_name         = "aDatacenter"
						enabled                  = true
						guidance_mode            = "off"
						import_existing_vms      = "basic"
						keyboard_layout          = "us"
						location                 = "somewhere"
						security_mode            = "off"
						visibility               = "public"

						config_hvm = {
							certificate_provider          = "internal"
							enable_network_type_selection = false
						}
					}`,
				ExpectNonEmptyPlan: true,
				PlanOnly:           true,
			},
			{
				// checks plan detects keyboard_layout change
				Config: providerConfig + `
					resource "hpe_morpheus_cloud" "example" {
						name      = "` + name + `"
						tenant_id = 1
						group_id  = 1

						code                     = "` + code + `"
						external_id              = "` + code + `"
						labels                   = ["Label1", "Label2"]
						agent_install_mode       = "ssh"
						appliance_url            = "https://somewhere.com"
						auto_recover_power_state = true
						costing_mode             = "costing"
						data_center_name         = "aDatacenter"
						enabled                  = true
						guidance_mode            = "off"
						import_existing_vms      = "off"
						keyboard_layout          = "uk"
						location                 = "somewhere"
						security_mode            = "off"
						visibility               = "public"

						config_hvm = {
							certificate_provider          = "internal"
							enable_network_type_selection = false
						}
					}`,
				ExpectNonEmptyPlan: true,
				PlanOnly:           true,
			},
			{
				// checks plan detects labels change
				Config: providerConfig + `
					resource "hpe_morpheus_cloud" "example" {
						name      = "` + name + `"
						tenant_id = 1
						group_id  = 1

						code                     = "` + code + `"
						external_id              = "` + code + `"
						labels                   = ["Label1", "Label2", "Label3"]
						agent_install_mode       = "ssh"
						appliance_url            = "https://somewhere.com"
						auto_recover_power_state = true
						costing_mode             = "costing"
						data_center_name         = "aDatacenter"
						enabled                  = true
						guidance_mode            = "off"
						import_existing_vms      = "off"
						keyboard_layout          = "uk"
						location                 = "somewhere"
						security_mode            = "off"
						visibility               = "public"

						config_hvm = {
							certificate_provider          = "internal"
							enable_network_type_selection = false
						}
					}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					[]resource.TestCheckFunc{
						resource.TestCheckResourceAttr(
							"hpe_morpheus_cloud.example",
							"labels.#",
							"3",
						),
						resource.TestCheckTypeSetElemAttr(
							"hpe_morpheus_cloud.example",
							"labels.*",
							"Label1",
						),
						resource.TestCheckTypeSetElemAttr(
							"hpe_morpheus_cloud.example",
							"labels.*",
							"Label2",
						),
						resource.TestCheckTypeSetElemAttr(
							"hpe_morpheus_cloud.example",
							"labels.*",
							"Label3",
						),
					}...),
				ExpectNonEmptyPlan: true,
				PlanOnly:           true,
			},
			{
				// checks plan detects location change
				Config: providerConfig + `
					resource "hpe_morpheus_cloud" "example" {
						name      = "` + name + `"
						tenant_id = 1
						group_id  = 1

						code                     = "` + code + `"
						external_id              = "` + code + `"
						labels                   = ["Label1", "Label2"]
						agent_install_mode       = "ssh"
						appliance_url            = "https://somewhere.com"
						auto_recover_power_state = true
						costing_mode             = "costing"
						data_center_name         = "aDatacenter"
						enabled                  = true
						guidance_mode            = "off"
						import_existing_vms      = "off"
						keyboard_layout          = "us"
						location                 = "changed"
						security_mode            = "off"
						visibility               = "public"

						config_hvm = {
							certificate_provider          = "internal"
							enable_network_type_selection = false
						}
					}`,
				ExpectNonEmptyPlan: true,
				PlanOnly:           true,
			},
			{
				// checks plan detects security_mode change
				Config: providerConfig + `
					resource "hpe_morpheus_cloud" "example" {
						name      = "` + name + `"
						tenant_id = 1
						group_id  = 1

						code                     = "` + code + `"
						external_id              = "` + code + `"
						labels                   = ["Label1", "Label2"]
						agent_install_mode       = "ssh"
						appliance_url            = "https://somewhere.com"
						auto_recover_power_state = true
						costing_mode             = "costing"
						data_center_name         = "aDatacenter"
						enabled                  = true
						guidance_mode            = "off"
						import_existing_vms      = "off"
						keyboard_layout          = "us"
						location                 = "somewhere"
						security_mode            = "internal"
						visibility               = "public"

						config_hvm = {
							certificate_provider          = "internal"
							enable_network_type_selection = false
						}
					}`,
				ExpectNonEmptyPlan: true,
				PlanOnly:           true,
			},
			{
				// checks plan detects visibility change
				Config: providerConfig + `
					resource "hpe_morpheus_cloud" "example" {
						name      = "` + name + `"
						tenant_id = 1
						group_id  = 1

						code                     = "` + code + `"
						external_id              = "` + code + `"
						labels                   = ["Label1", "Label2"]
						agent_install_mode       = "ssh"
						appliance_url            = "https://somewhere.com"
						auto_recover_power_state = true
						costing_mode             = "costing"
						data_center_name         = "aDatacenter"
						enabled                  = true
						guidance_mode            = "off"
						import_existing_vms      = "off"
						keyboard_layout          = "us"
						location                 = "somewhere"
						security_mode            = "off"
						visibility               = "private"

						config_hvm = {
							certificate_provider          = "internal"
							enable_network_type_selection = false
						}
					}`,
				ExpectNonEmptyPlan: true,
				PlanOnly:           true,
			},
			{
				// checks plan detects certificate_provider change
				Config: providerConfig + `
					resource "hpe_morpheus_cloud" "example" {
						name      = "` + name + `"
						tenant_id = 1
						group_id  = 1

						code                     = "` + code + `"
						external_id              = "` + code + `"
						labels                   = ["Label1", "Label2"]
						agent_install_mode       = "ssh"
						appliance_url            = "https://somewhere.com"
						auto_recover_power_state = true
						costing_mode             = "costing"
						data_center_name         = "aDatacenter"
						enabled                  = true
						guidance_mode            = "off"
						import_existing_vms      = "off"
						keyboard_layout          = "us"
						location                 = "somewhere"
						security_mode            = "off"
						visibility               = "public"

						config_hvm = {
							certificate_provider          = "manual"
							enable_network_type_selection = false
						}
					}`,
				ExpectNonEmptyPlan: true,
				PlanOnly:           true,
			},
			{
				// checks plan detects enable_network_type_selection change
				Config: providerConfig + `
					resource "hpe_morpheus_cloud" "example" {
						name      = "` + name + `"
						tenant_id = 1
						group_id  = 1

						code                     = "` + code + `"
						external_id              = "` + code + `"
						labels                   = ["Label1", "Label2"]
						agent_install_mode       = "ssh"
						appliance_url            = "https://somewhere.com"
						auto_recover_power_state = true
						costing_mode             = "costing"
						data_center_name         = "aDatacenter"
						enabled                  = true
						guidance_mode            = "off"
						import_existing_vms      = "off"
						keyboard_layout          = "us"
						location                 = "somewhere"
						security_mode            = "off"
						visibility               = "public"

						config_hvm = {
							certificate_provider          = "internal"
							enable_network_type_selection = true
						}
					}`,
				ExpectNonEmptyPlan: true,
				PlanOnly:           true,
			},
		},
	})
}

func TestAccMorpheusCloudValidationOneOf(t *testing.T) {
	defer testhelpers.RecordResult(t)

	t.Parallel()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// checks plan fails when agent_install_mode has invalid value
				Config: `
					resource "hpe_morpheus_cloud" "example" {
						name      = "TestCloud"
						tenant_id = 1
						group_id  = 1

						agent_install_mode       = "invalid"
					}`,
				ExpectError: regexp.MustCompile(`Attribute agent_install_mode value must be one of:`),
			},
			{
				// checks plan fails when costing_mode has invalid value
				Config: `
					resource "hpe_morpheus_cloud" "example" {
						name      = "TestCloud"
						tenant_id = 1
						group_id  = 1

						costing_mode             = "invalid"
					}`,
				ExpectError: regexp.MustCompile(`Attribute costing_mode value must be one of:`),
			},
			{
				// checks plan fails when guidance_mode has invalid value
				Config: `
					resource "hpe_morpheus_cloud" "example" {
						name      = "TestCloud"
						tenant_id = 1
						group_id  = 1

						guidance_mode            = "invalid"
					}`,
				ExpectError: regexp.MustCompile(`Attribute guidance_mode value must be one of:`),
			},
			{
				// checks plan fails when import_existing_vms has invalid value
				Config: `
					resource "hpe_morpheus_cloud" "example" {
						name      = "TestCloud"
						tenant_id = 1
						group_id  = 1

						import_existing_vms      = "invalid"
					}`,
				ExpectError: regexp.MustCompile(`Attribute import_existing_vms value must be one of:`),
			},
			{
				// checks plan fails when security_mode has invalid value
				Config: `
					resource "hpe_morpheus_cloud" "example" {
						name      = "TestCloud"
						tenant_id = 1
						group_id  = 1

						security_mode            = "invalid"
					}`,
				ExpectError: regexp.MustCompile(`Attribute security_mode value must be one of:`),
			},
			{
				// checks plan fails when visibility has invalid value
				Config: `
					resource "hpe_morpheus_cloud" "example" {
						name      = "TestCloud"
						tenant_id = 1
						group_id  = 1

						visibility               = "invalid"
					}`,
				ExpectError: regexp.MustCompile(`Attribute visibility value must be one of:`),
			},
		},
	})
}

func TestAccMorpheusCloudValidationRequiredAttrs(t *testing.T) {
	defer testhelpers.RecordResult(t)

	t.Parallel()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// checks plan fails when TenantId is removed
				Config: `
					resource "hpe_morpheus_cloud" "example" {
						name      = "cloud9"
						group_id  = 1
					}`,
				ExpectError: regexp.MustCompile(`The argument "tenant_id" is required`),
			},
			{
				// checks plan fails when GroupId is removed
				Config: `
					resource "hpe_morpheus_cloud" "example" {
						name      = "cloud9"
						tenant_id  = 1
					}`,
				ExpectError: regexp.MustCompile(`The argument "group_id" is required`),
			},
			{
				// checks plan fails when Name is removed
				Config: `
					resource "hpe_morpheus_cloud" "example" {
						tenant_id  = 1
				    group_id = 1
					}`,
				ExpectError: regexp.MustCompile(`The argument "name" is required`),
			},
		},
	})
}
