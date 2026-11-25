// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package host_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

// TestAccMorpheusHostWriteOnlySSHPassword is the write-only change-detection
// matrix for ssh_password_wo on an Omega stub host: the value is sent on
// create, re-sent in place when it alone changes, not re-sent when another
// attribute changes, and the plan is empty after each; then the import forms
// (plain, with digest, malformed digest).
func TestAccMorpheusHostWriteOnlySSHPassword(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.BareMetalStub)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	// Not parallel: the request recorder is process-wide.
	cloudID, groupID, poolID := bmaasFixture(t)
	name := acctest.RandomWithPrefix("TestAccMorpheusHost")

	c := testhelpers.WriteOnlyCase{ //nolint:gosec // attribute names, not credentials
		Address:   "hpe_morpheus_host.test",
		Attribute: "ssh_password_wo",
		APIPath:   "/api/servers",
		APIKey:    "server.sshPassword",
		Config: func(value, extra string, unrelated bool) string {
			description := "write-only matrix"
			if unrelated {
				description = "write-only matrix, updated"
			}

			return testhelpers.ProviderBlock() + omegaHostConfigWith(name, cloudID, groupID, poolID, `
  description         = "`+description+`"
  ssh_username        = "tfacc"
  ssh_password_wo     = "`+value+`"`+extra+`
`)
		},
	}

	rec := testhelpers.RecordRequests(t)

	// No marker attribute (a 2.3 resource): the value-only matrix.
	steps := testhelpers.WriteOnlyValueSteps(c, rec, "Secret123!", "Secret456!")
	// Create-only inputs the API does not return, as in the Omega import step.
	steps = append(steps, testhelpers.WriteOnlyImportSteps(c, "Secret456!",
		"config_omega_bmaas", "delete_options", "timeouts")...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps:                    steps,
	})
}
