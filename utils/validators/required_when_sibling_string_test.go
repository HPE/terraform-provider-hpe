// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package validators_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/HPE/terraform-provider-hpe/utils/validators"
)

// TestUnitRequiredWhenSiblingString exercises the RequiredWhenSiblingString
// validator directly, covering both constructors: the values that satisfy it,
// the null/empty/whitespace values it must reject once the sibling matches,
// the null sibling that only the MatchWhenNull variant treats as a match, and
// the unknown values it must defer on.
func TestUnitRequiredWhenSiblingString(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// Both validators mirror the way the generated option_list schema wires
	// them: api_type is required when type == "api"; source_url is required
	// when type == "rest" or type is unset, because the server defaults an
	// unset type to "rest".
	requireAPIType := validators.RequiredWhenSiblingEquals("type", "api")
	requireSourceURL := validators.RequiredWhenSiblingEqualsOrNull("type", "rest")

	// A two-value variant exercises the match loop past its first entry.
	requireAPITypeMulti := validators.RequiredWhenSiblingEquals("type", "api", "rest")

	cases := []struct {
		name      string
		validator validator.String
		attribute string
		sibling   types.String
		value     types.String
		wantError bool
	}{
		// MatchWhenNull == false.
		{
			name:      "matching sibling with a value passes",
			validator: requireAPIType,
			attribute: "api_type",
			sibling:   types.StringValue("api"),
			value:     types.StringValue("clouds"),
		},
		{
			name:      "matching sibling with a null value fails",
			validator: requireAPIType,
			attribute: "api_type",
			sibling:   types.StringValue("api"),
			value:     types.StringNull(),
			wantError: true,
		},
		{
			name:      "matching sibling with an empty value fails",
			validator: requireAPIType,
			attribute: "api_type",
			sibling:   types.StringValue("api"),
			value:     types.StringValue(""),
			wantError: true,
		},
		{
			name:      "matching sibling with a whitespace value fails",
			validator: requireAPIType,
			attribute: "api_type",
			sibling:   types.StringValue("api"),
			value:     types.StringValue("   "),
			wantError: true,
		},
		{
			name:      "matching sibling with an unknown value defers",
			validator: requireAPIType,
			attribute: "api_type",
			sibling:   types.StringValue("api"),
			value:     types.StringUnknown(),
		},
		{
			name:      "non-matching sibling passes",
			validator: requireAPIType,
			attribute: "api_type",
			sibling:   types.StringValue("manual"),
			value:     types.StringNull(),
		},
		{
			name:      "null sibling passes when match-when-null is off",
			validator: requireAPIType,
			attribute: "api_type",
			sibling:   types.StringNull(),
			value:     types.StringNull(),
		},
		{
			name:      "unknown sibling defers",
			validator: requireAPIType,
			attribute: "api_type",
			sibling:   types.StringUnknown(),
			value:     types.StringNull(),
		},
		{
			name:      "a later match value still matches",
			validator: requireAPITypeMulti,
			attribute: "api_type",
			sibling:   types.StringValue("rest"),
			value:     types.StringNull(),
			wantError: true,
		},

		// MatchWhenNull == true.
		{
			name:      "null sibling with a null value fails",
			validator: requireSourceURL,
			attribute: "source_url",
			sibling:   types.StringNull(),
			value:     types.StringNull(),
			wantError: true,
		},
		{
			name:      "null sibling with an empty value fails",
			validator: requireSourceURL,
			attribute: "source_url",
			sibling:   types.StringNull(),
			value:     types.StringValue(""),
			wantError: true,
		},
		{
			name:      "null sibling with a value passes",
			validator: requireSourceURL,
			attribute: "source_url",
			sibling:   types.StringNull(),
			value:     types.StringValue("https://example.com/list.json"),
		},
		{
			name:      "matching sibling with a null value fails when match-when-null is on",
			validator: requireSourceURL,
			attribute: "source_url",
			sibling:   types.StringValue("rest"),
			value:     types.StringNull(),
			wantError: true,
		},
		{
			name:      "matching sibling with a value passes when match-when-null is on",
			validator: requireSourceURL,
			attribute: "source_url",
			sibling:   types.StringValue("rest"),
			value:     types.StringValue("https://example.com/list.json"),
		},
		{
			name:      "non-matching sibling passes when match-when-null is on",
			validator: requireSourceURL,
			attribute: "source_url",
			sibling:   types.StringValue("manual"),
			value:     types.StringNull(),
		},
		{
			name:      "unknown sibling defers when match-when-null is on",
			validator: requireSourceURL,
			attribute: "source_url",
			sibling:   types.StringUnknown(),
			value:     types.StringNull(),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			req := validator.StringRequest{
				Path:        path.Root(tc.attribute),
				ConfigValue: tc.value,
				Config:      newOptionListConfig(tc.sibling, tc.attribute, tc.value),
			}
			resp := &validator.StringResponse{}
			tc.validator.ValidateString(ctx, req, resp)

			if got := resp.Diagnostics.HasError(); got != tc.wantError {
				t.Fatalf("%s=%v (type=%v): got error=%v, want error=%v (diagnostics: %v)",
					tc.attribute, tc.value, tc.sibling, got, tc.wantError, resp.Diagnostics)
			}
		})
	}
}

// TestUnitRequiredWhenSiblingStringDescription pins the schema description
// text. The MatchWhenNull variant must advertise the unset sibling as a
// trigger, otherwise the generated documentation understates when the
// attribute is required.
func TestUnitRequiredWhenSiblingStringDescription(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	cases := []struct {
		name      string
		validator validator.String
		want      string
	}{
		{
			name:      "equals lists only the match values",
			validator: validators.RequiredWhenSiblingEquals("type", "api"),
			want:      `must be set when "type" is one of [api]`,
		},
		{
			name:      "equals lists every match value",
			validator: validators.RequiredWhenSiblingEquals("type", "api", "rest"),
			want:      `must be set when "type" is one of [api rest]`,
		},
		{
			name:      "equals-or-null also documents the unset sibling",
			validator: validators.RequiredWhenSiblingEqualsOrNull("type", "rest"),
			want:      `must be set when "type" is one of [rest] or is unset`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if got := tc.validator.Description(ctx); got != tc.want {
				t.Errorf("Description() = %q, want %q", got, tc.want)
			}

			if got := tc.validator.MarkdownDescription(ctx); got != tc.want {
				t.Errorf("MarkdownDescription() = %q, want %q", got, tc.want)
			}
		})
	}
}

// TestUnitRequiredWhenSiblingStringDiagnostic pins the diagnostic raised on
// failure. The acceptance tests match on this text, so the summary and detail
// are part of the validator's contract.
func TestUnitRequiredWhenSiblingStringDiagnostic(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	const wantSummary = "Missing required attribute"

	cases := []struct {
		name       string
		validator  validator.String
		attribute  string
		sibling    types.String
		wantDetail string
	}{
		{
			name:       "equals reports the matched values",
			validator:  validators.RequiredWhenSiblingEquals("type", "api"),
			attribute:  "api_type",
			sibling:    types.StringValue("api"),
			wantDetail: `Attribute "api_type" must be set to a non-empty value when "type" is one of [api].`,
		},
		{
			name:       "equals-or-null reports the unset sibling",
			validator:  validators.RequiredWhenSiblingEqualsOrNull("type", "rest"),
			attribute:  "source_url",
			sibling:    types.StringNull(),
			wantDetail: `Attribute "source_url" must be set to a non-empty value when "type" is one of [rest] or is unset.`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			req := validator.StringRequest{
				Path:        path.Root(tc.attribute),
				ConfigValue: types.StringNull(),
				Config:      newOptionListConfig(tc.sibling, tc.attribute, types.StringNull()),
			}
			resp := &validator.StringResponse{}
			tc.validator.ValidateString(ctx, req, resp)

			errs := resp.Diagnostics.Errors()
			if len(errs) != 1 {
				t.Fatalf("got %d error diagnostics, want 1 (diagnostics: %v)", len(errs), resp.Diagnostics)
			}

			if got := errs[0].Summary(); got != wantSummary {
				t.Errorf("summary = %q, want %q", got, wantSummary)
			}

			if got := errs[0].Detail(); got != tc.wantDetail {
				t.Errorf("detail = %q, want %q", got, tc.wantDetail)
			}
		})
	}
}

// newOptionListConfig builds a config shaped like the option_list resource, so
// the validator resolves its "type" sibling exactly as it does against the
// generated schema. Only the attribute under validation carries a value; the
// other conditionally required attribute stays null.
func newOptionListConfig(sibling types.String, attribute string, value types.String) tfsdk.Config {
	values := map[string]tftypes.Value{
		"type":       rawString(sibling),
		"source_url": rawString(types.StringNull()),
		"api_type":   rawString(types.StringNull()),
	}
	values[attribute] = rawString(value)

	return tfsdk.Config{
		Schema: schema.Schema{
			Attributes: map[string]schema.Attribute{
				"type":       schema.StringAttribute{Optional: true},
				"source_url": schema.StringAttribute{Optional: true},
				"api_type":   schema.StringAttribute{Optional: true},
			},
		},
		Raw: tftypes.NewValue(
			tftypes.Object{
				AttributeTypes: map[string]tftypes.Type{
					"type":       tftypes.String,
					"source_url": tftypes.String,
					"api_type":   tftypes.String,
				},
			},
			values,
		),
	}
}

// rawString converts a types.String into the tftypes.Value a raw config needs,
// preserving its null and unknown states.
func rawString(v types.String) tftypes.Value {
	switch {
	case v.IsUnknown():
		return tftypes.NewValue(tftypes.String, tftypes.UnknownValue)
	case v.IsNull():
		return tftypes.NewValue(tftypes.String, nil)
	default:
		return tftypes.NewValue(tftypes.String, v.ValueString())
	}
}
