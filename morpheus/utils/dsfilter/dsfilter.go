// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

// Package dsfilter provides the client-side filter block shared by plural data
// sources.
//
// Several list endpoints support only a few query parameters, so the provider
// offers a repeatable filter block that narrows results after fetching. Every
// plural data source that does this needs the same two things: turn the
// configured blocks into compiled regular expressions, and decide whether an
// item satisfies them. Only two parts differ per data source — the generated
// block type, and how to read a named field off the item — so those are the two
// things callers supply.
//
// Semantics, which callers should document consistently: blocks are ANDed, and
// the values within one block are ORed. An item with no value for a filtered
// field does not match.
package dsfilter

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// Block is one configured filter, already extracted from the data source's own
// generated block type.
type Block struct {
	// Name is the field to filter on.
	Name string
	// Values are regular expressions; the block matches if any of them do.
	Values []string
}

// Compiled is a filter with its expressions compiled.
type Compiled struct {
	field string
	res   []*regexp.Regexp
}

// Field reports which field the filter applies to, for error messages.
func (c Compiled) Field() string { return c.field }

// Compile turns configured blocks into compiled filters.
//
// summary is used as the diagnostic summary so each data source keeps its own
// wording. Returns nil if any expression fails to compile, having recorded the
// reason in diags.
func Compile(blocks []Block, summary string, diags *diag.Diagnostics) []Compiled {
	compiled := make([]Compiled, 0, len(blocks))

	for _, b := range blocks {
		res := make([]*regexp.Regexp, 0, len(b.Values))

		for _, v := range b.Values {
			re, err := regexp.Compile(v)
			if err != nil {
				diags.AddError(summary, fmt.Sprintf(
					"invalid regular expression %q for filter %q: %s", v, b.Name, err))

				return nil
			}

			res = append(res, re)
		}

		compiled = append(compiled, Compiled{field: b.Name, res: res})
	}

	return compiled
}

// Matches reports whether item satisfies every filter.
//
// fieldValue reads the named field off the item and reports whether it has one.
// An item with no value for a filtered field does not match: filtering on a
// field an item does not carry is a request to exclude it, not to ignore the
// filter.
func Matches[T any](
	item T,
	filters []Compiled,
	fieldValue func(item T, field string) (string, bool),
) bool {
	for _, f := range filters {
		value, ok := fieldValue(item, f.field)
		if !ok {
			return false
		}

		matched := false

		for _, re := range f.res {
			if re.MatchString(value) {
				matched = true

				break
			}
		}

		if !matched {
			return false
		}
	}

	return true
}
