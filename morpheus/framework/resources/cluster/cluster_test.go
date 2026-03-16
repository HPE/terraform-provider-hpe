// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package cluster_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/cluster"
	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/systemoverride"
)

func TestMain(m *testing.M) {
	systemoverride.ParseFlags()
	code := m.Run()
	testhelpers.WriteMergedResults()

	os.Exit(code)
}

// Tests that our HVM example file template used for docs is a valid config
func TestAccMorpheusClusterHVMExampleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	dataSourcesConfig := `
data "hpe_morpheus_cloud" "test" {
	name = "hvm"
}

data "hpe_morpheus_group" "test" {
	name = "HVM"
}

data "hpe_morpheus_key_pair" "test" {
	name = "hvmserviceuser-keypair"
}

data "hpe_morpheus_service_plan" "test" {
	name = "Default Manual"
	provision_type_code = "manual"
}
`

	resourceConfig, err := cluster.RenderClusterHvmConfig(t, map[string]string{
		"Name":             name,
		"CloudId":          "data.hpe_morpheus_cloud.test.id",
		"GroupId":          "data.hpe_morpheus_group.test.id",
		"LayoutId":         "219", // HVM 1.2 Cluster on HVM/Ubuntu 24.04
		"Label1":           "terraform",
		"Label2":           "test",
		"CreateUser":       "false",
		"DynamicPlacement": "false",
		"CpuArch":          "x86_64",
		"CpuModel":         "host-model",
		"PowerPolicy":      "balanced",
		"ServicePlanId":    "data.hpe_morpheus_service_plan.test.id",
		"SshUsername":      "hvmserviceuser",
		"SshKeyPairId":     "data.hpe_morpheus_key_pair.test.id",
		"SshHost1Name":     name + "-worker-1",
		"SshHost1Ip":       "10.118.67.20",
		"SshHost2Name":     name + "-worker-2",
		"SshHost2Ip":       "10.118.67.21",
		"SshHost3Name":     name + "-worker-3",
		"SshHost3Ip":       "10.118.67.22",
		"Visibility":       "private",
		"Tag1Name":         "source",
		"Tag1Value":        "terraform",
		"Tag2Name":         "environment",
		"Tag2Value":        "test",
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster.example_hvm",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster.example_hvm",
			"description",
			"A test HVM cluster",
		),
		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_cluster.example_hvm",
			"cloud_id",
			"data.hpe_morpheus_cloud.test",
			"id",
		),
		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_cluster.example_hvm",
			"group_id",
			"data.hpe_morpheus_group.test",
			"id",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster.example_hvm",
			"layout_id",
			"219",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster.example_hvm",
			"labels.#",
			"2",
		),
		resource.TestCheckTypeSetElemAttr(
			"hpe_morpheus_cluster.example_hvm",
			"labels.*",
			"terraform",
		),
		resource.TestCheckTypeSetElemAttr(
			"hpe_morpheus_cluster.example_hvm",
			"labels.*",
			"test",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster.example_hvm",
			"config_hvm.create_user",
			"false",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster.example_hvm",
			"config_hvm.dynamic_placement",
			"false",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster.example_hvm",
			"config_hvm.cpu_arch",
			"x86_64",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster.example_hvm",
			"config_hvm.cpu_model",
			"host-model",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster.example_hvm",
			"config_hvm.power_policy",
			"balanced",
		),
		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_cluster.example_hvm",
			"server.service_plan_id",
			"data.hpe_morpheus_service_plan.test",
			"id",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster.example_hvm",
			"server.ssh_username",
			"hvmserviceuser",
		),
		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_cluster.example_hvm",
			"server.ssh_key_pair_id",
			"data.hpe_morpheus_key_pair.test",
			"id",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster.example_hvm",
			"server.ssh_hosts.#",
			"3",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster.example_hvm",
			"server.visibility",
			"private",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config:             providerConfig + dataSourcesConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				Check:              checkFn,
				PlanOnly:           false,
			},
			{
				Config:             providerConfig + dataSourcesConfig + resourceConfig,
				ExpectNonEmptyPlan: true,
				PlanOnly:           true,
			},
		},
	})
}

// Tests that our generic example file template used for docs is a valid config
func TestAccMorpheusClusterGenericExampleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	dataSourcesConfig := `
data "hpe_morpheus_cloud" "test" {
	name = "hvm"
}

data "hpe_morpheus_group" "test" {
	name = "HVM"
}

data "hpe_morpheus_key_pair" "test" {
	name = "hvmserviceuser-keypair"
}

data "hpe_morpheus_service_plan" "test" {
	name = "Default Manual"
	provision_type_code = "manual"
}
`

	resourceConfig, err := cluster.RenderClusterGenericConfig(t, map[string]string{
		"Name":                   name,
		"CloudId":                "data.hpe_morpheus_cloud.test.id",
		"GroupId":                "data.hpe_morpheus_group.test.id",
		"LayoutId":               "219", // HVM 1.2 Cluster on HVM/Ubuntu 24.04
		"Label1":                 "terraform",
		"Label2":                 "test",
		"ClusterTypeCode":        "mvm-cluster",
		"ServicePlanId":          "data.hpe_morpheus_service_plan.test.id",
		"SshPort":                "22",
		"SshUsername":            "hvmserviceuser",
		"SshKeyPairId":           "data.hpe_morpheus_key_pair.test.id",
		"ManagementNetInterface": "eth0",
		"SshHost1Name":           name + "-worker-1",
		"SshHost1Ip":             "10.118.67.20",
		"SshHost2Name":           name + "-worker-2",
		"SshHost2Ip":             "10.118.67.21",
		"SshHost3Name":           name + "-worker-3",
		"SshHost3Ip":             "10.118.67.22",
		"Visibility":             "private",
		"Tag1Name":               "source",
		"Tag1Value":              "terraform",
		"Tag2Name":               "environment",
		"Tag2Value":              "test",
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster.example_generic_hvm",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster.example_generic_hvm",
			"description",
			"A test generic cluster",
		),
		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_cluster.example_generic_hvm",
			"cloud_id",
			"data.hpe_morpheus_cloud.test",
			"id",
		),
		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_cluster.example_generic_hvm",
			"group_id",
			"data.hpe_morpheus_group.test",
			"id",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster.example_generic_hvm",
			"layout_id",
			"219",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster.example_generic_hvm",
			"cluster_type_code",
			"mvm-cluster",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster.example_generic_hvm",
			"labels.#",
			"2",
		),
		resource.TestCheckTypeSetElemAttr(
			"hpe_morpheus_cluster.example_generic_hvm",
			"labels.*",
			"terraform",
		),
		resource.TestCheckTypeSetElemAttr(
			"hpe_morpheus_cluster.example_generic_hvm",
			"labels.*",
			"test",
		),
		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_cluster.example_generic_hvm",
			"server.service_plan_id",
			"data.hpe_morpheus_service_plan.test",
			"id",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster.example_generic_hvm",
			"server.ssh_port",
			"22",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster.example_generic_hvm",
			"server.ssh_username",
			"hvmserviceuser",
		),
		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_cluster.example_generic_hvm",
			"server.ssh_key_pair_id",
			"data.hpe_morpheus_key_pair.test",
			"id",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster.example_generic_hvm",
			"server.management_net_interface",
			"eth0",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster.example_generic_hvm",
			"server.ssh_hosts.#",
			"3",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster.example_generic_hvm",
			"server.visibility",
			"private",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config:             providerConfig + dataSourcesConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				Check:              checkFn,
				PlanOnly:           false,
			},
			{
				Config:             providerConfig + dataSourcesConfig + resourceConfig,
				ExpectNonEmptyPlan: true,
				PlanOnly:           true,
			},
		},
	})
}
