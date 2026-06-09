// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package compare

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/iancoleman/strcase"
	"github.com/stretchr/testify/assert"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/sdkfuncs"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

const (
	testKey1 = "key1"
	testKey2 = "key2"

	testValue1 = "value1"
	testValue2 = 42.0

	testCertificateProvider        = "certificate"
	testApplianceUrl               = "https://morpheus.example.com"
	testExternalId                 = "external-id"
	testInventoryLevel             = "level1"
	testDatacenterName             = "datacenter1"
	testConsoleKeymap              = "us"
	testEnableNetworkTypeSelection = true

	testCertificateKey                = "certificateProvider"
	testApplianceUrlKey               = "applianceUrl"
	testExternalIdKey                 = "externalId"
	testInventoryLevelKey             = "inventoryLevel"
	testDatacenterNameKey             = "datacenterName"
	testConsoleKeymapKey              = "consoleKeymap"
	testEnableNetworkTypeSelectionKey = "enableNetworkTypeSelection"
)

// examplePlanObject is a struct that implements the necessary interface for testing
type examplePlanObject struct {
	objectValue  basetypes.ObjectValue
	dynamicValue basetypes.DynamicValue
}

func (e *examplePlanObject) getObjectValue() basetypes.ObjectValue {
	return e.objectValue
}

func (e *examplePlanObject) getDynamicValue() basetypes.DynamicValue {
	return e.dynamicValue
}

// exampleFromApi is a struct that simulates the API response for testing
type exampleFromApi struct {
	data       map[string]interface{}
	configType sdk.MappedNullable
}

func (e *exampleFromApi) getData() map[string]interface{} {
	return e.data
}

func (e *exampleFromApi) getConfigType() sdk.MappedNullable {
	return e.configType
}

// createExamplePlanObjectCloudHVMCamelCase creates an example plan object for testing
func createExamplePlanObjectCloudHVMCamelCase(t *testing.T) *examplePlanObject {
	// Define a sample plan as basetypes.ObjectValue
	planAttrTypes := map[string]attr.Type{
		testApplianceUrlKey:               basetypes.StringType{},
		testCertificateKey:                basetypes.StringType{},
		testConsoleKeymapKey:              basetypes.StringType{},
		testDatacenterNameKey:             basetypes.StringType{},
		testEnableNetworkTypeSelectionKey: basetypes.BoolType{},
		testExternalIdKey:                 basetypes.StringType{},
		testInventoryLevelKey:             basetypes.StringType{},
	}
	planAttrValues := map[string]attr.Value{
		testApplianceUrlKey:               basetypes.NewStringValue(testApplianceUrl),
		testCertificateKey:                basetypes.NewStringValue(testCertificateProvider),
		testConsoleKeymapKey:              basetypes.NewStringValue(testConsoleKeymap),
		testDatacenterNameKey:             basetypes.NewStringValue(testDatacenterName),
		testEnableNetworkTypeSelectionKey: basetypes.NewBoolValue(testEnableNetworkTypeSelection),
		testExternalIdKey:                 basetypes.NewStringValue(testExternalId),
		testInventoryLevelKey:             basetypes.NewStringValue(testInventoryLevel),
	}

	planObject, diags := basetypes.NewObjectValue(planAttrTypes, planAttrValues)
	assert.False(t, diags.HasError())

	dynamicValue := basetypes.NewDynamicValue(planObject)

	return &examplePlanObject{
		objectValue:  planObject,
		dynamicValue: dynamicValue,
	}
}

// createExamplePlanObjectCloudHVMSnakeCase creates an example plan object for testing
func createExamplePlanObjectCloudHVMSnakeCase(t *testing.T) *examplePlanObject {
	// Define a sample plan as basetypes.ObjectValue
	planAttrTypes := map[string]attr.Type{
		strcase.ToSnake(testApplianceUrlKey):               basetypes.StringType{},
		strcase.ToSnake(testCertificateKey):                basetypes.StringType{},
		strcase.ToSnake(testConsoleKeymapKey):              basetypes.StringType{},
		strcase.ToSnake(testDatacenterNameKey):             basetypes.StringType{},
		strcase.ToSnake(testEnableNetworkTypeSelectionKey): basetypes.BoolType{},
		strcase.ToSnake(testExternalIdKey):                 basetypes.StringType{},
		strcase.ToSnake(testInventoryLevelKey):             basetypes.StringType{},
	}
	planAttrValues := map[string]attr.Value{
		strcase.ToSnake(testApplianceUrlKey):               basetypes.NewStringValue(testApplianceUrl),
		strcase.ToSnake(testCertificateKey):                basetypes.NewStringValue(testCertificateProvider),
		strcase.ToSnake(testConsoleKeymapKey):              basetypes.NewStringValue(testConsoleKeymap),
		strcase.ToSnake(testDatacenterNameKey):             basetypes.NewStringValue(testDatacenterName),
		strcase.ToSnake(testEnableNetworkTypeSelectionKey): basetypes.NewBoolValue(testEnableNetworkTypeSelection),
		strcase.ToSnake(testExternalIdKey):                 basetypes.NewStringValue(testExternalId),
		strcase.ToSnake(testInventoryLevelKey):             basetypes.NewStringValue(testInventoryLevel),
	}

	planObject, diags := basetypes.NewObjectValue(planAttrTypes, planAttrValues)
	assert.False(t, diags.HasError())

	dynamicValue := basetypes.NewDynamicValue(planObject)

	return &examplePlanObject{
		objectValue:  planObject,
		dynamicValue: dynamicValue,
	}
}

// createExamplePlanObjectGeneric creates an example plan object for testing
func createExamplePlanObjectGeneric(t *testing.T) *examplePlanObject {
	// Define a sample plan as basetypes.ObjectValue
	bigF := &big.Float{}
	planAttrTypes := map[string]attr.Type{
		testKey1: basetypes.StringType{},
		testKey2: basetypes.NumberType{},
	}
	planAttrValues := map[string]attr.Value{
		testKey1: basetypes.NewStringValue(testValue1),
		testKey2: basetypes.NewNumberValue(bigF.SetFloat64(testValue2)),
	}
	planObject, diags := basetypes.NewObjectValue(planAttrTypes, planAttrValues)
	assert.False(t, diags.HasError())

	dynamicValue := basetypes.NewDynamicValue(planObject)

	return &examplePlanObject{
		objectValue:  planObject,
		dynamicValue: dynamicValue,
	}
}

// createExampleFromApiCloudHVM creates an example API response for testing
func createExampleFromApiCloudHVM(t *testing.T) *exampleFromApi {
	// Define a sample configType
	enableNetworkTypeSelection := convert.BoolToStringOnOff(testEnableNetworkTypeSelection)
	enableNTS := enableNetworkTypeSelection.ValueString()
	certProvider := testCertificateProvider
	applianceUrl := testApplianceUrl
	consoleKeymap := testConsoleKeymap
	datacenterName := testDatacenterName
	externalId := testExternalId
	inventoryLevel := testInventoryLevel
	hvmConfig := sdkfuncs.NewHvmCloudConfig()
	hvmConfig.CertificateProvider = &certProvider
	hvmConfig.ApplianceUrl = &applianceUrl
	hvmConfig.ConsoleKeymap = &consoleKeymap
	hvmConfig.DatacenterName = &datacenterName
	hvmConfig.EnableNetworkTypeSelection = *sdk.NewNullableString(&enableNTS)
	hvmConfig.ExternalId = *sdk.NewNullableString(&externalId)
	hvmConfig.InventoryLevel = &inventoryLevel

	fromApiData, err := hvmConfig.ToMap()
	assert.NoError(t, err)

	return &exampleFromApi{
		data:       fromApiData,
		configType: hvmConfig,
	}
}

func createPlanKeyMapHVM() map[string]string {
	// Create a map of plan keys to API keys for HVM config
	return map[string]string{
		strcase.ToSnake(testCertificateKey):                testCertificateKey,
		strcase.ToSnake(testApplianceUrlKey):               testApplianceUrlKey,
		strcase.ToSnake(testConsoleKeymapKey):              testConsoleKeymapKey,
		strcase.ToSnake(testDatacenterNameKey):             testDatacenterNameKey,
		strcase.ToSnake(testEnableNetworkTypeSelectionKey): testEnableNetworkTypeSelectionKey,
		strcase.ToSnake(testExternalIdKey):                 testExternalIdKey,
		strcase.ToSnake(testInventoryLevelKey):             testInventoryLevelKey,
	}
}

func diagnosticsKeysPlanFromApi(plan basetypes.ObjectValue, fromApi map[string]any) diag.Diagnostics {
	planAttrs := plan.Attributes()
	var diags diag.Diagnostics
	for k := range planAttrs {
		diags.AddWarning(
			"check config",
			fmt.Sprintf("key '%s' in plan not found in API response", k),
		)
	}

	for k := range fromApi {
		diags.AddWarning(
			"check config",
			fmt.Sprintf("key '%s' in API response not found in plan", k),
		)
	}

	return diags
}

func keysNotInPlan(plan basetypes.ObjectValue, fromApi map[string]any) map[string]any {
	planAttrs := plan.Attributes()
	// Build map of keys in planAttrs
	planKeys := make(map[string]struct{})
	for k := range planAttrs {
		planKeys[k] = struct{}{}
	}
	// Build map of keys in fromApiMap
	fromApiKeys := make(map[string]struct{})
	for k := range fromApi {
		fromApiKeys[k] = struct{}{}
	}

	// Now compare the keys in both maps
	// Find keys in fromApiMap that are not in planAttribute
	apiKeysNotInPlan := make(map[string]any)
	for k, v := range fromApiKeys {
		if _, ok := planKeys[k]; !ok {
			apiKeysNotInPlan[k] = v
		}
	}

	return apiKeysNotInPlan
}

func TestCheckPlanAttributeAgainstAPIAttribute(t *testing.T) {
	// create example plan and fromApi
	planGeneric := createExamplePlanObjectGeneric(t)
	planHVMCamelCase := createExamplePlanObjectCloudHVMCamelCase(t)
	planHVMSnakeCase := createExamplePlanObjectCloudHVMSnakeCase(t)
	fromApiHVM := createExampleFromApiCloudHVM(t)

	tests := []struct {
		name          string
		plan          any
		fromApi       any
		keyMap        map[string]string
		keysNotInPlan map[string]any
		wantDiags     diag.Diagnostics
	}{
		{
			name:          "HVM Camel case plan and fromApi match exactly",
			plan:          planHVMCamelCase.getObjectValue(),
			fromApi:       fromApiHVM.getConfigType(),
			keyMap:        nil,
			keysNotInPlan: map[string]any{},
			wantDiags:     nil,
		},
		{
			name:          "HVM Camel case plan Dynamic and fromApi match exactly",
			plan:          planHVMCamelCase.getDynamicValue(),
			fromApi:       fromApiHVM.getConfigType(),
			keyMap:        nil,
			keysNotInPlan: map[string]any{},
			wantDiags:     nil,
		},
		{
			name:          "HVM Camel case plan and fromApi Map match exactly",
			plan:          planHVMCamelCase.getObjectValue(),
			fromApi:       fromApiHVM.getData(),
			keyMap:        nil,
			keysNotInPlan: map[string]any{},
			wantDiags:     nil,
		},
		{
			name:          "HVM Camel case plan Dynamic and fromApi Map match exactly",
			plan:          planHVMCamelCase.getDynamicValue(),
			fromApi:       fromApiHVM.getData(),
			keyMap:        nil,
			keysNotInPlan: map[string]any{},
			wantDiags:     nil,
		},
		{
			name:          "HVM SnakeCase plan and fromApi match exactly with keyMap",
			plan:          planHVMSnakeCase.getObjectValue(),
			fromApi:       fromApiHVM.getConfigType(),
			keyMap:        createPlanKeyMapHVM(),
			keysNotInPlan: map[string]any{},
			wantDiags:     nil,
		},
		{
			name:          "HVM SnakeCase plan Dynamic and fromApi match exactly with keyMap",
			plan:          planHVMSnakeCase.getObjectValue(),
			fromApi:       fromApiHVM.getConfigType(),
			keyMap:        createPlanKeyMapHVM(),
			keysNotInPlan: map[string]any{},
			wantDiags:     nil,
		},
		{
			name:          "HVM SnakeCase plan and fromApi Map match exactly with keyMap",
			plan:          planHVMSnakeCase.getObjectValue(),
			fromApi:       fromApiHVM.getData(),
			keyMap:        createPlanKeyMapHVM(),
			keysNotInPlan: map[string]any{},
			wantDiags:     nil,
		},
		{
			name:          "HVM SnakeCase plan Dynamic and fromApi Map match exactly with keyMap",
			plan:          planHVMSnakeCase.getDynamicValue(),
			fromApi:       fromApiHVM.getData(),
			keyMap:        createPlanKeyMapHVM(),
			keysNotInPlan: map[string]any{},
			wantDiags:     nil,
		},
		{
			name:          "HVM SnakeCase plan Dynamic and fromApi Map no keyMap returns diags",
			plan:          planHVMSnakeCase.getDynamicValue(),
			fromApi:       fromApiHVM.getData(),
			keyMap:        nil,
			keysNotInPlan: keysNotInPlan(planHVMSnakeCase.getObjectValue(), fromApiHVM.getData()),
			wantDiags:     diagnosticsKeysPlanFromApi(planHVMSnakeCase.getObjectValue(), fromApiHVM.getData()),
		},
		{
			name:          "Generic plan Dynamic and fromApi Map no keyMap returns diags",
			plan:          planGeneric.getDynamicValue(),
			fromApi:       fromApiHVM.getData(),
			keyMap:        nil,
			keysNotInPlan: keysNotInPlan(planGeneric.getObjectValue(), fromApiHVM.getData()),
			wantDiags:     diagnosticsKeysPlanFromApi(planGeneric.getObjectValue(), fromApiHVM.getData()),
		},
		{
			name:          "Invalid plan type returns diag",
			plan:          12345,
			fromApi:       fromApiHVM.getConfigType(),
			keyMap:        nil,
			keysNotInPlan: nil,
			wantDiags: diag.Diagnostics{
				diag.NewWarningDiagnostic(
					"check config",
					"expected planAttribute to be basetypes.DynamicValue or basetypes.ObjectValuable, got int",
				),
			},
		},
		{
			name:          "Invalid fromApi type returns diag",
			plan:          planHVMCamelCase.getObjectValue(),
			fromApi:       12345,
			keyMap:        nil,
			keysNotInPlan: nil,
			wantDiags: diag.Diagnostics{
				diag.NewWarningDiagnostic(
					"check config",
					"expected apiAttribute to be map[string]any or sdk.MappedNullable, got int",
				),
			},
		},
		{
			name:          "Invalid plan Dynamic with non-object value returns diag",
			plan:          basetypes.NewDynamicValue(basetypes.NewStringValue("not-an-object")),
			fromApi:       fromApiHVM.getConfigType(),
			keyMap:        nil,
			keysNotInPlan: nil,
			wantDiags: diag.Diagnostics{
				diag.NewWarningDiagnostic(
					"check config",
					"expected planAttribute to be basetypes.DynamicValue of ObjectType,"+
						" got basetypes.StringValue",
				),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keysNotInPlan, gotDiags := CheckPlanAttributeAgainstAPIAttribute(
				context.Background(), tt.plan, tt.fromApi, tt.keyMap)
			assert.ElementsMatch(t, tt.wantDiags, gotDiags)
			assert.Equal(t, tt.keysNotInPlan, keysNotInPlan)
		})
	}
}
