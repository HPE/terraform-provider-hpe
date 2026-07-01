// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instancesnapshot_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/instance"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/instancesnapshot"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func TestAccMorpheusInstanceSnapshotResource(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping acceptance test in short mode")
	}

	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

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

	config, err := instancesnapshot.RenderInstanceSnapshotConfig(t, map[string]string{
		"Name":        snapshotName,
		"InstanceId":  "hpe_morpheus_instance.example.id",
		"Description": "Acceptance test snapshot",
	})
	if err != nil {
		t.Fatalf("failed to render config: %v", err)
	}

	resourceName := "hpe_morpheus_instance_snapshot.example"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: testhelpers.ProviderBlock() + instanceConfig + config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", snapshotName),
					resource.TestCheckResourceAttr(resourceName, "description", "Acceptance test snapshot"),
					resource.TestCheckResourceAttr(resourceName, "memory_snapshot", "false"),
					resource.TestCheckResourceAttr(resourceName, "retain_on_delete", "false"),
					resource.TestCheckResourceAttr(resourceName, "status", "complete"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "date_created"),
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

					return fmt.Sprintf("%s.%s",
						rs.Primary.Attributes["instance_id"],
						rs.Primary.Attributes["id"],
					), nil
				},
				ImportStateVerify: true,
				// retain_on_delete is not returned by the API
				ImportStateVerifyIgnore: []string{"retain_on_delete", "timeouts"},
			},
		},
	})
}
