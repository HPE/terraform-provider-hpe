// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package clusterlayout

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	providererrors "github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

const (
	summary                     = "read cluster layout data source"
	ErrorNoValidSearchTerms     = `no valid search terms - an id or name is required`
	ErrorRunningPreApply        = `Error running pre-apply plan: exit status 1`
	ErrorNoClusterLayoutFound   = `no cluster layout found`
	ErrorMultipleClusterLayouts = `multiple cluster layouts were returned`
)

// Ensure the implementation satisfies the expected interfaces.
var _ datasource.DataSource = &DataSource{}

// NewDataSource is a helper function to simplify the provider implementation.
func NewDataSource() datasource.DataSource {
	return &DataSource{}
}

// DataSource is the data source implementation.
type DataSource struct {
	configure.DataSourceWithMorpheusConfigure
	datasource.DataSource
}

// Metadata returns the data source type name.
func (d *DataSource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_" + "cluster_layout"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = ClusterLayoutDataSourceSchema(ctx)
}

func clusterLayoutAsState(
	_ context.Context,
	l *sdk.GetClusterLayout200ResponseLayout,
) (ClusterLayoutModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	state := ClusterLayoutModel{
		Id:                      convert.Int64ToType(l.Id),
		Name:                    convert.StrToType(l.Name),
		Code:                    convert.StrToType(l.Code),
		Description:             convert.StrToType(l.Description),
		InternalId:              convert.StrToType(l.InternalId),
		ClusterVersion:          convert.StrToType(l.ClusterVersion),
		ComputeVersion:          convert.StrToType(l.ComputeVersion),
		Creatable:               convert.BoolToType(l.Creatable),
		Enabled:                 convert.BoolToType(l.Enabled),
		HasAutoScale:            convert.BoolToType(l.HasAutoScale),
		HasConfig:               convert.BoolToType(l.HasConfig),
		HasSettings:             convert.BoolToType(l.HasSettings),
		InstallContainerRuntime: convert.BoolToType(l.InstallContainerRuntime),
		MemoryRequirement:       convert.Int64ToType(l.MemoryRequirement),
		ServerCount:             convert.Int64ToType(l.ServerCount),
		SortOrder:               convert.Int64ToType(l.SortOrder),
	}

	// DateCreated is *time.Time
	if l.DateCreated != nil {
		state.DateCreated = types.StringValue(l.DateCreated.String())
	} else {
		state.DateCreated = types.StringNull()
	}

	// LastUpdated is *time.Time
	if l.LastUpdated != nil {
		state.LastUpdated = types.StringValue(l.LastUpdated.String())
	} else {
		state.LastUpdated = types.StringNull()
	}

	// Labels is []string
	if l.Labels != nil {
		state.Labels = convert.StrSliceToSet(l.Labels)
	} else {
		state.Labels = convert.StrSliceToSet([]string{})
	}

	return state, diags
}

func getClusterLayoutByID(
	ctx context.Context,
	id int64,
	apiClient *sdk.APIClient,
) (*sdk.GetClusterLayout200ResponseLayout, error) {
	r, hresp, err := apiClient.ClusterLayoutsAPI.GetClusterLayout(ctx, id).Execute()
	if r == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET failed for cluster layout %d: %s", id, providererrors.ErrMsg(err, hresp))
	}

	if r.Layout == nil {
		return nil, fmt.Errorf("GET failed for cluster layout %d: response missing layout", id)
	}

	layout := *r.Layout

	return &layout, nil
}

func getClusterLayoutByName(
	ctx context.Context,
	name string,
	apiClient *sdk.APIClient,
) (*sdk.GetClusterLayout200ResponseLayout, error) {
	rs, hresp, err := apiClient.ClusterLayoutsAPI.ListClusterLayouts(ctx).Phrase(name).Execute()
	if rs == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET failed for cluster layout %s: %s", name, providererrors.ErrMsg(err, hresp))
	}

	// Use JSON round-trip for safe extraction since SDK list types may vary.
	raw, marshalErr := json.Marshal(rs.Layouts)
	if marshalErr != nil {
		return nil, fmt.Errorf("error marshaling cluster layouts: %w", marshalErr)
	}

	var layouts []struct {
		Id   *int64  `json:"id"`
		Name *string `json:"name"`
	}

	if unmarshalErr := json.Unmarshal(raw, &layouts); unmarshalErr != nil {
		return nil, fmt.Errorf("error decoding cluster layouts: %w", unmarshalErr)
	}

	// Client-side exact-match filter (Phrase is fuzzy).
	var matchedID int64

	var matchCount int

	for _, l := range layouts {
		if l.Name != nil && *l.Name == name {
			if l.Id != nil {
				matchedID = *l.Id
			}

			matchCount++
		}
	}

	if matchCount == 0 {
		return nil, errors.New(ErrorNoClusterLayoutFound)
	} else if matchCount > 1 {
		return nil, errors.New(ErrorMultipleClusterLayouts)
	}

	return getClusterLayoutByID(ctx, matchedID, apiClient)
}

func getClusterLayout(
	ctx context.Context,
	config *ClusterLayoutModel,
	apiClient *sdk.APIClient,
) (*sdk.GetClusterLayout200ResponseLayout, error) {
	if !config.Id.IsNull() {
		return getClusterLayoutByID(ctx, config.Id.ValueInt64(), apiClient)
	} else if !config.Name.IsNull() {
		return getClusterLayoutByName(ctx, config.Name.ValueString(), apiClient)
	}

	return nil, errors.New(ErrorNoValidSearchTerms)
}

// Read refreshes the Terraform state with the latest data.
func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config ClusterLayoutModel

	// Read config
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiClient, err := d.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			summary,
			"could not create sdk client",
		)

		return
	}

	layout, err := getClusterLayout(ctx, &config, apiClient)
	if err != nil {
		resp.Diagnostics.AddError(
			summary,
			err.Error(),
		)

		return
	}

	state, stateDiags := clusterLayoutAsState(ctx, layout)
	resp.Diagnostics.Append(stateDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
