// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package network_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
)

// TestAccMorpheusNetworkResourceValidationMissingName tests that the
// resource fails validation when the required 'name' field is missing
func TestAccMorpheusNetworkResourceValidationMissingName(t *testing.T) {
	if capabilities.Missing(t, capabilities.All) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)

	providerConfig := testhelpers.ProviderBlock()

	// Configuration missing the required 'name' field
	config := providerConfig + `
resource "hpe_morpheus_network" "test" {
	cloud_id = 1
	group_id = 1
	type_id  = 1
	labels   = ["terraform", "acctest", "hpe_morpheus_network", "sweepable"]
	# name is missing - this should cause validation error
}
`

	expected := `The argument "name" is required, but no definition was found`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      config,
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(expected),
			},
		},
	})
}

// TestAccMorpheusNetworkResourceValidationMissingCloudId tests that the
// resource fails validation when the required 'cloud_id' field is missing
func TestAccMorpheusNetworkResourceValidationMissingCloudId(t *testing.T) {
	if capabilities.Missing(t, capabilities.All) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)

	providerConfig := testhelpers.ProviderBlock()

	// Configuration missing the required 'cloud_id' field
	config := providerConfig + `
resource "hpe_morpheus_network" "test" {
	name     = "TestAccMorpheusNetworkResourceValidationMissingCloudId"
	group_id = 1
	type_id  = 1
	labels   = ["terraform", "acctest", "hpe_morpheus_network", "sweepable"]
	# cloud_id is missing - this should cause validation error
}
`

	expected := `The argument "cloud_id" is required, but no definition was found`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      config,
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(expected),
			},
		},
	})
}

// TestAccMorpheusNetworkResourceValidationMissingGroupId tests that the
// resource fails validation when the required 'group_id' field is missing
func TestAccMorpheusNetworkResourceValidationMissingGroupId(t *testing.T) {
	if capabilities.Missing(t, capabilities.All) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)

	providerConfig := testhelpers.ProviderBlock()

	// Configuration missing the required 'group_id' field
	config := providerConfig + `
resource "hpe_morpheus_network" "test" {
	name     = "TestAccMorpheusNetworkResourceValidationMissingGroupId"
	cloud_id = 1
	type_id  = 1
	labels   = ["terraform", "acctest", "hpe_morpheus_network", "sweepable"]
	# group_id is missing - this should cause validation error
}
`

	expected := `The argument "group_id" is required, but no definition was found`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      config,
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(expected),
			},
		},
	})
}

// TestAccMorpheusNetworkResourceValidationMissingTypeId tests that the
// resource fails validation when the required 'type_id' field is missing
func TestAccMorpheusNetworkResourceValidationMissingTypeId(t *testing.T) {
	if capabilities.Missing(t, capabilities.All) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)

	providerConfig := testhelpers.ProviderBlock()

	// Configuration missing the required 'type_id' field
	config := providerConfig + `
resource "hpe_morpheus_network" "test" {
	name     = "TestAccMorpheusNetworkResourceValidationMissingTypeId"
	cloud_id = 1
	group_id = 1
	labels   = ["terraform", "acctest", "hpe_morpheus_network", "sweepable"]
	# type_id is missing - this should cause validation error
}
`

	expected := `The argument "type_id" is required, but no definition was found`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      config,
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(expected),
			},
		},
	})
}

// TestAccMorpheusNetworkResourceValidationInvalidConfig tests that the
// resource fails validation when config field has invalid type (string
// instead of object)
func TestAccMorpheusNetworkResourceValidationInvalidConfig(t *testing.T) {
	if capabilities.Missing(t, capabilities.All) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)

	providerConfig := testhelpers.ProviderBlock()

	// Configuration with invalid config type
	config := providerConfig + `
resource "hpe_morpheus_network" "test" {
	name     = "TestAccMorpheusNetworkResourceValidationInvalidConfig"
	cloud_id = 1
	group_id = 1
	type_id  = 1
	labels   = ["terraform", "acctest", "hpe_morpheus_network", "sweepable"]
	config   = "invalid"
}
`

	expected := `attribute must be a valid object/map`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      config,
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(expected),
			},
		},
	})
}

// TestAccMorpheusNetworkResourceValidationInvalidPoolId tests that the
// resource fails validation when pool_id field has invalid type (string
// instead of int)
func TestAccMorpheusNetworkResourceValidationInvalidPoolId(t *testing.T) {
	if capabilities.Missing(t, capabilities.All) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)

	providerConfig := testhelpers.ProviderBlock()

	// Configuration with invalid pool_id type
	config := providerConfig + `
resource "hpe_morpheus_network" "test" {
	name     = "TestAccMorpheusNetworkResourceValidationInvalidPoolId"
	cloud_id = 1
	group_id = 1
	type_id  = 1
	labels   = ["terraform", "acctest", "hpe_morpheus_network", "sweepable"]
	pool_id  = "not_valid_int"
}
`

	expected := `Incorrect attribute value type`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      config,
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(expected),
			},
		},
	})
}

// TestAccMorpheusNetworkResourceValidationInvalidTenantIds tests that the
// resource fails validation when tenant_ids field has invalid type
// (string instead of set)
func TestAccMorpheusNetworkResourceValidationInvalidTenantIds(t *testing.T) {
	if capabilities.Missing(t, capabilities.All) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)

	providerConfig := testhelpers.ProviderBlock()

	// Configuration with invalid tenant_ids type
	config := providerConfig + `
resource "hpe_morpheus_network" "test" {
	name       = "TestAccMorpheusNetworkResourceValidationInvalidTenantIds"
	cloud_id   = 1
	group_id   = 1
	type_id    = 1
	labels     = ["terraform", "acctest", "hpe_morpheus_network", "sweepable"]
	tenant_ids = "not_valid_set"
}
`

	expected := `set of number required`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      config,
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(expected),
			},
		},
	})
}

// TestAccMorpheusNetworkResourceValidationValidConfigNull tests that the
// resource accepts null config value without validation error
func TestAccMorpheusNetworkResourceValidationValidConfigNull(t *testing.T) {
	if capabilities.Missing(t, capabilities.All) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	// Configuration with valid null config
	config := providerConfig + `
resource "hpe_morpheus_network" "test" {
	name     = "TestAccMorpheusNetworkResourceValidationValidConfigNull"
	cloud_id = 1
	group_id = 1
	type_id  = 1
	labels   = ["terraform", "acctest", "hpe_morpheus_network", "sweepable"]
	config   = null
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:             config,
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}
