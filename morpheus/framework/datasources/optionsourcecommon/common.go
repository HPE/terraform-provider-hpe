// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

// Package optionsourcecommon holds what the hpe_morpheus_option_source_value and
// hpe_morpheus_option_source_values data sources share: the input attributes
// and their model, turning a configuration into a request, running it with the
// right diagnostics, and rendering values into Terraform objects. Each data
// source adds only its own outputs and its own cardinality rule.
package optionsourcecommon

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/optionsource"
)

// InputsModel is the input half of both data sources' models. Embed it by value.
type InputsModel struct {
	Source            types.String `tfsdk:"source"`
	ServiceType       types.String `tfsdk:"service_type"`
	Parameters        types.Map    `tfsdk:"parameters"`
	ParameterStyle    types.String `tfsdk:"parameter_style"`
	Collection        types.String `tfsdk:"collection"`
	CloudID           types.Int64  `tfsdk:"cloud_id"`
	Phrase            types.String `tfsdk:"phrase"`
	ParentID          types.Int64  `tfsdk:"parent_id"`
	CredentialID      types.Int64  `tfsdk:"credential_id"`
	GroupID           types.Int64  `tfsdk:"group_id"`
	LayoutID          types.Int64  `tfsdk:"layout_id"`
	ProvisionTypeCode types.String `tfsdk:"provision_type_code"`
	ProvisionTypeID   types.Int64  `tfsdk:"provision_type_id"`
	PoolID            types.String `tfsdk:"pool_id"`
	PlanID            types.Int64  `tfsdk:"plan_id"`
	ServerGroupID     types.Int64  `tfsdk:"server_group_id"`
	ProviderID        types.Int64  `tfsdk:"provider_id"`
}

// Description is the shared preamble for both data sources' descriptions. It is
// deliberately short: the generated documentation page carries the narrative,
// and this text also appears in schema dumps and editor tooltips.
const Description = "Morpheus option sources are the lookups behind the UI's dropdowns — the resource " +
	"pools and networks the provisioning wizard offers, an SCVMM library share, the members of a " +
	"select list — served at `GET /api/options/{source}`."

// InputAttributes returns the input attributes shared by both data sources.
func InputAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"source": schema.StringAttribute{
			Required: true,
			Description: "The option source name, as the API names it, e.g. `zonePools`, `zoneNetworkOptions`, " +
				"`groups`, `scvmmLibraryShares`. A name not in the provider's catalogue of sources known for " +
				"Morpheus " + optionsource.KnownSourcesRelease + " produces a warning, not an error: plugins and " +
				"newer appliances add sources.",
			MarkdownDescription: "The option source name, as the API names it, e.g. `zonePools`, `zoneNetworkOptions`, " +
				"`groups`, `scvmmLibraryShares`. A name not in the provider's catalogue of sources known for " +
				"Morpheus " + optionsource.KnownSourcesRelease + " produces a warning, not an error: plugins and " +
				"newer appliances add sources.",
			Validators: []validator.String{stringvalidator.LengthAtLeast(1)},
		},
		"service_type": schema.StringAttribute{
			Optional: true,
			Description: "Selects the two-segment path `/api/options/{service_type}/{source}`, which is how a " +
				"source defined on a cloud-specific service is reached — for example `scvmm` for " +
				"`scvmmLibraryShares`. A source that needs one and is called without it fails with a generic " +
				"server error; the provider warns when the catalogue says so.",
			MarkdownDescription: "Selects the two-segment path `/api/options/{service_type}/{source}`, which is how a " +
				"source defined on a cloud-specific service is reached — for example `scvmm` for " +
				"`scvmmLibraryShares`. A source that needs one and is called without it fails with a generic " +
				"server error; the provider warns when the catalogue says so.",
		},
		"parameters": schema.MapAttribute{
			ElementType: types.StringType,
			Optional:    true,
			Description: "Parameters for the source, keyed as the API names them (`zoneId`, `layoutId`, ...). " +
				"Morpheus reads parameters from two places — the query string and `config[<key>]` — and a " +
				"parameter in the wrong place is silently ignored rather than rejected, so by default each key " +
				"is sent in both (see `parameter_style`). The typed attributes below are shortcuts for the common " +
				"keys; a key given both ways is an error.",
			MarkdownDescription: "Parameters for the source, keyed as the API names them (`zoneId`, `layoutId`, ...). " +
				"Morpheus reads parameters from two places — the query string and `config[<key>]` — and a " +
				"parameter in the wrong place is silently ignored rather than rejected, so by default each key " +
				"is sent in both (see `parameter_style`). The typed attributes below are shortcuts for the common " +
				"keys; a key given both ways is an error.",
		},
		"parameter_style": schema.StringAttribute{
			Optional: true,
			Description: "Where each parameter is placed: `auto` (the default) sends it both in the query string and " +
				"as `config[<key>]`; `top_level` and `config` send it in one place only. Use them only for a " +
				"source that misbehaves when a parameter is duplicated.",
			MarkdownDescription: "Where each parameter is placed: `auto` (the default) sends it both in the query string and " +
				"as `config[<key>]`; `top_level` and `config` send it in one place only. Use them only for a " +
				"source that misbehaves when a parameter is duplicated.",
			Validators: []validator.String{stringvalidator.OneOf(optionsource.Styles()...)},
		},
		"collection": schema.StringAttribute{
			Optional: true,
			Description: "For a source that returns an object of lists rather than a single list — " +
				"`zoneNetworkOptions` returns `networks`, `networkGroups`, `networkSubnets` and `networkTypes` — " +
				"the list to read. " +
				"Required for such a source; the error names the lists available. Not allowed for a source " +
				"that returns a list.",
			MarkdownDescription: "For a source that returns an object of lists rather than a single list — " +
				"`zoneNetworkOptions` returns `networks`, `networkGroups`, `networkSubnets` and `networkTypes` — " +
				"the list to read. " +
				"Required for such a source; the error names the lists available. Not allowed for a source " +
				"that returns a list.",
		},
		"cloud_id": schema.Int64Attribute{
			Optional: true,
			Description: "The cloud to evaluate the source in. Sent as `zoneId`, the API's name for a cloud (and as " +
				"`cloudId`, which some newer sources read instead).",
			MarkdownDescription: "The cloud to evaluate the source in. Sent as `zoneId`, the API's name for a cloud (and as " +
				"`cloudId`, which some newer sources read instead).",
		},
		"phrase": schema.StringAttribute{
			Optional:            true,
			Description:         "A search phrase, for sources that narrow by name server-side. Sent as `phrase`.",
			MarkdownDescription: "A search phrase, for sources that narrow by name server-side. Sent as `phrase`.",
		},
		"parent_id": schema.Int64Attribute{
			Optional: true,
			Description: "The id of the object a dependent source is scoped to — the load balancer for its pools, " +
				"the network server for its objects. Sent as `parentId` (and `parentid`, which some sources read).",
			MarkdownDescription: "The id of the object a dependent source is scoped to — the load balancer for its pools, " +
				"the network server for its objects. Sent as `parentId` (and `parentid`, which some sources read).",
		},
		"credential_id": schema.Int64Attribute{
			Optional: true,
			Description: "The id of an `hpe_morpheus_credential` for sources that authenticate to a target system " +
				"(SCVMM, Amazon, OpenStack). Sent as `credential.type`. This is the route for secrets; never put " +
				"them in `parameters`.",
			MarkdownDescription: "The id of an `hpe_morpheus_credential` for sources that authenticate to a target system " +
				"(SCVMM, Amazon, OpenStack). Sent as `credential.type`. This is the route for secrets; never put " +
				"them in `parameters`.",
		},
		"group_id": schema.Int64Attribute{
			Optional: true,
			Description: "The group (site) to evaluate the source for. Sent as both `groupId` and `siteId`, since " +
				"sources read one or the other.",
			MarkdownDescription: "The group (site) to evaluate the source for. Sent as both `groupId` and `siteId`, since " +
				"sources read one or the other.",
		},
		"layout_id": schema.Int64Attribute{
			Optional: true,
			Description: "The instance type layout that gives a provisioning source its context — the provisioning " +
				"wizard's sources (`zonePools`, `zoneNetworkOptions`) return an empty list without one. Sent as " +
				"`layoutId`.",
			MarkdownDescription: "The instance type layout that gives a provisioning source its context — the provisioning " +
				"wizard's sources (`zonePools`, `zoneNetworkOptions`) return an empty list without one. Sent as " +
				"`layoutId`.",
		},
		"provision_type_code": schema.StringAttribute{
			Optional: true,
			Description: "The provision type code (`kvm`, `vmware`, ...) as an alternative to `layout_id`. Sent as " +
				"both `provisionTypeCode` and `provisionType`, since sources read one or the other.",
			MarkdownDescription: "The provision type code (`kvm`, `vmware`, ...) as an alternative to `layout_id`. Sent as " +
				"both `provisionTypeCode` and `provisionType`, since sources read one or the other.",
		},
		"pool_id": schema.StringAttribute{
			Optional: true,
			Description: "The resource pool to scope a provisioning source to, as `pool-<id>` or a bare id — the " +
				"form `zonePools` returns as `value`. Sent as `poolId`.",
			MarkdownDescription: "The resource pool to scope a provisioning source to, as `pool-<id>` or a bare id — the " +
				"form `zonePools` returns as `value`. Sent as `poolId`.",
		},
		"provision_type_id": schema.Int64Attribute{
			Optional: true,
			Description: "The provision type by id, for the sources that read the id rather than the code " +
				"(`zonePools` reads both; `storageControllerTypes` reads only the id). Sent as `provisionTypeId`.",
			MarkdownDescription: "The provision type by id, for the sources that read the id rather than the code " +
				"(`zonePools` reads both; `storageControllerTypes` reads only the id). Sent as `provisionTypeId`.",
		},
		"plan_id": schema.Int64Attribute{
			Optional: true,
			Description: "The service plan, for sources whose options depend on it (plan-scoped sizing and " +
				"storage options). Sent as `planId`.",
			MarkdownDescription: "The service plan, for sources whose options depend on it (plan-scoped sizing and " +
				"storage options). Sent as `planId`.",
		},
		"server_group_id": schema.Int64Attribute{
			Optional: true,
			Description: "The cluster (server group), for sources scoped to one — its VMs, hosts, affinity groups, " +
				"datastores. Sent as `serverGroupId`.",
			MarkdownDescription: "The cluster (server group), for sources scoped to one — its VMs, hosts, affinity groups, " +
				"datastores. Sent as `serverGroupId`.",
		},
		"provider_id": schema.Int64Attribute{
			Optional: true,
			Description: "The backup provider, for the backup and replication sources scoped to one. Sent as " +
				"`providerId`.",
			MarkdownDescription: "The backup provider, for the backup and replication sources scoped to one. Sent as " +
				"`providerId`.",
		},
	}
}

// ValueAttrTypes is the attribute type of one option, shared by the plural
// list element and the singular data source's outputs.
func ValueAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":       types.StringType,
		"value":      types.StringType,
		"attributes": types.MapType{ElemType: types.StringType},
	}
}

// ValueAttributes returns the schema of one option, for the plural list element.
func ValueAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"name": schema.StringAttribute{
			Computed:            true,
			Description:         "The display label.",
			MarkdownDescription: "The display label.",
		},
		"value": schema.StringAttribute{
			Computed: true,
			Description: "The value to use in configuration — the source's `value`, or its `id` when it has no " +
				"`value` (as `zoneNetworkOptions` items do). Already in the form the API takes on input, e.g. " +
				"`pool-2`, `network-79`.",
			MarkdownDescription: "The value to use in configuration — the source's `value`, or its `id` when it has no " +
				"`value` (as `zoneNetworkOptions` items do). Already in the form the API takes on input, e.g. " +
				"`pool-2`, `network-79`.",
		},
		"attributes": schema.MapAttribute{
			ElementType: types.StringType,
			Computed:    true,
			Description: "Every other key the source returned for the option, as strings (numbers and booleans " +
				"in their textual form, nested objects as JSON). Source-specific: `zonePools` returns " +
				"`providerType` and `type`, `zoneNetworkOptions` returns `dhcpServer` and `pool`.",
			MarkdownDescription: "Every other key the source returned for the option, as strings (numbers and booleans " +
				"in their textual form, nested objects as JSON). Source-specific: `zonePools` returns " +
				"`providerType` and `type`, `zoneNetworkOptions` returns `dhcpServer` and `pool`.",
		},
	}
}

// Request turns the configured inputs into a request. Its diagnostics are
// configuration errors: an unknown parameter_style, or a key set both typed and
// in parameters.
func (m InputsModel) Request(ctx context.Context) (optionsource.Request, diag.Diagnostics) {
	var diags diag.Diagnostics

	style, err := optionsource.ParseStyle(m.ParameterStyle.ValueString())
	if err != nil {
		diags.AddAttributeError(attrPath("parameter_style"), "invalid parameter_style", err.Error())

		return optionsource.Request{}, diags
	}

	params := map[string]string{}

	if !m.Parameters.IsNull() && !m.Parameters.IsUnknown() {
		diags.Append(m.Parameters.ElementsAs(ctx, &params, false)...)
		if diags.HasError() {
			return optionsource.Request{}, diags
		}
	}

	merged, err := optionsource.MergeParameters(params, optionsource.Typed{
		CloudID:           m.CloudID.ValueInt64Pointer(),
		Phrase:            m.Phrase.ValueStringPointer(),
		ParentID:          m.ParentID.ValueInt64Pointer(),
		CredentialID:      m.CredentialID.ValueInt64Pointer(),
		GroupID:           m.GroupID.ValueInt64Pointer(),
		LayoutID:          m.LayoutID.ValueInt64Pointer(),
		ProvisionTypeCode: m.ProvisionTypeCode.ValueStringPointer(),
		ProvisionTypeID:   m.ProvisionTypeID.ValueInt64Pointer(),
		PoolID:            m.PoolID.ValueStringPointer(),
		PlanID:            m.PlanID.ValueInt64Pointer(),
		ServerGroupID:     m.ServerGroupID.ValueInt64Pointer(),
		ProviderID:        m.ProviderID.ValueInt64Pointer(),
	})
	if err != nil {
		diags.AddAttributeError(attrPath("parameters"), "conflicting parameters", err.Error())

		return optionsource.Request{}, diags
	}

	// A parameter the server overwrites is not an error — the call succeeds —
	// but it is silently not what the author asked for, which is exactly the
	// class of mistake these data sources exist to surface.
	for _, key := range optionsource.ServerSetParameters(merged) {
		diags.AddAttributeWarning(attrPath("parameters"), "parameter "+key+" is set by the server",
			optionsource.ServerSetAdvice(key))
	}

	return optionsource.Request{
		Source:      m.Source.ValueString(),
		ServiceType: m.ServiceType.ValueString(),
		Parameters:  merged,
		Style:       style,
		Collection:  m.Collection.ValueString(),
	}, diags
}

// Fetch runs the request, adding the catalogue's advice as a warning and
// translating a failure into one error diagnostic. It returns nil when the
// call failed; the diagnostics say why.
func Fetch(
	ctx context.Context,
	client *sdk.APIClient,
	req optionsource.Request,
	summary string,
	diags *diag.Diagnostics,
) *optionsource.Result {
	if advice := optionsource.Advice(req.Source, req.ServiceType); advice != "" {
		diags.AddAttributeWarning(attrPath("source"), "option source not in the catalogue", advice)
	}

	res, err := optionsource.Fetch(ctx, client, req)
	if err == nil {
		return res
	}

	var serverErr *optionsource.ServerError

	detail := err.Error()

	if errors.As(err, &serverErr) {
		detail = serverErr.Describe()

		if serverErr.Gasket() {
			// The server cannot say which; the catalogue can narrow it.
			switch {
			case !optionsource.Known(req.Source):
				detail += fmt.Sprintf(" The source is not in the provider's catalogue for Morpheus %s, so the "+
					"name is the most likely cause.", optionsource.KnownSourcesRelease)
			case optionsource.Advice(req.Source, req.ServiceType) != "":
				detail += " The catalogue's advice on service_type (see the warning) is the most likely cause."
			default:
				detail += fmt.Sprintf(" The source is known, so a missing or malformed parameter is the most "+
					"likely cause. Parameters sent: %s.", optionsource.DescribeParameters(req.Parameters))
			}
		}

		detail += "\n\n" + serverErr.RequestLine()
	}

	diags.AddError(summary, detail)

	return nil
}

// ValueObject renders one option as a Terraform object.
func ValueObject(v optionsource.Value) (types.Object, diag.Diagnostics) {
	attrsMap, diags := AttributesMap(v)
	if diags.HasError() {
		return types.ObjectNull(ValueAttrTypes()), diags
	}

	return types.ObjectValue(ValueAttrTypes(), map[string]attr.Value{
		"name":       types.StringValue(v.Name),
		"value":      types.StringValue(v.Value),
		"attributes": attrsMap,
	})
}

// AttributesMap renders an option's attributes as a map of strings.
func AttributesMap(v optionsource.Value) (types.Map, diag.Diagnostics) {
	elems := make(map[string]attr.Value, len(v.Attributes))
	for k, s := range v.Attributes {
		elems[k] = types.StringValue(s)
	}

	return types.MapValue(types.StringType, elems)
}

// DescribeCandidates renders values for a message, at most limit of them.
func DescribeCandidates(values []optionsource.Value, limit int) string {
	names := make([]string, 0, len(values))

	for _, v := range values {
		if v.Name != "" && v.Name != v.Value {
			names = append(names, fmt.Sprintf("%s (%s)", v.Name, v.Value))
		} else {
			names = append(names, v.Value)
		}
	}

	sort.Strings(names)

	if len(names) > limit {
		rest := len(names) - limit
		names = append(names[:limit], fmt.Sprintf("… and %d more", rest))
	}

	return strings.Join(names, "\n    ")
}

func attrPath(name string) path.Path { return path.Root(name) }
