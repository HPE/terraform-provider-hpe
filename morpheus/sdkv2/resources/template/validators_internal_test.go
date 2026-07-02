// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package template

import (
	"testing"

	"github.com/hashicorp/go-cty/cty"
)

// TestUnitCheckSpecTemplateSource verifies the source_type-conditional validation
// (MORPH-13329 / MORPH-13325). A required field is rejected when it is null
// (omitted) or known-but-empty, but an unknown value (e.g. derived from another
// resource or a data source) must be allowed through so the plan can proceed and
// resources can be chained.
func TestUnitCheckSpecTemplateSource(t *testing.T) {
	t.Parallel()

	raw := func(sourceType, specContent, specPath, repoID cty.Value) cty.Value {
		return cty.ObjectVal(map[string]cty.Value{
			"source_type":   sourceType,
			"spec_content":  specContent,
			"spec_path":     specPath,
			"repository_id": repoID,
		})
	}

	str := cty.StringVal
	nullStr := cty.NullVal(cty.String)
	unknownStr := cty.UnknownVal(cty.String)
	nullNum := cty.NullVal(cty.Number)
	unknownNum := cty.UnknownVal(cty.Number)

	tests := []struct {
		name    string
		raw     cty.Value
		wantErr bool
	}{
		// local -> spec_content
		{"local null content", raw(str("local"), nullStr, nullStr, nullNum), true},
		{"local empty content", raw(str("local"), str(""), nullStr, nullNum), true},
		{"local whitespace content", raw(str("local"), str("   "), nullStr, nullNum), true},
		{"local valid content", raw(str("local"), str("resource {}"), nullStr, nullNum), false},
		{"local unknown content passes", raw(str("local"), unknownStr, nullStr, nullNum), false},

		// url -> spec_path
		{"url null path", raw(str("url"), nullStr, nullStr, nullNum), true},
		{"url empty path", raw(str("url"), nullStr, str(""), nullNum), true},
		{"url valid path", raw(str("url"), nullStr, str("http://x/spec.tf"), nullNum), false},
		{"url unknown path passes", raw(str("url"), nullStr, unknownStr, nullNum), false},

		// repository -> repository_id
		{"repository null id", raw(str("repository"), nullStr, str("p.tf"), nullNum), true},
		{"repository zero id", raw(str("repository"), nullStr, str("p.tf"), cty.NumberIntVal(0)), true},
		{"repository valid id", raw(str("repository"), nullStr, str("p.tf"), cty.NumberIntVal(5)), false},
		{"repository unknown id passes", raw(str("repository"), nullStr, str("p.tf"), unknownNum), false},

		// edge cases: nothing to validate
		{"unknown source_type passes", raw(unknownStr, nullStr, nullStr, nullNum), false},
		{"null source_type passes", raw(nullStr, nullStr, nullStr, nullNum), false},
		{"null raw passes", cty.NullVal(cty.EmptyObject), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := checkSpecTemplateSource(tt.raw)
			if tt.wantErr && err == nil {
				t.Errorf("expected an error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}
}
