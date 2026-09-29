// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package tenant_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/tenant"
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

func TestAccMorpheusFindTenantByName(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name, resourceConfig := tenantFixture(t)

	dataSourceConfig, err := tenant.RenderExample(t, "example-name.tf.tmpl", map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	const dsName = "data.hpe_morpheus_tenant.example"

	checks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(dsName, "name", name),
		resource.TestCheckResourceAttr(dsName, "enabled", "true"),
		resource.TestCheckResourceAttr(dsName, "master", "false"),
		resource.TestCheckResourceAttr(dsName, "subdomain", name),
		resource.TestCheckResourceAttr(dsName, "currency", "USD"),
		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_tenant.example", "id", dsName, "id",
		),
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + resourceConfig,
			},
			{
				Config: providerConfig + resourceConfig + dataSourceConfig,
				Check:  checks,
			},
		},
	})
}

func TestAccMorpheusFindTenantById(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name, resourceConfig := tenantFixture(t)

	dataSourceConfig, err := tenant.RenderExample(t, "example-id.tf.tmpl", map[string]string{
		"Id": "hpe_morpheus_tenant.example.id",
	})
	if err != nil {
		t.Fatal(err)
	}

	const dsName = "data.hpe_morpheus_tenant.example"

	checks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(dsName, "name", name),
		resource.TestCheckResourceAttr(dsName, "enabled", "true"),
		resource.TestCheckResourceAttr(dsName, "master", "false"),
		resource.TestCheckResourceAttr(dsName, "subdomain", name),
		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_tenant.example", "id", dsName, "id",
		),
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + resourceConfig,
			},
			{
				Config: providerConfig + resourceConfig + dataSourceConfig,
				Check:  checks,
			},
		},
	})
}

func TestAccMorpheusFindTenantNoSearchAttrs(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	t.Parallel()

	// A real connection is used so the data source Read runs and returns the
	// "no valid search terms" error.
	config := testhelpers.ProviderBlock() + `
      data "hpe_morpheus_tenant" "test" {
      }`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile(tenant.ErrorNoValidSearchTerms),
			},
		},
	})
}

// TestAccMorpheusFindTenantEmptyName mirrors the MORPH-16725 regression guard:
// an empty name must be rejected at plan time by the LengthAtLeast(1) validator
// before any API call, not flow into the by-name lookup as a "not found".
func TestAccMorpheusFindTenantEmptyName(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	t.Parallel()

	config := testhelpers.ProviderBlock() + `
      data "hpe_morpheus_tenant" "test" {
        name = ""
      }`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: config,
				ExpectError: regexp.MustCompile(
					`Attribute name string length must be at least 1, got: 0`,
				),
			},
		},
	})
}

func TestAccMorpheusFindTenantNotFound(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	config := testhelpers.ProviderBlock() + `
      data "hpe_morpheus_tenant" "test" {
        name = "no-such-tenant-______"
      }`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile(tenant.ErrorNoTenantFound),
			},
		},
	})
}
