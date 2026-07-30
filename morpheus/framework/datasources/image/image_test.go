// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

//go:generate ../../../../bin/render example-id.tf.tmpl Id 99
//go:generate ../../../../bin/render example-name.tf.tmpl Name "\"Example name\""
//go:generate ../../../../bin/render example-name-type.tf.tmpl Name "\"Example name\"" Type "\"qcow2\""

package image_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/image"
	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

const providerConfigOffline = `
provider "hpe" {
  morpheus {
    url          = ""
    username     = ""
    password     = ""
  }
}
`

func TestMain(m *testing.M) {
	code := m.Run()
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

// TestAccMorpheusImageDatasource creates a single image and reads it back three
// ways in one apply: by id, by name, and by name+image_type. Consolidating the
// three former lookup variants (previously separate parallel tests that each
// created their own image) onto one created image removes the concurrent
// image-create contention that produced transient 403s, while keeping identical
// datasource coverage.
func TestAccMorpheusImageDatasource(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	name := acctest.RandomWithPrefix(t.Name())
	imageType := "qcow2"

	providerConfig := testhelpers.ProviderBlock()

	imageConfig, err := image.RenderImageConfig(t, map[string]string{
		"Name":              name,
		"OsTypeId":          "data.hpe_morpheus_os_type.test.id",
		"StorageProviderId": "data.hpe_morpheus_storage_bucket.test.id",
	})
	if err != nil {
		t.Fatal(err)
	}

	dependencyConfig := `
data "hpe_morpheus_os_type" "test" {
	name = "linux"
}

data "hpe_morpheus_storage_bucket" "test" {
	name = "Local Storage"
}
` + imageConfig

	// The example-*.tf.tmpl templates hardcode the "test" datasource label, so
	// they cannot be combined in a single config. Inline the three lookups here
	// so they all reference the one image created above.
	dataSourceConfig := `
data "hpe_morpheus_image" "by_id" {
  id = hpe_morpheus_image.example_image.id
}

data "hpe_morpheus_image" "by_name" {
  name = hpe_morpheus_image.example_image.name
}

data "hpe_morpheus_image" "by_name_type" {
  name       = hpe_morpheus_image.example_image.name
  image_type = "qcow2"
}
`

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr("data.hpe_morpheus_image.by_id", "name", name),
		resource.TestCheckResourceAttr("data.hpe_morpheus_image.by_id", "image_type", imageType),
		resource.TestCheckResourceAttr("data.hpe_morpheus_image.by_name", "name", name),
		resource.TestCheckResourceAttr("data.hpe_morpheus_image.by_name", "image_type", imageType),
		resource.TestCheckResourceAttr("data.hpe_morpheus_image.by_name_type", "name", name),
		resource.TestCheckResourceAttr("data.hpe_morpheus_image.by_name_type", "image_type", imageType),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + dependencyConfig + dataSourceConfig,
				Check:  checkFn,
			},
		},
	})
}

// verify that we get an error if `image_type` is specified without `name`
func TestAccMorpheusImageDatasourceByImageTypeOnly(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	t.Parallel()

	dataSourceConfig := `
data "hpe_morpheus_image" "test" {
  image_type = "qcow2"
}
`

	errMatch := regexp.MustCompile("Attribute \"name\" must be specified when \"image_type\" is specified")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), sdkv2morpheus.Provider()),
		IsUnitTest:               true,
		Steps: []resource.TestStep{
			{
				Config:      providerConfigOffline + dataSourceConfig,
				ExpectError: errMatch,
			},
		},
	})
}

// this should fail due to a conflict between id and name/image_type
func TestAccMorpheusImageDatasourceBothAttrs(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	t.Parallel()

	dataSourceConfig := `
data "hpe_morpheus_image" "test" {
  image_type = "qcow2"
  name = "Example Image"
  id = 5
}
`

	errMatch := regexp.MustCompile("Attribute \"(.*)\" cannot be specified when \"id\" is specified")
	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config:      providerConfigOffline + dataSourceConfig,
				ExpectError: errMatch,
			},
		},
	})
}
