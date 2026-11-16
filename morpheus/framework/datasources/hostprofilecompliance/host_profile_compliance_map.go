// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package hostprofilecompliance

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

// clusterComplianceAsState maps a cluster compliance response into the model.
func clusterComplianceAsState(
	ctx context.Context,
	config *HostProfileComplianceModel,
	r *sdk.GetClusterHostProfileCompliance200Response,
) (HostProfileComplianceModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	state := HostProfileComplianceModel{
		ClusterId:     config.ClusterId,
		OverallStatus: convert.StrToType(r.OverallStatus),
		ProfileId:     convert.Int64ToType(r.ProfileId),
		ProfileName:   convert.StrToType(r.ProfileName),
		// Server-only fields are not populated for a cluster lookup.
		ServerId:             types.Int64Null(),
		Compliant:            types.BoolNull(),
		Status:               types.StringNull(),
		LastChecked:          types.StringNull(),
		RemediatedDate:       types.StringNull(),
		RemediableViolations: types.Int64Null(),
		ManualViolations:     types.Int64Null(),
	}

	hosts, hostDiags := hostsAsSet(ctx, r.Hosts)
	diags.Append(hostDiags...)
	state.Hosts = hosts

	return state, diags
}

func hostsAsSet(
	ctx context.Context,
	hosts []sdk.GetClusterHostProfileCompliance200ResponseHostsInner,
) (types.Set, diag.Diagnostics) {
	elemType := HostsValue{}.Type(ctx)

	if len(hosts) == 0 {
		return types.SetNull(elemType), nil
	}

	elems := make([]HostsValue, 0, len(hosts))
	for i := range hosts {
		h := hosts[i]
		elems = append(elems, HostsValue{
			Compliant:            convert.BoolToType(h.Compliant),
			LastChecked:          convert.TimeToType(h.LastChecked.Get()),
			ManualViolations:     convert.Int32ToType(h.ManualViolations),
			RemediableViolations: convert.Int32ToType(h.RemediableViolations),
			RemediatedDate:       convert.TimeToType(h.RemediatedDate.Get()),
			ServerId:             convert.Int64ToType(h.ServerId),
			ServerName:           convert.StrToType(h.ServerName),
			Status:               convert.StrToType(h.Status),
			state:                attr.ValueStateKnown,
		})
	}

	return types.SetValueFrom(ctx, elemType, elems)
}

// serverComplianceAsState maps a server compliance response into the model. It
// returns diagnostics for symmetry with clusterComplianceAsState and so a
// future fallible conversion can surface errors without a signature change; the
// current conversions are all infallible, so the diagnostics are always empty.
func serverComplianceAsState(
	ctx context.Context,
	config *HostProfileComplianceModel,
	r *sdk.GetServerHostProfileCompliance200Response,
) (HostProfileComplianceModel, diag.Diagnostics) {
	return HostProfileComplianceModel{
		ServerId:             config.ServerId,
		ProfileId:            convert.Int64ToType(r.ProfileId),
		ProfileName:          convert.StrToType(r.ProfileName),
		Status:               convert.StrToType(r.Status),
		Compliant:            convert.BoolToType(r.Compliant),
		LastChecked:          convert.TimeToType(r.LastChecked.Get()),
		RemediatedDate:       convert.TimeToType(r.RemediatedDate.Get()),
		RemediableViolations: convert.Int32ToType(r.RemediableViolations),
		ManualViolations:     convert.Int32ToType(r.ManualViolations),
		// Cluster-only fields are not populated for a server lookup.
		ClusterId:     types.Int64Null(),
		OverallStatus: types.StringNull(),
		Hosts:         types.SetNull(HostsValue{}.Type(ctx)),
	}, nil
}
