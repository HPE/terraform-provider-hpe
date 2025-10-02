// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package utils

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

// CheckPlanAttributeAgainstAPIAttribute compares the keys in a plan Attribute (which is expected to be a
// basetypes.ObjectValuable or basetypes.DynamicValue of ObjectType) against the keys in the corresponding
// Attribute returned by the API (which is expected to be a map[string]any or sdk.MappedNullable). It returns
// a diag.Diagnostics containing warnings for any keys that are in one but not the other along with a map of
// keys that are in the API response but not in the plan.
// This can help identify misconfigurations or misunderstandings of the API.
//
// keyMap is an optional mapping of plan keys to fromApi keys, to account for any differences in naming.
// Only add differing keys to the map, identical keys do not need to be included.
//
// If either planAttribute or apiAttribute is nil, an empty diag.Diagnostics is returned.
// If either planAttribute or apiAttribute is of an unexpected type, a warning diag is returned.
//
// Note that this function only compares the presence of keys, not their values.
func CheckPlanAttributeAgainstAPIAttribute(
	ctx context.Context,
	planAttribute any,
	apiAttribute any,
	keyMap map[string]string,
) (map[string]any, diag.Diagnostics) {
	if planAttribute == nil || apiAttribute == nil {
		// Nothing to compare, return empty diags
		return nil, diag.Diagnostics{}
	}

	// Extract ObjectValue from planAttribute, or else return warning diag
	var planObject basetypes.ObjectValue
	switch plan := planAttribute.(type) {
	case basetypes.DynamicValue:
		// Need to check that this is an Object
		dynamicValue := plan.UnderlyingValue()
		dynamicObject, ok := dynamicValue.(basetypes.ObjectValue)
		if !ok {
			return nil, diag.Diagnostics{
				diag.NewWarningDiagnostic(
					"check config",
					fmt.Sprintf("expected planAttribute to be basetypes.DynamicValue of ObjectType, got %T", dynamicValue),
				),
			}
		}

		planObject = dynamicObject

	case basetypes.ObjectValuable:
		// Convert to ObjectValue and execute ToObjectValue to get Object
		planObjectValue, diagObject := plan.ToObjectValue(ctx)
		if diagObject.HasError() {
			return nil, diag.Diagnostics{
				diag.NewWarningDiagnostic(
					"check config",
					fmt.Sprintf("expected planAttribute to be basetypes.ObjectValuable, got error: %v", diagObject),
				),
			}
		}

		planObject = planObjectValue

	default:
		// Can't do anything with other types
		return nil, diag.Diagnostics{
			diag.NewWarningDiagnostic(
				"check config",
				fmt.Sprintf("expected planAttribute to be basetypes.DynamicValue or basetypes.ObjectValuable, got %T", plan),
			),
		}
	}

	// Extract map[string]any from apiAttribute, or else return warning diag
	var fromApiMap map[string]any
	switch fromApi := apiAttribute.(type) {
	case map[string]any:
		// Directly convert
		fromApiMap = fromApi

	case sdk.MappedNullable:
		// Convert using ToMap()
		apiMap, err := fromApi.ToMap()
		if err != nil {
			return nil, diag.Diagnostics{
				diag.NewWarningDiagnostic(
					"check config",
					fmt.Sprintf("expected apiAttribute to be sdk.MappedNullable, got error: %v", err),
				),
			}
		}

		fromApiMap = apiMap

	default:
		// Can't do anything with other types
		return nil, diag.Diagnostics{
			diag.NewWarningDiagnostic(
				"check config",
				fmt.Sprintf("expected apiAttribute to be map[string]any or sdk.MappedNullable, got %T", fromApi),
			),
		}
	}

	// Check the keys in the planAttribute against fromApiMap
	planAttrs := planObject.Attributes()
	// Build map of keys in planAttrs
	planKeys := make(map[string]struct{})
	for k := range planAttrs {
		planKeys[k] = struct{}{}
	}
	// Build map of keys in fromApiMap
	fromApiKeys := make(map[string]struct{})
	for k := range fromApiMap {
		fromApiKeys[k] = struct{}{}
	}

	// Apply keyMap to planKeys to convert to fromApiMap keys
	if keyMap != nil {
		convertedPlanKeys := make(map[string]struct{})
		for k := range planKeys {
			if newKey, ok := keyMap[k]; ok {
				convertedPlanKeys[newKey] = struct{}{}
			} else {
				convertedPlanKeys[k] = struct{}{}
			}
		}
		planKeys = convertedPlanKeys
	}

	// Now compare the keys in both maps
	// Warn if there are keys in one that are not in the other
	// This is not necessarily an error, as the API may return extra keys
	// or the planAttribute may have keys that are not returned by the API
	// But it is worth warning about in case of misconfiguration
	// or misunderstanding of the API
	// Find keys in planAttribute that are not in apiAttribute
	var diags diag.Diagnostics
	for k := range planKeys {
		if _, ok := fromApiKeys[k]; !ok {
			diags.AddWarning(
				"check config",
				fmt.Sprintf("key '%s' in plan not found in API response", k),
			)
		}
	}

	// Find keys in fromApiMap that are not in planAttribute
	apiKeysNotInPlan := make(map[string]any)
	for k, v := range fromApiKeys {
		if _, ok := planKeys[k]; !ok {
			diags.AddWarning(
				"check config",
				fmt.Sprintf("key '%s' in API response not found in plan", k),
			)
			apiKeysNotInPlan[k] = v
		}
	}

	return apiKeysNotInPlan, diags
}
