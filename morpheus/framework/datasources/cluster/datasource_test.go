// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package cluster_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/cluster"
	clusterresource "github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/cluster"
	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/systemoverride"
)

const providerConfigOffline = `
provider "hpe" {
  morpheus {
    url          = ""
    username     = ""
    password     = ""
  }
}
`

// nolint:goconst
const clusterPrereqConfig = `
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

func TestMain(m *testing.M) {
	systemoverride.ParseFlags()
	code := m.Run()
	testhelpers.WriteMergedResults()

	os.Exit(code)
}

func TestAccMorpheusFindClusterByID(t *testing.T) {
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	name := acctest.RandomWithPrefix(t.Name())
	providerConfig := testhelpers.ProviderBlock()

	clusterConfig, err := clusterresource.RenderClusterHvmConfig(t, map[string]string{
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

	dataSourceConfig, err := testhelpers.RenderExample(t,
		"example-id.tf.tmpl", "Id", "hpe_morpheus_cluster.example_hvm.id")
	if err != nil {
		t.Fatal(err)
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr("data.hpe_morpheus_cluster.example", "name", name),
		resource.TestCheckResourceAttrPair(
			"data.hpe_morpheus_cluster.example",
			"id",
			"hpe_morpheus_cluster.example_hvm",
			"id",
		),
		resource.TestCheckResourceAttrPair(
			"data.hpe_morpheus_cluster.example",
			"cloud_id",
			"data.hpe_morpheus_cloud.test",
			"id",
		),
		resource.TestCheckResourceAttrPair(
			"data.hpe_morpheus_cluster.example",
			"layout_id",
			"hpe_morpheus_cluster.example_hvm",
			"layout_id",
		),
		resource.TestCheckResourceAttrPair(
			"data.hpe_morpheus_cluster.example",
			"description",
			"hpe_morpheus_cluster.example_hvm",
			"description",
		),
		resource.TestCheckResourceAttrPair(
			"data.hpe_morpheus_cluster.example",
			"service_url",
			"hpe_morpheus_cluster.example_hvm",
			"service_url",
		),
		resource.TestCheckResourceAttrPair(
			"data.hpe_morpheus_cluster.example",
			"group_id",
			"data.hpe_morpheus_group.test",
			"id",
		),
		resource.TestCheckResourceAttr("data.hpe_morpheus_cluster.example", "labels.#", "2"),
		resource.TestCheckTypeSetElemAttr("data.hpe_morpheus_cluster.example", "labels.*", "terraform"),
		resource.TestCheckTypeSetElemAttr("data.hpe_morpheus_cluster.example", "labels.*", "test"),
		resource.TestCheckResourceAttrPair(
			"data.hpe_morpheus_cluster.example",
			"uuid",
			"hpe_morpheus_cluster.example_hvm",
			"uuid",
		),
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + clusterPrereqConfig + clusterConfig + dataSourceConfig,
				Check:  checkFn,
			},
		},
	})
}

func TestAccMorpheusFindClusterByName(t *testing.T) {
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	name := acctest.RandomWithPrefix(t.Name())
	providerConfig := testhelpers.ProviderBlock()

	clusterConfig, err := clusterresource.RenderClusterHvmConfig(t, map[string]string{
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

	dataSourceConfig, err := testhelpers.RenderExample(t,
		"example-name.tf.tmpl", "Name", "hpe_morpheus_cluster.example_hvm.name")
	if err != nil {
		t.Fatal(err)
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrPair(
			"data.hpe_morpheus_cluster.example",
			"name",
			"hpe_morpheus_cluster.example_hvm",
			"name",
		),
		resource.TestCheckResourceAttrPair(
			"data.hpe_morpheus_cluster.example",
			"id",
			"hpe_morpheus_cluster.example_hvm",
			"id",
		),
		resource.TestCheckResourceAttrPair(
			"data.hpe_morpheus_cluster.example",
			"layout_id",
			"hpe_morpheus_cluster.example_hvm",
			"layout_id",
		),
		resource.TestCheckResourceAttrPair(
			"data.hpe_morpheus_cluster.example",
			"description",
			"hpe_morpheus_cluster.example_hvm",
			"description",
		),
		resource.TestCheckResourceAttrPair(
			"data.hpe_morpheus_cluster.example",
			"service_url",
			"hpe_morpheus_cluster.example_hvm",
			"service_url",
		),
		resource.TestCheckResourceAttr("data.hpe_morpheus_cluster.example", "labels.#", "2"),
		resource.TestCheckTypeSetElemAttr("data.hpe_morpheus_cluster.example", "labels.*", "terraform"),
		resource.TestCheckTypeSetElemAttr("data.hpe_morpheus_cluster.example", "labels.*", "test"),
		resource.TestCheckResourceAttrPair(
			"data.hpe_morpheus_cluster.example",
			"uuid",
			"hpe_morpheus_cluster.example_hvm",
			"uuid",
		),
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + clusterPrereqConfig + clusterConfig + dataSourceConfig,
				Check:  checkFn,
			},
		},
	})
}

func TestAccMorpheusFindClusterNotFound(t *testing.T) {
	t.Parallel()
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	config := providerConfig + `
      data "hpe_morpheus_cluster" "test" {
        name = "______"
      }`

	checkFn := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckNoResourceAttr("data.hpe_morpheus_cluster.test", "id"),
	)

	expected := cluster.ErrorNoClusterFound

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config:      config,
				Check:       checkFn,
				ExpectError: regexp.MustCompile(expected),
			},
		},
	})
}

func TestAccMorpheusFindClusterNoSearchAttrs(t *testing.T) {
	t.Parallel()
	defer testhelpers.RecordResult(t)

	config := providerConfigOffline + `
      data "hpe_morpheus_cluster" "test" {
      }`

	checkFn := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckNoResourceAttr("data.hpe_morpheus_cluster.test", "id"),
	)

	expected := cluster.ErrorNoValidSearchTerms

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config:      config,
				Check:       checkFn,
				ExpectError: regexp.MustCompile(expected),
			},
		},
	})
}

func TestAccMorpheusFindClusterBothSearchAttrs(t *testing.T) {
	t.Parallel()
	defer testhelpers.RecordResult(t)

	config := providerConfigOffline + `
      data "hpe_morpheus_cluster" "test" {
        id = 1
        name = "______"
      }`

	checkFn := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckNoResourceAttr("data.hpe_morpheus_cluster.test", "id"),
	)

	expected := cluster.ErrorRunningPreApply

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config:      config,
				Check:       checkFn,
				ExpectError: regexp.MustCompile(expected),
			},
		},
	})
}
