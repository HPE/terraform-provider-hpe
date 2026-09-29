// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package user_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

// Removing first_name or last_name clears it on the user (MORPH-17990). Both
// are Optional+Computed with a NullableStringUpdateModifier, so the plan
// carries the null with nothing else changing, and Update sends it as null.
// (TestAccMorpheusUserResourceRemoveOsUsernamesOk covers the OS usernames,
// which go as empty strings.)
func TestAccMorpheusUserRemovedAttributesAreCleared(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())
	rolesConfig := roleResourceConfig("example", name+"-a")
	const rn = "hpe_morpheus_user.example"

	config := func(names string) string {
		return providerConfig + testhelpers.WhoamiBlock() + rolesConfig + `
resource "hpe_morpheus_user" "example" {
	tenant_id           = data.hpe_morpheus_whoami.current.tenant_id
	username            = "` + name + `"
	email               = "` + name + `@testacc.com"
	password_wo         = "Secret123!"
	password_wo_version = 1
	role_ids            = [hpe_morpheus_role.example.id]
` + names + `
}`
	}

	withNames := `
	first_name = "Ada"
	last_name  = "Lovelace"
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: config(withNames),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(rn, "first_name", "Ada"),
					resource.TestCheckResourceAttr(rn, "last_name", "Lovelace"),
				),
			},
			{
				Config: config(""),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(rn, plancheck.ResourceActionUpdate),
						plancheck.ExpectKnownValue(rn, tfjsonpath.New("first_name"), knownvalue.Null()),
						plancheck.ExpectKnownValue(rn, tfjsonpath.New("last_name"), knownvalue.Null()),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckNoResourceAttr(rn, "first_name"),
					resource.TestCheckNoResourceAttr(rn, "last_name"),
				),
			},
			{
				Config:             config(""),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}
