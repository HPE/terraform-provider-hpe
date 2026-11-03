// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkdomain_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

const (
	woAPIPath = "/api/networks/domains"
	woAPIKey  = "networkDomain.domainPassword"
	woMask    = "************"
)

// woConfig renders the write-only test resource. extra is appended inside the
// block (the deprecated marker line, or ""); unrelated toggles the description.
func woConfig(name, value, extra string, unrelated bool) string {
	description := "wo test"
	if unrelated {
		description = "wo test changed"
	}

	return testhelpers.ProviderBlock() + `
resource "hpe_morpheus_network_domain" "wo" {
  name               = "` + name + `"
  description        = "` + description + `"
  visibility         = "private"
  active             = true
  domain_username    = "joiner"
  domain_password_wo = "` + value + `"` + extra + `
}
`
}

// expectDomainPasswordValue asserts that, since the last Reset, a PUT/POST to
// the domains path sent networkDomain.domainPassword equal to want. This is
// what makes the mask rule observable: an unrelated update must send the mask,
// never the secret and never nothing (omitting it would clear the password).
func expectDomainPasswordValue(rec *testhelpers.RequestRecorder, want string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		writes := rec.Writes(woAPIPath)

		for _, w := range writes {
			nd, ok := w.Body["networkDomain"].(map[string]any)
			if !ok {
				continue
			}

			got, ok := nd["domainPassword"].(string)
			if !ok {
				return fmt.Errorf("domainPassword missing from a write to %q (omitting it clears the password)", woAPIPath)
			}

			if got != want {
				return fmt.Errorf("domainPassword sent %q, want %q", got, want)
			}

			return nil
		}

		return fmt.Errorf("no write to %q carried networkDomain.domainPassword", woAPIPath)
	}
}

// TestAccMorpheusNetworkDomainResourceWriteOnlyOk: changes to
// domain_password_wo are detected by digest; the deprecated marker still works;
// an unrelated update sends the mask (never the secret, never nothing); both
// import forms leave a clean plan.
func TestAccMorpheusNetworkDomainResourceWriteOnlyOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)
	testhelpers.SkipUnlessAcc(t)

	name := strings.ToLower(acctest.RandomWithPrefix(t.Name()))
	addr := "hpe_morpheus_network_domain.wo"

	const (
		v1 = "Secret123!"
		v2 = "Secret456!"
	)

	rec := testhelpers.RecordRequests(t)

	update := resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{
		plancheck.ExpectResourceAction(addr, plancheck.ResourceActionUpdate),
	}}

	steps := []resource.TestStep{
		// 1. Create: the configured value is sent.
		{
			PreConfig: rec.Reset,
			Config:    woConfig(name, v1, "", false),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet(addr, "write_only_revision"),
				expectDomainPasswordValue(rec, v1),
			),
		},
		{Config: woConfig(name, v1, "", false), PlanOnly: true},
		// 2. Change the value only: an update that re-sends the new value.
		{
			PreConfig:        rec.Reset,
			Config:           woConfig(name, v2, "", false),
			ConfigPlanChecks: update,
			Check:            expectDomainPasswordValue(rec, v2),
		},
		{Config: woConfig(name, v2, "", false), PlanOnly: true},
		// 3. Change something else: the MASK is sent, not the secret and not
		//    nothing. This is the critical difference from the generic matrix.
		{
			PreConfig:        rec.Reset,
			Config:           woConfig(name, v2, "", true),
			ConfigPlanChecks: update,
			Check:            expectDomainPasswordValue(rec, woMask),
		},
		{Config: woConfig(name, v2, "", true), PlanOnly: true},
		// 4. Bump the deprecated marker: the value is re-sent once.
		{
			PreConfig:        rec.Reset,
			Config:           woConfig(name, v2, "\n  domain_password_wo_version = 1\n", true),
			ConfigPlanChecks: update,
			Check:            expectDomainPasswordValue(rec, v2),
		},
		{Config: woConfig(name, v2, "\n  domain_password_wo_version = 1\n", true), PlanOnly: true},
		// 5. Remove the marker: one bookkeeping update, then clean.
		{
			Config:           woConfig(name, v2, "", true),
			ConfigPlanChecks: update,
			Check:            resource.TestCheckNoResourceAttr(addr, "domain_password_wo_version"),
		},
		{Config: woConfig(name, v2, "", true), PlanOnly: true},
	}

	// Import forms, on the resource created above.
	c := testhelpers.WriteOnlyCase{
		Address:   addr,
		Attribute: "domain_password_wo",
		Marker:    "domain_password_wo_version",
		APIPath:   woAPIPath,
		APIKey:    woAPIKey,
		Config: func(value, extra string, unrelated bool) string {
			return woConfig(name, value, extra, unrelated)
		},
	}
	steps = append(steps, testhelpers.WriteOnlyImportSteps(c, v2)...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps:                    steps,
	})
}
