// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package instance

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/cenkalti/backoff/v5"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/versioncheck"
)

var (
	_ resource.Resource                = &Resource{}
	_ resource.ResourceWithImportState = &Resource{}
	_ resource.ResourceWithModifyPlan  = &Resource{}
)

func NewResource() resource.Resource {
	return &Resource{}
}

type Resource struct {
	configure.ResourceWithMorpheusConfigure
}

// Metadata implements resource.Resource.
func (g *Resource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_" + "instance"
	resp.TypeName = req.ProviderTypeName + "_" + "instance"
}

// Schema implements resource.Resource.
func (g *Resource) Schema(
	ctx context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = InstanceResourceSchema(ctx)
}

func checkStatusDone(status string, targetStatuses []string, errorStatuses []string) error {
	switch {
	case slices.Contains(errorStatuses, status):
		return backoff.Permanent(errors.New("reached error status: " + status))
	case slices.Contains(targetStatuses, status):
		return nil
	default:
		return backoff.RetryAfter(5)
	}
}

func (g *Resource) ModifyPlan(
	ctx context.Context,
	req resource.ModifyPlanRequest,
	resp *resource.ModifyPlanResponse,
) {
	// Only run validation if we have a plan and a state
	if req.Plan.Raw.IsNull() || req.State.Raw.IsNull() {
		return
	}

	// Put prior state back into the plan for computed attributes whose
	// triggering configuration has not changed, so that an unrelated edit does
	// not show them as "(known after apply)". Done before the checks below,
	// which return early when the appliance version cannot be determined.
	restoreUnchangedComputedAttributes(ctx, req, resp)
	if resp.Diagnostics.HasError() {
		return
	}

	var plan, state InstanceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get API client - provider is configured at this point
	client, err := g.NewClient(ctx)
	if err != nil {
		// If we can't get a client, skip validation
		return
	}

	// Determine the appliance version to decide whether attribute-shape changes
	// require a resource replacement. Skip the check (rather than block the plan)
	// if the version cannot be determined.
	morphVersion, err := versioncheck.Appliance(ctx, client)
	if err != nil {
		return
	}

	// Check for network updates that require a resource replacement
	networkConstraint := &morpheusConstraint{
		morphVersion: morphVersion,
		plan:         plan.NetworkInterfaces,
		state:        state.NetworkInterfaces,
		constraint:   ">= 8.1.2",
		hclAttribute: "network_interfaces",
		mnemonic:     "network",
	}
	networkConstraint.checkForAttributeUpdate(resp)

	// Check for service plan options updates that require a resource replacement
	servicePlanOptionsConstraint := &morpheusConstraint{
		morphVersion: morphVersion,
		plan:         plan.ServicePlanOptions,
		state:        state.ServicePlanOptions,
		constraint:   ">= 8.1.2",
		hclAttribute: "service_plan_options",
		mnemonic:     "Service Plan Options",
	}
	servicePlanOptionsConstraint.checkForAttributeUpdate(resp)
}

type morpheusConstraint struct {
	morphVersion *versioncheck.Version
	plan         attr.Value
	state        attr.Value
	constraint   string
	hclAttribute string
	mnemonic     string
}

func (m *morpheusConstraint) checkForAttributeUpdate(
	resp *resource.ModifyPlanResponse,
) {
	// Has there been a change in the "shape" of the attribute?
	if m.plan.Equal(m.state) {
		return
	}

	// Build constraint
	ok, err := versioncheck.Satisfies(m.morphVersion, m.constraint)
	if err != nil || ok {
		return
	}

	// Attribute shape changed on an older appliance version;
	// force a new resource.
	resp.RequiresReplace = append(
		resp.RequiresReplace, path.Root(m.hclAttribute),
	)
	resp.Diagnostics.AddWarning(
		fmt.Sprintf("%s change will trigger replace",
			cases.Title(language.English, cases.NoLower).String(m.mnemonic)),
		fmt.Sprintf("Morpheus version must be %s to allow %s updates without instance replacement",
			m.constraint, m.mnemonic),
	)
}
