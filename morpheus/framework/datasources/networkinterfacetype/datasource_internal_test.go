// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkinterfacetype

import (
	"testing"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

func strPtr(s string) *string { return &s }

func int64Ptr(i int64) *int64 { return &i }

func TestMatchNetworkInterfaceType(t *testing.T) {
	nicTypes := []sdk.ZoneNetworkOptionsResponseNetworkTypesInner{
		{Id: int64Ptr(10), Name: strPtr("E1000"), Code: strPtr("e1000")},
		{Id: int64Ptr(11), Name: strPtr("VMXNET 3"), Code: strPtr("vmxnet3")},
		{Id: int64Ptr(12), Name: strPtr("E1000e"), Code: strPtr("e1000e")},
	}

	tests := []struct {
		name       string
		lookupName string
		input      []sdk.ZoneNetworkOptionsResponseNetworkTypesInner
		wantID     int64
		wantCode   string
		wantErr    string
	}{
		{
			name:       "exact match returns id and code",
			lookupName: "VMXNET 3",
			input:      nicTypes,
			wantID:     11,
			wantCode:   "vmxnet3",
		},
		{
			name:       "no match errors",
			lookupName: "does-not-exist",
			input:      nicTypes,
			wantErr:    ErrorNoNetworkInterfaceTypeFound,
		},
		{
			name:       "empty list errors as not found",
			lookupName: "E1000",
			input:      nil,
			wantErr:    ErrorNoNetworkInterfaceTypeFound,
		},
		{
			name:       "match is case sensitive",
			lookupName: "vmxnet 3",
			input:      nicTypes,
			wantErr:    ErrorNoNetworkInterfaceTypeFound,
		},
		{
			name:       "multiple matches errors",
			lookupName: "E1000",
			input: []sdk.ZoneNetworkOptionsResponseNetworkTypesInner{
				{Id: int64Ptr(10), Name: strPtr("E1000"), Code: strPtr("e1000")},
				{Id: int64Ptr(99), Name: strPtr("E1000"), Code: strPtr("e1000-dup")},
			},
			wantErr: ErrorMultipleNetworkInterfaceTypes,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := matchNetworkInterfaceType(tt.input, tt.lookupName)

			if tt.wantErr != "" {
				if err == nil {
					t.Fatalf("expected error %q, got nil", tt.wantErr)
				}

				if err.Error() != tt.wantErr {
					t.Fatalf("expected error %q, got %q", tt.wantErr, err.Error())
				}

				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got.id == nil || *got.id != tt.wantID {
				t.Errorf("id: want %d, got %v", tt.wantID, got.id)
			}

			if got.code == nil || *got.code != tt.wantCode {
				t.Errorf("code: want %q, got %v", tt.wantCode, got.code)
			}
		})
	}
}
