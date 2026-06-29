# AddImageBuild200ResponseAllOfImageBuild

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Account** | Pointer to [**NullableAddImageBuild200ResponseAllOfImageBuildAccount**](AddImageBuild200ResponseAllOfImageBuildAccount.md) |  | [optional] 
**Type** | Pointer to [**AddImageBuild200ResponseAllOfImageBuildType**](AddImageBuild200ResponseAllOfImageBuildType.md) |  | [optional] 
**Site** | Pointer to [**NullableAddImageBuild200ResponseAllOfImageBuildSite**](AddImageBuild200ResponseAllOfImageBuildSite.md) |  | [optional] 
**Zone** | Pointer to [**NullableAddImageBuild200ResponseAllOfImageBuildZone**](AddImageBuild200ResponseAllOfImageBuildZone.md) |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**BootScript** | Pointer to [**AddImageBuild200ResponseAllOfImageBuildBootScript**](AddImageBuild200ResponseAllOfImageBuildBootScript.md) |  | [optional] 
**BootCommand** | Pointer to **NullableString** |  | [optional] 
**PreseedScript** | Pointer to [**AddImageBuild200ResponseAllOfImageBuildPreseedScript**](AddImageBuild200ResponseAllOfImageBuildPreseedScript.md) |  | [optional] 
**Scripts** | Pointer to [**[]AddImageBuild200ResponseAllOfImageBuildScriptsInner**](AddImageBuild200ResponseAllOfImageBuildScriptsInner.md) |  | [optional] 
**SshUsername** | Pointer to **string** |  | [optional] 
**SshPassword** | Pointer to **string** |  | [optional] 
**StorageProvider** | Pointer to **NullableString** |  | [optional] 
**BuildOutputName** | Pointer to **NullableString** |  | [optional] 
**ConversionFormats** | Pointer to **NullableString** |  | [optional] 
**IsCloudInit** | Pointer to **bool** |  | [optional] 
**VmToolsInstalled** | Pointer to **bool** |  | [optional] 
**KeepResults** | Pointer to **NullableInt64** |  | [optional] 
**Config** | Pointer to [**AddImageBuild200ResponseAllOfImageBuildConfig**](AddImageBuild200ResponseAllOfImageBuildConfig.md) |  | [optional] 
**LastResult** | Pointer to [**AddImageBuild200ResponseAllOfImageBuildLastResult**](AddImageBuild200ResponseAllOfImageBuildLastResult.md) |  | [optional] 
**ExecutionCount** | Pointer to **int64** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &AddImageBuild200ResponseAllOfImageBuild{
    // Set fields directly
}
```

### Account (Nullable)

Use the Nullable wrapper methods:
- `obj.Account.IsSet()` — check if set
- `obj.Account.Get()` — get the inner value (returns pointer)
- `obj.Account.Set(&val)` — set the value
- `obj.Account.Unset()` — clear the value
### Site (Nullable)

Use the Nullable wrapper methods:
- `obj.Site.IsSet()` — check if set
- `obj.Site.Get()` — get the inner value (returns pointer)
- `obj.Site.Set(&val)` — set the value
- `obj.Site.Unset()` — clear the value
### Zone (Nullable)

Use the Nullable wrapper methods:
- `obj.Zone.IsSet()` — check if set
- `obj.Zone.Get()` — get the inner value (returns pointer)
- `obj.Zone.Set(&val)` — set the value
- `obj.Zone.Unset()` — clear the value
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


