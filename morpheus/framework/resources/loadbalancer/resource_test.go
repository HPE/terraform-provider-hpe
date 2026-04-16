// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package loadbalancer_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/loadbalancer"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/systemoverride"
)

func TestMain(m *testing.M) {
	systemoverride.ParseFlags()
	code := m.Run()
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func TestAccMorpheusLoadBalancerHAProxyExampleOk(t *testing.T) {
	t.Parallel()
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow acceptance test in short mode")
	}

	name := acctest.RandomWithPrefix(t.Name())
	// This resource only allows name to be 32 characters maximum.
	name = name[0:16] + name[len(name)-16:]

	config, err := loadbalancer.RenderLoadBalancerHAProxyConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatalf("failed to render config: %s", err)
	}

	resourceName := "hpe_morpheus_load_balancer.haproxy"
	config = testhelpers.ProviderBlock() + config

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "visibility", "public"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"config",
					"config_haproxy",
					"group_id",
					"cloud_id",
					"network_server_id",
				},
			},
		},
	})
}
