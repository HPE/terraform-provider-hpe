// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package testhelpers

import (
	"os"
	"testing"
)

// NSX-T acceptance tests need a realized tier-0 gateway, an edge cluster, and a
// group and network server to attach a tier-1 to. None of that can be created
// from Terraform, so it must already exist on the appliance under test.
//
// These identifiers were previously hard-coded (tier-0 router 28, group 3,
// network server 5, a named edge cluster). Those values exist on no appliance
// the tests are currently run against, so every NSX-T test failed at the first
// data source with "Router not found with id 28" -- a failure that looks like a
// provider bug but is only a stale fixture. Reading them from the environment
// keeps the tests appliance-agnostic, and skipping when they are absent makes
// the missing coverage explicit instead of reporting a false failure.
const (
	EnvNsxtTier0RouterID   = "TF_ACC_NSXT_TIER0_ROUTER_ID"
	EnvNsxtGroupID         = "TF_ACC_NSXT_GROUP_ID"
	EnvNsxtNetworkServerID = "TF_ACC_NSXT_NETWORK_SERVER_ID"
	EnvNsxtEdgeCluster     = "TF_ACC_NSXT_EDGE_CLUSTER"
)

// NsxtFixture identifies pre-existing NSX-T infrastructure that acceptance
// tests build on. All values are rendered into HCL as-is, so they are strings.
type NsxtFixture struct {
	// Tier0RouterID is the id of a realized tier-0 gateway, read via the
	// hpe_morpheus_network_router data source to obtain its provider_id.
	Tier0RouterID string
	// GroupID is the group the tier-1 gateway is created in.
	GroupID string
	// NetworkServerID is the NSX-T network server (integration) id.
	NetworkServerID string
	// EdgeCluster is the NSX-T edge cluster external id. A tier-1 needs one
	// before a load balancer service can deploy on it.
	EdgeCluster string
}

// RequireNsxtFixture returns the NSX-T fixture for the appliance under test,
// skipping the test when it has not been configured.
//
// It skips rather than fails: an appliance without NSX-T is a legitimate
// target for the rest of the suite, and turning that into a failure would
// train people to ignore red results.
func RequireNsxtFixture(t *testing.T) NsxtFixture {
	t.Helper()

	fixture := NsxtFixture{
		Tier0RouterID:   os.Getenv(EnvNsxtTier0RouterID),
		GroupID:         os.Getenv(EnvNsxtGroupID),
		NetworkServerID: os.Getenv(EnvNsxtNetworkServerID),
		EdgeCluster:     os.Getenv(EnvNsxtEdgeCluster),
	}

	var missing []string

	for _, v := range []struct {
		name  string
		value string
	}{
		{EnvNsxtTier0RouterID, fixture.Tier0RouterID},
		{EnvNsxtGroupID, fixture.GroupID},
		{EnvNsxtNetworkServerID, fixture.NetworkServerID},
		{EnvNsxtEdgeCluster, fixture.EdgeCluster},
	} {
		if v.value == "" {
			missing = append(missing, v.name)
		}
	}

	if len(missing) > 0 {
		msg := "NSX-T fixture not configured; set:"
		for _, name := range missing {
			msg += " " + name
		}

		t.Skip(msg)
	}

	return fixture
}
