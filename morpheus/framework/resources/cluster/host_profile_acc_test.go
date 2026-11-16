// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package cluster

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	testsweep "github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/sweep"
)

// hostProfileDefaultCode marks the immutable seeded default host profile, which
// the API refuses to attach to a cluster. The attach/detach test must therefore
// pick an enabled, non-default profile as its target.
const hostProfileDefaultCode = "host-profile.hvm.default"

// TestAccMorpheusClusterHostProfileAttachDetach exercises the cluster
// host-profile reconcile path (applyClusterHostProfile / readClusterHostProfile)
// directly against a live HVM cluster, covering the attach and detach
// acceptance criteria that the resource-level HVM provisioning tests do not
// reach (those are gated behind opt-in HVM provisioning with unreachable
// workers).
//
// The cluster-level attach/detach is a pure database operation in Morpheus: it
// sets the host-profile foreign key on the cluster (and any hypervisor hosts),
// and never runs a compliance SSH diff or remediation against running hosts, so
// it is non-disruptive to workloads on the target cluster.
func TestAccMorpheusClusterHostProfileAttachDetach(t *testing.T) {
	defer testhelpers.RecordResult(t)

	// This white-box test drives the SDK directly rather than through the
	// resource.Test harness, so it does not inherit the harness's TF_ACC gate;
	// apply it explicitly, along with the short-mode skip the sibling tests use.
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance test skipped unless TF_ACC is set")
	}

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	capabilities.MustHaveOrSkip(t, capabilities.HostProfile)

	clusterIDStr := testhelpers.HVMClusterID(t)
	clusterID, err := strconv.ParseInt(clusterIDStr, 10, 64)
	if err != nil {
		t.Fatalf("invalid %s %q: %v", testhelpers.EnvHVMClusterID, clusterIDStr, err)
	}

	ctx := context.Background()

	client, err := testhelpers.NewClientForServer(ctx, "")
	if err != nil {
		t.Fatalf("failed to create test client: %v", err)
	}

	// Capture the current attachment so cleanup can restore it.
	// readClusterHostProfile is unconditional: it performs the GET and reports
	// the actual attached profile (or null when none is attached).
	baseline, diags := readClusterHostProfile(ctx, clusterID, client)
	if diags.HasError() {
		t.Fatalf("reading baseline host profile: %s", diags.Errors())
	}

	// The test requires a cluster with no profile attached: attaching a second,
	// different profile is rejected with a 409. Skip rather than disturb a
	// cluster an operator has already configured.
	if !baseline.IsNull() {
		t.Skipf(
			"cluster %d already has host profile %d attached; "+
				"this test requires a cluster with no profile attached",
			clusterID, baseline.ValueInt64(),
		)
	}

	profileIDs := enabledNonDefaultHostProfileIDs(ctx, t, client, 2)

	// To exercise the A->B switch the appliance needs two enabled, non-default
	// profiles. When it is short and this is the master tenant (the only tenant
	// permitted to create host profiles), self-provision the shortfall so the
	// switch path is covered rather than silently skipped. Each created profile
	// is torn down via t.Cleanup, and carries the sweeper prefix as a backstop.
	if len(profileIDs) < 2 && capabilities.Has(capabilities.MasterTenant) {
		for len(profileIDs) < 2 {
			profileIDs = append(
				profileIDs, createHostProfile(ctx, t, client),
			)
		}
	}

	if len(profileIDs) == 0 {
		t.Skipf(
			"no enabled, non-default host profile on the appliance and not the "+
				"master tenant (cannot create one); the attach/detach test "+
				"requires at least one (default %q is not attachable)",
			hostProfileDefaultCode,
		)
	}

	profileID := profileIDs[0]

	// Restore the baseline (detached) state regardless of where the test exits.
	t.Cleanup(func() {
		if d := applyClusterHostProfile(ctx, clusterID, client, baseline); d.HasError() {
			t.Errorf("cleanup: restoring baseline host profile: %s", d.Errors())
		}
	})

	// A profile-less cluster must read back as null (no perpetual diff for an
	// unmanaged cluster): applyClusterHostProfile(null) is a no-op detach here.
	if d := applyClusterHostProfile(
		ctx, clusterID, client, types.Int64Null(),
	); d.HasError() {
		t.Fatalf("no-op detach on profile-less cluster: %s", d.Errors())
	}
	if got, d := readClusterHostProfile(ctx, clusterID, client); d.HasError() {
		t.Fatalf("reading profile-less cluster: %s", d.Errors())
	} else if !got.IsNull() {
		t.Fatalf("profile-less cluster: want null, got %d", got.ValueInt64())
	}

	// Attach A.
	if d := applyClusterHostProfile(
		ctx, clusterID, client, types.Int64Value(profileID),
	); d.HasError() {
		t.Fatalf("attaching host profile %d: %s", profileID, d.Errors())
	}

	got, diags := readClusterHostProfile(ctx, clusterID, client)
	if diags.HasError() {
		t.Fatalf("reading host profile after attach: %s", diags.Errors())
	}
	if got.IsNull() || got.ValueInt64() != profileID {
		t.Fatalf("after attach: want host profile %d, got %v", profileID, got)
	}

	// Re-attaching the same profile is an idempotent no-op (no 409).
	if d := applyClusterHostProfile(
		ctx, clusterID, client, types.Int64Value(profileID),
	); d.HasError() {
		t.Fatalf("idempotent re-attach of profile %d: %s", profileID, d.Errors())
	}

	// Switch A -> B. The attach endpoint 409s on replace, so this exercises the
	// detach-before-attach path. Only runs when a second distinct profile
	// exists on the appliance.
	if len(profileIDs) > 1 {
		profileB := profileIDs[1]
		if d := applyClusterHostProfile(
			ctx, clusterID, client, types.Int64Value(profileB),
		); d.HasError() {
			t.Fatalf("switching host profile %d -> %d: %s", profileID, profileB, d.Errors())
		}

		got, diags = readClusterHostProfile(ctx, clusterID, client)
		if diags.HasError() {
			t.Fatalf("reading host profile after switch: %s", diags.Errors())
		}
		if got.IsNull() || got.ValueInt64() != profileB {
			t.Fatalf("after switch: want host profile %d, got %v", profileB, got)
		}
	} else {
		t.Logf(
			"only one enabled non-default host profile on the appliance; " +
				"skipping the A->B switch assertion",
		)
	}

	// Removal detaches: a null plan must clear whatever is attached.
	if d := applyClusterHostProfile(
		ctx, clusterID, client, types.Int64Null(),
	); d.HasError() {
		t.Fatalf("detaching host profile: %s", d.Errors())
	}

	got, diags = readClusterHostProfile(ctx, clusterID, client)
	if diags.HasError() {
		t.Fatalf("reading host profile after detach: %s", diags.Errors())
	}
	if !got.IsNull() {
		t.Fatalf("after detach: want no host profile, got %d", got.ValueInt64())
	}
}

// enabledNonDefaultHostProfileIDs returns up to max IDs of enabled, non-default
// host profiles already present on the appliance. It may return fewer than max
// (including zero); the caller decides whether to self-provision the shortfall
// or skip.
func enabledNonDefaultHostProfileIDs(
	ctx context.Context,
	t *testing.T,
	client *sdk.APIClient,
	max int,
) []int64 {
	t.Helper()

	resp, httpResp, err := client.HostProfilesAPI.
		ListHostProfiles(ctx).Max(100).Execute()
	if err != nil {
		t.Fatalf("listing host profiles: %v", err)
	}
	if httpResp == nil || httpResp.StatusCode != 200 {
		t.Fatalf("listing host profiles: unexpected status %d", statusOf(httpResp))
	}

	var ids []int64
	for _, p := range resp.HostProfiles {
		if p.Id == nil || p.Enabled == nil || !*p.Enabled {
			continue
		}
		if code := p.Code.Get(); code != nil && *code == hostProfileDefaultCode {
			continue
		}

		ids = append(ids, *p.Id)
		if len(ids) == max {
			break
		}
	}

	return ids
}

// createHostProfile creates a non-default host profile on the appliance and
// registers its teardown, returning the new profile's id. The Morpheus API
// hard-codes enabled=true on create and ignores any enabled field in the
// request body, so a plain create yields an attachable (enabled) profile; the
// result is read back to guard that always-enabled contract against a future
// API regression. Only valid on the master tenant.
func createHostProfile(
	ctx context.Context,
	t *testing.T,
	client *sdk.APIClient,
) int64 {
	t.Helper()

	name := testsweep.TestResourcePrefix + "ClusterHostProfile" +
		acctest.RandString(8)

	createReq := sdk.CreateHostProfileRequest{
		Name: name,
	}

	createResp, hresp, err := client.HostProfilesAPI.
		CreateHostProfile(ctx).
		CreateHostProfileRequest(createReq).
		Execute()
	if err != nil || hresp == nil ||
		(hresp.StatusCode != http.StatusOK && hresp.StatusCode != http.StatusCreated) {
		t.Fatalf("creating host profile %q: %v (status %d)", name, err, statusOf(hresp))
	}
	if createResp == nil || createResp.Id == nil {
		t.Fatalf("creating host profile %q: response missing id", name)
	}

	id := *createResp.Id

	t.Cleanup(func() {
		if _, _, derr := client.HostProfilesAPI.
			DeleteHostProfile(ctx, id).Execute(); derr != nil {
			t.Errorf("cleanup: deleting host profile %d: %v", id, derr)
		}
	})

	// Guard the API's always-enabled-on-create contract: if a created profile
	// ever came back disabled, the A->B switch this fixture exists to cover
	// would not run.
	got, _, err := client.HostProfilesAPI.GetHostProfile(ctx, id).Execute()
	if err != nil {
		t.Fatalf("reading back created host profile %d: %v", id, err)
	}
	if got.HostProfile == nil || got.HostProfile.Enabled == nil ||
		!*got.HostProfile.Enabled {
		t.Fatalf(
			"created host profile %d is not enabled; the API is expected to "+
				"create every host profile enabled", id,
		)
	}

	return id
}

// statusOf returns the status code of an HTTP response, tolerating a nil
// response so error formatting never panics.
func statusOf(resp *http.Response) int {
	if resp == nil {
		return 0
	}

	return resp.StatusCode
}
