// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package serviceplan

import (
	"testing"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
)

// TestServicePlanInCloud verifies the cloud_id disambiguation helper used by the
// data source. ListServicePlans with IncludeZones(true) returns, for each plan,
// the clouds (zones) it is available in; servicePlanInZone reports whether a
// given cloud id is among them.
func TestServicePlanInCloud(t *testing.T) {
	t.Parallel()

	newPlan := func(zoneIDs ...int64) sdk.ListServicePlans200ResponseAllOfServicePlansInner {
		zones := make(
			[]sdk.ListServicePlans200ResponseAllOfServicePlansInnerZonesInner, 0, len(zoneIDs),
		)
		for _, id := range zoneIDs {
			zones = append(zones, sdk.ListServicePlans200ResponseAllOfServicePlansInnerZonesInner{
				Id: sdk.PtrInt64(id),
			})
		}

		return sdk.ListServicePlans200ResponseAllOfServicePlansInner{Zones: zones}
	}

	tests := []struct {
		name    string
		plan    sdk.ListServicePlans200ResponseAllOfServicePlansInner
		cloudID int64
		want    bool
	}{
		{"cloud present in zones", newPlan(3, 5, 7), 5, true},
		{"cloud absent from zones", newPlan(3, 7), 5, false},
		{"plan has no zones", newPlan(), 5, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := servicePlanInZone(tt.plan, tt.cloudID); got != tt.want {
				t.Errorf("servicePlanInZone(%d) = %v, want %v", tt.cloudID, got, tt.want)
			}
		})
	}
}
