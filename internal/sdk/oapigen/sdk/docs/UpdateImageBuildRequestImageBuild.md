# UpdateImageBuildRequestImageBuild

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | A name for the image build | [optional] 
**Description** | Pointer to **NullableString** | A description for the image build | [optional] 
**Type** | Pointer to **string** | The image builder type. | [optional] 
**Site** | Pointer to [**UpdateImageBuildRequestImageBuildSite**](UpdateImageBuildRequestImageBuildSite.md) |  | [optional] 
**Zone** | Pointer to [**UpdateImageBuildRequestImageBuildZone**](UpdateImageBuildRequestImageBuildZone.md) |  | [optional] 
**Config** | Pointer to **map[string]interface{}** | A map of config values. This is the instance config that is used for provisioning. See Provisioning Types. | [optional] 
**BootScript** | Pointer to [**UpdateImageBuildRequestImageBuildBootScript**](UpdateImageBuildRequestImageBuildBootScript.md) |  | [optional] 
**PreseedScript** | Pointer to [**UpdateImageBuildRequestImageBuildPreseedScript**](UpdateImageBuildRequestImageBuildPreseedScript.md) |  | [optional] 
**SshUsername** | Pointer to **string** | SSH Username | [optional] 
**SshPassword** | Pointer to **string** | SSH Password | [optional] 
**StorageProvider** | Pointer to **NullableString** | Storage Provider | [optional] 
**IsCloudInit** | Pointer to **string** | Cloud Init | [optional] 
**BuildOutputName** | Pointer to **NullableString** | Build Output Name | [optional] 
**ConversionFormats** | Pointer to **NullableString** |  | [optional] 
**KeepResults** | Pointer to **int64** | Keep Results - Keep only the most recent builds. Older executions will be deleted along with their associated Virtual Images. The value 0 disables this functionality. | [optional] [default to 0]

## Usage

Instantiate with a Go composite literal:

```go
obj := &UpdateImageBuildRequestImageBuild{
    // Set fields directly
}
```

### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value
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

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


