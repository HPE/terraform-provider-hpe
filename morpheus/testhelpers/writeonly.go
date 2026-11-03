// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package testhelpers

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// WriteOnlyCase describes one resource with one write-only attribute for
// WriteOnlyMatrix. Config renders the resource (and whatever it depends on)
// for a given write-only value and an "extra" HCL fragment (the deprecated
// marker line, or ""); UnrelatedChange renders it with the value and some
// other attribute changed, used to show the value is not re-sent.
type WriteOnlyCase struct {
	// Address is the resource address, e.g. "hpe_morpheus_backup_host.example".
	Address string
	// Attribute is the write-only attribute, e.g. "ssh_password_wo".
	Attribute string
	// Marker is the deprecated version marker, e.g. "ssh_password_wo_version".
	Marker string
	// APIPath is a substring of the API path the resource's writes go to,
	// e.g. "/api/backups"; APIKey the dotted key of the value in an update
	// request body, e.g. "backup.targetPassword"; CreateAPIKey the key in the
	// create body when it differs (default APIKey).
	APIPath, APIKey, CreateAPIKey string
	// Config renders the configuration. extra is appended inside the resource
	// block (the marker line, or ""); unrelated toggles some other attribute.
	Config func(value, extra string, unrelated bool) string
	// ImportIDFunc produces the import ID from state; default: the resource's
	// id. Composite IDs set it.
	ImportIDFunc func(s *terraform.State) (string, error)
	// ImportDigest, when set, is the digest to put in the import ID for the
	// value passed to WriteOnlyImportSteps -- for a non-string attribute it is
	// sha256 of the value's JSON form; default sha256 of the string itself.
	ImportDigest func(value string) string
	// ValueFromState, when set, resolves the value to digest for the import
	// step from the state at that point, for a write-only value that the
	// configuration takes from another resource (a reference such as
	// hpe_morpheus_instance.source.id) rather than a literal. The value
	// passed to WriteOnlyImportSteps is then the reference, which the digest
	// must not be taken of.
	ValueFromState func(s *terraform.State) (string, error)
}

// Sha256Hex is the digest an importer puts in the import ID for value.
func Sha256Hex(value string) string {
	sum := sha256.Sum256([]byte(value))

	return hex.EncodeToString(sum[:])
}

// WriteOnlyMatrixSteps is the change-detection matrix for a write-only
// attribute, as test steps:
//
//  1. create with the value; the value is sent; then the plan is empty
//  2. change the value only, no marker: an update is planned and the new
//     value is sent; then the plan is empty
//  3. change something else: an update, and the value is not sent
//  4. add the deprecated marker: an update that re-sends the value once
//  5. bump the marker with a new value: the value is sent
//  6. remove the marker: one update; then the plan is empty
//
// The value sent to the API is observed in-process with RequestRecorder, so
// the matrix does not need an appliance-side oracle.
func WriteOnlyMatrixSteps(c WriteOnlyCase, rec *RequestRecorder, v1, v2, v3 string) []resource.TestStep {
	update := resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{
		plancheck.ExpectResourceAction(c.Address, plancheck.ResourceActionUpdate),
	}}
	// The marker line goes wherever the case's Config places extra -- inside
	// the nested block for a nested attribute -- so only the leaf name is used.
	leaf := c.Marker
	if i := strings.LastIndex(leaf, "."); i >= 0 {
		leaf = leaf[i+1:]
	}

	marker := func(n int) string { return "\n  " + leaf + " = " + strconv.Itoa(n) + "\n" }
	revision := &revisionRecorder{address: c.Address}

	createKey := c.CreateAPIKey
	if createKey == "" {
		createKey = c.APIKey
	}

	return []resource.TestStep{
		{
			PreConfig: rec.Reset,
			Config:    c.Config(v1, "", false),
			Check: resource.ComposeAggregateTestCheckFunc(
				revision.record(),
				rec.ExpectSent(c.APIPath, createKey, true),
			),
		},
		{Config: c.Config(v1, "", false), PlanOnly: true},
		{
			PreConfig:        rec.Reset,
			Config:           c.Config(v2, "", false),
			ConfigPlanChecks: update,
			Check: resource.ComposeAggregateTestCheckFunc(
				revision.expect(true),
				rec.ExpectSent(c.APIPath, c.APIKey, true),
			),
		},
		{Config: c.Config(v2, "", false), PlanOnly: true},
		{
			PreConfig:        rec.Reset,
			Config:           c.Config(v2, "", true),
			ConfigPlanChecks: update,
			Check: resource.ComposeAggregateTestCheckFunc(
				revision.expect(false),
				rec.ExpectSent(c.APIPath, c.APIKey, false),
			),
		},
		{
			PreConfig:        rec.Reset,
			Config:           c.Config(v2, marker(1), true),
			ConfigPlanChecks: update,
			Check: resource.ComposeAggregateTestCheckFunc(
				revision.expect(true),
				rec.ExpectSent(c.APIPath, c.APIKey, true),
			),
		},
		{Config: c.Config(v2, marker(1), true), PlanOnly: true},
		{
			PreConfig:        rec.Reset,
			Config:           c.Config(v3, marker(2), true),
			ConfigPlanChecks: update,
			Check: resource.ComposeAggregateTestCheckFunc(
				revision.expect(true),
				rec.ExpectSent(c.APIPath, c.APIKey, true),
			),
		},
		{
			Config:           c.Config(v3, "", true),
			ConfigPlanChecks: update,
			Check:            resource.TestCheckNoResourceAttr(c.Address, c.Marker),
		},
		{Config: c.Config(v3, "", true), PlanOnly: true},
	}
}

// WriteOnlyValueSteps is WriteOnlyMatrixSteps for a resource with no marker
// attribute: create; change the value only (an update that re-sends it);
// change something else (the value is not re-sent); with an empty plan after
// each.
func WriteOnlyValueSteps(c WriteOnlyCase, rec *RequestRecorder, v1, v2 string) []resource.TestStep {
	update := resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{
		plancheck.ExpectResourceAction(c.Address, plancheck.ResourceActionUpdate),
	}}
	revision := &revisionRecorder{address: c.Address}

	createKey := c.CreateAPIKey
	if createKey == "" {
		createKey = c.APIKey
	}

	return []resource.TestStep{
		{
			PreConfig: rec.Reset,
			Config:    c.Config(v1, "", false),
			Check: resource.ComposeAggregateTestCheckFunc(
				revision.record(),
				rec.ExpectSent(c.APIPath, createKey, true),
			),
		},
		{Config: c.Config(v1, "", false), PlanOnly: true},
		{
			PreConfig:        rec.Reset,
			Config:           c.Config(v2, "", false),
			ConfigPlanChecks: update,
			Check: resource.ComposeAggregateTestCheckFunc(
				revision.expect(true),
				rec.ExpectSent(c.APIPath, c.APIKey, true),
			),
		},
		{Config: c.Config(v2, "", false), PlanOnly: true},
		{
			PreConfig:        rec.Reset,
			Config:           c.Config(v2, "", true),
			ConfigPlanChecks: update,
			Check: resource.ComposeAggregateTestCheckFunc(
				revision.expect(false),
				rec.ExpectSent(c.APIPath, c.APIKey, false),
			),
		},
		{Config: c.Config(v2, "", true), PlanOnly: true},
	}
}

// WriteOnlyImportSteps checks the two import forms on a resource created by
// the steps before them (they follow WriteOnlyMatrixSteps or a create step,
// with value being the resource's current write-only value):
//
//  1. import by plain id (ImportStateVerify, write_only_revision ignored):
//     the plan after import is clean
//  2. import with "<id>#attr=<sha256>": likewise clean, and the framework's
//     import verification passes
//  3. a malformed digest fragment is an error
func WriteOnlyImportSteps(c WriteOnlyCase, value string, ignore ...string) []resource.TestStep {
	ignore = append(ignore, "write_only_revision", c.Attribute)
	if c.Marker != "" {
		ignore = append(ignore, c.Marker)
	}

	idFunc := c.ImportIDFunc
	if idFunc == nil {
		idFunc = func(s *terraform.State) (string, error) {
			rs, ok := s.RootModule().Resources[c.Address]
			if !ok {
				return "", fmt.Errorf("%s not in state", c.Address)
			}

			return rs.Primary.ID, nil
		}
	}

	withDigest := func(s *terraform.State) (string, error) {
		id, err := idFunc(s)
		if err != nil {
			return "", err
		}

		digested := value
		if c.ValueFromState != nil {
			digested, err = c.ValueFromState(s)
			if err != nil {
				return "", err
			}
		}

		digest := Sha256Hex(digested)
		if c.ImportDigest != nil {
			digest = c.ImportDigest(digested)
		}

		return id + "#" + c.Attribute + "=" + digest, nil
	}

	bad := func(s *terraform.State) (string, error) {
		id, err := idFunc(s)
		if err != nil {
			return "", err
		}

		return id + "#not_an_attribute=abc", nil
	}

	return []resource.TestStep{
		{
			ResourceName:            c.Address,
			ImportState:             true,
			ImportStateIdFunc:       idFunc,
			ImportStateVerify:       true,
			ImportStateVerifyIgnore: ignore,
		},
		{
			ResourceName:            c.Address,
			ImportState:             true,
			ImportStateIdFunc:       withDigest,
			ImportStateVerify:       true,
			ImportStateVerifyIgnore: ignore,
		},
		{
			ResourceName:      c.Address,
			ImportState:       true,
			ImportStateIdFunc: bad,
			ExpectError:       regexp.MustCompile(`(?s)not_an_attribute.*is\s+not\s+a\s+write-only\s+attribute`),
		},
	}
}

// revisionRecorder remembers write_only_revision across steps.
type revisionRecorder struct {
	address string
	last    string
}

func (r *revisionRecorder) record() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[r.address]
		if !ok {
			return fmt.Errorf("%s not in state", r.address)
		}

		v := rs.Primary.Attributes["write_only_revision"]
		if v == "" {
			return fmt.Errorf("%s: write_only_revision is not set", r.address)
		}

		r.last = v

		return nil
	}
}

func (r *revisionRecorder) expect(changed bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		prev := r.last

		if err := r.record()(s); err != nil {
			return err
		}

		if changed && r.last == prev {
			return fmt.Errorf("%s: write_only_revision did not change (still %s)", r.address, prev)
		}

		if !changed && r.last != prev {
			return fmt.Errorf("%s: write_only_revision changed from %s to %s", r.address, prev, r.last)
		}

		return nil
	}
}

// SkipUnlessAcc is a convenience for write-only matrix tests.
func SkipUnlessAcc(t *testing.T) {
	t.Helper()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}
}

// WriteOnlyReplaceSteps is the change-detection matrix for a write-only
// attribute that Morpheus takes only at creation, so that a change must plan a
// replacement:
//
//  1. create with the value; then the plan is empty
//  2. change something else: an in-place update, and the value is not re-sent
//  3. change the value only, no marker: a replacement is planned; then the
//     plan is empty
//  4. bump the deprecated marker: a replacement
//  5. remove the marker: one in-place update (bookkeeping), then clean
//
// v1 and v2 are two values for the attribute, passed to the case's Config as
// value. Step 2 needs an attribute the resource updates in place; a case whose
// Config ignores unrelated has none, and the step is left out.
func WriteOnlyReplaceSteps(c WriteOnlyCase, rec *RequestRecorder, v1, v2 string) []resource.TestStep {
	action := func(a plancheck.ResourceActionType) resource.ConfigPlanChecks {
		return resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(c.Address, a)}}
	}

	leaf := c.Marker
	if i := strings.LastIndex(leaf, "."); i >= 0 {
		leaf = leaf[i+1:]
	}

	marker := func(n int) string { return "\n  " + leaf + " = " + strconv.Itoa(n) + "\n" }

	createKey := c.CreateAPIKey
	if createKey == "" {
		createKey = c.APIKey
	}

	steps := []resource.TestStep{
		{
			PreConfig: rec.Reset,
			Config:    c.Config(v1, "", false),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet(c.Address, "write_only_revision"),
				rec.ExpectSent(c.APIPath, createKey, true),
			),
		},
		{Config: c.Config(v1, "", false), PlanOnly: true},
	}

	if c.Config(v1, "", true) != c.Config(v1, "", false) {
		steps = append(steps, resource.TestStep{
			PreConfig:        rec.Reset,
			Config:           c.Config(v1, "", true),
			ConfigPlanChecks: action(plancheck.ResourceActionUpdate),
			Check:            rec.ExpectSent(c.APIPath, c.APIKey, false),
		})
	}

	return append(steps, []resource.TestStep{
		{
			PreConfig:        rec.Reset,
			Config:           c.Config(v2, "", true),
			ConfigPlanChecks: action(plancheck.ResourceActionReplace),
			Check:            rec.ExpectSent(c.APIPath, createKey, true),
		},
		{Config: c.Config(v2, "", true), PlanOnly: true},
		{
			Config:           c.Config(v2, marker(1), true),
			ConfigPlanChecks: action(plancheck.ResourceActionReplace),
		},
		{
			Config:           c.Config(v2, "", true),
			ConfigPlanChecks: action(plancheck.ResourceActionUpdate),
			Check:            resource.TestCheckNoResourceAttr(c.Address, c.Marker),
		},
		{Config: c.Config(v2, "", true), PlanOnly: true},
	}...)
}
