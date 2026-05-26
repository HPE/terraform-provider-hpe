// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package ostypeimage_test

import (
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/image"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/ostypeimage"
	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
)

func TestMain(m *testing.M) {
	code := m.Run()
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

// Tests that our example file template used for docs is a valid config.
func TestAccMorpheusOsTypeImageResourceExampleOk(t *testing.T) {
	if capabilities.Missing(t, capabilities.All) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	t.Parallel()
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())
	code := strings.ToLower(name)

	imageConfig, err := image.RenderImageConfig(t, map[string]string{
		"Name":              name,
		"OsTypeId":          "hpe_morpheus_os_type.test.id",
		"StorageProviderId": "1",
	})
	if err != nil {
		t.Fatal(err)
	}

	datasourceConfig := `
data "hpe_morpheus_provision_type" "test" {
	name = "KVM"
}

data "hpe_morpheus_cloud" "test" {
	name = "hvm"
}
`
	// The virtual images API will NOT update the underlying osType.Id of the virtual image.
	// So we set up a clean room scenario with the virtual image we wish to create
	// an OS Type image from, with the virtual image's os_type set correctly to an os type
	// that we create.
	dependencyConfig := `
resource "hpe_morpheus_os_type" "test" {
  name      = "` + name + `"
  code      = "` + code + `"
  platform  = "linux"
  bit_count = 64
}
` + imageConfig

	resourceConfig, err := ostypeimage.RenderOsTypeImageConfig(t, map[string]string{
		"OsTypeId":        "hpe_morpheus_os_type.test.id",
		"VirtualImageId":  "hpe_morpheus_image.example_image.id",
		"CloudId":         "data.hpe_morpheus_cloud.test.id",
		"ProvisionTypeId": "data.hpe_morpheus_provision_type.test.id",
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet(
			"hpe_morpheus_os_type_image.example",
			"id",
		),
		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_os_type_image.example",
			"os_type_id",
			"hpe_morpheus_os_type.test",
			"id",
		),
		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_os_type_image.example",
			"virtual_image_id",
			"hpe_morpheus_image.example_image",
			"id",
		),
		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_os_type_image.example",
			"cloud_id",
			"data.hpe_morpheus_cloud.test",
			"id",
		),
		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_os_type_image.example",
			"provision_type_id",
			"data.hpe_morpheus_provision_type.test",
			"id",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config:             providerConfig + datasourceConfig + dependencyConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				Check:              checkFn,
				PlanOnly:           false,
			},
			{
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      "hpe_morpheus_os_type_image.example",
				Check:             checkFn,
			},
		},
	})
}

// Tests creating with only the required attributes.
func TestAccMorpheusOsTypeImageResourceRequiredAttrsOk(t *testing.T) {
	if capabilities.Missing(t, capabilities.All) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	t.Parallel()
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())
	code := strings.ToLower(name)

	imageConfig, err := image.RenderImageConfig(t, map[string]string{
		"Name":              name,
		"OsTypeId":          "hpe_morpheus_os_type.test.id",
		"StorageProviderId": "1",
	})
	if err != nil {
		t.Fatal(err)
	}

	// The virtual images API will NOT update the underlying osType.Id of the virtual image.
	// So we set up a clean room scenario with the virtual image we wish to create
	// an OS Type image from, with the virtual image's os_type set correctly to an os type
	// that we create.
	dependencyConfig := `
resource "hpe_morpheus_os_type" "test" {
  name      = "` + name + `"
  code      = "` + code + `"
  platform  = "linux"
  bit_count = 64
}
` + imageConfig

	resourceConfig := `
resource "hpe_morpheus_os_type_image" "required_only" {
  os_type_id       = resource.hpe_morpheus_os_type.test.id
  virtual_image_id = resource.hpe_morpheus_image.example_image.id
}
`

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet(
			"hpe_morpheus_os_type_image.required_only",
			"id",
		),
		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_os_type_image.required_only",
			"os_type_id",
			"hpe_morpheus_os_type.test",
			"id",
		),
		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_os_type_image.required_only",
			"virtual_image_id",
			"hpe_morpheus_image.example_image",
			"id",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config:             providerConfig + dependencyConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				Check:              checkFn,
				PlanOnly:           false,
			},
			{
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      "hpe_morpheus_os_type_image.required_only",
				Check:             checkFn,
			},
		},
	})
}
