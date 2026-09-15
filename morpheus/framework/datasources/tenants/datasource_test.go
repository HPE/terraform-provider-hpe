// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package tenants_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/role"
	restenant "github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/tenant"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

// tenantFixture renders a tenant-scoped base role plus a tenant resource that
// references it. The returned name is the tenant's name.
func tenantFixture(t *testing.T) (string, string) {
	t.Helper()

	name := acctest.RandomWithPrefix(t.Name())

	roleConfig, err := role.RenderRoleTenantConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	tenantConfig, err := restenant.RenderTenantConfig(t, map[string]string{
		"Name":       name,
		"Subdomain":  name,
		"BaseRoleId": "hpe_morpheus_role.example.id",
	})
	if err != nil {
		t.Fatal(err)
	}

	return name, roleConfig + tenantConfig
}

func TestAccMorpheusTenantsUnfiltered(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	config := testhelpers.ProviderBlock() + `
      data "hpe_morpheus_tenants" "all" {
      }`

	const dsName = "data.hpe_morpheus_tenants.all"

	checks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrSet(dsName, "tenants.#"),
		resource.TestCheckResourceAttrSet(dsName, "ids.#"),
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check:  checks,
			},
		},
	})
}

func TestAccMorpheusTenantsFilterByName(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	name, resourceConfig := tenantFixture(t)

	dataSourceConfig := `
      data "hpe_morpheus_tenants" "filtered" {
        filter {
          name   = "name"
          values = ["^` + name + `$"]
        }
      }`

	const dsName = "data.hpe_morpheus_tenants.filtered"

	checks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(dsName, "tenants.#", "1"),
		resource.TestCheckResourceAttr(dsName, "ids.#", "1"),
		resource.TestCheckResourceAttr(dsName, "tenants.0.name", name),
		resource.TestCheckResourceAttr(dsName, "tenants.0.enabled", "true"),
		resource.TestCheckResourceAttr(dsName, "tenants.0.master", "false"),
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: testhelpers.ProviderBlock() + resourceConfig,
			},
			{
				Config: testhelpers.ProviderBlock() + resourceConfig + dataSourceConfig,
				Check:  checks,
			},
		},
	})
}

func TestAccMorpheusTenantsSortDescending(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	config := testhelpers.ProviderBlock() + `
      data "hpe_morpheus_tenants" "desc" {
        sort_ascending = false
      }`

	const dsName = "data.hpe_morpheus_tenants.desc"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check:  resource.TestCheckResourceAttrSet(dsName, "tenants.#"),
			},
		},
	})
}

func TestAccMorpheusTenantsInvalidFilterField(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	t.Parallel()

	config := testhelpers.ProviderBlock() + `
      data "hpe_morpheus_tenants" "bad" {
        filter {
          name   = "not_a_field"
          values = ["x"]
        }
      }`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile(`Attribute filter\[\S*\]\.name value must be one of`),
			},
		},
	})
}
