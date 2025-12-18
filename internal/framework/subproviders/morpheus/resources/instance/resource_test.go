// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package instance_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/internal/framework/provider"
	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus"
	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
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
func TestAccMorpheusInstanceExampleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())
	instanceTypeID := "9"
	resourcePool := "pool-62299"

	resourceConfig, err := testhelpers.RenderExample(t, "example.tf.tmpl",
		"Name", name,
		"InstanceType", instanceTypeID,
		"ResourcePool", resourcePool,
	)
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_instance.example",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_instance.example",
			"instance_type_id",
			instanceTypeID,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_instance.example",
			"config.resourcePoolId",
			resourcePool,
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           false,
				Check:              checkFn,
			},
		},
	})
}

// func TestAccMorpheusInstanceUpdate(t *testing.T) {
// 	defer testhelpers.RecordResult(t)
//
// 	if testing.Short() {
// 		t.Skip("Skipping slow test in short mode")
// 	}
//
// 	t.Parallel()
//
// 	providerConfig := testhelpers.ProviderBlock()
// 	name := acctest.RandomWithPrefix(t.Name())
//
// 	checks := []resource.TestCheckFunc{
// 		resource.TestCheckResourceAttr(
// 			"hpe_morpheus_instance.example",
// 			"name",
// 			name,
// 		),
// 		resource.TestCheckResourceAttr(
// 			"hpe_morpheus_instance.example",
// 			"instance_type_id",
// 			"9",
// 		),
// 		resource.TestCheckResourceAttr(
// 			"hpe_morpheus_instance.example",
// 			"layout_id",
// 			"5385",
// 		),
// 		resource.TestCheckResourceAttr(
// 			"hpe_morpheus_instance.example",
// 			"instance_context",
// 			"dev",
// 		),
// 		resource.TestCheckResourceAttr(
// 			"hpe_morpheus_instance.example",
// 			"layout_size",
// 			"1",
// 		),
// 	}
//
// 	checkFn := resource.ComposeAggregateTestCheckFunc(
// 		checks...,
// 	)
//
// 	resource.Test(t, resource.TestCase{
// 		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: providerConfig + `
// 					data "hpe_morpheus_cloud" "vme_cloud" {
// 						name = "HPE Alletra VME"
// 					}
//
// 					data "hpe_morpheus_service_plan" "vme_512mb" {
// 						name                = "1 CPU, 1GB Memory"
// 						provision_type_code = "kvm"
// 					}
//
// 					resource "hpe_morpheus_instance" "example" {
// 						name               = "` + name + `"
// 						cloud_id           = data.hpe_morpheus_cloud.vme_cloud.id
// 						layout_id          = 5385
// 						instance_type_id   = 9
// 						layout_size        = 1
//
// 						group_id           = 1
// 						plan_id            = data.hpe_morpheus_service_plan.vme_512mb.id
//
// 						instance_context   = "dev"
// 						network_interfaces = [
// 							{
// 								network_id = 103481
// 								ip_mode    = "dhcp"
// 							}
// 						]
//
// 						volumes = [
// 							{
// 								root_volume  = true
// 								name         = "root"
// 								size         = 10
// 								storage_type = 1
// 								datastore_id = 38658
// 							},
// 							{
// 								root_volume  = false
// 								name         = "data"
// 								size         = 10
// 								storage_type = 1
// 								datastore_id = 38658
// 							}
// 						]
//
// 						tags = [
// 							{
// 								name  = "managed_by"
// 								value = "terraform"
// 							}
// 						]
//
// 						config = {
// 							resourcePoolId       = "pool-62299"
// 							poolProviderType     = "mvm"
// 							nestedVirtualization = "off"
// 							noAgent              = true
// 							createUser           = false
// 						}
// 					}`,
// 				Check:    checkFn,
// 				PlanOnly: false,
// 			},
// 			{
// 				// checks plan has no effect
// 				Config: providerConfig + `
// 					data "hpe_morpheus_cloud" "vme_cloud" {
// 						name = "HPE Alletra VME"
// 					}
//
// 					data "hpe_morpheus_service_plan" "vme_512mb" {
// 						name                = "1 CPU, 1GB Memory"
// 						provision_type_code = "kvm"
// 					}
//
// 					resource "hpe_morpheus_instance" "example" {
// 						name               = "` + name + `"
// 						cloud_id           = data.hpe_morpheus_cloud.vme_cloud.id
// 						layout_id          = 5385
// 						instance_type_id   = 9
// 						layout_size        = 1
//
// 						group_id           = 1
// 						plan_id            = data.hpe_morpheus_service_plan.vme_512mb.id
//
// 						instance_context   = "dev"
// 						network_interfaces = [
// 							{
// 								network_id = 103481
// 								ip_mode    = "dhcp"
// 							}
// 						]
//
// 						volumes = [
// 							{
// 								root_volume  = true
// 								name         = "root"
// 								size         = 10
// 								storage_type = 1
// 								datastore_id = 38658
// 							}
// 						]
//
// 						tags = [
// 							{
// 								name  = "managed_by"
// 								value = "terraform"
// 							}
// 						]
//
// 						config = {
// 							resourcePoolId       = "pool-62299"
// 							poolProviderType     = "mvm"
// 							nestedVirtualization = "off"
// 							noAgent              = true
// 							createUser           = false
// 						}
// 					}`,
// 				ExpectNonEmptyPlan: false,
// 				PlanOnly:           true,
// 			},
// 			{
// 				// checks plan detects name change
// 				Config: providerConfig + `
// 					data "hpe_morpheus_cloud" "vme_cloud" {
// 						name = "HPE Alletra VME"
// 					}
//
// 					data "hpe_morpheus_service_plan" "vme_512mb" {
// 						name                = "1 CPU, 1GB Memory"
// 						provision_type_code = "kvm"
// 					}
//
// 					resource "hpe_morpheus_instance" "example" {
// 						name               = "changed" # changed
// 						cloud_id           = data.hpe_morpheus_cloud.vme_cloud.id
// 						layout_id          = 5385
// 						instance_type_id   = 9
// 						layout_size        = 1
//
// 						group_id           = 1
// 						plan_id            = data.hpe_morpheus_service_plan.vme_512mb.id
//
// 						instance_context   = "dev"
// 						network_interfaces = [
// 							{
// 								network_id = 103481
// 								ip_mode    = "dhcp"
// 							}
// 						]
//
// 						volumes = [
// 							{
// 								root_volume  = true
// 								name         = "root"
// 								size         = 10
// 								storage_type = 1
// 								datastore_id = 38658
// 							}
// 						]
//
// 						tags = [
// 							{
// 								name  = "managed_by"
// 								value = "terraform"
// 							}
// 						]
//
// 						config = {
// 							resourcePoolId       = "pool-62299"
// 							poolProviderType     = "mvm"
// 							nestedVirtualization = "off"
// 							noAgent              = true
// 							createUser           = false
// 						}
// 					}`,
// 				ExpectNonEmptyPlan: true,
// 				PlanOnly:           true,
// 			},
// 			{
// 				// checks plan detects instance_context change
// 				Config: providerConfig + `
// 					data "hpe_morpheus_cloud" "vme_cloud" {
// 						name = "HPE Alletra VME"
// 					}
//
// 					data "hpe_morpheus_service_plan" "vme_512mb" {
// 						name                = "1 CPU, 1GB Memory"
// 						provision_type_code = "kvm"
// 					}
//
// 					resource "hpe_morpheus_instance" "example" {
// 						name               = "` + name + `"
// 						cloud_id           = data.hpe_morpheus_cloud.vme_cloud.id
// 						layout_id          = 5385
// 						instance_type_id   = 9
// 						layout_size        = 1
//
// 						group_id           = 1
// 						plan_id            = data.hpe_morpheus_service_plan.vme_512mb.id
//
// 						instance_context   = "prod" # changed
// 						network_interfaces = [
// 							{
// 								network_id = 103481
// 								ip_mode    = "dhcp"
// 							}
// 						]
//
// 						volumes = [
// 							{
// 								root_volume  = true
// 								name         = "root"
// 								size         = 10
// 								storage_type = 1
// 								datastore_id = 38658
// 							}
// 						]
//
// 						tags = [
// 							{
// 								name  = "managed_by"
// 								value = "terraform"
// 							}
// 						]
//
// 						config = {
// 							resourcePoolId       = "pool-62299"
// 							poolProviderType     = "mvm"
// 							nestedVirtualization = "off"
// 							noAgent              = true
// 							createUser           = false
// 						}
// 					}`,
// 				ExpectNonEmptyPlan: true,
// 				PlanOnly:           true,
// 			},
// 			{
// 				// checks plan detects layout_size change
// 				Config: providerConfig + `
// 					data "hpe_morpheus_cloud" "vme_cloud" {
// 						name = "HPE Alletra VME"
// 					}
//
// 					data "hpe_morpheus_service_plan" "vme_512mb" {
// 						name                = "1 CPU, 1GB Memory"
// 						provision_type_code = "kvm"
// 					}
//
// 					resource "hpe_morpheus_instance" "example" {
// 						name               = "` + name + `"
// 						cloud_id           = data.hpe_morpheus_cloud.vme_cloud.id
// 						layout_id          = 5385
// 						instance_type_id   = 9
// 						layout_size        = 2 # changed
//
// 						group_id           = 1
// 						plan_id            = data.hpe_morpheus_service_plan.vme_512mb.id
//
// 						instance_context   = "dev"
// 						network_interfaces = [
// 							{
// 								network_id = 103481
// 								ip_mode    = "dhcp"
// 							}
// 						]
//
// 						volumes = [
// 							{
// 								root_volume  = true
// 								name         = "root"
// 								size         = 10
// 								storage_type = 1
// 								datastore_id = 38658
// 							}
// 						]
//
// 						tags = [
// 							{
// 								name  = "managed_by"
// 								value = "terraform"
// 							}
// 						]
//
// 						config = {
// 							resourcePoolId       = "pool-62299"
// 							poolProviderType     = "mvm"
// 							nestedVirtualization = "off"
// 							noAgent              = true
// 							createUser           = false
// 						}
// 					}`,
// 				ExpectNonEmptyPlan: true,
// 				PlanOnly:           true,
// 			},
// 			{
// 				// checks plan detects config change
// 				Config: providerConfig + `
// 					data "hpe_morpheus_cloud" "vme_cloud" {
// 						name = "HPE Alletra VME"
// 					}
//
// 					data "hpe_morpheus_service_plan" "vme_512mb" {
// 						name                = "1 CPU, 1GB Memory"
// 						provision_type_code = "kvm"
// 					}
//
// 					resource "hpe_morpheus_instance" "example" {
// 						name               = "` + name + `"
// 						cloud_id           = data.hpe_morpheus_cloud.vme_cloud.id
// 						layout_id          = 5385
// 						instance_type_id   = 9
// 						layout_size        = 1
//
// 						group_id           = 1
// 						plan_id            = data.hpe_morpheus_service_plan.vme_512mb.id
//
// 						instance_context   = "dev"
// 						network_interfaces = [
// 							{
// 								network_id = 103481
// 								ip_mode    = "dhcp"
// 							}
// 						]
//
// 						volumes = [
// 							{
// 								root_volume  = true
// 								name         = "root"
// 								size         = 10
// 								storage_type = 1
// 								datastore_id = 38658
// 							}
// 						]
//
// 						tags = [
// 							{
// 								name  = "managed_by"
// 								value = "terraform"
// 							}
// 						]
//
// 						config = {
// 							resourcePoolId       = "pool-62299"
// 							poolProviderType     = "mvm"
// 							nestedVirtualization = "on" # changed
// 							noAgent              = true
// 							createUser           = false
// 						}
// 					}`,
// 				ExpectNonEmptyPlan: true,
// 				PlanOnly:           true,
// 			},
// 		},
// 	})
// }
