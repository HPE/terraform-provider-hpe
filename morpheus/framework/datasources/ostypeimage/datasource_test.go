// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package ostypeimage_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/ostypeimage"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/systemoverride"
)

func TestMain(m *testing.M) {
	systemoverride.ParseFlags()
	code := m.Run()
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func TestAccMorpheusOsTypeImageDataSourceExampleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)
	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	testSystem := systemoverride.GetPreferred(t, "zodiac")
	providerConfig := testhelpers.ProviderBlockForServer(testSystem)

	dataSourceConfig, err := ostypeimage.RenderOsTypeImageDataSourceConfig(t, nil)
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_os_type_image.example",
			"id",
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_os_type_image.example",
			"tenant_id",
			"1",
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_os_type_image.example",
			"os_type_id",
			"65",
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_os_type_image.example",
			"virtual_image_name",
			"Debian 12",
		),
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_os_type_image.example",
			"virtual_image_id",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + dataSourceConfig,
				Check:  checkFn,
			},
		},
	})
}

func TestAccMorpheusOsTypeImageDataSourceSystemImageOk(t *testing.T) {
	defer testhelpers.RecordResult(t)
	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	testSystem := systemoverride.GetPreferred(t, "zodiac")
	providerConfig := testhelpers.ProviderBlockForServer(testSystem)

	dataSourceConfig, err := ostypeimage.RenderOsTypeImageDataSourceConfig(t,
		map[string]string{
			"OsTypeId":         "65",
			"VirtualImageName": "Morpheus Debian 12 20260119",
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_os_type_image.example",
			"id",
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_os_type_image.example",
			"tenant_id",
			"0",
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_os_type_image.example",
			"os_type_id",
			"65",
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_os_type_image.example",
			"virtual_image_name",
			"Morpheus Debian 12 20260119",
		),
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_os_type_image.example",
			"virtual_image_id",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + dataSourceConfig,
				Check:  checkFn,
			},
		},
	})
}

func TestAccMorpheusOsTypeImageDataSourceNotFound(t *testing.T) {
	defer testhelpers.RecordResult(t)
	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	testSystem := systemoverride.GetPreferred(t, "zodiac")
	providerConfig := testhelpers.ProviderBlockForServer(testSystem)

	dataSourceConfig, err := ostypeimage.RenderOsTypeImageDataSourceConfig(t, map[string]string{
		"VirtualImageName": "______nonexistent______",
	})
	if err != nil {
		t.Fatal(err)
	}

	expected := regexp.MustCompile(`no image with name`)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config:      providerConfig + dataSourceConfig,
				ExpectError: expected,
			},
		},
	})
}
