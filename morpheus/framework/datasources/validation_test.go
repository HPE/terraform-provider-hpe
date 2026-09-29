// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

// Package datasources_test holds cross-cutting acceptance tests that span every
// framework data source, as opposed to the per-data-source tests that live
// beside each implementation.
package datasources_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

// emptyLookupKeyCase is one data source that looks an object up by a free-text
// key (almost always `name`), together with any required parent-id attribute
// needed for the configuration to reach attribute validation at all.
type emptyLookupKeyCase struct {
	typeName string // e.g. hpe_morpheus_group
	key      string // the lookup attribute under test, e.g. name
	extra    string // required sibling attributes, e.g. `cloud_id = 1`, or ""
}

// emptyLookupKeyCases is every data source guarded by a minimum-length
// validator on its lookup key (MORPH-17381). Each row was derived from the data
// source specs carrying `string_length_at_least: 1`; keep it in step with them.
var emptyLookupKeyCases = []emptyLookupKeyCase{
	{"hpe_morpheus_backup_host", "name", ""},
	{"hpe_morpheus_backup_instance", "name", ""},
	{"hpe_morpheus_backup_job", "name", ""},
	{"hpe_morpheus_backup_type", "name", ""},
	{"hpe_morpheus_backup", "name", ""},
	{"hpe_morpheus_certificate", "name", ""},
	{"hpe_morpheus_cloud_affinity_group", "name", "cloud_id = 1"},
	{"hpe_morpheus_cloud", "name", ""},
	{"hpe_morpheus_cluster_affinity_group", "name", "cluster_id = 1"},
	{"hpe_morpheus_cluster_layout", "name", ""},
	{"hpe_morpheus_cluster_namespace", "name", "cluster_id = 1"},
	{"hpe_morpheus_cluster", "name", ""},
	{"hpe_morpheus_compute_server", "name", ""},
	{"hpe_morpheus_container_script", "name", ""},
	{"hpe_morpheus_datastore", "name", ""},
	{"hpe_morpheus_deployment", "name", ""},
	{"hpe_morpheus_environment", "name", ""},
	{"hpe_morpheus_group", "name", ""},
	{"hpe_morpheus_image", "name", ""},
	{"hpe_morpheus_instance_type_layout", "name", ""},
	{"hpe_morpheus_instance", "name", ""},
	{"hpe_morpheus_load_balancer_monitor", "name", "load_balancer_id = 1"},
	{"hpe_morpheus_load_balancer_pool", "name", "load_balancer_id = 1"},
	{"hpe_morpheus_load_balancer_profile", "name", "load_balancer_id = 1"},
	{"hpe_morpheus_load_balancer_virtual_server", "vip_name", "load_balancer_id = 1"},
	{"hpe_morpheus_load_balancer", "name", ""},
	{"hpe_morpheus_monitoring_alert", "name", ""},
	{"hpe_morpheus_monitoring_check_type", "name", ""},
	{"hpe_morpheus_monitoring_check", "name", ""},
	{"hpe_morpheus_monitoring_group", "name", ""},
	{"hpe_morpheus_network_dhcp_server", "name", "network_integration_id = 1"},
	{"hpe_morpheus_network_domain", "name", ""},
	{"hpe_morpheus_network_edge_cluster", "name", "network_server_id = 1"},
	{"hpe_morpheus_network_firewall_rule_group", "name", "network_integration_id = 1"},
	{"hpe_morpheus_network_firewall_rule", "name", "network_integration_id = 1"},
	{"hpe_morpheus_network_pool_server_type", "name", ""},
	{"hpe_morpheus_network_pool_server", "name", ""},
	{"hpe_morpheus_network_pool", "name", ""},
	{"hpe_morpheus_network_proxy", "name", ""},
	{"hpe_morpheus_network_router_bgp_neighbor", "ip_address", "router_id = 1"},
	{"hpe_morpheus_network_router_firewall_rule_group", "name", "router_id = 1"},
	{"hpe_morpheus_network_router_firewall_rule", "name", "router_id = 1"},
	{"hpe_morpheus_network_router_nat", "name", "router_id = 1"},
	{"hpe_morpheus_network_router_route", "name", "router_id = 1"},
	{"hpe_morpheus_network_router_type", "name", ""},
	{"hpe_morpheus_network_router", "name", ""},
	{"hpe_morpheus_network_server_group", "name", ""},
	{"hpe_morpheus_network_server", "name", ""},
	{"hpe_morpheus_network_transport_zone", "name", "network_server_id = 1"},
	{"hpe_morpheus_network_type", "name", ""},
	{"hpe_morpheus_network", "name", ""},
	{"hpe_morpheus_os_type", "name", ""},
	{"hpe_morpheus_policy", "name", ""},
	{"hpe_morpheus_provisioning_license", "name", ""},
	{"hpe_morpheus_role", "name", ""},
	{"hpe_morpheus_security_group_rule", "name", "security_group_id = 1"},
	{"hpe_morpheus_security_group", "name", ""},
	{"hpe_morpheus_service_plan", "name", ""},
	{"hpe_morpheus_storage_server", "name", ""},
	{"hpe_morpheus_storage_volume", "name", ""},
	{"hpe_morpheus_subnet_type", "name", ""},
	{"hpe_morpheus_tenant", "name", ""},
	{"hpe_morpheus_vdi_app", "name", ""},
	{"hpe_morpheus_vdi_gateway", "name", ""},
}

// TestAccMorpheusDataSourcesRejectEmptyLookupKey is the MORPH-17381 regression
// guard, generalising the MORPH-16725 fix across every by-name data source.
//
// A data source used to accept an empty lookup key: `name = ""` is non-null, so
// it flowed into the by-name search, where the API returned an empty list and
// the data source reported "not found" -- indistinguishable from a legitimate
// lookup that simply matched nothing. Each data source now carries a
// LengthAtLeast(1) validator on its lookup key, so the empty string is rejected
// at plan time, before any API call.
//
// The expected error is the validator's own message, asserted precisely: a
// regression to the old "not found" path FAILS these cases rather than passing
// them by accident. Because the failure is at plan time nothing is ever created,
// and the parent ids in `extra` never have to exist -- the id is only there so
// the configuration is complete enough to reach validation.
func TestAccMorpheusDataSourcesRejectEmptyLookupKey(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	for _, tc := range emptyLookupKeyCases {
		t.Run(tc.typeName, func(t *testing.T) {
			t.Parallel()

			config := providerConfig + fmt.Sprintf(`
      data %q "test" {
        %s = ""
        %s
      }`, tc.typeName, tc.key, tc.extra)

			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
				Steps: []resource.TestStep{
					{
						Config: config,
						// stringvalidator.LengthAtLeast(1) reports:
						//   Attribute <key> string length must be at least 1, got: 0
						ExpectError: regexp.MustCompile(
							fmt.Sprintf(`Attribute %s string length must be at least 1, got: 0`, tc.key),
						),
					},
				},
			})
		})
	}
}
