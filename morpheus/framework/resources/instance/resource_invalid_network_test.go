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
// pool it belongs to and which pool was requested.
//
// Validation runs before any provisioning, so this fails in about a second and
// creates nothing on the appliance. It targets the shared appliance's HVM
// cluster "Duck" (pool 1, network 1 "Management") and a second cluster's pool
// (1702), as the hpe_morpheus_cluster tests do.
func TestAccMorpheusInstanceInvalidNetworkExplained(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	name := acctest.RandomWithPrefix("tfacc-wrong-pool")

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
					// The API's own error is kept, then the explanation follows.
					`(?s)"networkInterfaces":"Invalid\s+network".*` +
						`network\s+1\s+"Management"\s+belongs\s+to\s+resource\s+pool\s+1\s+"Duck";\s+` +
						`this\s+instance\s+targets\s+resource\s+pool\s+1702\..*` +
						`Provision\s+into\s+the\s+cluster's\s+pool\s+\(hpe_morpheus_cluster:\s+permissions\.resource_pool\.id\)`,
				),
			},
		},
	})
}
