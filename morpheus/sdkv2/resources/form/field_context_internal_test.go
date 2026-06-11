package form

import (
	"testing"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

// TestBuildOptionTypeRowFieldContext verifies field_context is read back from the
// API option type into state.
func TestBuildOptionTypeRowFieldContext(t *testing.T) {
	t.Parallel()

	opt := morpheus.Option{
		Type:         "text",
		Name:         "myfield",
		FieldName:    "myfield",
		FieldContext: "instance",
	}

	row := buildOptionTypeRow(opt, false)
	if got := row["field_context"]; got != "instance" {
		t.Errorf("field_context = %v, want %q", got, "instance")
	}
}

// TestApplyOptionTypeConfigFieldContext verifies field_context is sent to the API
// as fieldContext only when set, so Morpheus keeps its default ("config") when
// the attribute is omitted.
func TestApplyOptionTypeConfigFieldContext(t *testing.T) {
	t.Parallel()

	// field_context set -> sent as fieldContext
	row := map[string]any{}
	if diags := applyOptionTypeConfigByType(row, map[string]any{
		"type":          "text",
		"field_context": "instance",
	}); diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	if got := row["fieldContext"]; got != "instance" {
		t.Errorf("fieldContext = %v, want %q", got, "instance")
	}

	// field_context empty -> not sent (preserve the Morpheus default)
	row2 := map[string]any{}
	if diags := applyOptionTypeConfigByType(row2, map[string]any{
		"type":          "text",
		"field_context": "",
	}); diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	if _, ok := row2["fieldContext"]; ok {
		t.Errorf("fieldContext should be unset when field_context is empty, got %v", row2["fieldContext"])
	}
}
