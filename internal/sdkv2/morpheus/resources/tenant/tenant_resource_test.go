// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package tenant_test

import (
	"context"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/internal/framework/provider"
	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus"
	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/resources/role"
	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus"
	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/resources/tenant"
)

func TestMain(m *testing.M) {
	code := m.Run()

	testhelpers.WriteMergedResults()

	os.Exit(code)
}

func newProviderWithError() (tfprotov6.ProviderServer, error) {
	frameworkProvider := provider.New("test", morpheus.New())()

	frameworkServer, err := providerserver.NewProtocol6WithError(frameworkProvider)()
	if err != nil {
		return nil, err
	}

	sdkv2Server, err := tf5to6server.UpgradeServer(context.Background(), sdkv2morpheus.Provider().GRPCProvider)
	if err != nil {
		return nil, err
	}

	providers := []func() tfprotov6.ProviderServer{
		func() tfprotov6.ProviderServer { return frameworkServer },
		func() tfprotov6.ProviderServer { return sdkv2Server },
	}

	return tf6muxserver.NewMuxServer(context.Background(), providers...)
}

var testAccProtoV6ProviderFactories = map[string]func() (
	tfprotov6.ProviderServer, error,
){
	"hpe": newProviderWithError,
}

func TestAccMorpheusTenantExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	dependencyResourceConfig, err := role.RenderRoleTenantConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	resourceConfig, err := tenant.RenderTenantConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_tenant.tf_example_tenant",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_tenant.tf_example_tenant",
			"description",
			"Terraform example tenant",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_tenant.tf_example_tenant",
			"enabled",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_tenant.tf_example_tenant",
			"subdomain",
			"tfexample",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_tenant.tf_example_tenant",
			"currency",
			"USD",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_tenant.tf_example_tenant",
			"account_number",
			"12345",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_tenant.tf_example_tenant",
			"account_name",
			"tenant 12345",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_tenant.tf_example_tenant",
			"customer_number",
			"12345",
		),
		resource.TestCheckResourceAttrSet(
			"hpe_morpheus_tenant.tf_example_tenant",
			"base_role_id",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create dependencies
			{
				Config:             providerConfig + dependencyResourceConfig,
				ExpectNonEmptyPlan: false,
			},
			// Plan
			{
				Config:             providerConfig + dependencyResourceConfig + resourceConfig,
				ExpectNonEmptyPlan: true,
				Check:              checkFn,
				PlanOnly:           true,
			},
			// Apply
			{
				Config: providerConfig + dependencyResourceConfig + resourceConfig,
				Check:  checkFn,
			},
			// Plan after apply
			{
				Config:             providerConfig + dependencyResourceConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
		},
	})
}
