// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package librarycontainerscript_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

func TestMain(m *testing.M) {
	code := m.Run()
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func TestAccMorpheusLibraryContainerScriptBasic(t *testing.T) {
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix("tf-acc")

	createConfig := `
resource "hpe_morpheus_library_container_script" "test" {
  name        = "` + name + `"
  script_type = "bash"
  script      = "#!/bin/bash\necho hello"
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + createConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("hpe_morpheus_library_container_script.test", "id"),
					resource.TestCheckResourceAttr("hpe_morpheus_library_container_script.test", "name", name),
					resource.TestCheckResourceAttr("hpe_morpheus_library_container_script.test", "script_type", "bash"),
				),
			},
			{
				ResourceName:      "hpe_morpheus_library_container_script.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccMorpheusLibraryContainerScriptUpdate(t *testing.T) {
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix("tf-acc")

	createConfig := `
resource "hpe_morpheus_library_container_script" "test" {
  name        = "` + name + `"
  script_type = "bash"
  script      = "#!/bin/bash\necho hello"
}
`

	updateConfig := `
resource "hpe_morpheus_library_container_script" "test" {
  name        = "` + name + `"
  script_type = "bash"
  script      = "#!/bin/bash\necho updated"
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + createConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_library_container_script.test", "script", "#!/bin/bash\necho hello"),
				),
			},
			{
				Config: providerConfig + updateConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"hpe_morpheus_library_container_script.test", "script", "!/bin/bash\necho updated"),
				),
			},
		},
	})
}
