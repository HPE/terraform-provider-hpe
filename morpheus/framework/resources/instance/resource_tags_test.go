// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instance_test

import (
	"context"
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/instance"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

// The tags Morpheus holds beyond the configured ones -- added by a tag policy,
// a provisioning workflow or the cloud's tag sync -- must not fail the apply or
// show as drift, and an update must not strip them (MORPH-20973). The tests add
// a tag out of band, as such a mechanism would, and then: the plan is empty;
// an update of the managed tags leaves the out-of-band tag in place; an empty
// `tags` is an empty set, not null, and still leaves it in place.

const tagsAddr = "hpe_morpheus_instance.example"

// TestAccMorpheusInstanceTagsVMware runs the scenario on the VMware cloud the
// customer report concerns (vCenter tags round-trip through Morpheus).
func TestAccMorpheusInstanceTagsVMware(t *testing.T) {
	defer testhelpers.RecordResult(t)
	capabilities.MustHaveOrSkip(t, capabilities.VMware)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	name := acctest.RandomWithPrefix("TestAccMorpheusInstance-tags")
	render := func(tags string) string {
		cfg, err := instance.RenderInstanceVMwareConfig(t, map[string]string{"Name": name})
		if err != nil {
			t.Fatal(err)
		}

		return testhelpers.ProviderBlock() + fmt.Sprintf(tagsPolicy, name) + replaceTags(t, cfg, tags)
	}

	runTagsScenario(t, render, true)
}

// TestAccMorpheusInstanceTagsHVM runs the same scenario on an HVM cloud.
func TestAccMorpheusInstanceTagsHVM(t *testing.T) {
	defer testhelpers.RecordResult(t)
	capabilities.MustHaveOrSkip(t, capabilities.HVM)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	name := acctest.RandomWithPrefix("TestAccMorpheusInstance-tags")
	render := func(tags string) string {
		cfg, err := instance.RenderInstanceHVMConfig(t, map[string]string{"Name": name})
		if err != nil {
			t.Fatal(err)
		}

		return testhelpers.ProviderBlock() + fmt.Sprintf(tagsPolicy, name) + replaceTags(t, cfg, tags)
	}

	runTagsScenario(t, render, false)
}

const (
	tagsInitial = `
  tags = [
    { name = "sweepable",  value = "true" },
    { name = "managed_by", value = "terraform" },
    { name = "env",        value = "prod" },
    { name = "costcenter", value = "42" },
  ]`
	// managed_by changed, costcenter dropped.
	tagsUpdated = `
  tags = [
    { name = "sweepable",  value = "true" },
    { name = "managed_by", value = "tf" },
    { name = "env",        value = "prod" },
  ]`
)

// tagsPolicy is a strict Morpheus Tags policy requiring the `env` tag on every
// instance the test user provisions -- the kind of policy a customer appliance
// has. Scoped to the user so nobody else's provisioning is affected.
const tagsPolicy = `
data "hpe_morpheus_whoami" "me" {}

resource "hpe_morpheus_policy" "require_env" {
  name                     = "%s-require-env"
  description              = "acceptance test: instances must carry env"
  associated_resource_type = "User"
  associated_resource_id   = data.hpe_morpheus_whoami.me.id
  enabled                  = true

  policy_type = { code = "tags" }

  config = {
    strict = true
    key    = "env"
  }
}
`

// tagsNoEnv drops the policy-required tag: Morpheus must refuse the removal.
const tagsNoEnv = `
  tags = [
    { name = "sweepable",  value = "true" },
    { name = "managed_by", value = "tf" },
  ]`

// tagsAdopted is what a configuration manages after an import: two of the tags
// Morpheus holds (one re-valued), not the policy's env nor the out-of-band one.
const tagsAdopted = `
  tags = [
    { name = "sweepable",  value = "true" },
    { name = "managed_by", value = "adopted" },
  ]`

// runTagsScenario runs the steps. policyEnforced says whether Morpheus enforces
// the Tags policy on this cloud: its provision type must have canEnforceTags,
// which the VMware, AWS, Azure and Google provision types do and HVM/KVM does
// not -- there the policy exists but never refuses anything.
func runTagsScenario(t *testing.T, render func(tags string) string, policyEnforced bool) {
	t.Helper()

	var instanceID int64

	refusal := []resource.TestStep{
		{
			// The policy-required tag cannot be removed: Morpheus refuses
			// the removeTags, the error is surfaced, nothing is lost. The
			// 400 body carries "must have tag property env" after a large
			// instance dump; the output is wrapped, so the pattern matches a
			// fragment that does not straddle a line break.
			Config:      render(tagsNoEnv),
			ExpectError: regexp.MustCompile(`(?s)error updating instance.*400.*tag property\s+env`),
		},
		{
			// ... and the state still holds what Morpheus holds.
			Config: render(tagsUpdated),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(tagsAddr, "tags.#", "3"),
				hasTag("tags_all", "env", "prod"),
				hasTag("tags_all", "owner", "ops"),
			),
		},
	}
	if !policyEnforced {
		refusal = nil
	}

	steps := []resource.TestStep{
		{
			// Four tags configured; Morpheus may add its own on provision,
			// which must not fail the apply.
			Config: render(tagsInitial),
			Check: resource.ComposeAggregateTestCheckFunc(
				recordID(&instanceID),
				resource.TestCheckResourceAttr(tagsAddr, "tags.#", "4"),
				// The managed set keeps the configured casing exactly,
				// whatever case the cloud reports.
				resource.TestCheckTypeSetElemNestedAttrs(tagsAddr, "tags.*", map[string]string{"name": "env", "value": "prod"}),
				// tags_all is verbatim: a VMware cloud reports "Prod".
				hasTagNamed("tags_all", "env"),
			),
		},
		{
			// A tag policy / workflow adds a tag out of band. Re-applying
			// the same configuration must be a no-op (the plan is empty);
			// the tag is visible in tags_all only.
			PreConfig: func() { addTagOutOfBand(t, instanceID, "owner", "ops") },
			Config:    render(tagsInitial),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(tagsAddr, "tags.#", "4"),
				hasTag("tags_all", "owner", "ops"),
				notTag("tags", "owner"),
			),
		},
		{Config: render(tagsInitial), PlanOnly: true},
		{
			// Managed tags change: only the difference is sent, so the
			// out-of-band tag survives and the dropped one is gone.
			Config: render(tagsUpdated),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(tagsAddr, "tags.#", "3"),
				hasTag("tags", "managed_by", "tf"),
				notTag("tags", "costcenter"),
				hasTag("tags_all", "managed_by", "tf"),
				hasTag("tags_all", "owner", "ops"),
				notTag("tags_all", "costcenter"),
			),
		},
		{Config: render(tagsUpdated), PlanOnly: true},
		// (the policy refusal steps, where enforced)
		{
			// Import: tags is not populated (the resource cannot know which
			// of Morpheus's tags the configuration manages); tags_all is.
			ResourceName:      tagsAddr,
			ImportState:       true,
			ImportStateVerify: true,
			ImportStateVerifyIgnore: []string{
				"tags", "timeouts", "config_vmware", "config_hvm", "network_interfaces",
				"volumes", "task_set_id", "service_plan_options", "ports", "evars", "labels",
			},
		},
		{
			// The import flow proper, part 1: forget the instance without
			// destroying it (a `removed` block), as a user who adopts an
			// existing instance starts with one Terraform does not manage.
			Config: forget(render("")),
		},
		{
			// Part 2: import it into the working state (ImportStatePersist)
			// under a configuration managing two of the tags Morpheus holds.
			Config:             render(tagsAdopted),
			ImportState:        true,
			ImportStatePersist: true,
			ResourceName:       tagsAddr,
			ImportStateIdFunc: func(_ *terraform.State) (string, error) {
				return strconv.FormatInt(instanceID, 10), nil
			},
		},
		{
			// Part 3: the plan is an in-place update, not a replacement;
			// the two become managed (one re-valued through addTags), and
			// the tags not configured -- env, which the policy requires,
			// and the out-of-band owner -- are left exactly as they were.
			Config: render(tagsAdopted),
			ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{
				plancheck.ExpectResourceAction(tagsAddr, plancheck.ResourceActionUpdate),
			}},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(tagsAddr, "tags.#", "2"),
				hasTag("tags", "managed_by", "adopted"),
				hasTag("tags_all", "managed_by", "adopted"),
				hasTag("tags_all", "env", "prod"),
				hasTag("tags_all", "owner", "ops"),
			),
		},
		{Config: render(tagsAdopted), PlanOnly: true},
	}

	// The refusal steps go after the managed-tag update, before the import.
	at := slices.IndexFunc(steps, func(st resource.TestStep) bool { return st.ImportState && st.ImportStateVerify })
	steps = slices.Insert(steps, at, refusal...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps:                    steps,
	})
}

// replaceTags swaps the example's tags block for the given one ("" removes it).
func replaceTags(t *testing.T, cfg, tags string) string {
	t.Helper()

	re := regexp.MustCompile(`(?s)\n  tags = \[.*?\n  \]`)
	if !re.MatchString(cfg) {
		t.Fatal("example configuration has no tags block to replace")
	}

	return re.ReplaceAllLiteralString(cfg, tags)
}

// forget replaces the instance resource in cfg with a `removed` block that
// drops it from state without destroying it; the data sources and the policy
// stay.
func forget(cfg string) string {
	re := regexp.MustCompile(`(?s)\nresource "hpe_morpheus_instance" "example" \{.*?\n\}\n`)

	return re.ReplaceAllLiteralString(cfg, `
removed {
  from = hpe_morpheus_instance.example

  lifecycle {
    destroy = false
  }
}
`)
}

func recordID(id *int64) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[tagsAddr]
		if !ok {
			return fmt.Errorf("%s not in state", tagsAddr)
		}

		v, err := strconv.ParseInt(rs.Primary.ID, 10, 64)
		if err != nil {
			return err
		}

		*id = v

		return nil
	}
}

// addTagOutOfBand does what a tag policy or a workflow task does: adds a tag
// to the instance behind Terraform's back, through the API's addTags.
func addTagOutOfBand(t *testing.T, id int64, name, value string) {
	t.Helper()

	ctx := context.Background()

	client, err := testhelpers.NewClientForServer(ctx, "")
	if err != nil {
		t.Fatal(err)
	}

	req := sdk.UpdateInstanceRequest{Instance: &sdk.UpdateInstanceRequestInstance{
		AddTags: []sdk.UpdateInstanceRequestInstanceAddTagsInner{{Name: &name, Value: &value}},
	}}

	_, hresp, err := client.InstancesAPI.UpdateInstance(ctx, id).UpdateInstanceRequest(req).Execute()
	if err != nil {
		body := ""
		if hresp != nil {
			body = hresp.Status
		}

		t.Fatalf("out-of-band addTags on instance %d: %v %s", id, err, body)
	}
}

// hasTagNamed passes when an element of the set attribute carries the name
// (case-insensitively), whatever its value.
func hasTagNamed(attr, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[tagsAddr]
		if !ok {
			return fmt.Errorf("%s not in state", tagsAddr)
		}

		for k, v := range rs.Primary.Attributes {
			if strings.HasPrefix(k, attr+".") && strings.HasSuffix(k, ".name") && strings.EqualFold(v, name) {
				return nil
			}
		}

		return fmt.Errorf("%s has no tag named %q", attr, name)
	}
}

// hasTag passes when an element of the set attribute carries the name and
// value, compared case-insensitively: a VMware cloud reports a tag in the case
// of the vCenter tag it matched ("owner" as "Owner", "prod" as "Prod"), and
// tags_all is verbatim.
func hasTag(attr, name, value string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[tagsAddr]
		if !ok {
			return fmt.Errorf("%s not in state", tagsAddr)
		}

		for k, v := range rs.Primary.Attributes {
			if !strings.HasPrefix(k, attr+".") || !strings.HasSuffix(k, ".name") || !strings.EqualFold(v, name) {
				continue
			}

			if got := rs.Primary.Attributes[strings.TrimSuffix(k, ".name")+".value"]; strings.EqualFold(got, value) {
				return nil
			}
		}

		return fmt.Errorf("%s has no tag %s=%s (case-insensitively)", attr, name, value)
	}
}

// notTag fails when any element of the set attribute carries the name.
func notTag(attr, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[tagsAddr]
		if !ok {
			return fmt.Errorf("%s not in state", tagsAddr)
		}

		for k, v := range rs.Primary.Attributes {
			if strings.HasPrefix(k, attr+".") && strings.HasSuffix(k, ".name") && strings.EqualFold(v, name) {
				return fmt.Errorf("%s contains a tag named %q (%s)", attr, name, k)
			}
		}

		return nil
	}
}
