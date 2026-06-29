# ListClusterLayouts200ResponseAllOfLayoutsInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**InternalId** | Pointer to **string** |  | [optional] 
**ServerCount** | Pointer to **int64** |  | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 
**Code** | Pointer to **string** |  | [optional] 
**LastUpdated** | Pointer to **time.Time** |  | [optional] 
**HasAutoScale** | Pointer to **bool** |  | [optional] 
**MemoryRequirement** | Pointer to **int64** |  | [optional] 
**ClusterVersion** | Pointer to **string** |  | [optional] 
**ComputeVersion** | Pointer to **string** |  | [optional] 
**HasSettings** | Pointer to **bool** |  | [optional] 
**SortOrder** | Pointer to **int64** |  | [optional] 
**HasConfig** | Pointer to **bool** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Creatable** | Pointer to **bool** |  | [optional] 
**Enabled** | Pointer to **bool** |  | [optional] 
**Description** | Pointer to **string** |  | [optional] 
**GroupType** | Pointer to [**ListClusterLayouts200ResponseAllOfLayoutsInnerGroupType**](ListClusterLayouts200ResponseAllOfLayoutsInnerGroupType.md) |  | [optional] 
**Labels** | Pointer to **[]string** |  | [optional] 
**EnvironmentVariables** | Pointer to **[]map[string]interface{}** |  | [optional] 
**OptionTypes** | Pointer to [**[]ListClusterLayouts200ResponseAllOfLayoutsInnerOptionTypesInner**](ListClusterLayouts200ResponseAllOfLayoutsInnerOptionTypesInner.md) |  | [optional] 
**Actions** | Pointer to **[]map[string]interface{}** |  | [optional] 
**ComputeServers** | Pointer to [**[]ListClusterLayouts200ResponseAllOfLayoutsInnerComputeServersInner**](ListClusterLayouts200ResponseAllOfLayoutsInnerComputeServersInner.md) |  | [optional] 
**InstallContainerRuntime** | Pointer to **bool** |  | [optional] 
**ProvisionType** | Pointer to [**ListClusterLayouts200ResponseAllOfLayoutsInnerProvisionType**](ListClusterLayouts200ResponseAllOfLayoutsInnerProvisionType.md) |  | [optional] 
**SpecTemplates** | Pointer to [**[]ListClusterLayouts200ResponseAllOfLayoutsInnerSpecTemplatesInner**](ListClusterLayouts200ResponseAllOfLayoutsInnerSpecTemplatesInner.md) |  | [optional] 
**TaskSets** | Pointer to **[]map[string]interface{}** |  | [optional] 
**Type** | Pointer to [**ListClusterLayouts200ResponseAllOfLayoutsInnerType**](ListClusterLayouts200ResponseAllOfLayoutsInnerType.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ListClusterLayouts200ResponseAllOfLayoutsInner{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


