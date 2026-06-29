// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instancesnapshot_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/instance"
	snapshotresource "github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/instancesnapshot"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/utils/adapter"
)

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

// TestAccMorpheusInstanceSnapshotDataSource creates a snapshot with the
// resource and then reads it back through the data source by its ID.
func TestAccMorpheusInstanceSnapshotDataSource(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping acceptance test in short mode")
	}

	defer testhelpers.RecordResult(t)

	if capabilities.Missing(t, capabilities.All) {
		t.Log("Skipping test due to missing capabilities")
	}

	t.Parallel()

	instanceName := acctest.RandomWithPrefix(t.Name())
	snapshotName := acctest.RandomWithPrefix(t.Name())

	// Provision a throwaway instance to snapshot so the test is self-contained.
	instanceConfig, err := instance.RenderInstanceConfig(t, map[string]string{
		"Name": instanceName,
	})
	if err != nil {
		t.Fatalf("failed to render instance config: %v", err)
	}

	resourceConfig, err := snapshotresource.RenderInstanceSnapshotConfig(t, map[string]string{
		"Name":        snapshotName,
		"InstanceId":  "hpe_morpheus_instance.example.id",
		"Description": "Acceptance test snapshot",
	})
	if err != nil {
		t.Fatalf("failed to render config: %v", err)
	}

	dataSourceConfig := `
data "hpe_morpheus_instance_snapshot" "example" {
  id = hpe_morpheus_instance_snapshot.example.id
}`

	dsName := "data.hpe_morpheus_instance_snapshot.example"
	resourceName := "hpe_morpheus_instance_snapshot.example"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewAdaptedMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: testhelpers.ProviderBlock() + instanceConfig + resourceConfig + dataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(dsName, "id", resourceName, "id"),
					resource.TestCheckResourceAttr(dsName, "name", snapshotName),
					resource.TestCheckResourceAttr(dsName, "description", "Acceptance test snapshot"),
					resource.TestCheckResourceAttrSet(dsName, "status"),
					resource.TestCheckResourceAttrSet(dsName, "date_created"),
				),
			},
		},
	})
}
