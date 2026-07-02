// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package template

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// validateSpecTemplateSource is a CustomizeDiff function for spec template
// resources. It enforces the source_type-conditional input requirements that the
// Morpheus UI enforces but the API silently accepts, so an invalid configuration
// fails at plan time instead of creating a broken spec template:
//
//   - source_type "local"      requires a non-empty spec_content
//   - source_type "url"        requires a non-empty spec_path
//   - source_type "repository" requires repository_id
//
// It reads the raw config so it reflects exactly what the user set and never
// fires on a computed read-back of these Optional+Computed fields.
func validateSpecTemplateSource(
	_ context.Context, d *schema.ResourceDiff, _ any,
) error {
	return checkSpecTemplateSource(d.GetRawConfig())
}

// checkSpecTemplateSource holds the source_type validation logic against the raw
// config so it can be unit tested directly.
//
// A required field is rejected only when it is explicitly null (omitted) or
// known-but-empty. An unknown value - e.g. spec_content, spec_path or
// repository_id derived from another resource or data source - is allowed
// through so the plan can proceed and resources can be chained; its concrete
// value is validated by the API at apply time.
func checkSpecTemplateSource(raw cty.Value) error {
	if raw.IsNull() || !raw.IsKnown() {
		return nil
	}

	sourceType := raw.GetAttr("source_type")
	if sourceType.IsNull() || !sourceType.IsKnown() {
		return nil
	}

	switch sourceType.AsString() {
	case sourceTypeLocal:
		v := raw.GetAttr("spec_content")
		if v.IsNull() || (v.IsKnown() && strings.TrimSpace(v.AsString()) == "") {
			return fmt.Errorf(
				"spec_content is required and must not be empty when source_type is %q",
				sourceTypeLocal,
			)
		}
	case sourceTypeURL:
		v := raw.GetAttr("spec_path")
		if v.IsNull() || (v.IsKnown() && strings.TrimSpace(v.AsString()) == "") {
			return fmt.Errorf(
				"spec_path is required when source_type is %q",
				sourceTypeURL,
			)
		}
	case sourceTypeRepository:
		v := raw.GetAttr("repository_id")
		if v.IsNull() || (v.IsKnown() && v.AsBigFloat().Sign() == 0) {
			return fmt.Errorf(
				"repository_id is required when source_type is %q",
				sourceTypeRepository,
			)
		}
	}

	return nil
}
