# GetInstanceTypeProvisioning200ResponseAllOfInstanceTypeInstanceTypeLayoutsInnerContainerTypesInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int32** |  | [optional] 
**Account** | Pointer to [**GetInstanceTypeProvisioning200ResponseAllOfInstanceTypeInstanceTypeLayoutsInnerContainerTypesInnerAccount**](GetInstanceTypeProvisioning200ResponseAllOfInstanceTypeInstanceTypeLayoutsInnerContainerTypesInnerAccount.md) |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Labels** | Pointer to **[]string** |  | [optional] 
**ShortName** | Pointer to **string** |  | [optional] 
**Code** | Pointer to **string** |  | [optional] 
**ContainerVersion** | Pointer to **string** |  | [optional] 
**ProvisionType** | Pointer to [**GetInstanceTypeProvisioning200ResponseAllOfInstanceTypeInstanceTypeLayoutsInnerContainerTypesInnerProvisionType**](GetInstanceTypeProvisioning200ResponseAllOfInstanceTypeInstanceTypeLayoutsInnerContainerTypesInnerProvisionType.md) |  | [optional] 
**VirtualImage** | Pointer to [**GetInstanceTypeProvisioning200ResponseAllOfInstanceTypeInstanceTypeLayoutsInnerContainerTypesInnerVirtualImage**](GetInstanceTypeProvisioning200ResponseAllOfInstanceTypeInstanceTypeLayoutsInnerContainerTypesInnerVirtualImage.md) |  | [optional] 
**OsType** | Pointer to [**GetInstanceTypeProvisioning200ResponseAllOfInstanceTypeInstanceTypeLayoutsInnerContainerTypesInnerOsType**](GetInstanceTypeProvisioning200ResponseAllOfInstanceTypeInstanceTypeLayoutsInnerContainerTypesInnerOsType.md) |  | [optional] 
**Category** | Pointer to **NullableString** |  | [optional] 
**Config** | Pointer to **map[string]interface{}** |  | [optional] 
**ContainerPorts** | Pointer to [**[]GetInstanceTypeProvisioning200ResponseAllOfInstanceTypeInstanceTypeLayoutsInnerContainerTypesInnerContainerPortsInner**](GetInstanceTypeProvisioning200ResponseAllOfInstanceTypeInstanceTypeLayoutsInnerContainerTypesInnerContainerPortsInner.md) |  | [optional] 
**ContainerScripts** | Pointer to **[]map[string]interface{}** |  | [optional] 
**ContainerTemplates** | Pointer to **[]map[string]interface{}** |  | [optional] 
**EnvironmentVariables** | Pointer to **[]map[string]interface{}** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetInstanceTypeProvisioning200ResponseAllOfInstanceTypeInstanceTypeLayoutsInnerContainerTypesInner{
    // Set fields directly
}
```

### Category (Nullable)

Use the Nullable wrapper methods:
- `obj.Category.IsSet()` — check if set
- `obj.Category.Get()` — get the inner value (returns pointer)
- `obj.Category.Set(&val)` — set the value
- `obj.Category.Unset()` — clear the value
### Config (Nullable)

Use the Nullable wrapper methods:
- `obj.Config.IsSet()` — check if set
- `obj.Config.Get()` — get the inner value (returns pointer)
- `obj.Config.Set(&val)` — set the value
- `obj.Config.Unset()` — clear the value
### ContainerPorts (Nullable)

Use the Nullable wrapper methods:
- `obj.ContainerPorts.IsSet()` — check if set
- `obj.ContainerPorts.Get()` — get the inner value (returns pointer)
- `obj.ContainerPorts.Set(&val)` — set the value
- `obj.ContainerPorts.Unset()` — clear the value
### ContainerScripts (Nullable)

Use the Nullable wrapper methods:
- `obj.ContainerScripts.IsSet()` — check if set
- `obj.ContainerScripts.Get()` — get the inner value (returns pointer)
- `obj.ContainerScripts.Set(&val)` — set the value
- `obj.ContainerScripts.Unset()` — clear the value
### ContainerTemplates (Nullable)

Use the Nullable wrapper methods:
- `obj.ContainerTemplates.IsSet()` — check if set
- `obj.ContainerTemplates.Get()` — get the inner value (returns pointer)
- `obj.ContainerTemplates.Set(&val)` — set the value
- `obj.ContainerTemplates.Unset()` — clear the value
### EnvironmentVariables (Nullable)

Use the Nullable wrapper methods:
- `obj.EnvironmentVariables.IsSet()` — check if set
- `obj.EnvironmentVariables.Get()` — get the inner value (returns pointer)
- `obj.EnvironmentVariables.Set(&val)` — set the value
- `obj.EnvironmentVariables.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


