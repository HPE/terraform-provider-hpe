// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package translate

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	_ "unsafe"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Marshal converts a Terraform model struct into a nested map[string]any
// suitable for JSON serialization as an API request body.
//
// It uses the tfsdk struct tags to build a flat map, applies inverse
// transforms (TF→API direction), then unflattens into the API shape.
func Marshal(ctx context.Context, model any, cc *CompiledConfig) (map[string]any, error) {
	// Step 1: Reflect over struct → flat map (snake_case keys from tfsdk tags)
	flat, err := structToFlat(ctx, model, "")
	if err != nil {
		return nil, fmt.Errorf("marshal: flattening model: %w", err)
	}

	// Step 2: Infer discriminator value from active variant if not explicitly set
	if cc.raw.Discriminator != nil {
		inferDiscriminatorForWrite(flat, cc)
	}

	// Step 3: Apply type conversions (e.g., bool → "on"/"off")
	cc.ApplyTypeConversionsForWrite(flat)

	// Step 4: Apply inverse transforms (TF names → API names, still snake_case)
	apiFlat := cc.TransformForWrite(flat)

	// Step 5: Convert snake_case keys to camelCase for the API
	apiFlat = SnakeToCamelKeys(apiFlat)

	// Step 6: Unflatten into nested map
	nested := Unflatten(apiFlat)

	// Step 7: Wrap in envelope if configured
	if cc.raw.Envelope != nil && cc.raw.Envelope.Request != "" {
		nested = map[string]any{cc.raw.Envelope.Request: nested}
	}

	// Step 8: Run PostWriteHook if configured
	if cc.postWrite != nil {
		if err := cc.postWrite(ctx, nested, model); err != nil {
			return nil, fmt.Errorf("post-write hook: %w", err)
		}
	}

	return nested, nil
}

// structToFlat reflects over a struct with tfsdk tags and builds a flat map.
// Null and unknown values are skipped (not included in the output).
func structToFlat(ctx context.Context, model any, prefix string) (map[string]any, error) {
	flat := make(map[string]any)

	v := reflect.ValueOf(model)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return flat, nil
		}
		v = v.Elem()
	}

	t := v.Type()
	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected struct, got %s", t.Kind())
	}

	for i := range t.NumField() {
		field := t.Field(i)
		tag := field.Tag.Get("tfsdk")
		if tag == "" || tag == "-" {
			continue
		}

		// Skip unexported fields (like 'state' in value types)
		if !field.IsExported() {
			continue
		}

		fieldVal := v.Field(i)
		key := tag
		if prefix != "" {
			key = prefix + "." + tag
		}

		// Get the interface value
		iface := fieldVal.Interface()

		// Check for nested object value types FIRST (before attr.Value check)
		// These have tfsdk fields + a state field, and may not satisfy attr.Value interface.
		if isNestedObjectValue(fieldVal) {
			// Check null/unknown via the state field (unexported, use unsafe)
			if isNestedNull(fieldVal) || isNestedUnknown(fieldVal) {
				continue
			}
			nestedFlat, err := structToFlat(ctx, iface, key)
			if err != nil {
				return nil, fmt.Errorf("field %s: %w", tag, err)
			}
			for k, val := range nestedFlat {
				flat[k] = val
			}

			continue
		}

		// Check if it's a Terraform attr.Value
		attrVal, isAttrValue := iface.(attr.Value)
		if isAttrValue {
			if attrVal.IsNull() || attrVal.IsUnknown() {
				continue
			}

			// Convert TF value to Go native
			native, err := attrValueToNative(ctx, attrVal)
			if err != nil {
				return nil, fmt.Errorf("field %s: %w", tag, err)
			}

			if native == nil {
				continue
			}

			// If the native value is a map (e.g., from types.Dynamic containing an object),
			// flatten it into dot-notation paths under this key.
			if m, ok := native.(map[string]any); ok {
				for k, v := range m {
					flat[key+"."+k] = v
				}

				continue
			}

			flat[key] = native

			continue
		}

		// Non-attr.Value fields (unlikely in TF models, but handle gracefully)
		flat[key] = iface
	}

	return flat, nil
}

// isNestedObjectValue checks if a reflect.Value represents a custom nested
// object type (like ConfigAwsValue) that should be recursively flattened.
// These types embed basetypes.ObjectValue or have a 'state' field and
// implement attr.Value, and their concrete struct has tfsdk-tagged fields.
func isNestedObjectValue(v reflect.Value) bool {
	t := v.Type()
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return false
	}

	// Check if it has tfsdk-tagged fields (characteristic of nested value types)
	hasTfsdkFields := false
	hasStateField := false

	for i := range t.NumField() {
		f := t.Field(i)
		if f.Tag.Get("tfsdk") != "" {
			hasTfsdkFields = true
		}
		if f.Name == "state" || f.Name == "State" {
			hasStateField = true
		}
	}

	return hasTfsdkFields && hasStateField
}

// Unmarshal populates a Terraform model struct from an API response map.
// It flattens the response, applies forward transforms (API→TF),
// then maps values into the struct fields using tfsdk tags.
//
// The plan parameter is used for plan preservation: fields not present
// in the API response will retain their plan values.
func Unmarshal(ctx context.Context, apiJSON map[string]any, model any, cc *CompiledConfig, plan any) error {
	// Step 1: Unwrap envelope if configured
	data := apiJSON
	if cc.raw.Envelope != nil && cc.raw.Envelope.Response != "" {
		if wrapped, ok := data[cc.raw.Envelope.Response]; ok {
			if m, ok := wrapped.(map[string]any); ok {
				data = m
			}
		}
	}

	// Step 2: Flatten API response (gives camelCase keys)
	flat := Flatten(data)

	// Step 3: Convert camelCase keys to snake_case (matching config.yaml conventions)
	flat = CamelToSnakeKeys(flat)

	// Step 4: Apply the first move (envelope unnest) before discriminator injection
	// This strips the "zone." prefix so discriminator can find fields at expected paths.
	if len(cc.forwardMoves) > 0 && cc.forwardMoves[0].to == "" {
		applyMove(cc.forwardMoves[0].from, cc.forwardMoves[0].to, flat)
	}

	// Step 5: Apply discriminator-based anyOf routing
	// Only inject for the variant that's active in the plan (if plan is provided)
	if cc.raw.Discriminator != nil {
		var planFlat map[string]any
		if plan != nil {
			planFlat, _ = structToFlat(ctx, plan, "")
		}

		flat = applyDiscriminatorForRead(flat, cc, planFlat)
	}

	// Step 6: Apply remaining forward transforms (skip the first if already applied)
	startIdx := 0
	if len(cc.forwardMoves) > 0 && cc.forwardMoves[0].to == "" {
		startIdx = 1
	}

	for _, m := range cc.forwardMoves[startIdx:] {
		applyMove(m.from, m.to, flat)
	}

	// Apply removes
	for path := range cc.removes {
		removeByPrefix(path, flat)
	}

	tfFlat := flat

	// Step 7: Apply type conversions AFTER transforms (e.g., "on"/"off" → bool)
	cc.ApplyTypeConversionsForRead(tfFlat)

	// Step 8: Get plan flat map for preservation
	var planFlat map[string]any
	if plan != nil {
		var err error
		planFlat, err = structToFlatRaw(ctx, plan, "")
		if err != nil {
			planFlat = nil // Fall through without plan preservation
		}
	}

	// Step 9: Apply preserve_on_read — force plan values for write-only fields
	applyPreserveOnRead(tfFlat, planFlat, cc)

	// Step 10: Apply computed_fields — extract values from different response paths
	applyComputedFields(tfFlat, data, cc)

	// Step 11: Map flat values into model struct
	if err := flatToStruct(ctx, tfFlat, model, "", planFlat); err != nil {
		return err
	}

	// Step 12: Run PostReadHook if configured
	if cc.postRead != nil {
		if err := cc.postRead(ctx, apiJSON, model, plan); err != nil {
			return fmt.Errorf("post-read hook: %w", err)
		}
	}

	return nil
}

// applyPreserveOnRead forces plan values for fields marked as preserve_on_read.
// These are write-only fields that the API never returns.
func applyPreserveOnRead(tfFlat map[string]any, planFlat map[string]any, cc *CompiledConfig) {
	if planFlat == nil || len(cc.preserveOnRead) == 0 {
		return
	}

	for path := range cc.preserveOnRead {
		// Preserve exact match
		preserveField(tfFlat, planFlat, path)

		// Also handle as a prefix — preserve all child fields
		prefix := path + "."
		for key := range planFlat {
			if strings.HasPrefix(key, prefix) && !strings.Contains(key, "__null__") && !strings.Contains(key, "__unknown__") {
				preserveField(tfFlat, planFlat, key)
			}
		}
	}
}

// preserveField copies a plan value into the tfFlat map if the API didn't provide it.
func preserveField(tfFlat map[string]any, planFlat map[string]any, key string) {
	if tfFlat[key] != nil {
		return // API provided a value, don't override
	}

	planVal, ok := planFlat[key]
	if !ok {
		return
	}

	// planFlat stores raw attr.Value objects — extract native value
	if av, ok := planVal.(interface{ IsNull() bool }); ok {
		if av.IsNull() {
			return // Don't preserve null values
		}
	}

	// Convert attr.Value to native if possible
	native, err := attrValueToNative(context.Background(), planVal.(attr.Value))
	if err == nil && native != nil {
		tfFlat[key] = native
	}
}

// applyComputedFields extracts field values from alternative response paths.
// This handles structural differences where the API returns data in a different
// shape than the request.
func applyComputedFields(tfFlat map[string]any, rawResponse map[string]any, cc *CompiledConfig) {
	if len(cc.raw.ComputedFields) == 0 {
		return
	}

	for tfField, cfg := range cc.raw.ComputedFields {
		// Only compute if not already set
		if tfFlat[tfField] != nil {
			continue
		}

		// Extract from the raw response using the dot-notation path
		val := extractFromPath(rawResponse, cfg.From)
		if val == nil {
			continue
		}

		// Type coerce if needed
		switch cfg.Type {
		case "int64":
			if f, ok := val.(float64); ok {
				tfFlat[tfField] = int64(f)
			} else {
				tfFlat[tfField] = val
			}
		default:
			tfFlat[tfField] = val
		}
	}
}

// extractFromPath extracts a value from a nested map using dot-notation with array indices.
// e.g., "groups.0.id" extracts response["groups"][0]["id"]
func extractFromPath(data map[string]any, path string) any {
	parts := strings.Split(path, ".")
	var current any = data

	for _, part := range parts {
		switch v := current.(type) {
		case map[string]any:
			current = v[part]
		case []any:
			// Parse array index
			idx := 0
			for _, c := range part {
				if c >= '0' && c <= '9' {
					idx = idx*10 + int(c-'0')
				} else {
					return nil
				}
			}

			if idx >= len(v) {
				return nil
			}

			current = v[idx]
		default:
			return nil
		}

		if current == nil {
			return nil
		}
	}

	return current
}

// structToFlatRaw is like structToFlat but preserves null/unknown values for plan comparison.
func structToFlatRaw(ctx context.Context, model any, prefix string) (map[string]any, error) {
	flat := make(map[string]any)

	v := reflect.ValueOf(model)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return flat, nil
		}
		v = v.Elem()
	}

	t := v.Type()
	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected struct, got %s", t.Kind())
	}

	for i := range t.NumField() {
		field := t.Field(i)
		tag := field.Tag.Get("tfsdk")
		if tag == "" || tag == "-" {
			continue
		}
		if !field.IsExported() {
			continue
		}

		key := tag
		if prefix != "" {
			key = prefix + "." + tag
		}

		fieldVal := v.Field(i)
		iface := fieldVal.Interface()

		// Check for nested object value types FIRST
		if isNestedObjectValue(fieldVal) {
			nestedFlat, err := structToFlatRaw(ctx, iface, key)
			if err != nil {
				return nil, err
			}
			for k, val := range nestedFlat {
				flat[k] = val
			}

			// Mark null/unknown state for plan preservation
			if isNestedNull(fieldVal) {
				flat[key+".__null__"] = true
			} else if isNestedUnknown(fieldVal) {
				flat[key+".__unknown__"] = true
			}

			continue
		}

		attrVal, isAttrValue := iface.(attr.Value)
		if isAttrValue {
			// Store the raw attr.Value for plan preservation decisions
			flat[key] = attrVal

			continue
		}

		flat[key] = iface
	}

	return flat, nil
}

// flatToStruct populates a struct from a flat map using tfsdk tags.
func flatToStruct(ctx context.Context, flat map[string]any, model any, prefix string, planFlat map[string]any) error {
	v := reflect.ValueOf(model)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	t := v.Type()

	for i := range t.NumField() {
		field := t.Field(i)
		tag := field.Tag.Get("tfsdk")
		if tag == "" || tag == "-" {
			continue
		}
		if !field.IsExported() {
			continue
		}

		key := tag
		if prefix != "" {
			key = prefix + "." + tag
		}

		fieldVal := v.Field(i)

		// Handle nested object value types
		if isNestedObjectValue(fieldVal) {
			err := populateNestedValue(ctx, flat, fieldVal, key, planFlat)
			if err != nil {
				return fmt.Errorf("field %s: %w", key, err)
			}

			continue
		}

		// Check if value exists in the flat map
		apiVal, exists := flat[key]
		if !exists || apiVal == nil {
			// Plan preservation: keep the plan value if available
			if planFlat != nil {
				if planVal, planExists := planFlat[key]; planExists {
					if av, ok := planVal.(attr.Value); ok {
						if fieldVal.CanSet() && reflect.TypeOf(av).AssignableTo(fieldVal.Type()) {
							fieldVal.Set(reflect.ValueOf(av))
						}
					}
				}
			}

			continue
		}

		// Convert the API value to the appropriate TF type
		if err := setFieldValue(ctx, fieldVal, apiVal); err != nil {
			return fmt.Errorf("field %s: %w", key, err)
		}
	}

	return nil
}

// setFieldValue sets a struct field to a value, converting from native Go types
// to Terraform framework types as needed.
func setFieldValue(ctx context.Context, fieldVal reflect.Value, value any) error {
	if !fieldVal.CanSet() {
		return nil
	}

	fieldType := fieldVal.Type()

	// Handle Set type specially — needs element type info
	if fieldType == reflect.TypeOf(types.Set{}) || fieldType == reflect.TypeOf(basetypes.SetValue{}) {
		setVal, err := nativeToSet(value)
		if err != nil {
			return err
		}

		fieldVal.Set(reflect.ValueOf(setVal))

		return nil
	}

	// Handle Dynamic type specially
	if fieldType == reflect.TypeOf(types.Dynamic{}) || fieldType == reflect.TypeOf(basetypes.DynamicValue{}) {
		// For dynamic fields, store as-is (they'll be handled elsewhere)
		fieldVal.Set(reflect.ValueOf(types.DynamicNull()))

		return nil
	}

	// Determine the target attr.Type from the field type
	targetType := getAttrTypeForField(fieldType)
	if targetType == nil {
		// Not a recognized TF type, try direct assignment
		rv := reflect.ValueOf(value)
		if rv.Type().AssignableTo(fieldType) {
			fieldVal.Set(rv)
		}

		return nil
	}

	// Convert the native value to the TF attr.Value
	attrVal, err := nativeToAttrValue(ctx, value, targetType)
	if err != nil {
		return err
	}

	if attrVal != nil && reflect.TypeOf(attrVal).AssignableTo(fieldType) {
		fieldVal.Set(reflect.ValueOf(attrVal))
	}

	return nil
}

// getAttrTypeForField returns the attr.Type for a TF framework field type.
func getAttrTypeForField(t reflect.Type) attr.Type {
	switch t {
	case reflect.TypeOf(types.String{}), reflect.TypeOf(basetypes.StringValue{}):
		return types.StringType
	case reflect.TypeOf(types.Bool{}), reflect.TypeOf(basetypes.BoolValue{}):
		return types.BoolType
	case reflect.TypeOf(types.Int64{}), reflect.TypeOf(basetypes.Int64Value{}):
		return types.Int64Type
	case reflect.TypeOf(types.Float64{}), reflect.TypeOf(basetypes.Float64Value{}):
		return types.Float64Type
	case reflect.TypeOf(types.Number{}), reflect.TypeOf(basetypes.NumberValue{}):
		return types.NumberType
	}

	// For Set, List, Dynamic — we need element type info which we can't get from reflect alone.
	// These will be handled specially.
	return nil
}

// populateNestedValue populates a nested object value type (like ConfigAwsValue).
func populateNestedValue(
	ctx context.Context, flat map[string]any, fieldVal reflect.Value,
	prefix string, planFlat map[string]any,
) error {
	// Check if any non-nil fields with this prefix exist in the flat map
	hasData := false
	childPrefix := prefix + "."

	for key, val := range flat {
		if strings.HasPrefix(key, childPrefix) && val != nil {
			hasData = true

			break
		}
	}

	if !hasData {
		// No data for this nested block.
		// If plan says null, or no plan available, keep null.
		if planFlat != nil {
			if _, unknownMarked := planFlat[prefix+".__unknown__"]; unknownMarked {
				setNestedState(fieldVal, attr.ValueStateUnknown)

				return nil
			}
		}

		setNestedState(fieldVal, attr.ValueStateNull)

		return nil
	}

	// Plan says null but API has data — only honor plan-null if this isn't an import.
	// During import, the plan is empty (only ID set), so we should populate from API.
	if planFlat != nil && !isImportContext(planFlat) {
		if _, nullMarked := planFlat[prefix+".__null__"]; nullMarked {
			setNestedState(fieldVal, attr.ValueStateNull)

			return nil
		}
	}

	// Recursively populate the nested struct's fields
	iface := fieldVal.Addr().Interface()
	if err := flatToStruct(ctx, flat, iface, prefix, planFlat); err != nil {
		return err
	}

	// Set the state to "known"
	setNestedState(fieldVal, attr.ValueStateKnown)

	return nil
}

// setNestedState sets the state field on a nested value type.
// Handles unexported 'state' fields using unsafe pointer access.
func setNestedState(v reflect.Value, state attr.ValueState) {
	// Try exported "State" field first
	stateField := v.FieldByName("State")
	if stateField.IsValid() && stateField.CanSet() {
		stateField.Set(reflect.ValueOf(state))

		return
	}

	// For unexported "state" field, we need to use the address
	stateField = v.FieldByName("state")
	if !stateField.IsValid() {
		return
	}

	// Use reflect.NewAt to get a settable pointer to the unexported field
	ptr := stateField.Addr()
	// This works because we have the address of the struct (it's addressable)
	statePtr := (*attr.ValueState)(ptr.UnsafePointer())
	*statePtr = state
}

// isImportContext detects if the plan represents an import operation.
// During import, the plan typically only has the "id" field set and everything else is null.
func isImportContext(planFlat map[string]any) bool {
	nonNullCount := 0

	for key, val := range planFlat {
		if strings.Contains(key, "__null__") || strings.Contains(key, "__unknown__") {
			continue
		}
		// Check if the attr.Value is null
		if av, ok := val.(interface{ IsNull() bool }); ok {
			if av.IsNull() {
				continue
			}
		}

		nonNullCount++
	}

	// If only 0-2 fields are non-null (typically just "id"), it's likely an import
	return nonNullCount <= 2
}

// isNestedNull checks if a nested value type is null.
// It tries the attr.Value interface first, then falls back to checking the state field.
func isNestedNull(v reflect.Value) bool {
	// Try IsNull via interface
	if iface, ok := v.Interface().(interface{ IsNull() bool }); ok {
		return iface.IsNull()
	}

	return false
}

// isNestedUnknown checks if a nested value type is unknown.
func isNestedUnknown(v reflect.Value) bool {
	if iface, ok := v.Interface().(interface{ IsUnknown() bool }); ok {
		return iface.IsUnknown()
	}

	return false
}

// inferDiscriminatorForWrite sets the discriminator field value based on which
// variant is active (non-null) in the flat map. For example, if config_hvm.* keys
// exist, it sets cloud_type_code to "standard".
func inferDiscriminatorForWrite(flat map[string]any, cc *CompiledConfig) {
	disc := cc.raw.Discriminator
	if disc == nil || disc.Field == "" {
		return
	}

	// If discriminator is already set, don't override
	if _, exists := flat[disc.Field]; exists {
		return
	}

	// Check which variant has data
	for discValue, variantField := range disc.Variants {
		prefix := variantField + "."
		for key := range flat {
			if strings.HasPrefix(key, prefix) {
				// Found data for this variant — set the discriminator
				flat[disc.Field] = discValue

				return
			}
		}
	}
}

// applyDiscriminatorForRead injects synthetic anyofN/oneofN path segments
// into the flat map so that forward moves can match them.
//
// The API response doesn't contain anyofN/oneofN path segments (they're codegen artifacts).
// But the forward moves reference them (e.g., "config.anyof0: config_aws",
// "zone_type.anyof1.code: cloud_type_code"). This function examines the moves
// and injects the missing segments where the base path matches.
func applyDiscriminatorForRead(flat map[string]any, cc *CompiledConfig, planFlat map[string]any) map[string]any {
	disc := cc.raw.Discriminator
	if disc == nil {
		return flat
	}

	// Find the discriminator value to determine which config variant is active
	var discValue string

	// Look for zone_type.code in the flat map (after zone unnest)
	if val, ok := flat["zone_type.code"]; ok {
		if s, ok := val.(string); ok {
			discValue = s
		}
	}

	// Find which config variant this maps to
	activeVariant := ""
	if discValue != "" {
		if v, ok := disc.Variants[discValue]; ok {
			activeVariant = v
		}
	}

	// If plan info is available, check which variant the user actually configured.
	// If the user used the generic "config" field instead of a typed variant,
	// we should NOT route config.* fields to the typed variant.
	if planFlat != nil {
		planActiveVariant := ""

		for _, variantField := range disc.Variants {
			prefix := variantField + "."
			for key := range planFlat {
				if strings.HasPrefix(key, prefix) {
					planActiveVariant = variantField

					break
				}
			}

			if planActiveVariant != "" {
				break
			}
		}

		// If plan has a typed variant, use that
		if planActiveVariant != "" {
			activeVariant = planActiveVariant
		} else {
			// Check if plan has the generic "config" field with data
			hasGenericConfig := false

			for key := range planFlat {
				if strings.HasPrefix(key, "config.") {
					hasGenericConfig = true

					break
				}
			}

			if hasGenericConfig {
				// User used generic config — don't route to a typed variant
				activeVariant = ""
			}
			// Otherwise (import case: plan is mostly empty), use API-inferred variant
		}
	}

	// Now inject anyofN segments for all forward moves that contain them.
	// But first, collect which config fields have their own direct moves
	// (like "config.appliance_url: appliance_url") — these should NOT be
	// injected with anyofN since they're handled by their own move.
	directConfigMoves := make(map[string]bool)

	for _, m := range cc.raw.Moves {
		for from, to := range m {
			if strings.HasPrefix(from, "template-") {
				continue
			}
			// A direct config move: starts with "config." and has no anyof segment
			if strings.HasPrefix(from, "config.") && !strings.Contains(from, "anyof") && !strings.Contains(from, "oneof") {
				// This is a direct config field move (e.g., config.appliance_url)
				// Store the base form so we can exclude it from anyof injection
				directConfigMoves[from] = true
				_ = to
			}
		}
	}

	result := make(map[string]any, len(flat))

	for k, v := range flat {
		result[k] = v
	}

	for _, m := range cc.raw.Moves {
		for from, to := range m {
			if strings.HasPrefix(from, "template-") {
				continue
			}

			// Check if this move has an anyofN/oneofN segment
			parts := strings.Split(from, ".")
			anyofIdx := -1

			for i, part := range parts {
				if isCodegenArtifact(part) {
					anyofIdx = i

					break
				}
			}

			if anyofIdx < 0 {
				continue
			}

			// Build the base path (without anyofN) and the full path (with anyofN)
			baseParts := append(append([]string{}, parts[:anyofIdx]...), parts[anyofIdx+1:]...)
			basePrefix := strings.Join(baseParts, ".")

			// For config variants, only inject if this is the active variant.
			// For non-config moves (like zone_type), always inject.
			isConfigMove := strings.HasPrefix(from, "config.")
			shouldInject := !isConfigMove ||
				to == activeVariant ||
				strings.HasPrefix(to, activeVariant+".") ||
				activeVariant == ""

			if !shouldInject {
				continue
			}

			// Check if we have a matching key in the flat map
			if val, exists := result[basePrefix]; exists {
				// Don't inject if this field has its own direct move
				if !directConfigMoves[basePrefix] {
					delete(result, basePrefix)
					result[from] = val
				}
			}

			// Also check for child keys
			basePrefixDot := basePrefix + "."
			for key, val := range result {
				if strings.HasPrefix(key, basePrefixDot) {
					// Don't inject if this child field has its own direct move
					if directConfigMoves[key] {
						continue
					}

					suffix := strings.TrimPrefix(key, basePrefixDot)
					fullKey := strings.Join(parts[:anyofIdx+1], ".") + "." + suffix
					delete(result, key)
					result[fullKey] = val
				}
			}
		}
	}

	return result
}

// attrValueToNative converts a Terraform attr.Value to a native Go value.
// Handles all primitive types and collections used in TF models.
func attrValueToNative(ctx context.Context, v attr.Value) (any, error) {
	if v == nil || v.IsNull() || v.IsUnknown() {
		return nil, nil
	}

	switch val := v.(type) {
	case basetypes.StringValue:
		return val.ValueString(), nil
	case basetypes.BoolValue:
		return val.ValueBool(), nil
	case basetypes.Int64Value:
		return val.ValueInt64(), nil
	case basetypes.Float64Value:
		return val.ValueFloat64(), nil
	case basetypes.NumberValue:
		f, _ := val.ValueBigFloat().Float64()

		return f, nil
	case basetypes.ListValue:
		return listToNative(ctx, val)
	case basetypes.SetValue:
		return setToNative(ctx, val)
	case basetypes.MapValue:
		return mapToNative(ctx, val)
	case basetypes.ObjectValue:
		return objectToNative(ctx, val)
	case basetypes.DynamicValue:
		inner := val.UnderlyingValue()
		if inner == nil || inner.IsNull() || inner.IsUnknown() {
			return nil, nil
		}

		return attrValueToNative(ctx, inner)
	default:
		return nil, fmt.Errorf("unsupported type for conversion: %T", v)
	}
}

func listToNative(ctx context.Context, l basetypes.ListValue) ([]any, error) {
	elems := l.Elements()
	result := make([]any, len(elems))
	for i, elem := range elems {
		var err error
		result[i], err = attrValueToNative(ctx, elem)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func setToNative(ctx context.Context, s basetypes.SetValue) ([]any, error) {
	elems := s.Elements()
	result := make([]any, len(elems))
	for i, elem := range elems {
		var err error
		result[i], err = attrValueToNative(ctx, elem)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func mapToNative(ctx context.Context, m basetypes.MapValue) (map[string]any, error) {
	elems := m.Elements()
	result := make(map[string]any, len(elems))
	for k, v := range elems {
		var err error
		result[k], err = attrValueToNative(ctx, v)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func objectToNative(ctx context.Context, o basetypes.ObjectValue) (map[string]any, error) {
	attrs := o.Attributes()
	result := make(map[string]any, len(attrs))
	for k, v := range attrs {
		var err error
		result[k], err = attrValueToNative(ctx, v)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

// nativeToAttrValue converts a native Go value to a Terraform attr.Value.
func nativeToAttrValue(_ context.Context, value any, targetType attr.Type) (attr.Value, error) {
	if value == nil {
		return nullForType(targetType), nil
	}

	switch targetType {
	case types.StringType:
		s, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("expected string, got %T", value)
		}

		return types.StringValue(s), nil

	case types.BoolType:
		b, ok := value.(bool)
		if !ok {
			return nil, fmt.Errorf("expected bool, got %T", value)
		}

		return types.BoolValue(b), nil

	case types.Int64Type:
		switch n := value.(type) {
		case float64:
			return types.Int64Value(int64(n)), nil
		case int64:
			return types.Int64Value(n), nil
		case int:
			return types.Int64Value(int64(n)), nil
		default:
			return nil, fmt.Errorf("expected number for int64, got %T", value)
		}

	case types.Float64Type:
		switch n := value.(type) {
		case float64:
			return types.Float64Value(n), nil
		case int64:
			return types.Float64Value(float64(n)), nil
		case int:
			return types.Float64Value(float64(n)), nil
		default:
			return nil, fmt.Errorf("expected number for float64, got %T", value)
		}

	default:
		return nil, fmt.Errorf("unsupported target type: %T", targetType)
	}
}

// nullForType returns the null value for a given attr.Type.
func nullForType(t attr.Type) attr.Value {
	switch t {
	case types.StringType:
		return types.StringNull()
	case types.BoolType:
		return types.BoolNull()
	case types.Int64Type:
		return types.Int64Null()
	case types.Float64Type:
		return types.Float64Null()
	default:
		return nil
	}
}

// nativeToSet converts a native slice to a types.Set.
// It infers the element type from the first element.
func nativeToSet(value any) (types.Set, error) {
	slice, ok := value.([]any)
	if !ok {
		return types.SetNull(types.StringType), fmt.Errorf("expected slice for set, got %T", value)
	}

	if len(slice) == 0 {
		return types.SetNull(types.StringType), nil
	}

	// Infer element type from first element
	var elemType attr.Type

	var elems []attr.Value

	for _, item := range slice {
		switch v := item.(type) {
		case string:
			elemType = types.StringType
			elems = append(elems, types.StringValue(v))
		case float64:
			elemType = types.Int64Type
			elems = append(elems, types.Int64Value(int64(v)))
		case int64:
			elemType = types.Int64Type
			elems = append(elems, types.Int64Value(v))
		case bool:
			elemType = types.BoolType
			elems = append(elems, types.BoolValue(v))
		default:
			return types.SetNull(types.StringType), fmt.Errorf("unsupported set element type: %T", item)
		}
	}

	set, diags := types.SetValue(elemType, elems)
	if diags.HasError() {
		return types.SetNull(types.StringType), fmt.Errorf("creating set: %s", diags.Errors()[0].Detail())
	}

	return set, nil
}
