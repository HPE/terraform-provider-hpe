# GetInstanceType200ResponseInstanceType

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Account** | Pointer to [**GetInstanceType200ResponseInstanceTypeAccount**](GetInstanceType200ResponseInstanceTypeAccount.md) |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Labels** | Pointer to **[]string** |  | [optional] 
**Code** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**ProvisionTypeCode** | Pointer to **NullableString** |  | [optional] 
**Category** | Pointer to **string** |  | [optional] 
**Active** | Pointer to **bool** |  | [optional] 
**HasProvisioningStep** | Pointer to **bool** |  | [optional] 
**HasDeployment** | Pointer to **bool** |  | [optional] 
**HasConfig** | Pointer to **bool** |  | [optional] 
**HasSettings** | Pointer to **bool** |  | [optional] 
**HasAutoScale** | Pointer to **bool** |  | [optional] 
**ProxyType** | Pointer to **NullableString** |  | [optional] 
**ProxyPort** | Pointer to **NullableString** |  | [optional] 
**ProxyProtocol** | Pointer to **NullableString** |  | [optional] 
**EnvironmentPrefix** | Pointer to **string** |  | [optional] 
**BackupType** | Pointer to **NullableString** |  | [optional] 
**Config** | Pointer to **map[string]interface{}** |  | [optional] 
**Visibility** | Pointer to **string** |  | [optional] 
**Featured** | Pointer to **bool** |  | [optional] 
**Versions** | Pointer to **[]string** |  | [optional] 
**InstanceTypeLayouts** | Pointer to [**[]GetInstanceType200ResponseInstanceTypeInstanceTypeLayoutsInner**](GetInstanceType200ResponseInstanceTypeInstanceTypeLayoutsInner.md) |  | [optional] 
**OptionTypes** | Pointer to [**[]GetInstanceType200ResponseInstanceTypeOptionTypesInner**](GetInstanceType200ResponseInstanceTypeOptionTypesInner.md) |  | [optional] 
**EnvironmentVariables** | Pointer to **[]map[string]interface{}** |  | [optional] 
**PriceSets** | Pointer to **[]map[string]interface{}** |  | [optional] 
**ImagePath** | Pointer to **NullableString** | Logo image URL | [optional] 
**DarkImagePath** | Pointer to **NullableString** | Dark logo image URL | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetInstanceType200ResponseInstanceType{
    // Set fields directly
}
```

### Labels (Nullable)

Use the Nullable wrapper methods:
- `obj.Labels.IsSet()` — check if set
- `obj.Labels.Get()` — get the inner value (returns pointer)
- `obj.Labels.Set(&val)` — set the value
- `obj.Labels.Unset()` — clear the value
### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value
### ProvisionTypeCode (Nullable)

Use the Nullable wrapper methods:
- `obj.ProvisionTypeCode.IsSet()` — check if set
- `obj.ProvisionTypeCode.Get()` — get the inner value (returns pointer)
- `obj.ProvisionTypeCode.Set(&val)` — set the value
- `obj.ProvisionTypeCode.Unset()` — clear the value
### ProxyType (Nullable)

Use the Nullable wrapper methods:
- `obj.ProxyType.IsSet()` — check if set
- `obj.ProxyType.Get()` — get the inner value (returns pointer)
- `obj.ProxyType.Set(&val)` — set the value
- `obj.ProxyType.Unset()` — clear the value
### ProxyPort (Nullable)

Use the Nullable wrapper methods:
- `obj.ProxyPort.IsSet()` — check if set
- `obj.ProxyPort.Get()` — get the inner value (returns pointer)
- `obj.ProxyPort.Set(&val)` — set the value
- `obj.ProxyPort.Unset()` — clear the value
### ProxyProtocol (Nullable)

Use the Nullable wrapper methods:
- `obj.ProxyProtocol.IsSet()` — check if set
- `obj.ProxyProtocol.Get()` — get the inner value (returns pointer)
- `obj.ProxyProtocol.Set(&val)` — set the value
- `obj.ProxyProtocol.Unset()` — clear the value
### BackupType (Nullable)

Use the Nullable wrapper methods:
- `obj.BackupType.IsSet()` — check if set
- `obj.BackupType.Get()` — get the inner value (returns pointer)
- `obj.BackupType.Set(&val)` — set the value
- `obj.BackupType.Unset()` — clear the value
### Config (Nullable)

Use the Nullable wrapper methods:
- `obj.Config.IsSet()` — check if set
- `obj.Config.Get()` — get the inner value (returns pointer)
- `obj.Config.Set(&val)` — set the value
- `obj.Config.Unset()` — clear the value
### EnvironmentVariables (Nullable)

Use the Nullable wrapper methods:
- `obj.EnvironmentVariables.IsSet()` — check if set
- `obj.EnvironmentVariables.Get()` — get the inner value (returns pointer)
- `obj.EnvironmentVariables.Set(&val)` — set the value
- `obj.EnvironmentVariables.Unset()` — clear the value
### PriceSets (Nullable)

Use the Nullable wrapper methods:
- `obj.PriceSets.IsSet()` — check if set
- `obj.PriceSets.Get()` — get the inner value (returns pointer)
- `obj.PriceSets.Set(&val)` — set the value
- `obj.PriceSets.Unset()` — clear the value
### ImagePath (Nullable)

Use the Nullable wrapper methods:
- `obj.ImagePath.IsSet()` — check if set
- `obj.ImagePath.Get()` — get the inner value (returns pointer)
- `obj.ImagePath.Set(&val)` — set the value
- `obj.ImagePath.Unset()` — clear the value
### DarkImagePath (Nullable)

Use the Nullable wrapper methods:
- `obj.DarkImagePath.IsSet()` — check if set
- `obj.DarkImagePath.Get()` — get the inner value (returns pointer)
- `obj.DarkImagePath.Set(&val)` — set the value
- `obj.DarkImagePath.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


