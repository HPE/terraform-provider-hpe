# ImageBuilds

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Account** | Pointer to [**ImageBuildsAccount**](ImageBuildsAccount.md) |  | [optional] 
**Type** | Pointer to [**ImageBuildsType**](ImageBuildsType.md) |  | [optional] 
**Site** | Pointer to [**ImageBuildsSite**](ImageBuildsSite.md) |  | [optional] 
**Zone** | Pointer to [**ImageBuildsZone**](ImageBuildsZone.md) |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **string** |  | [optional] 
**BootScript** | Pointer to [**ImageBuildsBootScript**](ImageBuildsBootScript.md) |  | [optional] 
**BootCommand** | Pointer to **NullableString** |  | [optional] 
**PreseedScript** | Pointer to [**ImageBuildsPreseedScript**](ImageBuildsPreseedScript.md) |  | [optional] 
**Scripts** | Pointer to [**[]ImageBuildsScriptsInner**](ImageBuildsScriptsInner.md) |  | [optional] 
**SshUsername** | Pointer to **string** |  | [optional] 
**SshPassword** | Pointer to **string** |  | [optional] 
**StorageProvider** | Pointer to **NullableString** |  | [optional] 
**BuildOutputName** | Pointer to **string** |  | [optional] 
**ConversionFormats** | Pointer to **NullableString** |  | [optional] 
**IsCloudInit** | Pointer to **bool** |  | [optional] 
**VmToolsInstalled** | Pointer to **bool** |  | [optional] 
**KeepResults** | Pointer to **int64** |  | [optional] 
**Config** | Pointer to [**ImageBuildsConfig**](ImageBuildsConfig.md) |  | [optional] 
**LastResult** | Pointer to [**ImageBuildsLastResult**](ImageBuildsLastResult.md) |  | [optional] 
**ExecutionCount** | Pointer to **int64** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ImageBuilds{
    // Set fields directly
}
```

### BootCommand (Nullable)

Use the Nullable wrapper methods:
- `obj.BootCommand.IsSet()` — check if set
- `obj.BootCommand.Get()` — get the inner value (returns pointer)
- `obj.BootCommand.Set(&val)` — set the value
- `obj.BootCommand.Unset()` — clear the value
### StorageProvider (Nullable)

Use the Nullable wrapper methods:
- `obj.StorageProvider.IsSet()` — check if set
- `obj.StorageProvider.Get()` — get the inner value (returns pointer)
- `obj.StorageProvider.Set(&val)` — set the value
- `obj.StorageProvider.Unset()` — clear the value
### ConversionFormats (Nullable)

Use the Nullable wrapper methods:
- `obj.ConversionFormats.IsSet()` — check if set
- `obj.ConversionFormats.Get()` — get the inner value (returns pointer)
- `obj.ConversionFormats.Set(&val)` — set the value
- `obj.ConversionFormats.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


