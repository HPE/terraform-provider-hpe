// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package securitygrouprule_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/securitygrouprule"
	ruleresource "github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/securitygrouprule"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

// securityGroupFixture renders a security group labelled
// hpe_morpheus_security_group.test.
func securityGroupFixture(name string) string {
	return `
resource "hpe_morpheus_security_group" "test" {
  name = "` + name + `"
}
`
}

func TestAccMorpheusFindSecurityGroupRuleByName(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())

	sgConfig := securityGroupFixture(name)

	resourceConfig, err := ruleresource.RenderSecurityGroupRuleConfig(t, map[string]string{
		"SecurityGroupId": "hpe_morpheus_security_group.test.id",
		"Name":            name + "-rule",
	})
	if err != nil {
		t.Fatal(err)
	}

	dataSourceConfig := `
data "hpe_morpheus_security_group_rule" "example" {
  name              = "` + name + `-rule"
  security_group_id = hpe_morpheus_security_group.test.id
  depends_on        = [hpe_morpheus_security_group_rule.example]
}
`

	checkFn := resource.ComposeAggregateTestCheckFunc(ruleChecks()...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + sgConfig + resourceConfig + dataSourceConfig,
				Check:  checkFn,
			},
		},
	})
}

func TestAccMorpheusFindSecurityGroupRuleById(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())

	sgConfig := securityGroupFixture(name)

	resourceConfig, err := ruleresource.RenderSecurityGroupRuleConfig(t, map[string]string{
		"SecurityGroupId": "hpe_morpheus_security_group.test.id",
		"Name":            name + "-rule",
	})
	if err != nil {
		t.Fatal(err)
	}

	dataSourceConfig, err := securitygrouprule.RenderRuleByIdConfig(t, map[string]string{
		"Id":              "hpe_morpheus_security_group_rule.example.id",
		"SecurityGroupId": "hpe_morpheus_security_group.test.id",
	})
	if err != nil {
		t.Fatal(err)
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(ruleChecks()...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + sgConfig + resourceConfig + dataSourceConfig,
				Check:  checkFn,
			},
		},
	})
}

func TestAccMorpheusFindSecurityGroupRuleNotFound(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())

	sgConfig := securityGroupFixture(name)

	dataSourceConfig := `
data "hpe_morpheus_security_group_rule" "example" {
  name              = "nonexistent-rule-name-that-should-not-exist"
  security_group_id = hpe_morpheus_security_group.test.id
}
`

	expected := regexp.MustCompile(`no security group rule found`)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config:      providerConfig + sgConfig + dataSourceConfig,
				ExpectError: expected,
			},
		},
	})
}

func TestAccMorpheusFindSecurityGroupRuleNoSearchAttrs(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	t.Parallel()

	// A real connection is used so the data source Read runs and returns the
	// "no valid search terms" error; with an unconfigured provider the mux
	// provider fails earlier with a connection error and the validation path is
	// never reached.
	config := testhelpers.ProviderBlock() + `
      data "hpe_morpheus_security_group_rule" "test" {
        security_group_id = 1
      }`

	expected := securitygrouprule.ErrorNoValidSearchTerms

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile(expected),
			},
		},
	})
}

func ruleChecks() []resource.TestCheckFunc {
	ds := "data.hpe_morpheus_security_group_rule.example"

	return []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet(ds, "id"),
		resource.TestCheckResourceAttrSet(ds, "security_group_id"),
		resource.TestCheckResourceAttrSet(ds, "name"),
		resource.TestCheckResourceAttrSet(ds, "direction"),
		resource.TestCheckResourceAttrSet(ds, "policy"),
		resource.TestCheckResourceAttrSet(ds, "protocol"),
	}
}
