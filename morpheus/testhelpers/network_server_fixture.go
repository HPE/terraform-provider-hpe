// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package testhelpers

import (
	"fmt"
	"testing"
)

// The network server fixture: the network servers (network-service
// integrations) the `network_server` tests read and create networks on, named
// by environment variables and resolved by data sources inside each test
// configuration. Nothing here is an id.
//
// Two kinds of network server are involved. A visible one -- an SDN or NSX-T
// integration, listed by GET /api/networks/servers -- is what the network
// server data-source tests look up. An HVM cluster's own network server is
// *hidden*: the list API returns it only with showHidden=true, so the typed
// HVM network tests, which need its id, can only resolve it once the data
// source can ask for hidden servers (MORPH-20056); until then they skip for
// want of EnvHVMNetworkServerName. There are no defaults: a guessed name turns
// a test that should skip into one that fails against the API.
const (
	// EnvNetworkServerName names a network server the appliance lists, for
	// the network server data-source tests.
	EnvNetworkServerName = "TF_VAR_testacc_morpheus_network_server_name"
	// EnvHVMNetworkServerName names the HVM cluster's network server
	// ("<cluster> HVM Network Server"), which every typed HVM network is
	// created on.
	EnvHVMNetworkServerName = "TF_VAR_testacc_morpheus_hvm_network_server_name"
	// EnvHVMSwitchName and EnvHVMSwitchName2 name two virtual switches of the
	// HVM cluster, as the hvm-vswitch-or-router option source lists them
	// (vs0, vs1, ...); an OVS port group is created on the first and moved
	// to the second, which replaces it.
	EnvHVMSwitchName  = "TF_VAR_testacc_morpheus_hvm_switch_name"
	EnvHVMSwitchName2 = "TF_VAR_testacc_morpheus_hvm_switch_name_2"
	// EnvHVMHostInterface names a host network interface the HVM hosts have
	// (bond0, ens1f1, ...), for an HVM data network.
	EnvHVMHostInterface = "TF_VAR_testacc_morpheus_hvm_host_interface"
)

// NetworkServerName returns the name of a listed network server, skipping the
// test when EnvNetworkServerName is unset.
func NetworkServerName(t *testing.T) string {
	t.Helper()

	return envOrSkip(t, EnvNetworkServerName, "a network server (network-service integration)")
}

// HVMNetworkServerName returns the name of the HVM cluster's network server,
// skipping the test when EnvHVMNetworkServerName is unset.
func HVMNetworkServerName(t *testing.T) string {
	t.Helper()

	return envOrSkip(t, EnvHVMNetworkServerName, "the HVM cluster's network server")
}

// HVMSwitchNames returns the names of two virtual switches of the HVM cluster,
// skipping the test when either variable is unset.
func HVMSwitchNames(t *testing.T) (first, second string) {
	t.Helper()

	first = envOrSkip(t, EnvHVMSwitchName, "a virtual switch of the HVM cluster")
	second = envOrSkip(t, EnvHVMSwitchName2, "a second virtual switch of the HVM cluster")

	return first, second
}

// HVMHostInterface returns the name of a host network interface of the HVM
// hosts, for the HOST_INTERFACE network layout, skipping the test when
// EnvHVMHostInterface is unset. Leave it unset for an HVM 2.0 cluster, whose
// networks attach to virtual switches instead.
func HVMHostInterface(t *testing.T) string {
	t.Helper()

	return envOrSkip(t, EnvHVMHostInterface, "a host network interface of the HVM hosts")
}

// NetworkServerBlock renders a data source that resolves a network server by
// name as data.hpe_morpheus_network_server.<label>. The name is quoted for
// HCL (HCLQuote), so one with a quote, a backslash or a `${` in it is looked
// up as written.
func NetworkServerBlock(label, name string) string {
	return fmt.Sprintf("data \"hpe_morpheus_network_server\" %q {\n  name = %s\n}\n", label, HCLQuote(name))
}

// HVMSwitchBlock renders the data source that resolves a virtual switch of
// the HVM cluster by name, as data.hpe_morpheus_option_source_value.<label>,
// whose value is what hpe_morpheus_network.switch_id takes. cloudIDRef and
// poolIDRef are the HCL expressions for the HVM cloud's id and the cluster's
// resource pool id.
func HVMSwitchBlock(label, name, cloudIDRef, poolIDRef string) string {
	return fmt.Sprintf(`data "hpe_morpheus_option_source_value" %q {
  source       = "hvm-vswitch-or-router"
  service_type = "com.hpe.morpheus.hvm"
  cloud_id     = %s
  name         = %s
  parameters = {
    "network.zonePool.id" = tostring(%s)
  }
}
`, label, cloudIDRef, HCLQuote(name), poolIDRef)
}
