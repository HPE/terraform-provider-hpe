// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkdomain_test

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

// legacyConfig uses the deprecated plaintext domain_password and the no-op
// auto_join_domain, exactly as a configuration written for the SDKv2 resource.
func legacyConfig(name, password, description string) string {
	return testhelpers.ProviderBlock() + `
resource "hpe_morpheus_network_domain" "legacy" {
  name             = "` + name + `"
  description      = "` + description + `"
  domain_username  = "joiner"
  domain_password  = "` + password + `"
  auto_join_domain = true
}
`
}

// TestAccMorpheusNetworkDomainResourceDeprecatedPasswordOk proves a
// configuration written for the SDKv2 resource still works: domain_password is
// sent on create and — because the API clears it when omitted — on every
// update, including one that changes nothing else; a changed value is sent too.
func TestAccMorpheusNetworkDomainResourceDeprecatedPasswordOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)
	testhelpers.SkipUnlessAcc(t)

	name := strings.ToLower(acctest.RandomWithPrefix(t.Name()))
	rec := testhelpers.RecordRequests(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				PreConfig: rec.Reset,
				Config:    legacyConfig(name, "Legacy123!", "legacy"),
				Check: resource.ComposeAggregateTestCheckFunc(
					expectDomainPasswordValue(rec, "Legacy123!"),
					resource.TestCheckResourceAttr("hpe_morpheus_network_domain.legacy", "domain_password", "Legacy123!"),
					resource.TestCheckResourceAttr("hpe_morpheus_network_domain.legacy", "auto_join_domain", "true"),
				),
			},
			{
				// Unrelated change: the plaintext value is re-sent, as before.
				PreConfig: rec.Reset,
				Config:    legacyConfig(name, "Legacy123!", "legacy changed"),
				Check:     expectDomainPasswordValue(rec, "Legacy123!"),
			},
			{
				PreConfig: rec.Reset,
				Config:    legacyConfig(name, "Legacy456!", "legacy changed"),
				Check:     expectDomainPasswordValue(rec, "Legacy456!"),
			},
			{
				// Switching to the write-only attribute is a plain update.
				PreConfig: rec.Reset,
				Config:    woConfig(name, "Secret789!", "", false),
				Check:     expectDomainPasswordValue(rec, "Secret789!"),
			},
		},
	})
}
