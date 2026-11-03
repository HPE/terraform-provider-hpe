// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package cluster

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/versioncheck"
)

// requireHvm2Layout refuses the operation when any HVM 2.0 network field is set
// but the selected cluster layout's clusterVersion does not satisfy the 2.0
// minimum. It fails open when the layout version cannot be determined: an
// unreadable layout version must not block an otherwise healthy operation
// (mirrors versioncheck.decide's fail-open rationale).
//
// It is the I/O half of the gate: it fetches the layout, then hands the pure
// decision to decideHvm2Layout.
func requireHvm2Layout(
	ctx context.Context,
	client *sdk.APIClient,
	layoutId int64,
	cfg ConfigHvmValue,
) diag.Diagnostics {
	var diags diag.Diagnostics

	anySet := hvm2FieldsSet(cfg)
	if !anySet {
		// No 2.0 fields configured — nothing to gate.
		return diags
	}

	resp, _, err := client.ClusterLayoutsAPI.GetClusterLayout(ctx, layoutId).Execute()
	if err != nil || resp == nil || resp.Layout == nil || resp.Layout.ClusterVersion == nil {
		reason := "cluster layout did not report a clusterVersion"
		if err != nil {
			reason = err.Error()
		}

		tflog.Warn(ctx, "skipping HVM 2.0 layout-version gate: version unavailable", map[string]any{
			"layout_id": layoutId,
			"reason":    reason,
		})

		return diags
	}

	refuse, detail := decideHvm2Layout(ctx, anySet, resp.Layout.ClusterVersion)
	if refuse {
		diags.AddAttributeError(
			path.Root("config_hvm"),
			"HVM 2.0 network configuration requires a newer cluster layout",
			detail,
		)
	}

	return diags
}

// decideHvm2Layout is the pure decision half of requireHvm2Layout, split out so
// the gate's behaviour can be unit tested without a live API. It returns
// refuse=true only when 2.0 fields are set AND the layout's clusterVersion is
// known and fails the 2.0 minimum. It fails open (refuse=false) when the
// version is nil or cannot be parsed/compared.
func decideHvm2Layout(
	ctx context.Context,
	anySet bool,
	clusterVersion *string,
) (refuse bool, detail string) {
	if !anySet {
		return false, ""
	}

	if clusterVersion == nil {
		tflog.Warn(ctx, "skipping HVM 2.0 layout-version gate: layout has no clusterVersion", nil)

		return false, ""
	}

	v, err := versioncheck.Parse(*clusterVersion)
	if err != nil {
		tflog.Warn(ctx, "skipping HVM 2.0 layout-version gate: unparseable clusterVersion", map[string]any{
			"cluster_version": *clusterVersion,
			"reason":          err.Error(),
		})

		return false, ""
	}

	ok, err := versioncheck.Satisfies(v, ClusterHvmVirtSwitchMinLayoutVersion)
	if err != nil {
		tflog.Warn(ctx, "skipping HVM 2.0 layout-version gate: malformed constraint", map[string]any{
			"constraint": ClusterHvmVirtSwitchMinLayoutVersion,
			"reason":     err.Error(),
		})

		return false, ""
	}

	if ok {
		return false, ""
	}

	return true, fmt.Sprintf(
		"The config_hvm network fields (virt_switch_name, dns1, dns2, dns3, "+
			"search_domain, ntp1, ntp2, proxy, no_proxy) are only honoured by "+
			"Morpheus when the selected cluster layout has clusterVersion 2.0 "+
			"or greater. The selected layout reports clusterVersion %q. Remove "+
			"these fields, or select a cluster layout with clusterVersion 2.0 "+
			"or greater.",
		*clusterVersion,
	)
}

// hvm2FieldsSet reports whether any of the nine HVM 2.0 network fields is set
// (known and non-null) on the given config.
func hvm2FieldsSet(cfg ConfigHvmValue) bool {
	fields := []interface {
		IsNull() bool
		IsUnknown() bool
	}{
		cfg.VirtSwitchName,
		cfg.Dns1,
		cfg.Dns2,
		cfg.Dns3,
		cfg.SearchDomain,
		cfg.Ntp1,
		cfg.Ntp2,
		cfg.Proxy,
		cfg.NoProxy,
	}

	for _, f := range fields {
		if !f.IsNull() && !f.IsUnknown() {
			return true
		}
	}

	return false
}
