// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package testhelpers

import "strings"

// HCLQuote returns s as an HCL quoted string literal, for a value taken from
// the environment -- a fixture name, a URL -- and spliced into a test
// configuration. Plain concatenation inside quotes breaks on a `"` or a `\`,
// and strconv.Quote is not enough either: HCL's quoted template treats `${`
// and `%{` as the start of an interpolation or a directive, so a value that
// contains them has to carry them as `$${` and `%%{` (hclsyntax: a `$` or `%`
// is literal unless the next character is `{`), and Go emits escapes HCL does
// not read (`\a`, `\v`, `\x..`). Backslash, double quote, newline, carriage
// return and tab are written the way HCL reads them; everything else is kept.
func HCLQuote(s string) string {
	return `"` + HCLEscape(s) + `"`
}

// HCLEscape is HCLQuote without the surrounding quotes: the form for a value
// that goes into a template which supplies the quotes itself, such as
// `url = "{{.Url}}"`.
func HCLEscape(s string) string {
	var b strings.Builder
	b.Grow(len(s))

	for i := 0; i < len(s); i++ {
		c := s[i]
		switch c {
		case '\\':
			b.WriteString(`\\`)
		case '"':
			b.WriteString(`\"`)
		case '\n':
			b.WriteString(`\n`)
		case '\r':
			b.WriteString(`\r`)
		case '\t':
			b.WriteString(`\t`)
		case '$', '%':
			b.WriteByte(c)
			if i+1 < len(s) && s[i+1] == '{' {
				b.WriteByte(c)
			}
		default:
			b.WriteByte(c)
		}
	}

	return b.String()
}
