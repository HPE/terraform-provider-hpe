// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package datastore_test

import (
	"context"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
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

func TestMain(m *testing.M) {
	code := m.Run()
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

// The datastore the HVM instance tests use, resolved by name within its
// cluster (TF_VAR_testacc_morpheus_hvm_datastore_name in
// TF_VAR_testacc_morpheus_hvm_cluster_name); then the same datastore by the id
// that lookup produced. Nothing here names an id.
func TestAccMorpheusFindDatastoreById(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	clusterName := testhelpers.HVMClusterName(t)
	datastoreName := testhelpers.HVMDatastoreName(t)

	providerConfig := testhelpers.ProviderBlock()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
data "hpe_morpheus_cluster" "hvm" {
  name = "` + clusterName + `"
}

data "hpe_morpheus_datastore" "by_name" {
  name       = "` + datastoreName + `"
  cluster_id = data.hpe_morpheus_cluster.hvm.id
}

data "hpe_morpheus_datastore" "test" {
  id = data.hpe_morpheus_datastore.by_name.id
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"data.hpe_morpheus_datastore.test", "id",
						"data.hpe_morpheus_datastore.by_name", "id",
					),
					resource.TestCheckResourceAttr("data.hpe_morpheus_datastore.test", "name", datastoreName),
					resource.TestCheckResourceAttr("data.hpe_morpheus_datastore.test", "associated_resource_type", "Cluster"),
					resource.TestCheckResourceAttrPair(
						"data.hpe_morpheus_datastore.test", "associated_resource_id",
						"data.hpe_morpheus_cluster.hvm", "id",
					),
					resource.TestCheckResourceAttrSet("data.hpe_morpheus_datastore.test", "type"),
				),
			},
		},
	})
}

// The same datastore by name alone. On an appliance with one HVM cluster that
// is unambiguous; with two, every cluster has a datastore of that name and the
// lookup must say so and point at cluster_id -- both outcomes are asserted.
func TestAccMorpheusFindDatastoreByName(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	datastoreName := testhelpers.HVMDatastoreName(t)

	providerConfig := testhelpers.ProviderBlock()

	dataSourceConfig, err := testhelpers.RenderExample(t,
		"example-name.tf.tmpl",
		"Name", `"`+datastoreName+`"`,
	)
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + dataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.hpe_morpheus_datastore.test", "id"),
					resource.TestCheckResourceAttr("data.hpe_morpheus_datastore.test", "name", datastoreName),
				),
				// More than one cluster: the name is ambiguous, and the error
				// says what to do about it.
				ExpectError: ambiguousDatastoreOrNil(t, datastoreName),
			},
		},
	})
}

// ambiguousDatastoreOrNil returns the pattern of the ambiguous-name error when
// the appliance has more than one datastore of that name, and nil (no error
// expected) when it has exactly one. It asks the appliance the same question
// the data source does.
func ambiguousDatastoreOrNil(t *testing.T, name string) *regexp.Regexp {
	t.Helper()

	count := countDatastoresNamed(t, name)
	if count > 1 {
		return regexp.MustCompile(`multiple\s+datastores\s+found\s+with\s+name\s+` + regexp.QuoteMeta(name) +
			`[\s\S]*Set\s+cluster_id\s+to\s+select\s+the\s+cluster's`)
	}

	return nil
}

func TestAccMorpheusFindDatastoreNotFound(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	config := providerConfig + `
      data "hpe_morpheus_datastore" "test" {
        name = "blahasdf"
      }`

	checks := []resource.TestCheckFunc{
		resource.TestCheckNoResourceAttr(
			"data.hpe_morpheus_datastore.test",
			"id",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config:      config,
				Check:       checkFn,
				ExpectError: regexp.MustCompile("datastore blahasdf not found"),
			},
		},
	})
}

func TestAccMorpheusFindDatastoreNoSearchAttrs(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	t.Parallel()

	// A real connection is used so the data source Read runs and returns the
	// "no valid search terms" error; with an unconfigured provider the mux
	// provider fails earlier with a connection error and the validation path is
	// never reached.
	config := testhelpers.ProviderBlock() + `
      data "hpe_morpheus_datastore" "test" {
      }`

	checks := []resource.TestCheckFunc{
		resource.TestCheckNoResourceAttr(
			"data.hpe_morpheus_datastore.test",
			"id",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config:      config,
				Check:       checkFn,
				ExpectError: regexp.MustCompile("either id or name must be specified"),
			},
		},
	})
}

func TestAccMorpheusFindDatastoreBothSearchAttrs(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	t.Parallel()

	config := providerConfigOffline + `
      data "hpe_morpheus_datastore" "test" {
        id = 1
        name = "blah"
      }`

	checks := []resource.TestCheckFunc{
		resource.TestCheckNoResourceAttr(
			"data.hpe_morpheus_datastore.test",
			"id",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config:      config,
				Check:       checkFn,
				ExpectError: regexp.MustCompile("Error running pre-apply plan: exit status 1"),
			},
		},
	})
}

// Every HVM cluster has datastores of its own (local, morpheus-cloud-init,
// morpheus-images), so on an appliance with more than one cluster a lookup by
// name alone finds one per cluster and fails. cluster_id limits the lookup to
// one cluster's datastores. The fixture is the HVM cluster the instance tests
// use, resolved by name.
func TestAccMorpheusFindDatastoreInCluster(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.HVM)

	clusterName := testhelpers.HVMClusterName(t)
	datastoreName := testhelpers.HVMDatastoreName(t)

	providerConfig := testhelpers.ProviderBlock()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				// By name within the cluster: found, associated with that
				// cluster, and the filter carried through as configured.
				Config: providerConfig + `
data "hpe_morpheus_cluster" "hvm" {
  name = "` + clusterName + `"
}

data "hpe_morpheus_datastore" "scoped" {
  name       = "` + datastoreName + `"
  cluster_id = data.hpe_morpheus_cluster.hvm.id
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.hpe_morpheus_datastore.scoped", "name", datastoreName),
					resource.TestCheckResourceAttr("data.hpe_morpheus_datastore.scoped", "associated_resource_type", "Cluster"),
					resource.TestCheckResourceAttrPair(
						"data.hpe_morpheus_datastore.scoped", "cluster_id",
						"data.hpe_morpheus_cluster.hvm", "id",
					),
				),
			},
			{
				// By id with the cluster: the same datastore, accepted.
				Config: providerConfig + `
data "hpe_morpheus_cluster" "hvm" {
  name = "` + clusterName + `"
}

data "hpe_morpheus_datastore" "scoped" {
  name       = "` + datastoreName + `"
  cluster_id = data.hpe_morpheus_cluster.hvm.id
}

data "hpe_morpheus_datastore" "by_id" {
  id         = data.hpe_morpheus_datastore.scoped.id
  cluster_id = data.hpe_morpheus_cluster.hvm.id
}
`,
				Check: resource.TestCheckResourceAttrPair(
					"data.hpe_morpheus_datastore.by_id", "id",
					"data.hpe_morpheus_datastore.scoped", "id",
				),
			},
			{
				// A cluster that does not exist has no datastores: not found,
				// and the error names the cluster. (The appliance answers the
				// cluster's listing with 404 "Cluster not found"; the data
				// source reports that as not found, so a listing that failed
				// for any other reason would not pass here.)
				Config: providerConfig + `
data "hpe_morpheus_datastore" "elsewhere" {
  name       = "` + datastoreName + `"
  cluster_id = 2147483647
}
`,
				ExpectError: regexp.MustCompile(`not\s+found\s+in\s+cluster\s+2147483647`),
			},
			{
				// By id against a cluster that does not have it: refused, and
				// the error names both. This is the cluster-scoped GET's 404,
				// the branch a real cluster without the datastore takes too.
				Config: providerConfig + `
data "hpe_morpheus_cluster" "hvm" {
  name = "` + clusterName + `"
}

data "hpe_morpheus_datastore" "scoped" {
  name       = "` + datastoreName + `"
  cluster_id = data.hpe_morpheus_cluster.hvm.id
}

data "hpe_morpheus_datastore" "misplaced" {
  id         = data.hpe_morpheus_datastore.scoped.id
  cluster_id = 2147483647
}
`,
				ExpectError: regexp.MustCompile(`is\s+not\s+one\s+of\s+cluster\s+2147483647's\s+datastores`),
			},
		},
	})
}

// countDatastoresNamed asks the appliance how many datastores carry a name,
// through the same listing the data source uses.
func countDatastoresNamed(t *testing.T, name string) int {
	t.Helper()

	client, err := testhelpers.NewClientForServer(context.Background(), "")
	if err != nil {
		t.Fatalf("failed to create test client: %v", err)
	}

	datastores, _, err := client.DatastoresAPI.ListDatastores(context.Background()).Name(name).Max(100).Execute()
	if err != nil {
		t.Fatalf("failed to list datastores named %s: %v", name, err)
	}

	count := 0
	for _, d := range datastores.Datastores {
		if d.Name == name {
			count++
		}
	}

	return count
}
