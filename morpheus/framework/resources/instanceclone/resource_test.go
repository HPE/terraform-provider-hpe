// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instanceclone_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	morpheus "github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/instanceclone"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
)

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func TestAccMorpheusInstanceCloneResource(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping acceptance test in short mode")
	}

	if capabilities.Missing(t, capabilities.All) {
		t.Log("Skipping test due to missing capabilities")

		return
	}

	t.Parallel()
	defer testhelpers.RecordResult(t)

	cloneName := acctest.RandomWithPrefix(t.Name())
	sourceInstanceID := os.Getenv("TF_ACC_INSTANCE_ID")
	networkID := os.Getenv("TF_ACC_NETWORK_ID")

	if sourceInstanceID == "" {
		t.Skip("TF_ACC_INSTANCE_ID must be set for instance clone tests")
	}

	if networkID == "" {
		networkID = "1"
	}

	config, err := instanceclone.RenderInstanceCloneConfig(t, map[string]string{
		"Name":             cloneName,
		"SourceInstanceId": sourceInstanceID,
		"NetworkId":        networkID,
	})
	if err != nil {
		t.Fatalf("failed to render config: %v", err)
	}

	resourceName := "hpe_morpheus_instance_clone.example"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: testhelpers.ProviderBlock() + config,
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
				// source_instance_id is not returned by the API
				ImportStateVerifyIgnore: []string{
					"source_instance_id", "timeouts",
				},
			},
		},
	})
}
