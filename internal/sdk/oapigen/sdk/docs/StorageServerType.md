# StorageServerType

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Code** | Pointer to **string** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **string** |  | [optional] 
**Enabled** | Pointer to **bool** |  | [optional] 
**Creatable** | Pointer to **bool** |  | [optional] 
**HasNamespaces** | Pointer to **bool** |  | [optional] 
**HasGroups** | Pointer to **bool** |  | [optional] 
**HasBlock** | Pointer to **bool** |  | [optional] 
**HasObject** | Pointer to **bool** |  | [optional] 
**HasFile** | Pointer to **bool** |  | [optional] 
**HasDatastore** | Pointer to **bool** |  | [optional] 
**HasDisks** | Pointer to **bool** |  | [optional] 
**HasHosts** | Pointer to **bool** |  | [optional] 
**CreateNamespaces** | Pointer to **bool** |  | [optional] 
**CreateGroup** | Pointer to **bool** |  | [optional] 
**CreateBlock** | Pointer to **bool** |  | [optional] 
**CreateObject** | Pointer to **bool** |  | [optional] 
**CreateFile** | Pointer to **bool** |  | [optional] 
**CreateDatastore** | Pointer to **bool** |  | [optional] 
**CreateDisk** | Pointer to **bool** |  | [optional] 
**CreateHost** | Pointer to **bool** |  | [optional] 
**IconCode** | Pointer to **NullableString** |  | [optional] 
**HasFileBrowser** | Pointer to **bool** |  | [optional] 
**OptionTypes** | Pointer to [**[]GetStorageServerTypes200ResponseStorageServerTypeOptionTypesInner**](GetStorageServerTypes200ResponseStorageServerTypeOptionTypesInner.md) |  | [optional] 
**GroupOptionTypes** | Pointer to [**[]GetStorageServerTypes200ResponseStorageServerTypeGroupOptionTypesInner**](GetStorageServerTypes200ResponseStorageServerTypeGroupOptionTypesInner.md) |  | [optional] 
**BucketOptionTypes** | Pointer to **[]map[string]interface{}** |  | [optional] 
**ShareOptionTypes** | Pointer to **[]map[string]interface{}** |  | [optional] 
**ShareAccessOptionTypes** | Pointer to **[]map[string]interface{}** |  | [optional] 
**StorageVolumeTypes** | Pointer to [**[]GetStorageServerTypes200ResponseStorageServerTypeStorageVolumeTypesInner**](GetStorageServerTypes200ResponseStorageServerTypeStorageVolumeTypesInner.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &StorageServerType{
    // Set fields directly
}
```

### IconCode (Nullable)

Use the Nullable wrapper methods:
- `obj.IconCode.IsSet()` — check if set
- `obj.IconCode.Get()` — get the inner value (returns pointer)
- `obj.IconCode.Set(&val)` — set the value
- `obj.IconCode.Unset()` — clear the value
### BucketOptionTypes (Nullable)

Use the Nullable wrapper methods:
- `obj.BucketOptionTypes.IsSet()` — check if set
- `obj.BucketOptionTypes.Get()` — get the inner value (returns pointer)
- `obj.BucketOptionTypes.Set(&val)` — set the value
- `obj.BucketOptionTypes.Unset()` — clear the value
### ShareOptionTypes (Nullable)

Use the Nullable wrapper methods:
- `obj.ShareOptionTypes.IsSet()` — check if set
- `obj.ShareOptionTypes.Get()` — get the inner value (returns pointer)
- `obj.ShareOptionTypes.Set(&val)` — set the value
- `obj.ShareOptionTypes.Unset()` — clear the value
### ShareAccessOptionTypes (Nullable)

Use the Nullable wrapper methods:
- `obj.ShareAccessOptionTypes.IsSet()` — check if set
- `obj.ShareAccessOptionTypes.Get()` — get the inner value (returns pointer)
- `obj.ShareAccessOptionTypes.Set(&val)` — set the value
- `obj.ShareAccessOptionTypes.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


