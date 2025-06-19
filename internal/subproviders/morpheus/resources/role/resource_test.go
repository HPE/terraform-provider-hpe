// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

//go:generate go run ../../../../../cmd/render example.tf.tmpl Name "ExampleRole" Multitenant "false" Description "An example role" RoleType "user"

package role_test

import (
	"os"
	"testing"

	"github.com/HPE/terraform-provider-hpe/internal/provider"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/testhelpers"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
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

// Check that we can create a role with only required attributes specified
func TestAccMorpheusRoleRequiredAttrsOk(t *testing.T) {
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	resourceConfig, err := testhelpers.RenderExample(t, "example-required.tf.tmpl",
		"Name", "TestAccMorpheusRoleRequiredAttrsOk")
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.example_required",
			"name",
			"TestAccMorpheusRoleRequiredAttrsOk",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.example_required",
			"multitenant",
			"false",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.example_required",
			"role_type",
			"user",
		),
		resource.TestCheckNoResourceAttr(
			"hpe_morpheus_role.example_required",
			"description",
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
				ResourceName:      "hpe_morpheus_role.example_required",
				Check:             checkFn,
			},
		},
	})
}

// Check that we can create a role with all attributes specified
func TestAccMorpheusRoleAllAttrsOk(t *testing.T) {
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	resourceConfig, err := testhelpers.RenderExample(t, "example-all.tf.tmpl",
		"Name", "TestAccMorpheusRoleAllAttrsOk",
		"Multitenant", "true",
		"Description", "test",
		"RoleType", "user")
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.example_all",
			"name",
			"TestAccMorpheusRoleAllAttrsOk",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.example_all",
			"description",
			"test",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.example_all",
			"multitenant",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.example_all",
			"role_type",
			"user",
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
				ResourceName:      "hpe_morpheus_role.example_all",
				Check:             checkFn,
			},
		},
	})
}

// Tests that our example file template used for docs is a valid config
func TestAccMorpheusRoleExampleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	resourceConfig, err := testhelpers.RenderExample(t, "example.tf.tmpl",
		"Name", "TestAccMorpheusRoleExampleOk",
		"Multitenant", "false",
		"Description", "a test of the example HCL config",
		"RoleType", "user")
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.example",
			"name",
			"TestAccMorpheusRoleExampleOk",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.example",
			"description",
			"a test of the example HCL config",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.example",
			"multitenant",
			"false",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_role.example",
			"role_type",
			"user",
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
				ResourceName:      "hpe_morpheus_role.example",
				Check:             checkFn,
			},
		},
	})
}

// TODO: Add more acceptance tests
