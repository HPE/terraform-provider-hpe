// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package template

import (
	"context"
	"fmt"
	"strings"

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
// spec_content, spec_path and repository_id are Optional+Computed, so the check
// reads the raw config: it reflects exactly what the user set and never fires on
// a computed read-back of these fields.
func validateSpecTemplateSource(
	_ context.Context, d *schema.ResourceDiff, _ any,
) error {
	raw := d.GetRawConfig()
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
		if v.IsNull() || !v.IsKnown() || strings.TrimSpace(v.AsString()) == "" {
			return fmt.Errorf(
				"spec_content is required and must not be empty when source_type is %q",
				sourceTypeLocal,
			)
		}
	case sourceTypeURL:
		v := raw.GetAttr("spec_path")
		if v.IsNull() || !v.IsKnown() || strings.TrimSpace(v.AsString()) == "" {
			return fmt.Errorf(
				"spec_path is required when source_type is %q",
				sourceTypeURL,
			)
		}
	case sourceTypeRepository:
		v := raw.GetAttr("repository_id")
		if v.IsNull() || !v.IsKnown() || v.AsBigFloat().Sign() == 0 {
			return fmt.Errorf(
				"repository_id is required when source_type is %q",
				sourceTypeRepository,
			)
		}
	}

	return nil
}
