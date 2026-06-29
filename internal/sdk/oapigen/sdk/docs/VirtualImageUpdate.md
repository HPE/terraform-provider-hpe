# VirtualImageUpdate

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | A name for the virtual image | [optional] 
**Labels** | Pointer to **[]string** | Array of label strings, can be used for filtering. | [optional] 
**ImageType** | Pointer to **string** | Code of image type. eg. vmware, ami, etc. | [optional] 
**StorageProvider** | Pointer to [**VirtualImageUpdateStorageProvider**](VirtualImageUpdateStorageProvider.md) |  | [optional] 
**IsCloudInit** | Pointer to **bool** | Cloud Init Enabled? | [optional] [default to false]
**UserData** | Pointer to **NullableString** | Cloud-Init User Data, a bash script | [optional] 
**Uefi** | Pointer to **bool** | UEFI enabled? | [optional] 
**FipsEnabled** | Pointer to **bool** | FIPS enabled? | [optional] 
**InstallAgent** | Pointer to **bool** | Install Agent? | [optional] [default to false]
**SshUsername** | Pointer to **NullableString** | SSH Username | [optional] 
**SshPassword** | Pointer to **NullableString** | SSH Password | [optional] 
**SshKey** | Pointer to **NullableString** | SSH Key | [optional] 
**OsType** | Pointer to **NullableInt64** | A Map containing the id of the OS Type. This can also be passed as a string (code or name) instead. | [optional] 
**Visibility** | Pointer to **string** | private or public | [optional] [default to "private"]
**Accounts** | Pointer to **[]int64** |  | [optional] 
**IsAutoJoinDomain** | Pointer to **bool** | Auto Join Domain? | [optional] [default to false]
**VirtioSupported** | Pointer to **bool** | VirtIO Drivers Loaded? | [optional] [default to true]
**VmToolsInstalled** | Pointer to **bool** | VM Tools Installed? | [optional] [default to true]
**IsForceCustomization** | Pointer to **bool** | Force Guest Customization? | [optional] [default to false]
**TrialVersion** | Pointer to **bool** | Trial Version | [optional] [default to false]
**IsSysprep** | Pointer to **bool** | Sysprep Enabled? | [optional] [default to false]
**Config** | Pointer to [**VirtualImageUpdateConfig**](VirtualImageUpdateConfig.md) |  | [optional] 
**Tags** | Pointer to [**[]VirtualImageUpdateTagsInner**](VirtualImageUpdateTagsInner.md) | Metadata tags, Array of objects having a name and value, this adds or updates the specified tags and removes any tags not specified. | [optional] 
**AddTags** | Pointer to [**[]VirtualImageUpdateAddTagsInner**](VirtualImageUpdateAddTagsInner.md) | Add or update value of Metadata tags, Array of objects having a name and value. | [optional] 
**RemoveTags** | Pointer to [**[]VirtualImageUpdateRemoveTagsInner**](VirtualImageUpdateRemoveTagsInner.md) | Remove Metadata tags, Array of objects having a name and an optional value. If value is passed, it must match to be removed. | [optional] 
**MinRamGB** | Pointer to **NullableInt64** |  | [optional] 
**MinDisk** | Pointer to **NullableInt64** |  | [optional] 
**MinDiskGB** | Pointer to **NullableInt64** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &VirtualImageUpdate{
    // Set fields directly
}
```

### Labels (Nullable)

Use the Nullable wrapper methods:
- `obj.Labels.IsSet()` — check if set
- `obj.Labels.Get()` — get the inner value (returns pointer)
- `obj.Labels.Set(&val)` — set the value
- `obj.Labels.Unset()` — clear the value
### UserData (Nullable)

Use the Nullable wrapper methods:
- `obj.UserData.IsSet()` — check if set
- `obj.UserData.Get()` — get the inner value (returns pointer)
- `obj.UserData.Set(&val)` — set the value
- `obj.UserData.Unset()` — clear the value
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
### SshKey (Nullable)

Use the Nullable wrapper methods:
- `obj.SshKey.IsSet()` — check if set
- `obj.SshKey.Get()` — get the inner value (returns pointer)
- `obj.SshKey.Set(&val)` — set the value
- `obj.SshKey.Unset()` — clear the value
### OsType (Nullable)

Use the Nullable wrapper methods:
- `obj.OsType.IsSet()` — check if set
- `obj.OsType.Get()` — get the inner value (returns pointer)
- `obj.OsType.Set(&val)` — set the value
- `obj.OsType.Unset()` — clear the value
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

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


