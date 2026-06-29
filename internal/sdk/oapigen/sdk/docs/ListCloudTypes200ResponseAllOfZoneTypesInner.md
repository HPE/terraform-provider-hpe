# ListCloudTypes200ResponseAllOfZoneTypesInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Code** | Pointer to **string** |  | [optional] 
**Enabled** | Pointer to **bool** |  | [optional] 
**Provision** | Pointer to **bool** |  | [optional] 
**AutoCapacity** | Pointer to **bool** |  | [optional] 
**MigrationTarget** | Pointer to **bool** |  | [optional] 
**HasAffinityGroups** | Pointer to **bool** |  | [optional] 
**HasDatastores** | Pointer to **bool** |  | [optional] 
**HasNetworks** | Pointer to **bool** |  | [optional] 
**HasResourcePools** | Pointer to **bool** |  | [optional] 
**HasSecurityGroups** | Pointer to **bool** |  | [optional] 
**HasContainers** | Pointer to **bool** |  | [optional] 
**HasBareMetal** | Pointer to **bool** |  | [optional] 
**HasServices** | Pointer to **bool** |  | [optional] 
**HasFunctions** | Pointer to **bool** |  | [optional] 
**HasJobs** | Pointer to **bool** |  | [optional] 
**HasDiscovery** | Pointer to **bool** |  | [optional] 
**HasCloudInit** | Pointer to **bool** |  | [optional] 
**HasFolders** | Pointer to **bool** |  | [optional] 
**HasMarketplace** | Pointer to **bool** |  | [optional] 
**HasNativePlans** | Pointer to **bool** |  | [optional] 
**CanCreateResourcePools** | Pointer to **bool** |  | [optional] 
**CanDeleteResourcePools** | Pointer to **bool** |  | [optional] 
**CanCreateDatastores** | Pointer to **bool** |  | [optional] 
**CanCreateNetworks** | Pointer to **bool** |  | [optional] 
**CanChooseContainerMode** | Pointer to **bool** |  | [optional] 
**ProvisionRequiresResourcePool** | Pointer to **bool** |  | [optional] 
**SupportsDistributedWorker** | Pointer to **bool** |  | [optional] 
**Cloud** | Pointer to **string** |  | [optional] 
**ProvisionTypes** | Pointer to **[]int64** |  | [optional] 
**ZoneInstanceTypeLayoutId** | Pointer to **int64** |  | [optional] 
**ServerTypes** | Pointer to [**[]ListCloudTypes200ResponseAllOfZoneTypesInnerServerTypesInner**](ListCloudTypes200ResponseAllOfZoneTypesInnerServerTypesInner.md) |  | [optional] 
**OptionTypes** | Pointer to [**[]ListCloudTypes200ResponseAllOfZoneTypesInnerOptionTypesInner**](ListCloudTypes200ResponseAllOfZoneTypesInnerOptionTypesInner.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ListCloudTypes200ResponseAllOfZoneTypesInner{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


