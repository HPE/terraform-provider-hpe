// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

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

// Test validation: incompatible policy type and resource type (motd does not support User)
func TestAccMorpheusPolicyValidationIncompatiblePolicyAndResourceType(t *testing.T) {
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlockMixed()
	name := acctest.RandomWithPrefix(t.Name())
	roleName := acctest.RandomWithPrefix(t.Name() + "-role")
	userName := acctest.RandomWithPrefix(t.Name() + "-user")

	resourceConfig := `
resource "hpe_morpheus_role" "test" {
  name = "` + roleName + `"
  role_type = "user"
}

resource "hpe_morpheus_user" "test" {
  username = "` + userName + `"
  email = "` + userName + `@test.com"
  role_ids = [hpe_morpheus_role.test.id]
  password_wo = "TestPassword123!"
}

resource "hpe_morpheus_policy" "validation_test" {
  name = "` + name + `"
  associated_resource_type = "User"
  associated_resource_id = hpe_morpheus_user.test.id
  
  policy_type = {
    code = "motd"
  }
  
  config = {
    enabled = true
    message = "Test message"
    fullPage = false
    title = "MOTD"
  }
}
`

	resource.Test(t, resource.TestCase{
		ExternalProviders: map[string]resource.ExternalProvider{
			"morpheus": {
				Source:            "gomorpheus/morpheus",
				VersionConstraint: "0.13.2",
			},
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      providerConfig + resourceConfig,
				ExpectError: regexp.MustCompile("Incompatible policy type and resource type"),
			},
		},
	})
}

// Test validation: tenants not supported for policy type
func TestAccMorpheusPolicyValidationTenantsNotSupported(t *testing.T) {
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())
	groupName := acctest.RandomWithPrefix(t.Name() + "-group")
	cloudName := acctest.RandomWithPrefix(t.Name() + "-cloud")

	resourceConfig := `
resource "hpe_morpheus_group" "test" {
  name = "` + groupName + `"
  location = "test"
}

resource "hpe_morpheus_cloud" "test" {
  name = "` + cloudName + `"
  tenant_id = 1
  group_id = hpe_morpheus_group.test.id
  code = "` + cloudName + `"
  cloud_type_code = "standard"
  
  config = {
    certificateProvider = "internal"
    enableNetworkTypeSelection = false
  }
}

resource "hpe_morpheus_policy" "validation_test" {
  name = "` + name + `"
  associated_resource_type = "Cloud"
  associated_resource_id = hpe_morpheus_cloud.test.id
  tenants = [1]
  
  policy_type = {
    code = "requiredNetwork"
  }
  
  config = {
    requiredNetworks = [1]
  }
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      providerConfig + resourceConfig,
				ExpectError: regexp.MustCompile("Tenants not supported for this policy type"),
			},
		},
	})
}
