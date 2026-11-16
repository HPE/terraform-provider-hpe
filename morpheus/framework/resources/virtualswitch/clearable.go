// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package virtualswitch

// This resource does not use the generic clearable.Attributes machinery: its
// clearable values are handled directly by the update diff planner, because the
// clears go through different endpoints and wire encodings than a simple PUT of
// null.
//
// Clearable on update:
//
//   - description — removing it sends an explicit description: null on the
//     PUT /{id} step (see buildDescriptionUpdate; the SDK field is an omitempty
//     *string, so the null goes through the additional-properties map).
//   - bond_mode — removing it sends bondMode: "none" in the compositeEdit
//     switchUpdate, which Morpheus treats as "clear the bond" (see
//     computeUpdatePlan and bondModeNone).
//   - a segment's vlan_id / host_ips — cleared via the compositeEdit
//     editedSegments entry, whose NullableInt32 field is always Set() (even to
//     nil) and hostIps is always a (possibly empty) list, so Morpheus's
//     editSegment containsKey check sees the field and writes an explicit
//     null/empty (see buildCompositeEdit).
//
// Not clearable: a segment's ip_address / netmask / gateway. They are
// Optional+Computed with UseStateForUnknown (Morpheus fills them from the
// management network when they are not given), so removing one from the
// configuration keeps the prior value rather than clearing it; set a new value
// to change it.
//
// Not clearable (documented for the reviewer and the docs):
//
//   - name — required; there is no empty name, and a change routes through the
//     rename endpoint rather than a clear.
//   - base_type — Optional+Computed with a default and RequiresReplace; it has
//     no empty value and is immutable after creation (replacesSwitch).
//   - code — Optional+Computed and RequiresReplace: set once at create; a
//     change replaces the switch, and when unset Morpheus's value is read.
//   - mtu — Optional+Computed with a default; Morpheus refuses to change it once
//     the switch carries networks (a refused composite edit), so it is never
//     cleared to a null, only defaulted.
