// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package storagevolume

import "testing"

func TestMaxStorageSizeError(t *testing.T) {
	cases := []struct {
		name          string
		size          int64
		typeCode      string
		typeCodeKnown bool
		wantErr       bool
	}{
		{"zero is invalid", 0, "", false, true},
		{"negative is invalid", -1, "hpealletraMPLUN", true, true},
		{"min boundary 1 is valid", 1, "hpealletraMPLUN", true, false},
		{"alletra in range", 100, "hpealletraMPLUN", true, false},
		{"alletra at max is valid", 65536, "hpealletraMPLUN-active-pp", true, false},
		{"alletra over max is invalid", 65537, "hpealletraMPLUN", true, true},
		{"alletra classic over max is invalid", 70000, "hpealletraMPLUN-classic-pp", true, true},
		{"non-alletra (3par) over alletra max is valid", 70000, "3par", true, false},
		{"unknown type_code does not apply alletra cap", 70000, "", false, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			title, _ := maxStorageSizeError(tc.size, tc.typeCode, tc.typeCodeKnown)
			if got := title != ""; got != tc.wantErr {
				t.Fatalf("maxStorageSizeError(%d, %q, %v) error=%v, want %v",
					tc.size, tc.typeCode, tc.typeCodeKnown, got, tc.wantErr)
			}
		})
	}
}
