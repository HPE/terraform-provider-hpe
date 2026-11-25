// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package testhelpers

import (
	"os"
	"testing"
)

// The HPE Morpheus SDN (Juniper CN2) surfaces are built on an HVM cluster whose
// cloud points at an SDN network server, and on that network server's floating
// IP pool. Of that, only the network server can be created from Terraform,
// and only a lab appliance carries the SDN plugin and a controller, so every
// SDN acceptance test builds on pre-existing infrastructure read from the
// environment.
//
// Reading the identifiers from the environment keeps the tests
// appliance-agnostic, and skipping when they are absent makes the missing
// coverage explicit instead of reporting a false failure.
const (
	// EnvSdnCloudID is the HVM cloud whose network service is the SDN network
	// server; SDN networks, routers and load balancers are created in it.
	EnvSdnCloudID = "TF_ACC_SDN_CLOUD_ID"
	// EnvSdnNetworkServerID is the SDN (Juniper CN2) network server id.
	EnvSdnNetworkServerID = "TF_ACC_SDN_NETWORK_SERVER_ID"
	// EnvSdnClusterPoolID is the HVM cluster's resource pool id -- the
	// zonePool.id an SDN network, router or load balancer is scoped to.
	EnvSdnClusterPoolID = "TF_ACC_SDN_CLUSTER_POOL_ID"
	// EnvSdnFloatingIpPoolID is a floating IP pool on the SDN network server,
	// from which floating IPs are allocated.
	EnvSdnFloatingIpPoolID = "TF_ACC_SDN_FLOATING_IP_POOL_ID"
	// EnvSdnLoadBalancer opts in to the SDN load balancer acceptance test,
	// held because any load balancer delete on Morpheus 9.1.x damages the
	// other load balancers' records (MORPH-18174).
	EnvSdnLoadBalancer = "TF_ACC_SDN_LOAD_BALANCER"
	// EnvSdnGroupID is the group SDN routers are created in. Optional; defaults
	// to "1".
	EnvSdnGroupID = "TF_ACC_SDN_GROUP_ID"
)

// SdnFixture identifies pre-existing HPE Morpheus SDN infrastructure that
// acceptance tests build on. All values are rendered into HCL as-is, so they
// are strings.
type SdnFixture struct {
	// CloudID is the HVM cloud whose network service is the SDN network server.
	CloudID string
	// NetworkServerID is the SDN (Juniper CN2) network server id.
	NetworkServerID string
	// ClusterPoolID is the HVM cluster resource pool id (zonePool.id) that SDN
	// resources are scoped to.
	ClusterPoolID string
	// FloatingIpPoolID is a floating IP pool on the SDN network server.
	FloatingIpPoolID string
	// GroupID is the group SDN routers are created in; defaults to "1".
	GroupID string
}

// RequireSdnFixture returns the SDN fixture for the appliance under test,
// skipping the test when any of the four required variables is unset.
//
// It skips rather than fails: an appliance without the SDN plugin is a
// legitimate target for the rest of the suite, and turning that into a failure
// would train people to ignore red results.
func RequireSdnFixture(t *testing.T) SdnFixture {
	t.Helper()

	fixture := SdnFixture{
		CloudID:          os.Getenv(EnvSdnCloudID),
		NetworkServerID:  os.Getenv(EnvSdnNetworkServerID),
		ClusterPoolID:    os.Getenv(EnvSdnClusterPoolID),
		FloatingIpPoolID: os.Getenv(EnvSdnFloatingIpPoolID),
		GroupID:          os.Getenv(EnvSdnGroupID),
	}

	if fixture.GroupID == "" {
		fixture.GroupID = "1"
	}

	var missing []string

	for _, v := range []struct {
		name  string
		value string
	}{
		{EnvSdnCloudID, fixture.CloudID},
		{EnvSdnNetworkServerID, fixture.NetworkServerID},
		{EnvSdnClusterPoolID, fixture.ClusterPoolID},
		{EnvSdnFloatingIpPoolID, fixture.FloatingIpPoolID},
	} {
		if v.value == "" {
			missing = append(missing, v.name)
		}
	}

	if len(missing) > 0 {
		msg := "SDN fixture not configured; set:"
		for _, name := range missing {
			msg += " " + name
		}

		t.Skip(msg)
	}

	return fixture
}
