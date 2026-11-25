// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package testhelpers_test

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

// Every value HCLQuote writes must read back as itself through HCL's own
// parser -- that, not the escaping, is the contract.
func TestUnitHCLQuoteRoundTrips(t *testing.T) {
	t.Parallel()

	values := []string{
		"",
		"zodiac-tower",
		"Demo Inventory",
		`a "quoted" name`,
		`back\slash`,
		"two\nlines\twith a tab\r",
		"${not.an.interpolation}",
		"%{ not a directive }",
		"$${already escaped}",
		"a lone $ and a lone % and $$ and %%",
		"unicode: ÄÖÜ – 日本",
	}

	for _, value := range values {
		t.Run(value, func(t *testing.T) {
			t.Parallel()

			quoted := testhelpers.HCLQuote(value)

			expr, diags := hclsyntax.ParseExpression([]byte(quoted), "test.hcl", hcl.InitialPos)
			if diags.HasErrors() {
				t.Fatalf("HCLQuote(%q) = %s, which HCL refuses: %s", value, quoted, diags.Error())
			}

			got, diags := expr.Value(&hcl.EvalContext{})
			if diags.HasErrors() {
				t.Fatalf("HCLQuote(%q) = %s, which HCL cannot evaluate: %s", value, quoted, diags.Error())
			}

			if got.Type() != cty.String || got.AsString() != value {
				t.Fatalf("HCLQuote(%q) = %s, which reads back as %#v", value, quoted, got)
			}

			if want := `"` + testhelpers.HCLEscape(value) + `"`; quoted != want {
				t.Fatalf("HCLQuote(%q) = %s, but HCLEscape gives %s", value, quoted, want)
			}
		})
	}
}
