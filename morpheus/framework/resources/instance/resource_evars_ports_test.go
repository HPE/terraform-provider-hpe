// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instance_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/instance"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

// evarsHCL renders evars list elements for the example template.
func evarsHCL(kv ...string) string {
	out := ""
	for i := 0; i+1 < len(kv); i += 2 {
		out += `    {
      name  = "` + kv[i] + `"
      value = "` + kv[i+1] + `"
    },
`
	}

	return out
}

// TestAccMorpheusInstanceResourceEvarsRemove is the MORPH-612 regression guard
// for removing environment variables: an instance created with two evars is
// re-applied with the evars block removed entirely, and the variables must be
// gone from the instance.
//
// Previously the update never wrote evars, so the read reported the two
// unchanged server values and the count stayed at 2. The checks assert the
// count read back from the API after the update.
func TestAccMorpheusInstanceResourceEvarsRemove(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())

	withEvars, err := instance.RenderInstanceConfig(t, map[string]string{
		"Name":  name,
		"Evars": evarsHCL("ENV", "test", "DEBUG", "true"),
	})
	if err != nil {
		t.Fatal(err)
	}

	withoutEvars, err := instance.RenderInstanceConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	inPlace := resource.ConfigPlanChecks{
		PreApply: []plancheck.PlanCheck{
			plancheck.ExpectResourceAction(instanceResourceName, plancheck.ResourceActionUpdate),
		},
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + withEvars,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "evars.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(instanceResourceName, "evars.*",
						map[string]string{"name": "ENV", "value": "test"}),
					resource.TestCheckTypeSetElemNestedAttrs(instanceResourceName, "evars.*",
						map[string]string{"name": "DEBUG", "value": "true"}),
				),
			},
			{
				Config:           providerConfig + withoutEvars,
				ConfigPlanChecks: inPlace,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Removed evars read back as absent -- null set, so no evars.# key.
					resource.TestCheckNoResourceAttr(instanceResourceName, "evars.#"),
				),
			},
			{
				Config:             providerConfig + withoutEvars,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

// TestAccMorpheusInstanceResourceEvarsModifyAndAdd is the modify-and-add
// scenario folded in from MORPH-610: two existing values change and a third
// variable is added in one update. Previously this failed with three
// "inconsistent result after apply" errors because the planned set never
// reached the instance. The checks assert the values read back from the API.
func TestAccMorpheusInstanceResourceEvarsModifyAndAdd(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())

	initial, err := instance.RenderInstanceConfig(t, map[string]string{
		"Name":  name,
		"Evars": evarsHCL("APP_ENV", "development", "LOG_LEVEL", "info"),
	})
	if err != nil {
		t.Fatal(err)
	}

	updated, err := instance.RenderInstanceConfig(t, map[string]string{
		"Name":  name,
		"Evars": evarsHCL("APP_ENV", "production", "LOG_LEVEL", "debug", "API_KEY", "test-key-123"),
	})
	if err != nil {
		t.Fatal(err)
	}

	inPlace := resource.ConfigPlanChecks{
		PreApply: []plancheck.PlanCheck{
			plancheck.ExpectResourceAction(instanceResourceName, plancheck.ResourceActionUpdate),
		},
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + initial,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "evars.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(instanceResourceName, "evars.*",
						map[string]string{"name": "APP_ENV", "value": "development"}),
				),
			},
			{
				Config:           providerConfig + updated,
				ConfigPlanChecks: inPlace,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "evars.#", "3"),
					resource.TestCheckTypeSetElemNestedAttrs(instanceResourceName, "evars.*",
						map[string]string{"name": "APP_ENV", "value": "production"}),
					resource.TestCheckTypeSetElemNestedAttrs(instanceResourceName, "evars.*",
						map[string]string{"name": "LOG_LEVEL", "value": "debug"}),
					resource.TestCheckTypeSetElemNestedAttrs(instanceResourceName, "evars.*",
						map[string]string{"name": "API_KEY", "value": "test-key-123"}),
				),
			},
			{
				Config:             providerConfig + updated,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

// TestAccMorpheusInstanceResourcePortsRequireReplace covers the ports half of
// MORPH-612. Ports are create-only in Morpheus, so removing them cannot be an
// in-place update; previously the change was accepted and silently never
// applied. The plan must now propose replacing the instance, and the
// replacement must complete with the ports gone.
func TestAccMorpheusInstanceResourcePortsRequireReplace(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())

	withPorts, err := instance.RenderInstanceConfig(t, map[string]string{
		"Name": name,
		"Ports": `    {
      name = "http"
      port = 8080
    },
`,
	})
	if err != nil {
		t.Fatal(err)
	}

	withoutPorts, err := instance.RenderInstanceConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + withPorts,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "ports.#", "1"),
				),
			},
			{
				Config: providerConfig + withoutPorts,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(instanceResourceName, plancheck.ResourceActionReplace),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckNoResourceAttr(instanceResourceName, "ports.#"),
				),
			},
		},
	})
}
