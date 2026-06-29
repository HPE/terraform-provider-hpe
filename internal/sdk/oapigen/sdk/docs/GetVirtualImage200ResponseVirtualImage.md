# GetVirtualImage200ResponseVirtualImage

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** | A name for the virtual image | [optional] 
**Description** | Pointer to **NullableString** | A description for the virtual image | [optional] 
**Labels** | Pointer to **[]string** | Array of label strings, can be used for filtering. | [optional] 
**OwnerId** | Pointer to **int64** | Owner of the image | [optional] 
**Tenant** | Pointer to [**GetVirtualImage200ResponseVirtualImageTenant**](GetVirtualImage200ResponseVirtualImageTenant.md) |  | [optional] 
**ImageType** | Pointer to **string** | Code of image type. eg. vmware, ami, etc. | [optional] 
**UserUploaded** | Pointer to **bool** | Is uploaded by an user? | [optional] 
**UserDefined** | Pointer to **bool** | Is defined by an user? | [optional] 
**SystemImage** | Pointer to **bool** | Is created by system? | [optional] 
**IsCloudInit** | Pointer to **bool** | Cloud Init Enabled? | [optional] 
**SshUsername** | Pointer to **NullableString** | SSH Username | [optional] 
**SshPassword** | Pointer to **NullableString** |  | [optional] 
**SshPasswordHash** | Pointer to **NullableString** |  | [optional] 
**SshKey** | Pointer to **NullableString** | SSH Key | [optional] 
**OsType** | Pointer to [**GetVirtualImage200ResponseVirtualImageOsType**](GetVirtualImage200ResponseVirtualImageOsType.md) |  | [optional] 
**MinRam** | Pointer to **NullableInt64** | Minimum required RAM in bytes | [optional] 
**MinRamGB** | Pointer to **NullableFloat64** | Minimum required RAM in gigabytes | [optional] 
**MinDisk** | Pointer to **NullableInt64** | Minimum required disk size in bytes | [optional] 
**MinDiskGB** | Pointer to **NullableString** | Minimum required disk size in gigabytes | [optional] 
**RawSize** | Pointer to **NullableInt64** | Size of image in bytes | [optional] 
**RawSizeGB** | Pointer to **NullableFloat32** | Size of image in gigabytes | [optional] 
**TrialVersion** | Pointer to **bool** | Is Trial Version? | [optional] 
**VirtioSupported** | Pointer to **bool** | VirtIO Drivers Loaded? | [optional] 
**Uefi** | Pointer to **NullableBool** | UEFI enabled? | [optional] 
**IsAutoJoinDomain** | Pointer to **bool** | Auto Join Domain? | [optional] 
**VmToolsInstalled** | Pointer to **bool** | VM Tools Installed? | [optional] 
**InstallAgent** | Pointer to **bool** | Install Agent? | [optional] 
**IsForceCustomization** | Pointer to **bool** | Force Guest Customization? | [optional] 
**IsSysprep** | Pointer to **bool** | Sysprep Enabled? | [optional] 
**FipsEnabled** | Pointer to **bool** | FIPS enabled? | [optional] 
**UserData** | Pointer to **NullableString** | Cloud-Init User Data, a bash script | [optional] 
**ConsoleKeymap** | Pointer to **NullableString** |  | [optional] 
**StorageProvider** | Pointer to [**GetVirtualImage200ResponseVirtualImageStorageProvider**](GetVirtualImage200ResponseVirtualImageStorageProvider.md) |  | [optional] 
**ExternalId** | Pointer to **NullableString** |  | [optional] 
**Visibility** | Pointer to **string** | private or public | [optional] 
**Accounts** | Pointer to [**[]GetVirtualImage200ResponseVirtualImageAccountsInner**](GetVirtualImage200ResponseVirtualImageAccountsInner.md) |  | [optional] 
**Config** | Pointer to [**GetVirtualImage200ResponseVirtualImageConfig**](GetVirtualImage200ResponseVirtualImageConfig.md) |  | [optional] 
**Volumes** | Pointer to **[]map[string]interface{}** |  | [optional] 
**StorageControllers** | Pointer to [**[]GetVirtualImage200ResponseVirtualImageStorageControllersInner**](GetVirtualImage200ResponseVirtualImageStorageControllersInner.md) |  | [optional] 
**NetworkInterfaces** | Pointer to **[]map[string]interface{}** |  | [optional] 
**Tags** | Pointer to [**[]GetVirtualImage200ResponseVirtualImageTagsInner**](GetVirtualImage200ResponseVirtualImageTagsInner.md) | Metadata tags, Array of objects having a name and value | [optional] 
**Locations** | Pointer to [**[]GetVirtualImage200ResponseVirtualImageLocationsInner**](GetVirtualImage200ResponseVirtualImageLocationsInner.md) |  | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 
**LastUpdated** | Pointer to **time.Time** |  | [optional] 
**Status** | Pointer to **string** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetVirtualImage200ResponseVirtualImage{
    // Set fields directly
}
```

### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value
### SshUsername (Nullable)

Use the Nullable wrapper methods:
- `obj.SshUsername.IsSet()` — check if set
- `obj.SshUsername.Get()` — get the inner value (returns pointer)
- `obj.SshUsername.Set(&val)` — set the value
- `obj.SshUsername.Unset()` — clear the value
### SshPassword (Nullable)

Use the Nullable wrapper methods:
- `obj.SshPassword.IsSet()` — check if set
- `obj.SshPassword.Get()` — get the inner value (returns pointer)
- `obj.SshPassword.Set(&val)` — set the value
- `obj.SshPassword.Unset()` — clear the value
### SshPasswordHash (Nullable)

Use the Nullable wrapper methods:
- `obj.SshPasswordHash.IsSet()` — check if set
- `obj.SshPasswordHash.Get()` — get the inner value (returns pointer)
- `obj.SshPasswordHash.Set(&val)` — set the value
- `obj.SshPasswordHash.Unset()` — clear the value
### SshKey (Nullable)

Use the Nullable wrapper methods:
- `obj.SshKey.IsSet()` — check if set
- `obj.SshKey.Get()` — get the inner value (returns pointer)
- `obj.SshKey.Set(&val)` — set the value
- `obj.SshKey.Unset()` — clear the value
### MinRam (Nullable)

Use the Nullable wrapper methods:
- `obj.MinRam.IsSet()` — check if set
- `obj.MinRam.Get()` — get the inner value (returns pointer)
- `obj.MinRam.Set(&val)` — set the value
- `obj.MinRam.Unset()` — clear the value
### MinRamGB (Nullable)

Use the Nullable wrapper methods:
- `obj.MinRamGB.IsSet()` — check if set
- `obj.MinRamGB.Get()` — get the inner value (returns pointer)
- `obj.MinRamGB.Set(&val)` — set the value
- `obj.MinRamGB.Unset()` — clear the value
### MinDisk (Nullable)

Use the Nullable wrapper methods:
- `obj.MinDisk.IsSet()` — check if set
- `obj.MinDisk.Get()` — get the inner value (returns pointer)
- `obj.MinDisk.Set(&val)` — set the value
- `obj.MinDisk.Unset()` — clear the value
### MinDiskGB (Nullable)

Use the Nullable wrapper methods:
- `obj.MinDiskGB.IsSet()` — check if set
- `obj.MinDiskGB.Get()` — get the inner value (returns pointer)
- `obj.MinDiskGB.Set(&val)` — set the value
- `obj.MinDiskGB.Unset()` — clear the value
### RawSize (Nullable)

Use the Nullable wrapper methods:
- `obj.RawSize.IsSet()` — check if set
- `obj.RawSize.Get()` — get the inner value (returns pointer)
- `obj.RawSize.Set(&val)` — set the value
- `obj.RawSize.Unset()` — clear the value
### RawSizeGB (Nullable)

Use the Nullable wrapper methods:
- `obj.RawSizeGB.IsSet()` — check if set
- `obj.RawSizeGB.Get()` — get the inner value (returns pointer)
- `obj.RawSizeGB.Set(&val)` — set the value
- `obj.RawSizeGB.Unset()` — clear the value
### Uefi (Nullable)

Use the Nullable wrapper methods:
- `obj.Uefi.IsSet()` — check if set
- `obj.Uefi.Get()` — get the inner value (returns pointer)
- `obj.Uefi.Set(&val)` — set the value
- `obj.Uefi.Unset()` — clear the value
### UserData (Nullable)

Use the Nullable wrapper methods:
- `obj.UserData.IsSet()` — check if set
- `obj.UserData.Get()` — get the inner value (returns pointer)
- `obj.UserData.Set(&val)` — set the value
- `obj.UserData.Unset()` — clear the value
### ConsoleKeymap (Nullable)

Use the Nullable wrapper methods:
- `obj.ConsoleKeymap.IsSet()` — check if set
- `obj.ConsoleKeymap.Get()` — get the inner value (returns pointer)
- `obj.ConsoleKeymap.Set(&val)` — set the value
- `obj.ConsoleKeymap.Unset()` — clear the value
### ExternalId (Nullable)

Use the Nullable wrapper methods:
- `obj.ExternalId.IsSet()` — check if set
- `obj.ExternalId.Get()` — get the inner value (returns pointer)
- `obj.ExternalId.Set(&val)` — set the value
- `obj.ExternalId.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


