// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

// Package tenancy provides shared plan-time guards for attributes that Morpheus
// only honours for the master tenant. Morpheus silently coerces or discards
// these values for sub-tenant callers rather than returning an error, so the
// provider fails fast (or, when the caller's tenancy cannot be determined,
// fails open with a warning) instead of leaving state disagreeing with config.
package tenancy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// MasterChecker reports whether the configured caller is the master tenant. It
// is satisfied by the embedded configure structs
// (ResourceWithMorpheusConfigure / DataSourceWithMorpheusConfigure).
type MasterChecker interface {
	CallerIsMaster(ctx context.Context) (bool, error)
}

// visibilityPublicDetail is the shared explanation appended to the visibility
// guard diagnostics.
const visibilityPublicDetail = "Morpheus silently stores \"private\" for objects " +
	"created or updated by sub-tenant users; set visibility = \"private\" or run " +
	"as a master-tenant user."

// GuardVisibilityPublic rejects, at plan time, visibility = "public" set by a
// non-master caller. It only calls whoami when visibility is a known "public"
// value, so the common (non-public) case incurs no API call. On an error
// determining tenancy it fails open with a warning; a confirmed non-master
// caller gets a blocking error.
func GuardVisibilityPublic(
	ctx context.Context,
	mc MasterChecker,
	plan tfsdk.Plan,
	diags *diag.Diagnostics,
) {
	var visibility types.String
	diags.Append(plan.GetAttribute(ctx, path.Root("visibility"), &visibility)...)
	if diags.HasError() {
		return
	}

	if visibility.IsNull() || visibility.IsUnknown() {
		return
	}
	if visibility.ValueString() != "public" {
		return
	}

	isMaster, err := mc.CallerIsMaster(ctx)
	if err != nil {
		diags.AddAttributeWarning(
			path.Root("visibility"),
			"Could not validate visibility at plan time",
			"Unable to determine whether the authenticated tenant is the master "+
				"tenant, so visibility was not validated; Morpheus will enforce it "+
				"on apply. Error: "+err.Error(),
		)

		return
	}

	if !isMaster {
		diags.AddAttributeError(
			path.Root("visibility"),
			"visibility = \"public\" requires the master tenant",
			visibilityPublicDetail,
		)
	}
}

// GuardMasterOnlyAttribute rejects, at plan time, any known non-null value for
// attr set by a non-master caller. It only calls whoami when attr has a known,
// non-null value. On an error determining tenancy it fails open with a warning.
func GuardMasterOnlyAttribute(
	ctx context.Context,
	mc MasterChecker,
	plan tfsdk.Plan,
	attr path.Path,
	diags *diag.Diagnostics,
) {
	var value types.String
	diags.Append(plan.GetAttribute(ctx, attr, &value)...)
	if diags.HasError() {
		return
	}

	if value.IsNull() || value.IsUnknown() {
		return
	}

	isMaster, err := mc.CallerIsMaster(ctx)
	if err != nil {
		diags.AddAttributeWarning(
			attr,
			"Could not validate a master-tenant-only attribute at plan time",
			"Unable to determine whether the authenticated tenant is the master "+
				"tenant, so this attribute was not validated; Morpheus will enforce "+
				"it on apply. Error: "+err.Error(),
		)

		return
	}

	if !isMaster {
		diags.AddAttributeError(
			attr,
			attr.String()+" can only be set by the master tenant",
			attr.String()+" can only be set by the master tenant; Morpheus "+
				"silently discards it for sub-tenant users.",
		)
	}
}

// CheckVisibilityApplied is the apply-time fallback for visibility. When the
// plan asked for a known "public" value but the API stored "private" (the
// sub-tenant coercion), it returns a blocking error diagnostic with the same
// wording as the plan-time guard, instead of letting Terraform emit an opaque
// "inconsistent result" error. It returns nil otherwise.
func CheckVisibilityApplied(planned, actual types.String) diag.Diagnostic {
	if planned.IsNull() || planned.IsUnknown() {
		return nil
	}
	if planned.ValueString() != "public" {
		return nil
	}
	if actual.ValueString() == "public" {
		return nil
	}

	return diag.NewAttributeErrorDiagnostic(
		path.Root("visibility"),
		"visibility = \"public\" requires the master tenant",
		visibilityPublicDetail,
	)
}
