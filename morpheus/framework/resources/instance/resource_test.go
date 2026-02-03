// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

//go:generate ../../../../../../bin/render example.tf.tmpl Name "TestInstance" InstanceType "9" ResourcePool "pool-62299"
//go:generate ../../../../../../bin/render example_twonetworks.tf.tmpl Name "TestInstance" InstanceType "9" ResourcePool "pool-62299"
//go:generate ../../../../../../bin/render example_timeouts.tf.tmpl Name "TestInstance" InstanceType "9" ResourcePool "pool-62299"
//go:generate ../../../../../../bin/render example_vmware.tf.tmpl Name "TestInstance" InstanceType "9" ResourcePool "pool-1"
//go:generate ../../../../../../bin/render example_metal.tf.tmpl Name "TestInstance" CloudName "aCloud" EnvironmentName "anEnvironment" GroupName "aGroup" InstanceTypeLayout "Single ILO Server" Role "aRole" PlanName "G3i"

package instance_test

import (
	"os"
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

// Tests that our example file template used for docs is a valid config
func TestAccMorpheusInstanceExampleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

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

func TestAccMorpheusInstanceUpdateName(t *testing.T) {
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())
	updatedName := name + "-updated"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
					data "hpe_morpheus_cloud" "vme_cloud" {
						name = "HPE Alletra VME"
					}

					data "hpe_morpheus_service_plan" "vme_512mb" {
						name                = "1 CPU, 1GB Memory"
						provision_type_code = "kvm"
					}

					resource "hpe_morpheus_instance" "example" {
						name               = "` + name + `"
						cloud_id           = data.hpe_morpheus_cloud.vme_cloud.id
						layout_id          = 5385
						instance_type_id   = 9

						group_id           = 1
						plan_id            = data.hpe_morpheus_service_plan.vme_512mb.id

						instance_context   = "dev"
						network_interfaces = [
							{
								network_id = 103481
							}
						]

						volumes = [
							{
								root_volume     = true
								name            = "root"
								size            = 10
								storage_type_id = 1
								datastore_id    = 38658
							}
						]

						tags = [
							{
								name  = "managed_by"
								value = "terraform"
							}
						]

						config = {
							resourcePoolId       = "pool-62299"
							poolProviderType     = "mvm"
							nestedVirtualization = "off"
							noAgent              = true
							createUser           = false
						}
					}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"hpe_morpheus_instance.example",
						"name",
						name,
					),
				),
			},
			{
				Config: providerConfig + `
					data "hpe_morpheus_cloud" "vme_cloud" {
						name = "HPE Alletra VME"
					}

					data "hpe_morpheus_service_plan" "vme_512mb" {
						name                = "1 CPU, 1GB Memory"
						provision_type_code = "kvm"
					}

					resource "hpe_morpheus_instance" "example" {
						name               = "` + updatedName + `"
						cloud_id           = data.hpe_morpheus_cloud.vme_cloud.id
						layout_id          = 5385
						instance_type_id   = 9

						group_id           = 1
						plan_id            = data.hpe_morpheus_service_plan.vme_512mb.id

						instance_context   = "dev"
						network_interfaces = [
							{
								network_id = 103481
							}
						]

						volumes = [
							{
								root_volume     = true
								name            = "root"
								size            = 10
								storage_type_id = 1
								datastore_id    = 38658
							}
						]

						tags = [
							{
								name  = "managed_by"
								value = "terraform"
							}
						]

						config = {
							resourcePoolId       = "pool-62299"
							poolProviderType     = "mvm"
							nestedVirtualization = "off"
							noAgent              = true
							createUser           = false
						}
					}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"hpe_morpheus_instance.example",
						"name",
						updatedName,
					),
				),
			},
		},
	})
}

func TestAccMorpheusInstanceUpdateInstanceContext(t *testing.T) {
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
					data "hpe_morpheus_cloud" "vme_cloud" {
						name = "HPE Alletra VME"
					}

					data "hpe_morpheus_service_plan" "vme_512mb" {
						name                = "1 CPU, 1GB Memory"
						provision_type_code = "kvm"
					}

					resource "hpe_morpheus_instance" "example" {
						name               = "` + name + `"
						cloud_id           = data.hpe_morpheus_cloud.vme_cloud.id
						layout_id          = 5385
						instance_type_id   = 9

						group_id           = 1
						plan_id            = data.hpe_morpheus_service_plan.vme_512mb.id

						instance_context   = "dev"
						network_interfaces = [
							{
								network_id = 103481
							}
						]

						volumes = [
							{
								root_volume     = true
								name            = "root"
								size            = 10
								storage_type_id = 1
								datastore_id    = 38658
							}
						]

						tags = [
							{
								name  = "managed_by"
								value = "terraform"
							}
						]

						config = {
							resourcePoolId       = "pool-62299"
							poolProviderType     = "mvm"
							nestedVirtualization = "off"
							noAgent              = true
							createUser           = false
						}
					}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"hpe_morpheus_instance.example",
						"instance_context",
						"dev",
					),
				),
			},
			{
				Config: providerConfig + `
					data "hpe_morpheus_cloud" "vme_cloud" {
						name = "HPE Alletra VME"
					}

					data "hpe_morpheus_service_plan" "vme_512mb" {
						name                = "1 CPU, 1GB Memory"
						provision_type_code = "kvm"
					}

					resource "hpe_morpheus_instance" "example" {
						name               = "` + name + `"
						cloud_id           = data.hpe_morpheus_cloud.vme_cloud.id
						layout_id          = 5385
						instance_type_id   = 9

						group_id           = 1
						plan_id            = data.hpe_morpheus_service_plan.vme_512mb.id

						instance_context   = "production"
						network_interfaces = [
							{
								network_id = 103481
							}
						]

						volumes = [
							{
								root_volume     = true
								name            = "root"
								size            = 10
								storage_type_id = 1
								datastore_id    = 38658
							}
						]

						tags = [
							{
								name  = "managed_by"
								value = "terraform"
							}
						]

						config = {
							resourcePoolId       = "pool-62299"
							poolProviderType     = "mvm"
							nestedVirtualization = "off"
							noAgent              = true
							createUser           = false
						}
					}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"hpe_morpheus_instance.example",
						"instance_context",
						"production",
					),
				),
			},
		},
	})
}

func TestAccMorpheusInstanceUpdateTags(t *testing.T) {
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
					data "hpe_morpheus_cloud" "vme_cloud" {
						name = "HPE Alletra VME"
					}

					data "hpe_morpheus_service_plan" "vme_512mb" {
						name                = "1 CPU, 1GB Memory"
						provision_type_code = "kvm"
					}

					resource "hpe_morpheus_instance" "example" {
						name               = "` + name + `"
						cloud_id           = data.hpe_morpheus_cloud.vme_cloud.id
						layout_id          = 5385
						instance_type_id   = 9

						group_id           = 1
						plan_id            = data.hpe_morpheus_service_plan.vme_512mb.id

						instance_context   = "dev"
						network_interfaces = [
							{
								network_id = 103481
							}
						]

						volumes = [
							{
								root_volume     = true
								name            = "root"
								size            = 10
								storage_type_id = 1
								datastore_id    = 38658
							}
						]

						tags = [
							{
								name  = "managed_by"
								value = "terraform"
							}
						]

						config = {
							resourcePoolId       = "pool-62299"
							poolProviderType     = "mvm"
							nestedVirtualization = "off"
							noAgent              = true
							createUser           = false
						}
					}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"hpe_morpheus_instance.example",
						"tags.#",
						"1",
					),
				),
			},
			{
				Config: providerConfig + `
					data "hpe_morpheus_cloud" "vme_cloud" {
						name = "HPE Alletra VME"
					}

					data "hpe_morpheus_service_plan" "vme_512mb" {
						name                = "1 CPU, 1GB Memory"
						provision_type_code = "kvm"
					}

					resource "hpe_morpheus_instance" "example" {
						name               = "` + name + `"
						cloud_id           = data.hpe_morpheus_cloud.vme_cloud.id
						layout_id          = 5385
						instance_type_id   = 9

						group_id           = 1
						plan_id            = data.hpe_morpheus_service_plan.vme_512mb.id

						instance_context   = "dev"
						network_interfaces = [
							{
								network_id = 103481
							}
						]

						volumes = [
							{
								root_volume     = true
								name            = "root"
								size            = 10
								storage_type_id = 1
								datastore_id    = 38658
							}
						]

						tags = [
							{
								name  = "managed_by"
								value = "terraform"
							},
							{
								name  = "environment"
								value = "test"
							}
						]

						config = {
							resourcePoolId       = "pool-62299"
							poolProviderType     = "mvm"
							nestedVirtualization = "off"
							noAgent              = true
							createUser           = false
						}
					}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"hpe_morpheus_instance.example",
						"tags.#",
						"2",
					),
				),
			},
		},
	})
}
