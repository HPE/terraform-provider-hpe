// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package images

import (
	"regexp"
	"strings"
	"testing"
)

// A filter block is applied after the fetch, so a configuration filtering only
// by block reads the whole library. Where the expression starts with a literal,
// that literal appears in every name it can match, and the server can be asked
// for names containing it.
func TestPhraseFromFilters(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		filters []compiledFilter
		want    string
		wantOK  bool
	}{
		{
			name:    "plain literal",
			filters: compiled(t, "name", "ubuntu"),
			want:    "ubuntu",
			wantOK:  true,
		},
		{
			name:    "anchored literal",
			filters: compiled(t, "name", "^ubuntu"),
			want:    "ubuntu",
			wantOK:  true,
		},
		{
			name:    "fully anchored literal",
			filters: compiled(t, "name", "^ubuntu-22$"),
			want:    "ubuntu-22",
			wantOK:  true,
		},
		{
			name:    "literal followed by a pattern",
			filters: compiled(t, "name", "ubuntu.*"),
			want:    "ubuntu",
			wantOK:  true,
		},
		{
			// The prefix stops at the first metacharacter — `.` matches any
			// character, so only what precedes it is guaranteed.
			name:    "prefix stops at a metacharacter",
			filters: compiled(t, "name", "ubuntu-22.04"),
			want:    "ubuntu-22",
			wantOK:  true,
		},
		{
			name:    "literal before a character class",
			filters: compiled(t, "name", "^rhel[0-9]"),
			want:    "rhel",
			wantOK:  true,
		},
		{
			// A quirk of the standard library rather than a rule: `^ubuntu`
			// and `ubuntu.*` both yield "ubuntu", but combining them does not.
			// Deriving nothing is always safe — it just means the read is not
			// narrowed — so this is a missed optimisation, not a fault.
			name:    "anchor followed by a wildcard yields nothing",
			filters: compiled(t, "name", "^ubuntu.*"),
			wantOK:  false,
		},
		{
			// No common prefix, so nothing can be asked of the server.
			name:    "alternation has no literal prefix",
			filters: compiled(t, "name", "ubuntu|debian"),
			wantOK:  false,
		},
		{
			name:    "leading wildcard has no literal prefix",
			filters: compiled(t, "name", ".*ubuntu"),
			wantOK:  false,
		},
		{
			name:    "character class has no literal prefix",
			filters: compiled(t, "name", "^[a-z]"),
			wantOK:  false,
		},
		{
			// Values within a block are ORed, and one phrase cannot express
			// that union.
			name:    "multi-valued block is skipped",
			filters: compiled(t, "name", "ubuntu", "debian"),
			wantOK:  false,
		},
		{
			// Only the name field maps to phrase.
			name:    "other fields are skipped",
			filters: compiled(t, "status", "active"),
			wantOK:  false,
		},
		{
			name:   "no filters",
			want:   "",
			wantOK: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, ok := phraseFromFilters(tt.filters)
			if ok != tt.wantOK {
				t.Fatalf("derived = %v, want %v (phrase %q)", ok, tt.wantOK, got)
			}

			if ok && got != tt.want {
				t.Errorf("phrase = %q, want %q", got, tt.want)
			}
		})
	}
}

// The pre-filter is only correct if the server returns everything the block
// would accept. phrase is a case-insensitive substring match on the name, so
// any name the expression matches must contain the derived phrase — otherwise a
// row the practitioner asked for would be dropped before the block ever saw it.
func TestPhraseFromFiltersNeverExcludesAMatch(t *testing.T) {
	t.Parallel()

	patterns := []string{
		"ubuntu", "^ubuntu", "^ubuntu-22$", "ubuntu.*", "ubuntu[0-9]*",
		"ubuntu-22.04", "^rhel[0-9]", "^ubuntu.*",
	}

	// Names chosen to include matches and non-matches for each pattern.
	names := []string{
		"ubuntu", "ubuntu-22", "ubuntu-22.04", "ubuntu-22X04", "xubuntu",
		"UBUNTU", "debian-12", "ubuntu2", "", "prefixed-ubuntu-suffix",
		"rhel9", "rhel", "RHEL9",
	}

	for _, p := range patterns {
		t.Run(p, func(t *testing.T) {
			t.Parallel()

			re := regexp.MustCompile(p)

			phrase, ok := phraseFromFilters(compiled(t, "name", p))
			if !ok {
				t.Skipf("no phrase derived for %q, so nothing is sent", p)
			}

			for _, n := range names {
				if !re.MatchString(n) {
					continue
				}

				// The server matches case-insensitively, so compare that way.
				if !strings.Contains(
					strings.ToLower(n), strings.ToLower(phrase),
				) {
					t.Errorf(
						"name %q matches %q but does not contain the derived "+
							"phrase %q — the server would exclude it before "+
							"the filter block ran",
						n, p, phrase)
				}
			}
		})
	}
}
