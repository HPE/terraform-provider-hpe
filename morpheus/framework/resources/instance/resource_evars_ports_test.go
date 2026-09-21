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

// evarOpt is one environment variable with its display options, for the
// MORPH-17657 tests. Export/Masked are rendered only when non-nil so a test can
// also exercise the "omitted -> default" path.
type evarOpt struct {
	name, value    string
	export, masked *bool
}

func boolp(b bool) *bool { return &b }

func evarsOptHCL(vars ...evarOpt) string {
	out := ""
	for _, v := range vars {
		out += "    {\n      name  = \"" + v.name + "\"\n      value = \"" + v.value + "\"\n"
		if v.export != nil {
			out += "      export = " + map[bool]string{true: "true", false: "false"}[*v.export] + "\n"
		}
		if v.masked != nil {
			out += "      masked = " + map[bool]string{true: "true", false: "false"}[*v.masked] + "\n"
		}
		out += "    },\n"
	}

	return out
}

// TestAccMorpheusInstanceResourceEvarsOptionsOnCreate covers MORPH-17657's
// create path: the create request drops export/masked, so the provider applies
// them with a follow-up add-or-update. A variable created masked and one
// created non-exported must read back with those options set, alongside a
// defaults-only variable, and the steady-state plan must be empty.
func TestAccMorpheusInstanceResourceEvarsOptionsOnCreate(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())

	config, err := instance.RenderInstanceConfig(t, map[string]string{
		"Name": name,
		"Evars": evarsOptHCL(
			evarOpt{name: "SECRET", value: "s3cr3t", masked: boolp(true)},
			evarOpt{name: "LOCAL_ONLY", value: "v", export: boolp(false)},
			evarOpt{name: "PLAIN", value: "p"},
		),
	})
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "evars.#", "3"),
					resource.TestCheckTypeSetElemNestedAttrs(instanceResourceName, "evars.*",
						map[string]string{"name": "SECRET", "value": "s3cr3t", "masked": "true"}),
					resource.TestCheckTypeSetElemNestedAttrs(instanceResourceName, "evars.*",
						map[string]string{"name": "LOCAL_ONLY", "value": "v", "export": "false"}),
					// Omitted options stay unset in state (the appliance applies its
					// defaults); the element is still present with its name and value.
					resource.TestCheckTypeSetElemNestedAttrs(instanceResourceName, "evars.*",
						map[string]string{"name": "PLAIN", "value": "p"}),
				),
			},
			{
				Config:             providerConfig + config,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

// TestAccMorpheusInstanceResourceEvarsOptionsUpdate covers MORPH-17657's update
// path: toggling masked on and export off for existing variables, with the
// values unchanged, must apply in place and read back, with an empty plan after.
func TestAccMorpheusInstanceResourceEvarsOptionsUpdate(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())

	initial, err := instance.RenderInstanceConfig(t, map[string]string{
		"Name": name,
		"Evars": evarsOptHCL(
			evarOpt{name: "TOKEN", value: "t0k3n"},
			evarOpt{name: "REGION", value: "eu"},
		),
	})
	if err != nil {
		t.Fatal(err)
	}

	toggled, err := instance.RenderInstanceConfig(t, map[string]string{
		"Name": name,
		"Evars": evarsOptHCL(
			evarOpt{name: "TOKEN", value: "t0k3n", masked: boolp(true)},
			evarOpt{name: "REGION", value: "eu", export: boolp(false)},
		),
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
					resource.TestCheckTypeSetElemNestedAttrs(instanceResourceName, "evars.*",
						map[string]string{"name": "TOKEN", "value": "t0k3n"}),
				),
			},
			{
				Config:           providerConfig + toggled,
				ConfigPlanChecks: inPlace,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "evars.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(instanceResourceName, "evars.*",
						map[string]string{"name": "TOKEN", "value": "t0k3n", "masked": "true"}),
					resource.TestCheckTypeSetElemNestedAttrs(instanceResourceName, "evars.*",
						map[string]string{"name": "REGION", "value": "eu", "export": "false"}),
				),
			},
			{
				Config:             providerConfig + toggled,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

// TestAccMorpheusInstanceResourceEvarsOnlyUpdateKeepsServicePlanOptions guards
// the interaction between the evars update path and the instance update gate.
//
// service_plan_options is not returned by the API and is carried forward
// through state. An evars-only change is applied through the environment
// variable operations, outside the instance update call, so whichever way the
// gate resolves for that update, the configured service_plan_options must
// survive it and the following plan must be empty. This pins that contract
// against both the update path and the skip path.
func TestAccMorpheusInstanceResourceEvarsOnlyUpdateKeepsServicePlanOptions(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())

	initial, err := instance.RenderInstanceConfig(t, map[string]string{
		"Name":      name,
		"MaxMemory": "2048",
		"Evars":     evarsHCL("STAGE", "one"),
	})
	if err != nil {
		t.Fatal(err)
	}

	// Only the evar value changes; service_plan_options is untouched.
	evarsOnly, err := instance.RenderInstanceConfig(t, map[string]string{
		"Name":      name,
		"MaxMemory": "2048",
		"Evars":     evarsHCL("STAGE", "two"),
	})
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + initial,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "service_plan_options.max_memory", "2048"),
					resource.TestCheckTypeSetElemNestedAttrs(instanceResourceName, "evars.*",
						map[string]string{"name": "STAGE", "value": "one"}),
				),
			},
			{
				Config: providerConfig + evarsOnly,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(instanceResourceName, plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckTypeSetElemNestedAttrs(instanceResourceName, "evars.*",
						map[string]string{"name": "STAGE", "value": "two"}),
					// The carried-forward value must survive the evars-only update.
					resource.TestCheckResourceAttr(instanceResourceName, "service_plan_options.max_memory", "2048"),
				),
			},
			{
				Config:             providerConfig + evarsOnly,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}
