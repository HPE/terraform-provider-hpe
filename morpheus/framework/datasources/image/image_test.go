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

	"github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/image"
	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/systemoverride"
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
	systemoverride.ParseFlags()
	code := m.Run()
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func TestAccMorpheusImageDatasourceById(t *testing.T) {
	defer testhelpers.RecordResult(t)
	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	name := acctest.RandomWithPrefix(t.Name())

	testSystem := systemoverride.GetPreferred(t, "zodiac")
	providerConfig := testhelpers.ProviderBlockForServer(testSystem)

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

	dataSourceConfig, err := testhelpers.RenderExample(t,
		"example-id.tf.tmpl", "Id", "hpe_morpheus_image.example_image.id")
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_image.test",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_image.test",
			"image_type",
			"qcow2",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + dependencyConfig + dataSourceConfig,
				Check:  checkFn,
			},
		},
	})
}

func TestAccMorpheusImageDatasourceByName(t *testing.T) {
	defer testhelpers.RecordResult(t)
	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	name := acctest.RandomWithPrefix(t.Name())

	testSystem := systemoverride.GetPreferred(t, "zodiac")
	providerConfig := testhelpers.ProviderBlockForServer(testSystem)

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

	dataSourceConfig, err := testhelpers.RenderExample(t,
		"example-name.tf.tmpl", "Name", "hpe_morpheus_image.example_image.name")
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_image.test",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_image.test",
			"image_type",
			"qcow2",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + dependencyConfig + dataSourceConfig,
				Check:  checkFn,
			},
		},
	})
}

func TestAccMorpheusImageDatasourceByNameAndImageType(t *testing.T) {
	defer testhelpers.RecordResult(t)
	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	name := acctest.RandomWithPrefix(t.Name())
	imageType := "qcow2"

	testSystem := systemoverride.GetPreferred(t, "zodiac")
	providerConfig := testhelpers.ProviderBlockForServer(testSystem)

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

	dataSourceConfig, err := testhelpers.RenderExample(t, "example-name-type.tf.tmpl",
		"Name", "hpe_morpheus_image.example_image.name",
		"Type", "\"qcow2\"",
	)
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_image.test",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_image.test",
			"image_type",
			imageType,
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), sdkv2morpheus.Provider()),
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
	t.Parallel()

	dataSourceConfig := `
data "hpe_morpheus_image" "test" {
  image_type = "qcow2"
}
`

	errMatch := regexp.MustCompile("Attribute \"name\" must be specified when \"image_type\" is specified")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), sdkv2morpheus.Provider()),
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
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config:      providerConfigOffline + dataSourceConfig,
				ExpectError: errMatch,
			},
		},
	})
}
