// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package clusterlayouts_test

import (
	"context"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

// filterConfig builds a plural cluster layouts data source with a single
// filter block. Filter values are regular expressions, so ".+" matches any
// non-empty field value and therefore selects every layout that has the given
// field populated -- enough to prove the field name is wired into the data
// source's field resolver.
func filterConfig(field, value string) string {
	return testhelpers.ProviderBlock() + `
      data "hpe_morpheus_cluster_layouts" "test" {
        filter {
          name   = "` + field + `"
          values = ["` + value + `"]
        }
      }`
}

const dsAddr = "data.hpe_morpheus_cluster_layouts.test"

// TestAccMorpheusFindClusterLayoutsNoFilter reads the data source with no
// filter block, which returns every cluster layout on the appliance -- every
// one, not the first page: the count is checked against the total the API
// reports, so a fetch truncated at the endpoint's default page size (25) fails
// here rather than surfacing later as "no such layout" from a filter.
func TestAccMorpheusFindClusterLayoutsNoFilter(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	ctx := context.Background()

	client, err := testhelpers.NewClientForServer(ctx, "")
	if err != nil {
		t.Fatalf("failed to create API client: %v", err)
	}

	rs, hresp, err := client.ClusterLayoutsAPI.ListClusterLayouts(ctx).Max(1).Execute()
	if err != nil || rs == nil || hresp.StatusCode != http.StatusOK ||
		rs.Meta == nil || rs.Meta.Total == nil {
		t.Fatalf("failed to read the cluster layout total from the API: %v", err)
	}

	total := strconv.FormatInt(*rs.Meta.Total, 10)

	config := testhelpers.ProviderBlock() + `
      data "hpe_morpheus_cluster_layouts" "test" {
      }`

	checks := resource.ComposeAggregateTestCheckFunc(
		// Every layout the appliance holds is returned, and the materialized
		// element attributes are populated.
		resource.TestCheckResourceAttr(dsAddr, "cluster_layouts.#", total),
		resource.TestCheckResourceAttrSet(dsAddr, "cluster_layouts.0.id"),
		resource.TestCheckResourceAttrSet(dsAddr, "cluster_layouts.0.name"),
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check:  checks,
			},
		},
	})
}

// TestAccMorpheusFindClusterLayoutsByFilterField exercises each supported
// filter field. A ".+" value (regex: any non-empty string) matches any layout
// that has the field set, which proves the field is accepted and resolved. For
// the fields that are populated on every layout we additionally assert at least
// one element is returned; compute_version can be sparse, so there we only
// assert the read succeeded and the set attribute is present.
func TestAccMorpheusFindClusterLayoutsByFilterField(t *testing.T) {
	cases := []struct {
		field         string
		expectElement bool
	}{
		{"name", true},
		{"code", true},
		{"cluster_version", true},
		{"compute_version", false},
		{"provision_type", true},
	}

	for _, tc := range cases {
		t.Run(tc.field, func(t *testing.T) {
			defer testhelpers.RecordResult(t)

			capabilities.MustHaveOrSkip(t, capabilities.All)

			if testing.Short() {
				t.Skip("Skipping slow test in short mode")
			}

			t.Parallel()

			var check resource.TestCheckFunc
			if tc.expectElement {
				check = resource.TestCheckResourceAttrSet(dsAddr, "cluster_layouts.0.id")
			} else {
				check = resource.TestCheckResourceAttrSet(dsAddr, "cluster_layouts.#")
			}

			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
				Steps: []resource.TestStep{
					{
						Config: filterConfig(tc.field, ".+"),
						Check:  resource.ComposeAggregateTestCheckFunc(check),
					},
				},
			})
		})
	}
}

// TestAccMorpheusFindClusterLayoutsNoMatch filters on a value no layout can
// have, so the data source returns an empty set rather than an error.
func TestAccMorpheusFindClusterLayoutsNoMatch(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	checks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(dsAddr, "cluster_layouts.#", "0"),
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: filterConfig("cluster_version", "__no_such_version__"),
				Check:  checks,
			},
		},
	})
}

// TestAccMorpheusFindClusterLayoutsHvm2 discovers the HVM cluster layout 2.0 by
// its version, which is what the data source exists for: the combination of a
// cluster_version filter and the "mvm-" code prefix of the HVM cluster layouts
// selects the layout a 2.0 HVM cluster must be built from. Morpheus ships that
// layout from 9.1.0, so the test is skipped on older appliances.
func TestAccMorpheusFindClusterLayoutsHvm2(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	testhelpers.SkipUnlessApplianceVersionAtLeast(context.Background(), t, ">= 9.1.0")

	t.Parallel()

	config := testhelpers.ProviderBlock() + `
      data "hpe_morpheus_cluster_layouts" "test" {
        filter {
          name   = "cluster_version"
          values = ["^2\\.0"]
        }
        filter {
          name   = "code"
          values = ["^mvm-"]
        }
      }`

	checks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrSet(dsAddr, "cluster_layouts.0.id"),
		resource.TestMatchResourceAttr(dsAddr, "cluster_layouts.0.cluster_version", regexp.MustCompile(`^2\.0`)),
		resource.TestMatchResourceAttr(dsAddr, "cluster_layouts.0.code", regexp.MustCompile(`^mvm-`)),
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check:  checks,
			},
		},
	})
}
