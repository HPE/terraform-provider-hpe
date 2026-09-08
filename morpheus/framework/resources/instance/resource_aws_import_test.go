// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instance_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

// TestAccMorpheusInstanceResourceAWSImport covers MORPH-16191 for config_aws:
// an imported instance must not plan changes nobody made.
//
// The config deliberately omits the four config_aws attributes that declare a
// schema default -- public_ip_type, is_ec2, create_user, no_agent -- so they
// materialise from their defaults at apply and then have to be recovered from
// the API on import. Import discards prior state, so if Read writes null (or,
// for is_ec2, dereferences a nil pointer and panics) the next plan is
// non-empty; and because config_aws forces replacement, that non-empty plan is
// a proposal to destroy and recreate the instance. isEc2OrDefault,
// publicIpTypeOrDefault, createUserOrDefault and noAgentOrDefault exist to
// prevent that.
//
// The assertion is the ImportStateVerify step: it imports by id into a fresh
// state and checks every attribute matches what apply produced, with config_aws
// deliberately NOT ignored. A mismatch there is exactly the non-empty plan the
// ticket is about -- if Read wrote null for public_ip_type, the imported value
// would differ from the "subnet" apply left in state and the step fails. This
// is the same guard the image_id test in this package relies on for the same
// class of defect.
//
// A caveat recorded honestly: a normally provisioned AWS instance has these
// fields populated in the API response (isEC2 comes back as the string
// "false", publicIpType as "subnet"), so this end-to-end test passes with or
// without the fix. It guards the integrated import path and the is_ec2 nil
// panic; the per-attribute null/absent/unparseable cases are pinned by the
// unit tests in instance_read_internal_test.go.
//
// Fixture values are appliance-specific and taken from a live AWS instance on
// the target system. layout_id is used rather than a data source because the
// "Amazon VM" layout name is not unique.
func TestAccMorpheusInstanceResourceAWSImport(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.AWS)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig := fmt.Sprintf(`
data "hpe_morpheus_cloud" "aws_cloud" {
  name = "QA Amazon"
}

resource "hpe_morpheus_instance" "aws_import_test" {
  name             = "%[1]s"
  cloud_id         = data.hpe_morpheus_cloud.aws_cloud.id
  group_id         = 1
  layout_id        = 1118
  instance_type_id = 9
  plan_id          = 548

  instance_context = "dev"

  network_interfaces = [
    {
      network_id = 28
    }
  ]

  volumes = [
    {
      root_volume              = true
      name                     = "root"
      size                     = 160
      storage_type_id          = 23
      datastore_auto_selection = "auto"
    }
  ]

  tags = [
    {
      name  = "sweepable"
      value = "true"
    },
    {
      name  = "managed_by"
      value = "terraform"
    }
  ]

  # public_ip_type, is_ec2, create_user and no_agent are intentionally omitted:
  # they are Optional+Computed with schema defaults, and leaving them out is
  # what forces the default-recovery path to run on import.
  config_aws = {
    resource_pool_id = "pool-12284"
    security_groups = [
      { id = "sg-fe44ed9a" },
    ]
  }

  timeouts = {
    create = "1h"
    delete = "20m"
  }
}
`, name)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + resourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// The omitted attributes resolve to their schema defaults.
					resource.TestCheckResourceAttr(
						"hpe_morpheus_instance.aws_import_test",
						"config_aws.public_ip_type", "subnet"),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_instance.aws_import_test",
						"config_aws.is_ec2", "false"),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_instance.aws_import_test",
						"config_aws.create_user", "false"),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_instance.aws_import_test",
						"config_aws.no_agent", "true"),
				),
			},
			{
				// Import discards prior state, so config_aws must be rebuilt
				// from the API. config_aws is deliberately NOT ignored -- it
				// carries the attributes under test. The ignored attributes are
				// pre-existing import gaps on this resource, unrelated to this
				// test, and match the affinity group and image_id tests in this
				// package.
				ResourceName:      "hpe_morpheus_instance.aws_import_test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"compute_servers",
					"config",
					"connection_info",
					"labels",
					"network_interfaces",
					"timeouts",
					"volumes",
				},
			},
		},
	})
}
