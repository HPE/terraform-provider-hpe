// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

// Package framework_test holds rules that every framework resource schema has
// to satisfy. They run against the registered resources, so a schema
// regenerated from a changed config.yaml is checked as a whole.
package framework_test

import (
	"context"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/dynamicplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/numberplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

// describer is what every plan modifier has in common, whatever its value
// type.
type describer interface {
	Description(context.Context) string
}

// modifierKinds are the three constructors this file tells apart, for each
// of the framework's value types. The modifiers are unexported types, and
// RequiresReplace and RequiresReplaceIfConfigured are the same type (both
// are RequiresReplaceIf with a different condition), so the type cannot tell
// them apart; the description each constructor gives its modifier can, and
// that is what the framework itself shows the practitioner. Every value type
// is listed so that nothing is assumed about the text being the same across
// types.
var modifierKinds = []struct{ requiresReplace, ifConfigured, useState describer }{
	{
		boolplanmodifier.RequiresReplace(),
		boolplanmodifier.RequiresReplaceIfConfigured(),
		boolplanmodifier.UseStateForUnknown(),
	},
	{
		dynamicplanmodifier.RequiresReplace(),
		dynamicplanmodifier.RequiresReplaceIfConfigured(),
		dynamicplanmodifier.UseStateForUnknown(),
	},
	{
		float32planmodifier.RequiresReplace(),
		float32planmodifier.RequiresReplaceIfConfigured(),
		float32planmodifier.UseStateForUnknown(),
	},
	{
		float64planmodifier.RequiresReplace(),
		float64planmodifier.RequiresReplaceIfConfigured(),
		float64planmodifier.UseStateForUnknown(),
	},
	{
		int32planmodifier.RequiresReplace(),
		int32planmodifier.RequiresReplaceIfConfigured(),
		int32planmodifier.UseStateForUnknown(),
	},
	{
		int64planmodifier.RequiresReplace(),
		int64planmodifier.RequiresReplaceIfConfigured(),
		int64planmodifier.UseStateForUnknown(),
	},
	{
		listplanmodifier.RequiresReplace(),
		listplanmodifier.RequiresReplaceIfConfigured(),
		listplanmodifier.UseStateForUnknown(),
	},
	{
		mapplanmodifier.RequiresReplace(),
		mapplanmodifier.RequiresReplaceIfConfigured(),
		mapplanmodifier.UseStateForUnknown(),
	},
	{
		numberplanmodifier.RequiresReplace(),
		numberplanmodifier.RequiresReplaceIfConfigured(),
		numberplanmodifier.UseStateForUnknown(),
	},
	{
		objectplanmodifier.RequiresReplace(),
		objectplanmodifier.RequiresReplaceIfConfigured(),
		objectplanmodifier.UseStateForUnknown(),
	},
	{
		setplanmodifier.RequiresReplace(),
		setplanmodifier.RequiresReplaceIfConfigured(),
		setplanmodifier.UseStateForUnknown(),
	},
	{
		stringplanmodifier.RequiresReplace(),
		stringplanmodifier.RequiresReplaceIfConfigured(),
		stringplanmodifier.UseStateForUnknown(),
	},
}

// modifierDescriptionSets returns the descriptions the framework gives
// RequiresReplace, RequiresReplaceIfConfigured and UseStateForUnknown, one
// set each across every value type. It fails the test if the three sets
// overlap, because the rule below could then not tell the modifiers apart.
func modifierDescriptionSets(ctx context.Context, t *testing.T) (requiresReplace, ifConfigured, useState map[string]bool) {
	t.Helper()

	requiresReplace, ifConfigured, useState = map[string]bool{}, map[string]bool{}, map[string]bool{}

	for _, k := range modifierKinds {
		requiresReplace[k.requiresReplace.Description(ctx)] = true
		ifConfigured[k.ifConfigured.Description(ctx)] = true
		useState[k.useState.Description(ctx)] = true
	}

	for d := range requiresReplace {
		if ifConfigured[d] || useState[d] {
			t.Fatalf("the framework describes RequiresReplace and another modifier alike (%q); the rule cannot tell them apart", d)
		}
	}

	for d := range ifConfigured {
		if useState[d] {
			t.Fatalf("the framework describes RequiresReplaceIfConfigured and UseStateForUnknown alike (%q)", d)
		}
	}

	return requiresReplace, ifConfigured, useState
}

// TestUnitReplaceOnChangeAttributesKeepStateWhenUnconfigured guards against
// an attribute shape that forces a replacement on every ordinary update.
//
// Before the plan modifiers run, the framework marks every Computed attribute
// the configuration leaves null as unknown. RequiresReplace then compares
// that unknown with the state, finds them different, and asks for a
// replacement — and Terraform replaces on an unknown at a RequiresReplace
// path. So an Optional+Computed attribute with RequiresReplace and nothing
// to put the state value back is a replacement waiting to happen: renaming
// an image created without a description destroyed and recreated it.
//
// A default keeps the attribute out of the framework's pass;
// UseStateForUnknown puts the state value back before RequiresReplace runs;
// RequiresReplaceIfConfigured does not fire for an unconfigured attribute.
// Every computed attribute that carries RequiresReplace must have one of the
// three. In config.yaml that is `state_for_unknown: true` beside
// `requires_replace: true`.
func TestUnitReplaceOnChangeAttributesKeepStateWhenUnconfigured(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	requiresReplace, ifConfigured, useState := modifierDescriptionSets(ctx, t)

	var offenders []string

	// How many attributes the rule recognised as carrying each modifier. The
	// provider has well over a hundred of each, so a zero means the
	// recognition is broken and the test would otherwise pass by finding
	// nothing.
	var seenRequiresReplace, seenUseState int

	for _, newResource := range adapter.NewMorpheus().Resources(ctx) {
		r := newResource()

		var md resource.MetadataResponse

		r.Metadata(ctx, resource.MetadataRequest{ProviderTypeName: "hpe"}, &md)

		var sr resource.SchemaResponse

		r.Schema(ctx, resource.SchemaRequest{}, &sr)

		walkAttributes(sr.Schema.Attributes, sr.Schema.Blocks, "", func(name string, a schema.Attribute) {
			descriptions := modifierDescriptions(ctx, a)

			hasRequiresReplace := containsAny(descriptions, requiresReplace)
			hasUseState := containsAny(descriptions, useState)

			if hasRequiresReplace {
				seenRequiresReplace++
			}

			if hasUseState {
				seenUseState++
			}

			if !a.IsComputed() || hasField(a, "Default") || !hasRequiresReplace {
				return
			}

			if hasUseState || containsAny(descriptions, ifConfigured) {
				return
			}

			offenders = append(offenders, md.TypeName+"."+name)
		})
	}

	if seenRequiresReplace == 0 || seenUseState == 0 {
		t.Fatalf("recognised %d RequiresReplace and %d UseStateForUnknown modifiers across the provider; "+
			"the rule no longer recognises the framework's modifiers and would pass vacuously",
			seenRequiresReplace, seenUseState)
	}

	sort.Strings(offenders)

	if len(offenders) > 0 {
		t.Errorf("computed attributes with RequiresReplace and nothing to keep the state value when unconfigured "+
			"(every update that leaves them unset forces a replacement); add `state_for_unknown: true` beside "+
			"`requires_replace: true` in the resource's config.yaml and regenerate:\n  %s",
			strings.Join(offenders, "\n  "))
	}
}

// walkAttributes visits every attribute in attrs and blocks, descending into
// nested attributes and nested blocks. Names are dotted paths; collection
// elements are written as [].
func walkAttributes(
	attrs map[string]schema.Attribute,
	blocks map[string]schema.Block,
	prefix string,
	visit func(name string, a schema.Attribute),
) {
	for name, a := range attrs {
		visit(prefix+name, a)

		switch n := a.(type) {
		case schema.SingleNestedAttribute:
			walkAttributes(n.Attributes, nil, prefix+name+".", visit)
		case schema.ListNestedAttribute:
			walkAttributes(n.NestedObject.Attributes, nil, prefix+name+"[].", visit)
		case schema.SetNestedAttribute:
			walkAttributes(n.NestedObject.Attributes, nil, prefix+name+"[].", visit)
		case schema.MapNestedAttribute:
			walkAttributes(n.NestedObject.Attributes, nil, prefix+name+"[].", visit)
		}
	}

	for name, b := range blocks {
		switch n := b.(type) {
		case schema.SingleNestedBlock:
			walkAttributes(n.Attributes, n.Blocks, prefix+name+".", visit)
		case schema.ListNestedBlock:
			walkAttributes(n.NestedObject.Attributes, n.NestedObject.Blocks, prefix+name+"[].", visit)
		case schema.SetNestedBlock:
			walkAttributes(n.NestedObject.Attributes, n.NestedObject.Blocks, prefix+name+"[].", visit)
		}
	}
}

// hasField reports whether the attribute's struct has the named field set to
// a non-nil value. Every attribute kind exposes Default and PlanModifiers as
// fields of those names, which spares a case per kind.
func hasField(a schema.Attribute, field string) bool {
	f := reflect.ValueOf(a).FieldByName(field)

	return f.IsValid() && !f.IsZero()
}

// modifierDescriptions returns the descriptions of the attribute's plan
// modifiers, whatever their value type.
func modifierDescriptions(ctx context.Context, a schema.Attribute) []string {
	f := reflect.ValueOf(a).FieldByName("PlanModifiers")
	if !f.IsValid() || f.IsZero() {
		return nil
	}

	descriptions := make([]string, 0, f.Len())

	for i := range f.Len() {
		m, ok := f.Index(i).Interface().(interface{ Description(context.Context) string })
		if !ok {
			descriptions = append(descriptions, fmt.Sprintf("%T", f.Index(i).Interface()))

			continue
		}

		descriptions = append(descriptions, m.Description(ctx))
	}

	return descriptions
}

// containsAny reports whether any of the descriptions is in the set.
func containsAny(descriptions []string, set map[string]bool) bool {
	for _, d := range descriptions {
		if set[d] {
			return true
		}
	}

	return false
}
