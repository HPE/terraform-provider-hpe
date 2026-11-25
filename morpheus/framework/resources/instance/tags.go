// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instance

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

// Tags: what the resource manages against what Morpheus holds (MORPH-20973).
//
// Morpheus does not hold only the tags it was sent: tag policies, provisioning
// workflows (addTags), layout option types and the cloud's own tag sync (vCenter
// tags, say) add tags to a new instance or overwrite a value. `tags` is
// Optional and not Computed, so the state Create and Update return must equal
// the plan exactly; storing the API's list verbatim therefore failed the apply
// ("Provider produced inconsistent result after apply ... .tags: length changed
// from 4 to 6") and tainted a healthy instance, and a configured empty set came
// back null.
//
// `tags` is therefore the set of tags this resource manages, matched by name
// (case-insensitively: Morpheus may change the case of a name, MORPH-9264; a
// VMware cloud may also change the case of a value, to that of the vCenter tag
// it matched -- a value that differs only in case is the configured one):
//
//   - a read keeps, of the API's tags, only those whose name is configured (or,
//     on a refresh, in the prior state), with the API's value -- so a value
//     Morpheus overwrites is drift, a tag Morpheus removes is drift, and a tag
//     Morpheus adds is neither;
//   - after a write, a configured tag the API does not report (or reports with
//     another value) is kept as configured, with a warning, so the apply ends
//     consistent; the next refresh shows whatever is really there;
//   - an omitted set stays null and an empty set stays empty;
//   - updates send addTags/removeTags for the difference between the prior
//     managed set and the plan, never the whole list, so Morpheus's own tags
//     survive;
//   - `tags_all` (computed) exposes everything Morpheus holds.

// apiTag is one tag as the instance read returns it.
type apiTag = sdk.AddInstance200ResponseAllOfOneOfInstanceTagsInner

// tagsAllFromAPI is the unfiltered list for tags_all (empty set, not null, when
// there are none: the attribute is computed and always known after a read).
func tagsAllFromAPI(ctx context.Context, apiTags []apiTag) types.Set {
	elems := make([]attr.Value, 0, len(apiTags))

	for _, t := range apiTags {
		if t.Name == nil {
			continue
		}

		elems = append(elems, TagsAllValue{
			Name:  types.StringPointerValue(t.Name),
			Value: types.StringPointerValue(t.Value),
			state: attr.ValueStateKnown,
		})
	}

	return types.SetValueMust(TagsAllValue{}.Type(ctx), elems)
}

// managedTags returns the tags to store: of apiTags, those whose name the
// reference set (the plan after a write, the prior state on a refresh) manages,
// named as the reference names them, valued as the API values them. After a
// write (refresh false) a reference tag the API lacks or values differently is
// kept as in the reference, and the differences are reported as a warning. A
// null reference stays null; an empty reference stays empty.
func managedTags(ctx context.Context, id int64, reference types.Set, apiTags []apiTag, refresh bool) (types.Set, diag.Diagnostics) {
	var diags diag.Diagnostics

	if reference.IsNull() || reference.IsUnknown() {
		return reference, diags
	}

	byLowerName := make(map[string]apiTag, len(apiTags))
	for _, t := range apiTags {
		if t.Name != nil {
			byLowerName[strings.ToLower(*t.Name)] = t
		}
	}

	elems := make([]attr.Value, 0, len(reference.Elements()))

	var kept []string

	for _, e := range reference.Elements() {
		ref, ok := e.(TagsValue)
		if !ok {
			continue
		}

		name := ref.Name.ValueString()
		api, found := byLowerName[strings.ToLower(name)]

		switch {
		case found && api.Value != nil && strings.EqualFold(*api.Value, ref.Value.ValueString()):
			// The same value, possibly re-cased: a VMware cloud reports a
			// value in the case of the vCenter tag it matched ("prod" comes
			// back "Prod"). The configured casing is kept, so it is neither an
			// inconsistency after the write nor drift on a refresh.
			elems = append(elems, ref)
		case found && !refresh:
			// Written, read back with another value -- or with none at all
			// (a tag Morpheus has not valued yet): keep what was configured
			// so the apply is consistent; the next refresh reads the real
			// value.
			held := "no value"
			if api.Value != nil {
				held = *api.Value
			}

			kept = append(kept, fmt.Sprintf("%s=%s (Morpheus holds %s)", name, ref.Value.ValueString(), held))
			elems = append(elems, ref)
		case found:
			elems = append(elems, TagsValue{
				Name:  types.StringValue(name),
				Value: types.StringPointerValue(api.Value),
				state: attr.ValueStateKnown,
			})
		case !refresh:
			// Written, not read back: Morpheus has not stored it (yet).
			kept = append(kept, fmt.Sprintf("%s=%s (not reported by Morpheus)", name, ref.Value.ValueString()))
			elems = append(elems, ref)
		default:
			// Refresh: gone from Morpheus, so gone from state -- drift.
		}
	}

	if len(kept) > 0 {
		slices.Sort(kept)
		diags.AddWarning(
			"Instance tags differ from the configuration after the write",
			fmt.Sprintf("Instance %d: Morpheus reports the following configured tags differently: %s. "+
				"The configured values are kept in state so the apply is consistent; the next refresh reads "+
				"what Morpheus holds, and a difference then is drift. Tags Morpheus adds itself are not managed "+
				"by this resource (see tags_all).", id, strings.Join(kept, ", ")),
		)
	}

	return types.SetValueMust(TagsValue{}.Type(ctx), elems), diags
}

// tagChanges is the difference between the managed tags in the prior state and
// in the plan, as the update request's incremental lists: removeTags for a
// name no longer configured (by name only, so a value Morpheus changed is
// removed too), addTags for a name that is new or whose configured value
// differs from the prior state (a case-only difference included). A null or
// empty plan removes every previously managed tag and nothing else.
func tagChanges(prior, plan types.Set) (
	[]sdk.UpdateInstanceRequestInstanceAddTagsInner,
	[]sdk.UpdateInstanceRequestInstanceRemoveTagsInner,
) {
	priorByName := tagMap(prior)
	planByName := tagMap(plan)

	var add []sdk.UpdateInstanceRequestInstanceAddTagsInner

	var remove []sdk.UpdateInstanceRequestInstanceRemoveTagsInner

	for lower, t := range planByName {
		// Names match case-insensitively (Morpheus may re-case a name);
		// values exactly, so a deliberate case change of a value is sent --
		// a cloud that re-cases it on read-back is tolerated by managedTags.
		if p, ok := priorByName[lower]; ok && p.Value.ValueString() == t.Value.ValueString() {
			continue
		}

		add = append(add, sdk.UpdateInstanceRequestInstanceAddTagsInner{
			Name:  t.Name.ValueStringPointer(),
			Value: t.Value.ValueStringPointer(),
		})
	}

	for lower, t := range priorByName {
		if _, ok := planByName[lower]; ok {
			continue
		}

		remove = append(remove, sdk.UpdateInstanceRequestInstanceRemoveTagsInner{Name: t.Name.ValueStringPointer()})
	}

	sortTags(add, func(t sdk.UpdateInstanceRequestInstanceAddTagsInner) *string { return t.Name })
	sortTags(remove, func(t sdk.UpdateInstanceRequestInstanceRemoveTagsInner) *string { return t.Name })

	return add, remove
}

func tagMap(set types.Set) map[string]TagsValue {
	m := map[string]TagsValue{}

	if set.IsNull() || set.IsUnknown() {
		return m
	}

	for _, e := range set.Elements() {
		if t, ok := e.(TagsValue); ok {
			m[strings.ToLower(t.Name.ValueString())] = t
		}
	}

	return m
}

// sortTags orders a request list by name, for a stable request body.
func sortTags[T any](list []T, name func(T) *string) {
	slices.SortFunc(list, func(a, b T) int {
		return strings.Compare(derefTagName(name(a)), derefTagName(name(b)))
	})
}

func derefTagName(s *string) string {
	if s == nil {
		return ""
	}

	return *s
}

// uniqueTagNamesValidator rejects a configuration whose tags repeat a name,
// ignoring case: tags are matched by name case-insensitively (Morpheus and the
// clouds re-case them), so two such entries would be one tag to Morpheus and
// one of them would be silently dropped from every update.
type uniqueTagNamesValidator struct{}

func (v uniqueTagNamesValidator) Description(context.Context) string {
	return "tag names must be unique, ignoring case"
}

func (v uniqueTagNamesValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v uniqueTagNamesValidator) ValidateResource(
	ctx context.Context,
	req resource.ValidateConfigRequest,
	resp *resource.ValidateConfigResponse,
) {
	var tags types.Set

	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("tags"), &tags)...)
	if resp.Diagnostics.HasError() || tags.IsNull() || tags.IsUnknown() {
		return
	}

	seen := map[string]string{}

	for _, e := range tags.Elements() {
		t, ok := e.(TagsValue)
		if !ok || t.Name.IsNull() || t.Name.IsUnknown() {
			continue
		}

		name := t.Name.ValueString()
		if first, dup := seen[strings.ToLower(name)]; dup {
			resp.Diagnostics.AddAttributeError(
				path.Root("tags"),
				"Duplicate tag name",
				fmt.Sprintf("tags %q and %q are the same tag to Morpheus: tag names are matched ignoring case. "+
					"Give each tag a distinct name.", first, name),
			)

			return
		}

		seen[strings.ToLower(name)] = name
	}
}
