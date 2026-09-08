// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package monitoringcheck

import (
	"testing"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

// TestMapGetResponseToModelNullOptionals guards the read-path nil-safety fix: an
// explicit JSON null (or an omitted field) for the optional check_interval and
// description must map to a Terraform null without panicking. The SDK decoder marks an
// explicit null as IsSet with a nil Get(), so the prior *Get() dereference panicked on
// import or refresh of a check that omitted these fields.
func TestMapGetResponseToModelNullOptionals(t *testing.T) {
	t.Parallel()

	tests := map[string]sdk.GetChecks200ResponseCheck{
		"explicit null optionals": {
			CheckInterval: *sdk.NewNullableInt64(nil),
			Description:   *sdk.NewNullableString(nil),
		},
		"absent optionals (zero value)": {},
	}

	for name, check := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var model MonitoringCheckModel
			mapGetResponseToModel(&model, &check) // must not panic

			if !model.CheckInterval.IsNull() {
				t.Errorf("check_interval = %v, want null", model.CheckInterval)
			}

			if !model.Description.IsNull() {
				t.Errorf("description = %v, want null", model.Description)
			}
		})
	}
}

// TestMapGetResponseToModelPopulatedOptionals confirms a present value still maps
// through, so the null guard did not over-null real data.
func TestMapGetResponseToModelPopulatedOptionals(t *testing.T) {
	t.Parallel()

	interval := int64(120)
	desc := "web check"
	check := sdk.GetChecks200ResponseCheck{
		CheckInterval: *sdk.NewNullableInt64(&interval),
		Description:   *sdk.NewNullableString(&desc),
	}

	var model MonitoringCheckModel
	mapGetResponseToModel(&model, &check)

	if model.CheckInterval.IsNull() || model.CheckInterval.ValueInt64() != interval {
		t.Errorf("check_interval = %v, want %d", model.CheckInterval, interval)
	}

	if model.Description.IsNull() || model.Description.ValueString() != desc {
		t.Errorf("description = %v, want %q", model.Description, desc)
	}
}
