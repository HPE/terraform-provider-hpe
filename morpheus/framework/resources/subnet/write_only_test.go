// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package subnet_test

import (
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

// TestAccMorpheusSubnetResourceWriteOnlyOk: the write-only config map is
// taken by Morpheus only at creation, so a change to it plans a replacement
// (MORPH-18005); the deprecated config_version marker still replaces when
// bumped and does nothing when removed; both import forms leave a clean plan.
// The digest of a map is sha256 of its JSON form (Terraform's
// sha256(jsonencode(config))).
//
// Needs an Azure VNet with free address space: set TF_ACC_SUBNET_NETWORK_ID
// and TF_ACC_SUBNET_CIDR_PREFIX (e.g. "172.17"); the existing subnet tests
// hardcode a VNet that no longer exists.
func TestAccMorpheusSubnetResourceWriteOnlyOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.Subnet)
	testhelpers.SkipUnlessAcc(t)

	networkID := os.Getenv("TF_ACC_SUBNET_NETWORK_ID")
	if networkID == "" {
		t.Skip("TF_ACC_SUBNET_NETWORK_ID not set")
	}

	// The VNet's address space, e.g. "172.17"; two free /24s are taken from
	// its upper range.
	prefix := os.Getenv("TF_ACC_SUBNET_CIDR_PREFIX")
	if prefix == "" {
		prefix = "10.0"
	}

	name := acctest.RandomWithPrefix("tfwo")
	// Two distinct /24s inside the VNet; the third octet is random to keep
	// parallel runs apart.
	octet := acctest.RandIntRange(100, 250)
	cidr := func(n int) string { return prefix + "." + itoa(octet+n) + ".0/24" }

	// The config map as HCL, and as the JSON Terraform's jsonencode() yields
	// (sorted keys), from which the import digest is taken.
	hcl := func(v string) string {
		return `{
    subnetName = "` + name + `"
    subnetCidr = "` + v + `"
  }`
	}
	jsonOf := func(v string) string { return `{"subnetCidr":"` + v + `","subnetName":"` + name + `"}` }

	c := testhelpers.WriteOnlyCase{ //nolint:gosec // attribute names, not credentials
		Address:   "hpe_morpheus_subnet.wo",
		Attribute: "config",
		Marker:    "config_version",
		APIPath:   "/api/subnets",
		APIKey:    "subnet.config.subnetCidr",
		Config: func(value, extra string, unrelated bool) string {
			desc := "write-only matrix"
			if unrelated {
				desc = "write-only matrix, changed"
			}

			return testhelpers.ProviderBlock() + `
resource "hpe_morpheus_subnet" "wo" {
  name        = "` + name + `"
  description = "` + desc + `"
  type_id     = 8
  network_id  = ` + networkID + `
  visibility  = "private"
  config      = ` + hcl(value) + extra + `
}
`
		},
		ImportDigest: func(value string) string { return testhelpers.Sha256Hex(jsonOf(value)) },
	}

	rec := testhelpers.RecordRequests(t)

	steps := testhelpers.WriteOnlyReplaceSteps(c, rec, cidr(0), cidr(1))
	steps = append(steps, testhelpers.WriteOnlyImportSteps(c, cidr(1))...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps:                    steps,
	})
}

func itoa(n int) string { return strconv.Itoa(n) }
