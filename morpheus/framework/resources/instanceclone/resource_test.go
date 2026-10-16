// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package instanceclone_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/instance"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/instanceclone"
	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/resources/usergroup"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

// cloneSourceConfig renders the instance a clone test clones: an HVM instance
// built from the HVM fixture (testhelpers/hvm_fixture.go), created by the test
// itself. Each test clones a source of its own rather than a shared,
// pre-existing one: cloning the same instance from several tests at once
// races the restore of the source's status on the appliance and can leave it
// stuck in "cloning", and a source named by the fixture would be exactly that
// shared instance. The rendered configuration declares the data sources the
// clone reuses, by the addresses the instance package exports
// (instance.ExampleNetworkData and, when the fixture names a datastore,
// instance.ExampleDatastoreData).
func cloneSourceConfig(t *testing.T, name string) string {
	t.Helper()

	cfg, err := instance.RenderInstanceConfig(t, map[string]string{
		"Name":         name,
		"SingleVolume": "true",
	})
	if err != nil {
		t.Fatalf("failed to render the source instance: %v", err)
	}

	return cfg
}

// cloneOverrides points a clone at the rendered source instance and at the
// network -- and datastore, when the fixture names one -- the source's
// configuration resolved by name.
func cloneOverrides(name string) map[string]string {
	overrides := map[string]string{
		"Name":               name,
		"SourceInstanceName": "",
		"SourceInstanceId":   instance.ExampleInstance + ".id",
		"ClusterName":        "",
		"NetworkName":        "",
		"NetworkId":          instance.ExampleNetworkData + ".id",
	}

	if os.Getenv(testhelpers.EnvHVMDatastoreName) != "" {
		overrides["DatastoreId"] = instance.ExampleDatastoreData + ".id"
	}

	return overrides
}

func TestAccMorpheusInstanceCloneResource(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping acceptance test in short mode")
	}

	capabilities.MustHaveOrSkip(t, capabilities.All)

	t.Parallel()
	defer testhelpers.RecordResult(t)

	name := acctest.RandomWithPrefix(t.Name())
	cloneName := name + "-clone"

	sourceConfig := cloneSourceConfig(t, name)

	config, err := instanceclone.RenderInstanceCloneConfig(t, cloneOverrides(cloneName))
	if err != nil {
		t.Fatalf("failed to render config: %v", err)
	}

	resourceName := "hpe_morpheus_instance_clone.example"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: testhelpers.ProviderBlock() + sourceConfig + config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", cloneName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
				),
			},
			// Import test
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", resourceName)
					}

					return rs.Primary.Attributes["id"], nil
				},
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"timeouts",
					"write_only_revision",
				},
			},
		},
	})
}

// TestAccMorpheusInstanceCloneResourceUserGroupStorageProfile clones an instance
// with a user_group and a volume storage_profile and asserts both round-trip
// into state. On import, user_group is pinned to what the API yields -- nothing,
// since the clone API records a user group only when the source instance was
// provisioned with one (the clone request's own userGroup is applied to the
// guest but never written to the clone's config) -- and the storage profile is
// recovered best-effort, so it is not verified.
func TestAccMorpheusInstanceCloneResourceUserGroupStorageProfile(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping acceptance test in short mode")
	}

	capabilities.MustHaveOrSkip(t, capabilities.All)

	t.Parallel()
	defer testhelpers.RecordResult(t)

	name := acctest.RandomWithPrefix(t.Name())
	cloneName := name + "-clone"

	// The user group is created for the run and referenced, so no id is
	// assumed. kvm-cache-none is a standard KVM/HVM storage profile code.
	storageProfile := "kvm-cache-none"

	userGroupConfig, err := usergroup.RenderUserGroupConfig(t, map[string]string{
		"Name":        name,
		"Description": "user group for " + name,
		"UserIds":     "[]", // no members: the example's default references a whoami data source this config does not declare
	})
	if err != nil {
		t.Fatalf("failed to render the user group: %v", err)
	}

	sourceConfig := cloneSourceConfig(t, name)

	overrides := cloneOverrides(cloneName)
	overrides["UserGroup"] = "hpe_morpheus_user_group.example.id"
	overrides["StorageProfile"] = storageProfile

	config, err := instanceclone.RenderInstanceCloneConfig(t, overrides)
	if err != nil {
		t.Fatalf("failed to render config: %v", err)
	}

	resourceName := "hpe_morpheus_instance_clone.example"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config: testhelpers.ProviderBlock() + userGroupConfig + sourceConfig + config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", cloneName),
					resource.TestCheckResourceAttrPair(
						resourceName, "user_group",
						"hpe_morpheus_user_group.example", "id",
					),
					resource.TestCheckResourceAttr(
						resourceName, "volumes.0.storage_profile", storageProfile),
				),
			},
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", resourceName)
					}

					return rs.Primary.Attributes["id"], nil
				},
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"timeouts",
					// user_group: pinned by ImportStateCheck below instead, to
					// the value the API actually yields.
					"user_group",
					// storage_profile is a write-mostly volume input recovered
					// best-effort from the API on import.
					"volumes.0.storage_profile",
					"write_only_revision",
				},
				// The clone API records config.userGroup on the clone only
				// when the SOURCE instance was provisioned with one; the clone
				// request's own userGroup is applied to the guest and never
				// written to the clone's config (InstancesController.clone,
				// identical on 9.0 and 9.1). The source here has none, so an
				// import yields no user_group - and that is asserted, so a
				// change in either direction (the API starting to report it,
				// or the read inventing one) shows up here.
				ImportStateCheck: func(states []*terraform.InstanceState) error {
					if len(states) != 1 {
						return fmt.Errorf("expected one imported state, got %d", len(states))
					}

					if got, ok := states[0].Attributes["user_group"]; ok && got != "" {
						return fmt.Errorf("import yielded user_group %q; the clone API does not report "+
							"the clone's own user group (only one inherited from the source), so none was expected", got)
					}

					return nil
				},
			},
		},
	})
}
