# ListImageBuilds200ResponseAllOfImageBuildsInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Account** | Pointer to [**NullableListImageBuilds200ResponseAllOfImageBuildsInnerAccount**](ListImageBuilds200ResponseAllOfImageBuildsInnerAccount.md) |  | [optional] 
**Type** | Pointer to [**ListImageBuilds200ResponseAllOfImageBuildsInnerType**](ListImageBuilds200ResponseAllOfImageBuildsInnerType.md) |  | [optional] 
**Site** | Pointer to [**NullableListImageBuilds200ResponseAllOfImageBuildsInnerSite**](ListImageBuilds200ResponseAllOfImageBuildsInnerSite.md) |  | [optional] 
**Zone** | Pointer to [**NullableListImageBuilds200ResponseAllOfImageBuildsInnerZone**](ListImageBuilds200ResponseAllOfImageBuildsInnerZone.md) |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **string** |  | [optional] 
**BootScript** | Pointer to [**ListImageBuilds200ResponseAllOfImageBuildsInnerBootScript**](ListImageBuilds200ResponseAllOfImageBuildsInnerBootScript.md) |  | [optional] 
**BootCommand** | Pointer to **NullableString** |  | [optional] 
**PreseedScript** | Pointer to [**ListImageBuilds200ResponseAllOfImageBuildsInnerPreseedScript**](ListImageBuilds200ResponseAllOfImageBuildsInnerPreseedScript.md) |  | [optional] 
**Scripts** | Pointer to [**[]ListImageBuilds200ResponseAllOfImageBuildsInnerScriptsInner**](ListImageBuilds200ResponseAllOfImageBuildsInnerScriptsInner.md) |  | [optional] 
**SshUsername** | Pointer to **string** |  | [optional] 
**SshPassword** | Pointer to **string** |  | [optional] 
**StorageProvider** | Pointer to **NullableString** |  | [optional] 
**BuildOutputName** | Pointer to **string** |  | [optional] 
**ConversionFormats** | Pointer to **NullableString** |  | [optional] 
**IsCloudInit** | Pointer to **bool** |  | [optional] 
**VmToolsInstalled** | Pointer to **bool** |  | [optional] 
**KeepResults** | Pointer to **int64** |  | [optional] 
**Config** | Pointer to [**ListImageBuilds200ResponseAllOfImageBuildsInnerConfig**](ListImageBuilds200ResponseAllOfImageBuildsInnerConfig.md) |  | [optional] 
**LastResult** | Pointer to [**ListImageBuilds200ResponseAllOfImageBuildsInnerLastResult**](ListImageBuilds200ResponseAllOfImageBuildsInnerLastResult.md) |  | [optional] 
**ExecutionCount** | Pointer to **int64** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ListImageBuilds200ResponseAllOfImageBuildsInner{
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


