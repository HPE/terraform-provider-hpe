// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package cluster_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/cluster"
	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

// rotateHostSSHPassword rewrites a host's SSH password directly through the API
// (PUT /api/servers/{id}), reproducing an out-of-band change. Morpheus has no
// importable UpdateHost/UpdateServer operation in the pruned SDK, so the PUT is
// made with a raw authenticated request. The sshPasswordHash changes as a
// result, which is what the provider's cluster drift detection compares.
func rotateHostSSHPassword(t *testing.T, serverID int64, value string) {
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
		t.Fatalf("rotateHostSSHPassword: marshal: %v", err)
	}

	url := base + "/api/servers/" + strconv.FormatInt(serverID, 10)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPut, url, bytes.NewReader(body))
	if err != nil {
		t.Fatalf("rotateHostSSHPassword: request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("rotateHostSSHPassword: PUT %s: %v", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("rotateHostSSHPassword: PUT %s returned %d: %s", url, resp.StatusCode, string(b))
	}
}

// firstClusterHostServerID returns the server id of the cluster's first host,
// read through the SDK, so the drift test can rotate its SSH password. It
// returns an error rather than calling t.Fatal, so it can be used inside a
// TestStep Check (where a fatal would abort before the cluster is destroyed).
func firstClusterHostServerID(clusterAddress string, state *terraform.State) (int64, error) {
	if state == nil {
		return 0, fmt.Errorf("no state stashed; add testhelpers.StashState to the preceding step")
	}

	rs, ok := state.RootModule().Resources[clusterAddress]
	if !ok {
		return 0, fmt.Errorf("resource not found in state: %s", clusterAddress)
	}

	clusterID, err := strconv.ParseInt(rs.Primary.Attributes["id"], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("cluster id %q is not an integer: %w", rs.Primary.Attributes["id"], err)
	}

	ctx := context.Background()

	client, err := testhelpers.NewClientForServer(ctx, "")
	if err != nil {
		return 0, fmt.Errorf("firstClusterHostServerID: client: %w", err)
	}

	clusterResp, _, err := client.ClustersAPI.GetCluster(ctx, clusterID).Execute()
	if err != nil || clusterResp == nil || clusterResp.Cluster == nil {
		return 0, fmt.Errorf("firstClusterHostServerID: GetCluster %d: %w", clusterID, err)
	}

	for _, srv := range clusterResp.Cluster.Servers {
		if srv.Id != nil && *srv.Id != 0 {
			return *srv.Id, nil
		}
	}

	return 0, fmt.Errorf("cluster %d has no host servers with a non-zero id", clusterID)
}

// clusterHostSSHPasswordHash reads the sshPasswordHash of a host, so the test
// can prove the rotate PUT actually changed it.
func clusterHostSSHPasswordHash(t *testing.T, serverID int64) string {
	t.Helper()

	ctx := context.Background()

	client, err := testhelpers.NewClientForServer(ctx, "")
	if err != nil {
		t.Fatalf("clusterHostSSHPasswordHash: client: %v", err)
	}

	idParam := sdk.GetHostIdParameter{Int64: &serverID}

	h, _, err := client.HostsAPI.GetHost(ctx, idParam).Execute()
	if err != nil || h == nil || h.Server == nil {
		t.Fatalf("clusterHostSSHPasswordHash: GetHost %d: %v", serverID, err)
	}

	if v := h.Server.SshPasswordHash.Get(); v != nil {
		return *v
	}

	return ""
}

// TestAccMorpheusClusterResourceWriteOnlyDriftDetected asserts appliance-side
// drift detection for the write-only server.ssh_password_wo. After the cluster
// is created, one host's SSH password is rotated through the API; the next plan
// must detect the drift. Because server.ssh_password_wo is immutable, the
// policy is forced to a warning (a re-send would replace the cluster, which the
// provider must never do silently), so the expected outcome is a clean plan
// accompanied by a drift warning -- not a replacement. The testing framework
// does not expose warning diagnostics to a plancheck, so this asserts the clean
// plan and documents the warning in a comment.
//
// Gated like the other HVM cluster tests.
func TestAccMorpheusClusterResourceWriteOnlyDriftDetected(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)
	testhelpers.SkipUnlessAcc(t)

	name := acctest.RandomWithPrefix(t.Name())

	dataSourcesConfig := `
data "hpe_morpheus_cloud" "test" {
	name = "hvm"
}

data "hpe_morpheus_group" "test" {
	name = "HVM"
}

data "hpe_morpheus_key_pair" "test" {
	name = "hvmserviceuser-keypair"
}

data "hpe_morpheus_service_plan" "test" {
	name = "Default Manual"
	provision_type_code = "manual"
}
`

	resourceConfig, err := cluster.RenderClusterHvmConfig(t, map[string]string{
		"Name":             name,
		"CloudId":          "data.hpe_morpheus_cloud.test.id",
		"GroupId":          "data.hpe_morpheus_group.test.id",
		"LayoutId":         "449",
		"Label1":           "terraform",
		"Label2":           "test",
		"CreateUser":       "false",
		"DynamicPlacement": "false",
		"CpuArch":          "x86_64",
		"CpuModel":         "host-model",
		"PowerPolicy":      "balanced",
		"ServicePlanId":    "data.hpe_morpheus_service_plan.test.id",
		"SshUsername":      "hvmserviceuser",
		"SshKeyPairId":     "data.hpe_morpheus_key_pair.test.id",
		"SshHost1Name":     name + "-worker-1",
		"SshHost1Ip":       "10.118.67.20",
		"SshHost2Name":     name + "-worker-2",
		"SshHost2Ip":       "10.118.67.21",
		"SshHost3Name":     name + "-worker-3",
		"SshHost3Ip":       "10.118.67.22",
		"Visibility":       "private",
		"Tag1Name":         "source",
		"Tag1Value":        "terraform",
		"Tag2Name":         "environment",
		"Tag2Value":        "test",
	})
	if err != nil {
		t.Fatal(err)
	}

	config := testhelpers.ProviderBlock() + dataSourcesConfig + resourceConfig

	var stashed *terraform.State

	var hostServerID int64

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					testhelpers.StashState(&stashed),
					// Capture the host server id here, in a Check, so a failure
					// to resolve it fails the step without leaving the cluster
					// undestroyed (a t.Fatal in PreConfig would).
					func(state *terraform.State) error {
						id, err := firstClusterHostServerID("hpe_morpheus_cluster.example_hvm", state)
						if err != nil {
							return err
						}

						hostServerID = id

						return nil
					},
				),
			},
			{
				// Rotate one host's SSH password on the appliance, then plan.
				// The host record belongs to the cluster this test creates and
				// is destroyed with it; the PUT changes the credential Morpheus
				// stores for that record, not the physical host, so nothing is
				// left rotated behind.
				// Drift is detected and reported as a warning; the plan stays
				// empty because the attribute is immutable and cannot be
				// re-sent without replacing the cluster. The testing framework
				// does not expose warning diagnostics to a plancheck, so the
				// clean plan is asserted here and the Warn path is covered by
				// TestUnitClusterHostDriftWarnPath.
				PreConfig: func() {
					if hostServerID == 0 {
						t.Fatal("host server id was not captured in step 1")
					}

					before := clusterHostSSHPasswordHash(t, hostServerID)
					rotateHostSSHPassword(t, hostServerID, "Rotated456!")
					if after := clusterHostSSHPasswordHash(t, hostServerID); after == before {
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
