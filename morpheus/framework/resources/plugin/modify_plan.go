// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package plugin

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// ConfigValidators enforces that exactly one of jar_path/jar_url is set. The
// generated schema already declares ExactlyOneOf on both attributes; this is a
// belt-and-braces resource-level validator kept for clarity and future-proofing.
func (r *Resource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.ExactlyOneOf(
			path.MatchRoot("jar_path"),
			path.MatchRoot("jar_url"),
		),
	}
}

// ModifyPlan validates jar_url when it is known at plan time: an http(s) URL
// with a host that is not a literal address the fetch policy refuses, and --
// for plain http -- paired with jar_sha256 (see validateJarURL). jar_sha256 is
// read from the configuration rather than the plan, where the plan modifier
// may have left it unknown.
func (r *Resource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() {
		// Destroy plan: nothing to validate.
		return
	}

	var plan, config PluginModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.JarUrl.IsNull() || plan.JarUrl.IsUnknown() {
		return
	}

	if err := validateJarURL(plan.JarUrl.ValueString(), !config.JarSha256.IsNull()); err != nil {
		resp.Diagnostics.AddAttributeError(path.Root("jar_url"), "invalid jar_url", err.Error())
	}
}
