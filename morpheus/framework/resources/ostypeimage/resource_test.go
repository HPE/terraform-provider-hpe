// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package ostypeimage_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/ostypeimage"
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

// Tests that our example file template used for docs is a valid config.
func TestAccMorpheusOsTypeImageExampleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	resourceConfig, err := ostypeimage.RenderOsTypeImageConfig(t, nil)
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet(
			"hpe_morpheus_os_type_image.example",
			"id",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_os_type_image.example",
			"os_type_id",
			"1",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_os_type_image.example",
			"virtual_image_id",
			"257",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_os_type_image.example",
			"cloud_id",
			"1",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_os_type_image.example",
			"provision_type_id",
			"22",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				Check:              checkFn,
				PlanOnly:           false,
			},
			{
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"os_type_id"},
				ResourceName:            "hpe_morpheus_os_type_image.example",
				Check:                   checkFn,
			},
		},
	})
}

// Tests creating with only the required attributes.
func TestAccMorpheusOsTypeImageRequiredAttrsOk(t *testing.T) {
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	resourceConfig := `
resource "hpe_morpheus_os_type_image" "required_only" {
  os_type_id       = 1
  virtual_image_id = 257
}
`

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet(
			"hpe_morpheus_os_type_image.required_only",
			"id",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_os_type_image.required_only",
			"os_type_id",
			"1",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_os_type_image.required_only",
			"virtual_image_id",
			"257",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				Check:              checkFn,
				PlanOnly:           false,
			},
			{
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"os_type_id"},
				ResourceName:            "hpe_morpheus_os_type_image.required_only",
				Check:                   checkFn,
			},
		},
	})
}
