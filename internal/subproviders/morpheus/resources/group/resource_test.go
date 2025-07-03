// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

//go:generate go run ../../../../../cmd/render example.tf.tmpl Name "TestGroup" Location "here" Code "aCode" Label "aLabel"

package group_test

import (
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/internal/provider"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/testhelpers"
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

// Tests that our example file template used for docs is a valid config
func TestAccMorpheusGroupExampleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())
	code := strings.ToLower(name)

	resourceConfig, err := testhelpers.RenderExample(t, "example.tf.tmpl",
		"Name", name,
		"Location", "here",
		"Code", code,
		"Label", "aLabel")
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_group.example",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_group.example",
			"code",
			code,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_group.example",
			"location",
			"here",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_group.example",
			"labels.#",
			"2",
		),
		resource.TestCheckTypeSetElemAttr(
			"hpe_morpheus_group.example",
			"labels.*",
			"aLabel1",
		),
		resource.TestCheckTypeSetElemAttr(
			"hpe_morpheus_group.example",
			"labels.*",
			"aLabel2",
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
				ImportState:       true,
				ImportStateVerify: true, // Check state post import
				ResourceName:      "hpe_morpheus_group.example",
				Check:             checkFn,
			},
		},
	})
}

func TestAccMorpheusGroupUpdateOk(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())
	code := strings.ToLower(name)

	baseChecks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_group.test",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_group.test",
			"code",
			code,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_group.test",
			"location",
			"here",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_group.test",
			"labels.#",
			"2",
		),
		resource.TestCheckTypeSetElemAttr(
			"hpe_morpheus_group.test",
			"labels.*",
			"Label1",
		),
		resource.TestCheckTypeSetElemAttr(
			"hpe_morpheus_group.test",
			"labels.*",
			"Label2",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(
		baseChecks...,
	)

	updateChecks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_group.test",
			"name",
			name+"-changed",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_group.test",
			"code",
			code+"-changed",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_group.test",
			"location",
			"here-changed",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_group.test",
			"labels.#",
			"3",
		),
	}

	checkUpdateFn := resource.ComposeAggregateTestCheckFunc(
		updateChecks...,
	)
	_ = checkUpdateFn

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "hpe_morpheus_group" "test" {
  name = "` + name + `"
  location = "here"
  code = "` + code + `"
  labels = ["Label1", "Label2"]
}`,
				Check:    checkFn,
				PlanOnly: false,
			},
			{
				Config: providerConfig + `
# checks plan has no effect
resource "hpe_morpheus_group" "test" {
	name = "` + name + `"
	location = "here"
	code = "` + code + `"
	labels = ["Label1", "Label2"]
}`,
				Check:              checkFn,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
			{
				Config: providerConfig + `
# checks plan detects change to labels
resource "hpe_morpheus_group" "test" {
	name = "` + name + `"
	location = "here"
	code = "` + code + `"
	labels = ["Label1", "Label2", "env:blah"]
}`,
				ExpectNonEmptyPlan: true,
				PlanOnly:           true,
			},
			{
				Config: providerConfig + `
# checks plan detects name change
resource "hpe_morpheus_group" "test" {
	name = "` + name + `-changed"
	location = "here"
	code = "` + code + `"
	labels = ["Label1", "Label2"]
}`,
				ExpectNonEmptyPlan: true,
				PlanOnly:           true,
			},
			{
				Config: providerConfig + `
# checks plan detects code change
resource "hpe_morpheus_group" "test" {
	name = "` + name + `"
	location = "here"
	code = "` + code + `-changed"
	labels = ["Label1", "Label2"]
}`,
				ExpectNonEmptyPlan: true,
				PlanOnly:           true,
			},
			{
				Config: providerConfig + `
# checks plan detects code set to null
resource "hpe_morpheus_group" "test" {
	name = "` + name + `"
	location = "here"
	# code = "` + code + `-changed"
	labels = ["Label1", "Label2"]
}`,
				ExpectNonEmptyPlan: true,
				PlanOnly:           true,
			},
			{
				Config: providerConfig + `
# checks plan detects location change
resource "hpe_morpheus_group" "test" {
	name = "` + name + `"
	location = "there"
	code = "` + code + `"
	labels = ["Label1", "Label2"]
}`,
				ExpectNonEmptyPlan: true,
				PlanOnly:           true,
			},
			{
				Config: providerConfig + `
# checks plan detects location set to null
resource "hpe_morpheus_group" "test" {
	name = "` + name + `"
	# location = "here"
	code = "` + code + `"
	labels = ["Label1", "Label2"]
}`,
				ExpectNonEmptyPlan: true,
				PlanOnly:           true,
			},
			{
				Config: providerConfig + `
# checks apply of changes to all changeable fields
resource "hpe_morpheus_group" "test" {
	name = "` + name + `-changed"
	location = "here-changed"
	code = "` + code + `-changed"
	labels = ["Label1", "Label2", "Label3"]
}`,
				Check:    checkUpdateFn,
				PlanOnly: false,
			},
			{
				Config: providerConfig + `
# checks plan has no effect
resource "hpe_morpheus_group" "test" {
	name = "` + name + `-changed"
	location = "here-changed"
	code = "` + code + `-changed"
	labels = ["Label1", "Label2", "Label3"]
}`,
				Check:              checkUpdateFn,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func TestAccMorpheusGroupRequiredAttrsOk(t *testing.T) {
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	name := acctest.RandomWithPrefix(t.Name())

	providerConfig := testhelpers.ProviderBlock()

	resourceConfig := `
resource "hpe_morpheus_group" "example_required" {
  name = "` + name + `"
}
`
	checks := []resource.TestCheckFunc{
		// required
		resource.TestCheckResourceAttr(
			"hpe_morpheus_group.example_required",
			"name",
			name,
		),
		// checks for optional
		resource.TestCheckNoResourceAttr(
			"hpe_morpheus_group.example_required",
			"location",
		),
		resource.TestCheckNoResourceAttr(
			"hpe_morpheus_group.example_required",
			"code",
		),
		resource.TestCheckNoResourceAttr(
			"hpe_morpheus_group.example_required",
			"labels",
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
				ImportState:       true,
				ImportStateVerify: true, // Check state post import
				ResourceName:      "hpe_morpheus_group.example_required",
				Check:             checkFn,
			},
		},
	})
}
