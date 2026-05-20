// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

//go:generate ../../../../bin/render example_alletramp_hvm.tf.tmpl Name "TestAlletraDatastore" AssociatedResourceID 1 StorageServerID 1 GroupID 1 TenantID 1
//go:generate ../../../../bin/render example_alletramp_bm.tf.tmpl Name "TestAlletraDatastore" CloudName "Metal" AssociatedResourceID 1 StorageServerID 1 GroupID 1 TenantID 1

package datastore_test

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
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/systemoverride"
	"github.com/HPE/terraform-provider-hpe/provider"
)

func TestMain(m *testing.M) {
	systemoverride.ParseFlags()
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
func TestAccMorpheusDatastoreExampleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	t.Skip("Skipping all Feature data-store tests")
	testSystem := systemoverride.GetPreferred(t, "feature")
	providerConfig := testhelpers.ProviderBlockForServer(testSystem)

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := testhelpers.RenderExample(
		t, "example_alletramp_hvm.tf.tmpl",
		"Name", name,
		"ProviderConfig", providerConfig,
		"AssociatedResourceID", "6032",
		"StorageServerID", "1489",
		"GroupID", "17830",
		"TenantID", "466",
	)
	if err != nil {
		t.Fatal(err)
	}

	checksNotNested := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_datastore.example",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_datastore.example",
			"associated_resource_type",
			"Cluster",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_datastore.example",
			"associated_resource_id",
			"6032",
		),
	}
	checksNested := []resource.TestCheckFunc{
		resource.TestCheckTypeSetElemNestedAttrs(
			"hpe_morpheus_datastore.example",
			"resource_permissions.groups.*",
			map[string]string{
				"id": "17830",
			},
		),
		resource.TestCheckTypeSetElemNestedAttrs(
			"hpe_morpheus_datastore.example",
			"tenants.*",
			map[string]string{
				"id": "466",
			},
		),
	}

	// When we import we'll see the root tenant added to the list of tenants
	importCheckRootTenant := resource.TestCheckTypeSetElemNestedAttrs(
		"hpe_morpheus_datastore.example",
		"tenants.*",
		map[string]string{
			"id": "1",
		},
	)

	checksCombined := append(checksNotNested, checksNested...)
	checksImport := append(checksCombined, importCheckRootTenant)

	checkFnCombined := resource.ComposeAggregateTestCheckFunc(checksCombined...)
	checkFnImport := resource.ComposeAggregateTestCheckFunc(checksImport...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				Check:              checkFnCombined,
				PlanOnly:           false,
			},
			{
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				Check:              checkFnCombined,
				PlanOnly:           true,
			},
			{
				ImportState: true,
				// Don't verify state after import as tenants list will have an extra entry for root tenant
				ImportStateVerify: false,
				ResourceName:      "hpe_morpheus_datastore.example",
				// Check that the root tenant is added upon import
				Check: checkFnImport,
			},
		},
	})
}

// Tests that our example file template used for docs is a valid config
func TestAccMorpheusDatastoreUpdateOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	t.Skip("Skipping all Feature data-store tests")
	testSystem := systemoverride.GetPreferred(t, "feature")
	providerConfig := testhelpers.ProviderBlockForServer(testSystem)
	name := acctest.RandomWithPrefix(t.Name())

	checksNotNested := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_datastore.example",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_datastore.example",
			"associated_resource_type",
			"Cluster",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_datastore.example",
			"associated_resource_id",
			"6032",
		),
	}
	checkNestedResoucePermissions := []resource.TestCheckFunc{
		resource.TestCheckTypeSetElemNestedAttrs(
			"hpe_morpheus_datastore.example",
			"resource_permissions.groups.*",
			map[string]string{
				"id": "17830",
			},
		),
	}
	checkNestedTenants := []resource.TestCheckFunc{
		resource.TestCheckTypeSetElemNestedAttrs(
			"hpe_morpheus_datastore.example",
			"tenants.*",
			map[string]string{
				"id": "466",
			},
		),
	}

	checksNotNestedAndResourcePermissions := append(checksNotNested, checkNestedResoucePermissions...)
	checksAll := append(checksNotNestedAndResourcePermissions, checkNestedTenants...)

	checkFnAll := resource.ComposeAggregateTestCheckFunc(checksAll...)
	checkFnNotNested := resource.ComposeAggregateTestCheckFunc(checksNotNested...)
	checkFnNotNestedAndResourcePermissions := resource.ComposeAggregateTestCheckFunc(
		checksNotNestedAndResourcePermissions...,
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Create starting config, no groups or tenants
				Config: providerConfig + `
                    resource "hpe_morpheus_datastore" "example" {
                      name = ` + "\"" + name + "\"" + `
                      datastore_type = {
                        id   = 8
                        code = "hpedatastore-alletra-mp"
                      }
                      associated_resource_type = "Cluster"
                      visibility               = "private"
                      active                   = true
                      associated_resource_id   = 6032
                    
                      config_alletramp_hvm = {
                        protocol_type     = "iSCSI"
                        enable_ransomware = false
                      }

                      storage_server = {
                        id = 1489
                      }
					}`,
				Check:    checkFnNotNested,
				PlanOnly: false,
			},
			{
				// Checks that plan results in no changes
				Config: providerConfig + `
                    resource "hpe_morpheus_datastore" "example" {
                      name = ` + "\"" + name + "\"" + `
                      datastore_type = {
                        id   = 8
                        code = "hpedatastore-alletra-mp"
                      }
                      associated_resource_type = "Cluster"
                      visibility               = "private"
                      active                   = true
                      associated_resource_id   = 6032
                    
                      config_alletramp_hvm = {
                        protocol_type     = "iSCSI"
                        enable_ransomware = false
                      }

                      storage_server = {
                        id = 1489
                      }
					}`,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
			{
				// Checks that plan detects name change
				Config: providerConfig + `
                    resource "hpe_morpheus_datastore" "example" {
                      name = "changed"
                      datastore_type = {
                        id   = 8
                        code = "hpedatastore-alletra-mp"
                      }
                      associated_resource_type = "Cluster"
                      visibility               = "private"
                      active                   = true
                      associated_resource_id   = 6032
                    
                      config_alletramp_hvm = {
                        protocol_type     = "iSCSI"
                        enable_ransomware = false
                      }

                      storage_server = {
                        id = 1489
                      }
					}`,
				ExpectNonEmptyPlan: true,
				PlanOnly:           true,
			},
			{
				// Checks that plan detects visibility change
				Config: providerConfig + `
                    resource "hpe_morpheus_datastore" "example" {
                      name = ` + "\"" + name + "\"" + `
                      datastore_type = {
                        id   = 8
                        code = "hpedatastore-alletra-mp"
                      }
                      associated_resource_type = "Cluster"
                      visibility               = "public"
                      active                   = true
                      associated_resource_id   = 6032
                    
                      config_alletramp_hvm = {
                        protocol_type     = "iSCSI"
                        enable_ransomware = false
                      }

                      storage_server = {
                        id = 1489
                      }
					}`,
				ExpectNonEmptyPlan: true,
				PlanOnly:           true,
			},
			{
				// Checks that plan detects associated_resource_type change
				Config: providerConfig + `
                    resource "hpe_morpheus_datastore" "example" {
                      name = ` + "\"" + name + "\"" + `
                      datastore_type = {
                        id   = 8
                        code = "hpedatastore-alletra-mp"
                      }
                      associated_resource_type = "Cloud"
                      visibility               = "private"
                      active                   = true
                      associated_resource_id   = 6032
                    
                      config_alletramp_hvm = {
                        protocol_type     = "iSCSI"
                        enable_ransomware = false
                      }

                      storage_server = {
                        id = 1489
                      }
					}`,
				ExpectNonEmptyPlan: true,
				PlanOnly:           true,
			},
			{
				// Checks that plan detects associated_resource_id change
				Config: providerConfig + `
                    resource "hpe_morpheus_datastore" "example" {
                      name = ` + "\"" + name + "\"" + `
                      datastore_type = {
                        id   = 8
                        code = "hpedatastore-alletra-mp"
                      }
                      associated_resource_type = "Cluster"
                      visibility               = "private"
                      active                   = true
                      associated_resource_id   = 1
                    
                      config_alletramp_hvm = {
                        protocol_type     = "iSCSI"
                        enable_ransomware = false
                      }

                      storage_server = {
                        id = 1489
                      }
					}`,
				ExpectNonEmptyPlan: true,
				PlanOnly:           true,
			},
			{
				// Check that plan detects datastore_type and config block changes
				Config: providerConfig + `
                    resource "hpe_morpheus_datastore" "example" {
                      name = ` + "\"" + name + "\"" + `

                      datastore_type = {
                        id   = 3
                        code = "libvirt-netfs-nfs"
                      }

                      associated_resource_type = "Cluster"
                      visibility               = "private"
                      active                   = true
                      associated_resource_id   = 6032

                      config_nfs = {
                        source_dir_path = "/tmp/dir/dir"
                        source_hostname = "nfs.example.com"
                        source_version = "3"
                      }
					}`,
				ExpectNonEmptyPlan: true,
				PlanOnly:           true,
			},
			{
				// Update with resource_permissions
				Config: providerConfig + `
                    resource "hpe_morpheus_datastore" "example" {
                      name = ` + "\"" + name + "\"" + `
                      datastore_type = {
                        id   = 8
                        code = "hpedatastore-alletra-mp"
                      }
                      associated_resource_type = "Cluster"
                      visibility               = "private"
                      active                   = true
                      associated_resource_id   = 6032

                      resource_permissions = {
                        groups = [
                          {
                            id = 17830
                          }
                        ]
                      }
                    
                      config_alletramp_hvm = {
                        protocol_type     = "iSCSI"
                        enable_ransomware = false
                      }

                      storage_server = {
                        id = 1489
                      }
					}`,
				Check:    checkFnNotNestedAndResourcePermissions,
				PlanOnly: false,
			},
			{
				// Check change in resource permissions detected
				Config: providerConfig + `
                    resource "hpe_morpheus_datastore" "example" {
                      name = ` + "\"" + name + "\"" + `
                      datastore_type = {
                        id   = 8
                        code = "hpedatastore-alletra-mp"
                      }
                      associated_resource_type = "Cluster"
                      visibility               = "private"
                      active                   = true
                      associated_resource_id   = 6032

                      config_alletramp_hvm = {
                        protocol_type     = "iSCSI"
                        enable_ransomware = false
                      }

                      storage_server = {
                        id = 1489
                      }
					}`,
				ExpectNonEmptyPlan: true,
				PlanOnly:           true,
			},
			{
				// Update with tenants
				Config: providerConfig + `
                    resource "hpe_morpheus_datastore" "example" {
                      name = ` + "\"" + name + "\"" + `
                      datastore_type = {
                        id   = 8
                        code = "hpedatastore-alletra-mp"
                      }
                      associated_resource_type = "Cluster"
                      visibility               = "private"
                      active                   = true
                      associated_resource_id   = 6032

                      resource_permissions = {
                        groups = [
                          {
                            id = 17830
                          }
                        ]
                      }
                      tenants = [
                        {
                          id = 466
                        }
                      ]
                    
                      config_alletramp_hvm = {
                        protocol_type     = "iSCSI"
                        enable_ransomware = false
                      }

                      storage_server = {
                        id = 1489
                      }
					}`,
				Check:    checkFnAll,
				PlanOnly: false,
			},
			{
				// Check change in tenants detected
				Config: providerConfig + `
                    resource "hpe_morpheus_datastore" "example" {
                      name = ` + "\"" + name + "\"" + `
                      datastore_type = {
                        id   = 8
                        code = "hpedatastore-alletra-mp"
                      }
                      associated_resource_type = "Cluster"
                      visibility               = "private"
                      active                   = true
                      associated_resource_id   = 6032

                      resource_permissions = {
                        groups = [
                          {
                            id = 17830
                          }
                        ]
                      }

                      config_alletramp_hvm = {
                        protocol_type     = "iSCSI"
                        enable_ransomware = false
                      }

                      storage_server = {
                        id = 1489
                      }
					}`,
				ExpectNonEmptyPlan: true,
				PlanOnly:           true,
			},
		},
	})
}

func TestAccMorpheusDatastoreValidationOneOf(t *testing.T) {
	defer testhelpers.RecordResult(t)

	t.Parallel()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// checks plan fails when visibility has invalid value
				Config: `
                    resource "hpe_morpheus_datastore" "example" {
                      name = "TestDatastore"
                      datastore_type = {
                        id   = 8
                        code = "hpedatastore-alletra-mp"
                      }
                      associated_resource_type = "Cluster"
                      visibility               = "invalid"
                      active                   = true
                      associated_resource_id   = 6032
                    
                      config_alletramp_hvm = {
                        protocol_type     = "iSCSI"
                        enable_ransomware = false
                      }
					}`,
				ExpectError: regexp.MustCompile(`Attribute visibility value must be one of:`),
			},
			{
				// checks plan fails when associated_resource_type has invalid value
				Config: `
                    resource "hpe_morpheus_datastore" "example" {
                      name = "TestDatastore"
                      datastore_type = {
                        id   = 8
                        code = "hpedatastore-alletra-mp"
                      }
                      associated_resource_type = "invalid"
                      visibility               = "private"
                      active                   = true
                      associated_resource_id   = 6032
                    
                      config_alletramp_hvm = {
                        protocol_type     = "iSCSI"
                        enable_ransomware = false
                      }
					}`,
				ExpectError: regexp.MustCompile(`Attribute associated_resource_type value must be one of:`),
			},
			{
				// checks plan fails when two config blocks are specified
				Config: `
                    resource "hpe_morpheus_datastore" "example" {
                      name = "TestDatastore"
                      datastore_type = {
                        id   = 8
                        code = "hpedatastore-alletra-mp"
                      }
                      associated_resource_type = "Cluster"
                      visibility               = "private"
                      active                   = true
                      associated_resource_id   = 6032
                    
                      config_alletramp_hvm = {
                        protocol_type     = "iSCSI"
                        enable_ransomware = false
                      }

                      config_nfs = {
                        source_dir_path = "/dir/tmp"
                        source_hostname = "nfs.example.com"
                      }
					}`,
				ExpectError: regexp.MustCompile(`Invalid Attribute Combination`),
			},
			{
				// checks plan fails when hvm has invaldi protocol_type
				Config: `
                    resource "hpe_morpheus_datastore" "example" {
                      name = "TestDatastore"
                      datastore_type = {
                        id   = 8
                        code = "hpedatastore-alletra-mp"
                      }
                      associated_resource_type = "Cluster"
                      visibility               = "private"
                      active                   = true
                      associated_resource_id   = 6032
                    
                      config_alletramp_hvm = {
                        protocol_type     = "invalid"
                        enable_ransomware = false
                      }
					}`,
				ExpectError: regexp.MustCompile(
					`Attribute config_alletramp_hvm.protocol_type value must be one of`,
				),
			},
			{
				// checks plan fails when nfs source_version is invalid
				Config: `
                    resource "hpe_morpheus_datastore" "example" {
                      name = "TestDatastore"
                      datastore_type = {
                        id   = 8
                        code = "hpedatastore-alletra-mp"
                      }
                      associated_resource_type = "Cluster"
                      visibility               = "private"
                      active                   = true
                      associated_resource_id   = 6032
 
                      config_nfs = {
                        source_dir_path = "/dir/tmp"
                        source_hostname = "nfs.example.com"
                        source_version  = "invalid"

                      }
					}`,
				ExpectError: regexp.MustCompile(`Attribute config_nfs.source_version value must be one of`),
			},
		},
	})
}

func TestAccMorpheusDatastoreValidationRequiredAttrs(t *testing.T) {
	defer testhelpers.RecordResult(t)

	t.Parallel()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// checks plan fails when associated_resource_id is removed
				Config: `
                    resource "hpe_morpheus_datastore" "example" {
                      name = "TestDatastore"
                      datastore_type = {
                        id   = 8
                        code = "hpedatastore-alletra-mp"
                      }
                      associated_resource_type = "Cluster"
                    
                      config_alletramp_hvm = {
                        protocol_type     = "iSCSI"
                        enable_ransomware = false
                      }
					}`,
				ExpectError: regexp.MustCompile(`The argument "associated_resource_id" is required`),
			},
			{
				// checks plan fails when associated_resource_type is removed
				Config: `
                    resource "hpe_morpheus_datastore" "example" {
                      name = "TestDatastore"
                      datastore_type = {
                        id   = 8
                        code = "hpedatastore-alletra-mp"
                      }
                      associated_resource_id   = 6032
                    
                      config_alletramp_hvm = {
                        protocol_type     = "iSCSI"
                        enable_ransomware = false
                      }
					}`,
				ExpectError: regexp.MustCompile(`The argument "associated_resource_type" is required`),
			},
			{
				// checks plan fails when Name is removed
				Config: `
                    resource "hpe_morpheus_datastore" "example" {
                      datastore_type = {
                        code = "hpedatastore-alletra-mp"
                      }
                      associated_resource_type = "Cluster"
                      visibility               = "private"
                      active                   = true
                      associated_resource_id   = 6032
                    
                      config_alletramp_hvm = {
                        protocol_type     = "iSCSI"
                        enable_ransomware = false
                      }
					}`,
				ExpectError: regexp.MustCompile(`The argument "name" is required`),
			},
			{
				// checks plan fails when datastore_type.id is removed
				Config: `
                    resource "hpe_morpheus_datastore" "example" {
                      name = "TestDatastore"
                      datastore_type = {
                        code = "hpedatastore-alletra-mp"
                      }
                      associated_resource_type = "Cluster"
                      associated_resource_id   = 6032
                    
                      config_alletramp_hvm = {
                        protocol_type     = "iSCSI"
                        enable_ransomware = false
                      }
					}`,
				ExpectError: regexp.MustCompile(`Incorrect attribute value type`),
			},
			{
				// checks plan fails when datastore_type.code is removed
				Config: `
                    resource "hpe_morpheus_datastore" "example" {
                      name = "TestDatastore"
                      datastore_type = {
                        id   = 8
                      }
                      associated_resource_type = "Cluster"
                      associated_resource_id   = 6032
                    
                      config_alletramp_hvm = {
                        protocol_type     = "iSCSI"
                        enable_ransomware = false
                      }
					}`,
				ExpectError: regexp.MustCompile(`Incorrect attribute value type`),
			},
		},
	})
}
