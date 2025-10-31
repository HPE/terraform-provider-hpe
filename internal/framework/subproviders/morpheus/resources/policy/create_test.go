package policy_test

import (
	"testing"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test creating a policy with required attributes only
func TestAccMorpheusPolicyRequiredAttrsOk(t *testing.T) {
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

resource "hpe_morpheus_policy" "required" {
  name = "` + name + `"
  associated_resource_type = "Group"
  associated_resource_id = hpe_morpheus_group.test.id
  
  policy_type = {
    code = "maxMemory"
  }
  
  config = {
    maxMemory = 1073741824
  }
}
`

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_policy.required",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_policy.required",
			"associated_resource_type",
			"Group",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_policy.required",
			"policy_type.code",
			"maxMemory",
		),
		// Check computed defaults
		resource.TestCheckResourceAttr(
			"hpe_morpheus_policy.required",
			"enabled",
			"true",
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
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"config"},
				ResourceName:            "hpe_morpheus_policy.required",
				Check:                   checkFn,
			},
		},
	})
}

// Test creating policies with different policy types
func TestAccMorpheusPolicyDifferentTypesOk(t *testing.T) {
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	nameMaxMemory := acctest.RandomWithPrefix(t.Name() + "-maxmem")
	nameMaxCores := acctest.RandomWithPrefix(t.Name() + "-maxcores")
	nameMaxStorage := acctest.RandomWithPrefix(t.Name() + "-maxstorage")
	groupName := acctest.RandomWithPrefix(t.Name() + "-group")

	resourceConfig := `
resource "hpe_morpheus_group" "test" {
  name = "` + groupName + `"
  location = "test"
}

resource "hpe_morpheus_policy" "max_memory" {
  name = "` + nameMaxMemory + `"
  description = "Max memory policy"
  associated_resource_type = "Group"
  associated_resource_id = hpe_morpheus_group.test.id
  enabled = true
  
  policy_type = {
    code = "maxMemory"
  }
  
  config = {
    maxMemory = 1073741824
  }
}

resource "hpe_morpheus_policy" "max_cores" {
  name = "` + nameMaxCores + `"
  description = "Max cores policy"
  associated_resource_type = "Group"
  associated_resource_id = hpe_morpheus_group.test.id
  enabled = true
  
  policy_type = {
    code = "maxCores"
  }
  
  config = {
    maxCores = 4
  }
}

resource "hpe_morpheus_policy" "max_storage" {
  name = "` + nameMaxStorage + `"
  description = "Max storage policy"
  associated_resource_type = "Group"
  associated_resource_id = hpe_morpheus_group.test.id
  enabled = true
  
  policy_type = {
    code = "maxStorage"
  }
  
  config = {
    maxStorage = 10737418240
  }
}
`

	checksMaxMemory := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr("hpe_morpheus_policy.max_memory", "name", nameMaxMemory),
		resource.TestCheckResourceAttr("hpe_morpheus_policy.max_memory", "policy_type.code", "maxMemory"),
		resource.TestCheckResourceAttr("hpe_morpheus_policy.max_memory", "description", "Max memory policy"),
	}

	checksMaxCores := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr("hpe_morpheus_policy.max_cores", "name", nameMaxCores),
		resource.TestCheckResourceAttr("hpe_morpheus_policy.max_cores", "policy_type.code", "maxCores"),
		resource.TestCheckResourceAttr("hpe_morpheus_policy.max_cores", "description", "Max cores policy"),
	}

	checksMaxStorage := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr("hpe_morpheus_policy.max_storage", "name", nameMaxStorage),
		resource.TestCheckResourceAttr("hpe_morpheus_policy.max_storage", "policy_type.code", "maxStorage"),
		resource.TestCheckResourceAttr("hpe_morpheus_policy.max_storage", "description", "Max storage policy"),
	}

	allChecks := append(checksMaxMemory, checksMaxCores...)
	allChecks = append(allChecks, checksMaxStorage...)
	checkFn := resource.ComposeAggregateTestCheckFunc(allChecks...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				Check:              checkFn,
				PlanOnly:           false,
			},
		},
	})
}

// Test creating policies scoped to different resource types (Group, Cloud, User, Role)
func TestAccMorpheusPolicyResourceTypesOk(t *testing.T) {
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	// Create dependency resources
	groupName := acctest.RandomWithPrefix(t.Name() + "-group")
	cloudName := acctest.RandomWithPrefix(t.Name() + "-cloud")
	roleName := acctest.RandomWithPrefix(t.Name() + "-role")
	userName := acctest.RandomWithPrefix(t.Name() + "-user")

	dependencyConfig := `
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
`

	policyNameGroup := acctest.RandomWithPrefix(t.Name() + "-group-policy")
	policyNameCloud := acctest.RandomWithPrefix(t.Name() + "-cloud-policy")
	policyNameRole := acctest.RandomWithPrefix(t.Name() + "-role-policy")
	policyNameUser := acctest.RandomWithPrefix(t.Name() + "-user-policy")

	resourceConfig, err := testhelpers.RenderExample(t, "example.tf.tmpl",
		"ResourceName", "group_policy",
		"Name", policyNameGroup,
		"Description", "Example group-scoped policy",
		"AssociatedResourceType", "Group",
		"AssociatedResourceID", "hpe_morpheus_group.test.id",
		"PolicyTypeCode", "maxMemory",
		"ConfigKey", "maxMemory",
		"ConfigValue", "1073741824",
	)
	if err != nil {
		t.Fatal(err)
	}

	cloudResourceConfig, err := testhelpers.RenderExample(t, "example.tf.tmpl",
		"ResourceName", "cloud_policy",
		"Name", policyNameCloud,
		"Description", "Example cloud-scoped policy",
		"AssociatedResourceType", "Cloud",
		"AssociatedResourceID", "hpe_morpheus_cloud.test.id",
		"PolicyTypeCode", "maxMemory",
		"ConfigKey", "maxMemory",
		"ConfigValue", "1073741824",
	)
	if err != nil {
		t.Fatal(err)
	}

	roleResourceConfig, err := testhelpers.RenderExample(t, "example.tf.tmpl",
		"ResourceName", "role_policy",
		"Name", policyNameRole,
		"Description", "Example role-scoped policy",
		"AssociatedResourceType", "Role",
		"AssociatedResourceID", "hpe_morpheus_role.test.id",
		"PolicyTypeCode", "maxMemory",
		"ConfigKey", "maxMemory",
		"ConfigValue", "1073741824",
	)
	if err != nil {
		t.Fatal(err)
	}

	userResourceConfig, err := testhelpers.RenderExample(t, "example.tf.tmpl",
		"ResourceName", "user_policy",
		"Name", policyNameUser,
		"Description", "Example user-scoped policy",
		"AssociatedResourceType", "User",
		"AssociatedResourceID", "hpe_morpheus_user.test.id",
		"PolicyTypeCode", "maxMemory",
		"ConfigKey", "maxMemory",
		"ConfigValue", "1073741824",
	)
	if err != nil {
		t.Fatal(err)
	}

	checksGroup := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr("hpe_morpheus_policy.group_policy", "name", policyNameGroup),
		resource.TestCheckResourceAttr("hpe_morpheus_policy.group_policy", "associated_resource_type", "Group"),
		resource.TestCheckResourceAttrSet("hpe_morpheus_policy.group_policy", "associated_resource_id"),
	}

	checksCloud := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr("hpe_morpheus_policy.cloud_policy", "name", policyNameCloud),
		resource.TestCheckResourceAttr("hpe_morpheus_policy.cloud_policy", "associated_resource_type", "Cloud"),
		resource.TestCheckResourceAttrSet("hpe_morpheus_policy.cloud_policy", "associated_resource_id"),
	}

	checksRole := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr("hpe_morpheus_policy.role_policy", "name", policyNameRole),
		resource.TestCheckResourceAttr("hpe_morpheus_policy.role_policy", "associated_resource_type", "Role"),
		resource.TestCheckResourceAttrSet("hpe_morpheus_policy.role_policy", "associated_resource_id"),
	}

	checksUser := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr("hpe_morpheus_policy.user_policy", "name", policyNameUser),
		resource.TestCheckResourceAttr("hpe_morpheus_policy.user_policy", "associated_resource_type", "User"),
		resource.TestCheckResourceAttrSet("hpe_morpheus_policy.user_policy", "associated_resource_id"),
	}

	allChecks := append(checksGroup, checksCloud...)
	allChecks = append(allChecks, checksRole...)
	allChecks = append(allChecks, checksUser...)
	checkFn := resource.ComposeAggregateTestCheckFunc(
		allChecks...,
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:             providerConfig + dependencyConfig + resourceConfig + cloudResourceConfig + roleResourceConfig + userResourceConfig,
				ExpectNonEmptyPlan: false,
				Check:              checkFn,
				PlanOnly:           false,
			},
		},
	})
}
