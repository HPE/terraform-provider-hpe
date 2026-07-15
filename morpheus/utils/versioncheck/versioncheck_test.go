// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package versioncheck_test

import (
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/versioncheck"
)

func TestParseAndSatisfies(t *testing.T) {
	cases := []struct {
		name       string
		build      string
		constraint string
		want       bool
	}{
		{"four-segment build meets release constraint", "9.0.2.18", ">= 9.0.2", true},
		{"four-segment build meets its own build constraint", "9.0.2.18", ">= 9.0.2.18", true},
		{"four-segment build below build constraint", "9.0.2.17", ">= 9.0.2.18", false},
		{"three-segment equals release", "9.0.2", ">= 9.0.2", true},
		{"below release", "9.0.1", ">= 9.0.2", false},
		{"patch above release", "9.0.3", ">= 9.0.2", true},
		{"prerelease suffix is stripped", "9.0.2.18-hotfix1", ">= 9.0.2.18", true},
		{"surrounding whitespace is trimmed", " 9.0.2 ", ">= 9.0.2", true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			v, err := versioncheck.Parse(tc.build)
			if err != nil {
				t.Fatalf("Parse(%q) unexpected error: %v", tc.build, err)
			}

			got, err := versioncheck.Satisfies(v, tc.constraint)
			if err != nil {
				t.Fatalf("Satisfies(%q, %q) unexpected error: %v",
					tc.build, tc.constraint, err)
			}

			if got != tc.want {
				t.Errorf("Satisfies(%q, %q) = %v, want %v",
					tc.build, tc.constraint, got, tc.want)
			}
		})
	}
}

func TestParseInvalid(t *testing.T) {
	if _, err := versioncheck.Parse("not-a-version"); err == nil {
		t.Error(`Parse("not-a-version") expected an error, got nil`)
	}
}

func TestSatisfiesInvalidConstraint(t *testing.T) {
	v, err := versioncheck.Parse("9.0.2")
	if err != nil {
		t.Fatalf("Parse unexpected error: %v", err)
	}

	if _, err := versioncheck.Satisfies(v, "not-a-constraint"); err == nil {
		t.Error("Satisfies with invalid constraint expected an error, got nil")
	}
}
