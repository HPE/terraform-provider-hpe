// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

//go:build experimental

package serviceplan

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/convert"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/errors"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                   = &Resource{}
	_ resource.ResourceWithImportState    = &Resource{}
	_ resource.ResourceWithValidateConfig = &Resource{}
)

func NewResource() resource.Resource {
	return &Resource{}
}

// Resource defines the resource implementation.
type Resource struct {
	configure.ResourceWithMorpheusConfigure
	resource.Resource
}

func (r *Resource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_morpheus_service_plan"
}

func (r *Resource) Schema(
	ctx context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = ServicePlanResourceSchema(ctx)
}

// populate service plan resource model with current API values
func getServicePlanAsState(
	ctx context.Context,
	id int64,
	client *sdk.APIClient,
) (ServicePlanModel, diag.Diagnostics) {
	var state ServicePlanModel
	var diags diag.Diagnostics

	sp, hresp, err := client.ServicePlansAPI.GetServicePlans(ctx, id).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK || sp == nil {
		diags.AddError(
			"populate service plan resource",
			fmt.Sprintf("service plan %d GET failed", id)+errors.ErrMsg(err, hresp),
		)
		return state, diags
	}

	pricesetIDValues := []attr.Value{}
	for _, v := range sp.ServicePlan.PriceSets {
		pricesetIDValues = append(pricesetIDValues, convert.Int64ToType(v.Id))
	}

	pricesetIDSet, diags := types.SetValue(types.Int64Type, pricesetIDValues)
	if diags.HasError() {
		return state, diags
	}

	// config is defaulted to empty map by api
	apiConfig := sp.ServicePlan.Config
	// config top level fields moved to service plan top level
	state.StorageSizeType = convert.StrToType(apiConfig.StorageSizeType.Get())
	state.MemorySizeType = convert.StrToType(apiConfig.MemorySizeType.Get())

	if apiConfig.Ranges != nil {
		configRanges := map[string]attr.Value{}

		configRanges["min_storage"] = convert.Int64ToType(apiConfig.Ranges.MinStorage.Get())
		configRanges["max_storage"] = convert.Int64ToType(apiConfig.Ranges.MaxStorage.Get())
		configRanges["min_memory"] = convert.Int64ToType(apiConfig.Ranges.MinMemory.Get())
		configRanges["max_memory"] = convert.Int64ToType(apiConfig.Ranges.MaxMemory.Get())
		configRanges["min_cores"] = convert.Int64ToType(apiConfig.Ranges.MinCores.Get())
		configRanges["max_cores"] = convert.Int64ToType(apiConfig.Ranges.MaxCores.Get())
		configRanges["min_sockets"] = convert.Int64ToType(apiConfig.Ranges.MinSockets.Get())
		configRanges["max_sockets"] = convert.Int64ToType(apiConfig.Ranges.MaxSockets.Get())
		configRanges["min_cores_per_socket"] = convert.Int64ToType(apiConfig.Ranges.MinCoresPerSocket.Get())
		configRanges["max_cores_per_socket"] = convert.Int64ToType(apiConfig.Ranges.MaxCoresPerSocket.Get())
		configRanges["min_per_disk_size"] = convert.Int64ToType(apiConfig.Ranges.MinPerDiskSize.Get())
		configRanges["max_per_disk_size"] = convert.Int64ToType(apiConfig.Ranges.MaxPerDiskSize.Get())

		configRangesValue, diags := NewConfigRangesValue(ConfigRangesValue{}.AttributeTypes(ctx), configRanges)
		if diags.HasError() {
			return state, diags
		}
		state.ConfigRanges = configRangesValue
	}

	state.PricesetIds = pricesetIDSet
	state.Id = convert.Int64ToType(sp.ServicePlan.Id)
	state.Name = convert.StrToType(sp.ServicePlan.Name)
	state.Code = convert.StrToType(sp.ServicePlan.Code)
	state.MaxMemory = convert.Int64ToType(sp.ServicePlan.MaxMemory)
	state.MaxStorage = convert.Int64ToType(sp.ServicePlan.MaxStorage)
	state.MaxCores = convert.Int64ToType(sp.ServicePlan.MaxCores.Get())
	state.MaxCpu = convert.Int64ToType(sp.ServicePlan.MaxCpu.Get())
	state.MaxDisks = convert.Int64ToType(sp.ServicePlan.MaxDisks.Get())
	state.CoresPerSocket = convert.Int64ToType(sp.ServicePlan.CoresPerSocket.Get())
	state.CustomCores = convert.BoolToType(sp.ServicePlan.CustomCores)
	state.CustomCpu = convert.BoolToType(sp.ServicePlan.CustomCpu)
	state.CustomMaxMemory = convert.BoolToType(sp.ServicePlan.CustomMaxMemory.Get())
	state.CustomMaxStorage = convert.BoolToType(sp.ServicePlan.CustomMaxStorage.Get())
	state.AddVolumes = convert.BoolToType(sp.ServicePlan.AddVolumes.Get())
	state.Description = convert.StrToType(sp.ServicePlan.Description)
	state.SortOrder = convert.Int64ToType(sp.ServicePlan.SortOrder)

	if sp.ServicePlan.ProvisionType != nil {
		state.ProvisionTypeCode = convert.StrToType(sp.ServicePlan.ProvisionType.Code)
	}

	return state, diags
}

// helper function to set provisionType in addServicePlan struct
func setProvisionTypeInCreate(
	ctx context.Context,
	client *sdk.APIClient,
	plan *ServicePlanModel,
	addServicePlan *sdk.AddServicePlansRequestServicePlan,
) error {
	provisionTypeCode := plan.ProvisionTypeCode.ValueString()

	pTypes, hresp, err := client.ProvisioningAPI.ListProvisionTypes(ctx).Code(
		provisionTypeCode).Execute()
	if pTypes == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return fmt.Errorf("GET failed for provision type code %s: %s", provisionTypeCode, errors.ErrMsg(err, hresp))
	}

	var matchingProvisionTypes []sdk.
		GetInstanceTypeProvisioning200ResponseAllOfInstanceTypeInstanceTypeLayoutsInnerProvisionType
	for _, pt := range pTypes.GetProvisionTypes() {
		if ptCode, ok := pt.GetCodeOk(); ok && *ptCode == provisionTypeCode {
			matchingProvisionTypes = append(matchingProvisionTypes, pt)
		}
	}

	if len(matchingProvisionTypes) == 0 {
		return fmt.Errorf("provision type with code %s not found", provisionTypeCode)
	}

	if len(matchingProvisionTypes) > 1 {
		return fmt.Errorf("multiple provision types with code %s found", provisionTypeCode)
	}

	pTypeID, ok := matchingProvisionTypes[0].GetIdOk()
	if !ok {
		return fmt.Errorf("id not found for provision type with code %s", provisionTypeCode)
	}

	provisionType := sdk.AddClusterLayoutsRequestLayoutProvisionType{}
	provisionType.Id = *pTypeID
	addServicePlan.ProvisionType = provisionType

	return nil
}

// helper function to nest schema values into config struct
func setConfigInCreate(
	ctx context.Context,
	plan *ServicePlanModel,
	addServicePlan *sdk.AddServicePlansRequestServicePlan,
) {
	config := sdk.NewAddServicePlansRequestServicePlanConfig()

	// top level fields first
	if !plan.StorageSizeType.IsNull() {
		storageSize := plan.StorageSizeType.ValueString()

		config.StorageSizeType = &storageSize
	}
	if !plan.MemorySizeType.IsNull() {
		memorySize := plan.MemorySizeType.ValueString()
		config.MemorySizeType = &memorySize
	}

	// ConfigRanges
	if !plan.ConfigRanges.IsNull() {
		ranges := sdk.NewAddServicePlansRequestServicePlanConfigRanges()

		if !plan.ConfigRanges.MinMemory.IsNull() {
			minMemory := plan.ConfigRanges.MinMemory.ValueInt64Pointer()
			ranges.MinMemory = minMemory
		}
		if !plan.ConfigRanges.MaxMemory.IsNull() {
			maxMemory := plan.ConfigRanges.MaxMemory.ValueInt64Pointer()
			ranges.MaxMemory = maxMemory
		}
		if !plan.ConfigRanges.MinStorage.IsNull() {
			minStorage := plan.ConfigRanges.MinStorage.ValueInt64Pointer()
			ranges.MinStorage = minStorage
		}
		if !plan.ConfigRanges.MaxStorage.IsNull() {
			maxStorage := plan.ConfigRanges.MaxStorage.ValueInt64Pointer()
			ranges.MaxStorage = maxStorage
		}
		if !plan.ConfigRanges.MinCores.IsNull() {
			minCores := plan.ConfigRanges.MinCores.ValueInt64Pointer()
			ranges.MinCores = minCores
		}
		if !plan.ConfigRanges.MaxCores.IsNull() {
			maxCores := plan.ConfigRanges.MaxCores.ValueInt64Pointer()
			ranges.MaxCores = maxCores
		}
		if !plan.ConfigRanges.MinCoresPerSocket.IsNull() {
			minCoresPerSocket := plan.ConfigRanges.MinCoresPerSocket.ValueInt64Pointer()
			ranges.MinCoresPerSocket = minCoresPerSocket
		}
		if !plan.ConfigRanges.MaxCoresPerSocket.IsNull() {
			maxCoresPerSocket := plan.ConfigRanges.MaxCoresPerSocket.ValueInt64Pointer()
			ranges.MaxCoresPerSocket = maxCoresPerSocket
		}
		if !plan.ConfigRanges.MinPerDiskSize.IsNull() {
			minPerDiskSize := plan.ConfigRanges.MinPerDiskSize.ValueInt64Pointer()
			ranges.MinPerDiskSize = minPerDiskSize
		}
		if !plan.ConfigRanges.MaxPerDiskSize.IsNull() {
			maxPerDiskSize := plan.ConfigRanges.MaxPerDiskSize.ValueInt64Pointer()
			ranges.MaxPerDiskSize = maxPerDiskSize
		}
		if !plan.ConfigRanges.MinSockets.IsNull() {
			minSockets := plan.ConfigRanges.MinSockets.ValueInt64Pointer()
			ranges.MinSockets = minSockets
		}
		if !plan.ConfigRanges.MaxSockets.IsNull() {
			maxSockets := plan.ConfigRanges.MaxSockets.ValueInt64Pointer()
			ranges.MaxSockets = maxSockets
		}
		config.Ranges = ranges
	}
	addServicePlan.Config = config
}

func (r *Resource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan ServicePlanModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := plan.Name.ValueString()
	addServicePlan := sdk.NewAddServicePlansRequestServicePlanWithDefaults()

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"create service plan resource",
			"service plan "+name+": failed to create client: "+err.Error(),
		)

		return
	}

	// required
	addServicePlan.SetName(name)
	addServicePlan.SetCode(plan.Code.ValueString())
	addServicePlan.SetMaxMemory(plan.MaxMemory.ValueInt64())
	addServicePlan.SetMaxStorage(plan.MaxStorage.ValueInt64())

	err = setProvisionTypeInCreate(ctx, client, &plan, addServicePlan)
	if err != nil {
		// setProvisionTypeInCreate checks for cert err
		resp.Diagnostics.AddError(
			"create service plan resource",
			"set provision type POST failed : "+err.Error(),
		)

		return
	}

	// optional
	if !plan.PricesetIds.IsNull() && !plan.PricesetIds.IsUnknown() {
		var pricesetIds []int64
		diags := plan.PricesetIds.ElementsAs(ctx, &pricesetIds, false)
		if resp.Diagnostics.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}

		var pricesets []sdk.AddServicePlansRequestServicePlanPriceSetsInner
		for _, v := range pricesetIds {
			priceset := sdk.AddServicePlansRequestServicePlanPriceSetsInner{
				Id: &v,
			}
			pricesets = append(pricesets, priceset)
		}

		addServicePlan.PriceSets = pricesets
	}

	if !plan.Description.IsNull() {
		description := plan.Description.ValueStringPointer()
		addServicePlan.Description = description
	}

	if !plan.MaxCores.IsNull() {
		maxCores := plan.MaxCores.ValueInt64Pointer()
		addServicePlan.MaxCores = maxCores
	}

	if !plan.MaxCpu.IsNull() {
		maxCpu := plan.MaxCpu.ValueInt64Pointer()
		addServicePlan.MaxCpu = maxCpu
	}

	if !plan.MaxDisks.IsNull() {
		maxDisks := plan.MaxDisks.ValueInt64Pointer()
		addServicePlan.MaxDisks = maxDisks
	}

	if !plan.CoresPerSocket.IsNull() && !plan.CoresPerSocket.IsUnknown() {
		coresPerSocket := plan.CoresPerSocket.ValueInt64()
		addServicePlan.CoresPerSocket = &coresPerSocket
	}

	if !plan.CustomCores.IsNull() && !plan.CustomCores.IsUnknown() {
		customCores := plan.CustomCores.ValueBool()
		addServicePlan.CustomCores = &customCores
	}

	if !plan.CustomCpu.IsNull() && !plan.CustomCpu.IsUnknown() {
		customCpu := plan.CustomCpu.ValueBool()
		addServicePlan.CustomCpu = &customCpu
	}

	if !plan.CustomMaxMemory.IsNull() && !plan.CustomMaxMemory.IsUnknown() {
		customMaxMemory := plan.CustomMaxMemory.ValueBool()
		addServicePlan.CustomMaxMemory = &customMaxMemory
	}

	if !plan.CustomMaxStorage.IsNull() && !plan.CustomMaxStorage.IsUnknown() {
		customMaxStorage := plan.CustomMaxStorage.ValueBool()
		addServicePlan.CustomMaxStorage = &customMaxStorage
	}

	if !plan.AddVolumes.IsNull() && !plan.AddVolumes.IsUnknown() {
		addVolumes := plan.AddVolumes.ValueBool()
		addServicePlan.AddVolumes = &addVolumes
	}

	if !plan.SortOrder.IsNull() {
		sortOrder := plan.SortOrder.ValueInt64()
		addServicePlan.SortOrder = &sortOrder
	}

	if !plan.MemorySizeType.IsNull() || !plan.StorageSizeType.IsNull() ||
		!plan.ConfigRanges.IsNull() {
		setConfigInCreate(ctx, &plan, addServicePlan)
	}

	addServicePlanRequest := sdk.NewAddServicePlansRequest(*addServicePlan)

	servicePlan, hresp, err := client.ServicePlansAPI.AddServicePlans(ctx).AddServicePlansRequest(*addServicePlanRequest).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"create service plan resource",
			"service plan "+name+" POST failed: "+errors.ErrMsg(err, hresp),
		)

		return
	}

	if servicePlan.Id == nil {
		resp.Diagnostics.AddError(
			"create service plan resource",
			"service plan"+name+" id is nil",
		)
		return
	}
	id := *servicePlan.Id
	plan.Id = types.Int64Value(id)

	// write id as soon as possible
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, diags := getServicePlanAsState(ctx, id, client)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// update not implemented for now
func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
}

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var plan ServicePlanModel

	diags := req.State.Get(ctx, &plan)
	if diags.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"read service plan resource",
			"new client call failed with "+err.Error(),
		)

		return
	}

	id := plan.Id.ValueInt64()
	state, diags := getServicePlanAsState(ctx, id, client)
	if diags.HasError() {
		resp.Diagnostics.AddError(
			"read service plan resource",
			fmt.Sprintf("service plan %d: failed to read from api", id),
		)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var data ServicePlanModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	id := data.Id.ValueInt64()

	client, _ := r.NewClient(ctx)

	_, hresp, err := client.ServicePlansAPI.RemoveServicePlans(ctx, id).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"delete service plan resource",
			fmt.Sprintf("service plan %d: DELETE failed ", id)+errors.ErrMsg(err, hresp),
		)

		return
	}
}

func (r *Resource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	id, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"import service plan resource",
			"provided import ID '"+req.ID+"' is invalid (non-number)",
		)

		return
	}

	diags := resp.State.SetAttribute(
		ctx, path.Root("id"), id,
	)
	resp.Diagnostics.Append(diags...)
}

// This method is called by Terraform's ValidateResourceConfig RPC.
// we use this to validate min and max ranges set by boolean values in the schema.
// If the boolean value is null or false, then setting the corresponding min/max range is redundant.
func (r *Resource) ValidateConfig(
	ctx context.Context,
	req resource.ValidateConfigRequest,
	resp *resource.ValidateConfigResponse,
) {
	var config ServicePlanModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if !config.ConfigRanges.IsNull() {

		// ValueBool returns false if value was set to false, or if state is null or unknown
		if !config.CustomMaxMemory.ValueBool() {

			if !config.ConfigRanges.MinMemory.IsNull() {
				resp.Diagnostics.AddAttributeError(
					path.Root("config_ranges.min_memory"),
					"Conflicting attributes in configuration",
					`min_memory set when custom_memory has not been `+
						`set custom_memory to true to add a min_memory value.`,
				)
				return
			}
			if !config.ConfigRanges.MaxMemory.IsNull() {
				resp.Diagnostics.AddAttributeError(
					path.Root("config_ranges.max_memory"),
					"Conflicting attributes in configuration",
					`max_memory set when custom_memory has not been `+
						`set custom_memory to true to add a max_memory value.`,
				)
				return
			}

		}

		if !config.CustomMaxStorage.ValueBool() {

			if !config.ConfigRanges.MinStorage.IsNull() {
				resp.Diagnostics.AddAttributeError(
					path.Root("config_ranges.min_storage"),
					"Conflicting attributes in configuration",
					`min_storage set when custom_storage has not been `+
						`set custom_storage to true to add a min_storage value.`,
				)
				return
			}
			if !config.ConfigRanges.MaxStorage.IsNull() {
				resp.Diagnostics.AddAttributeError(
					path.Root("config_ranges.max_storage"),
					"Conflicting attributes in configuration",
					`max_storage set when custom_storage has not been `+
						`set custom_storage to true to add a max_storage value.`,
				)
				return
			}

		}

		// customCores is used with minCores and maxCores
		// if customCores is false or null, then min max should not be specified
		if !config.CustomCores.ValueBool() {

			if !config.ConfigRanges.MinCores.IsNull() {
				resp.Diagnostics.AddAttributeError(
					path.Root("config_ranges.min_cores"),
					"Conflicting attributes in configuration",
					`min_cores set when custom_cores  has not been `+
						`set to true to add a min_cores value.`,
				)
				return
			}
			if !config.ConfigRanges.MaxCores.IsNull() {
				resp.Diagnostics.AddAttributeError(
					path.Root("config_ranges.max_cores"),
					"Conflicting attributes in configuration",
					`max_cores set when custom_cores has not been `+
						`set custom_cores to true to add a max_cores value.`,
				)
				return
			}

		}

		if !config.AddVolumes.ValueBool() {

			if !config.ConfigRanges.MinPerDiskSize.IsNull() {
				resp.Diagnostics.AddAttributeError(
					path.Root("config.min_per_disk_size"),
					"Conflicting attributes in configuration",
					`min_per_disk_size set when add_volumes has not been `+
						`set add_volumes to true to add a min_per_disk_size value.`,
				)
				return
			}
			if !config.ConfigRanges.MaxPerDiskSize.IsNull() {
				resp.Diagnostics.AddAttributeError(
					path.Root("config.max_per_disk_size"),
					"Conflicting attributes in configuration",
					`max_per_disk_size set when add_volumes has not been `+
						`set add_volumes to true to add a max_per_disk_size value.`,
				)
				return
			}

		}

	}
}
