// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

//go:build experimental

//go:generate go run ../../../../../../cmd/render example.tf.tmpl ResourceName "group_policy" Name "TestMaxMemoryGroupPolicy" Description "Example group-scoped policy" AssociatedResourceType "Group" AssociatedResourceID "1" PolicyTypeCode "maxMemory" ConfigKey "maxMemory" ConfigValue "8"
//go:generate go run ../../../../../../cmd/render example.tf.tmpl ResourceName "cloud_policy" Name "TestMaxMemoryCloudPolicy" Description "Example cloud-scoped policy" AssociatedResourceType "Cloud" AssociatedResourceID "1" PolicyTypeCode "maxMemory" ConfigKey "maxMemory" ConfigValue "8"
//go:generate go run ../../../../../../cmd/render example.tf.tmpl ResourceName "user_policy" Name "TestMaxMemoryUserPolicy" Description "Example user-scoped policy" AssociatedResourceType "User" AssociatedResourceID "1" PolicyTypeCode "maxMemory" ConfigKey "maxMemory" ConfigValue "8"
//go:generate go run ../../../../../../cmd/render example.tf.tmpl ResourceName "role_policy" Name "TestMaxMemoryRolePolicy" Description "Example role-scoped policy" AssociatedResourceType "Role" AssociatedResourceID "1" PolicyTypeCode "maxMemory" ConfigKey "maxMemory" ConfigValue "8"

package policy_test

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

// Test validation: associated_resource_id required when not Global
func TestAccMorpheusPolicyValidationResourceIdRequired(t *testing.T) {
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig := `
resource "hpe_morpheus_policy" "validation_test" {
  name = "` + name + `"
  associated_resource_type = "Group"
  
  policy_type = {
    code = "maxMemory"
  }
  
  config = {
    maxMemory = 8
  }
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      providerConfig + resourceConfig,
				ExpectError: regexp.MustCompile("associated_resource_id is required"),
			},
		},
	})
}

// Test validation: invalid policy type code
func TestAccMorpheusPolicyValidationInvalidPolicyType(t *testing.T) {
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())
	groupName := acctest.RandomWithPrefix(t.Name() + "-group")

	resourceConfig := `
resource "hpe_morpheus_group" "test" {
  name = "` + groupName + `"
  location = "test"
}

resource "hpe_morpheus_policy" "validation_test" {
  name = "` + name + `"
  associated_resource_type = "Group"
  associated_resource_id = hpe_morpheus_group.test.id
  
  policy_type = {
    code = "invalidPolicyType"
  }
  
  config = {
    maxMemory = 8
  }
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      providerConfig + resourceConfig,
				ExpectError: regexp.MustCompile("Attribute policy_type.code value must be one of"),
			},
		},
	})
}

// Test validation: invalid associated_resource_type
func TestAccMorpheusPolicyValidationInvalidResourceType(t *testing.T) {
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig := `
resource "hpe_morpheus_policy" "validation_test" {
  name = "` + name + `"
  associated_resource_type = "InvalidType"
  associated_resource_id = 1
  
  policy_type = {
    code = "maxMemory"
  }
  
  config = {
    maxMemory = 8
  }
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      providerConfig + resourceConfig,
				ExpectError: regexp.MustCompile("Attribute associated_resource_type value must be one of"),
			},
		},
	})
}
