// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

//go:generate go run ../../../../../../cmd/render example-id.tf.tmpl Id 99
//go:generate go run ../../../../../../cmd/render example-name.tf.tmpl Name "\"Example name\""
//go:generate go run ../../../../../../cmd/render example-name-type.tf.tmpl Name "\"Example name\"" Type "\"qcow2\""

package image_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/internal/framework/provider"
	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus"
	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
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

func TestAccMorpheusImageDatasourceImageDatasourceById(t *testing.T) {
	defer testhelpers.RecordResult(t)
	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	name := acctest.RandomWithPrefix(t.Name())

	providerConfig := testhelpers.ProviderBlock()

	imageResourceConfig := `
resource "hpe_morpheus_image" "test_image" {
  name       = "` + name + `"
  image_type = "qcow2"
}
`

	dataSourceConfig, err := testhelpers.RenderExample(t,
		"example-id.tf.tmpl", "Id", "hpe_morpheus_image.test_image.id")
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
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + imageResourceConfig + dataSourceConfig,
				Check:  checkFn,
			},
		},
	})
}

func TestAccMorpheusImageDatasourceByIdExisting(t *testing.T) {
	defer testhelpers.RecordResult(t)
	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	id := "105490" // ID matching "AlmaLinux 9" on the Feature system
	name := "AlmaLinux 9"
	imageType := "azure-reference"

	providerConfig := testhelpers.ProviderBlock()

	dataSourceConfig, err := testhelpers.RenderExample(t,
		"example-id.tf.tmpl", "Id", id)
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
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + dataSourceConfig,
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

	providerConfig := testhelpers.ProviderBlock()

	imageResourceConfig := `
resource "hpe_morpheus_image" "test_image" {
  name       = "` + name + `"
  image_type = "qcow2"
}
`

	dataSourceConfig, err := testhelpers.RenderExample(t,
		"example-name.tf.tmpl", "Name", "hpe_morpheus_image.test_image.name")
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
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + imageResourceConfig + dataSourceConfig,
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
	imageType := "iso"

	providerConfig := testhelpers.ProviderBlock()

	imageResourceConfig := `
resource "hpe_morpheus_image" "test_image" {
  name       = "` + name + `"
  image_type = "` + imageType + `"
}
`

	dataSourceConfig, err := testhelpers.RenderExample(t, "example-name-type.tf.tmpl",
		"Name", "hpe_morpheus_image.test_image.name",
		"Type", "\"iso\"",
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
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + imageResourceConfig + dataSourceConfig,
				Check:  checkFn,
			},
		},
	})
}

// verify that we get an error if `image_type` is specified without `name`
func TestAccMorpheusImageDatasourceByImageTypeOnly(t *testing.T) {
	defer testhelpers.RecordResult(t)
	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	dataSourceConfig := `
data "hpe_morpheus_image" "test" {
  image_type = "qcow2"
}
`

	errMatch := regexp.MustCompile("Attribute \"name\" must be specified when \"image_type\" is specified")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		IsUnitTest:               true,
		Steps: []resource.TestStep{
			{
				Config:      providerConfig + dataSourceConfig,
				ExpectError: errMatch,
			},
		},
	})
}

// this should fail due to a conflict between id and name/image_type
func TestAccMorpheusImageDatasourceBothAttrs(t *testing.T) {
	defer testhelpers.RecordResult(t)
	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

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
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      providerConfig + dataSourceConfig,
				ExpectError: errMatch,
			},
		},
	})
}
