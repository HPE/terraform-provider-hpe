// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package testhelpers

import (
	"os"
	"testing"
)

// EnvKubernetesClusterID is the environment variable that supplies the ID of a
// namespace-capable Kubernetes cluster for tests gated on the
// kubernetes_cluster capability (cluster namespaces, HKS provisioning).
const EnvKubernetesClusterID = "TF_VAR_testacc_morpheus_k8s_cluster_id"

// EnvAffinityCloudID is the environment variable that supplies the ID of a
// cloud that supports affinity groups, for tests gated on the affinity_group
// capability. Morpheus only seeds hasAffinityGroups for the vmware,
// vmwareCloudAws and macstadium cloud types, so most clouds on an appliance
// cannot serve these tests.
const EnvAffinityCloudID = "TF_VAR_testacc_morpheus_affinity_cloud_id"

// EnvAffinityClusterID is the environment variable that supplies the ID of a
// cluster that supports affinity groups, for tests gated on the affinity_group
// capability. Only the mvm-cluster (HVM) cluster type supports them.
const EnvAffinityClusterID = "TF_VAR_testacc_morpheus_affinity_cluster_id"

// EnvAffinityPoolID is the environment variable that supplies the ID of the
// resource pool to scope cloud affinity groups to. Morpheus requires a pool
// when creating a cloud affinity group, and it must be a resource pool of type
// Cluster because the underlying DRS rule is created on a cluster.
const EnvAffinityPoolID = "TF_VAR_testacc_morpheus_affinity_pool_id"

// EnvComputeServerID is the environment variable that supplies the ID of an
// existing compute server (host) for the compute server data source tests.
const EnvComputeServerID = "TF_VAR_testacc_morpheus_compute_server_id"

// EnvComputeServerName is the environment variable that supplies the name of an
// existing compute server (host) for the compute server data source tests. The
// name has to be unique on the appliance, because the data source refuses to
// resolve an ambiguous name.
const EnvComputeServerName = "TF_VAR_testacc_morpheus_compute_server_name"

// EnvComputeServerCloudID is the environment variable that supplies the ID of a
// cloud with compute servers on it, for the cloud_id filter of the compute
// servers data source.
const EnvComputeServerCloudID = "TF_VAR_testacc_morpheus_compute_server_cloud_id"

// EnvComputeServerInstanceID is the environment variable that supplies the ID
// of an instance that owns at least one compute server, for the instance_id
// filter of the compute servers data source. Most hosts on an appliance are
// not owned by an instance, so this is a separate variable from
// EnvComputeServerID rather than something derivable from it.
const EnvComputeServerInstanceID = "TF_VAR_testacc_morpheus_compute_server_instance_id"

// KubernetesClusterID returns the Kubernetes cluster ID to use for
// namespace/HKS tests, taking TF_VAR_testacc_morpheus_cluster_id when set and
// otherwise falling back to the supplied value. These tests only execute when
// the kubernetes_cluster capability is enabled, so whoever enables it is
// expected to point this at a real cluster on the target appliance.
func KubernetesClusterID(fallback string) string {
	if v := os.Getenv(EnvKubernetesClusterID); v != "" {
		return v
	}

	return fallback
}

// AffinityCloudID returns the cloud ID to run cloud affinity group tests
// against, skipping the test when EnvAffinityCloudID is unset.
//
// There is deliberately no fallback. Affinity groups only exist on a handful of
// cloud types, so a guessed ID such as "1" turns a test that should skip into a
// test that fails against the API on every appliance whose first cloud happens
// to be some other type.
func AffinityCloudID(t *testing.T) string {
	t.Helper()

	return envOrSkip(t, EnvAffinityCloudID,
		"a VMware cloud that supports affinity groups")
}

// AffinityClusterID returns the cluster ID to run cluster affinity group tests
// against, skipping the test when EnvAffinityClusterID is unset. As with
// AffinityCloudID there is no fallback: only HVM clusters support affinity
// groups, so guessing an ID produces a failure rather than a skip.
func AffinityClusterID(t *testing.T) string {
	t.Helper()

	return envOrSkip(t, EnvAffinityClusterID,
		"an HVM cluster that supports affinity groups")
}

// AffinityPoolID returns the resource pool ID to scope cloud affinity group
// tests to, skipping the test when EnvAffinityPoolID is unset. Morpheus
// rejects a cloud affinity group created without a pool, and the pool has to
// be a resource pool of type Cluster on the cloud named by EnvAffinityCloudID,
// so there is no sensible fallback here either.
func AffinityPoolID(t *testing.T) string {
	t.Helper()

	return envOrSkip(t, EnvAffinityPoolID,
		"a resource pool of type Cluster on the affinity group cloud")
}

// ComputeServerID returns the ID of the compute server to look up, skipping the
// test when EnvComputeServerID is unset.
//
// As with the affinity group helpers there is no fallback. Host IDs are not
// predictable — they are allocated per appliance and low IDs are frequently
// absent — so a guessed ID such as "1" turns a test that should skip into a
// test that fails with a 404.
func ComputeServerID(t *testing.T) string {
	t.Helper()

	return envOrSkip(t, EnvComputeServerID, "an existing compute server")
}

// ComputeServerName returns the name of the compute server to look up, skipping
// the test when EnvComputeServerName is unset. There is no fallback for the
// same reason as ComputeServerID: no host name is guaranteed to exist.
func ComputeServerName(t *testing.T) string {
	t.Helper()

	return envOrSkip(t, EnvComputeServerName,
		"an existing compute server with a unique name")
}

// ComputeServerCloudID returns the cloud ID to filter compute servers by,
// skipping the test when EnvComputeServerCloudID is unset.
func ComputeServerCloudID(t *testing.T) string {
	t.Helper()

	return envOrSkip(t, EnvComputeServerCloudID, "a cloud with compute servers")
}

// ComputeServerInstanceID returns the instance ID to filter compute servers by,
// skipping the test when EnvComputeServerInstanceID is unset. It must name an
// instance that owns at least one compute server, otherwise the filter matches
// nothing and the test asserts against an empty result.
func ComputeServerInstanceID(t *testing.T) string {
	t.Helper()

	return envOrSkip(t, EnvComputeServerInstanceID,
		"an instance that owns at least one compute server")
}

// envOrSkip returns the value of the named environment variable, or skips the
// test with a message naming both the variable and the infrastructure it is
// expected to point at.
func envOrSkip(t *testing.T, name, requirement string) string {
	t.Helper()

	value := os.Getenv(name)
	if value == "" {
		t.Skip(name + " not set; skipping test requiring " + requirement)
	}

	return value
}

//nolint:lll
const providerConfig = `
variable "testacc_morpheus_url" {
  default = null
}
variable "testacc_morpheus_username" {
  default = null
}
variable "testacc_morpheus_password" {
  default = null
}
variable "testacc_morpheus_access_token" {
  default = null
}
variable "testacc_morpheus_insecure" {
  default = false
}

provider "hpe" {
        morpheus {
                url = var.testacc_morpheus_url
                access_token    = var.testacc_morpheus_access_token
                username = var.testacc_morpheus_access_token == null ? var.testacc_morpheus_username : null
                password = var.testacc_morpheus_access_token == null ? var.testacc_morpheus_password : null
                insecure = var.testacc_morpheus_insecure
        }
}
`

//nolint:lll
const providerConfigLegacy = `
variable "testacc_morpheus_url" {
  default = null
}
variable "testacc_morpheus_username" {
  default = null
}
variable "testacc_morpheus_password" {
  default = null
}
variable "testacc_morpheus_access_token" {
  default = null
}
variable "testacc_morpheus_insecure" {
  default = false
}

provider "morpheus" {
  url          = var.testacc_morpheus_url
  access_token = var.testacc_morpheus_access_token
  username     = var.testacc_morpheus_access_token == null ? var.testacc_morpheus_username : null
  password     = var.testacc_morpheus_access_token == null ? var.testacc_morpheus_password : null
}
`

//nolint:lll
const providerConfigLegacyProviderBlockOnly = `
provider "morpheus" {
  url          = var.testacc_morpheus_url
  access_token = var.testacc_morpheus_access_token
  username     = var.testacc_morpheus_access_token == null ? var.testacc_morpheus_username : null
  password     = var.testacc_morpheus_access_token == null ? var.testacc_morpheus_password : null
}
`

// ProviderBlock returns a provider block that can be used for acceptance testing
func ProviderBlock() string {
	return providerConfig
}

// ProviderBlockLegacy returns a provider block for the legacy morpheus provider
func ProviderBlockLegacy() string {
	return providerConfigLegacy
}

// ProviderBlockMixed returns a provider block for mixed usage of the new and old providers
func ProviderBlockMixed() string {
	return providerConfig + providerConfigLegacyProviderBlockOnly
}

//nolint:lll
const providerConfigUnitTest = `
provider "hpe" {
  morpheus {
    url          = "http://localhost"
    access_token = "unit-test"
  }
}
`

// ProviderBlockUnitTest returns a provider block with a placeholder URL and
// token so that plan-time validation tests (CustomizeDiff / schema) can run as
// unit tests (IsUnitTest, no TF_ACC) without real credentials. The Morpheus
// client is created lazily, so no connection is made before the validation
// under test fires.
func ProviderBlockUnitTest() string {
	return providerConfigUnitTest
}

// If broker_url is null, it'll default to the cloud broker.
// So don't set testacc_pce_identity_broker_url if we wish to test the cloud broker.
const providerConfigPceIdentity = `
variable "testacc_pce_identity_client_id" {
  default = null
}

variable "testacc_pce_identity_client_secret" {
  default = null
}

variable "testacc_pce_identity_issuer_url" {
  default = null
}

variable "testacc_pce_identity_location" {
  default = null
}

variable "testacc_pce_identity_space" {
  default = null
}

variable "testacc_pce_identity_broker_url" {
  default = null
}

provider "hpe" {
  morpheus {
    pce_identity {
      client_id     = var.testacc_pce_identity_client_id
      client_secret = var.testacc_pce_identity_client_secret
      issuer_url    = var.testacc_pce_identity_issuer_url
      location      = var.testacc_pce_identity_location
      space         = var.testacc_pce_identity_space
      broker_url    = var.testacc_pce_identity_broker_url
    }
  }
}
`

func ProviderBlockPceIdentity() string {
	return providerConfigPceIdentity
}
