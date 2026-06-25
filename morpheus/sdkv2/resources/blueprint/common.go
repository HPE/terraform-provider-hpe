// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package blueprint

import (
	"log"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

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

// applyBlueprintPermissions calls PUT /api/blueprints/{id}/update-permissions to set group access.
func applyBlueprintPermissions(client *morpheus.Client, id int64, d *schema.ResourceData) error {
	groupAccessAll := d.Get("group_access_all").(bool)
	var sites []map[string]any
	if v, ok := d.GetOk("group_ids"); ok {
		for _, g := range v.(*schema.Set).List() {
			sites = append(sites, map[string]any{"id": g.(int)})
		}
	}
	req := &morpheus.Request{
		Body: map[string]any{
			"resourcePermission": map[string]any{
				"all":   groupAccessAll,
				"sites": sites,
			},
		},
	}
	resp, err := client.UpdateBlueprintPermissions(id, req)
	if err != nil {
		log.Printf("API FAILURE (blueprint permissions): %s - %s", resp, err)

		return err
	}

	return nil
}

// applyLayoutPermissions calls PUT /api/library/layouts/{id}/permissions to set group access.
func applyLayoutPermissions(client *morpheus.Client, id int64, d *schema.ResourceData) error {
	groupAccessAll := d.Get("group_access_all").(bool)
	var sites []map[string]any
	if v, ok := d.GetOk("group_ids"); ok {
		for _, g := range v.(*schema.Set).List() {
			sites = append(sites, map[string]any{"id": g.(int)})
		}
	}
	req := &morpheus.Request{
		Body: map[string]any{
			"instanceTypeLayout": map[string]any{
				"permissions": map[string]any{
					"resourcePermissions": map[string]any{
						"all":   groupAccessAll,
						"sites": sites,
					},
				},
			},
		},
	}
	resp, err := client.UpdateInstanceLayoutPermissions(id, req)
	if err != nil {
		log.Printf("API FAILURE (layout permissions): %s - %s", resp, err)

		return err
	}

	return nil
}

// setBlueprintPermissionsInState sets visibility, group_access_all, and group_ids in state
// from the blueprint API response values.
func setBlueprintPermissionsInState(d *schema.ResourceData, visibility string, all bool, sites []any) {
	d.Set("visibility", visibility)
	d.Set("group_access_all", all)

	var groupIDs []int
	for _, s := range sites {
		if siteMap, ok := s.(map[string]any); ok {
			switch v := siteMap["id"].(type) {
			case float64:
				groupIDs = append(groupIDs, int(v))
			case int64:
				groupIDs = append(groupIDs, int(v))
			case int:
				groupIDs = append(groupIDs, v)
			}
		}
	}
	d.Set("group_ids", groupIDs)
}

// setLayoutPermissionsInState sets group_access_all and group_ids in state
// from the layout API response values.
func setLayoutPermissionsInState(d *schema.ResourceData, all bool, siteIDs []int64) {
	d.Set("group_access_all", all)
	groupIDs := make([]int64, 0)
	groupIDs = append(groupIDs, siteIDs...)
	d.Set("group_ids", groupIDs)
}
