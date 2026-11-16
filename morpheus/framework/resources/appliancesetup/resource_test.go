// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package appliancesetup_test

import (
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/auth"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

const setupResourceName = "hpe_morpheus_appliance_setup.test"

// TestAccMorpheusApplianceSetupAdopts: against an appliance that is already
// set up (every test appliance is), the resource adopts it without sending a
// setup request, reports already_initialised, reads the build version with
// the provider's credentials, plans no change afterwards and imports by the
// appliance URL. The fresh-appliance path (the setup request itself) cannot
// be exercised on a shared appliance; the fake-appliance unit tests cover it.
func TestAccMorpheusApplianceSetupAdopts(t *testing.T) {
	defer testhelpers.RecordResult(t)

	testhelpers.SkipUnlessAcc(t)

	applianceURL := os.Getenv("TF_VAR_testacc_morpheus_url")
	if applianceURL == "" {
		t.Skip("TF_VAR_testacc_morpheus_url not set")
	}

	config := testhelpers.ProviderBlock() + `
resource "hpe_morpheus_appliance_setup" "test" {
  appliance_name = "Adopted appliance"
  first_name   = "Adopted"
  last_name    = "Admin"
  account_name = "Adopted"
  username     = "adopted-admin"
  password_wo  = "never-sent-Because-the-appliance-is-set-up-1!"
  email        = "adopted@example.com"
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(setupResourceName, "already_initialised", "true"),
					resource.TestCheckResourceAttr(setupResourceName, "hub_mode", "skip"),
					resource.TestCheckResourceAttrSet(setupResourceName, "build_version"),
					resource.TestCheckResourceAttrSet(setupResourceName, "id"),
					resource.TestCheckResourceAttrPair(setupResourceName, "id", setupResourceName, "appliance_url"),
					resource.TestCheckNoResourceAttr(setupResourceName, "password_wo"),
				),
			},
			{
				Config: config,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
			},
			{
				Config:      strings.Replace(config, `username     = "adopted-admin"`, `username     = "changed-admin"`, 1),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("Setup inputs cannot be changed"),
			},
			{
				Config:       config,
				ResourceName: setupResourceName,
				ImportState:  true,
				ImportStateIdFunc: func(*terraform.State) (string, error) {
					return auth.NormalizeBaseURL(applianceURL), nil
				},
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "id",
				ImportStateVerifyIgnore: []string{
					"account_name", "username", "email", "appliance_name", "first_name", "last_name", "timeouts",
				},
			},
		},
	})
}
