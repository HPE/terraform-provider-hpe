// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package plan_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	morpheus "github.com/HPE/terraform-provider-hpe/internal/sdk/legacy"
	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/client"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

// priceRegressionConfig renders a minimal fixed price with a random code and no
// tenant_id, mirroring the QA HCL shape (price_type=fixed, price_unit=hour,
// incur_charges=always, currency=USD, cost=10.00).
func priceRegressionConfig(resourceName, name, code string) string {
	return fmt.Sprintf(`
resource "hpe_morpheus_price" %q {
  name          = %q
  code          = %q
  price_type    = "fixed"
  price_unit    = "hour"
  incur_charges = "always"
  currency      = "USD"
  cost          = 10.00
}
`, resourceName, name, code)
}

// newLegacyPriceTestClient builds a legacy Morpheus client from the standard
// TF_VAR_testacc_morpheus_* environment, used to assert appliance state
// (row count / active flag) directly in the regression tests.
func newLegacyPriceTestClient(t *testing.T) *morpheus.Client {
	t.Helper()

	url := os.Getenv("TF_VAR_testacc_morpheus_url")
	token := os.Getenv("TF_VAR_testacc_morpheus_access_token")
	insecure := os.Getenv("TF_VAR_testacc_morpheus_insecure") == "true"

	opts := []morpheus.ClientOption{morpheus.SkipLogin(), morpheus.WithInsecure(insecure)}

	return client.NewLegacyClient(context.Background(), url, "", "", "", token, opts...)
}

// pricesForCode returns all rows (including inactive) for an exact code.
func pricesForCode(t *testing.T, c *morpheus.Client, code string) []morpheus.Price {
	t.Helper()

	resp, err := c.ListPrices(&morpheus.Request{
		QueryParams: map[string]string{
			"code":            code,
			"includeInactive": "true",
		},
	})
	if err != nil {
		t.Fatalf("ListPrices for code %q failed: %v", code, err)
	}
	result, ok := resp.Result.(*morpheus.ListPricesResult)
	if !ok || result.Prices == nil {
		return nil
	}

	var out []morpheus.Price
	for _, p := range *result.Prices {
		if p.Code == code {
			out = append(out, p)
		}
	}

	return out
}

// TestAccMorpheusPriceDuplicateCodeRejected verifies that creating a second
// price with an already-active code fails cleanly rather than setting an
// invalid id (MORPH-15913) or leaving an inconsistent result. It also covers
// the MORPH-13247 active-orphan scenario.
func TestAccMorpheusPriceDuplicateCodeRejected(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	providerConfig := testhelpers.ProviderBlock()

	code := acctest.RandomWithPrefix("tf-acc-dupcode")
	firstName := acctest.RandomWithPrefix("tf-acc-first")
	dupName := acctest.RandomWithPrefix("tf-acc-dup")

	firstConfig := providerConfig + priceRegressionConfig("first", firstName, code)
	dupConfig := firstConfig + priceRegressionConfig("duplicate", dupName, code)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			// Step 1: create the first price and confirm a real id is set.
			{
				Config: firstConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("hpe_morpheus_price.first", "id"),
					resource.TestCheckResourceAttrWith(
						"hpe_morpheus_price.first", "id",
						func(value string) error {
							if value == "" || value == "0" {
								return fmt.Errorf("expected a non-zero id, got %q", value)
							}

							return nil
						},
					),
				),
			},
			// Step 2: a second price with the same active code must be rejected
			// by the pre-flight lookup with import guidance, before any POST.
			{
				Config:      dupConfig,
				ExpectError: regexp.MustCompile(`(?s)already exists.*terraform import`),
			},
			// Step 3: the original config still plans clean (first is intact).
			{
				Config:             firstConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

// TestAccMorpheusPriceRecreateAfterDestroy verifies that destroying a price
// (soft delete / deactivate) and re-creating it with the same code re-activates
// the existing row instead of creating a duplicate row or throwing HTTP 500
// (MORPH-13247, MORPH-15915).
func TestAccMorpheusPriceRecreateAfterDestroy(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	providerConfig := testhelpers.ProviderBlock()
	legacyClient := newLegacyPriceTestClient(t)

	code := acctest.RandomWithPrefix("tf-acc-recreate")
	name := acctest.RandomWithPrefix("tf-acc-recreate-name")

	priceConfig := providerConfig + priceRegressionConfig("p", name, code)
	emptyConfig := providerConfig + "terraform {}\n"

	var capturedID string

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			// Step 1: create and capture the id.
			{
				Config: priceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("hpe_morpheus_price.p", "id"),
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources["hpe_morpheus_price.p"]
						if !ok {
							return fmt.Errorf("resource hpe_morpheus_price.p not found in state")
						}
						capturedID = rs.Primary.ID
						if capturedID == "" || capturedID == "0" {
							return fmt.Errorf("expected a non-zero id, got %q", capturedID)
						}

						return nil
					},
				),
			},
			// Step 2: destroy the price (deactivate).
			{
				Config: emptyConfig,
			},
			// Step 3: re-create with the same code -> re-activates the same row.
			{
				Config: priceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrWith(
						"hpe_morpheus_price.p", "id",
						func(value string) error {
							if value != capturedID {
								return fmt.Errorf(
									"expected re-created price to reuse id %q, got %q",
									capturedID, value)
							}

							return nil
						},
					),
					func(_ *terraform.State) error {
						rows := pricesForCode(t, legacyClient, code)
						if len(rows) != 1 {
							return fmt.Errorf(
								"expected exactly 1 row for code %q, got %d", code, len(rows))
						}
						if !rows[0].Active {
							return fmt.Errorf("expected row for code %q to be active", code)
						}

						return nil
					},
				),
			},
		},
	})
}

// TestAccMorpheusPriceRecreateActiveCodeErrors verifies the pre-flight
// import-guidance path: creating a new resource whose code matches an existing
// ACTIVE price is rejected with actionable import guidance and never adopts the
// active price (MORPH-13247).
func TestAccMorpheusPriceRecreateActiveCodeErrors(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	providerConfig := testhelpers.ProviderBlock()

	code := acctest.RandomWithPrefix("tf-acc-activecode")
	aName := acctest.RandomWithPrefix("tf-acc-a")
	bName := acctest.RandomWithPrefix("tf-acc-b")

	aConfig := providerConfig + priceRegressionConfig("a", aName, code)
	abConfig := aConfig + priceRegressionConfig("b", bName, code)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			// Step 1: create price a.
			{
				Config: aConfig,
				Check:  resource.TestCheckResourceAttrSet("hpe_morpheus_price.a", "id"),
			},
			// Step 2: adding b with a's active code must error with import guidance.
			{
				Config:      abConfig,
				ExpectError: regexp.MustCompile(`(?s)already exists.*terraform import`),
			},
			// Step 3: a is still intact and plans clean.
			{
				Config:             aConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}
