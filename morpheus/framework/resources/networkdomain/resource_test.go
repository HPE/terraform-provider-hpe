// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkdomain_test

import (
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/networkdomain"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

// TestAccMorpheusNetworkDomainResourceExampleOk exercises the ported example:
// create, an in-place update, and both import forms, with an empty plan after
// apply (no ExpectNonEmptyPlan).
func TestAccMorpheusNetworkDomainResourceExampleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	// The server lower-cases the domain name; use a lower-case fixture name so
	// there is no perpetual plan diff.
	name := strings.ToLower(acctest.RandomWithPrefix(t.Name()))

	resourceName := "hpe_morpheus_network_domain.example"

	createConfig, err := networkdomain.RenderNetworkDomainConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatalf("failed to render config: %s", err)
	}

	updateConfig, err := networkdomain.RenderNetworkDomainConfig(t, map[string]string{
		"Name":        name,
		"Description": "Updated by acceptance test",
	})
	if err != nil {
		t.Fatalf("failed to render config: %s", err)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + createConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", "Terraform example network domain"),
					resource.TestCheckResourceAttr(resourceName, "public_zone", "false"),
					resource.TestCheckResourceAttr(resourceName, "visibility", "private"),
					resource.TestCheckResourceAttr(resourceName, "active", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "write_only_revision"),
				),
			},
			{Config: providerConfig + createConfig, PlanOnly: true},
			{
				Config: providerConfig + updateConfig,
				Check: resource.TestCheckResourceAttr(
					resourceName, "description", "Updated by acceptance test",
				),
			},
			{Config: providerConfig + updateConfig, PlanOnly: true},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"domain_password_wo",
					"domain_password_wo_version",
					"write_only_revision",
				},
			},
		},
	})
}
