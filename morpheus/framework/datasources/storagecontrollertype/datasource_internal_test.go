// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package storagecontrollertype

import (
	"testing"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

func strPtr(s string) *string { return &s }

func int64Ptr(i int64) *int64 { return &i }

// TestBuildControllerMountPoint pins the exact string format, which is the
// parity contract with hpegl. The -1:1:6:0 case is the canonical example from
// the plan and the instance documentation.
func TestBuildControllerMountPoint(t *testing.T) {
	tests := []struct {
		name       string
		busNumber  int64
		typeID     int64
		unitNumber int64
		want       string
	}{
		{
			name:       "scsi vmware paravirtual bus 1 unit 0",
			busNumber:  1,
			typeID:     6,
			unitNumber: 0,
			want:       "-1:1:6:0",
		},
		{
			name:       "bus 0 unit 0",
			busNumber:  0,
			typeID:     6,
			unitNumber: 0,
			want:       "-1:0:6:0",
		},
		{
			name:       "non-zero unit number",
			busNumber:  1,
			typeID:     104,
			unitNumber: 3,
			want:       "-1:1:104:3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildControllerMountPoint(tt.busNumber, tt.typeID, tt.unitNumber)
			if got != tt.want {
				t.Errorf("buildControllerMountPoint(%d, %d, %d) = %q, want %q",
					tt.busNumber, tt.typeID, tt.unitNumber, got, tt.want)
			}
		})
	}
}

func ct(
	id int64,
	name string,
	displayOrder int64,
	category string,
	maxDevices int64,
) sdk.ListProvisionTypes200ResponseAllOfProvisionTypesInnerControllerTypesInner {
	return sdk.ListProvisionTypes200ResponseAllOfProvisionTypesInnerControllerTypesInner{
		Id:           int64Ptr(id),
		Name:         strPtr(name),
		DisplayOrder: int64Ptr(displayOrder),
		Category:     strPtr(category),
		MaxDevices:   int64Ptr(maxDevices),
	}
}

func TestMatchControllerType(t *testing.T) {
	controllers := []sdk.ListProvisionTypes200ResponseAllOfProvisionTypesInnerControllerTypesInner{
		ct(4, "IDE", 1, "ide", 2),
		ct(6, "SCSI VMware Paravirtual", 2, "scsi", 16),
		ct(104, "SATA", 3, "sata", 30),
	}

	tests := []struct {
		name          string
		lookupName    string
		input         []sdk.ListProvisionTypes200ResponseAllOfProvisionTypesInnerControllerTypesInner
		wantID        int64
		wantCategory  string
		wantMaxDevice int64
		wantErr       string
	}{
		{
			name:          "exact match",
			lookupName:    "SCSI VMware Paravirtual",
			input:         controllers,
			wantID:        6,
			wantCategory:  "scsi",
			wantMaxDevice: 16,
		},
		{
			name:          "case-insensitive and whitespace-trimmed match (hpegl parity)",
			lookupName:    "  scsi vmware paravirtual  ",
			input:         controllers,
			wantID:        6,
			wantCategory:  "scsi",
			wantMaxDevice: 16,
		},
		{
			name:       "no match errors",
			lookupName: "NVMe",
			input:      controllers,
			wantErr:    ErrorNoStorageControllerTypeFound,
		},
		{
			name:       "empty list errors as not found",
			lookupName: "SCSI VMware Paravirtual",
			input:      nil,
			wantErr:    ErrorNoStorageControllerTypeFound,
		},
		{
			name:       "multiple matches errors",
			lookupName: "SCSI",
			input: []sdk.ListProvisionTypes200ResponseAllOfProvisionTypesInnerControllerTypesInner{
				ct(6, "scsi", 2, "scsi", 16),
				ct(7, "SCSI", 5, "scsi", 16),
			},
			wantErr: ErrorMultipleStorageControllerTypes,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := matchControllerType(tt.input, tt.lookupName)

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

			if got.id != tt.wantID {
				t.Errorf("id: want %d, got %d", tt.wantID, got.id)
			}

			if got.category != tt.wantCategory {
				t.Errorf("category: want %q, got %q", tt.wantCategory, got.category)
			}

			if got.maxDevices != tt.wantMaxDevice {
				t.Errorf("maxDevices: want %d, got %d", tt.wantMaxDevice, got.maxDevices)
			}
		})
	}
}
