// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

// Package capabilities provides capability-based test activation.
// Tests declare their infrastructure requirements using capabilities,
// and only run if the target system has those capabilities available.
//
// Basic Morpheus appliance access is always assumed (not a declared capability).
// Only declare capabilities for additional infrastructure requirements.
//
// Usage:
//
//	func TestAccMorpheusVMwareCloud(t *testing.T) {
//	    capabilities.MustHaveOrSkip(t, capabilities.VMware)
//	    // ... rest of test
//	}
//
// Set available capabilities via environment variable:
//
//	export TF_ACC_CAPABILITIES="vmware,nsxt,alletra"
//
// If TF_ACC_CAPABILITIES is not set, no capabilities are available and
// no tests will run. This is a safe default since tests can be destructive.
//
// The All capability marks tests that don't target specific shared infrastructure
// (VMware, Alletra, NSXT, NSXV). On shared hardware, exclude these tests by not
// including "all" in TF_ACC_CAPABILITIES.
//
// See README.md in this package for full documentation.
package capabilities

// Capability represents an infrastructure capability required by tests.
// Note: Basic Morpheus appliance access is always assumed and not declared.
type Capability string

const (
	// All marks tests that don't target specific shared infrastructure.
	// Tests are marked with All if they don't target VMware, Alletra, or NSX*.
	// On shared hardware, exclude these by not including "all" in TF_ACC_CAPABILITIES.
	All Capability = "all"

	Settings Capability = "settings"

	// Cloud Types
	VMware    Capability = "vmware"
	AWS       Capability = "aws"
	Azure     Capability = "azure"
	GCP       Capability = "gcp"
	HVM       Capability = "hvm"
	OpenStack Capability = "openstack"
	Hyperv    Capability = "hyperv"

	// Network Integrations
	NSXT Capability = "nsxt"
	NSXV Capability = "nsxv"
	ACI  Capability = "aci"

	// Network Features
	Network             Capability = "network"
	NetworkDHCP         Capability = "network_dhcp"
	NetworkPool         Capability = "network_pool"
	NetworkServer       Capability = "network_server"
	NetworkRouter       Capability = "network_router"
	NetworkFirewall     Capability = "network_firewall"
	NetworkLoadBalancer Capability = "network_loadbalancer"
	Subnet              Capability = "subnet"

	// NetworkLoadBalancerHAProxy marks tests that provision an HAProxy container
	// load balancer. This needs the HAProxy load-balancer instance type/layout
	// (load-balancer-haproxy-1.7) seeded on the appliance and a cloud able to
	// provision the backing container. It is distinct from NetworkLoadBalancer,
	// which only signals general load-balancer API support (e.g. NSX-T, config
	// validation, data sources). Gate HAProxy-provisioning tests on this so they
	// skip on environments lacking the HAProxy layout/infrastructure.
	NetworkLoadBalancerHAProxy Capability = "network_loadbalancer_haproxy"

	// Automation Integrations
	Ansible      Capability = "ansible"
	AnsibleTower Capability = "ansible_tower"
	ARM          Capability = "arm"
	Chef         Capability = "chef"
	Git          Capability = "git"
	LDAP         Capability = "ldap"
	Puppet       Capability = "puppet"
	ServiceNow   Capability = "service_now"
	Task         Capability = "task"
	VRO          Capability = "vro"

	// Container/Orchestration
	Kubernetes Capability = "kubernetes"
	Docker     Capability = "docker"

	// KubernetesCluster marks tests that require a live, healthy Kubernetes
	// cluster on the target appliance -- one that can provision HKS workers and
	// supports cluster namespaces. This is distinct from Kubernetes, which only
	// signals that Kubernetes library artifacts (spec templates, blueprints,
	// cluster layouts/types) can be exercised without a running cluster. Gate
	// namespace and HKS-provisioning tests on this so they skip on environments
	// that advertise "kubernetes" but have no usable cluster.
	KubernetesCluster Capability = "kubernetes_cluster"

	// Storage
	Alletra      Capability = "alletra"
	Backup       Capability = "backup"
	ResourcePool Capability = "resource_pool"

	// VDI
	VDI Capability = "vdi"

	// Licensing
	License Capability = "license"
)

// String returns the string representation of the capability.
func (c Capability) String() string {
	return string(c)
}
