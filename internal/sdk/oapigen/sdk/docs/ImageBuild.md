# ImageBuild

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Account** | Pointer to [**GetImageBuild200ResponseImageBuildAccount**](GetImageBuild200ResponseImageBuildAccount.md) |  | [optional] 
**Type** | Pointer to [**GetImageBuild200ResponseImageBuildType**](GetImageBuild200ResponseImageBuildType.md) |  | [optional] 
**Site** | Pointer to [**GetImageBuild200ResponseImageBuildSite**](GetImageBuild200ResponseImageBuildSite.md) |  | [optional] 
**Zone** | Pointer to [**GetImageBuild200ResponseImageBuildZone**](GetImageBuild200ResponseImageBuildZone.md) |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**BootScript** | Pointer to [**GetImageBuild200ResponseImageBuildBootScript**](GetImageBuild200ResponseImageBuildBootScript.md) |  | [optional] 
**BootCommand** | Pointer to **NullableString** |  | [optional] 
**PreseedScript** | Pointer to [**GetImageBuild200ResponseImageBuildPreseedScript**](GetImageBuild200ResponseImageBuildPreseedScript.md) |  | [optional] 
**Scripts** | Pointer to [**[]GetImageBuild200ResponseImageBuildScriptsInner**](GetImageBuild200ResponseImageBuildScriptsInner.md) |  | [optional] 
**SshUsername** | Pointer to **string** |  | [optional] 
**SshPassword** | Pointer to **string** |  | [optional] 
**StorageProvider** | Pointer to **NullableString** |  | [optional] 
**BuildOutputName** | Pointer to **NullableString** |  | [optional] 
**ConversionFormats** | Pointer to **NullableString** |  | [optional] 
**IsCloudInit** | Pointer to **bool** |  | [optional] 
**VmToolsInstalled** | Pointer to **bool** |  | [optional] 
**KeepResults** | Pointer to **NullableInt64** |  | [optional] 
**Config** | Pointer to [**GetImageBuild200ResponseImageBuildConfig**](GetImageBuild200ResponseImageBuildConfig.md) |  | [optional] 
**LastResult** | Pointer to [**GetImageBuild200ResponseImageBuildLastResult**](GetImageBuild200ResponseImageBuildLastResult.md) |  | [optional] 
**ExecutionCount** | Pointer to **int64** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ImageBuild{
    // Set fields directly
}
```

### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value
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
### BuildOutputName (Nullable)

Use the Nullable wrapper methods:
- `obj.BuildOutputName.IsSet()` — check if set
- `obj.BuildOutputName.Get()` — get the inner value (returns pointer)
- `obj.BuildOutputName.Set(&val)` — set the value
- `obj.BuildOutputName.Unset()` — clear the value
### ConversionFormats (Nullable)

Use the Nullable wrapper methods:
- `obj.ConversionFormats.IsSet()` — check if set
- `obj.ConversionFormats.Get()` — get the inner value (returns pointer)
- `obj.ConversionFormats.Set(&val)` — set the value
- `obj.ConversionFormats.Unset()` — clear the value
### KeepResults (Nullable)

Use the Nullable wrapper methods:
- `obj.KeepResults.IsSet()` — check if set
- `obj.KeepResults.Get()` — get the inner value (returns pointer)
- `obj.KeepResults.Set(&val)` — set the value
- `obj.KeepResults.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


