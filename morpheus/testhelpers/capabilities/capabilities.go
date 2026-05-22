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
//	    if capabilities.Missing(t, capabilities.VMware) {
//	        return
//	    }
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

	// Cloud Types
	VMware    Capability = "vmware"
	AWS       Capability = "aws"
	Azure     Capability = "azure"
	GCP       Capability = "gcp"
	OpenStack Capability = "openstack"
	Hyperv    Capability = "hyperv"

	// Network Integrations
	NSXT Capability = "nsxt"
	NSXV Capability = "nsxv"
	ACI  Capability = "aci"

	// Network Features
	NetworkDHCP         Capability = "network_dhcp"
	NetworkRouter       Capability = "network_router"
	NetworkFirewall     Capability = "network_firewall"
	NetworkLoadBalancer Capability = "network_loadbalancer"

	// Automation Integrations
	Ansible      Capability = "ansible"
	AnsibleTower Capability = "ansible_tower"
	Chef         Capability = "chef"
	Puppet       Capability = "puppet"

	// Container/Orchestration
	Kubernetes Capability = "kubernetes"
	Docker     Capability = "docker"

	// Storage
	Alletra Capability = "alletra"

	// VDI
	VDI Capability = "vdi"
)

// String returns the string representation of the capability.
func (c Capability) String() string {
	return string(c)
}
