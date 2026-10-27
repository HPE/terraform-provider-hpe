// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package testhelpers

import "testing"

// The bare-metal fixture: the objects the HPE Alletra MP Bare Metal (BMaaS)
// datastore test creates against. As with the HVM fixture, everything is named
// and resolved by data sources inside the test configuration, so the
// configuration works on any appliance that has the objects. Nothing here is
// an id. The resource pool has no data source of its own, but the zonePools
// option source lists a cloud's pools, so hpe_morpheus_option_source_value
// resolves it by name.
//
// There are no defaults. Each accessor skips the test, naming the variable and
// what it should point at, when its variable is unset.
const (
	// EnvMetalCloudName names the HPE Bare Metal cloud the datastore is
	// associated with.
	EnvMetalCloudName = "TF_VAR_testacc_morpheus_metal_cloud_name"

	// EnvMetalStorageServerName names the HPE Alletra MP storage server the
	// datastore is backed by (type code "alletra-mp").
	EnvMetalStorageServerName = "TF_VAR_testacc_morpheus_metal_storage_server_name"

	// EnvMetalPoolName names a resource pool of that cloud, as listed by
	// GET /api/clouds/<cloud>/resource-pools -- "Default", the pool an HPE
	// Bare Metal cloud creates when it is initialised, unless another was
	// added. BMaaS datastores must be created in a resource pool of their
	// cloud; Morpheus rejects an id that names no pool with "zonePool.id:
	// Required fields are missing".
	EnvMetalPoolName = "TF_VAR_testacc_morpheus_metal_pool_name"

	// EnvMetalGroupName names a group to grant access to the datastore.
	EnvMetalGroupName = "TF_VAR_testacc_morpheus_metal_group_name"

	// EnvMetalTenantName names a tenant to share the datastore with.
	EnvMetalTenantName = "TF_VAR_testacc_morpheus_metal_tenant_name"
)

// MetalCloudName returns the bare-metal cloud's name, skipping the test when
// EnvMetalCloudName is unset.
func MetalCloudName(t *testing.T) string {
	t.Helper()

	return envOrSkip(t, EnvMetalCloudName, "an HPE Bare Metal cloud")
}

// MetalStorageServerName returns the Alletra MP storage server's name,
// skipping the test when EnvMetalStorageServerName is unset.
func MetalStorageServerName(t *testing.T) string {
	t.Helper()

	return envOrSkip(t, EnvMetalStorageServerName, "an HPE Alletra MP storage server")
}

// MetalPoolName returns the name of a resource pool of the bare-metal cloud,
// skipping the test when EnvMetalPoolName is unset.
func MetalPoolName(t *testing.T) string {
	t.Helper()

	return envOrSkip(t, EnvMetalPoolName, "a resource pool of the HPE Bare Metal cloud")
}

// MetalGroupName returns the name of a group to grant access to, skipping the
// test when EnvMetalGroupName is unset.
func MetalGroupName(t *testing.T) string {
	t.Helper()

	return envOrSkip(t, EnvMetalGroupName, "a group to grant datastore access to")
}

// MetalTenantName returns the name of a tenant to share with, skipping the
// test when EnvMetalTenantName is unset.
func MetalTenantName(t *testing.T) string {
	t.Helper()

	return envOrSkip(t, EnvMetalTenantName, "a tenant to share the datastore with")
}
