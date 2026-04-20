// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package instance

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/cenkalti/backoff/v5"
	semver "github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/constants"
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
	resp.TypeName = req.ProviderTypeName + "_" + constants.SubProviderName + "_instance"
	resp.TypeName = strings.Join(
		[]string{req.ProviderTypeName, constants.SubProviderName, "instance"},
		"_",
	)
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

	// Get appliance version from the health API
	apiResp, hresp, err := client.HealthAPI.ListHealth(ctx).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		// If we can't get a health response, skip validation
		return
	}

	health, ok := apiResp.GetHealthOk()
	if !ok {
		return
	}

	apiVersion, ok := health.GetBuildVersionOk()
	if !ok {
		return
	}

	versionParts := strings.Split(*apiVersion, ".")
	if len(versionParts) < 3 {
		resp.Diagnostics.AddError(
			"Unable to Parse Appliance Version",
			"Not enough components - expect at least major.minor.patch, got %s"+
				*apiVersion,
		)
	}

	trimmedVersion := strings.Join(versionParts[:3], ".")

	version, err := semver.NewVersion(trimmedVersion)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Parse Appliance Version",
			fmt.Sprintf(
				"Unable to parse appliance version %s: %v",
				*apiVersion, err.Error(),
			),
		)

		return
	}

	// Check for network updates that require a resource replacement
	networkConstraint := &morpheusConstraint{
		morphVersion: version,
		plan:         plan.NetworkInterfaces,
		state:        state.NetworkInterfaces,
		constraint:   ">= 8.1.2",
		hclAttribute: "network_interfaces",
		mnemonic:     "network",
	}
	networkConstraint.checkForAttributeUpdate(resp)

	// Check for service plan options updates that require a resource replacement
	servicePlanOptionsConstraint := &morpheusConstraint{
		morphVersion: version,
		plan:         plan.ServicePlanOptions,
		state:        state.ServicePlanOptions,
		constraint:   ">= 8.1.2",
		hclAttribute: "service_plan_options",
		mnemonic:     "Service Plan Options",
	}
	servicePlanOptionsConstraint.checkForAttributeUpdate(resp)
}

type morpheusConstraint struct {
	morphVersion *semver.Version
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
	versionCheck, err := semver.NewConstraint(m.constraint)
	if err != nil {
		return
	}

	if versionCheck.Check(m.morphVersion) {
		return
	}

	// Attribute shape changed on an older appliance version;
	// force a new resource.
	resp.RequiresReplace = append(
		resp.RequiresReplace, path.Root(m.hclAttribute),
	)
	resp.Diagnostics.AddWarning(fmt.Sprintf("%s change will trigger replace",
		cases.Title(language.English, cases.NoLower).String(m.mnemonic)),
		fmt.Sprintf("Morpheus version must be %s to allow %s updates without instance replacement",
			m.constraint, m.mnemonic),
	)
}
