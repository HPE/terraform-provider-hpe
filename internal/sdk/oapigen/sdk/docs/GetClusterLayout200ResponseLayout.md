# GetClusterLayout200ResponseLayout

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
**GroupType** | Pointer to [**GetClusterLayout200ResponseLayoutGroupType**](GetClusterLayout200ResponseLayoutGroupType.md) |  | [optional] 
**Labels** | Pointer to **[]string** |  | [optional] 
**EnvironmentVariables** | Pointer to **[]map[string]interface{}** |  | [optional] 
**OptionTypes** | Pointer to [**[]GetClusterLayout200ResponseLayoutOptionTypesInner**](GetClusterLayout200ResponseLayoutOptionTypesInner.md) |  | [optional] 
**Actions** | Pointer to **[]map[string]interface{}** |  | [optional] 
**ComputeServers** | Pointer to [**[]GetClusterLayout200ResponseLayoutComputeServersInner**](GetClusterLayout200ResponseLayoutComputeServersInner.md) |  | [optional] 
**InstallContainerRuntime** | Pointer to **bool** |  | [optional] 
**ProvisionType** | Pointer to [**GetClusterLayout200ResponseLayoutProvisionType**](GetClusterLayout200ResponseLayoutProvisionType.md) |  | [optional] 
**SpecTemplates** | Pointer to [**[]GetClusterLayout200ResponseLayoutSpecTemplatesInner**](GetClusterLayout200ResponseLayoutSpecTemplatesInner.md) |  | [optional] 
**TaskSets** | Pointer to **[]map[string]interface{}** |  | [optional] 
**Type** | Pointer to [**GetClusterLayout200ResponseLayoutType**](GetClusterLayout200ResponseLayoutType.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetClusterLayout200ResponseLayout{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


