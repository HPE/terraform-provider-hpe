// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package cluster_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

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

	// nolint:goconst
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

	// nolint:goconst
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

func TestAccMorpheusClusterHVMUpdateOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())
	updatedName := name + "-updated"

	// nolint:goconst
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

	serverConfig := `
	server = {
		service_plan_id = data.hpe_morpheus_service_plan.test.id

		ssh_username = "hvmserviceuser"
		ssh_key_pair_id = data.hpe_morpheus_key_pair.test.id

		ssh_hosts = [
			{
				name = "` + name + `-worker-1"
				ip   = "10.118.67.20"
			},
			{
				name = "` + name + `-worker-2"
				ip   = "10.118.67.21"
			},
			{
				name = "` + name + `-worker-3"
				ip   = "10.118.67.22"
			}
		]

		management_net_interface = "eth0"

		visibility = "private"

		tags = [
			{
				name  = "source"
				value = "terraform"
			},
			{
				name  = "environment"
				value = "test"
			},
		]
	}
`

	createConfig := providerConfig + dataSourcesConfig + `
resource "hpe_morpheus_cluster" "test" {
	name        = "` + name + `"
	description = "Initial description"
	cloud_id    = data.hpe_morpheus_cloud.test.id
	group_id    = data.hpe_morpheus_group.test.id
	layout_id   = 219

	labels = ["terraform", "test"]

	config_hvm = {
		create_user       = false
		dynamic_placement = false
		cpu_arch          = "x86_64"
		cpu_model         = "host-model"
		power_policy      = "balanced"
	}

` + serverConfig + `
}
`

	// Update all updatable fields in a single operation:
	// top-level: name, description, labels
	// config_hvm: cpu_arch, cpu_model, dynamic_placement, vcpu_placement_mode, power_policy
	updateConfig := providerConfig + dataSourcesConfig + `
resource "hpe_morpheus_cluster" "test" {
	name        = "` + updatedName + `"
	description = "Updated description"
	cloud_id    = data.hpe_morpheus_cloud.test.id
	group_id    = data.hpe_morpheus_group.test.id
	layout_id   = 219

	labels = ["terraform-updated", "test-updated"]

	config_hvm = {
		create_user         = false
		dynamic_placement   = true
		cpu_arch            = "x86_64"
		cpu_model           = "host-passthrough"
		power_policy        = "performance"
		vcpu_placement_mode = "auto"
	}

` + serverConfig + `
}
`

	resourceName := "hpe_morpheus_cluster.test"

	createChecks := resource.ComposeAggregateTestCheckFunc(
		// top-level fields
		resource.TestCheckResourceAttr(resourceName, "name", name),
		resource.TestCheckResourceAttr(resourceName, "description", "Initial description"),
		resource.TestCheckResourceAttr(resourceName, "labels.#", "2"),
		resource.TestCheckTypeSetElemAttr(resourceName, "labels.*", "terraform"),
		resource.TestCheckTypeSetElemAttr(resourceName, "labels.*", "test"),
		// config_hvm fields
		resource.TestCheckResourceAttr(resourceName, "config_hvm.dynamic_placement", "false"),
		resource.TestCheckResourceAttr(resourceName, "config_hvm.cpu_arch", "x86_64"),
		resource.TestCheckResourceAttr(resourceName, "config_hvm.cpu_model", "host-model"),
		resource.TestCheckResourceAttr(resourceName, "config_hvm.power_policy", "balanced"),
	)

	updateChecks := resource.ComposeAggregateTestCheckFunc(
		// top-level fields
		resource.TestCheckResourceAttr(resourceName, "name", updatedName),
		resource.TestCheckResourceAttr(resourceName, "description", "Updated description"),
		resource.TestCheckResourceAttr(resourceName, "labels.#", "2"),
		resource.TestCheckTypeSetElemAttr(resourceName, "labels.*", "terraform-updated"),
		resource.TestCheckTypeSetElemAttr(resourceName, "labels.*", "test-updated"),
		// config_hvm fields
		resource.TestCheckResourceAttr(resourceName, "config_hvm.dynamic_placement", "true"),
		resource.TestCheckResourceAttr(resourceName, "config_hvm.cpu_arch", "x86_64"),
		resource.TestCheckResourceAttr(resourceName, "config_hvm.cpu_model", "host-passthrough"),
		resource.TestCheckResourceAttr(resourceName, "config_hvm.power_policy", "performance"),
		resource.TestCheckResourceAttr(resourceName, "config_hvm.vcpu_placement_mode", "auto"),
	)

	checkInPlaceUpdate := resource.ConfigPlanChecks{
		PreApply: []plancheck.PlanCheck{
			plancheck.ExpectResourceAction(
				resourceName,
				plancheck.ResourceActionUpdate,
			),
		},
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config: createConfig,
				Check:  createChecks,
			},
			{
				Config:           updateConfig,
				Check:            updateChecks,
				ConfigPlanChecks: checkInPlaceUpdate,
			},
			{
				Config:             updateConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
		},
	})
}

func TestAccMorpheusClusterGenericUpdateOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())
	updatedName := name + "-updated"

	// nolint:goconst
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

	serverConfig := `
	server = {
		service_plan_id = data.hpe_morpheus_service_plan.test.id

		ssh_username = "hvmserviceuser"
		ssh_key_pair_id = data.hpe_morpheus_key_pair.test.id

		ssh_hosts = [
			{
				name = "` + name + `-worker-1"
				ip   = "10.118.67.20"
			},
			{
				name = "` + name + `-worker-2"
				ip   = "10.118.67.21"
			},
			{
				name = "` + name + `-worker-3"
				ip   = "10.118.67.22"
			}
		]

		management_net_interface = "eth0"

		visibility = "private"

		tags = [
			{
				name  = "source"
				value = "terraform"
			},
			{
				name  = "environment"
				value = "test"
			},
		]
	}
`

	createConfig := providerConfig + dataSourcesConfig + `
resource "hpe_morpheus_cluster" "test" {
	name              = "` + name + `"
	description       = "Initial description"
	cloud_id          = data.hpe_morpheus_cloud.test.id
	group_id          = data.hpe_morpheus_group.test.id
	layout_id         = 219
	cluster_type_code = "mvm-cluster"

	labels = ["terraform", "test"]

	config = {
		cpuArch              = "x86_64"
		cpuModel             = "host-model"
		dynamicPlacementMode = "off"
		powerPolicy          = "balanced"
	}

` + serverConfig + `
}
`

	updateConfig := providerConfig + dataSourcesConfig + `
resource "hpe_morpheus_cluster" "test" {
	name              = "` + updatedName + `"
	description       = "Updated description"
	cloud_id          = data.hpe_morpheus_cloud.test.id
	group_id          = data.hpe_morpheus_group.test.id
	layout_id         = 219
	cluster_type_code = "mvm-cluster"

	labels = ["terraform-updated", "test-updated"]

	config = {
		cpuArch              = "x86_64"
		cpuModel             = "host-passthrough"
		dynamicPlacementMode = "on"
		powerPolicy          = "performance"
		vcpuPlacementMode    = "auto"
	}

` + serverConfig + `
}
`

	resourceName := "hpe_morpheus_cluster.test"

	createChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "name", name),
		resource.TestCheckResourceAttr(resourceName, "description", "Initial description"),
		resource.TestCheckResourceAttr(resourceName, "cluster_type_code", "mvm-cluster"),
		resource.TestCheckResourceAttr(resourceName, "labels.#", "2"),
		resource.TestCheckTypeSetElemAttr(resourceName, "labels.*", "terraform"),
		resource.TestCheckTypeSetElemAttr(resourceName, "labels.*", "test"),
		resource.TestCheckResourceAttr(resourceName, "config.cpuArch", "x86_64"),
		resource.TestCheckResourceAttr(resourceName, "config.cpuModel", "host-model"),
		resource.TestCheckResourceAttr(resourceName, "config.dynamicPlacementMode", "off"),
		resource.TestCheckResourceAttr(resourceName, "config.powerPolicy", "balanced"),
	)

	updateChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "name", updatedName),
		resource.TestCheckResourceAttr(resourceName, "description", "Updated description"),
		resource.TestCheckResourceAttr(resourceName, "cluster_type_code", "mvm-cluster"),
		resource.TestCheckResourceAttr(resourceName, "labels.#", "2"),
		resource.TestCheckTypeSetElemAttr(resourceName, "labels.*", "terraform-updated"),
		resource.TestCheckTypeSetElemAttr(resourceName, "labels.*", "test-updated"),
		resource.TestCheckResourceAttr(resourceName, "config.cpuArch", "x86_64"),
		resource.TestCheckResourceAttr(resourceName, "config.cpuModel", "host-passthrough"),
		resource.TestCheckResourceAttr(resourceName, "config.dynamicPlacementMode", "on"),
		resource.TestCheckResourceAttr(resourceName, "config.powerPolicy", "performance"),
		resource.TestCheckResourceAttr(resourceName, "config.vcpuPlacementMode", "auto"),
	)

	checkInPlaceUpdate := resource.ConfigPlanChecks{
		PreApply: []plancheck.PlanCheck{
			plancheck.ExpectResourceAction(
				resourceName,
				plancheck.ResourceActionUpdate,
			),
		},
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config: createConfig,
				Check:  createChecks,
			},
			{
				Config:           updateConfig,
				Check:            updateChecks,
				ConfigPlanChecks: checkInPlaceUpdate,
			},
			{
				Config:             updateConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
		},
	})
}
