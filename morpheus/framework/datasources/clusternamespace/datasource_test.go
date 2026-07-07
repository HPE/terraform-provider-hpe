// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package clusternamespace_test

import (
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/clusternamespace"
	nsresource "github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/clusternamespace"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

func TestMain(m *testing.M) {
	code := m.Run()
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

const providerConfigOffline = `
provider "hpe" {
  morpheus {
    url          = ""
    username     = ""
    password     = ""
  }
}
`

const clusterID = "571"

func TestAccMorpheusFindClusterNamespaceByName(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.Kubernetes)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := strings.ToLower(acctest.RandomWithPrefix(t.Name()))
	// to get around the 63-character name limitation
	if len(name) > 63 {
		name = name[:63]
	}

	resourceConfig, err := nsresource.RenderClusterNamespaceConfig(t, map[string]string{
		"ClusterId": clusterID,
		"Name":      name,
	})
	if err != nil {
		t.Fatal(err)
	}

	dataSourceConfig := `
data "hpe_morpheus_cluster_namespace" "example" {
  name       = "` + name + `"
  cluster_id = ` + clusterID + `
  depends_on = [hpe_morpheus_cluster_namespace.example]
}
`

	checkFn := resource.ComposeAggregateTestCheckFunc(namespaceChecks()...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + resourceConfig + dataSourceConfig,
				Check:  checkFn,
			},
		},
	})
}

func TestAccMorpheusFindClusterNamespaceById(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.Kubernetes)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := strings.ToLower(acctest.RandomWithPrefix(t.Name()))
	if len(name) > 63 {
		name = name[:63]
	}

	resourceConfig, err := nsresource.RenderClusterNamespaceConfig(t, map[string]string{
		"ClusterId": clusterID,
		"Name":      name,
	})
	if err != nil {
		t.Fatal(err)
	}

	dataSourceConfig, err := clusternamespace.RenderNamespaceByIdConfig(t, map[string]string{
		"Id":        "hpe_morpheus_cluster_namespace.example.id",
		"ClusterId": clusterID,
	})
	if err != nil {
		t.Fatal(err)
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(namespaceChecks()...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + resourceConfig + dataSourceConfig,
				Check:  checkFn,
			},
		},
	})
}

func TestAccMorpheusFindClusterNamespaceNotFound(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.Kubernetes)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	dataSourceConfig := `
data "hpe_morpheus_cluster_namespace" "example" {
  name       = "nonexistent-namespace-name-that-should-not-exist"
  cluster_id = ` + clusterID + `
}
`

	expected := regexp.MustCompile(`no cluster namespace found`)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config:      providerConfig + dataSourceConfig,
				ExpectError: expected,
			},
		},
	})
}

func TestAccMorpheusFindClusterNamespaceNoSearchAttrs(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.Kubernetes)

	t.Parallel()

	config := providerConfigOffline + `
      data "hpe_morpheus_cluster_namespace" "test" {
        cluster_id = 571
      }`

	expected := clusternamespace.ErrorNoValidSearchTerms

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile(expected),
			},
		},
	})
}

func namespaceChecks() []resource.TestCheckFunc {
	ds := "data.hpe_morpheus_cluster_namespace.example"

	return []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet(ds, "id"),
		resource.TestCheckResourceAttrSet(ds, "cluster_id"),
		resource.TestCheckResourceAttrSet(ds, "name"),
		resource.TestCheckResourceAttrSet(ds, "visibility"),
	}
}
