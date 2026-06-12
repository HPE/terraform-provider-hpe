// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package serviceplan

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	internalErrors "github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

const (
	summary                 = "read service plan data source"
	ErrorNoServicePlanFound = `no service plan found`
	ErrorNoValidSearchTerms = "no valid search terms - an id or (name and provision_type_code) " +
		"is required"
	ErrorRunningPreApply      = `Error running pre-apply plan: exit status 1`
	ErrorMultipleServicePlans = `multiple service plans were returned`
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
	resp.TypeName = req.ProviderTypeName + "_morpheus_service_plan"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = ServicePlanDataSourceSchema(ctx)
}

func getServicePlanByID(
	ctx context.Context,
	id int64,
	apiClient *sdk.APIClient,
) (*sdk.GetServicePlans200ResponseServicePlan, error) {
	sp, hresp, err := apiClient.ServicePlansAPI.GetServicePlans(ctx, id).Execute()
	if sp == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"GET failed for service plan %d: %s", id, internalErrors.ErrMsg(err, hresp))
	}

	if sp.ServicePlan == nil {
		return nil, fmt.Errorf("service plan %d is nil", id)
	}

	return sp.ServicePlan, nil
}

// servicePlanInZone reports whether the service plan is available in the cloud
// (zone) with the given id, based on the zones returned by ListServicePlans when
// IncludeZones(true) is set.
func servicePlanInZone(
	sp sdk.ListServicePlans200ResponseAllOfServicePlansInner,
	cloudID int64,
) bool {
	for _, z := range sp.Zones {
		if z.Id != nil && *z.Id == cloudID {
			return true
		}
	}

	return false
}

func getServicePlanByName(
	ctx context.Context,
	name string,
	provisionTypeCode string,
	cloudID *int64,
	apiClient *sdk.APIClient,
) (*sdk.GetServicePlans200ResponseServicePlan, error) {
	pTypes, hresp, err := apiClient.ProvisioningAPI.ListProvisionTypes(ctx).Code(
		provisionTypeCode).Execute()
	if pTypes == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET failed for service plan , provision type code %s: %s",
			provisionTypeCode, internalErrors.ErrMsg(err, hresp))
	}

	var matchingProvisionTypes []sdk.
		ListProvisionTypes200ResponseAllOfProvisionTypesInner
	for _, pt := range pTypes.ProvisionTypes {
		if pt.Code != nil && *pt.Code == provisionTypeCode {
			matchingProvisionTypes = append(matchingProvisionTypes, pt)
		}
	}

	if len(matchingProvisionTypes) == 0 {
		return nil, fmt.Errorf("provision type with code %s not found", provisionTypeCode)
	}

	if len(matchingProvisionTypes) > 1 {
		return nil, fmt.Errorf("multiple provision types with code %s found", provisionTypeCode)
	}

	if matchingProvisionTypes[0].Id == nil {
		return nil, fmt.Errorf("id not found for provision type with code %s", provisionTypeCode)
	}

	// IncludeZones(true) returns, for each plan, the clouds (zones) it is available
	// in, so we can disambiguate plans that share a name across clouds/regions
	// (e.g. Azure) using the optional cloud_id filter.
	ps, hresp, err := apiClient.ServicePlansAPI.ListServicePlans(ctx).Name(
		name).ProvisionTypeId(*matchingProvisionTypes[0].Id).IncludeZones(true).Execute()
	if ps == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"GET failed for service_plan %s: %s", name, internalErrors.ErrMsg(err, hresp))
	}

	var matchingServicePlans []sdk.ListServicePlans200ResponseAllOfServicePlansInner

	for _, sp := range ps.ServicePlans {
		if sp.Name != nil && sp.ProvisionType != nil && sp.ProvisionType.Code != nil {
			// now check name and ProvisionType match getplanByName() params
			if *sp.Name == name && *sp.ProvisionType.Code == provisionTypeCode {
				// when cloud_id is set, only keep plans available in that cloud
				if cloudID != nil && !servicePlanInZone(sp, *cloudID) {
					continue
				}
				matchingServicePlans = append(matchingServicePlans, sp)
			}
		}
	}
	if len(matchingServicePlans) == 1 {
		if matchingServicePlans[0].Id != nil {
			// same return types as GetPlanByID
			return getServicePlanByID(ctx, *matchingServicePlans[0].Id, apiClient)
		}

		return nil, fmt.Errorf("service plan %s, id not found", name)
	} else if len(matchingServicePlans) > 1 {
		return nil, errors.New(ErrorMultipleServicePlans)
	}

	return nil, errors.New(ErrorNoServicePlanFound)
}

func getServicePlan(
	ctx context.Context,
	data ServicePlanModel,
	apiClient *sdk.APIClient,
) (*sdk.GetServicePlans200ResponseServicePlan, error) {
	if !data.Id.IsNull() {
		return getServicePlanByID(ctx, data.Id.ValueInt64(), apiClient)
	} else if !data.Name.IsNull() && !data.ProvisionTypeCode.IsNull() {
		var cloudID *int64
		if !data.CloudId.IsNull() {
			cloudID = data.CloudId.ValueInt64Pointer()
		}

		return getServicePlanByName(
			ctx, data.Name.ValueString(), data.ProvisionTypeCode.ValueString(), cloudID, apiClient)
	}

	return nil, errors.New(ErrorNoValidSearchTerms)
}

func servicePlanAsState(
	ctx context.Context,
	servicePlan *sdk.GetServicePlans200ResponseServicePlan,
) (ServicePlanModel, diag.Diagnostics) {
	var state ServicePlanModel
	var diags diag.Diagnostics

	priceSetIDValues := []attr.Value{}
	for _, v := range servicePlan.PriceSets {
		priceSetIDValues = append(priceSetIDValues, convert.Int64ToType(v.Id))
	}

	priceSetIDSet, diags := types.SetValue(types.Int64Type, priceSetIDValues)
	if diags.HasError() {
		return state, diags
	}

	if servicePlan.Config.Ranges != nil {
		configRangesValue, diags := NewConfigRangesValue(
			ConfigRangesValue{}.AttributeTypes(ctx),
			map[string]attr.Value{
				"min_storage":          convert.Int64ToType(servicePlan.Config.Ranges.MinStorage.Get()),
				"max_storage":          convert.Int64ToType(servicePlan.Config.Ranges.MaxStorage.Get()),
				"min_memory":           convert.Int64ToType(servicePlan.Config.Ranges.MinMemory.Get()),
				"max_memory":           convert.Int64ToType(servicePlan.Config.Ranges.MaxMemory.Get()),
				"min_cores":            convert.Int64ToType(servicePlan.Config.Ranges.MinCores.Get()),
				"max_cores":            convert.Int64ToType(servicePlan.Config.Ranges.MaxCores.Get()),
				"min_sockets":          convert.Int64ToType(servicePlan.Config.Ranges.MinSockets.Get()),
				"max_sockets":          convert.Int64ToType(servicePlan.Config.Ranges.MaxSockets.Get()),
				"min_cores_per_socket": convert.Int64ToType(servicePlan.Config.Ranges.MinCoresPerSocket.Get()),
				"max_cores_per_socket": convert.Int64ToType(servicePlan.Config.Ranges.MaxCoresPerSocket.Get()),
				"min_per_disk_size":    convert.Int64ToType(servicePlan.Config.Ranges.MinPerDiskSize.Get()),
				"max_per_disk_size":    convert.Int64ToType(servicePlan.Config.Ranges.MaxPerDiskSize.Get()),
			},
		)
		if diags.HasError() {
			return state, diags
		}

		state.ConfigRanges = configRangesValue
	}

	state.AddVolumes = convert.BoolToType(servicePlan.AddVolumes.Get())
	state.Code = convert.StrToType(servicePlan.Code)
	state.CoresPerSocket = convert.Int64ToType(servicePlan.CoresPerSocket.Get())
	state.CustomCores = convert.BoolToType(servicePlan.CustomCores)
	state.CustomCpu = convert.BoolToType(servicePlan.CustomCpu)
	state.CustomMaxMemory = convert.BoolToType(servicePlan.CustomMaxMemory.Get())
	state.CustomMaxStorage = convert.BoolToType(servicePlan.CustomMaxStorage.Get())
	state.Description = convert.StrToType(servicePlan.Description)
	state.Id = convert.Int64ToType(servicePlan.Id)
	state.MaxCores = convert.Int64ToType(servicePlan.MaxCores.Get())
	state.MaxCpu = convert.Int64ToType(servicePlan.MaxCpu.Get())
	state.MaxDisks = convert.Int64ToType(servicePlan.MaxDisks.Get())
	state.MaxMemory = convert.Int64ToType(servicePlan.MaxMemory)
	state.MaxStorage = convert.Int64ToType(servicePlan.MaxStorage)
	state.MemorySizeType = convert.StrToType(servicePlan.Config.MemorySizeType.Get())
	state.Name = convert.StrToType(servicePlan.Name)
	state.PriceSetIds = priceSetIDSet
	state.ProvisionTypeCode = convert.StrToType(servicePlan.ProvisionType.Code)
	state.SortOrder = convert.Int64ToType(servicePlan.SortOrder)
	state.StorageSizeType = convert.StrToType(servicePlan.Config.StorageSizeType.Get())

	return state, diags
}

// Read refreshes the Terraform state with the latest data.
func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var data ServicePlanModel

	// Read config
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

	plan, err := getServicePlan(ctx, data, apiClient)
	if err != nil {
		resp.Diagnostics.AddError(
			summary,
			err.Error(),
		)

		return
	}

	apiState, diags := servicePlanAsState(ctx, plan)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)

		return
	}

	// cloud_id is a filter-only input; the API does not return it, so echo the
	// configured value back into state to keep config and state consistent.
	apiState.CloudId = data.CloudId

	diags = resp.State.Set(ctx, &apiState)
	resp.Diagnostics.Append(diags...)
}
