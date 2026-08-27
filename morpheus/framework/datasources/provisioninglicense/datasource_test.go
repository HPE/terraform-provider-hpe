// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package provisioninglicense_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/provisioninglicense"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func TestAccMorpheusFindProvisioningLicenseById(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	// A provisioning license can only be created with a genuine license key.
	// The resource's own acceptance test uses PlanOnly for the same reason.
	// Provide a real key via TF_VAR_testacc_provisioning_license_key to exercise
	// the create-then-lookup path; otherwise skip.
	licenseKey := os.Getenv("TF_VAR_testacc_provisioning_license_key")
	if licenseKey == "" {
		t.Skip("Skipping: set TF_VAR_testacc_provisioning_license_key " +
			"(a valid provisioning license key) to run this test")
	}

	name := acctest.RandomWithPrefix(t.Name())
	providerConfig := testhelpers.ProviderBlock()

	resourceConfig := `
resource "hpe_morpheus_provisioning_license" "test" {
  name            = "` + name + `"
  license_type    = "win"
  license_key_wo  = "` + licenseKey + `"
  description     = "test license"
}
`
	dataSourceConfig, err := testhelpers.RenderExample(t,
		"example-id.tf.tmpl", "Id", "hpe_morpheus_provisioning_license.test.id")
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_provisioning_license.example", "name", name,
		),
		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_provisioning_license.test", "id",
			"data.hpe_morpheus_provisioning_license.example", "id",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + resourceConfig,
			},
			{
				Config: providerConfig + resourceConfig + dataSourceConfig,
				Check:  checkFn,
			},
		},
	})
}

func TestAccMorpheusFindProvisioningLicenseByName(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	// See TestAccMorpheusFindProvisioningLicenseById: a real license key is
	// required to actually create a provisioning license.
	licenseKey := os.Getenv("TF_VAR_testacc_provisioning_license_key")
	if licenseKey == "" {
		t.Skip("Skipping: set TF_VAR_testacc_provisioning_license_key " +
			"(a valid provisioning license key) to run this test")
	}

	name := acctest.RandomWithPrefix(t.Name())
	providerConfig := testhelpers.ProviderBlock()

	resourceConfig := `
resource "hpe_morpheus_provisioning_license" "test" {
  name            = "` + name + `"
  license_type    = "win"
  license_key_wo  = "` + licenseKey + `"
  description     = "test license"
}
`
	dataSourceConfig, err := testhelpers.RenderExample(t,
		"example-name.tf.tmpl", "Name", name)
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_provisioning_license.example", "name", name,
		),
		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_provisioning_license.test", "id",
			"data.hpe_morpheus_provisioning_license.example", "id",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + resourceConfig,
			},
			{
				Config: providerConfig + resourceConfig + dataSourceConfig,
				Check:  checkFn,
			},
		},
	})
}

func TestAccMorpheusFindProvisioningLicenseNotFound(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	config := providerConfig + `
      data "hpe_morpheus_provisioning_license" "test" {
        name = "______"
      }`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile(provisioninglicense.ErrorNoLicenseFound),
			},
		},
	})
}

func TestAccMorpheusFindProvisioningLicenseNoSearchAttrs(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	t.Parallel()

	// A real connection is used so the data source Read runs and returns the
	// "no valid search terms" error; with an unconfigured provider the mux
	// provider fails earlier with a connection error and the validation path is
	// never reached.
	config := testhelpers.ProviderBlock() + `
      data "hpe_morpheus_provisioning_license" "test" {
      }`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile(provisioninglicense.ErrorNoValidSearchTerms),
			},
		},
	})
}
