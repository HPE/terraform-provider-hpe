// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package hostprofile_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	hostprofile "github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/hostprofile"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func TestAccMorpheusHostProfileResourceExampleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.HostProfile, capabilities.MasterTenant)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := hostprofile.RenderHostProfileConfig(
		t,
		map[string]string{
			"Name": name,
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_host_profile.example", "name", name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_host_profile.example", "description",
			"An example host profile",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_host_profile.example",
			"config.storage.iscsi.transport", "tcp",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_host_profile.example",
			"config.system.cpu.power_management", "performance",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_host_profile.example",
			"config.system.logging.level", "info",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_host_profile.example",
			"config.baseline.additional_packages", "openvswitch=2.17.0",
		),
		resource.TestCheckResourceAttrSet(
			"hpe_morpheus_host_profile.example", "id",
		),
		resource.TestCheckResourceAttrSet(
			"hpe_morpheus_host_profile.example", "enabled",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(
			t, adapter.NewMorpheus(), nil,
		),
		Steps: []resource.TestStep{
			{
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				Check:              checkFn,
				PlanOnly:           false,
			},
			{
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      "hpe_morpheus_host_profile.example",
			},
		},
	})
}

func TestAccMorpheusHostProfileResourceUpdateOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.HostProfile, capabilities.MasterTenant)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())
	updatedName := name + "-updated"

	createConfig := `
resource "hpe_morpheus_host_profile" "test" {
  name        = "` + name + `"
  description = "Initial description"
}
`

	updateConfig := `
resource "hpe_morpheus_host_profile" "test" {
  name        = "` + updatedName + `"
  description = "Updated description"

  config = {
    system = {
      logging = {
        level = "debug"
      }
    }
  }
}
`

	createChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(
			"hpe_morpheus_host_profile.test", "name", name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_host_profile.test", "description",
			"Initial description",
		),
	)

	// An explicitly configured empty description must round-trip as "" rather
	// than being coerced to null: Create and Update send it as-is, so the
	// read-back has to preserve it or the apply fails with "provider produced
	// inconsistent result". (Coercion of an echoed "" to null only applies when
	// the description is absent from config -- exercised by clearConfig below.)
	emptyDescConfig := `
resource "hpe_morpheus_host_profile" "test" {
  name        = "` + name + `"
  description = ""
}
`

	emptyDescChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(
			"hpe_morpheus_host_profile.test", "name", name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_host_profile.test", "description", "",
		),
	)

	updateChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(
			"hpe_morpheus_host_profile.test", "name", updatedName,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_host_profile.test", "description",
			"Updated description",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_host_profile.test",
			"config.system.logging.level", "debug",
		),
	)

	// Removing the description must actually clear it on the appliance (the
	// PUT ignores an explicit null, so the provider sends an empty string and
	// coerces the echoed "" back to null). A follow-up plan must be empty.
	clearConfig := `
resource "hpe_morpheus_host_profile" "test" {
  name = "` + updatedName + `"

  config = {
    system = {
      logging = {
        level = "debug"
      }
    }
  }
}
`

	clearChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(
			"hpe_morpheus_host_profile.test", "name", updatedName,
		),
		resource.TestCheckNoResourceAttr(
			"hpe_morpheus_host_profile.test", "description",
		),
	)

	// Add a second config leaf alongside the existing one. Every config field
	// is plain-optional, so this proves a leaf can be added without disturbing
	// the sibling that is already set.
	addLeafConfig := `
resource "hpe_morpheus_host_profile" "test" {
  name = "` + updatedName + `"

  config = {
    system = {
      timezone = "UTC"
      logging = {
        level = "debug"
      }
    }
  }
}
`

	addLeafChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(
			"hpe_morpheus_host_profile.test",
			"config.system.timezone", "UTC",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_host_profile.test",
			"config.system.logging.level", "debug",
		),
	)

	// Remove the timezone leaf again. With config fields plain-optional (not
	// computed) the removed leaf is dropped from the plan and the full config
	// map we send fully replaces the stored one, so it is actually cleared --
	// not re-populated from prior state.
	removeLeafChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckNoResourceAttr(
			"hpe_morpheus_host_profile.test",
			"config.system.timezone",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_host_profile.test",
			"config.system.logging.level", "debug",
		),
	)

	// config.system.ntp.servers is a plain-Optional set. A populated set and an
	// explicitly empty one must both round-trip: Create/Update send the set as
	// configured ("servers":[] for the empty case) and the read-back has to
	// return the same known value, or the apply fails with "provider produced
	// inconsistent result". The empty case is the regression: the read-back
	// used to fold an empty list into null. Both steps are followed by a plan
	// that must be empty, which also proves the appliance echoes the empty list
	// rather than dropping the key (an absent key would read back as null and
	// plan a change on every run).
	ntpServersConfig := `
resource "hpe_morpheus_host_profile" "test" {
  name = "` + updatedName + `"

  config = {
    system = {
      ntp = {
        enabled = true
        servers = ["0.pool.ntp.org", "1.pool.ntp.org"]
      }
    }
  }
}
`

	ntpServersChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(
			"hpe_morpheus_host_profile.test",
			"config.system.ntp.enabled", "true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_host_profile.test",
			"config.system.ntp.servers.#", "2",
		),
		resource.TestCheckTypeSetElemAttr(
			"hpe_morpheus_host_profile.test",
			"config.system.ntp.servers.*", "0.pool.ntp.org",
		),
		resource.TestCheckTypeSetElemAttr(
			"hpe_morpheus_host_profile.test",
			"config.system.ntp.servers.*", "1.pool.ntp.org",
		),
	)

	ntpEmptyServersConfig := `
resource "hpe_morpheus_host_profile" "test" {
  name = "` + updatedName + `"

  config = {
    system = {
      ntp = {
        enabled = false
        servers = []
      }
    }
  }
}
`

	ntpEmptyServersChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(
			"hpe_morpheus_host_profile.test",
			"config.system.ntp.enabled", "false",
		),
		// "servers.#" = "0" is a known empty set; a null set has no "#" key.
		resource.TestCheckResourceAttr(
			"hpe_morpheus_host_profile.test",
			"config.system.ntp.servers.#", "0",
		),
	)

	// Drop the entire config block. Update sends an explicit empty config
	// object so the appliance clears the stored document rather than
	// preserving it, and the follow-up plan must be empty.
	noConfigConfig := `
resource "hpe_morpheus_host_profile" "test" {
  name = "` + updatedName + `"
}
`

	noConfigChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckNoResourceAttr(
			"hpe_morpheus_host_profile.test",
			"config.system.logging.level",
		),
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(
			t, adapter.NewMorpheus(), nil,
		),
		Steps: []resource.TestStep{
			{
				Config:   providerConfig + createConfig,
				Check:    createChecks,
				PlanOnly: false,
			},
			{
				Config:             providerConfig + createConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
			{
				Config:   providerConfig + emptyDescConfig,
				Check:    emptyDescChecks,
				PlanOnly: false,
			},
			{
				Config:             providerConfig + emptyDescConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
			{
				Config:   providerConfig + updateConfig,
				Check:    updateChecks,
				PlanOnly: false,
			},
			{
				Config:             providerConfig + updateConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
			{
				Config:   providerConfig + clearConfig,
				Check:    clearChecks,
				PlanOnly: false,
			},
			{
				Config:             providerConfig + clearConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
			{
				Config:   providerConfig + addLeafConfig,
				Check:    addLeafChecks,
				PlanOnly: false,
			},
			{
				Config:             providerConfig + addLeafConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
			{
				// Removing the timezone leaf reuses clearConfig (logging.level
				// only), proving a nested leaf clears rather than sticking.
				Config:   providerConfig + clearConfig,
				Check:    removeLeafChecks,
				PlanOnly: false,
			},
			{
				Config:             providerConfig + clearConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
			{
				Config:   providerConfig + ntpServersConfig,
				Check:    ntpServersChecks,
				PlanOnly: false,
			},
			{
				Config:             providerConfig + ntpServersConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
			{
				Config:   providerConfig + ntpEmptyServersConfig,
				Check:    ntpEmptyServersChecks,
				PlanOnly: false,
			},
			{
				Config:             providerConfig + ntpEmptyServersConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
			{
				Config:   providerConfig + noConfigConfig,
				Check:    noConfigChecks,
				PlanOnly: false,
			},
			{
				Config:             providerConfig + noConfigConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
		},
	})
}
