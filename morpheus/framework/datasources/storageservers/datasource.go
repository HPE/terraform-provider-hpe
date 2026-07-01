// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package storageservers

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	providererrors "github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

const (
	summary = "read storage servers data source"
	// listMax bounds the number of records fetched from the API in one call.
	listMax            = 250
	cloudListMax       = 1000
	refTypeComputeZone = "ComputeZone"
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
	resp.TypeName = req.ProviderTypeName + "_" + "storage_servers"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = StorageServersDataSourceSchema(ctx)
	resp.Schema.Description = "Retrieves a list of Morpheus storage servers, optionally " +
		"filtered using one or more filter blocks."
	resp.Schema.MarkdownDescription = "Retrieves a list of Morpheus storage servers, optionally " +
		"filtered using one or more filter blocks."
}

// storageServersFilterModel decodes a single filter block.
type storageServersFilterModel struct {
	Name   types.String `tfsdk:"name"`
	Values types.Set    `tfsdk:"values"`
}

// compiledFilter is a filter block with its values pre-compiled as regular
// expressions.
type compiledFilter struct {
	field string
	res   []*regexp.Regexp
}

// Read refreshes the Terraform state with the latest data.
func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config StorageServersModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var blocks []storageServersFilterModel
	resp.Diagnostics.Append(config.Filter.ElementsAs(ctx, &blocks, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	filters := compileFilters(ctx, blocks, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	apiClient, err := d.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(summary, "could not create sdk client")

		return
	}

	cloudNames, err := listCloudNames(ctx, apiClient)
	if err != nil {
		resp.Diagnostics.AddError(summary, err.Error())

		return
	}

	rs, hresp, err := apiClient.StorageAPI.ListStorageServers(ctx).
		Max(listMax).
		Sort("id").
		Direction("asc").
		Execute()
	if rs == nil || err != nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(summary, fmt.Sprintf("LIST failed for storage servers: %s",
			providererrors.ErrMsg(err, hresp)))

		return
	}

	objs := make([]attr.Value, 0, len(rs.StorageServers))

	for i := range rs.StorageServers {
		ss := rs.StorageServers[i]
		if !storageServerMatchesFilters(&ss, filters, cloudNames) {
			continue
		}

		v, diags := storageServerInnerToValue(ctx, &ss, cloudNames)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		objs = append(objs, v)
	}

	setVal, diags := types.SetValue(StorageServersValue{}.Type(ctx), objs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	config.StorageServers = setVal

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

// compileFilters converts the filter blocks from configuration into compiled
// regular expressions. Invalid patterns are reported as diagnostics.
func compileFilters(
	ctx context.Context,
	blocks []storageServersFilterModel,
	diags *diag.Diagnostics,
) []compiledFilter {
	filters := make([]compiledFilter, 0, len(blocks))

	for _, b := range blocks {
		field := b.Name.ValueString()

		var values []string
		diags.Append(b.Values.ElementsAs(ctx, &values, false)...)
		if diags.HasError() {
			return nil
		}

		res := make([]*regexp.Regexp, 0, len(values))
		for _, v := range values {
			re, err := regexp.Compile(v)
			if err != nil {
				diags.AddError(summary,
					fmt.Sprintf("invalid regular expression %q for filter %q: %s", v, field, err))

				return nil
			}
			res = append(res, re)
		}

		filters = append(filters, compiledFilter{field: field, res: res})
	}

	return filters
}

// storageServerMatchesFilters reports whether ss satisfies every filter block.
// Within a block, the field must match ANY value (OR); across blocks all must
// match (AND).
func storageServerMatchesFilters(
	ss *sdk.ListStorageServers200ResponseAllOfStorageServersInner,
	filters []compiledFilter,
	cloudNames map[int64]string,
) bool {
	for _, f := range filters {
		val, ok := storageServerFieldValue(ss, f.field, cloudNames)
		if !ok {
			return false
		}

		matched := false
		for _, re := range f.res {
			if re.MatchString(val) {
				matched = true

				break
			}
		}
		if !matched {
			return false
		}
	}

	return true
}

// storageServerFieldValue returns the string representation of the named filter
// field for regex matching, and whether the field is present.
func storageServerFieldValue(
	ss *sdk.ListStorageServers200ResponseAllOfStorageServersInner,
	field string,
	cloudNames map[int64]string,
) (string, bool) {
	switch field {
	case "name":
		if ss.Name != nil {
			return *ss.Name, true
		}
	case "type_id":
		if ss.Type != nil && ss.Type.Id != nil {
			return strconv.FormatInt(*ss.Type.Id, 10), true
		}
	case "type_code":
		if ss.Type != nil && ss.Type.Code != nil {
			return *ss.Type.Code, true
		}
	case "cloud_id":
		if id, ok := cloudID(ss); ok {
			return strconv.FormatInt(id, 10), true
		}
	case "cloud_name":
		if id, ok := cloudID(ss); ok {
			if n, found := cloudNames[id]; found {
				return n, true
			}
		}
	case "visibility":
		if ss.Visibility != nil {
			return *ss.Visibility, true
		}
	case "enabled":
		if ss.Enabled != nil {
			return strconv.FormatBool(*ss.Enabled), true
		}
	case "status":
		if s := ss.Status.Get(); s != nil {
			return *s, true
		}
	}

	return "", false
}

// storageServerInnerToValue maps an API storage server into the generated
// custom object value used as the element of the storage_servers set.
func storageServerInnerToValue(
	ctx context.Context,
	ss *sdk.ListStorageServers200ResponseAllOfStorageServersInner,
	cloudNames map[int64]string,
) (StorageServersValue, diag.Diagnostics) {
	typeID := types.Int64Null()
	typeCode := types.StringNull()
	typeName := types.StringNull()
	if ss.Type != nil {
		typeID = convert.Int64ToType(ss.Type.Id)
		typeCode = convert.StrToType(ss.Type.Code)
		typeName = convert.StrToType(ss.Type.Name)
	}

	cloudIDVal := types.Int64Null()
	cloudNameVal := types.StringNull()
	if id, ok := cloudID(ss); ok {
		cloudIDVal = types.Int64Value(id)
		if n, found := cloudNames[id]; found {
			cloudNameVal = types.StringValue(n)
		}
	}

	attrs := map[string]attr.Value{
		"admin_port":       int32PtrToType(ss.AdminPort.Get()),
		"api_port":         int32PtrToType(ss.ApiPort.Get()),
		"category":         convert.StrToType(ss.Category.Get()),
		"cloud_id":         cloudIDVal,
		"cloud_name":       cloudNameVal,
		"date_created":     timeToType(ss.DateCreated),
		"description":      convert.StrToType(ss.Description.Get()),
		"disk_count":       int32PtrToType(ss.DiskCount.Get()),
		"enabled":          convert.BoolToType(ss.Enabled),
		"error_message":    convert.StrToType(ss.ErrorMessage.Get()),
		"external_id":      convert.StrToType(ss.ExternalId.Get()),
		"external_ip":      convert.StrToType(ss.ExternalIp.Get()),
		"id":               convert.Int64ToType(ss.Id),
		"internal_id":      convert.StrToType(ss.InternalId.Get()),
		"internal_ip":      convert.StrToType(ss.InternalIp.Get()),
		"last_updated":     timeToType(ss.LastUpdated),
		"max_storage":      convert.Int64ToType(ss.MaxStorage.Get()),
		"name":             convert.StrToType(ss.Name),
		"ref_id":           convert.Int64ToType(ss.RefId),
		"ref_type":         convert.StrToType(ss.RefType),
		"serial_number":    convert.StrToType(ss.SerialNumber.Get()),
		"server_model":     convert.StrToType(ss.ServerModel.Get()),
		"server_vendor":    convert.StrToType(ss.ServerVendor.Get()),
		"service_host":     convert.StrToType(ss.ServiceHost.Get()),
		"service_path":     convert.StrToType(ss.ServicePath.Get()),
		"service_url":      convert.StrToType(ss.ServiceUrl.Get()),
		"service_username": convert.StrToType(ss.ServiceUsername.Get()),
		"service_version":  convert.StrToType(ss.ServiceVersion.Get()),
		"status":           convert.StrToType(ss.Status.Get()),
		"status_date":      timeToType(ss.StatusDate),
		"status_message":   convert.StrToType(ss.StatusMessage.Get()),
		"type_code":        typeCode,
		"type_id":          typeID,
		"type_name":        typeName,
		"used_storage":     convert.Int64ToType(ss.UsedStorage.Get()),
		"visibility":       convert.StrToType(ss.Visibility),
	}

	return NewStorageServersValue(StorageServersValue{}.AttributeTypes(ctx), attrs)
}

// cloudID returns the cloud (zone) id of a storage server, valid only when the
// server is scoped to a ComputeZone via refType/refId.
func cloudID(
	ss *sdk.ListStorageServers200ResponseAllOfStorageServersInner,
) (int64, bool) {
	if ss.RefType == nil || *ss.RefType != refTypeComputeZone || ss.RefId == nil {
		return 0, false
	}

	return *ss.RefId, true
}

// listCloudNames returns a map of cloud id to cloud name, used to resolve the
// cloud_id / cloud_name attributes and filters from a storage server's
// refType/refId.
func listCloudNames(
	ctx context.Context,
	apiClient *sdk.APIClient,
) (map[int64]string, error) {
	rs, hresp, err := apiClient.CloudsAPI.ListClouds(ctx).Max(cloudListMax).Execute()
	if rs == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("LIST failed for clouds: %s", providererrors.ErrMsg(err, hresp))
	}

	names := make(map[int64]string, len(rs.Zones))
	for i := range rs.Zones {
		z := rs.Zones[i]
		if z.Id != nil && z.Name != nil {
			names[*z.Id] = *z.Name
		}
	}

	return names, nil
}

// int32PtrToType widens a nullable int32 (Grails Integer) to a Terraform Int64.
func int32PtrToType(i *int32) types.Int64 {
	if i == nil {
		return types.Int64Null()
	}

	return types.Int64Value(int64(*i))
}

// timeToType formats a nullable timestamp as an RFC 3339 string.
func timeToType(t *time.Time) types.String {
	if t == nil {
		return types.StringNull()
	}

	return types.StringValue(t.Format(time.RFC3339))
}
