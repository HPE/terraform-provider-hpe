// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instance_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/instance"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

// TestAccMorpheusInstanceInvalidNetworkExplained reproduces a customer failure:
// an HVM instance whose network belongs to one cluster's resource pool while
// the request names a different pool. Morpheus rejects it during validation
// with a bare "Invalid network"; the provider must say which network, which
// pool it belongs to and which pool was requested — and, since the requested
// pool has no selectable network at all, say that outright rather than invite
// the user to pick another network.
//
// Validation runs before any provisioning, so this fails in about a second and
// creates nothing on the appliance. The network is the HVM fixture cluster's
// (TF_VAR_testacc_morpheus_hvm_network_name); the pool targeted is a second
// cluster's (TF_VAR_testacc_morpheus_hvm_other_cluster_name). Both are looked
// up by name, so nothing here names an id -- the explanation's ids are matched
// as numbers.
func TestAccMorpheusInstanceInvalidNetworkExplained(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	name := acctest.RandomWithPrefix(t.Name())
	networkName := testhelpers.HVMNetworkName(t)

	cfg, err := instance.RenderInstanceHVMWrongPoolConfig(t, map[string]string{
		"Name": `"` + name + `"`,
	})
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: testhelpers.ProviderBlock() + cfg,
				// Terraform re-wraps diagnostics at about 80 columns, so every
				// space in the expected text is matched as \s+.
				ExpectError: regexp.MustCompile(
					// The API's own error is kept, then the explanation follows:
					// the network, the pool it belongs to (named after its
					// cluster), the pool the instance targets, and the fix.
					`(?s)"networkInterfaces":"Invalid\s+network".*` +
						`network\s+\d+\s+"` + regexp.QuoteMeta(networkName) + `"\s+belongs\s+to\s+resource\s+pool\s+\d+\s+"[^"]*";\s+` +
						`this\s+instance\s+targets\s+resource\s+pool\s+\d+\..*` +
						`Provision\s+into\s+the\s+cluster's\s+pool\s+\(hpe_morpheus_cluster:\s+permissions\.resource_pool\.id\).*` +
						// Then what the targeted pool offers instead, in the exact
						// scope the check was made in: nothing at all, or a list
						// (the second cluster usually has a Management network).
						`(No\s+network\s+at\s+all\s+is\s+selectable\s+for\s+this\s+instance|` +
						`Networks\s+Morpheus\s+would\s+accept\s+for\s+this\s+instance)\s+` +
						`\(cloud\s+\d+,\s+group\s+\d+,\s+layout\s+\d+,\s+resource\s+pool\s+\d+\)`,
				),
			},
		},
	})
}

// TestAccMorpheusInstanceInvalidNetworkListsSelectable is the other shape of
// the same failure: the requested pool does have selectable networks, and the
// rejected one is not among them. The explanation must list what Morpheus
// would have accepted. The rejected network is one of the cloud's that belongs
// to no resource pool (TF_VAR_testacc_morpheus_hvm_poolless_network_name, a
// docker bridge on the QA appliance) against the fixture cluster's own pool,
// whose selectable networks include the fixture network; a pool-less network
// passes the pool filter, so the line for it must not blame the pool.
func TestAccMorpheusInstanceInvalidNetworkListsSelectable(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	name := acctest.RandomWithPrefix(t.Name())
	poollessNetwork := testhelpers.HVMPoollessNetworkName(t)
	clusterNetwork := testhelpers.HVMNetworkName(t)

	cfg, err := instance.RenderInstanceHVMWrongPoolConfig(t, map[string]string{
		"Name":             `"` + name + `"`,
		"NetworkName":      poollessNetwork,
		"NetworkPoolless":  "true",
		"OtherClusterName": testhelpers.HVMClusterName(t), // target the fixture cluster's own pool
	})
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: testhelpers.ProviderBlock() + cfg,
				ExpectError: regexp.MustCompile(
					`(?s)"networkInterfaces":"Invalid\s+network".*` +
						`network\s+\d+\s+"` + regexp.QuoteMeta(poollessNetwork) + `"\s+belongs\s+to\s+no\s+resource\s+pool,\s+which\s+does\s+not\s+` +
						`exclude\s+it\s+from\s+resource\s+pool\s+\d+;.*` +
						`Networks\s+Morpheus\s+would\s+accept\s+for\s+this\s+instance\s+` +
						`\(cloud\s+\d+,\s+group\s+\d+,\s+layout\s+\d+,\s+resource\s+pool\s+\d+\):.*` +
						`network\s+\d+\s+"` + regexp.QuoteMeta(clusterNetwork) + `"`,
				),
			},
		},
	})
}
