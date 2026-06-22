// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package blueprint

import (
	"fmt"
	"log"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// applyBlueprintPermissions sends POST /api/blueprints/{id}/permissions to set
// visibility and resourcePermissions (group access). Call after Create and Update.
func applyBlueprintPermissions(client *morpheus.Client, id int64, d *schema.ResourceData) diag.Diagnostics {
	var diags diag.Diagnostics

	visibility := "private"
	if v, ok := d.Get("visibility").(string); ok && v != "" {
		visibility = v
	}

	resourcePermissions := map[string]any{
		"all": false,
	}
	if v, ok := d.Get("all_group_access").(bool); ok {
		resourcePermissions["all"] = v
	}

	if attr, ok := d.GetOk("group_access_ids"); ok {
		if groupSet, ok := attr.(*schema.Set); ok {
			sites := make([]map[string]any, 0, groupSet.Len())
			for _, s := range groupSet.List() {
				if gid, ok := s.(int); ok {
					sites = append(sites, map[string]any{"id": gid})
				}
			}
			resourcePermissions["sites"] = sites
		}
	}

	resp, err := client.Execute(&morpheus.Request{
		Method: "POST",
		Path:   fmt.Sprintf("/api/blueprints/%d/permissions", id),
		Body: map[string]any{
			"blueprint": map[string]any{
				"visibility": visibility,
			},
			"resourcePermissions": resourcePermissions,
		},
	})
	if err != nil {
		log.Printf("API FAILURE (blueprint permissions): %s - %s", resp, err)

		return append(diags, diag.FromErr(err)...)
	}
	log.Printf("API RESPONSE (blueprint permissions): %s", resp)

	return diags
}

// This cannot currently be handled efficiently by a DiffSuppressFunc.
// See: https://github.com/hashicorp/terraform-plugin-sdk/issues/477
//
//nolint:unused
func matchTemplatesWithSchema(templates []int64, declaredTemplates []any) []int64 {
	result := make([]int64, len(declaredTemplates))

	rMap := make(map[int64]int64, len(templates))
	for _, template := range templates {
		rMap[template] = template
	}

	for i, definedTemplate := range declaredTemplates {
		// skip if type assertion failed
		if definedTemplateInt64, ok := definedTemplate.(int64); ok {
			if v, ok := rMap[definedTemplateInt64]; ok {
				// matched node type declared by ID
				result[i] = v
				delete(rMap, v)
			}
		}
	}
	// append unmatched node type to the result
	for _, rcpt := range rMap {
		result = append(result, rcpt)
	}

	return result
}
