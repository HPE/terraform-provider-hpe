// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

//go:generate ../../../../bin/render example_alletramp_hvm.tf.tmpl Name "TestAlletraDatastore" AssociatedResourceID 1 StorageServerID 1 GroupID 1 TenantID 1
//go:generate ../../../../bin/render example_alletramp_bm.tf.tmpl Name "TestAlletraDatastore" CloudName "Metal" StorageServerName "alletra-mp-1" ResourcePoolName "Default" GroupName "Default" TenantName "Tenant A"

package datastore_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
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

// Tests that our example file template used for docs is a valid config
func TestAccMorpheusDatastoreResourceExampleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.Alletra)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

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

// TestAccMorpheusDatastoreResourceBmaasExampleOk creates an HPE Alletra MP Bare
// Metal (BMaaS) datastore from the docs example, so the example is proven valid
// and the Cloud-scoped create path -- storage_server, resource_pool (sent as
// zonePool) and config_alletramp_bmaas -- is exercised. The HVM example test
// covers only the Cluster path, which is how the published BMaaS example with
// hard-coded ids went unverified until QA copied it.
//
// Everything is resolved by name from the bare-metal fixture
// (testhelpers/metal_fixture.go), the resource pool through the zonePools
// option source. The storage plugin refuses a second online BMaaS datastore
// for the same cloud, storage server and protocol, so the fixture's cloud must
// not already have an iSCSI BMaaS datastore on that storage server.
func TestAccMorpheusDatastoreResourceBmaasExampleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.Metal, capabilities.Alletra)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	cloudName := testhelpers.MetalCloudName(t)
	storageServerName := testhelpers.MetalStorageServerName(t)
	poolName := testhelpers.MetalPoolName(t)
	groupName := testhelpers.MetalGroupName(t)
	tenantName := testhelpers.MetalTenantName(t)

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := testhelpers.RenderExample(
		t, "example_alletramp_bm.tf.tmpl",
		"Name", name,
		"CloudName", cloudName,
		"StorageServerName", storageServerName,
		"ResourcePoolName", poolName,
		"GroupName", groupName,
		"TenantName", tenantName,
	)
	if err != nil {
		t.Fatal(err)
	}

	const rn = "hpe_morpheus_datastore.example"

	checks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(rn, "name", name),
		resource.TestCheckResourceAttr(rn, "associated_resource_type", "Cloud"),
		resource.TestCheckResourceAttrPair(rn, "associated_resource_id", "data.hpe_morpheus_cloud.metal", "id"),
		resource.TestCheckResourceAttr(rn, "datastore_type.code", "hpedatastore-alletra-mp-bmaas"),
		resource.TestCheckResourceAttrSet(rn, "datastore_type.id"),
		resource.TestCheckResourceAttr(rn, "config_alletramp_bmaas.protocol_type", "iSCSI"),
		resource.TestCheckResourceAttrPair(rn, "storage_server.id", "data.hpe_morpheus_storage_server.alletra", "id"),
		resource.TestCheckResourceAttrPair(rn, "resource_pool.id", "data.hpe_morpheus_option_source_value.pool", "attributes.id"),
		resource.TestCheckTypeSetElemAttrPair(rn, "resource_permissions.groups.*.id", "data.hpe_morpheus_group.example", "id"),
		resource.TestCheckTypeSetElemAttrPair(rn, "tenants.*.id", "data.hpe_morpheus_tenant.example", "id"),
	)

	// Import has no plan to fall back on, so what it reads back from the API
	// is checked attribute by attribute rather than against the whole prior
	// state (tenants gains the root tenant on read, as the HVM test notes).
	importCheck := func(states []*terraform.InstanceState) error {
		for _, is := range states {
			if is.Attributes["name"] != name {
				continue
			}

			for attr, want := range map[string]string{
				"associated_resource_type":             "Cloud",
				"datastore_type.code":                  "hpedatastore-alletra-mp-bmaas",
				"config_alletramp_bmaas.protocol_type": "iSCSI",
			} {
				if got := is.Attributes[attr]; got != want {
					return fmt.Errorf("%s = %q after import, want %q", attr, got, want)
				}
			}

			for _, attr := range []string{"storage_server.id", "resource_pool.id"} {
				if is.Attributes[attr] == "" {
					return fmt.Errorf("%s is empty after import", attr)
				}
			}

			return nil
		}

		return fmt.Errorf("imported state for datastore %q not found", name)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + resourceConfig,
				Check:  checks,
			},
			{
				Config:             providerConfig + resourceConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
				Check:              checks,
			},
			{
				ImportState:      true,
				ResourceName:     rn,
				ImportStateCheck: importCheck,
			},
		},
	})
}

// Tests that our example file template used for docs is a valid config
func TestAccMorpheusDatastoreResourceUpdateOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.Alletra)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
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
				// The rename is an in-place update that must keep id known: planned
				// as unknown it forces the replacement of anything that references
				// the datastore through a create-only attribute (MORPH-18309).
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("hpe_morpheus_datastore.example", plancheck.ResourceActionUpdate),
						plancheck.ExpectKnownValue("hpe_morpheus_datastore.example", tfjsonpath.New("id"), knownvalue.NotNull()),
					},
				},
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

func TestAccMorpheusDatastoreResourceValidationOneOf(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

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

func TestAccMorpheusDatastoreResourceValidationRequiredAttrs(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

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
