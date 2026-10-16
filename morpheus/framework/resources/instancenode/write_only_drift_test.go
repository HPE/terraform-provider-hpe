// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instancenode_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

// rotateServerSSHPassword rewrites a server's SSH password directly through the
// API (PUT /api/servers/{id}), reproducing an out-of-band change. Morpheus has
// no importable UpdateHost/UpdateServer operation in the pruned SDK, so the PUT
// is made with a raw authenticated request rather than an SDK call. The
// sshPasswordHash Morpheus renders beside the secret changes as a result, which
// is what the provider's drift detection compares.
func rotateServerSSHPassword(t *testing.T, serverID int64, value string) {
	t.Helper()

	base, ok := testhelpers.LookupProviderEnv("", "url")
	if !ok {
		t.Fatal("TF_VAR_testacc_morpheus_url not set")
	}

	token, ok := testhelpers.LookupProviderEnv("", "access_token")
	if !ok {
		t.Fatal("TF_VAR_testacc_morpheus_access_token not set")
	}

	body, err := json.Marshal(map[string]any{
		"server": map[string]any{"sshPassword": value},
	})
	if err != nil {
		t.Fatalf("rotateServerSSHPassword: marshal: %v", err)
	}

	url := base + "/api/servers/" + strconv.FormatInt(serverID, 10)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPut, url, bytes.NewReader(body))
	if err != nil {
		t.Fatalf("rotateServerSSHPassword: request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("rotateServerSSHPassword: PUT %s: %v", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("rotateServerSSHPassword: PUT %s returned %d: %s", url, resp.StatusCode, string(b))
	}
}

// serverIDFromState reads the int64 "server_id" attribute of a resource from
// stashed state.
func serverIDFromState(t *testing.T, state *terraform.State, address string) int64 {
	t.Helper()

	if state == nil {
		t.Fatal("no state stashed; add testhelpers.StashState to the preceding step")
	}

	rs, ok := state.RootModule().Resources[address]
	if !ok {
		t.Fatalf("resource not found in state: %s", address)
	}

	id, err := strconv.ParseInt(rs.Primary.Attributes["server_id"], 10, 64)
	if err != nil {
		t.Fatalf("resource %s server_id %q is not an integer: %v",
			address, rs.Primary.Attributes["server_id"], err)
	}

	if id == 0 {
		t.Fatalf("resource %s has server_id 0; refusing to PUT /api/servers/0", address)
	}

	return id
}

// serverSSHPasswordHash reads the sshPasswordHash Morpheus renders beside a
// server's stored SSH password, so the test can prove the rotate PUT actually
// changed it (a PUT that silently no-ops is a real finding, not a reason to
// skip).
func serverSSHPasswordHash(t *testing.T, serverID int64) string {
	t.Helper()

	ctx := context.Background()

	client, err := testhelpers.NewClientForServer(ctx, "")
	if err != nil {
		t.Fatalf("serverSSHPasswordHash: client: %v", err)
	}

	idParam := sdk.GetHostIdParameter{Int64: &serverID}

	h, _, err := client.HostsAPI.GetHost(ctx, idParam).Execute()
	if err != nil || h == nil || h.Server == nil {
		t.Fatalf("serverSSHPasswordHash: GetHost %d: %v", serverID, err)
	}

	if v := h.Server.SshPasswordHash.Get(); v != nil {
		return *v
	}

	return ""
}

// TestAccMorpheusInstanceNodeWriteOnlyDriftDetected asserts appliance-side
// drift detection for the write-only ssh_password. After the node is
// provisioned, the server's SSH password is rotated through the API; the next
// plan must detect the drift. Because ssh_password is immutable, the policy is
// forced to a warning (a re-send would be a replacement, which the provider
// must never do silently), so the expected outcome is a clean plan accompanied
// by a drift warning -- not an update. The testing framework does not expose
// warning diagnostics to a plancheck, so this asserts the clean plan and
// documents the warning in a comment.
//
// Needs a provisionable bare-metal instance, like the other node tests, plus an
// ssh_password to detect drift on.
func TestAccMorpheusInstanceNodeWriteOnlyDriftDetected(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.Metal)
	testhelpers.SkipUnlessAcc(t)

	instanceID := os.Getenv("TF_VAR_metal_instance_id")
	poolID := os.Getenv("TF_VAR_metal_pool_id")

	if instanceID == "" || poolID == "" {
		t.Skip("TF_VAR_metal_instance_id and TF_VAR_metal_pool_id not set")
	}

	config := testhelpers.ProviderBlock() + fmt.Sprintf(`
resource "hpe_morpheus_instance_node" "drift" {
  instance_id      = %s
  resource_pool_id = %s
  ssh_password     = "Secret123!"
}
`, instanceID, poolID)

	var stashed *terraform.State

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check:  testhelpers.StashState(&stashed),
			},
			{
				// Rotate the SSH password on the appliance, then plan. Drift is
				// detected and reported as a warning; the plan stays empty
				// because ssh_password is immutable and cannot be re-sent in
				// place.
				PreConfig: func() {
					serverID := serverIDFromState(t, stashed, "hpe_morpheus_instance_node.drift")
					before := serverSSHPasswordHash(t, serverID)
					rotateServerSSHPassword(t, serverID, "Rotated456!")
					if after := serverSSHPasswordHash(t, serverID); after == before {
						t.Fatalf("rotate PUT did not change sshPasswordHash (%q); "+
							"drift cannot be detected", before)
					}
				},
				Config:             config,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}
