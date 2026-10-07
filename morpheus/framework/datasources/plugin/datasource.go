// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package plugin

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sort"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

const (
	summary = "read plugin data source"

	ErrorNoValidSearchTerms = "no valid search terms - an id or code is required"
)

var _ datasource.DataSource = &DataSource{}

func NewDataSource() datasource.DataSource {
	return &DataSource{}
}

type DataSource struct {
	configure.DataSourceWithMorpheusConfigure
	datasource.DataSource
}

func (d *DataSource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_" + "plugin"
}

func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = PluginDataSourceSchema(ctx)
}

func populatePluginData(
	ctx context.Context,
	data *PluginModel,
	plugin *sdk.GetPlugin200ResponseAllOfPlugin,
) error {
	data.Id = convert.Int64ToType(plugin.Id)
	data.Name = convert.StrToType(plugin.Name)
	data.Code = convert.StrToType(plugin.Code)
	data.Description = convert.StrToType(plugin.Description)
	data.Version = convert.StrToType(plugin.Version)
	data.Enabled = convert.BoolToType(plugin.Enabled)
	data.Author = convert.StrToType(plugin.Author.Get())
	data.WebsiteUrl = convert.StrToType(plugin.WebsiteUrl.Get())
	data.SourceCodeLocationUrl = convert.StrToType(plugin.SourceCodeLocationUrl.Get())
	data.IssueTrackerUrl = convert.StrToType(plugin.IssueTrackerUrl.Get())
	data.Valid = convert.BoolToType(plugin.Valid)
	data.HasValidUpdate = convert.BoolToType(plugin.HasValidUpdate)
	data.IsPrepackaged = convert.BoolToType(plugin.IsPrepackaged)
	data.Status = convert.StrToType(plugin.Status)
	data.StatusMessage = convert.StrToType(plugin.StatusMessage.Get())
	data.DateCreated = convert.TimeToType(plugin.DateCreated)
	data.LastUpdated = convert.TimeToType(plugin.LastUpdated)

	providers, err := buildProvidersList(ctx, plugin.Providers)
	if err != nil {
		return err
	}
	data.Providers = providers

	config, err := convert.MapToDynamic(ctx, plugin.Config)
	if err != nil {
		return fmt.Errorf("plugin config mapping failed: %w", err)
	}
	data.Config = config

	return nil
}

// buildProvidersList maps the plugin's providers to a list with a stable order
// (type, then name): the API returns them in a different order on every read,
// and {type, name} is not unique, so a set cannot hold them.
func buildProvidersList(
	ctx context.Context,
	providers []sdk.GetPlugin200ResponseAllOfPluginProvidersInner,
) (types.List, error) {
	elemType := ProvidersValue{}.Type(ctx)

	sorted := make([]ProvidersValue, 0, len(providers))
	for _, p := range providers {
		sorted = append(sorted, ProvidersValue{
			Name:          convert.StrToType(p.Name),
			ProvidersType: convert.StrToType(p.Type),
			state:         attr.ValueStateKnown,
		})
	}
	sort.SliceStable(sorted, func(i, j int) bool {
		ti, tj := sorted[i].ProvidersType.ValueString(), sorted[j].ProvidersType.ValueString()
		if ti != tj {
			return ti < tj
		}

		return sorted[i].Name.ValueString() < sorted[j].Name.ValueString()
	})

	values := make([]attr.Value, 0, len(sorted))
	for _, v := range sorted {
		values = append(values, v)
	}

	list, diags := types.ListValueFrom(ctx, elemType, values)
	if diags.HasError() {
		return types.ListNull(elemType), fmt.Errorf("plugin providers mapping failed: %v", diags)
	}

	return list, nil
}

func getPluginByID(
	ctx context.Context,
	id int64,
	data *PluginModel,
	apiClient *sdk.APIClient,
) error {
	pluginResp, hresp, err := apiClient.PluginsAPI.GetPlugin(ctx, id).Execute()
	if hresp != nil && hresp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("plugin with id %d not found", id)
	}
	if err != nil || hresp == nil || hresp.StatusCode != http.StatusOK {
		return fmt.Errorf("plugin %d GET failed: %s", id, errfmt.ErrMsg(err, hresp))
	}

	if pluginResp == nil || pluginResp.Plugin == nil {
		return fmt.Errorf("plugin with id %d not found", id)
	}

	return populatePluginData(ctx, data, pluginResp.Plugin)
}

func getPluginByCode(
	ctx context.Context,
	code string,
	data *PluginModel,
	apiClient *sdk.APIClient,
) error {
	pluginsResp, hresp, err := apiClient.PluginsAPI.ListPlugins(ctx).Code(code).Execute()
	if err != nil || hresp == nil || hresp.StatusCode != http.StatusOK {
		return fmt.Errorf("plugin with code %s GET failed: %s", code, errfmt.ErrMsg(err, hresp))
	}

	var matching []sdk.ListPlugins200ResponseAllOfPluginsInner
	if pluginsResp != nil {
		for _, p := range pluginsResp.Plugins {
			if p.Code != nil && *p.Code == code {
				matching = append(matching, p)
			}
		}
	}

	if len(matching) == 0 {
		return fmt.Errorf("no plugin found with code %s", code)
	}

	if len(matching) > 1 {
		ids := make([]int64, 0, len(matching))
		for _, p := range matching {
			if p.Id != nil {
				ids = append(ids, *p.Id)
			}
		}

		return fmt.Errorf("multiple plugins found with code %s: ids %v", code, ids)
	}

	id := matching[0].Id
	if id == nil {
		return fmt.Errorf("no plugin found with code %s", code)
	}

	return getPluginByID(ctx, *id, data, apiClient)
}

func getPlugin(
	ctx context.Context,
	data *PluginModel,
	apiClient *sdk.APIClient,
) error {
	if !data.Id.IsNull() {
		return getPluginByID(ctx, data.Id.ValueInt64(), data, apiClient)
	}

	if !data.Code.IsNull() {
		return getPluginByCode(ctx, data.Code.ValueString(), data, apiClient)
	}

	return errors.New(ErrorNoValidSearchTerms)
}

func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var data PluginModel

	diags := req.Config.Get(ctx, &data)
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

	if err := getPlugin(ctx, &data, apiClient); err != nil {
		resp.Diagnostics.AddError(
			summary,
			err.Error(),
		)

		return
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
