// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package user_test

import (
	"context"
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

// Change detection for the write-only passwords (MORPH-18005): a change to
// password_wo is detected without the deprecated password_wo_version marker,
// the marker keeps working for the 2.2 line, an import leaves a clean plan,
// and a resource created by an earlier provider upgrades with one in-place
// update.

// passwordWorks reports whether username/password authenticate against the
// appliance the tests run against -- the proof that a password was, or was
// not, sent to Morpheus.
func passwordWorks(t *testing.T, username, password string) bool {
	t.Helper()

	base, ok := testhelpers.LookupProviderEnv("", "url")
	if !ok {
		t.Fatalf("%s not set", testhelpers.ProviderEnvName("", "url"))
	}

	form := url.Values{"username": {username}, "password": {password}}
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost,
		strings.TrimRight(base, "/")+"/oauth/token?grant_type=password&scope=write&client_id=morph-api",
		strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	insecure, _ := testhelpers.LookupProviderEnv("", "insecure")
	client := &http.Client{Timeout: 30 * time.Second, Transport: &http.Transport{
		// Test appliances use self-signed certificates.
		TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure == "true"}, //nolint:gosec
	}}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("authenticating as %s: %v", username, err)
	}

	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

// checkPassword asserts, as a TestCheckFunc, that password does (want=true) or
// does not (want=false) authenticate the user.
func checkPassword(t *testing.T, username, password string, want bool) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		if got := passwordWorks(t, username, password); got != want {
			return fmt.Errorf("password %q for %s: authenticates=%v, want %v", password, username, got, want)
		}

		return nil
	}
}

// revisionRecorder remembers write_only_revision across steps so a step can
// assert that it changed, or did not.
type revisionRecorder struct {
	last string
}

func (r *revisionRecorder) record(address string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[address]
		if !ok {
			return fmt.Errorf("%s not in state", address)
		}

		v := rs.Primary.Attributes["write_only_revision"]
		if v == "" {
			return fmt.Errorf("%s: write_only_revision is not set", address)
		}

		r.last = v

		return nil
	}
}

func (r *revisionRecorder) expect(address string, changed bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		prev := r.last

		if err := r.record(address)(s); err != nil {
			return err
		}

		if changed && r.last == prev {
			return fmt.Errorf("%s: write_only_revision did not change (still %s)", address, prev)
		}

		if !changed && r.last != prev {
			return fmt.Errorf("%s: write_only_revision changed from %s to %s", address, prev, r.last)
		}

		return nil
	}
}

func userConfig(name, email, password, extra string) string {
	return testhelpers.ProviderBlock() + testhelpers.WhoamiBlock() + roleResourceConfig("example", name+"-role") + `
resource "hpe_morpheus_user" "wo" {
	tenant_id   = data.hpe_morpheus_whoami.current.tenant_id
	username    = "` + name + `"
	email       = "` + email + `"
	password_wo = "` + password + `"
	role_ids    = [hpe_morpheus_role.example.id]
` + extra + `
}`
}

// TestAccMorpheusUserResourceWriteOnlyDigestOk is the change-detection matrix
// for password_wo on one resource: create; change the value with no marker;
// change something else; bump the marker; remove the marker; and the plan is
// empty after each.
func TestAccMorpheusUserResourceWriteOnlyDigestOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	name := acctest.RandomWithPrefix(t.Name())
	addr := "hpe_morpheus_user.wo"
	rev := &revisionRecorder{}

	update := resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{
		plancheck.ExpectResourceAction(addr, plancheck.ResourceActionUpdate),
	}}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				// 1. create: the digest is recorded, the password is live.
				Config: userConfig(name, name+"@testacc.com", "Secret123!", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					rev.record(addr),
					checkPassword(t, name, "Secret123!", true),
				),
			},
			{
				Config:   userConfig(name, name+"@testacc.com", "Secret123!", ""),
				PlanOnly: true,
			},
			{
				// 2. change the password only, no marker: an update is planned
				// and the new password is live.
				Config:           userConfig(name, name+"@testacc.com", "Secret456!", ""),
				ConfigPlanChecks: update,
				Check: resource.ComposeAggregateTestCheckFunc(
					rev.expect(addr, true),
					checkPassword(t, name, "Secret456!", true),
					checkPassword(t, name, "Secret123!", false),
				),
			},
			{
				Config:   userConfig(name, name+"@testacc.com", "Secret456!", ""),
				PlanOnly: true,
			},
			{
				// 3. change something else: the password is not re-sent.
				Config:           userConfig(name, name+"-2@testacc.com", "Secret456!", ""),
				ConfigPlanChecks: update,
				Check: resource.ComposeAggregateTestCheckFunc(
					rev.expect(addr, false),
					checkPassword(t, name, "Secret456!", true),
				),
			},
			{
				// 4. the deprecated marker still works: adding it re-sends the
				// (unchanged) password once.
				Config:           userConfig(name, name+"-2@testacc.com", "Secret456!", "\tpassword_wo_version = 1\n"),
				ConfigPlanChecks: update,
				Check:            rev.expect(addr, true),
			},
			{
				Config:   userConfig(name, name+"-2@testacc.com", "Secret456!", "\tpassword_wo_version = 1\n"),
				PlanOnly: true,
			},
			{
				// 5. bumping the marker with a new password re-sends it.
				Config:           userConfig(name, name+"-2@testacc.com", "Secret789!", "\tpassword_wo_version = 2\n"),
				ConfigPlanChecks: update,
				Check: resource.ComposeAggregateTestCheckFunc(
					rev.expect(addr, true),
					checkPassword(t, name, "Secret789!", true),
				),
			},
			{
				// 6. removing the marker is one in-place update, then clean.
				Config:           userConfig(name, name+"-2@testacc.com", "Secret789!", ""),
				ConfigPlanChecks: update,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckNoResourceAttr(addr, "password_wo_version"),
					checkPassword(t, name, "Secret789!", true),
				),
			},
			{
				Config:   userConfig(name, name+"-2@testacc.com", "Secret789!", ""),
				PlanOnly: true,
			},
		},
	})
}

// importFixture creates a role and a user directly through the API -- a user
// that exists before Terraform sees it -- and returns their ids.
func importFixture(t *testing.T, name, password string) (userID, roleID int64) {
	t.Helper()

	ctx := context.Background()

	client, err := testhelpers.NewClientForServer(ctx, "")
	if err != nil {
		t.Fatal(err)
	}

	roleResp, _, err := client.RolesAPI.AddRoles(ctx).AddRolesRequest(sdk.AddRolesRequest{
		Role: sdk.AddRolesRequestRole{Authority: name + "-role", RoleType: sdk.PtrString("user")},
	}).Execute()
	if err != nil || roleResp == nil || roleResp.Role == nil || roleResp.Role.Id == nil {
		t.Fatalf("creating fixture role: %v", err)
	}

	roleID = *roleResp.Role.Id
	t.Cleanup(func() { _, _, _ = client.RolesAPI.DeleteRole(ctx, roleID).Execute() })

	userResp, _, err := client.UsersAPI.AddUser(ctx).AddUserRequest(sdk.AddUserRequest{
		User: sdk.AddUserRequestUser{
			Username: name, Email: name + "@testacc.com", Password: password,
			Roles: []sdk.AddUserRequestUserRolesInner{{
				AddUserRequestUserRolesInnerOneOf: &sdk.AddUserRequestUserRolesInnerOneOf{Id: sdk.PtrInt64(roleID)},
			}},
		},
	}).Execute()
	if err != nil || userResp == nil || userResp.User == nil || userResp.User.Id == nil {
		t.Fatalf("creating fixture user: %v", err)
	}

	userID = *userResp.User.Id
	t.Cleanup(func() { _, _, _ = client.UsersAPI.DeleteUser(ctx, userID).Execute() })

	return userID, roleID
}

func importedConfig(importID, name, email, password string, roleID int64) string {
	return testhelpers.ProviderBlock() + testhelpers.WhoamiBlock() + `
import {
	to = hpe_morpheus_user.imported
	id = "` + importID + `"
}

resource "hpe_morpheus_user" "imported" {
	tenant_id   = data.hpe_morpheus_whoami.current.tenant_id
	username    = "` + name + `"
	email       = "` + email + `"
	password_wo = "` + password + `"
	role_ids    = [` + strconv.FormatInt(roleID, 10) + `]
}`
}

func expectAction(addr string, action plancheck.ResourceActionType) resource.ConfigPlanChecks {
	return resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{
		plancheck.ExpectResourceAction(addr, action),
	}}
}

// TestAccMorpheusUserResourceWriteOnlyImportPlainOk: importing with a plain
// id leaves a clean plan, whatever password_wo says. The configured value is
// taken as applied, so a change to it is not seen until the user is next
// updated for another reason -- after which changes are detected.
func TestAccMorpheusUserResourceWriteOnlyImportPlainOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	name := acctest.RandomWithPrefix(t.Name())
	addr := "hpe_morpheus_user.imported"
	userID, roleID := importFixture(t, name, "Original1!")
	id := strconv.FormatInt(userID, 10)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				// Import; the configured password differs from the real one and
				// nothing is sent.
				Config: importedConfig(id, name, name+"@testacc.com", "Configured1!", roleID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(addr, "id", id),
					resource.TestCheckResourceAttrSet(addr, "write_only_revision"),
					checkPassword(t, name, "Original1!", true),
				),
			},
			{
				Config:   importedConfig(id, name, name+"@testacc.com", "Configured1!", roleID),
				PlanOnly: true,
			},
			{
				// The blind spot: a changed password_wo is not detected.
				Config:   importedConfig(id, name, name+"@testacc.com", "Changed2!", roleID),
				PlanOnly: true,
			},
			{
				// An unrelated update records the digests of the configured
				// values (and sends nothing write-only).
				Config:           importedConfig(id, name, name+"-2@testacc.com", "Changed2!", roleID),
				ConfigPlanChecks: expectAction(addr, plancheck.ResourceActionUpdate),
				Check:            checkPassword(t, name, "Original1!", true),
			},
			{
				// From here on, changes are detected.
				Config:           importedConfig(id, name, name+"-2@testacc.com", "Changed3!", roleID),
				ConfigPlanChecks: expectAction(addr, plancheck.ResourceActionUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkPassword(t, name, "Changed3!", true),
					checkPassword(t, name, "Original1!", false),
				),
			},
			{
				Config:   importedConfig(id, name, name+"-2@testacc.com", "Changed3!", roleID),
				PlanOnly: true,
			},
		},
	})
}

// TestAccMorpheusUserResourceWriteOnlyImportWithDigestOk: importing with the
// password's sha256 in the import id ("<id>#password_wo=<hex>") leaves a
// clean plan and detects the very next change to password_wo. This is the
// form a migration tool that knows the value uses.
func TestAccMorpheusUserResourceWriteOnlyImportWithDigestOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	name := acctest.RandomWithPrefix(t.Name())
	addr := "hpe_morpheus_user.imported"
	userID, roleID := importFixture(t, name, "Original1!")
	sum := sha256.Sum256([]byte("Original1!"))
	id := strconv.FormatInt(userID, 10) + "#password_wo=" + hex.EncodeToString(sum[:])

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: importedConfig(id, name, name+"@testacc.com", "Original1!", roleID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(addr, "id", strconv.FormatInt(userID, 10)),
					resource.TestCheckResourceAttrSet(addr, "write_only_revision"),
				),
			},
			{
				Config:   importedConfig(id, name, name+"@testacc.com", "Original1!", roleID),
				PlanOnly: true,
			},
			{
				// The first change after import is detected and applied.
				Config:           importedConfig(id, name, name+"@testacc.com", "Changed2!", roleID),
				ConfigPlanChecks: expectAction(addr, plancheck.ResourceActionUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkPassword(t, name, "Changed2!", true),
					checkPassword(t, name, "Original1!", false),
				),
			},
			{
				Config:   importedConfig(id, name, name+"@testacc.com", "Changed2!", roleID),
				PlanOnly: true,
			},
		},
	})
}

// TestAccMorpheusUserResourceWriteOnlyImportBadDigest: a malformed digest
// fragment is an error, never a silent "not set".
func TestAccMorpheusUserResourceWriteOnlyImportBadDigest(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	name := acctest.RandomWithPrefix(t.Name())
	userID, roleID := importFixture(t, name, "Original1!")
	id := strconv.FormatInt(userID, 10) + "#passwd=abc"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config:      importedConfig(id, name, name+"@testacc.com", "Original1!", roleID),
				ExpectError: regexp.MustCompile(`"passwd" is not a write-only attribute`),
			},
		},
	})
}

// TestAccMorpheusUserResourceWriteOnlyUpgradeOk: a user created by a released
// provider (which has no digests) upgrades to this build with exactly one
// in-place update that records the digests, after which the plan is clean and
// a password change is detected.
func TestAccMorpheusUserResourceWriteOnlyUpgradeOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	name := acctest.RandomWithPrefix(t.Name())
	addr := "hpe_morpheus_user.wo"

	// The released provider has the marker and no write_only_revision (and no
	// hpe_morpheus_whoami data source, so tenant_id is left to default).
	releasedConfig := func(password string) string {
		return testhelpers.ProviderBlock() + roleResourceConfig("example", name+"-role") + `
resource "hpe_morpheus_user" "wo" {
	username            = "` + name + `"
	email               = "` + name + `@testacc.com"
	password_wo         = "` + password + `"
	password_wo_version = 1
	role_ids            = [hpe_morpheus_role.example.id]
}`
	}
	released := releasedConfig("Secret123!")

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"hpe": {Source: "HPE/hpe", VersionConstraint: "2.0.0"},
				},
				Config: released,
				Check:  checkPassword(t, name, "Secret123!", true),
			},
			{
				// Same configuration, this build: one in-place update, which
				// only records the digests and sets write_only_revision.
				ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
				Config:                   released,
				ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(addr, plancheck.ResourceActionUpdate),
				}},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(addr, "write_only_revision"),
					checkPassword(t, name, "Secret123!", true),
				),
			},
			{
				ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
				Config:                   released,
				PlanOnly:                 true,
			},
			{
				// A password change with the marker left alone is detected.
				ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
				Config:                   releasedConfig("Secret456!"),
				ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(addr, plancheck.ResourceActionUpdate),
				}},
				Check: checkPassword(t, name, "Secret456!", true),
			},
		},
	})
}
