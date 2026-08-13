// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package cloudaffinitygroupmember_test

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

const resourceName = "hpe_morpheus_cloud_affinity_group_member.test"

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

// envSeedServer names an existing compute server to place in the group.
//
// The tests deliberately do not provision an instance. Membership does not
// require one, and provisioning is by far the slowest and most fragile part of
// this appliance — an unrelated IPAM outage would otherwise fail tests that
// have nothing to do with provisioning.
const envSeedServer = "TF_VAR_testacc_morpheus_affinity_member_server_id"

func seedServerID(t *testing.T) string {
	t.Helper()

	v := os.Getenv(envSeedServer)
	if v == "" {
		t.Skip(envSeedServer + " not set; skipping test requiring an existing compute " +
			"server to place in an affinity group")
	}

	return v
}

// stateInt64 pulls an attribute out of the resource's Terraform state.
func stateInt64(s *terraform.State, attr string) (int64, error) {
	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return 0, fmt.Errorf("resource %s not found in state", resourceName)
	}

	return strconv.ParseInt(rs.Primary.Attributes[attr], 10, 64)
}

// checkMemberPresent asserts the server really is in the group, read from the
// API rather than from state.
func checkMemberPresent() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		cloudID, err := stateInt64(s, "cloud_id")
		if err != nil {
			return err
		}

		groupID, err := stateInt64(s, "affinity_group_id")
		if err != nil {
			return err
		}

		serverID, err := stateInt64(s, "server_id")
		if err != nil {
			return err
		}

		//nolint:thelper // TestCheckFunc has no *testing.T to forward.
		members := groupMembersNoT(cloudID, groupID)
		for _, m := range members {
			if m == serverID {
				return nil
			}
		}

		return fmt.Errorf("server %d is not in affinity group %d; members are %v",
			serverID, groupID, members)
	}
}

// groupMembersNoT is groupMembers for use inside a TestCheckFunc, which has no
// *testing.T available.
func groupMembersNoT(cloudID, groupID int64) []int64 {
	client, err := testhelpers.NewClientForServer(context.Background(), "")
	if err != nil {
		return nil
	}

	result, _, err := client.CloudsAPI.
		GetCloudAffinityGroup(context.Background(), cloudID, groupID).Execute()
	if err != nil || result == nil || result.AffinityGroup == nil {
		return nil
	}

	ids := make([]int64, 0, len(result.AffinityGroup.Servers))
	for _, s := range result.AffinityGroup.Servers {
		if s.Id != nil {
			ids = append(ids, *s.Id)
		}
	}

	return ids
}

// config builds a group owned by the test plus one membership in it.
//
// The group is created by the test rather than reused, so that anything left
// behind is removed by the affinity group sweeper. A membership has no name,
// tags or metadata of its own, so a leaked one cannot be identified later —
// which is why there is no membership sweeper.
func config(name, cloudID, serverID, poolID string) string {
	return fmt.Sprintf(`
resource "hpe_morpheus_cloud_affinity_group" "test" {
  cloud_id      = %[2]s
  pool_id       = %[4]s
  name          = "%[1]s"
  affinity_type = "KEEP_TOGETHER"
  active        = true
}

resource "hpe_morpheus_cloud_affinity_group_member" "test" {
  cloud_id        = %[2]s
  affinity_group_id = hpe_morpheus_cloud_affinity_group.test.id
  server_id         = %[3]s
}
`, name, cloudID, serverID, poolID)
}

// TestAccMorpheusCloudAffinityGroupMember covers the lifecycle and import.
func TestAccMorpheusCloudAffinityGroupMember(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.VMware, capabilities.AffinityGroup)

	cloudID := testhelpers.AffinityCloudID(t)
	serverID := seedServerID(t)
	poolID := testhelpers.AffinityPoolID(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + config(name, cloudID, serverID, poolID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "server_id", serverID),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					checkMemberPresent(),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccMorpheusCloudAffinityGroupMemberCoexists is the reason this resource
// exists.
//
// A server is added to the group by something other than this resource, exactly
// as happens when an instance is provisioned into the group. Terraform must
// leave it alone: the previous approach, where the group resource owned the
// whole membership set, removed it on the next apply.
func TestAccMorpheusCloudAffinityGroupMemberCoexists(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.VMware, capabilities.AffinityGroup)

	cloudID := testhelpers.AffinityCloudID(t)
	serverID := seedServerID(t)
	poolID := testhelpers.AffinityPoolID(t)

	otherServer := os.Getenv(envSeedServer + "_2")
	if otherServer == "" {
		t.Skip(envSeedServer + "_2 not set; skipping test requiring a second existing " +
			"compute server to join the group out of band")
	}

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())
	cfg := providerConfig + config(name, cloudID, serverID, poolID)

	otherID, err := strconv.ParseInt(otherServer, 10, 64)
	if err != nil {
		t.Fatalf("%s_2 is not a number: %v", envSeedServer, err)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: cfg,
				Check: resource.ComposeAggregateTestCheckFunc(
					checkMemberPresent(),
					// Join the group behind Terraform's back, standing in for an
					// instance being provisioned into it. Done from a check
					// because this is where the group id is available.
					joinOutOfBand(otherID),
				),
			},
			{
				// Re-apply unchanged. The member added out of band by the
				// previous step's check must still be there afterwards.
				Config: cfg,
				Check: resource.ComposeAggregateTestCheckFunc(
					checkMemberPresent(),
					checkOtherMemberSurvived(otherID),
				),
			},
		},
	})
}

// joinOutOfBand adds a server to the group without going through Terraform,
// standing in for an instance being provisioned into it.
func joinOutOfBand(serverID int64) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		cloudID, err := stateInt64(s, "cloud_id")
		if err != nil {
			return err
		}

		groupID, err := stateInt64(s, "affinity_group_id")
		if err != nil {
			return err
		}

		client, err := testhelpers.NewClientForServer(context.Background(), "")
		if err != nil {
			return fmt.Errorf("failed to create API client: %w", err)
		}

		servers := append(groupMembersNoT(cloudID, groupID), serverID)

		body := sdk.UpdateCloudAffinityGroupRequest{
			AffinityGroup: &sdk.UpdateCloudAffinityGroupRequestAffinityGroup{
				Servers: servers,
			},
		}

		if _, _, err := client.CloudsAPI.
			UpdateCloudAffinityGroup(context.Background(), cloudID, groupID).
			UpdateCloudAffinityGroupRequest(body).Execute(); err != nil {
			return fmt.Errorf("failed to add server %d out of band: %w", serverID, err)
		}

		return nil
	}
}

// checkOtherMemberSurvived asserts Terraform did not evict the out-of-band
// member.
func checkOtherMemberSurvived(otherID int64) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		cloudID, err := stateInt64(s, "cloud_id")
		if err != nil {
			return err
		}

		groupID, err := stateInt64(s, "affinity_group_id")
		if err != nil {
			return err
		}

		members := groupMembersNoT(cloudID, groupID)
		for _, m := range members {
			if m == otherID {
				return nil
			}
		}

		return fmt.Errorf(
			"server %d was added to affinity group %d outside Terraform but is no "+
				"longer a member after apply; members are %v. Membership managed "+
				"elsewhere must not be evicted", otherID, groupID, members)
	}
}
